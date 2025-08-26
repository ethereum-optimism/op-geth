// Copyright 2025 The op-geth Authors
// This file is part of the op-geth library.
//
// The op-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The op-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the op-geth library. If not, see <http://www.gnu.org/licenses/>.

package txpool

import (
	"math/big"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/core/types"
)

type RollupCostFunc func(tx types.RollupTransaction) *uint256.Int

type RollupTransaction interface {
	types.RollupTransaction
	Cost() *big.Int
}

// TotalTxCost returns the transaction's total cost, including the regular execution cost and
// the rollup costs (L1 costs and operator costs).
// Only regular costs apply if rollupCostFn is nil.
func TotalTxCost(tx RollupTransaction, rollupCostFn RollupCostFunc) (*uint256.Int, bool) {
	cost, overflow := uint256.FromBig(tx.Cost())
	if overflow {
		return nil, true
	} else if rollupCostFn == nil {
		return cost, false
	}

	rollupCost := rollupCostFn(tx)
	if rollupCost == nil {
		return cost, false
	}
	return cost.AddOverflow(cost, rollupCost)
}
