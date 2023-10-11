package txpool

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

var NonWhitelistedFeeCurrencyError = errors.New("Fee currency given is not whitelisted at current block")

// FeeCurrencyValidator validates currency whitelisted status at the specified
// block number.
type FeeCurrencyValidator interface {
	IsWhitelisted(feeCurrency *common.Address, at *big.Int) bool
}

// AcceptMap is a set of accepted transaction types for a transaction subpool.
type AcceptMap = map[uint8]struct{}

// CeloValidationOptions define certain differences between transaction validation
// across the different pools without having to duplicate those checks.
// In comparison to the standard ValidationOptions, the Accept field has been
// changed to allow to test for CeloDynamicFeeTx types.
type CeloValidationOptions struct {
	Config *params.ChainConfig // Chain configuration to selectively validate based on current fork rules

	AcceptMap AcceptMap // Map of transaction types that should be accepted for the calling pool
	MaxSize   uint64    // Maximum size of a transaction that the caller can meaningfully handle
	MinTip    *big.Int  // Minimum gas tip needed to allow a transaction into the caller pool
}

// Accepts returns true iff txType is accepted by this CeloValidationOptions.
func (cvo *CeloValidationOptions) Accepts(txType uint8) bool {
	_, ok := cvo.AcceptMap[txType]
	return ok
}

// CeloValidateTransaction is a helper method to check whether a transaction is valid
// according to the consensus rules, but does not check state-dependent validation
// (balance, nonce, etc).
//
// This check is public to allow different transaction pools to check the basic
// rules without duplicating code and running the risk of missed updates.
func CeloValidateTransaction(tx *types.Transaction, head *types.Header,
	signer types.Signer, opts *CeloValidationOptions, fcv FeeCurrencyValidator) error {

	if err := ValidateTransaction(tx, head, signer, opts); err != nil {
		return err
	}
	if FeeCurrencyTx(tx) {
		if !fcv.IsWhitelisted(tx.FeeCurrency(), head.Number) {
			return NonWhitelistedFeeCurrencyError
		}
	}
	return nil
}

// NewAcceptMap creates a new AcceptMap with the types provided as keys.
func NewAcceptMap(types ...uint8) AcceptMap {
	m := make(AcceptMap, len(types))
	for _, t := range types {
		m[t] = struct{}{}
	}
	return m
}

// FeeCurrencyTxType returns true if and only if the transaction type
// given can handle custom gas fee currencies.
func FeeCurrencyTxType(t uint8) bool {
	return t == types.CeloDynamicFeeTxType
}

// FeeCurrencyTx returns true if this transaction specifies a custom
// gas fee currency.
func FeeCurrencyTx(tx *types.Transaction) bool {
	return FeeCurrencyTxType(tx.Type()) && tx.FeeCurrency() != nil
}
