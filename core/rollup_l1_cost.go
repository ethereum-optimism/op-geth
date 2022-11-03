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

var (
	L1BaseFeeSlot = common.BigToHash(big.NewInt(1))
	OverheadSlot  = common.BigToHash(big.NewInt(5))
	ScalarSlot    = common.BigToHash(big.NewInt(6))
)

var L1BlockAddr = common.HexToAddress("0x4200000000000000000000000000000000000015")

// NewL1CostFunc returns a function used for calculating L1 fee cost.
// This depends on the oracles because gas costs can change over time.
// It returns nil if there is no applicable cost function.
func NewL1CostFunc(config *params.ChainConfig, statedb vm.StateDB) vm.L1CostFunc {
	cacheBlockNum := ^uint64(0)
	var l1BaseFee, overhead, scalar *big.Int
	return func(blockNum uint64, msg vm.RollupMessage) *big.Int {
		rollupDataGas := msg.RollupDataGas() // Only fake txs for RPC view-calls are 0.
		if config.Optimism == nil || msg.IsDepositTx() || rollupDataGas == 0 {
			return nil
		}
		if blockNum != cacheBlockNum {
			l1BaseFee = statedb.GetState(L1BlockAddr, L1BaseFeeSlot).Big()
			overhead = statedb.GetState(L1BlockAddr, OverheadSlot).Big()
			scalar = statedb.GetState(L1BlockAddr, ScalarSlot).Big()
			cacheBlockNum = blockNum
		}
		return types.L1Cost(rollupDataGas, l1BaseFee, overhead, scalar)
	}
}
