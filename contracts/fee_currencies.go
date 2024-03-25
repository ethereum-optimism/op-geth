package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
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

// Debits transaction fees from the transaction sender and stores them in the temporary address
func DebitFees(evm *vm.EVM, feeCurrency *common.Address, address common.Address, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) == 0 {
		return nil
	}
	tokenAbi, err := abigen.FeeCurrencyMetaData.GetAbi()
	if err != nil {
		return err
	}
	// Solidity: function debitGasFees(address from, uint256 value)
	input, err := tokenAbi.Pack("debitGasFees", address, amount)
	if err != nil {
		return fmt.Errorf("pack debitGasFees: %w", err)
	}

	caller := vm.AccountRef(common.ZeroAddress)

	ret, leftoverGas, err := evm.Call(caller, *feeCurrency, input, maxGasForDebitGasFeesTransactions, new(uint256.Int))
	gasUsed := maxGasForDebitGasFeesTransactions - leftoverGas
	log.Trace("DebitFees called", "feeCurrency", *feeCurrency, "gasUsed", gasUsed)
	if err != nil {
		revertReason, err2 := abi.UnpackRevert(ret)
		if err2 == nil {
			return fmt.Errorf("DebitFees reverted: %s", revertReason)
		}
	}
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

	tokenAbi, err := abigen.FeeCurrencyMetaData.GetAbi()
	if err != nil {
		return err
	}
	// Solidity:
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
	input, err := tokenAbi.Pack("creditGasFees", txSender, tipReceiver, common.ZeroAddress, baseFeeReceiver, refund, feeTip, common.Big0, baseFee)
	if err != nil {
		return fmt.Errorf("pack creditGasFees: %w", err)
	}

	caller := vm.AccountRef(common.ZeroAddress)
	ret, leftoverGas, err := evm.Call(caller, *feeCurrency, input, maxGasForCreditGasFeesTransactions, new(uint256.Int))
	gasUsed := maxGasForCreditGasFeesTransactions - leftoverGas
	log.Trace("CreditFees called", "feeCurrency", *feeCurrency, "gasUsed", gasUsed)
	if err != nil {
		revertReason, err2 := abi.UnpackRevert(ret)
		if err2 == nil {
			return fmt.Errorf("CreditFees reverted: %s", revertReason)
		}
	}
	return err
}
