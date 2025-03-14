package interop

import (
	"context"
	"errors"
	"sync"

	"github.com/ethereum/go-ethereum/common"
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
	if cl.closed {
		return errors.New("client is closed")
	}
	if cl.client != nil {
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
		cl.client.Close()
	}
	cl.closed = true
}

func NewInteropClient(rpcEndpoint string) *InteropClient {
	return &InteropClient{endpoint: rpcEndpoint}
}

func (cl *InteropClient) CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, executingDescriptor interoptypes.ExecutingDescriptor) error {
	if err := cl.maybeDial(ctx); err != nil {
		return err
	}
	err := cl.client.CallContext(ctx, nil, "supervisor_checkAccessList", inboxEntries, minSafety, executingDescriptor)
	return err
}
