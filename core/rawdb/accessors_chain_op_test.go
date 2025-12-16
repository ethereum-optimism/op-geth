package rawdb

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
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

func TestDecodeRawLegacyReceipt(t *testing.T) {
	// On op-mainnet:
	// cast rpc debug_dbAncient -w '["receipts",10]'
	rawDBReceiptHex := "0xf90104f90101018301c60df8faf89c9425e1c58040f27ecf20bbd4ca83a09290326896b3f884a092e98423f8adac6e64d0608e519fd1cefb861498385c6dee70d58fc926ddc68ca00000000000000000000000000000000000000000000000000000006ea50c77c9a0000000000000000000000000000000000000000000000000000000000000c026a00000000000000000000000002539ffc9ded82926a5aaee065e800c7d1de0245480f85a9425e1c58040f27ecf20bbd4ca83a09290326896b3f842a0fe25c73e3b9089fac37d55c4c7efcba6f04af04cebd2fc4d6d7dbb07e1e5234fa000000000000000000000000000000000000000000000005b3636a3506234800080"

	rawDBReceiptBytes, err := hexutil.Decode(rawDBReceiptHex)
	require.NoError(t, err)

	sr := new([]receiptLogs)
	err = rlp.DecodeBytes(rawDBReceiptBytes, sr)
	require.NoError(t, err)
	require.Len(t, *sr, 1)
	require.Equal(t, (*sr)[0].Logs[0].Address, common.HexToAddress("0x25e1c58040f27ecf20bbd4ca83a09290326896b3"))
}
