package types

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
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

// Tests that the effective gas price is correctly derived for different transaction types, in different scenarios.
func TestReceiptEffectiveGasPriceDerivation(t *testing.T) {
	gasPrice := big.NewInt(1000)
	gasFeeCap := big.NewInt(800)
	gasTipCap := big.NewInt(100)
	// Receipt base fee is the base fee encoded in the receipt which will be set post cel2 for CeloDynamicFeeTxV2 types.
	receiptBaseFee := big.NewInt(50)

	t.Run("LegacyTx", func(t *testing.T) {
		testNonDynamic(t, NewTransaction(0, common.Address{}, big.NewInt(0), 0, gasPrice, nil), gasPrice)
	})
	t.Run("AccessListTx", func(t *testing.T) {
		testNonDynamic(t, NewTx(&AccessListTx{GasPrice: gasPrice}), gasPrice)
	})
	t.Run("DynamicFeeTx", func(t *testing.T) {
		tx := NewTx(&DynamicFeeTx{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap})
		testDynamic(t, tx, nil)
	})
	t.Run("BlobTx", func(t *testing.T) {
		tx := NewTx(&BlobTx{GasFeeCap: uint256.MustFromBig(gasFeeCap), GasTipCap: uint256.MustFromBig(gasTipCap)})
		testDynamic(t, tx, nil)
	})
	t.Run("CeloDynamicFeeTx", func(t *testing.T) {
		tx := NewTx(&CeloDynamicFeeTx{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap})
		testDynamic(t, tx, nil)
		tx = NewTx(&CeloDynamicFeeTx{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap, FeeCurrency: &common.Address{}})
		testDynamicWithFeeCurrency(t, tx, nil)
	})
	t.Run("CeloDynamicFeeTxV2", func(t *testing.T) {
		tx := NewTx(&CeloDynamicFeeTxV2{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap})
		testDynamic(t, tx, nil)
		testDynamic(t, tx, receiptBaseFee)
		tx = NewTx(&CeloDynamicFeeTxV2{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap, FeeCurrency: &common.Address{}})
		testDynamicWithFeeCurrency(t, tx, nil)
		testDynamicWithFeeCurrency(t, tx, receiptBaseFee)
	})
	t.Run("CeloDenominatedTx", func(t *testing.T) {
		tx := NewTx(&CeloDenominatedTx{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap})
		testDynamic(t, tx, nil)
		tx = NewTx(&CeloDenominatedTx{GasFeeCap: gasFeeCap, GasTipCap: gasTipCap, FeeCurrency: &common.Address{}})
		testDynamicWithFeeCurrency(t, tx, nil)
	})
}

func testNonDynamic(t *testing.T, tx *Transaction, receiptBaseFee *big.Int) {
	// Non dynamic txs should always have the gas price defined in the tx.
	config := params.TestChainConfig
	config.GingerbreadBlock = big.NewInt(1)
	config.LondonBlock = big.NewInt(3)
	preGingerbreadBlock := uint64(0)
	postGingerbreadBlock := uint64(2)

	receipts := []*Receipt{{BaseFee: receiptBaseFee}}
	txs := []*Transaction{tx}

	// Pre-gingerbread
	err := Receipts(receipts).DeriveFields(config, blockHash, preGingerbreadBlock, blockTime, nil, nil, txs)
	require.NoError(t, err)
	require.Equal(t, tx.GasPrice(), receipts[0].EffectiveGasPrice)

	// Post-gingerbread
	err = Receipts(receipts).DeriveFields(config, blockHash, postGingerbreadBlock, blockTime, baseFee, nil, txs)
	require.NoError(t, err)
	require.Equal(t, tx.GasPrice(), receipts[0].EffectiveGasPrice)
}

// Dynamic txs with no fee currency should have nil for the effective gas price pre-gingerbread and the correct
// effective gas price post-gingerbread, if the receipt base fee is set then the post-gingerbread effective gas price
// should be calculated with that.
func testDynamic(t *testing.T, tx *Transaction, receiptBaseFee *big.Int) {
	config := params.TestChainConfig
	config.GingerbreadBlock = big.NewInt(1)
	config.LondonBlock = big.NewInt(3)
	preGingerbreadBlock := uint64(0)
	postGingerbreadBlock := uint64(2)
	receipts := []*Receipt{{BaseFee: receiptBaseFee}}
	txs := []*Transaction{tx}

	// Pre-gingerbread
	err := Receipts(receipts).DeriveFields(config, blockHash, preGingerbreadBlock, blockTime, nil, nil, txs)
	require.NoError(t, err)
	var nilBigInt *big.Int
	require.Equal(t, nilBigInt, receipts[0].EffectiveGasPrice)

	// Post-gingerbread
	err = Receipts(receipts).DeriveFields(config, blockHash, postGingerbreadBlock, blockTime, baseFee, nil, txs)
	require.NoError(t, err)
	if receiptBaseFee != nil {
		require.Equal(t, tx.inner.effectiveGasPrice(new(big.Int), receiptBaseFee), receipts[0].EffectiveGasPrice)
	} else {
		require.Equal(t, tx.inner.effectiveGasPrice(new(big.Int), baseFee), receipts[0].EffectiveGasPrice)
	}
}

// Dynamic txs with a fee currency set should have nil for the effective gas price pre and post gingerbread, unless
// the receiptBaseFee is set, in which case the post-gingerbread effective gas price should be calculated with the
// receiptBaseFee.
func testDynamicWithFeeCurrency(t *testing.T, tx *Transaction, receiptBaseFee *big.Int) {
	config := params.TestChainConfig
	config.GingerbreadBlock = big.NewInt(1)
	config.LondonBlock = big.NewInt(3)
	preGingerbreadBlock := uint64(0)
	postGingerbreadBlock := uint64(2)
	receipts := []*Receipt{{BaseFee: receiptBaseFee}}
	txs := []*Transaction{tx}

	// Pre-gingerbread
	err := Receipts(receipts).DeriveFields(config, blockHash, preGingerbreadBlock, blockTime, nil, nil, txs)
	require.NoError(t, err)
	var nilBigInt *big.Int
	require.Equal(t, nilBigInt, receipts[0].EffectiveGasPrice)

	// Post-gingerbread
	err = Receipts(receipts).DeriveFields(config, blockHash, postGingerbreadBlock, blockTime, baseFee, nil, txs)
	require.NoError(t, err)
	if receiptBaseFee != nil {
		require.Equal(t, tx.inner.effectiveGasPrice(new(big.Int), receiptBaseFee), receipts[0].EffectiveGasPrice)
	} else if tx.Type() == CeloDenominatedTxType {
		require.Equal(t, tx.inner.effectiveGasPrice(new(big.Int), baseFee), receipts[0].EffectiveGasPrice)
	} else {
		require.Equal(t, nilBigInt, receipts[0].EffectiveGasPrice)
	}
}
