package txpool

import (
	"context"
	"time"

	"github.com/holiman/uint256"

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
	CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, execDesc interoptypes.ExecutingDescriptor) error
}

type interopAccessFilter struct {
	api     interopFilterAPI
	timeout uint64
	chainID uint256.Int
}

// NewInteropFilter creates a new IngressFilter that filters transactions based on the interop access list.
// the timeout is set to 1 day, the specified preverifier window
func NewInteropFilter(api interopFilterAPI, chainID uint256.Int) IngressFilter {
	return &interopAccessFilter{
		api:     api,
		timeout: 86400,
		chainID: chainID,
	}
}

// FilterTx implements IngressFilter.FilterTx
// it uses provided functions to get the access list from the transaction
// and check it against the supervisor
func (f *interopAccessFilter) FilterTx(ctx context.Context, tx *types.Transaction) bool {
	hashes := f.api.TxToInteropAccessList(tx)
	// if there are no interop access list entries, allow the transaction (there is no interop check to perform)
	if len(hashes) == 0 {
		return true
	}
	t, err := f.api.CurrentInteropBlockTime()
	// if there are interop access list entries, but the interop API is not available, reject the transaction
	if err != nil {
		return false
	}
	// if the transaction is older than the preverifier window, reject it eagerly
	expireTime := time.Unix(int64(t), 0).Add(time.Duration(-f.timeout) * time.Second)
	if tx.Time().Compare(expireTime) < 0 {
		return false
	}
	exDesc := interoptypes.ExecutingDescriptor{Timestamp: t, Timeout: f.timeout, ChainID: f.chainID}

	// perform the interop check and update internal failsafe bool if needed
	return f.api.CheckAccessList(ctx, hashes, interoptypes.CrossUnsafe, exDesc) == nil
}
