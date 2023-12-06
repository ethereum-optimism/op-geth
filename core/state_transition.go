// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	cmath "github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"golang.org/x/crypto/sha3"
)

var (
	BVM_ETH_ADDR = common.HexToAddress("0xdEAddEaDdeadDEadDEADDEAddEADDEAddead1111")
)

// ExecutionResult includes all output after executing given evm
// message no matter the execution itself is successful or not.
type ExecutionResult struct {
	UsedGas    uint64 // Total used gas but include the refunded gas
	Err        error  // Any error encountered during the execution(listed in core/vm/errors.go)
	ReturnData []byte // Returned data from evm(function result or data supplied with revert opcode)
}

// Unwrap returns the internal evm error which allows us for further
// analysis outside.
func (result *ExecutionResult) Unwrap() error {
	return result.Err
}

// Failed returns the indicator whether the execution is successful or not
func (result *ExecutionResult) Failed() bool { return result.Err != nil }

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (result *ExecutionResult) Return() []byte {
	if result.Err != nil {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (result *ExecutionResult) Revert() []byte {
	if result.Err != vm.ErrExecutionReverted {
		return nil
	}
	return common.CopyBytes(result.ReturnData)
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, accessList types.AccessList, isContractCreation bool, isHomestead, isEIP2028 bool, isEIP3860 bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if isContractCreation && isHomestead {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	dataLen := uint64(len(data))
	// Bump the required gas by the amount of transactional data
	if dataLen > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		nonZeroGas := params.TxDataNonZeroGasFrontier
		if isEIP2028 {
			nonZeroGas = params.TxDataNonZeroGasEIP2028
		}
		if (math.MaxUint64-gas)/nonZeroGas < nz {
			return 0, ErrGasUintOverflow
		}
		gas += nz * nonZeroGas

		z := dataLen - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, ErrGasUintOverflow
		}
		gas += z * params.TxDataZeroGas

		if isContractCreation && isEIP3860 {
			lenWords := toWordSize(dataLen)
			if (math.MaxUint64-gas)/params.InitCodeWordGas < lenWords {
				return 0, ErrGasUintOverflow
			}
			gas += lenWords * params.InitCodeWordGas
		}
	}
	if accessList != nil {
		gas += uint64(len(accessList)) * params.TxAccessListAddressGas
		gas += uint64(accessList.StorageKeys()) * params.TxAccessListStorageKeyGas
	}
	return gas, nil
}

// toWordSize returns the ceiled word size required for init code payment calculation.
func toWordSize(size uint64) uint64 {
	if size > math.MaxUint64-31 {
		return math.MaxUint64/32 + 1
	}

	return (size + 31) / 32
}

type RunMode uint8

const (
	CommitMode RunMode = iota
	GasEstimationMode
	GasEstimationWithSkipCheckBalanceMode
	EthcallMode
)

// A Message contains the data derived from a single transaction that is relevant to state
// processing.
type Message struct {
	To         *common.Address
	From       common.Address
	Nonce      uint64
	Value      *big.Int
	GasLimit   uint64
	GasPrice   *big.Int
	GasFeeCap  *big.Int
	GasTipCap  *big.Int
	Data       []byte
	AccessList types.AccessList

	// When SkipAccountCheckss is true, the message nonce is not checked against the
	// account nonce in state. It also disables checking that the sender is an EOA.
	// This field will be set to true for operations like RPC eth_call.
	SkipAccountChecks bool

	IsSystemTx    bool                // IsSystemTx indicates the message, if also a deposit, does not emit gas usage.
	IsDepositTx   bool                // IsDepositTx indicates the message is force-included and can persist a mint.
	Mint          *big.Int            // Mint is the amount to mint before EVM processing, or nil if there is no minting.
	ETHValue      *big.Int            // ETHValue is the amount to mint BVM_ETH before EVM processing, or nil if there is no minting.
	MetaTxParams  *types.MetaTxParams // MetaTxParams contains necessary parameter to sponsor gas fee for msg.From.
	RollupDataGas types.RollupGasData // RollupDataGas indicates the rollup cost of the message, 0 if not a rollup or no cost.

	// runMode
	RunMode RunMode
}

// TransactionToMessage converts a transaction into a Message.
func TransactionToMessage(tx *types.Transaction, s types.Signer, baseFee *big.Int) (*Message, error) {
	metaTxParams, err := types.DecodeAndVerifyMetaTxParams(tx)
	if err != nil {
		return nil, err
	}
	msg := &Message{
		Nonce:             tx.Nonce(),
		GasLimit:          tx.Gas(),
		GasPrice:          new(big.Int).Set(tx.GasPrice()),
		GasFeeCap:         new(big.Int).Set(tx.GasFeeCap()),
		GasTipCap:         new(big.Int).Set(tx.GasTipCap()),
		To:                tx.To(),
		Value:             tx.Value(),
		Data:              tx.Data(),
		AccessList:        tx.AccessList(),
		IsSystemTx:        tx.IsSystemTx(),
		IsDepositTx:       tx.IsDepositTx(),
		Mint:              tx.Mint(),
		RollupDataGas:     tx.RollupDataGas(),
		ETHValue:          tx.ETHValue(),
		MetaTxParams:      metaTxParams,
		SkipAccountChecks: false,
		RunMode:           CommitMode,
	}
	// If baseFee provided, set gasPrice to effectiveGasPrice.
	if baseFee != nil {
		msg.GasPrice = cmath.BigMin(msg.GasPrice.Add(msg.GasTipCap, baseFee), msg.GasFeeCap)
	}
	msg.From, err = types.Sender(s, tx)
	return msg, err
}

// CalculateRollupGasDataFromMessage calculate RollupGasData from message.
func (st *StateTransition) CalculateRollupGasDataFromMessage() {
	tx := types.NewTx(&types.DynamicFeeTx{
		Nonce:     st.msg.Nonce,
		Value:     st.msg.Value,
		Gas:       st.msg.GasLimit,
		GasTipCap: st.msg.GasTipCap,
		GasFeeCap: st.msg.GasFeeCap,
		Data:      st.msg.Data,
	})

	st.msg.RollupDataGas = tx.RollupDataGas()

	// add a constant to cover sigs(V,R,S) and other data to make sure that the gasLimit from eth_estimateGas can cover L1 cost
	// just used for estimateGas and the actual L1 cost depends on users' tx when executing
	st.msg.RollupDataGas.Ones += 80
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(evm *vm.EVM, msg *Message, gp *GasPool) (*ExecutionResult, error) {
	return NewStateTransition(evm, msg, gp).TransitionDb()
}

// StateTransition represents a state transition.
//
// == The State Transitioning Model
//
// A state transition is a change made when a transaction is applied to the current world
// state. The state transitioning model does all the necessary work to work out a valid new
// state root.
//
//  1. Nonce handling
//  2. Pre pay gas
//  3. Create a new state object if the recipient is nil
//  4. Value transfer
//
// == If contract creation ==
//
//	4a. Attempt to run transaction data
//	4b. If valid, use result as code for the new state object
//
// == end ==
//
//  5. Run Script section
//  6. Derive new state root
type StateTransition struct {
	gp           *GasPool
	msg          *Message
	gasRemaining uint64
	initialGas   uint64
	state        vm.StateDB
	evm          *vm.EVM
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(evm *vm.EVM, msg *Message, gp *GasPool) *StateTransition {
	return &StateTransition{
		gp:    gp,
		evm:   evm,
		msg:   msg,
		state: evm.StateDB,
	}
}

// to returns the recipient of the message.
func (st *StateTransition) to() common.Address {
	if st.msg == nil || st.msg.To == nil /* contract creation */ {
		return common.Address{}
	}
	return *st.msg.To
}

func (st *StateTransition) buyGas() (*big.Int, error) {
	if err := st.applyMetaTransaction(); err != nil {
		return nil, err
	}
	mgval := new(big.Int).SetUint64(st.msg.GasLimit)
	mgval = mgval.Mul(mgval, st.msg.GasPrice)
	var l1Cost *big.Int
	if st.msg.RunMode == GasEstimationMode || st.msg.RunMode == GasEstimationWithSkipCheckBalanceMode {
		st.CalculateRollupGasDataFromMessage()
	}
	if st.evm.Context.L1CostFunc != nil && st.msg.RunMode != EthcallMode {
		l1Cost = st.evm.Context.L1CostFunc(st.evm.Context.BlockNumber.Uint64(), st.evm.Context.Time, st.msg.RollupDataGas, st.msg.IsDepositTx)
	}
	if l1Cost != nil && (st.msg.RunMode == GasEstimationMode || st.msg.RunMode == GasEstimationWithSkipCheckBalanceMode) {
		mgval = mgval.Add(mgval, l1Cost)
	}
	balanceCheck := mgval
	if st.msg.GasFeeCap != nil {
		balanceCheck = new(big.Int).SetUint64(st.msg.GasLimit)
		balanceCheck = balanceCheck.Mul(balanceCheck, st.msg.GasFeeCap)
		balanceCheck.Add(balanceCheck, st.msg.Value)
		if l1Cost != nil && st.msg.RunMode == GasEstimationMode {
			balanceCheck.Add(balanceCheck, l1Cost)
		}
	}
	if st.msg.RunMode != GasEstimationWithSkipCheckBalanceMode && st.msg.RunMode != EthcallMode {
		if st.msg.MetaTxParams != nil {
			pureGasFeeValue := new(big.Int).Sub(balanceCheck, st.msg.Value)
			sponsorAmount, selfPayAmount := types.CalculateSponsorPercentAmount(st.msg.MetaTxParams, pureGasFeeValue)
			if have, want := st.state.GetBalance(st.msg.MetaTxParams.GasFeeSponsor), sponsorAmount; have.Cmp(want) < 0 {
				return nil, fmt.Errorf("%w: gas fee sponsor %v have %v want %v", ErrInsufficientFunds, st.msg.MetaTxParams.GasFeeSponsor.Hex(), have, want)
			}
			selfPayAmount = new(big.Int).Add(selfPayAmount, st.msg.Value)
			if have, want := st.state.GetBalance(st.msg.From), selfPayAmount; have.Cmp(want) < 0 {
				return nil, fmt.Errorf("%w: address %v have %v want %v", ErrInsufficientFunds, st.msg.From.Hex(), have, want)
			}
		} else {
			if have, want := st.state.GetBalance(st.msg.From), balanceCheck; have.Cmp(want) < 0 {
				return nil, fmt.Errorf("%w: address %v have %v want %v", ErrInsufficientFunds, st.msg.From.Hex(), have, want)
			}
		}
	}

	if err := st.gp.SubGas(st.msg.GasLimit); err != nil {
		return nil, err
	}
	st.gasRemaining += st.msg.GasLimit

	st.initialGas = st.msg.GasLimit
	if st.msg.RunMode != GasEstimationWithSkipCheckBalanceMode && st.msg.RunMode != EthcallMode {
		if st.msg.MetaTxParams != nil {
			sponsorAmount, selfPayAmount := types.CalculateSponsorPercentAmount(st.msg.MetaTxParams, mgval)
			st.state.SubBalance(st.msg.MetaTxParams.GasFeeSponsor, sponsorAmount)
			st.state.SubBalance(st.msg.From, selfPayAmount)
			log.Debug("BuyGas for metaTx",
				"sponsor", st.msg.MetaTxParams.GasFeeSponsor.String(), "amount", sponsorAmount.String(),
				"user", st.msg.From.String(), "amount", selfPayAmount.String())
		} else {
			st.state.SubBalance(st.msg.From, mgval)
		}
	}
	return l1Cost, nil
}

func (st *StateTransition) applyMetaTransaction() error {
	if st.msg.MetaTxParams == nil {
		return nil
	}
	if st.msg.MetaTxParams.ExpireHeight < st.evm.Context.BlockNumber.Uint64() {
		return types.ErrExpiredMetaTx
	}

	st.msg.Data = st.msg.MetaTxParams.Payload
	return nil
}

func (st *StateTransition) preCheck() (*big.Int, error) {
	if st.msg.IsDepositTx {
		// No fee fields to check, no nonce to check, and no need to check if EOA (L1 already verified it for us)
		// Gas is free, but no refunds!
		st.initialGas = st.msg.GasLimit
		st.gasRemaining += st.msg.GasLimit // Add gas here in order to be able to execute calls.
		// Don't touch the gas pool for system transactions
		if st.msg.IsSystemTx {
			if st.evm.ChainConfig().IsOptimismRegolith(st.evm.Context.Time) {
				return nil, fmt.Errorf("%w: address %v", ErrSystemTxNotSupported,
					st.msg.From.Hex())
			}
			return common.Big0, nil
		}
		if err := st.gp.SubGas(st.msg.GasLimit); err != nil {
			return nil, err
		}
		return common.Big0, nil // gas used by deposits may not be used by other txs
	}
	// Only check transactions that are not fake
	msg := st.msg
	if !msg.SkipAccountChecks {
		// Make sure this transaction's nonce is correct.
		stNonce := st.state.GetNonce(msg.From)
		if msgNonce := msg.Nonce; stNonce < msgNonce {
			return nil, fmt.Errorf("%w: address %v, tx: %d state: %d", ErrNonceTooHigh,
				msg.From.Hex(), msgNonce, stNonce)
		} else if stNonce > msgNonce {
			return nil, fmt.Errorf("%w: address %v, tx: %d state: %d", ErrNonceTooLow,
				msg.From.Hex(), msgNonce, stNonce)
		} else if stNonce+1 < stNonce {
			return nil, fmt.Errorf("%w: address %v, nonce: %d", ErrNonceMax,
				msg.From.Hex(), stNonce)
		}
		// Make sure the sender is an EOA
		codeHash := st.state.GetCodeHash(msg.From)
		if codeHash != (common.Hash{}) && codeHash != types.EmptyCodeHash {
			return nil, fmt.Errorf("%w: address %v, codehash: %s", ErrSenderNoEOA,
				msg.From.Hex(), codeHash)
		}
	}

	// Make sure that transaction gasFeeCap is greater than the baseFee (post london)
	if st.evm.ChainConfig().IsLondon(st.evm.Context.BlockNumber) {
		// Skip the checks if gas fields are zero and baseFee was explicitly disabled (eth_call)
		if !st.evm.Config.NoBaseFee || msg.GasFeeCap.BitLen() > 0 || msg.GasTipCap.BitLen() > 0 {
			if l := msg.GasFeeCap.BitLen(); l > 256 {
				return nil, fmt.Errorf("%w: address %v, maxFeePerGas bit length: %d", ErrFeeCapVeryHigh,
					msg.From.Hex(), l)
			}
			if l := msg.GasTipCap.BitLen(); l > 256 {
				return nil, fmt.Errorf("%w: address %v, maxPriorityFeePerGas bit length: %d", ErrTipVeryHigh,
					msg.From.Hex(), l)
			}
			if msg.GasFeeCap.Cmp(msg.GasTipCap) < 0 {
				return nil, fmt.Errorf("%w: address %v, maxPriorityFeePerGas: %s, maxFeePerGas: %s", ErrTipAboveFeeCap,
					msg.From.Hex(), msg.GasTipCap, msg.GasFeeCap)
			}
			// This will panic if baseFee is nil, but basefee presence is verified
			// as part of header validation.
			if msg.GasFeeCap.Cmp(st.evm.Context.BaseFee) < 0 {
				return nil, fmt.Errorf("%w: address %v, maxFeePerGas: %s baseFee: %s", ErrFeeCapTooLow,
					msg.From.Hex(), msg.GasFeeCap, st.evm.Context.BaseFee)
			}
		}
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the evm execution result with following fields.
//
//   - used gas: total gas used (including gas being refunded)
//   - returndata: the returned data from evm
//   - concrete execution error: various EVM errors which abort the execution, e.g.
//     ErrOutOfGas, ErrExecutionReverted
//
// However if any consensus issue encountered, return the error directly with
// nil evm execution result.
func (st *StateTransition) TransitionDb() (*ExecutionResult, error) {
	if mint := st.msg.Mint; mint != nil {
		st.state.AddBalance(st.msg.From, mint)
	}
	//add eth value
	if ethValue := st.msg.ETHValue; ethValue != nil && ethValue.Cmp(big.NewInt(0)) != 0 {
		st.addBVMETHBalance(ethValue)
		st.addBVMETHTotalSupply(ethValue)
		st.generateBVMETHMintEvent(*st.msg.To, ethValue)
	}
	snap := st.state.Snapshot()

	result, err := st.innerTransitionDb()
	// Failed deposits must still be included. Unless we cannot produce the block at all due to the gas limit.
	// On deposit failure, we rewind any state changes from after the minting, and increment the nonce.
	if err != nil && err != ErrGasLimitReached && st.msg.IsDepositTx {
		st.state.RevertToSnapshot(snap)
		// Even though we revert the state changes, always increment the nonce for the next deposit transaction
		st.state.SetNonce(st.msg.From, st.state.GetNonce(st.msg.From)+1)
		// Record deposits as using all their gas (matches the gas pool)
		// System Transactions are special & are not recorded as using any gas (anywhere)
		// Regolith changes this behaviour so the actual gas used is reported.
		// In this case the tx is invalid so is recorded as using all gas.
		gasUsed := st.msg.GasLimit
		if st.msg.IsSystemTx && !st.evm.ChainConfig().IsRegolith(st.evm.Context.Time) {
			gasUsed = 0
		}
		result = &ExecutionResult{
			UsedGas:    gasUsed,
			Err:        fmt.Errorf("failed deposit: %w", err),
			ReturnData: nil,
		}
		err = nil
	}
	return result, err
}

func (st *StateTransition) innerTransitionDb() (*ExecutionResult, error) {
	// First check this message satisfies all consensus rules before
	// applying the message. The rules include these clauses
	//
	// 1. the nonce of the message caller is correct
	// 2. caller has enough balance to cover transaction fee(gaslimit * gasprice)
	// 3. the amount of gas required is available in the block
	// 4. the purchased gas is enough to cover intrinsic usage
	// 5. there is no overflow when calculating intrinsic gas
	// 6. caller has enough balance to cover asset transfer for **topmost** call

	// Check clauses 1-3, buy gas if everything is correct
	tokenRatio := st.state.GetState(types.L1BlockAddr, types.TokenRatioSlot).Big().Uint64()
	l1Cost, err := st.preCheck()
	if err != nil {
		return nil, err
	}

	if st.evm.Config.Debug {
		st.evm.Config.Tracer.CaptureTxStart(st.initialGas)
		defer func() {
			st.evm.Config.Tracer.CaptureTxEnd(st.gasRemaining)
		}()
	}

	var (
		msg              = st.msg
		sender           = vm.AccountRef(msg.From)
		rules            = st.evm.ChainConfig().Rules(st.evm.Context.BlockNumber, st.evm.Context.Random != nil, st.evm.Context.Time)
		contractCreation = msg.To == nil
	)

	// Check clauses 4-5, subtract intrinsic gas if everything is correct
	gas, err := IntrinsicGas(msg.Data, msg.AccessList, contractCreation, rules.IsHomestead, rules.IsIstanbul, rules.IsShanghai)
	if err != nil {
		return nil, err
	}
	if !st.msg.IsDepositTx && !st.msg.IsSystemTx {
		gas = gas * tokenRatio
	}
	if st.gasRemaining < gas {
		return nil, fmt.Errorf("%w: have %d, want %d", ErrIntrinsicGas, st.gasRemaining, gas)
	}
	st.gasRemaining -= gas

	var l1Gas uint64
	if !st.msg.IsDepositTx && !st.msg.IsSystemTx {
		if st.msg.GasPrice.Cmp(common.Big0) > 0 && l1Cost != nil {
			l1Gas = new(big.Int).Div(l1Cost, st.msg.GasPrice).Uint64()
			if st.msg.GasLimit < l1Gas {
				return nil, fmt.Errorf("%w: have %d, want %d", ErrIntrinsicGas, st.gasRemaining, l1Gas)
			}
		}
		if st.gasRemaining < l1Gas {
			return nil, fmt.Errorf("%w: have %d, want %d", ErrIntrinsicGas, st.gasRemaining, l1Gas)
		}
		st.gasRemaining -= l1Gas
		if tokenRatio > 0 {
			st.gasRemaining = st.gasRemaining / tokenRatio
		}
	}

	// Check clause 6
	if msg.Value.Sign() > 0 && !st.evm.Context.CanTransfer(st.state, msg.From, msg.Value) {
		return nil, fmt.Errorf("%w: address %v", ErrInsufficientFundsForTransfer, msg.From.Hex())
	}

	// Check whether the init code size has been exceeded.
	if rules.IsShanghai && contractCreation && len(msg.Data) > params.MaxInitCodeSize {
		return nil, fmt.Errorf("%w: code size %v limit %v", ErrMaxInitCodeSizeExceeded, len(msg.Data), params.MaxInitCodeSize)
	}

	// Execute the preparatory steps for state transition which includes:
	// - prepare accessList(post-berlin)
	// - reset transient storage(eip 1153)
	st.state.Prepare(rules, msg.From, st.evm.Context.Coinbase, msg.To, vm.ActivePrecompiles(rules), msg.AccessList)

	var (
		ret   []byte
		vmerr error // vm errors do not effect consensus and are therefore not assigned to err
	)
	if contractCreation {
		ret, _, st.gasRemaining, vmerr = st.evm.Create(sender, msg.Data, st.gasRemaining, msg.Value)
	} else {
		// Increment the nonce for the next transaction
		st.state.SetNonce(msg.From, st.state.GetNonce(sender.Address())+1)
		ret, st.gasRemaining, vmerr = st.evm.Call(sender, st.to(), msg.Data, st.gasRemaining, msg.Value)
	}

	// if deposit: skip refunds, skip tipping coinbase
	// Regolith changes this behaviour to report the actual gasUsed instead of always reporting all gas used.
	if st.msg.IsDepositTx && !rules.IsOptimismRegolith {
		// Record deposits as using all their gas (matches the gas pool)
		// System Transactions are special & are not recorded as using any gas (anywhere)
		gasUsed := st.msg.GasLimit
		if st.msg.IsSystemTx {
			gasUsed = 0
		}
		return &ExecutionResult{
			UsedGas:    gasUsed,
			Err:        vmerr,
			ReturnData: ret,
		}, nil
	}
	// Note for deposit tx there is no ETH refunded for unused gas, but that's taken care of by the fact that gasPrice
	// is always 0 for deposit tx. So calling refundGas will ensure the gasUsed accounting is correct without actually
	// changing the sender's balance
	if !st.msg.IsDepositTx && !st.msg.IsSystemTx {
		if !rules.IsLondon {
			// Before EIP-3529: refunds were capped to gasUsed / 2
			st.refundGas(params.RefundQuotient, tokenRatio)
		} else {
			// After EIP-3529: refunds are capped to gasUsed / 5
			st.refundGas(params.RefundQuotientEIP3529, tokenRatio)
		}
	}

	if st.msg.IsDepositTx && rules.IsOptimismRegolith {
		// Skip coinbase payments for deposit tx in Regolith
		return &ExecutionResult{
			UsedGas:    st.gasUsed(),
			Err:        vmerr,
			ReturnData: ret,
		}, nil
	}
	effectiveTip := msg.GasPrice
	if rules.IsLondon {
		effectiveTip = cmath.BigMin(msg.GasTipCap, new(big.Int).Sub(msg.GasFeeCap, st.evm.Context.BaseFee))
	}

	if st.evm.Config.NoBaseFee && msg.GasFeeCap.Sign() == 0 && msg.GasTipCap.Sign() == 0 {
		// Skip fee payment when NoBaseFee is set and the fee fields
		// are 0. This avoids a negative effectiveTip being applied to
		// the coinbase when simulating calls.
	} else {
		fee := new(big.Int).SetUint64(st.gasUsed())
		fee.Mul(fee, effectiveTip)
		st.state.AddBalance(st.evm.Context.Coinbase, fee)
	}

	// Check that we are post bedrock to enable op-geth to be able to create pseudo pre-bedrock blocks (these are pre-bedrock, but don't follow l2 geth rules)
	// Note optimismConfig will not be nil if rules.IsOptimismBedrock is true
	if optimismConfig := st.evm.ChainConfig().Optimism; optimismConfig != nil && rules.IsOptimismBedrock {
		st.state.AddBalance(params.OptimismBaseFeeRecipient, new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.evm.Context.BaseFee))
		// Can not collect l1 fee here again, all l1 fee has been collected by CoinBase & OptimismBaseFeeRecipient
		//if cost := st.evm.Context.L1CostFunc(st.evm.Context.BlockNumber.Uint64(), st.evm.Context.Time, st.msg.RollupDataGas, st.msg.IsDepositTx); cost != nil {
		//	st.state.AddBalance(params.OptimismL1FeeRecipient, cost)
		//}
	}

	return &ExecutionResult{
		UsedGas:    st.gasUsed(),
		Err:        vmerr,
		ReturnData: ret,
	}, nil
}

func (st *StateTransition) refundGas(refundQuotient, tokenRatio uint64) {
	if st.msg.RunMode == GasEstimationWithSkipCheckBalanceMode || st.msg.RunMode == EthcallMode {
		st.gasRemaining = st.gasRemaining * tokenRatio
		st.gp.AddGas(st.gasRemaining)
		return
	}
	// Apply refund counter, capped to a refund quotient
	refund := st.gasUsed() / refundQuotient
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gasRemaining += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	st.gasRemaining = st.gasRemaining * tokenRatio
	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gasRemaining), st.msg.GasPrice)
	if st.msg.MetaTxParams != nil {
		sponsorRefundAmount, selfRefundAmount := types.CalculateSponsorPercentAmount(st.msg.MetaTxParams, remaining)
		st.state.AddBalance(st.msg.MetaTxParams.GasFeeSponsor, sponsorRefundAmount)
		st.state.AddBalance(st.msg.From, selfRefundAmount)
		log.Debug("RefundGas for metaTx",
			"sponsor", st.msg.MetaTxParams.GasFeeSponsor.String(), "refundAmount", sponsorRefundAmount.String(),
			"user", st.msg.From.String(), "refundAmount", selfRefundAmount.String())
	} else {
		st.state.AddBalance(st.msg.From, remaining)
	}

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gasRemaining)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	return st.initialGas - st.gasRemaining
}

func (st *StateTransition) addBVMETHBalance(ethValue *big.Int) {
	key := getBVMETHBalanceKey(*st.msg.To)
	value := st.state.GetState(BVM_ETH_ADDR, key)
	bal := value.Big()
	bal = bal.Add(bal, ethValue)
	st.state.SetState(BVM_ETH_ADDR, key, common.BigToHash(bal))
}

func (st *StateTransition) addBVMETHTotalSupply(ethValue *big.Int) {
	key := getBVMETHTotalSupplyKey()
	value := st.state.GetState(BVM_ETH_ADDR, key)
	bal := value.Big()
	bal = bal.Add(bal, ethValue)
	st.state.SetState(BVM_ETH_ADDR, key, common.BigToHash(bal))
}

func getBVMETHBalanceKey(addr common.Address) common.Hash {
	position := common.Big0
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(common.LeftPadBytes(addr.Bytes(), 32))
	hasher.Write(common.LeftPadBytes(position.Bytes(), 32))
	digest := hasher.Sum(nil)
	return common.BytesToHash(digest)
}

func getBVMETHTotalSupplyKey() common.Hash {
	return common.BytesToHash(common.Big2.Bytes())
}

func (st *StateTransition) generateBVMETHMintEvent(mintAddress common.Address, mintValue *big.Int) {
	// keccak("Mint(address,uint256)") = "0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885"
	methodHash := common.HexToHash("0x0f6798a560793a54c3bcfe86a93cde1e73087d944c0ea20544137d4121396885")
	topics := make([]common.Hash, 2)
	topics[0] = methodHash
	topics[1] = mintAddress.Hash()
	//data means the mint amount in MINT EVENT.
	d := common.HexToHash(common.Bytes2Hex(mintValue.Bytes())).Bytes()
	st.evm.StateDB.AddLog(&types.Log{
		Address: BVM_ETH_ADDR,
		Topics:  topics,
		Data:    d,
		// This is a non-consensus field, but assigned here because
		// core/state doesn't know the current block number.
		BlockNumber: st.evm.Context.BlockNumber.Uint64(),
	})
}
