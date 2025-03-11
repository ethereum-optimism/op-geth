package eth

import (
	"context"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/miner"
)

func (s *Ethereum) CheckMessages(ctx context.Context, messages []interoptypes.Message, minSafety interoptypes.SafetyLevel, executingTimestamp uint64) error {
	return errors.New("deprecated. use CheckAccessList")
}

func (s *Ethereum) CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, executingTimestamp uint64) error {
	if s.interopRPC == nil {
		return errors.New("cannot check interop access list, no RPC available")
	}
	return s.interopRPC.CheckAccessList(ctx, inboxEntries, minSafety, interoptypes.ExecutingDescriptor{Timestamp: executingTimestamp})
}

var _ miner.BackendWithInterop = (*Ethereum)(nil)
