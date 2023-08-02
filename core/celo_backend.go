package core

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// CeloBackend provide a partial ContractBackend implementation, so that we can
// access core contracts during block processing.
type CeloBackend struct {
	chainConfig *params.ChainConfig
	state       *state.StateDB
}

// ContractCaller implementation

func (b *CeloBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.state.GetCode(contract), nil
}

func (b *CeloBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	// Ensure message is initialized properly.
	if call.Gas == 0 {
		// Chosen to be the same as ethconfig.Defaults.RPCGasCap
		call.Gas = 50000000
	}
	if call.Value == nil {
		call.Value = new(big.Int)
	}

	// Minimal initialization, might need to be extended when CeloBackend
	// is used in more places. Also initializing blockNumber and time with
	// 0 works now, but will break once we add hardforks at a later time.
	if blockNumber == nil {
		blockNumber = common.Big0
	}
	blockCtx := vm.BlockContext{BlockNumber: blockNumber, Time: 0}
	txCtx := vm.TxContext{}
	vmConfig := vm.Config{}

	evm := vm.NewEVM(blockCtx, txCtx, b.state, b.chainConfig, vmConfig)
	ret, _, err := evm.StaticCall(vm.AccountRef(evm.Origin), *call.To, call.Data, call.Gas)

	return ret, err
}
