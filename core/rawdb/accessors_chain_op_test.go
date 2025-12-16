package rawdb

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestParseLegacyReceiptRLP(t *testing.T) {
	// Create a gasUsed value greater than a uint64 can represent
	gasUsed := big.NewInt(0)
	gasUsed = gasUsed.SetUint64(math.MaxUint64)
	gasUsed = gasUsed.Add(gasUsed, big.NewInt(math.MaxInt64))
	sanityCheck := (&big.Int{}).SetUint64(gasUsed.Uint64())
	require.NotEqual(t, gasUsed, sanityCheck)
	receipt := types.LegacyOptimismStoredReceiptRLP{
		CumulativeGasUsed: 1,
		Logs: []*types.LogForStorage{
			{Address: common.BytesToAddress([]byte{0x11})},
			{Address: common.BytesToAddress([]byte{0x01, 0x11})},
		},
		L1GasUsed:  gasUsed,
		L1GasPrice: gasUsed,
		L1Fee:      gasUsed,
		FeeScalar:  "6",
	}

	data, err := rlp.EncodeToBytes(receipt)
	require.NoError(t, err)
	var result storedReceiptRLP
	err = rlp.DecodeBytes(data, &result)
	require.NoError(t, err)
	require.Equal(t, receipt.L1GasUsed, result.L1GasUsed)
	require.Equal(t, receipt.L1GasPrice, result.L1GasPrice)
	require.Equal(t, receipt.L1Fee, result.L1Fee)
	require.Equal(t, receipt.FeeScalar, result.FeeScalar)
}
