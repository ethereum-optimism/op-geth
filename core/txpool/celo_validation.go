package txpool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/exchange"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// AcceptSet is a set of accepted transaction types for a transaction subpool.
type AcceptSet = map[uint8]struct{}

// CeloValidationOptions define certain differences between transaction validation
// across the different pools without having to duplicate those checks.
// In comparison to the standard ValidationOptions, the Accept field has been
// changed to allow to test for CeloDynamicFeeTx types.
type CeloValidationOptions struct {
	Config *params.ChainConfig // Chain configuration to selectively validate based on current fork rules

	AcceptSet AcceptSet // Set of transaction types that should be accepted for the calling pool
	MaxSize   uint64    // Maximum size of a transaction that the caller can meaningfully handle
	MinTip    *big.Int  // Minimum gas tip needed to allow a transaction into the caller pool

	EffectiveGasCeil uint64 // if non-zero, a gas ceiling to enforce independent of the header's gaslimit value
}

// NewAcceptSet creates a new AcceptSet with the types provided.
func NewAcceptSet(types ...uint8) AcceptSet {
	m := make(AcceptSet, len(types))
	for _, t := range types {
		m[t] = struct{}{}
	}
	return m
}

// Accepts returns true iff txType is accepted by this CeloValidationOptions.
func (cvo *CeloValidationOptions) Accepts(txType uint8) bool {
	_, ok := cvo.AcceptSet[txType]
	return ok
}

// CeloValidateTransaction is a helper method to check whether a transaction is valid
// according to the consensus rules, but does not check state-dependent validation
// (balance, nonce, etc).
//
// This check is public to allow different transaction pools to check the basic
// rules without duplicating code and running the risk of missed updates.
func CeloValidateTransaction(tx *types.Transaction, head *types.Header,
	signer types.Signer, opts *CeloValidationOptions, currencyCtx common.FeeCurrencyContext) error {
	if err := ValidateTransaction(tx, head, signer, opts, currencyCtx); err != nil {
		return err
	}
	if !common.IsCurrencyAllowed(currencyCtx.ExchangeRates, tx.FeeCurrency()) {
		return exchange.ErrUnregisteredFeeCurrency
	}

	return nil
}
