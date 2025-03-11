package txpool

import (
	"context"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/log"
)

// IngressFilter is an interface that allows filtering of transactions before they are added to the transaction pool.
// Implementations of this interface can be used to filter transactions based on various criteria.
// FilterTx will return true if the transaction should be allowed, and false if it should be rejected.
type IngressFilter interface {
	FilterTx(ctx context.Context, tx *types.Transaction) bool
}

type interopSimFilter struct {
	logsFn  func(tx *types.Transaction) (logs []*types.Log, logTimestamp uint64, err error)
	checkFn func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error
}

func NewInteropFilter(
	logsFn func(tx *types.Transaction) ([]*types.Log, uint64, error),
	checkFn func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error) IngressFilter {
	return &interopSimFilter{
		logsFn:  logsFn,
		checkFn: checkFn,
	}
}

// FilterTx implements IngressFilter.FilterTx
// it gets logs checks for message safety based on the function provided
func (f *interopSimFilter) FilterTx(ctx context.Context, tx *types.Transaction) bool {
	logs, logTimestamp, err := f.logsFn(tx)
	if err != nil {
		log.Debug("Failed to retrieve logs of tx", "txHash", tx.Hash(), "err", err)
		return false // default to deny if logs cannot be retrieved
	}
	if len(logs) == 0 {
		return true // default to allow if there are no logs
	}
	ems, err := interoptypes.ExecutingMessagesFromLogs(logs)
	if err != nil {
		log.Debug("Failed to parse executing messages of tx", "txHash", tx.Hash(), "err", err)
		return false // default to deny if logs cannot be parsed
	}
	if len(ems) == 0 {
		return true // default to allow if there are no executing messages
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	// check with the supervisor if the transaction should be allowed given the executing messages
	// the message can be unsafe (discovered only via P2P unsafe blocks), but it must be cross-valid
	// so CrossUnsafe is used here
	return f.checkFn(ctx, ems, interoptypes.CrossUnsafe, logTimestamp) == nil
}

type interopAccessFilter struct {
}

// FilterTx implements IngressFilter.FilterTx
// it takes the access list from the transaction and checks it against the supervisor
func (f *interopAccessFilter) FilterTx(ctx context.Context, tx *types.Transaction) bool {
	return true
}
