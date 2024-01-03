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

package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

type RollupCostData struct {
	zeroes, ones uint64
}

func NewRollupCostData(data []byte) (out RollupCostData) {
	for _, b := range data {
		if b == 0 {
			out.zeroes++
		} else {
			out.ones++
		}
	}
	return out
}

type StateGetter interface {
	GetState(common.Address, common.Hash) common.Hash
}

// L1CostFunc is used in the state transition to determine the L1 data fee charged to the sender of
// non-Deposit transactions.
// It returns nil if no L1 data fee is charged.
type L1CostFunc func(rcd RollupCostData, blockTime uint64) *big.Int

// l1CostFunc is an internal version of L1CostFunc that also returns the gasUsed for use in
// receipts.
type l1CostFunc func(rcd RollupCostData) (fee, gasUsed *big.Int)

var (
	L1BasefeeSlot = common.BigToHash(big.NewInt(1))
	OverheadSlot  = common.BigToHash(big.NewInt(5))
	ScalarSlot    = common.BigToHash(big.NewInt(6))
)

var L1BlockAddr = common.HexToAddress("0x4200000000000000000000000000000000000015")

// NewL1CostFunc returns a function used for calculating L1 fee cost, or nil if this is not an
// op-stack chain.
func NewL1CostFunc(config *params.ChainConfig, statedb StateGetter) L1CostFunc {
	if config.Optimism == nil {
		return nil
	}
	forBlock := ^uint64(0)
	var cachedFunc l1CostFunc
	return func(rollupCostData RollupCostData, blockTime uint64) *big.Int {
		if rollupCostData == (RollupCostData{}) {
			return nil // Do not charge if there is no rollup cost-data (e.g. RPC call or deposit).
		}
		if forBlock != blockTime {
			// Note: The following variables are not initialized from the state DB until this point
			// to allow deposit transactions from the block to be processed first by state
			// transition.  This behavior is consensus critical!
			l1Basefee := statedb.GetState(L1BlockAddr, L1BasefeeSlot).Big()
			overhead := statedb.GetState(L1BlockAddr, OverheadSlot).Big()
			scalar := statedb.GetState(L1BlockAddr, ScalarSlot).Big()
			isRegolith := config.IsRegolith(blockTime)
			cachedFunc = newL1CostFunc(l1Basefee, overhead, scalar, isRegolith)
			if forBlock != ^uint64(0) {
				// best practice is not to re-use l1 cost funcs across different blocks, but we
				// make it work just in case.
				log.Info("l1 cost func re-used for different L1 block", "oldTime", forBlock, "newTime", blockTime)
			}
			forBlock = blockTime
		}
		fee, _ := cachedFunc(rollupCostData)
		return fee
	}
}

var (
	oneMillion = big.NewInt(1_000_000)
)

func newL1CostFunc(l1Basefee, overhead, scalar *big.Int, isRegolith bool) l1CostFunc {
	return func(rollupCostData RollupCostData) (fee, gasUsed *big.Int) {
		if rollupCostData == (RollupCostData{}) {
			return nil, nil // Do not charge if there is no rollup cost-data (e.g. RPC call or deposit)
		}
		gas := rollupCostData.zeroes * params.TxDataZeroGas
		if isRegolith {
			gas += rollupCostData.ones * params.TxDataNonZeroGasEIP2028
		} else {
			gas += (rollupCostData.ones + 68) * params.TxDataNonZeroGasEIP2028
		}
		gasWithOverhead := new(big.Int).SetUint64(gas)
		gasWithOverhead.Add(gasWithOverhead, overhead)
		l1Cost := l1CostHelper(gasWithOverhead, l1Basefee, scalar)
		return l1Cost, gasWithOverhead
	}
}

// extractL1GasParams extracts the gas parameters necessary to compute gas costs from L1 block info
// calldata.
func extractL1GasParams(config *params.ChainConfig, time uint64, data []byte) (l1Basefee *big.Int, costFunc l1CostFunc, feeScalar *big.Float, err error) {
	// data consists of func selector followed by 7 ABI-encoded parameters (32 bytes each)
	if len(data) < 4+32*8 {
		return nil, nil, nil, fmt.Errorf("expected at least %d L1 info bytes, got %d", 4+32*8, len(data))
	}
	data = data[4:]                                      // trim function selector
	l1Basefee = new(big.Int).SetBytes(data[32*2 : 32*3]) // arg index 2
	overhead := new(big.Int).SetBytes(data[32*6 : 32*7]) // arg index 6
	scalar := new(big.Int).SetBytes(data[32*7 : 32*8])   // arg index 7
	fscalar := new(big.Float).SetInt(scalar)             // legacy: format fee scalar as big Float
	fdivisor := new(big.Float).SetUint64(1_000_000)      // 10**6, i.e. 6 decimals
	feeScalar = new(big.Float).Quo(fscalar, fdivisor)
	costFunc = newL1CostFunc(l1Basefee, overhead, scalar, config.IsRegolith(time))
	return
}

// L1Cost computes the the L1 data fee. It is used by e2e tests so must remain exported.
func L1Cost(rollupDataGas uint64, l1Basefee, overhead, scalar *big.Int) *big.Int {
	l1GasUsed := new(big.Int).SetUint64(rollupDataGas)
	l1GasUsed.Add(l1GasUsed, overhead)
	return l1CostHelper(l1GasUsed, l1Basefee, scalar)
}

func l1CostHelper(gasWithOverhead, l1Basefee, scalar *big.Int) *big.Int {
	fee := new(big.Int).Set(gasWithOverhead)
	fee.Mul(fee, l1Basefee).Mul(fee, scalar).Div(fee, oneMillion)
	return fee
}
