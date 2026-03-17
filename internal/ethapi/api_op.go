package ethapi

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

// The HeaderByNumberOrHash method returns a nil error and nil header
// if the header is not found, but only for nonexistent block numbers. This is
// different from StateAndHeaderByNumberOrHash. To account for this discrepancy,
// headerOrNumberByHash will properly convert the error into an ethereum.NotFound.
func headerByNumberOrHash(ctx context.Context, b Backend, blockNrOrHash rpc.BlockNumberOrHash) (*types.Header, error) {
	header, err := b.HeaderByNumberOrHash(ctx, blockNrOrHash)
	if header == nil {
		return nil, fmt.Errorf("header %w", ethereum.NotFound)
	}
	return header, err
}
