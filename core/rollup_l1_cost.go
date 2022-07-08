// Copyright 2022 The go-ethereum Authors
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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

var big10 = big.NewInt(10)

var (
	L1BaseFeeSlot = common.BigToHash(big.NewInt(1))
	OverheadSlot  = common.BigToHash(big.NewInt(3))
	ScalarSlot    = common.BigToHash(big.NewInt(4))
	DecimalsSlot  = common.BigToHash(big.NewInt(5))
)

var (
	OVM_GasPriceOracleAddr = common.HexToAddress("0x420000000000000000000000000000000000000F")
	L1BlockAddr            = common.HexToAddress("0x4200000000000000000000000000000000000015")
)

// NewL1CostFunc returns a function used for calculating L1 fee cost.
// This depends on the oracles because gas costs can change over time.
// It returns nil if there is no applicable cost function.
func NewL1CostFunc(config *params.ChainConfig, statedb vm.StateDB) vm.L1CostFunc {
	cacheBlockNum := ^uint64(0)
	var l1BaseFee, overhead, scalar, decimals, divisor *big.Int
	return func(blockNum uint64, msg vm.RollupMessage) *big.Int {
		rollupDataGas := msg.RollupDataGas() // Only fake txs for RPC view-calls are 0.
		if config.Optimism == nil || msg.Nonce() == types.DepositsNonce || rollupDataGas == 0 {
			return nil
		}
		if blockNum != cacheBlockNum {
			l1BaseFee = statedb.GetState(L1BlockAddr, L1BaseFeeSlot).Big()
			overhead = statedb.GetState(OVM_GasPriceOracleAddr, OverheadSlot).Big()
			scalar = statedb.GetState(OVM_GasPriceOracleAddr, ScalarSlot).Big()
			decimals = statedb.GetState(OVM_GasPriceOracleAddr, DecimalsSlot).Big()
			divisor = new(big.Int).Exp(big10, decimals, nil)
			cacheBlockNum = blockNum
		}
		l1GasUsed := new(big.Int).SetUint64(rollupDataGas)
		l1GasUsed = l1GasUsed.Add(l1GasUsed, overhead)
		l1Cost := l1GasUsed.Mul(l1GasUsed, l1BaseFee)
		l1Cost = l1Cost.Mul(l1Cost, scalar)
		l1Cost = l1Cost.Div(l1Cost, divisor)
		return l1Cost
	}
}
