package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
)

var (
	bedrockGenesisTestConfig = func() *params.ChainConfig {
		conf := *params.AllCliqueProtocolChanges // copy the config
		conf.Clique = nil
		conf.BedrockBlock = big.NewInt(0)
		conf.Optimism = &params.OptimismConfig{EIP1559Elasticity: 50, EIP1559Denominator: 10}
		return &conf
	}()
	ecotoneTestConfig = func() *params.ChainConfig {
		conf := *bedrockGenesisTestConfig // copy the config
		time := uint64(0)
		conf.EcotoneTime = &time
		return &conf
	}()
	isthmusTestConfig = func() *params.ChainConfig {
		conf := *ecotoneTestConfig // copy the config
		time := uint64(0)
		conf.FjordTime = &time
		conf.GraniteTime = &time
		conf.HoloceneTime = &time
		conf.IsthmusTime = &time
		return &conf
	}()
	jovianTestConfig = func() *params.ChainConfig {
		conf := *isthmusTestConfig // copy the config
		time := uint64(0)
		conf.JovianTime = &time
		return &conf
	}()

	depositReceiptNoNonce = &Receipt{
		Status:            ReceiptStatusFailed,
		CumulativeGasUsed: 1,
		Logs: []*Log{
			{
				Address: common.BytesToAddress([]byte{0x11}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
			{
				Address: common.BytesToAddress([]byte{0x01, 0x11}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
		},
		Type: DepositTxType,
	}
	nonce                   = uint64(1234)
	depositReceiptWithNonce = &Receipt{
		Status:                ReceiptStatusFailed,
		CumulativeGasUsed:     1,
		DepositNonce:          &nonce,
		DepositReceiptVersion: nil,
		Logs: []*Log{
			{
				Address: common.BytesToAddress([]byte{0x11}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
			{
				Address: common.BytesToAddress([]byte{0x01, 0x11}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
		},
		Type: DepositTxType,
	}
	version                           = CanyonDepositReceiptVersion
	depositReceiptWithNonceAndVersion = &Receipt{
		Status:                ReceiptStatusFailed,
		CumulativeGasUsed:     1,
		DepositNonce:          &nonce,
		DepositReceiptVersion: &version,
		Logs: []*Log{
			{
				Address: common.BytesToAddress([]byte{0x11}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
			{
				Address: common.BytesToAddress([]byte{0x01, 0x11}),
				Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
				Data:    []byte{0x01, 0x00, 0xff},
			},
		},
		Type: DepositTxType,
	}

	daFootprintGasScalar = uint16(400)
)

func clearComputedFieldsOnOPStackReceipts(receipts []*Receipt) []*Receipt {
	receipts = clearComputedFieldsOnReceipts(receipts)
	for _, receipt := range receipts {
		receipt.L1GasPrice = nil
		receipt.L1BlobBaseFee = nil
		receipt.L1GasUsed = nil
		receipt.L1Fee = nil
		receipt.FeeScalar = nil
		receipt.L1BaseFeeScalar = nil
		receipt.L1BlobBaseFeeScalar = nil
		receipt.OperatorFeeScalar = nil
		receipt.OperatorFeeConstant = nil
		receipt.DAFootprintGasScalar = nil
	}
	return receipts
}

func getOptimismTxReceipts(l1AttributesPayload []byte, l1GasPrice, l1GasUsed, l1Fee *big.Int, feeScalar *big.Float) ([]*Transaction, []*Receipt) {
	// Create a few transactions to have receipts for
	txs := Transactions{
		NewTx(&DepositTx{
			To:    nil, // contract creation
			Value: big.NewInt(6),
			Gas:   50,
			Data:  l1AttributesPayload,
		}),
		emptyTx,
	}

	// Create the corresponding receipts
	receipts := Receipts{
		&Receipt{
			Type:              DepositTxType,
			PostState:         common.Hash{5}.Bytes(),
			CumulativeGasUsed: 50 + 15,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x33}),
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[0].Hash(),
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
				},
				{
					Address: common.BytesToAddress([]byte{0x03, 0x33}),
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[0].Hash(),
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       1,
				},
			},
			TxHash:            txs[0].Hash(),
			ContractAddress:   common.HexToAddress("0x3bb898b4bbe24f68a4e9be46cfe72d1787fd74f4"),
			GasUsed:           65,
			EffectiveGasPrice: big.NewInt(0),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  0,
			DepositNonce:      &depNonce1,
		},
		&Receipt{
			Type:              LegacyTxType,
			EffectiveGasPrice: big.NewInt(0),
			PostState:         common.Hash{4}.Bytes(),
			CumulativeGasUsed: 10,
			Logs:              []*Log{},
			// derived fields:
			TxHash:           txs[1].Hash(),
			GasUsed:          18446744073709551561,
			BlockHash:        blockHash,
			BlockNumber:      blockNumber,
			TransactionIndex: 1,
			L1GasPrice:       l1GasPrice,
			L1GasUsed:        l1GasUsed,
			L1Fee:            l1Fee,
			FeeScalar:        feeScalar,
		},
	}
	for _, receipt := range receipts {
		receipt.Bloom = CreateBloom(receipt)
	}

	return txs, receipts
}

func getOptimismEcotoneTxReceipts(l1AttributesPayload []byte, l1GasPrice, l1BlobBaseFee, l1GasUsed, l1Fee *big.Int, baseFeeScalar, blobBaseFeeScalar *uint64) ([]*Transaction, []*Receipt) {
	txs, receipts := getOptimismTxReceipts(l1AttributesPayload, l1GasPrice, l1GasUsed, l1Fee, nil)
	receipts[1].L1BlobBaseFee = l1BlobBaseFee
	receipts[1].L1BaseFeeScalar = baseFeeScalar
	receipts[1].L1BlobBaseFeeScalar = blobBaseFeeScalar
	return txs, receipts
}

func getOptimismIsthmusTxReceipts(l1AttributesPayload []byte, l1GasPrice, l1BlobBaseFee, l1GasUsed, l1Fee *big.Int, baseFeeScalar, blobBaseFeeScalar, operatorFeeScalar, operatorFeeConstant *uint64) ([]*Transaction, []*Receipt) {
	txs, receipts := getOptimismEcotoneTxReceipts(l1AttributesPayload, l1GasPrice, l1BlobBaseFee, l1GasUsed, l1Fee, baseFeeScalar, blobBaseFeeScalar)
	receipts[1].OperatorFeeScalar = operatorFeeScalar
	receipts[1].OperatorFeeConstant = operatorFeeConstant
	return txs, receipts
}

func getOptimismJovianTxReceipts(l1AttributesPayload []byte, l1GasPrice, l1BlobBaseFee, l1GasUsed, l1Fee *big.Int, baseFeeScalar, blobBaseFeeScalar, operatorFeeScalar, operatorFeeConstant, daFootprintGasScalar *uint64) ([]*Transaction, []*Receipt) {
	txs, receipts := getOptimismIsthmusTxReceipts(l1AttributesPayload, l1GasPrice, l1BlobBaseFee, l1GasUsed, l1Fee, baseFeeScalar, blobBaseFeeScalar, operatorFeeScalar, operatorFeeConstant)
	receipts[1].DAFootprintGasScalar = daFootprintGasScalar
	if daFootprintGasScalar != nil {
		receipts[1].BlobGasUsed = *daFootprintGasScalar * txs[1].RollupCostData().EstimatedDASize().Uint64()
	}
	return txs, receipts
}

func TestDeriveOptimismBedrockTxReceipts(t *testing.T) {
	// Bedrock style l1 attributes with L1Scalar=7_000_000 (becomes 7 after division), L1Overhead=50, L1BaseFee=1000*1e6
	payload := common.Hex2Bytes("015d8eb900000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2000000000000000000000000000000000000000000000000000000003b9aca0000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2000000000000000000000000000000000000000000000000000000000000003200000000000000000000000000000000000000000000000000000000006acfc0015d8eb900000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2000000000000000000000000000000000000000000000000000000003b9aca0000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2000000000000000000000000000000000000000000000000000000000000003200000000000000000000000000000000000000000000000000000000006acfc0")
	// the parameters we use below are defined in rollup_test.go
	l1GasPrice := baseFee
	l1GasUsed := bedrockGas
	feeScalar := big.NewFloat(float64(scalar.Uint64() / 1e6))
	l1Fee := bedrockFee
	txs, receipts := getOptimismTxReceipts(payload, l1GasPrice, l1GasUsed, l1Fee, feeScalar)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnOPStackReceipts(receipts)
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.NoError(t, err)
	checkBedrockReceipts(t, receipts, derivedReceipts)

	// Should get same result with the Ecotone config because it will assume this is "first ecotone block"
	// if it sees the bedrock style L1 attributes.
	err = Receipts(derivedReceipts).DeriveFields(ecotoneTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.NoError(t, err)
	checkBedrockReceipts(t, receipts, derivedReceipts)
}

func TestDeriveOptimismEcotoneTxReceipts(t *testing.T) {
	// Ecotone style l1 attributes with baseFeeScalar=2, blobBaseFeeScalar=3, baseFee=1000*1e6, blobBaseFee=10*1e6
	payload := common.Hex2Bytes("440a5e20000000020000000300000000000004d200000000000004d200000000000004d2000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2")
	// the parameters we use below are defined in rollup_test.go
	baseFeeScalarUint64 := baseFeeScalar.Uint64()
	blobBaseFeeScalarUint64 := blobBaseFeeScalar.Uint64()
	txs, receipts := getOptimismEcotoneTxReceipts(payload, baseFee, blobBaseFee, ecotoneGas, ecotoneFee, &baseFeeScalarUint64, &blobBaseFeeScalarUint64)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnOPStackReceipts(receipts)
	// Should error out if we try to process this with a pre-Ecotone config
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.Error(t, err)

	err = Receipts(derivedReceipts).DeriveFields(ecotoneTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.NoError(t, err)
	diffReceipts(t, receipts, derivedReceipts)
}

func TestDeriveOptimismIsthmusTxReceipts(t *testing.T) {
	// Isthmus style l1 attributes with baseFeeScalar=2, blobBaseFeeScalar=3, baseFee=1000*1e6, blobBaseFee=10*1e6, operatorFeeScalar=1439103868, operatorFeeConstant=1256417826609331460
	payload := common.Hex2Bytes("098999be000000020000000300000000000004d200000000000004d200000000000004d2000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d255c6fb7c116fb15b44847d04")
	// the parameters we use below are defined in rollup_test.go
	baseFeeScalarUint64 := baseFeeScalar.Uint64()
	blobBaseFeeScalarUint64 := blobBaseFeeScalar.Uint64()
	operatorFeeScalarUint64 := operatorFeeScalar.Uint64()
	operatorFeeConstantUint64 := operatorFeeConstant.Uint64()
	txs, receipts := getOptimismIsthmusTxReceipts(payload, baseFee, blobBaseFee, minimumFjordGas, fjordFee, &baseFeeScalarUint64, &blobBaseFeeScalarUint64, &operatorFeeScalarUint64, &operatorFeeConstantUint64)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnOPStackReceipts(receipts)
	// Should error out if we try to process this with a pre-Isthmus config
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.Error(t, err)

	err = Receipts(derivedReceipts).DeriveFields(isthmusTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.NoError(t, err)
	diffReceipts(t, receipts, derivedReceipts)
}

func TestDeriveOptimismIsthmusTxReceiptsNoOperatorFee(t *testing.T) {
	// Isthmus style l1 attributes with baseFeeScalar=2, blobBaseFeeScalar=3, baseFee=1000*1e6, blobBaseFee=10*1e6, operatorFeeScalar=0, operatorFeeConstant=0
	payload := common.Hex2Bytes("098999be000000020000000300000000000004d200000000000004d200000000000004d2000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2000000000000000000000000")
	// the parameters we use below are defined in rollup_test.go
	baseFeeScalarUint64 := baseFeeScalar.Uint64()
	blobBaseFeeScalarUint64 := blobBaseFeeScalar.Uint64()
	txs, receipts := getOptimismIsthmusTxReceipts(payload, baseFee, blobBaseFee, minimumFjordGas, fjordFee, &baseFeeScalarUint64, &blobBaseFeeScalarUint64, nil, nil)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnOPStackReceipts(receipts)
	// Should error out if we try to process this with a pre-Isthmus config
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.Error(t, err)

	err = Receipts(derivedReceipts).DeriveFields(isthmusTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.NoError(t, err)
	diffReceipts(t, receipts, derivedReceipts)
}

func TestDeriveOptimismJovianTxReceipts(t *testing.T) {
	// Jovian style l1 attributes with baseFeeScalar=2, blobBaseFeeScalar=3, baseFee=1000*1e6, blobBaseFee=10*1e6, operatorFeeScalar=1439103868, operatorFeeConstant=1256417826609331460, daFootprintGasScalar=400
	payload := common.Hex2Bytes("3db6be2b000000020000000300000000000004d200000000000004d200000000000004d2000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d255c6fb7c116fb15b44847d040190")
	// the parameters we use below are defined in rollup_test.go
	baseFeeScalarUint64 := baseFeeScalar.Uint64()
	blobBaseFeeScalarUint64 := blobBaseFeeScalar.Uint64()
	operatorFeeScalarUint64 := operatorFeeScalar.Uint64()
	operatorFeeConstantUint64 := operatorFeeConstant.Uint64()
	daFootprintGasScalarUint64 := uint64(daFootprintGasScalar)
	txs, receipts := getOptimismJovianTxReceipts(payload, baseFee, blobBaseFee, minimumFjordGas, fjordFee, &baseFeeScalarUint64, &blobBaseFeeScalarUint64, &operatorFeeScalarUint64, &operatorFeeConstantUint64, &daFootprintGasScalarUint64)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnOPStackReceipts(receipts)
	// Should error out if we try to process this with a pre-Jovian config
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.Error(t, err)

	err = Receipts(derivedReceipts).DeriveFields(jovianTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	require.NoError(t, err)
	diffReceipts(t, receipts, derivedReceipts)
}

func diffReceipts(t *testing.T, receipts, derivedReceipts []*Receipt) {
	// Check diff of receipts against derivedReceipts.
	r1, err := json.MarshalIndent(receipts, "", "  ")
	if err != nil {
		t.Fatal("error marshaling input receipts:", err)
	}
	r2, err := json.MarshalIndent(derivedReceipts, "", "  ")
	if err != nil {
		t.Fatal("error marshaling derived receipts:", err)
	}
	d := diff.Diff(string(r1), string(r2))
	if d != "" {
		t.Fatal("receipts differ:", d)
	}
}

func checkBedrockReceipts(t *testing.T, receipts, derivedReceipts []*Receipt) {
	diffReceipts(t, receipts, derivedReceipts)

	// Check that we preserved the invariant: l1Fee = l1GasPrice * l1GasUsed * l1FeeScalar
	// but with more difficult int math...
	l2Rcpt := derivedReceipts[1]
	l1GasCost := new(big.Int).Mul(l2Rcpt.L1GasPrice, l2Rcpt.L1GasUsed)
	l1Fee := new(big.Float).Mul(new(big.Float).SetInt(l1GasCost), l2Rcpt.FeeScalar)
	require.Equal(t, new(big.Float).SetInt(l2Rcpt.L1Fee), l1Fee)
}

func TestBedrockDepositReceiptUnchanged(t *testing.T) {
	expectedRlp := common.FromHex("7EF90156A003000000000000000000000000000000000000000000000000000000000000000AB9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000F0D7940000000000000000000000000000000000000033C001D7940000000000000000000000000000000000000333C002")
	// Deposit receipt with no nonce
	receipt := &Receipt{
		Type:              DepositTxType,
		PostState:         common.Hash{3}.Bytes(),
		CumulativeGasUsed: 10,
		Logs: []*Log{
			{Address: common.BytesToAddress([]byte{0x33}), Data: []byte{1}, Topics: []common.Hash{}},
			{Address: common.BytesToAddress([]byte{0x03, 0x33}), Data: []byte{2}, Topics: []common.Hash{}},
		},
		TxHash:          common.Hash{},
		ContractAddress: common.BytesToAddress([]byte{0x03, 0x33, 0x33}),
		GasUsed:         4,
	}

	rlp, err := receipt.MarshalBinary()
	require.NoError(t, err)
	require.Equal(t, expectedRlp, rlp)

	// Consensus values should be unchanged after reparsing
	parsed := new(Receipt)
	err = parsed.UnmarshalBinary(rlp)
	require.NoError(t, err)
	require.Equal(t, receipt.Status, parsed.Status)
	require.Equal(t, receipt.CumulativeGasUsed, parsed.CumulativeGasUsed)
	require.Equal(t, receipt.Bloom, parsed.Bloom)
	require.EqualValues(t, receipt.Logs, parsed.Logs)
	// And still shouldn't have a nonce
	require.Nil(t, parsed.DepositNonce)
	// ..or a deposit nonce
	require.Nil(t, parsed.DepositReceiptVersion)
}

// Regolith introduced an inconsistency in behavior between EncodeIndex and MarshalBinary for a
// deposit transaction receipt. TestReceiptEncodeIndexBugIsEnshrined makes sure this difference is
// preserved for backwards compatibility purposes, but also that there is no discrepancy for the
// post-Canyon encoding.
func TestReceiptEncodeIndexBugIsEnshrined(t *testing.T) {
	// Check that a post-Regolith, pre-Canyon receipt produces the expected difference between
	// EncodeIndex and MarshalBinary.
	buf := new(bytes.Buffer)
	receipts := Receipts{depositReceiptWithNonce}
	receipts.EncodeIndex(0, buf)
	indexBytes := buf.Bytes()

	regularBytes, _ := receipts[0].MarshalBinary()

	require.NotEqual(t, indexBytes, regularBytes)

	// Confirm the buggy encoding is as expected, which means it should encode as if it had no
	// nonce specified (like that of a non-deposit receipt, whose encoding would differ only in the
	// type byte).
	buf.Reset()
	tempReceipt := *depositReceiptWithNonce
	tempReceipt.Type = eip1559Receipt.Type
	buggyBytes, _ := tempReceipt.MarshalBinary()

	require.Equal(t, indexBytes[1:], buggyBytes[1:])

	// check that the post-Canyon encoding has no differences between EncodeIndex and
	// MarshalBinary.
	buf.Reset()
	receipts = Receipts{depositReceiptWithNonceAndVersion}
	receipts.EncodeIndex(0, buf)
	indexBytes = buf.Bytes()

	regularBytes, _ = receipts[0].MarshalBinary()

	require.Equal(t, indexBytes, regularBytes)

	// Check that bumping the nonce post-canyon changes the hash
	bumpedReceipt := *depositReceiptWithNonceAndVersion
	bumpedNonce := nonce + 1
	bumpedReceipt.DepositNonce = &bumpedNonce
	bumpedBytes, _ := bumpedReceipt.MarshalBinary()
	require.NotEqual(t, regularBytes, bumpedBytes)
}

func TestRoundTripReceipt(t *testing.T) {
	tests := []struct {
		name string
		rcpt *Receipt
	}{
		{name: "Legacy", rcpt: legacyReceipt},
		{name: "AccessList", rcpt: accessListReceipt},
		{name: "EIP1559", rcpt: eip1559Receipt},
		{name: "DepositNoNonce", rcpt: depositReceiptNoNonce},
		{name: "DepositWithNonce", rcpt: depositReceiptWithNonce},
		{name: "DepositWithNonceAndVersion", rcpt: depositReceiptWithNonceAndVersion},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := test.rcpt.MarshalBinary()
			require.NoError(t, err)

			d := &Receipt{}
			err = d.UnmarshalBinary(data)
			require.NoError(t, err)
			require.Equal(t, test.rcpt, d)
			require.Equal(t, test.rcpt.DepositNonce, d.DepositNonce)
			require.Equal(t, test.rcpt.DepositReceiptVersion, d.DepositReceiptVersion)
		})

		t.Run(fmt.Sprintf("%sRejectExtraData", test.name), func(t *testing.T) {
			data, err := test.rcpt.MarshalBinary()
			require.NoError(t, err)
			data = append(data, 1, 2, 3, 4)
			d := &Receipt{}
			err = d.UnmarshalBinary(data)
			require.Error(t, err)
		})
	}
}

func TestRoundTripReceiptForStorage(t *testing.T) {
	tests := []struct {
		name string
		rcpt *Receipt
	}{
		{name: "Legacy", rcpt: legacyReceipt},
		{name: "AccessList", rcpt: accessListReceipt},
		{name: "EIP1559", rcpt: eip1559Receipt},
		{name: "DepositNoNonce", rcpt: depositReceiptNoNonce},
		{name: "DepositWithNonce", rcpt: depositReceiptWithNonce},
		{name: "DepositWithNonceAndVersion", rcpt: depositReceiptWithNonceAndVersion},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := rlp.EncodeToBytes((*ReceiptForStorage)(test.rcpt))
			require.NoError(t, err)

			d := &ReceiptForStorage{}
			err = rlp.DecodeBytes(data, d)
			require.NoError(t, err)
			// Only check the stored fields - the others are derived later
			require.Equal(t, test.rcpt.Status, d.Status)
			require.Equal(t, test.rcpt.CumulativeGasUsed, d.CumulativeGasUsed)
			require.Equal(t, test.rcpt.Logs, d.Logs)
			require.Equal(t, test.rcpt.DepositNonce, d.DepositNonce)
			require.Equal(t, test.rcpt.DepositReceiptVersion, d.DepositReceiptVersion)
		})
	}
}
