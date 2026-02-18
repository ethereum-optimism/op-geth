package downloader

import (
	"math/big"
	"math/rand"
	"slices"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

func TestCorrectReceipts(t *testing.T) {
	type testcase struct {
		blockNum uint64
		nonces   []uint64
		txTypes  []uint8
		validate func([]*types.ReceiptForStorage, []*types.ReceiptForStorage)
	}

	validateNonceDiff := func(diffIdxs ...int) func(original []*types.ReceiptForStorage, corrected []*types.ReceiptForStorage) {
		return func(original, corrected []*types.ReceiptForStorage) {
			for i, orig := range original {
				corr := corrected[i]
				if slices.Contains(diffIdxs, i) {
					// expect different deposit nonce for this index
					assert.NotEqual(t, orig.DepositNonce, corr.DepositNonce)
					// but all other fields that are in RLP storage should equal
					assert.Equal(t, orig.CumulativeGasUsed, corr.CumulativeGasUsed)
					assert.Equal(t, orig.Status, corr.Status)
					assert.Equal(t, orig.Logs, corr.Logs)
				} else {
					assert.Equal(t, orig, corr)
				}
			}
		}
	}

	// Tests use the real reference data, so block numbers and chainIDs are selected for different test cases
	testcases := []testcase{
		// Test case 1: No receipts
		{
			blockNum: 6825767,
			nonces:   []uint64{},
			txTypes:  []uint8{},
			validate: func(original []*types.ReceiptForStorage, corrected []*types.ReceiptForStorage) {
				assert.Empty(t, original)
				assert.Empty(t, corrected)
			},
		},
		// Test case 2: No deposits
		{
			blockNum: 6825767,
			nonces:   []uint64{1, 2, 3},
			txTypes:  []uint8{1, 1, 1},
			validate: func(original []*types.ReceiptForStorage, corrected []*types.ReceiptForStorage) {
				assert.Equal(t, original, corrected)
			},
		},
		// Test case 3: all deposits with no correction
		{
			blockNum: 8835769,
			nonces:   []uint64{78756, 78757, 78758, 78759, 78760, 78761, 78762, 78763, 78764},
			txTypes:  []uint8{126, 126, 126, 126, 126, 126, 126, 126, 126},
			validate: func(original []*types.ReceiptForStorage, corrected []*types.ReceiptForStorage) {
				assert.Equal(t, original, corrected)
			},
		},
		// Test case 4: all deposits with a correction
		{
			blockNum: 8835769,
			nonces:   []uint64{78756, 78757, 78758, 12345, 78760, 78761, 78762, 78763, 78764},
			txTypes:  []uint8{126, 126, 126, 126, 126, 126, 126, 126, 126},
			validate: validateNonceDiff(3),
		},
		// Test case 5: deposits with several corrections and non-deposits
		{
			blockNum: 8835769,
			nonces:   []uint64{0, 1, 2, 78759, 78760, 78761, 6, 78763, 78764, 9, 10, 11},
			txTypes:  []uint8{126, 126, 126, 126, 126, 126, 126, 126, 126, 1, 1, 1},
			// indexes 0, 1, 2, 6 were modified
			// indexes 9, 10, 11 were added too, but they are not user deposits
			validate: validateNonceDiff(0, 1, 2, 6),
		},
	}

	goerliCID := big.NewInt(420)

	rng := rand.New(rand.NewSource(10))
	for _, tc := range testcases {
		// Create original receipts and transactions
		receipts := make([]*types.ReceiptForStorage, len(tc.nonces))
		transactions := make(types.Transactions, len(tc.nonces))
		for i := range tc.nonces {
			rlog := &types.Log{
				Topics: make([]common.Hash, 1),
				Data:   make([]byte, rng.Intn(16)+1),
			}
			rng.Read(rlog.Address[:])
			rng.Read(rlog.Topics[0][:])
			rng.Read(rlog.Data)
			receipts[i] = &types.ReceiptForStorage{
				CumulativeGasUsed: uint64(rng.Intn(1000)),
				Status:            1,
				Logs:              []*types.Log{rlog},
				DepositNonce:      &tc.nonces[i],
			}
			switch tc.txTypes[i] {
			case types.DepositTxType:
				transactions[i] = types.NewTx(&types.DepositTx{})
			case types.AccessListTxType:
				transactions[i] = types.NewTx(&types.AccessListTx{})
			}
		}

		// Encode original receipts to RLP
		originalRLP, err := rlp.EncodeToBytes(receipts)
		assert.NoError(t, err)

		// Call correctReceiptsRLP
		correctedRLP := correctReceipts(originalRLP, transactions, tc.blockNum, goerliCID)

		// Validate the results
		var corrected []*types.ReceiptForStorage
		assert.NoError(t, rlp.DecodeBytes(correctedRLP, &corrected))
		for _, r := range corrected {
			// Nuke Bloom field which isn't available in the original
			r.Bloom = types.Bloom{}
		}
		tc.validate(receipts, corrected)
	}
}
