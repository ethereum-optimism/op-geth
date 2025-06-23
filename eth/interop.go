package eth

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/miner"
)

func (s *Ethereum) CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, execDesc interoptypes.ExecutingDescriptor) error {
	if s.interopRPC == nil {
		return errors.New("cannot check interop access list, no RPC available")
	}
	return s.interopRPC.CheckAccessList(ctx, inboxEntries, minSafety, execDesc)
}

func (s *Ethereum) inferBlockTime(current *types.Header) (uint64, error) {
	if current.Number.Uint64() == 0 {
		return 0, errors.New("current head is at genesis: penultimate header is nil")
	}
	penultimate := s.BlockChain().GetHeaderByHash(current.ParentHash)
	if penultimate == nil {
		// We could use a for loop and retry, but this function is used
		// in the ingress filters, which should fail fast to maintain uptime.
		return 0, errors.New("penultimate header is nil")
	}
	return current.Time - penultimate.Time, nil
}

// CurrentInteropBlockTime returns the current block time,
// or an error if Interop is not enabled.
func (s *Ethereum) CurrentInteropBlockTime() (uint64, error) {
	chainConfig := s.APIBackend.ChainConfig()
	if !chainConfig.IsOptimism() {
		return 0, errors.New("chain is not an Optimism chain")
	}
	if chainConfig.InteropTime == nil {
		return 0, errors.New("interop time not set in chain config")
	}
	// The pending block may be aliased to the current block in op-geth. Infer the pending time instead.
	currentHeader := s.BlockChain().CurrentHeader()
	blockTime, err := s.inferBlockTime(currentHeader)
	if err != nil {
		return 0, fmt.Errorf("infer block time: %v", err)
	}
	return currentHeader.Time + blockTime, nil
}

// TxToInteropAccessList returns the interop specific access list storage keys for a transaction.
func (s *Ethereum) TxToInteropAccessList(tx *types.Transaction) []common.Hash {
	return interoptypes.TxToInteropAccessList(tx)
}

var _ miner.BackendWithInterop = (*Ethereum)(nil)
