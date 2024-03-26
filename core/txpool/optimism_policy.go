package txpool

import (
	"github.com/ethereum/go-ethereum/core/types"
)

type OptimismTxPolicyStatus uint

const (
	OptimismTxPolicyInvalid OptimismTxPolicyStatus = iota
	OptimismTxPolicyValid
)

type OptimismTxPoolPolicy interface {
	// Run validation logic on the transaction prior to pool submission
	//
	// TOOD: Look into taking in the entire batch such that validation
	//       can be parallelized by the implementor if necessary
	ValidateTx(tx *types.Transaction) (OptimismTxPolicyStatus, error)
}

var _ OptimismTxPoolPolicy = &NoOpTxPoolPolicy{}

type NoOpTxPoolPolicy struct{}

func (p *NoOpTxPoolPolicy) ValidateTx(tx *types.Transaction) (OptimismTxPolicyStatus, error) {
	return OptimismTxPolicyValid, nil
}
