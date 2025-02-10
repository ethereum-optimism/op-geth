package types

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/assert"
)

// mockOldBeforeGingerbreadHeader is same as BeforeGingerbreadHeader
// but doesn't implement EncodeRLP and DecodeRLP
type mockOldBeforeGingerbreadHeader BeforeGingerbreadHeader

// mockOldAfterGingerbreadHeader is also same as AfterGingerbreadHeader
type mockOldAfterGingerbreadHeader Header

var (
	mockBeforeGingerbreadHeader = &BeforeGingerbreadHeader{
		ParentHash:  common.HexToHash("0x112233445566778899001122334455667788990011223344556677889900aabb"),
		Coinbase:    common.HexToAddress("0x8888f1f195afa192cfee860698584c030f4c9db1"),
		Root:        EmptyRootHash,
		TxHash:      EmptyTxsHash,
		ReceiptHash: EmptyReceiptsHash,
		Bloom:       Bloom(common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")),
		Number:      math.BigPow(2, 9),
		GasUsed:     1476322,
		Time:        9876543,
		Extra:       []byte("test before gingerbread header extra"),
	}

	mockWithdrawalHash         = common.HexToHash("0x4585754a71d14791295bc094dc53eb0b32f21d92e58350a4140163a047b854a7")
	mockExcessBlobGas          = uint64(123456789)
	mockBlobGasUsed            = uint64(12345678)
	mockParentBeaconRoot       = common.HexToHash("0x9229c626ebd6328b3ddc7fe8636f2fd9a344f4c02e2e281f59a3b7e4e46833e5")
	mockAfterGingerbreadHeader = &AfterGingerbreadHeader{
		ParentHash:       mockBeforeGingerbreadHeader.ParentHash,
		UncleHash:        EmptyUncleHash,
		Coinbase:         mockBeforeGingerbreadHeader.Coinbase,
		Root:             EmptyRootHash,
		TxHash:           EmptyTxsHash,
		ReceiptHash:      EmptyReceiptsHash,
		Bloom:            mockBeforeGingerbreadHeader.Bloom,
		Difficulty:       big.NewInt(17179869184),
		Number:           mockBeforeGingerbreadHeader.Number,
		GasLimit:         12345678,
		GasUsed:          mockBeforeGingerbreadHeader.GasUsed,
		Time:             mockBeforeGingerbreadHeader.Time,
		Extra:            []byte("test after gingerbread header extra"),
		MixDigest:        common.HexToHash("0x036a0a7a3611ecd974ef274e603ceab81246fb50dc350519b9f47589e8fe3014"),
		Nonce:            EncodeNonce(12345),
		BaseFee:          math.BigPow(10, 8),
		WithdrawalsHash:  &mockWithdrawalHash,
		ExcessBlobGas:    &mockExcessBlobGas,
		BlobGasUsed:      &mockBlobGasUsed,
		ParentBeaconRoot: &mockParentBeaconRoot,
	}
)

func ToMockOldBeforeGingerbreadHeader(h *BeforeGingerbreadHeader) *mockOldBeforeGingerbreadHeader {
	return &mockOldBeforeGingerbreadHeader{
		ParentHash:  h.ParentHash,
		Coinbase:    h.Coinbase,
		Root:        h.Root,
		TxHash:      h.TxHash,
		ReceiptHash: h.ReceiptHash,
		Bloom:       h.Bloom,
		Number:      h.Number,
		GasUsed:     h.GasUsed,
		Time:        h.Time,
		Extra:       h.Extra,
	}
}

func BeforeGingerbreadHeaderToHeader(h *BeforeGingerbreadHeader) *Header {
	return &Header{
		ParentHash:     h.ParentHash,
		Coinbase:       h.Coinbase,
		Root:           h.Root,
		TxHash:         h.TxHash,
		ReceiptHash:    h.ReceiptHash,
		Bloom:          h.Bloom,
		Number:         h.Number,
		GasUsed:        h.GasUsed,
		Time:           h.Time,
		Extra:          h.Extra,
		Difficulty:     new(big.Int),
		PreGingerbread: true,
	}
}

func TestRLPDecodeHeaderCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		oldHeader interface{}
		newHeader *Header
	}{
		{
			name:      "BeforeGingerbreadHeader",
			oldHeader: ToMockOldBeforeGingerbreadHeader(mockBeforeGingerbreadHeader),
			newHeader: BeforeGingerbreadHeaderToHeader(mockBeforeGingerbreadHeader),
		},
		{
			name:      "AfterGingerbreadHeader",
			oldHeader: (*mockOldAfterGingerbreadHeader)(mockAfterGingerbreadHeader),
			newHeader: (*Header)(mockAfterGingerbreadHeader),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewBuffer([]byte{})

			// encode by reflection style
			err := rlp.Encode(r, test.oldHeader)
			assert.NoError(t, err, "failed RLP encode")

			// decode by generated code
			decodedHeader := &Header{}
			rlp.DecodeBytes(r.Bytes(), decodedHeader)
			assert.NoError(t, err, "failed RLP decode")

			assert.Equal(t, test.newHeader, decodedHeader)
		})
	}
}

func TestRlpEncodeHeaderCompatibility(t *testing.T) {
	tests := []struct {
		name      string
		oldHeader interface{} // header type which doesn't implement EncodeRLP
		newHeader *Header     // header type which implements EncodeRLP
	}{
		{
			name:      "BeforeGingerbreadHeader",
			oldHeader: ToMockOldBeforeGingerbreadHeader(mockBeforeGingerbreadHeader),
			newHeader: BeforeGingerbreadHeaderToHeader(mockBeforeGingerbreadHeader),
		},
		{
			name:      "AfterGingerbreadHeader",
			oldHeader: (*mockOldAfterGingerbreadHeader)(mockAfterGingerbreadHeader),
			newHeader: (*Header)(mockAfterGingerbreadHeader),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			r := bytes.NewBuffer([]byte{})

			// old RLP encoding
			err := rlp.Encode(r, test.oldHeader)
			assert.NoError(t, err, "failed RLP encode by reflection style")
			oldEncodedData := r.Bytes()

			r.Reset()

			// new RLP encoding
			err = rlp.Encode(r, test.newHeader)
			assert.NoError(t, err, "failed RLP encode by generated code")
			newEncodedData := r.Bytes()

			assert.Equal(t, oldEncodedData, newEncodedData)
		})
	}
}
