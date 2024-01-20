package beacon

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type OpLegacy struct{}

func (o *OpLegacy) Author(header *types.Header) (common.Address, error) {
	return header.Coinbase, nil
}

func (o *OpLegacy) VerifyHeader(chain consensus.ChainHeaderReader, header *types.Header) error {
	// redundant check to guarantee DB consistency
	parent := chain.GetHeader(header.ParentHash, header.Number.Uint64()-1)
	if parent == nil {
		return consensus.ErrUnknownAncestor
	}
	return nil // legacy chain is verified by block-hash reverse sync otherwise
}

func (o *OpLegacy) VerifyHeaders(chain consensus.ChainHeaderReader, headers []*types.Header) (chan<- struct{}, <-chan error) {
	quit := make(chan struct{}, 1)
	result := make(chan error, len(headers))
	for _, h := range headers {
		result <- o.VerifyHeader(chain, h)
	}
	return quit, result
}

func (o *OpLegacy) VerifyUncles(chain consensus.ChainReader, block *types.Block) error {
	return nil
}

func (o *OpLegacy) Prepare(chain consensus.ChainHeaderReader, header *types.Header) error {
	return fmt.Errorf("cannot prepare for legacy block header: %s (num %d)", header.Hash(), header.Number)
}

func (o *OpLegacy) Finalize(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, withdrawals []*types.Withdrawal) {
	panic(fmt.Errorf("cannot finalize legacy block header: %s (num %d)", header.Hash(), header.Number))
}

func (o *OpLegacy) FinalizeAndAssemble(chain consensus.ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction, uncles []*types.Header, receipts []*types.Receipt, withdrawals []*types.Withdrawal) (*types.Block, error) {
	return nil, fmt.Errorf("cannot finalize and assemble for legacy block header: %s (num %d)", header.Hash(), header.Number)
}

func (o *OpLegacy) Seal(chain consensus.ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return fmt.Errorf("cannot seal legacy block header: %s (num %d)", block.Hash(), block.Number())
}

func (o *OpLegacy) SealHash(header *types.Header) common.Hash {
	panic(fmt.Errorf("cannot compute pow/poa seal-hash for legacy block header: %s (num %d)", header.Hash(), header.Number))
}

func (o *OpLegacy) CalcDifficulty(chain consensus.ChainHeaderReader, time uint64, parent *types.Header) *big.Int {
	return big.NewInt(0)
}

func (o *OpLegacy) APIs(chain consensus.ChainHeaderReader) []rpc.API {
	return nil
}

func (o *OpLegacy) Close() error {
	return nil
}

var _ consensus.Engine = (*OpLegacy)(nil)
