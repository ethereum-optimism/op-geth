package txpool

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
)

// IngressFilter is an interface that allows filtering of transactions before they are added to the transaction pool.
// Implementations of this interface can be used to filter transactions based on various criteria.
// FilterTx will return true if the transaction should be allowed, and false if it should be rejected.
type IngressFilter interface {
	FilterTx(ctx context.Context, tx *types.Transaction) bool
}

type interopFilterAPI interface {
	CurrentInteropBlockTime() (uint64, error)
	TxToInteropAccessList(tx *types.Transaction) []common.Hash
	CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, executingTimestamp uint64) error
}
type interopAccessFilter struct {
	api interopFilterAPI
}

func NewInteropFilter(api interopFilterAPI) IngressFilter {
	return &interopAccessFilter{
		api: api,
	}
}

// FilterTx implements IngressFilter.FilterTx
// it uses provided functions to get the access list from the transaction
// and check it against the supervisor
func (f *interopAccessFilter) FilterTx(ctx context.Context, tx *types.Transaction) bool {
	// if CurrentInteropBlockTime returns an error, we assume that the chain is not set up for interop
	// in which case we allow all transactions
	time, err := f.api.CurrentInteropBlockTime()
	if err != nil {
		return true
	}
	hashes := f.api.TxToInteropAccessList(tx)
	if len(hashes) == 0 {
		return true
	}
	return f.api.CheckAccessList(ctx, hashes, interoptypes.CrossUnsafe, time) == nil
}
