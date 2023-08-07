package vm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/params"
)

type CeloPrecompiledContract interface {
	RequiredGas(input []byte) uint64                              // RequiredGas calculates the contract gas use
	Run(input []byte, ctx *celoPrecompileContext) ([]byte, error) // Run runs the precompiled contract
}

type wrap struct {
	PrecompiledContract
}

func (pw *wrap) Run(input []byte, ctx *celoPrecompileContext) ([]byte, error) {
	return pw.PrecompiledContract.Run(input)
}

type celoPrecompileContext struct {
	*BlockContext
	*params.Rules

	caller common.Address
	evm    *EVM
}

func NewContext(caller common.Address, evm *EVM) *celoPrecompileContext {
	return &celoPrecompileContext{
		BlockContext: &evm.Context,
		Rules:        &evm.chainRules,
		caller:       caller,
		evm:          evm,
	}
}

func celoPrecompileAddress(index byte) common.Address {
	celoPrecompiledContractsAddressOffset := byte(0xff)
	return common.BytesToAddress(append([]byte{0}, (celoPrecompiledContractsAddressOffset - index)))
}

func (ctx *celoPrecompileContext) IsCallerGoldToken() (bool, error) {
	return ctx.BlockContext.GoldTokenAddress == ctx.caller, nil
}

// Native transfer contract to make Celo Gold ERC20 compatible.
type transfer struct{}

func (c *transfer) RequiredGas(input []byte) uint64 {
	return params.CallValueTransferGas
}

func (c *transfer) Run(input []byte, ctx *celoPrecompileContext) ([]byte, error) {
	if isGoldToken, err := ctx.IsCallerGoldToken(); err != nil {
		return nil, err
	} else if !isGoldToken {
		return nil, fmt.Errorf("Unable to call transfer from unpermissioned address")
	}

	// input is comprised of 3 arguments:
	//   from:  32 bytes representing the address of the sender
	//   to:    32 bytes representing the address of the recipient
	//   value: 32 bytes, a 256 bit integer representing the amount of Celo Gold to transfer
	// 3 arguments x 32 bytes each = 96 bytes total input
	if len(input) != 96 {
		return nil, ErrInputLength
	}

	from := common.BytesToAddress(input[0:32])
	to := common.BytesToAddress(input[32:64])

	var parsed bool
	value, parsed := math.ParseBig256(hexutil.Encode(input[64:96]))
	if !parsed {
		return nil, fmt.Errorf("Error parsing transfer: unable to parse value from " + hexutil.Encode(input[64:96]))
	}

	if from == common.ZeroAddress {
		// Mint case: Create cGLD out of thin air
		ctx.evm.StateDB.AddBalance(to, value)
	} else {
		// Fail if we're trying to transfer more than the available balance
		if !ctx.CanTransfer(ctx.evm.StateDB, from, value) {
			return nil, ErrInsufficientBalance
		}

		ctx.Transfer(ctx.evm.StateDB, from, to, value)
	}

	return input, nil
}
