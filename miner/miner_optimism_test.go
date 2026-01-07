package miner

import (
	"encoding/binary"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

const testDAFootprintGasScalar = 400

// TestDAFootprintMining tests that the miner correctly limits the DA footprint of the block.
// It builds a block via the miner from txpool
// transactions and then imports the block into the chain, asserting that
// execution succeeds.
func TestDAFootprintMining(t *testing.T) {
	requirePreJovianBehavior := func(t *testing.T, block *types.Block, receipts []*types.Receipt) {
		var txGas uint64
		for _, receipt := range receipts {
			txGas += receipt.GasUsed
		}
		require.Equal(t, txGas, block.GasUsed(), "total tx gas used should be equal to block gas used")
		require.Zero(t, *block.Header().BlobGasUsed, "expected 0 blob gas used")
	}

	requireLargeDAFootprintBehavior := func(t *testing.T, block *types.Block, receipts []*types.Receipt) {
		var (
			txGas       uint64
			daFootprint uint64
			txs         = block.Transactions()
		)

		require.Equal(t, len(receipts), len(txs))

		for i, receipt := range receipts {
			txGas += receipt.GasUsed
			if txs[i].IsDepositTx() {
				continue
			}
			daFootprint += txs[i].RollupCostData().EstimatedDASize().Uint64() * testDAFootprintGasScalar
		}
		require.Equal(t, txGas, block.GasUsed(), "total tx gas used should be equal to block gas used")
		require.Greater(t, daFootprint, block.GasUsed(), "total DA footprint used should be greater than block gas used")
		require.LessOrEqual(t, daFootprint, block.GasLimit(), "total DA footprint used should be less or equal block gas limit")
	}
	t.Run("jovian-one-min-tx", func(t *testing.T) {
		testMineAndExecute(t, 0, jovianConfig(), func(t *testing.T, _ *core.BlockChain, block *types.Block, receipts []*types.Receipt) {
			require.Len(t, receipts, 2) // 1 test pending tx and 1 deposit tx
			requireLargeDAFootprintBehavior(t, block, receipts)

			// Double-confirm DA footprint calculation manually in this simple transaction case.
			daFootprint, err := types.CalcDAFootprint(block.Transactions())
			require.NoError(t, err, "failed to calculate DA footprint")
			require.Equal(t, daFootprint, *block.Header().BlobGasUsed,
				"header blob gas used should match calculated DA footprint")
			require.Equal(t, testDAFootprintGasScalar*types.MinTransactionSize.Uint64(), daFootprint,
				"simple pending transaction should lead to min DA footprint")
		})
	})
	t.Run("jovian-at-limit", func(t *testing.T) {
		testMineAndExecute(t, 17, jovianConfig(), func(t *testing.T, _ *core.BlockChain, block *types.Block, receipts []*types.Receipt) {
			require.Len(t, receipts, 19) // including 1 test pending tx and 1 deposit tx
			requireLargeDAFootprintBehavior(t, block, receipts)
		})
	})
	t.Run("jovian-above-limit", func(t *testing.T) {
		testMineAndExecute(t, 18, jovianConfig(), func(t *testing.T, _ *core.BlockChain, block *types.Block, receipts []*types.Receipt) {
			require.Len(t, receipts, 19) // same as for 17, because 18th tx from pool shouldn't have been included
			requireLargeDAFootprintBehavior(t, block, receipts)
		})
	})
	t.Run("isthmus", func(t *testing.T) {
		testMineAndExecute(t, 39, isthmusConfig(), func(t *testing.T, _ *core.BlockChain, block *types.Block, receipts []*types.Receipt) {
			require.Len(t, receipts, 41) // including 1 test pending tx and 1 deposit tx
			requirePreJovianBehavior(t, block, receipts)
		})
	})

	t.Run("jovian-invalid-blobGasUsed", func(t *testing.T) {
		testMineAndExecute(t, 0, jovianConfig(), func(t *testing.T, bc *core.BlockChain, block *types.Block, receipts []*types.Receipt) {
			require.Len(t, receipts, 2) // 1 test pending tx and 1 deposit tx
			header := block.Header()
			*header.BlobGasUsed += 1 // invalidate blobGasUsed
			invalidBlock := block.WithSeal(header)
			_, err := bc.InsertChain(types.Blocks{invalidBlock})
			require.ErrorContains(t, err, "invalid DA footprint in blobGasUsed field (remote: 40001 local: 40000)")
		})
	})
}

func testMineAndExecute(t *testing.T, numTxs uint64, cfg *params.ChainConfig, assertFn func(*testing.T, *core.BlockChain, *types.Block, []*types.Receipt)) {
	db := rawdb.NewMemoryDatabase()
	w, b := newTestWorker(t, cfg, beacon.New(ethash.NewFaker()), db, 0)

	// Start from nonce 1 to avoid colliding with the preloaded pending tx.
	txs := genTxs(1, numTxs)

	// Add to txpool for the miner to pick up.
	if errs := b.txPool.Add(txs, false); len(errs) > 0 {
		for _, err := range errs {
			require.NoError(t, err, "failed adding tx to pool")
		}
	}

	parent := b.chain.CurrentBlock()
	ts := parent.Time + 12
	dtx := new(types.DepositTx)
	if cfg.IsJovian(parent.Time) {
		dtx = jovianDepositTx(testDAFootprintGasScalar)
	}

	genParams := &generateParams{
		parentHash:    b.chain.CurrentBlock().Hash(),
		timestamp:     ts,
		withdrawals:   types.Withdrawals{},
		beaconRoot:    new(common.Hash),
		gasLimit:      ptr(uint64(1e6)), // Small gas limit to easily fill block
		txs:           types.Transactions{types.NewTx(dtx)},
		eip1559Params: eip1559.EncodeHolocene1559Params(250, 6),
	}
	if cfg.IsJovian(ts) {
		genParams.minBaseFee = new(uint64)
	}
	r := w.generateWork(genParams, false)
	require.NoError(t, r.err, "block generation failed")
	require.NotNil(t, r.block, "no block generated")

	assertFn(t, b.chain, r.block, r.receipts)

	// Import the block into the chain, which executes it via StateProcessor.
	_, err := b.chain.InsertChain(types.Blocks{r.block})
	require.NoError(t, err, "block import/execution failed")
}

func jovianDepositTx(daFootprintGasScalar uint16) *types.DepositTx {
	data := make([]byte, types.JovianL1AttributesLen)
	copy(data[0:4], types.JovianL1AttributesSelector)
	binary.BigEndian.PutUint16(data[types.JovianL1AttributesLen-2:types.JovianL1AttributesLen], daFootprintGasScalar)
	return &types.DepositTx{Data: data}
}

func ptr[T any](v T) *T {
	return &v
}
