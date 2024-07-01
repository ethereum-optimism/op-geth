package types

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

func TestCeloDynamicFeeTxReceiptEncodeDecode(t *testing.T) {
	checkEncodeDecodeConsistency(createTypedReceipt(CeloDynamicFeeTxType), t)
}

func TestCeloDynamicFeeTxV2ReceiptEncodeDecode(t *testing.T) {
	t.Run("NoBaseFee", func(t *testing.T) {
		checkEncodeDecodeConsistency(createTypedReceipt(CeloDynamicFeeTxV2Type), t)
	})

	t.Run("WithBaseFee", func(t *testing.T) {
		r := createTypedReceipt(CeloDynamicFeeTxV2Type)
		r.BaseFee = big.NewInt(1000)
		checkEncodeDecodeConsistency(r, t)
	})
}

func createTypedReceipt(receiptType uint8) *Receipt {
	// Note this receipt and logs lack lots of fields, those fields are derived from the
	// block and transaction and so are not part of encoding/decoding.
	r := &Receipt{
		Type:              receiptType,
		PostState:         common.Hash{3}.Bytes(),
		CumulativeGasUsed: 6,
		Logs: []*Log{
			{
				Address: common.BytesToAddress([]byte{0x33}),
				Topics:  []common.Hash{common.HexToHash("dead")},
				Data:    []byte{0x01, 0x02, 0x03},
			},
			{
				Address: common.BytesToAddress([]byte{0x03, 0x33}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x02},
			},
		},
	}
	r.Bloom = CreateBloom(Receipts{r})
	return r
}

// checkEncodeDecodeConsistency checks both RLP and binary encoding/decoding consistency.
func checkEncodeDecodeConsistency(r *Receipt, t *testing.T) {
	checkRLPEncodeDecodeConsistency(r, t)
	checkStorageRLPEncodeDecodeConsistency((*ReceiptForStorage)(r), t)
	checkBinaryEncodeDecodeConsistency(r, t)
}

// checkRLPEncodeDecodeConsistency encodes and decodes the receipt and checks that they are equal.
func checkRLPEncodeDecodeConsistency(r *Receipt, t *testing.T) {
	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, r)
	require.NoError(t, err)

	var r2 Receipt
	err = rlp.Decode(buf, &r2)
	require.NoError(t, err)

	require.EqualValues(t, r, &r2)
}

// checkRLPEncodeDecodeConsistency encodes and decodes the receipt and checks that they are equal.
func checkBinaryEncodeDecodeConsistency(r *Receipt, t *testing.T) {
	bytes, err := r.MarshalBinary()
	require.NoError(t, err)

	r2 := &Receipt{}
	err = r2.UnmarshalBinary(bytes)
	require.NoError(t, err)

	require.EqualValues(t, r, r2)
}

// checkStorageRLPEncodeDecodeConsistency encodes and decodes the receipt and checks that they are equal.
func checkStorageRLPEncodeDecodeConsistency(r *ReceiptForStorage, t *testing.T) {
	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, r)
	require.NoError(t, err)

	// Stored receipts do not encode the type, (although they do require it to be set during encoding)
	// since it is derived from the associated transaction. So for the sake of the comparison we set it
	// to 0 and restore it after the comparison.
	receiptType := r.Type
	defer func() { r.Type = receiptType }()
	r.Type = 0

	var r2 ReceiptForStorage
	err = rlp.Decode(buf, &r2)
	require.NoError(t, err)

	require.EqualValues(t, r, &r2)
}
