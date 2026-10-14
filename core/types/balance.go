package types

//go:generate stringer -type=BalanceChangeReason

// BalanceChangeReason describes the reason why an account's balance was changed in a transaction
type BalanceChangeReason uint8

const (
	// ValidatorReward is the balance change caused by the block validator reward
	ValidatorReward BalanceChangeReason = iota

	// EmptyAccountByCode0 is the balance change caused by the empty account by code0 rule
	EmptyAccountByCode0

	// StartNonce is the balance change caused by the start nonce
	StartNonce

	// GasRefund is the balance change caused by refunding gas
	GasRefund

	// BalanceMint is for balance minting operations
	BalanceMint

	// BalanceBurn is for balance burning operations
	BalanceBurn

	// BlockReward is for block mining reward
	BlockReward

	// FeeTip is for transaction fee tips
	FeeTip

	// UncleReward is for uncle block rewards
	UncleReward

	// TransactionFee is for transaction fees
	TransactionFee

	// L1Fee is for Layer 1 fees
	L1Fee

	// Deposit is for deposits
	Deposit

	// BaseFeeBurn is for base fee burning
	BaseFeeBurn
)