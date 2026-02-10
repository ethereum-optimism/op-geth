package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

var errSequencerClosed = errors.New("sequencer client is closed")

var (
	// Metrics for monitoring sequencer client behavior
	sequencerCallsTotal    = metrics.NewRegisteredCounter("sequencer/calls/total", nil)
	sequencerCallsSuccess  = metrics.NewRegisteredCounter("sequencer/calls/success", nil)
	sequencerCallsFailed   = metrics.NewRegisteredCounter("sequencer/calls/failed", nil)
	sequencerFailovers     = metrics.NewRegisteredCounter("sequencer/failovers", nil)
	sequencerEndpointGauge = metrics.NewRegisteredGauge("sequencer/active_endpoint_index", nil)
)

// SequencerClient manages multiple sequencer RPC endpoints with automatic failover.
// It tries the preferred endpoint first, and on connection failures, automatically
// fails over to the next available endpoint.
type SequencerClient struct {
	mu             sync.RWMutex
	endpoints      []string
	clients        []*Client
	preferredIndex int // index of currently preferred endpoint
	closed         bool
	dialTimeout    time.Duration
	requestTimeout time.Duration
}

// NewSequencerClient creates a new SequencerClient with the given endpoints.
// At least one endpoint must be provided.
func NewSequencerClient(endpoints []string, dialTimeout, requestTimeout time.Duration) (*SequencerClient, error) {
	if len(endpoints) == 0 {
		return nil, errors.New("at least one sequencer endpoint required")
	}

	for i, ep := range endpoints {
		if ep == "" {
			return nil, fmt.Errorf("sequencer endpoint %d is empty", i)
		}
	}

	return &SequencerClient{
		endpoints:      endpoints,
		clients:        make([]*Client, len(endpoints)),
		preferredIndex: 0,
		dialTimeout:    dialTimeout,
		requestTimeout: requestTimeout,
	}, nil
}

// getClient returns an available RPC client, starting from the preferred endpoint.
// If the preferred endpoint's client is not connected, it dials it. If dialing fails,
// it cycles through the remaining endpoints. Returns the client and its endpoint index.
func (sc *SequencerClient) getClient(ctx context.Context) (*Client, int, error) {
	sc.mu.RLock()
	if sc.closed {
		sc.mu.RUnlock()
		return nil, 0, errSequencerClosed
	}
	startIdx := sc.preferredIndex
	numEndpoints := len(sc.endpoints)
	sc.mu.RUnlock()

	for i := range numEndpoints {
		idx := (startIdx + i) % numEndpoints

		if ctx.Err() != nil {
			return nil, 0, ctx.Err()
		}

		client, err := sc.connectEndpoint(ctx, idx)
		if errors.Is(err, errSequencerClosed) {
			return nil, 0, err
		}
		if err != nil {
			log.Warn("Failed to connect to sequencer endpoint", "endpoint", sc.endpoints[idx], "err", err)
			continue
		}
		return client, idx, nil
	}

	return nil, 0, fmt.Errorf("all %d sequencer endpoints are unreachable", numEndpoints)
}

// connectEndpoint returns a cached client for the given index or dials a new one.
func (sc *SequencerClient) connectEndpoint(ctx context.Context, idx int) (*Client, error) {
	sc.mu.RLock()
	if sc.closed {
		sc.mu.RUnlock()
		return nil, errSequencerClosed
	}
	if client := sc.clients[idx]; client != nil {
		sc.mu.RUnlock()
		return client, nil
	}
	sc.mu.RUnlock()

	dialCtx, cancel := context.WithTimeout(ctx, sc.dialTimeout)
	client, err := DialContext(dialCtx, sc.endpoints[idx])
	cancel()
	if err != nil {
		return nil, fmt.Errorf("failed to dial %s: %w", sc.endpoints[idx], err)
	}

	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		client.Close()
		return nil, errSequencerClosed
	}
	if existing := sc.clients[idx]; existing != nil {
		client.Close()
		return existing, nil
	}
	sc.clients[idx] = client
	log.Info("Connected to sequencer endpoint", "endpoint", sc.endpoints[idx], "index", idx)
	return client, nil
}

