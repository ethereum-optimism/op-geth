package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"syscall"
	"testing"
	"time"
)

// mockRPCServer creates a test HTTP server that responds to JSON-RPC requests
func mockRPCServer(t *testing.T, handler func(method string, params []json.RawMessage) (interface{}, error)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			JSONRPC string            `json:"jsonrpc"`
			ID      json.RawMessage   `json:"id"`
			Method  string            `json:"method"`
			Params  []json.RawMessage `json:"params"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := handler(req.Method, req.Params)

		var resp struct {
			JSONRPC string           `json:"jsonrpc"`
			ID      json.RawMessage  `json:"id"`
			Result  interface{}      `json:"result,omitempty"`
			Error   *json.RawMessage `json:"error,omitempty"`
		}
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		if err != nil {
			errObj := map[string]interface{}{
				"code":    -32000,
				"message": err.Error(),
			}
			errBytes, _ := json.Marshal(errObj)
			errRaw := json.RawMessage(errBytes)
			resp.Error = &errRaw
		} else {
			resp.Result = result
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestNewSequencerClient(t *testing.T) {
	t.Run("empty endpoints", func(t *testing.T) {
		_, err := NewSequencerClient(nil, time.Second, time.Second)
		if err == nil {
			t.Fatal("expected error for empty endpoints")
		}
	})

	t.Run("empty endpoint string", func(t *testing.T) {
		_, err := NewSequencerClient([]string{""}, time.Second, time.Second)
		if err == nil {
			t.Fatal("expected error for empty endpoint string")
		}
	})

	t.Run("zero dial timeout", func(t *testing.T) {
		_, err := NewSequencerClient([]string{"http://localhost:8545"}, 0, time.Second)
		if err == nil {
			t.Fatal("expected error for zero dial timeout")
		}
	})

	t.Run("zero request timeout", func(t *testing.T) {
		_, err := NewSequencerClient([]string{"http://localhost:8545"}, time.Second, 0)
		if err == nil {
			t.Fatal("expected error for zero request timeout")
		}
	})

	t.Run("single endpoint", func(t *testing.T) {
		client, err := NewSequencerClient([]string{"http://localhost:8545"}, time.Second, time.Second)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.EndpointCount() != 1 {
			t.Fatalf("expected 1 endpoint, got %d", client.EndpointCount())
		}
		client.Close()
	})

	t.Run("multiple endpoints", func(t *testing.T) {
		client, err := NewSequencerClient([]string{"http://localhost:8545", "http://localhost:8546"}, time.Second, time.Second)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.EndpointCount() != 2 {
			t.Fatalf("expected 2 endpoints, got %d", client.EndpointCount())
		}
		client.Close()
	})
}

func TestSequencerClient_CallContext_SingleEndpoint_Success(t *testing.T) {
	server := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		if method == "eth_sendRawTransaction" {
			return "0x1234567890abcdef", nil
		}
		return nil, fmt.Errorf("unknown method: %s", method)
	})
	defer server.Close()

	client, err := NewSequencerClient([]string{server.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "0x1234567890abcdef" {
		t.Fatalf("unexpected result: %s", result)
	}
}

func TestSequencerClient_CallContext_SingleEndpoint_ApplicationError(t *testing.T) {
	server := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		return nil, errors.New("nonce too low")
	})
	defer server.Close()

	client, err := NewSequencerClient([]string{server.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected error")
	}
	if err.Error() != "nonce too low" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSequencerClient_CallContext_Failover_FirstFails_SecondSucceeds(t *testing.T) {
	var callCount atomic.Int32

	// First server always fails with connection-like error
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount.Add(1)
		// Simulate server error (5xx triggers failover)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server1.Close()

	// Second server succeeds
	server2 := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		callCount.Add(1)
		return "0xsuccess", nil
	})
	defer server2.Close()

	client, err := NewSequencerClient([]string{server1.URL, server2.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "0xsuccess" {
		t.Fatalf("unexpected result: %s", result)
	}

	// Verify failover occurred (both servers were called)
	if callCount.Load() != 2 {
		t.Fatalf("expected 2 calls (failover), got %d", callCount.Load())
	}

	// Verify preferred endpoint changed to second server
	if client.PreferredEndpoint() != server2.URL {
		t.Fatalf("expected preferred endpoint to be %s, got %s", server2.URL, client.PreferredEndpoint())
	}
}

func TestSequencerClient_CallContext_AllEndpointsFail(t *testing.T) {
	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server2.Close()

	client, err := NewSequencerClient([]string{server1.URL, server2.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected error when all endpoints fail")
	}
	if !strings.Contains(err.Error(), "all 2 sequencer endpoints failed") {
		t.Fatalf("unexpected error message: %v", err)
	}
}

func TestSequencerClient_CallContext_ApplicationError_NoFailover(t *testing.T) {
	var callCount atomic.Int32

	// First server returns application error (should NOT trigger failover)
	server1 := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		callCount.Add(1)
		return nil, errors.New("nonce too low")
	})
	defer server1.Close()

	// Second server should never be called
	server2 := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		callCount.Add(1)
		return "0xsuccess", nil
	})
	defer server2.Close()

	client, err := NewSequencerClient([]string{server1.URL, server2.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected error")
	}

	// Only first server should be called (no failover for application errors)
	if callCount.Load() != 1 {
		t.Fatalf("expected 1 call (no failover for app error), got %d", callCount.Load())
	}
}

func TestSequencerClient_CallContext_StickyPreference(t *testing.T) {
	var server1Calls, server2Calls atomic.Int32

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server1Calls.Add(1)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}))
	defer server1.Close()

	server2 := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		server2Calls.Add(1)
		return "0xsuccess", nil
	})
	defer server2.Close()

	client, err := NewSequencerClient([]string{server1.URL, server2.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// First call - should failover to server2
	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0x1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Second call - should go directly to server2 (sticky preference)
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0x2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Third call - should go directly to server2 (sticky preference)
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0x3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Server1 should only be called once (first call failover)
	if server1Calls.Load() != 1 {
		t.Fatalf("expected server1 to be called 1 time, got %d", server1Calls.Load())
	}

	// Server2 should be called 3 times (all successful calls)
	if server2Calls.Load() != 3 {
		t.Fatalf("expected server2 to be called 3 times, got %d", server2Calls.Load())
	}
}

func TestSequencerClient_CallContext_ClosedClient(t *testing.T) {
	server := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		return "0xsuccess", nil
	})
	defer server.Close()

	client, err := NewSequencerClient([]string{server.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}

	client.Close()

	var result string
	err = client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected error for closed client")
	}
	if !strings.Contains(err.Error(), "closed") {
		t.Fatalf("expected closed error, got: %v", err)
	}
}

func TestSequencerClient_CallContext_ContextCanceled(t *testing.T) {
	server := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		time.Sleep(5 * time.Second) // Slow response
		return "0xsuccess", nil
	})
	defer server.Close()

	client, err := NewSequencerClient([]string{server.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result string
	err = client.CallContext(ctx, &result, "eth_sendRawTransaction", "0xdeadbeef")
	if err == nil {
		t.Fatal("expected error for canceled context")
	}
}

func TestSequencerClient_CallContext_Concurrent(t *testing.T) {
	var callCount atomic.Int32

	server := mockRPCServer(t, func(method string, params []json.RawMessage) (interface{}, error) {
		callCount.Add(1)
		return "0xsuccess", nil
	})
	defer server.Close()

	client, err := NewSequencerClient([]string{server.URL}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Run 10 concurrent calls
	done := make(chan error, 10)
	for range 10 {
		go func() {
			var result string
			err := client.CallContext(context.Background(), &result, "eth_sendRawTransaction", "0xdeadbeef")
			done <- err
		}()
	}

	// Wait for all to complete
	for range 10 {
		if err := <-done; err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if callCount.Load() != 10 {
		t.Fatalf("expected 10 calls, got %d", callCount.Load())
	}
}

func TestShouldFailover(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"nil error", nil, false},
		{"context deadline exceeded", context.DeadlineExceeded, true},
		{"context canceled", context.Canceled, false},
		{"rpc client quit", ErrClientQuit, true},
		{"io.EOF", io.EOF, true},
		{"io.ErrUnexpectedEOF", io.ErrUnexpectedEOF, true},
		{"wrapped EOF", fmt.Errorf("read failed: %w", io.EOF), true},
		{"connection refused", &net.OpError{Op: "dial", Err: &net.DNSError{Err: "connection refused"}}, true},
		{"dns error", &net.DNSError{Err: "no such host", Name: "example.com"}, true},
		{"syscall ECONNREFUSED", syscall.ECONNREFUSED, true},
		{"syscall ECONNRESET", syscall.ECONNRESET, true},
		{"syscall EPIPE", syscall.EPIPE, true},
		{"syscall ETIMEDOUT", syscall.ETIMEDOUT, true},
		{"wrapped net error", fmt.Errorf("request failed: %w", &net.OpError{Op: "read", Err: syscall.ECONNRESET}), true},
		{"http 500", HTTPError{StatusCode: 500, Status: "500 Internal Server Error"}, true},
		{"http 503", HTTPError{StatusCode: 503, Status: "503 Service Unavailable"}, true},
		{"http 400", HTTPError{StatusCode: 400, Status: "400 Bad Request"}, false},
		{"application error", errors.New("nonce too low"), false},
		{"insufficient funds", errors.New("insufficient funds"), false},
		{"gas too low", errors.New("intrinsic gas too low"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldFailover(tt.err)
			if result != tt.expected {
				t.Errorf("shouldFailover(%v) = %v, expected %v", tt.err, result, tt.expected)
			}
		})
	}
}

func TestSequencerClient_PreferredEndpoint(t *testing.T) {
	client, err := NewSequencerClient([]string{"http://endpoint1:8545", "http://endpoint2:8545"}, time.Second, time.Second)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Initially should be first endpoint
	if client.PreferredEndpoint() != "http://endpoint1:8545" {
		t.Fatalf("expected first endpoint as preferred, got %s", client.PreferredEndpoint())
	}
}
