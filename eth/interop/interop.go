package interop

import (
	"context"
	"sync"

	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/rpc"
)

type InteropClient struct {
	mu       sync.Mutex
	client   *rpc.Client
	endpoint string
	closed   bool // don't allow lazy-dials after Close
}

// maybeDial dials the endpoint if it was not already.
func (cl *InteropClient) maybeDial(ctx context.Context) error {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if cl.client != nil || cl.closed {
		return nil
	}
	rpcClient, err := rpc.DialContext(ctx, cl.endpoint)
	if err != nil {
		return err
	}
	cl.client = rpcClient
	return nil
}

func (cl *InteropClient) Close() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if cl.client != nil {
		cl.Close()
	}
	cl.closed = true
}

func NewInteropClient(rpcEndpoint string) *InteropClient {
	return &InteropClient{endpoint: rpcEndpoint}
}

// CheckMessages checks if the given messages meet the given minimum safety level.
func (cl *InteropClient) CheckMessages(ctx context.Context, messages []interoptypes.Message, minSafety interoptypes.SafetyLevel) error {
	// we lazy-dial the endpoint, so we can start geth, and build blocks, without supervisor endpoint availability.
	if err := cl.maybeDial(ctx); err != nil { // a single dial attempt is made, the next call may retry.
		return err
	}
	return cl.client.CallContext(ctx, nil, "supervisor_checkMessages", messages, minSafety)
}
