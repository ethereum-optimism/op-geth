package rpc

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

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

// getClient returns an RPC client for the given index
func (sc *SequencerClient) getClient(ctx context.Context, idx int) (*Client, error) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if sc.closed {
		return nil, errors.New("sequencer client is closed")
	}

	if sc.clients[idx] != nil {
		return sc.clients[idx], nil
	}

	// Dial with timeout
	dialCtx, cancel := context.WithTimeout(ctx, sc.dialTimeout)
	defer cancel()

	client, err := DialContext(dialCtx, sc.endpoints[idx])
	if err != nil {
		return nil, fmt.Errorf("failed to dial sequencer endpoint %s: %w", sc.endpoints[idx], err)
	}

	sc.clients[idx] = client
	log.Info("Connected to sequencer endpoint", "endpoint", sc.endpoints[idx], "index", idx)
	return client, nil
}

// shouldFailover determines if the error warrants trying another endpoint.
// Returns true for connection/transport errors, false for application-level errors.
func shouldFailover(err error) bool {
	if err == nil {
		return false
	}

	// Context errors - these indicate connection issues
	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return true
	}

	// RPC client closed/quit
	if errors.Is(err, ErrClientQuit) {
		return true
	}

	// HTTP errors - only failover on 5xx (server errors)
	var httpErr HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode >= http.StatusInternalServerError
	}

	// Check for connection/network errors by message content
	errStr := strings.ToLower(err.Error())
	connectionErrors := []string{
		"connection refused",
		"connection reset",
		"no such host",
		"network is unreachable",
		"i/o timeout",
		"eof",
		"broken pipe",
		"connection timed out",
	}
	for _, ce := range connectionErrors {
		if strings.Contains(errStr, ce) {
			return true
		}
	}

	return false
}

// CallContext performs an RPC call with automatic failover on connection errors.
// It tries the preferred endpoint first, then fails over to other endpoints if needed.
// Application-level errors (like nonce too low) are returned immediately without failover.
func (sc *SequencerClient) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	sequencerCallsTotal.Inc(1)

	sc.mu.RLock()
	if sc.closed {
		sc.mu.RUnlock()
		return errors.New("sequencer client is closed")
	}
	startIdx := sc.preferredIndex
	numEndpoints := len(sc.endpoints)
	sc.mu.RUnlock()

	var lastErr error

	for i := range numEndpoints {
		idx := (startIdx + i) % numEndpoints

		// Check if parent context is already done
		if ctx.Err() != nil {
			return ctx.Err()
		}

		client, err := sc.getClient(ctx, idx)
		if err != nil {
			lastErr = err
			log.Warn("Failed to get sequencer client", "endpoint", sc.endpoints[idx], "err", err)
			if i < numEndpoints-1 {
				sequencerFailovers.Inc(1)
				log.Info("Failing over to next sequencer endpoint",
					"failed", sc.endpoints[idx],
					"next", sc.endpoints[(idx+1)%numEndpoints])
			}
			continue
		}

		// Create timeout context for this specific request
		reqCtx, cancel := context.WithTimeout(ctx, sc.requestTimeout)
		err = client.CallContext(reqCtx, result, method, args...)
		cancel()

		if err == nil {
			sequencerCallsSuccess.Inc(1)
			// Update preferred endpoint if different
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

		// Check if we should failover
		if !shouldFailover(err) {
			// Application-level error, don't failover
			sequencerCallsFailed.Inc(1)
			return err
		}

		// Mark this client as potentially bad (close and nil it for re-dial)
		sc.mu.Lock()
		if sc.clients[idx] != nil {
			sc.clients[idx].Close()
			sc.clients[idx] = nil
		}
		sc.mu.Unlock()

		log.Warn("Sequencer endpoint failed",
			"endpoint", sc.endpoints[idx],
			"err", err,
			"remaining", numEndpoints-i-1)

		if i < numEndpoints-1 {
			sequencerFailovers.Inc(1)
			log.Info("Failing over to next sequencer endpoint",
				"failed", sc.endpoints[idx],
				"next", sc.endpoints[(idx+1)%numEndpoints])
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
