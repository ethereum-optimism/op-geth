package downloader

import (
	"testing"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// makeCorrection is a helper function to create a slice of receipts and a slice of corrected receipts
func makeCorrection(bn uint64, cid uint64, ns []uint64, ty []uint8) (types.Receipts, types.Receipts) {
	receipts := make(types.Receipts, len(ns))
	correctedReceipts := make(types.Receipts, len(ns))
	transactions := make(types.Transactions, len(ns))
	for i := range ns {
		receipts[i] = &types.Receipt{Type: ty[i], DepositNonce: &ns[i]}
		correctedReceipts[i] = &types.Receipt{Type: ty[i], DepositNonce: &ns[i]}
		transactions[i] = types.NewTx(&types.DepositTx{})
	}

	correctedReceipts = correctReceipts(correctedReceipts, transactions, bn, cid)

	return receipts, correctedReceipts
}

func TestCorrectReceipts(t *testing.T) {
	type testcase struct {
		blockNum uint64
		chainID  uint64
		nonces   []uint64
		txTypes  []uint8
		validate func(types.Receipts, types.Receipts)
	}

	// Tests use the real reference data, so block numbers and chainIDs are selected for different test cases
	testcases := []testcase{
		// Test case 1: No receipts
		{
			blockNum: 6825767,
			chainID:  420,
			nonces:   []uint64{},
			txTypes:  []uint8{},
			validate: func(receipts types.Receipts, correctedReceipts types.Receipts) {
				assert.Empty(t, correctedReceipts)
			},
		},
		// Test case 2: No deposits
		{
			blockNum: 6825767,
			chainID:  420,
			nonces:   []uint64{1, 2, 3},
			txTypes:  []uint8{1, 1, 1},
			validate: func(receipts types.Receipts, correctedReceipts types.Receipts) {
				assert.Equal(t, receipts, correctedReceipts)
			},
		},
		// Test case 3: all deposits with no correction
		{
			blockNum: 8835769,
			chainID:  420,
			nonces:   []uint64{78756, 78757, 78758, 78759, 78760, 78761, 78762, 78763, 78764},
			txTypes:  []uint8{126, 126, 126, 126, 126, 126, 126, 126, 126},
			validate: func(receipts types.Receipts, correctedReceipts types.Receipts) {
				assert.Equal(t, receipts, correctedReceipts)
			},
		},
		// Test case 4: all deposits with a correction
		{
			blockNum: 8835769,
			chainID:  420,
			nonces:   []uint64{78756, 78757, 78758, 12345, 78760, 78761, 78762, 78763, 78764},
			txTypes:  []uint8{126, 126, 126, 126, 126, 126, 126, 126, 126},
			validate: func(receipts types.Receipts, correctedReceipts types.Receipts) {
				assert.NotEqual(t, receipts[3], correctedReceipts[3])
				for i := range receipts {
					if i != 3 {
						assert.Equal(t, receipts[i], correctedReceipts[i])
					}
				}
			},
		},
		// Test case 5: deposits with several corrections and non-deposits
		{
			blockNum: 8835769,
			chainID:  420,
			nonces:   []uint64{0, 1, 2, 78759, 78760, 78761, 6, 78763, 78764, 9, 10, 11},
			txTypes:  []uint8{126, 126, 126, 126, 126, 126, 126, 126, 126, 1, 1, 1},
			validate: func(receipts types.Receipts, correctedReceipts types.Receipts) {
				// indexes 0, 1, 2, 6 were modified
				// indexes 9, 10, 11 were added too, but they are not user deposits
				assert.NotEqual(t, receipts[0], correctedReceipts[0])
				assert.NotEqual(t, receipts[1], correctedReceipts[1])
				assert.NotEqual(t, receipts[2], correctedReceipts[2])
				assert.NotEqual(t, receipts[6], correctedReceipts[6])
				for i := range receipts {
					if i != 0 && i != 1 && i != 2 && i != 6 {
						assert.Equal(t, receipts[i], correctedReceipts[i])
					}
				}
			},
		},
	}

	for _, tc := range testcases {
		receipts, correctedReceipts := makeCorrection(tc.blockNum, tc.chainID, tc.nonces, tc.txTypes)
		tc.validate(receipts, correctedReceipts)
	}
}
