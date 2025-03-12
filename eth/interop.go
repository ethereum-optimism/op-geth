package eth

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/miner"
	"github.com/ethereum/go-ethereum/params"
)

func (s *Ethereum) CheckMessages(ctx context.Context, messages []interoptypes.Message, minSafety interoptypes.SafetyLevel, executingTimestamp uint64) error {
	return errors.New("deprecated. use CheckAccessList")
}

func (s *Ethereum) CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, execDesc interoptypes.ExecutingDescriptor) error {
	if s.interopRPC == nil {
		return errors.New("cannot check interop access list, no RPC available")
	}
	return s.interopRPC.CheckAccessList(ctx, inboxEntries, minSafety, execDesc)
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
	return s.BlockChain().CurrentBlock().Time, nil
}

// TxToInteropAccessList returns the interop specific access list storage keys for a transaction.
func (s *Ethereum) TxToInteropAccessList(tx *types.Transaction) []common.Hash {
	return txToInteropAccessList(tx)
}

func txToInteropAccessList(tx *types.Transaction) []common.Hash {
	if tx == nil {
		return nil
	}
	al := tx.AccessList()
	if len(al) == 0 {
		return nil
	}
	var hashes []common.Hash
	for i := range al {
		if al[i].Address == params.InteropCrossL2InboxAddress {
			hashes = append(hashes, al[i].StorageKeys...)
		}
	}
	return hashes
}

var _ miner.BackendWithInterop = (*Ethereum)(nil)
