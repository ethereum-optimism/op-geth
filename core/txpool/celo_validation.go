package txpool

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
)

var NonWhitelistedFeeCurrencyError = errors.New("Fee currency given is not whitelisted at current block")

// FeeCurrencyValidator validates currency whitelisted status at the specified
// block number.
type FeeCurrencyValidator interface {
	IsWhitelisted(st *state.StateDB, feeCurrency *common.Address) bool
	// Balance returns the feeCurrency balance of the address specified, in the given state.
	// If feeCurrency is nil, the native currency balance has to be returned.
	Balance(st *state.StateDB, address common.Address, feeCurrency *common.Address) *uint256.Int
}

func NewFeeCurrencyValidator() FeeCurrencyValidator {
	return &feeval{}
}

type feeval struct {
}

func (f *feeval) IsWhitelisted(st *state.StateDB, feeCurrency *common.Address) bool {
	// TODO: implement proper validation for all currencies
	// Hardcoded for the moment
	return true
	//return feeCurrency == nil
}

func (f *feeval) Balance(st *state.StateDB, address common.Address, feeCurrency *common.Address) *uint256.Int {
	// TODO: implement proper balance retrieval for fee currencies
	return st.GetBalance(address)
}

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
	signer types.Signer, opts *CeloValidationOptions, st *state.StateDB, fcv FeeCurrencyValidator) error {
	if err := ValidateTransaction(tx, head, signer, opts); err != nil {
		return err
	}
	if IsFeeCurrencyTx(tx) {
		if !fcv.IsWhitelisted(st, tx.FeeCurrency()) {
			return NonWhitelistedFeeCurrencyError
		}
	}
	return nil
}

// IsFeeCurrencyTxType returns true if and only if the transaction type
// given can handle custom gas fee currencies.
func IsFeeCurrencyTxType(t uint8) bool {
	return t == types.CeloDynamicFeeTxType
}

// IsFeeCurrencyTx returns true if this transaction specifies a custom
// gas fee currency.
func IsFeeCurrencyTx(tx *types.Transaction) bool {
	return IsFeeCurrencyTxType(tx.Type()) && tx.FeeCurrency() != nil
}

// See: txpool.ValidationOptionsWithState
type CeloValidationOptionsWithState struct {
	ValidationOptionsWithState

	// FeeCurrencyValidator allows for balance check of non native fee currencies.
	FeeCurrencyValidator FeeCurrencyValidator
}
