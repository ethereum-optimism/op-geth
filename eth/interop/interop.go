package interop

import (
	"context"

	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/rpc"
)

type InteropClient struct {
	rpcClient *rpc.Client
}

func DialClient(ctx context.Context, rpcEndpoint string) (*InteropClient, error) {
	cl, err := rpc.DialContext(ctx, rpcEndpoint)
	if err != nil {
		return nil, err
	}
	return &InteropClient{rpcClient: cl}, nil
}

func (cl *InteropClient) CheckMessages(ctx context.Context, messages []interoptypes.Message, minSafety interoptypes.SafetyLevel) error {
	return cl.rpcClient.CallContext(ctx, nil, "supervisor_checkMessages", messages, minSafety)
}
