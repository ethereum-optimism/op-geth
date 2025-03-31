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

package legacypool

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"
)

func setupOPStackPool() (*LegacyPool, *ecdsa.PrivateKey) {
	return setupPoolWithConfig(params.OptimismTestConfig)
}

func setupTestL1FeeParams(t *testing.T, pool *LegacyPool, l1FeeScalar byte) {
	l1FeeScalars := common.Hash{19: l1FeeScalar}
	// sanity check
	l1BaseFeeScalar, l1BlobBaseFeeScalar := types.ExtractEcotoneFeeParams(l1FeeScalars[:])
	require.EqualValues(t, l1FeeScalar, l1BaseFeeScalar.Uint64())
	require.Zero(t, l1BlobBaseFeeScalar.Sign())
	pool.currentState.SetState(types.L1BlockAddr, types.L1FeeScalarsSlot, l1FeeScalars)
	l1BaseFee := big.NewInt(1e6) // to account for division by 1e12 in L1 cost
	pool.currentState.SetState(types.L1BlockAddr, types.L1BaseFeeSlot, common.BigToHash(l1BaseFee))
	// sanity checks
	require.Equal(t, l1BaseFee, pool.currentState.GetState(types.L1BlockAddr, types.L1BaseFeeSlot).Big())
}

func setupTestOperatorFeeParams(t *testing.T, pool *LegacyPool, opFeeConst byte) {
	opFeeParams := common.Hash{31: opFeeConst} // 0 scalar
	// sanity check
	s, c := types.ExtractOperatorFeeParams(opFeeParams)
	require.Zero(t, s.Sign())
	require.EqualValues(t, opFeeConst, c.Uint64())
	pool.currentState.SetState(types.L1BlockAddr, types.OperatorFeeParamsSlot, opFeeParams)
}

func TestInvalidRollupTransactions(t *testing.T) {
	t.Run("zero-rollup-cost", func(t *testing.T) {
		testInvalidRollupTransactions(t, nil)
	})

	t.Run("l1-cost", func(t *testing.T) {
		testInvalidRollupTransactions(t, setupTestL1FeeParams)
	})

	t.Run("operator-cost", func(t *testing.T) {
		testInvalidRollupTransactions(t, setupTestOperatorFeeParams)
	})
}

func testInvalidRollupTransactions(t *testing.T, stateMod func(t *testing.T, pool *LegacyPool, x byte)) {
	t.Parallel()

	pool, key := setupOPStackPool()
	defer pool.Close()

	const gasLimit = 100_000
	tx := transaction(0, gasLimit, key)
	from, _ := deriveSender(tx)

	// base fee is 1
	testAddBalance(pool, from, new(big.Int).Add(big.NewInt(gasLimit), tx.Value()))
	// we add the test variant with zero rollup cost as a sanity check that the tx would indeed be valid
	if stateMod == nil {
		require.NoError(t, pool.addRemote(tx))
		return
	}

	// Now we cause insufficient funds error due to rollup cost
	stateMod(t, pool, 1)

	rcost := pool.rollupCostFn(tx)
	require.Equal(t, 1, rcost.Sign(), "rollup cost must be >0")

	require.ErrorIs(t, pool.addRemote(tx), core.ErrInsufficientFunds)
}

func TestRollupTransactionCostAccounting(t *testing.T) {
	t.Run("zero-rollup-cost", func(t *testing.T) {
		testRollupTransactionCostAccounting(t, nil)
	})

	t.Run("l1-cost", func(t *testing.T) {
		testRollupTransactionCostAccounting(t, setupTestL1FeeParams)
	})

	t.Run("operator-cost", func(t *testing.T) {
		testRollupTransactionCostAccounting(t, setupTestOperatorFeeParams)
	})
}

func testRollupTransactionCostAccounting(t *testing.T, stateMod func(t *testing.T, pool *LegacyPool, x byte)) {
	t.Parallel()

	pool, key := setupOPStackPool()
	defer pool.Close()

	const gasLimit = 100_000
	gasPrice0, gasPrice1 := big.NewInt(100), big.NewInt(110)
	tx0 := pricedTransaction(0, gasLimit, gasPrice0, key)
	tx1 := pricedTransaction(0, gasLimit, gasPrice1, key)
	from, _ := deriveSender(tx0)

	require.NotNil(t, pool.rollupCostFn)

	if stateMod != nil {
		stateMod(t, pool, 1)
	}

	cost0, of := txpool.TotalTxCost(tx0, pool.rollupCostFn)
	require.False(t, of)

	if stateMod != nil {
		require.Greater(t, cost0.Uint64(), tx0.Cost().Uint64(), "tx0 total cost should be greater than regular cost")
	}

	// we add the initial tx to the pool
	testAddBalance(pool, from, cost0.ToBig())
	require.NoError(t, pool.addRemoteSync(tx0))
	_, ok := pool.queue[from]
	require.False(t, ok, "tx0 should not be in queue, but pending")
	pending, ok := pool.pending[from]
	require.True(t, ok, "tx0 should be pending")
	require.Equal(t, cost0, pending.totalcost, "tx0 total pending cost should match")

	pool.reset(nil, nil) // reset the rollup cost function, simulates a head change
	if stateMod != nil {
		stateMod(t, pool, 2) // change rollup params to cause higher cost
		cost0r, of := txpool.TotalTxCost(tx0, pool.rollupCostFn)
		require.False(t, of)
		require.Greater(t, cost0r.Uint64(), cost0.Uint64(), "new tx0 cost should be larger")
	}
	cost1, of := txpool.TotalTxCost(tx1, pool.rollupCostFn)
	require.False(t, of)
	// add just enough for the replacement tx
	testAddBalance(pool, from, new(uint256.Int).Sub(cost1, cost0).ToBig())
	// now we add the replacement and check the accounting
	require.NoError(t, pool.addRemoteSync(tx1))
	_, ok = pool.queue[from]
	require.False(t, ok, "tx1 should not be in queue, but pending")
	pending, ok = pool.pending[from]
	require.True(t, ok, "tx1 should be pending")
	require.Equal(t, cost1, pending.totalcost, "tx1 total pending cost should match")
}

// TestRollupCostFuncChange tests that changes in the underlying rollup cost parameters
// are correctly picked up by the transaction pool and the underlying list implementation.
func TestRollupCostFuncChange(t *testing.T) {
	t.Parallel()

	pool, key := setupOPStackPool()
	defer pool.Close()

	const gasLimit = 100_000
	gasPrice := big.NewInt(100)
	tx0 := pricedTransaction(0, gasLimit, gasPrice, key)
	tx1 := pricedTransaction(1, gasLimit, gasPrice, key)
	from, _ := deriveSender(tx0)

	require.NotNil(t, pool.rollupCostFn)

	setupTestOperatorFeeParams(t, pool, 10)

	cost0, of := txpool.TotalTxCost(tx0, pool.rollupCostFn)
	require.False(t, of)

	// 1st add tx0, consuming all balance
	testAddBalance(pool, from, cost0.ToBig())
	require.NoError(t, pool.addRemoteSync(tx0))

	// 2nd add same balance but increase op fee const by 10
	// so adding 2nd tx should fail with 10 missing.
	testAddBalance(pool, from, cost0.ToBig())
	pool.reset(nil, nil) // reset the rollup cost function, simulates a head change
	setupTestOperatorFeeParams(t, pool, 20)
	require.ErrorContains(t, pool.addRemoteSync(tx1), "overshot 10")

	// 3rd now add missing 10, adding tx1 should succeed
	testAddBalance(pool, from, big.NewInt(10))
	require.NoError(t, pool.addRemoteSync(tx1))
}
