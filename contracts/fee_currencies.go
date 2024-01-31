package contracts

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/celo/abigen"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

const (
	Thousand = 1000
	Million  = 1000 * 1000

	maxGasForDebitGasFeesTransactions  uint64 = 1 * Million
	maxGasForCreditGasFeesTransactions uint64 = 1 * Million
	// Default intrinsic gas cost of transactions paying for gas in alternative currencies.
	// Calculated to estimate 1 balance read, 1 debit, and 4 credit transactions.
	IntrinsicGasForAlternativeFeeCurrency uint64 = 50 * Thousand
)

var (
	tmpAddress = common.HexToAddress("0xce106a5")

	// ErrNonWhitelistedFeeCurrency is returned if the currency specified to use for the fees
	// isn't one of the currencies whitelisted for that purpose.
	ErrNonWhitelistedFeeCurrency = errors.New("non-whitelisted fee currency address")
)

// GetBalanceOf returns an account's balance on a given ERC20 currency
func GetBalanceOf(caller bind.ContractCaller, accountOwner common.Address, contractAddress common.Address) (result *big.Int, err error) {
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

func ConvertGoldToCurrency(exchangeRates map[common.Address]*big.Rat, feeCurrency *common.Address, goldAmount *big.Int) (*big.Int, error) {
	exchangeRate, ok := exchangeRates[*feeCurrency]
	if !ok {
		return nil, ErrNonWhitelistedFeeCurrency
	}
	return new(big.Int).Div(new(big.Int).Mul(goldAmount, exchangeRate.Num()), exchangeRate.Denom()), nil
}

// Debits transaction fees from the transaction sender and stores them in the temporary address
func DebitFees(evm *vm.EVM, feeCurrency *common.Address, address common.Address, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}
	abi, err := abigen.FeeCurrencyMetaData.GetAbi()
	if err != nil {
		return err
	}
	// Solidity: function transfer(address to, uint256 amount) returns(bool)
	input, err := abi.Pack("transfer", tmpAddress, amount)
	if err != nil {
		return err
	}

	caller := vm.AccountRef(address)

	_, leftoverGas, err := evm.Call(caller, *feeCurrency, input, maxGasForDebitGasFeesTransactions, new(uint256.Int))
	gasUsed := maxGasForDebitGasFeesTransactions - leftoverGas
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
	caller := vm.AccountRef(tmpAddress)
	leftoverGas := maxGasForCreditGasFeesTransactions

	abi, err := abigen.FeeCurrencyMetaData.GetAbi()
	if err != nil {
		return err
	}

	leftoverGas, err = transferErc20(evm, abi, caller, feeCurrency, baseFeeReceiver, baseFee, leftoverGas)
	if err != nil {
		return fmt.Errorf("base fee: %w", err)
	}

	// If the tip is non-zero, but the coinbase of the miner is not set, send the tip back to the tx sender
	unusedFee, leftoverGas, err := creditFee(evm, abi, caller, feeCurrency, tipReceiver, feeTip, leftoverGas, "fee tip")
	if err != nil {
		return err
	}
	refund.Add(refund, unusedFee)

	// If the data fee is non-zero, but the data fee receiver is not set, send the tip back to the tx sender
	unusedFee, leftoverGas, err = creditFee(evm, abi, caller, feeCurrency, l1DataFeeReceiver, l1DataFee, leftoverGas, "l1 data fee")
	if err != nil {
		return err
	}
	refund.Add(refund, unusedFee)

	unusedFee, leftoverGas, err = creditFee(evm, abi, caller, feeCurrency, txSender, refund, leftoverGas, "refund")
	if err != nil {
		return err
	}
	if unusedFee.Cmp(common.Big0) != 0 {
		return errors.New("could not refund remaining fees to sender")
	}

	gasUsed := maxGasForCreditGasFeesTransactions - leftoverGas
	log.Trace("creditFees called", "feeCurrency", *feeCurrency, "gasUsed", gasUsed)
	return nil
}

func creditFee(evm *vm.EVM, abi *abi.ABI, caller vm.AccountRef, feeCurrency *common.Address, receiver common.Address, amount *big.Int, gasLimit uint64, action string) (*big.Int, uint64, error) {
	if amount != nil && amount.Cmp(common.Big0) == 1 {
		if receiver != common.ZeroAddress {
			leftoverGas, err := transferErc20(evm, abi, caller, feeCurrency, receiver, amount, gasLimit)
			if err != nil {
				return common.Big0, leftoverGas, fmt.Errorf("%s: %w", action, err)
			}
		} else {
			return amount, gasLimit, nil
		}
	}
	return common.Big0, gasLimit, nil
}

func transferErc20(evm *vm.EVM, abi *abi.ABI, caller vm.AccountRef, feeCurrency *common.Address, receiver common.Address, amount *big.Int, gasLimit uint64) (uint64, error) {
	// Solidity: function transfer(address to, uint256 amount) returns(bool)
	transferData, err := abi.Pack("transfer", receiver, amount)
	if err != nil {
		return 0, fmt.Errorf("pack transfer: %w", err)
	}
	_, leftoverGas, err := evm.Call(caller, *feeCurrency, transferData, gasLimit, new(uint256.Int))
	if err != nil {
		return 0, fmt.Errorf("call transfer: %w", err)
	}
	return leftoverGas, nil
}
