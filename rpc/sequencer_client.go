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
	mu             sync.Mutex
	endpoints      []string
	client         *Client
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

	if dialTimeout == 0 || requestTimeout == 0 {
		return nil, errors.New("dialTimeout and requestTimeout must be greater than 0")
	}

	for i, ep := range endpoints {
		if ep == "" {
			return nil, fmt.Errorf("sequencer endpoint %d is empty", i)
		}
	}

	return &SequencerClient{
		endpoints:      endpoints,
		client:         nil,
		preferredIndex: 0,
		dialTimeout:    dialTimeout,
		requestTimeout: requestTimeout,
	}, nil
}

// getClient returns the current RPC client or dials a new one.
// If the current client is nil or connection fails, try the next endpoint.
func (sc *SequencerClient) getClient(ctx context.Context) (*Client, error) {
	if sc.closed {
		return nil, errSequencerClosed
	}

	if sc.client != nil {
		return sc.client, nil
	}

	// Try to connect to an endpoint, starting from preferredIndex
	startIdx := sc.preferredIndex
	numEndpoints := len(sc.endpoints)

	for i := range numEndpoints {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}

		idx := (startIdx + i) % numEndpoints
		endpoint := sc.endpoints[idx]

		dialCtx, cancel := context.WithTimeout(ctx, sc.dialTimeout)
		defer cancel()
		client, err := DialContext(dialCtx, endpoint)

		if err != nil {
			log.Warn("Failed to connect to sequencer endpoint", "endpoint", endpoint, "err", err)
			continue
		}

		sc.client = client
		sc.preferredIndex = idx
		log.Info("Connected to sequencer endpoint", "endpoint", endpoint, "index", idx)
		return client, nil
	}

	return nil, fmt.Errorf("all %d sequencer endpoints are unreachable", numEndpoints)
}

// invalidateCurrentClient closes the current client and moves to the next endpoint.
func (sc *SequencerClient) invalidateCurrentClient() {
	if sc.client != nil {
		sc.client.Close()
		sc.client = nil
	}

	// Move to next endpoint for the next connection attempt
	sc.preferredIndex = (sc.preferredIndex + 1) % len(sc.endpoints)
	log.Info("Invalidated sequencer client, will try next endpoint",
		"new_index", sc.preferredIndex,
		"endpoint", sc.endpoints[sc.preferredIndex])
	sequencerEndpointGauge.Update(int64(sc.preferredIndex))
	sequencerFailovers.Inc(1)
}

// shouldFailover determines if the error warrants trying another endpoint.
// Returns true for connection/transport errors, false for application-level errors.
func shouldFailover(err error) bool {
	if err == nil {
		return false
	}

	// Context errors — timeout
	if errors.Is(err, context.DeadlineExceeded) {
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
// If the call fails with a network error, it invalidates the current client and retries
// with the next endpoint. Application-level errors are returned immediately without failover.
func (sc *SequencerClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sequencerCallsTotal.Inc(1)

	numEndpoints := len(sc.endpoints)
	var lastErr error

	for attempt := range numEndpoints {
		client, err := sc.getClient(ctx)
		if err != nil {
			lastErr = err
			sequencerCallsFailed.Inc(1)
			return fmt.Errorf("failed to get sequencer client: %w", err)
		}

		reqCtx, cancel := context.WithTimeout(ctx, sc.requestTimeout)
		defer cancel()
		err = client.CallContext(reqCtx, result, method, args...)

		// Success case
		if err == nil {
			sequencerCallsSuccess.Inc(1)
			return nil
		}

		lastErr = err

		// If the caller's context was canceled, don't blame the endpoint
		if ctx.Err() != nil {
			sequencerCallsFailed.Inc(1)
			return ctx.Err()
		}

		if !shouldFailover(err) {
			sequencerCallsFailed.Inc(1)
			return err
		}

		log.Warn("Sequencer endpoint failed with network error",
			"endpoint", sc.endpoints[sc.preferredIndex],
			"err", err,
			"remaining_attempts", numEndpoints-attempt-1)

		sc.invalidateCurrentClient()
	}

	// All endpoints failed
	sequencerCallsFailed.Inc(1)
	return fmt.Errorf("all %d sequencer endpoints failed, last error: %w", numEndpoints, lastErr)
}

// Closes the underlying RPC client if connected.
func (sc *SequencerClient) Close() {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return
	}
	sc.closed = true

	if sc.client != nil {
		sc.client.Close()
		sc.client = nil
	}
	log.Info("Sequencer client closed", "endpoints", len(sc.endpoints))
}

// EndpointCount returns the number of configured endpoints.
func (sc *SequencerClient) EndpointCount() int {
	return len(sc.endpoints)
}

// PreferredEndpoint returns the currently preferred endpoint URL.
func (sc *SequencerClient) PreferredEndpoint() string {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.endpoints[sc.preferredIndex]
}