// invalidateClient closes the client at the given index and advances the preferred
// endpoint so subsequent calls start from the next endpoint.
func (sc *SequencerClient) invalidateClient(idx int) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.clients[idx] != nil {
		sc.clients[idx].Close()
		sc.clients[idx] = nil
	}
	if sc.preferredIndex == idx {
		sc.preferredIndex = (idx + 1) % len(sc.endpoints)
		sequencerEndpointGauge.Update(int64(sc.preferredIndex))
	}
}

// shouldFailover determines if the error warrants trying another endpoint.
// Returns true for connection/transport errors, false for application-level errors.
func shouldFailover(err error) bool {
	if err == nil {
		return false
	}

	// Context errors — timeout or cancellation
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	// RPC client closed/quit
	if errors.Is(err, ErrClientQuit) {
		return true
	}

	// EOF — connection dropped
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}

	// Network-level errors (connection refused, reset, DNS failure, timeout, etc.)
	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return true
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return true
	}

	// Syscall errors (broken pipe, connection reset)
	if errors.Is(err, syscall.ECONNREFUSED) ||
		errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.ECONNABORTED) ||
		errors.Is(err, syscall.EPIPE) ||
		errors.Is(err, syscall.ETIMEDOUT) {
		return true
	}

	// HTTP errors — only failover on 5xx (server errors)
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode >= http.StatusInternalServerError // >= 500
	}

	return false
}

// CallContext performs an RPC call with automatic failover on connection errors.
// It gets a client via getClient, makes the call, and on connection failure
// invalidates the client and retries with the next available endpoint.
// Application-level errors (like nonce too low) are returned immediately without failover.
func (sc *SequencerClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	sequencerCallsTotal.Inc(1)

	numEndpoints := len(sc.endpoints)
	var lastErr error

	for attempt := range numEndpoints {
		client, idx, err := sc.getClient(ctx)
		if err != nil {
			sequencerCallsFailed.Inc(1)
			if lastErr != nil {
				return fmt.Errorf("all %d sequencer endpoints failed, last error: %w", numEndpoints, lastErr)
			}
			return err
		}

		reqCtx, cancel := context.WithTimeout(ctx, sc.requestTimeout)
		err = client.CallContext(reqCtx, result, method, args...)
		cancel()

		if err == nil {
			sequencerCallsSuccess.Inc(1)
			// Update preferred endpoint on success
			sc.mu.Lock()
			if sc.preferredIndex != idx {
				log.Info("Switching preferred sequencer endpoint",
					"from", sc.endpoints[sc.preferredIndex],
					"to", sc.endpoints[idx])
				sc.preferredIndex = idx
				sequencerEndpointGauge.Update(int64(idx))
			}
			sc.mu.Unlock()
			return nil
		}

		lastErr = err

		if !shouldFailover(err) {
			sequencerCallsFailed.Inc(1)
			return err
		}

		log.Warn("Sequencer endpoint failed", "endpoint", sc.endpoints[idx], "err", err, "remaining", numEndpoints-attempt-1)
		sc.invalidateClient(idx)

		if attempt < numEndpoints-1 {
			sequencerFailovers.Inc(1)
		}
	}

	sequencerCallsFailed.Inc(1)
	return fmt.Errorf("all %d sequencer endpoints failed, last error: %w", numEndpoints, lastErr)
}

// Close closes all underlying RPC clients.
func (sc *SequencerClient) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return
	}
	sc.closed = true

	for i, client := range sc.clients {
		if client != nil {
			client.Close()
			sc.clients[i] = nil
		}
	}
	log.Info("Sequencer client closed", "endpoints", len(sc.endpoints))
}

// EndpointCount returns the number of configured endpoints.
func (sc *SequencerClient) EndpointCount() int {
	return len(sc.endpoints)
}

// PreferredEndpoint returns the currently preferred endpoint URL.
func (sc *SequencerClient) PreferredEndpoint() string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return sc.endpoints[sc.preferredIndex]
}
