package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/addresses"
	"github.com/ethereum/go-ethereum/contracts/celo/abigen"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
)

const (
	Thousand = 1000

	// Default intrinsic gas cost of transactions paying for gas in alternative currencies.
	// Calculated to estimate 1 balance read, 1 debit, and 4 credit transactions.
	IntrinsicGasForAlternativeFeeCurrency uint64 = 50 * Thousand
	maxAllowedGasForDebitAndCredit        uint64 = 3 * IntrinsicGasForAlternativeFeeCurrency
)

var feeCurrencyABI *abi.ABI

func init() {
	var err error
	feeCurrencyABI, err = abigen.FeeCurrencyMetaData.GetAbi()
	if err != nil {
		panic(err)
	}
}

// Returns nil if debit is possible, used in tx pool validation
func TryDebitFees(tx *types.Transaction, from common.Address, backend *CeloBackend) error {
	amount := new(big.Int).SetUint64(tx.Gas())
	amount.Mul(amount, tx.GasFeeCap())

	snapshot := backend.State.Snapshot()
	err := DebitFees(backend.NewEVM(), tx.FeeCurrency(), from, amount)
	backend.State.RevertToSnapshot(snapshot)
	return err
}

// Debits transaction fees from the transaction sender and stores them in the temporary address
func DebitFees(evm *vm.EVM, feeCurrency *common.Address, address common.Address, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}

	leftoverGas, err := evm.CallWithABI(
		feeCurrencyABI, "debitGasFees", *feeCurrency, maxAllowedGasForDebitAndCredit,
		// debitGasFees(address from, uint256 value) parameters
		address, amount,
	)
	gasUsed := maxAllowedGasForDebitAndCredit - leftoverGas
	evm.Context.GasUsedForDebit = gasUsed
	log.Trace("DebitFees called", "feeCurrency", *feeCurrency, "gasUsed", gasUsed)
	return err
}

// Credits fees to the respective parties
// - the base fee goes to the fee handler
// - the transaction tip goes to the miner
// - the l1 data fee goes the the data fee receiver, is the node runs in rollup mode
// - remaining funds are refunded to the transaction sender
func CreditFees(
	evm *vm.EVM,
	feeCurrency *common.Address,
	txSender, tipReceiver, baseFeeReceiver, l1DataFeeReceiver common.Address,
	refund, feeTip, baseFee, l1DataFee *big.Int,
) error {
	// Our old `creditGasFees` function does not accept an l1DataFee and
	// the fee currencies do not implement the new interface yet. Since tip
	// and data fee both go to the sequencer, we can work around that for
	// now by addint the l1DataFee to the tip.
	if l1DataFee != nil {
		feeTip = new(big.Int).Add(feeTip, l1DataFee)
	}

	// Not all fee currencies can handle a receiver being the zero address.
	// In that case send the fee to the base fee recipient, which we know is non-zero.
	if tipReceiver.Cmp(common.ZeroAddress) == 0 {
		tipReceiver = baseFeeReceiver
	}

	maxAllowedGasForCredit := maxAllowedGasForDebitAndCredit - evm.Context.GasUsedForDebit
	leftoverGas, err := evm.CallWithABI(
		feeCurrencyABI, "creditGasFees", *feeCurrency, maxAllowedGasForCredit,
		// function creditGasFees(
		// 	address from,
		// 	address feeRecipient,
		// 	address, // gatewayFeeRecipient, unused
		// 	address communityFund,
		// 	uint256 refund,
		// 	uint256 tipTxFee,
		// 	uint256, // gatewayFee, unused
		// 	uint256 baseTxFee
		// )
		txSender, tipReceiver, common.ZeroAddress, baseFeeReceiver, refund, feeTip, common.Big0, baseFee,
	)

	gasUsed := maxAllowedGasForCredit - leftoverGas
	log.Trace("CreditFees called", "feeCurrency", *feeCurrency, "gasUsed", gasUsed)

	gasUsedForDebitAndCredit := evm.Context.GasUsedForDebit + gasUsed
	if gasUsedForDebitAndCredit > IntrinsicGasForAlternativeFeeCurrency {
		log.Info("Gas usage for debit+credit exceeds intrinsic gas!", "gasUsed", gasUsedForDebitAndCredit, "intrinsicGas", IntrinsicGasForAlternativeFeeCurrency, "feeCurrency", feeCurrency)
	}
	return err
}

// GetExchangeRates returns the exchange rates for all gas currencies from CELO
func GetExchangeRates(caller bind.ContractCaller) (common.ExchangeRates, error) {
	exchangeRates := map[common.Address]*big.Rat{}
	directory, err := abigen.NewFeeCurrencyDirectoryCaller(addresses.FeeCurrencyDirectoryAddress, caller)
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to access FeeCurrencyDirectory: %w", err)
	}

	registeredTokens, err := directory.GetCurrencies(&bind.CallOpts{})
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to get whitelisted tokens: %w", err)
	}
	for _, tokenAddress := range registeredTokens {
		rate, err := directory.GetExchangeRate(&bind.CallOpts{}, tokenAddress)
		if err != nil {
			log.Error("Failed to get medianRate for gas currency!", "err", err, "tokenAddress", tokenAddress.Hex())
			continue
		}
		if rate.Numerator.Sign() <= 0 || rate.Denominator.Sign() <= 0 {
			log.Error("Bad exchange rate for fee currency", "tokenAddress", tokenAddress.Hex(), "numerator", rate.Numerator, "denominator", rate.Denominator)
			continue
		}
		exchangeRates[tokenAddress] = new(big.Rat).SetFrac(rate.Numerator, rate.Denominator)
	}

	return exchangeRates, nil
}

// GetBalanceERC20 returns an account's balance on a given ERC20 currency
func GetBalanceERC20(caller bind.ContractCaller, accountOwner common.Address, contractAddress common.Address) (result *big.Int, err error) {
	token, err := abigen.NewFeeCurrencyCaller(contractAddress, caller)
	if err != nil {
		return nil, fmt.Errorf("failed to access FeeCurrency: %w", err)
	}

	balance, err := token.BalanceOf(&bind.CallOpts{}, accountOwner)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// GetFeeBalance returns the account's balance from the specified feeCurrency
// (if feeCurrency is nil or ZeroAddress, native currency balance is returned).
func GetFeeBalance(backend *CeloBackend, account common.Address, feeCurrency *common.Address) *big.Int {
	if feeCurrency == nil || *feeCurrency == common.ZeroAddress {
		return backend.State.GetBalance(account).ToBig()
	}
	balance, err := GetBalanceERC20(backend, account, *feeCurrency)
	if err != nil {
		log.Error("Error while trying to get ERC20 balance:", "cause", err, "contract", feeCurrency.Hex(), "account", account.Hex())
	}
	return balance
}
