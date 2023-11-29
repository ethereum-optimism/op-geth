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
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

type RollupGasData struct {
	Zeroes, Ones uint64
}

func (r RollupGasData) DataGas(time uint64, cfg *params.ChainConfig) (gas uint64) {
	gas = r.Zeroes * params.TxDataZeroGas
	if cfg.IsRegolith(time) {
		gas += r.Ones * params.TxDataNonZeroGasEIP2028
	} else {
		gas += (r.Ones + 68) * params.TxDataNonZeroGasEIP2028
	}
	return gas
}

type StateGetter interface {
	GetState(common.Address, common.Hash) common.Hash
}

// L1CostFunc is used in the state transition to determine the cost of a rollup message.
// Returns nil if there is no cost.
type L1CostFunc func(blockNum uint64, blockTime uint64, dataGas RollupGasData, isDepositTx bool) *big.Int

var (
	L1BaseFeeSlot  = common.BigToHash(big.NewInt(1))
	OverheadSlot   = common.BigToHash(big.NewInt(5))
	ScalarSlot     = common.BigToHash(big.NewInt(6))
	TokenRatioSlot = common.BigToHash(big.NewInt(0))
)

var (
	L1BlockAddr   = common.HexToAddress("0x4200000000000000000000000000000000000015")
	GasOracleAddr = common.HexToAddress("0x420000000000000000000000000000000000000F")
	Decimals      = big.NewInt(1_000_000)
)

// NewL1CostFunc returns a function used for calculating L1 fee cost.
// This depends on the oracles because gas costs can change over time.
// It returns nil if there is no applicable cost function.
func NewL1CostFunc(config *params.ChainConfig, statedb StateGetter) L1CostFunc {
	cacheBlockNum := ^uint64(0)
	var l1BaseFee, overhead, scalar, tokenRatio *big.Int
	return func(blockNum uint64, blockTime uint64, dataGas RollupGasData, isDepositTx bool) *big.Int {
		rollupDataGas := dataGas.DataGas(blockTime, config) // Only fake txs for RPC view-calls are 0.
		if config.Optimism == nil || isDepositTx || rollupDataGas == 0 {
			return common.Big0
		}
		if blockNum != cacheBlockNum {
			l1BaseFee = statedb.GetState(L1BlockAddr, L1BaseFeeSlot).Big()
			overhead = statedb.GetState(L1BlockAddr, OverheadSlot).Big()
			scalar = statedb.GetState(L1BlockAddr, ScalarSlot).Big()
			tokenRatio = statedb.GetState(GasOracleAddr, TokenRatioSlot).Big()
			cacheBlockNum = blockNum
		}
		return L1Cost(rollupDataGas, l1BaseFee, overhead, scalar, tokenRatio)
	}
}

func L1Cost(rollupDataGas uint64, l1BaseFee, overhead, scalar, tokenRatio *big.Int) *big.Int {
	l1GasUsed := new(big.Int).SetUint64(rollupDataGas)
	l1GasUsed = l1GasUsed.Add(l1GasUsed, overhead)
	l1Cost := l1GasUsed.Mul(l1GasUsed, l1BaseFee)
	l1Cost = l1Cost.Mul(l1Cost, scalar)
	l1Cost = l1Cost.Mul(l1Cost, tokenRatio)
	return l1Cost.Div(l1Cost, Decimals)
}

// DeriveL1GasInfo reads L1 gas related information to be included
// on the receipt
func DeriveL1GasInfo(state StateGetter) (*big.Int, *big.Int, *big.Int, *big.Float, *big.Int) {
	l1BaseFee, overhead, scalar, scaled := readL1BlockStorageSlots(L1BlockAddr, state)
	tokenRatio := readGPOStorageSlots(GasOracleAddr, state)
	return l1BaseFee, overhead, scalar, scaled, tokenRatio
}

func readL1BlockStorageSlots(addr common.Address, state StateGetter) (*big.Int, *big.Int, *big.Int, *big.Float) {
	l1BaseFee := state.GetState(addr, L1BaseFeeSlot)
	overhead := state.GetState(addr, OverheadSlot)
	scalar := state.GetState(addr, ScalarSlot)
	scaled := scaleDecimals(scalar.Big(), Decimals)
	return l1BaseFee.Big(), overhead.Big(), scalar.Big(), scaled
}

func readGPOStorageSlots(addr common.Address, state StateGetter) *big.Int {
	tokenRatio := state.GetState(addr, TokenRatioSlot)
	return tokenRatio.Big()
}

// scaleDecimals will scale a value by decimals
func scaleDecimals(scalar, divisor *big.Int) *big.Float {
	fscalar := new(big.Float).SetInt(scalar)
	fdivisor := new(big.Float).SetInt(divisor)
	// fscalar / fdivisor
	return new(big.Float).Quo(fscalar, fdivisor)
}
