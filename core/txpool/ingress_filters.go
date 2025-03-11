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

type checkFn func(context.Context, []common.Hash, interoptypes.SafetyLevel, uint64) error

func NewInteropFilter(fn checkFn) IngressFilter {
	return &interopAccessFilter{
		checkAccess: fn,
	}
}

type interopAccessFilter struct {
	checkAccess checkFn
}

// FilterTx implements IngressFilter.FilterTx
// it takes the access list from the transaction and checks it against the supervisor
func (f *interopAccessFilter) FilterTx(ctx context.Context, tx *types.Transaction) bool {
	al := tx.AccessList()

	if len(al) == 0 {
		return true
	}

	// I don't really think this is the right way to turn an access list into a list of hashes
	// I'm just plugging it in for now since it returns the correct type.
	hashes := make([]common.Hash, 0, len(al))
	for i := range al {
		hashes[i] = al[i].StorageKeys[0]
	}

	// Note - the ExecutingDescriptor is not used here, but it is required by the interop client
	// I'm just passing in 0 for now since it's not used. Not sure what the effect of this will be.
	err := f.checkAccess(ctx, hashes, interoptypes.CrossUnsafe, 0)
	return err == nil
}
