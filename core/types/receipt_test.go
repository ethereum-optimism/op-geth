// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/holiman/uint256"
	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
)

var (
	bedrockGenesisTestConfig = func() *params.ChainConfig {
		conf := *params.AllCliqueProtocolChanges // copy the config
		conf.Clique = nil
		conf.TerminalTotalDifficultyPassed = true
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

	legacyReceipt = &Receipt{
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
	}
	accessListReceipt = &Receipt{
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
		Type: AccessListTxType,
	}
	eip1559Receipt = &Receipt{
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
		Type: DynamicFeeTxType,
	}
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

	// Create a few transactions to have receipts for
	to2 = common.HexToAddress("0x2")
	to3 = common.HexToAddress("0x3")
	to4 = common.HexToAddress("0x4")
	to5 = common.HexToAddress("0x5")
	to6 = common.HexToAddress("0x6")
	to7 = common.HexToAddress("0x7")
	txs = Transactions{
		NewTx(&LegacyTx{
			Nonce:    1,
			Value:    big.NewInt(1),
			Gas:      1,
			GasPrice: big.NewInt(11),
		}),
		NewTx(&LegacyTx{
			To:       &to2,
			Nonce:    2,
			Value:    big.NewInt(2),
			Gas:      2,
			GasPrice: big.NewInt(22),
		}),
		NewTx(&AccessListTx{
			To:       &to3,
			Nonce:    3,
			Value:    big.NewInt(3),
			Gas:      3,
			GasPrice: big.NewInt(33),
		}),
		// EIP-1559 transactions.
		NewTx(&DynamicFeeTx{
			To:        &to4,
			Nonce:     4,
			Value:     big.NewInt(4),
			Gas:       4,
			GasTipCap: big.NewInt(44),
			GasFeeCap: big.NewInt(1044),
		}),
		NewTx(&DynamicFeeTx{
			To:        &to5,
			Nonce:     5,
			Value:     big.NewInt(5),
			Gas:       5,
			GasTipCap: big.NewInt(55),
			GasFeeCap: big.NewInt(1055),
		}),
		// EIP-4844 transactions.
		NewTx(&BlobTx{
			To:         to6,
			Nonce:      6,
			Value:      uint256.NewInt(6),
			Gas:        6,
			GasTipCap:  uint256.NewInt(66),
			GasFeeCap:  uint256.NewInt(1066),
			BlobFeeCap: uint256.NewInt(100066),
			BlobHashes: []common.Hash{{}},
		}),
		NewTx(&BlobTx{
			To:         to7,
			Nonce:      7,
			Value:      uint256.NewInt(7),
			Gas:        7,
			GasTipCap:  uint256.NewInt(77),
			GasFeeCap:  uint256.NewInt(1077),
			BlobFeeCap: uint256.NewInt(100077),
			BlobHashes: []common.Hash{{}, {}, {}},
		}),
		NewTx(&DepositTx{
			To:    nil, // contract creation
			Value: big.NewInt(6),
			Gas:   50,
		}),
		NewTx(&DepositTx{
			To:    nil, // contract creation
			Value: big.NewInt(6),
			Gas:   60,
		}),
	}
	depNonce1                   = uint64(7)
	depNonce2                   = uint64(8)
	blockNumber                 = big.NewInt(1)
	blockTime                   = uint64(2)
	blockHash                   = common.BytesToHash([]byte{0x03, 0x14})
	canyonDepositReceiptVersion = CanyonDepositReceiptVersion

	// Create the corresponding receipts
	receipts = Receipts{
		&Receipt{
			Status:            ReceiptStatusFailed,
			CumulativeGasUsed: 1,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x11}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[0].Hash(),
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
				},
				{
					Address: common.BytesToAddress([]byte{0x01, 0x11}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[0].Hash(),
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       1,
				},
			},
			// derived fields:
			TxHash:            txs[0].Hash(),
			ContractAddress:   common.HexToAddress("0x5a443704dd4b594b382c22a083e2bd3090a6fef3"),
			GasUsed:           1,
			EffectiveGasPrice: big.NewInt(11),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  0,
		},
		&Receipt{
			PostState:         common.Hash{2}.Bytes(),
			CumulativeGasUsed: 3,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x22}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[1].Hash(),
					TxIndex:     1,
					BlockHash:   blockHash,
					Index:       2,
				},
				{
					Address: common.BytesToAddress([]byte{0x02, 0x22}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[1].Hash(),
					TxIndex:     1,
					BlockHash:   blockHash,
					Index:       3,
				},
			},
			// derived fields:
			TxHash:            txs[1].Hash(),
			GasUsed:           2,
			EffectiveGasPrice: big.NewInt(22),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  1,
		},
		&Receipt{
			Type:              AccessListTxType,
			PostState:         common.Hash{3}.Bytes(),
			CumulativeGasUsed: 6,
			Logs:              []*Log{},
			// derived fields:
			TxHash:            txs[2].Hash(),
			GasUsed:           3,
			EffectiveGasPrice: big.NewInt(33),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  2,
		},
		&Receipt{
			Type:              DynamicFeeTxType,
			PostState:         common.Hash{4}.Bytes(),
			CumulativeGasUsed: 10,
			Logs:              []*Log{},
			// derived fields:
			TxHash:            txs[3].Hash(),
			GasUsed:           4,
			EffectiveGasPrice: big.NewInt(1044),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  3,
		},
		&Receipt{
			Type:              DynamicFeeTxType,
			PostState:         common.Hash{5}.Bytes(),
			CumulativeGasUsed: 15,
			Logs:              []*Log{},
			// derived fields:
			TxHash:            txs[4].Hash(),
			GasUsed:           5,
			EffectiveGasPrice: big.NewInt(1055),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  4,
		},
		&Receipt{
			Type:              BlobTxType,
			PostState:         common.Hash{6}.Bytes(),
			CumulativeGasUsed: 21,
			Logs:              []*Log{},
			// derived fields:
			TxHash:            txs[5].Hash(),
			GasUsed:           6,
			EffectiveGasPrice: big.NewInt(1066),
			BlobGasUsed:       params.BlobTxBlobGasPerBlob,
			BlobGasPrice:      big.NewInt(920),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  5,
		},
		&Receipt{
			Type:              BlobTxType,
			PostState:         common.Hash{7}.Bytes(),
			CumulativeGasUsed: 28,
			Logs:              []*Log{},
			// derived fields:
			TxHash:            txs[6].Hash(),
			GasUsed:           7,
			EffectiveGasPrice: big.NewInt(1077),
			BlobGasUsed:       3 * params.BlobTxBlobGasPerBlob,
			BlobGasPrice:      big.NewInt(920),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  6,
		},
		&Receipt{
			Type:              DepositTxType,
			PostState:         common.Hash{5}.Bytes(),
			CumulativeGasUsed: 50 + 28,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x33}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[7].Hash(),
					TxIndex:     7,
					BlockHash:   blockHash,
					Index:       4,
				},
				{
					Address: common.BytesToAddress([]byte{0x03, 0x33}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[7].Hash(),
					TxIndex:     7,
					BlockHash:   blockHash,
					Index:       5,
				},
			},
			TxHash:                txs[7].Hash(),
			ContractAddress:       common.HexToAddress("0x3bb898b4bbe24f68a4e9be46cfe72d1787fd74f4"),
			GasUsed:               50,
			EffectiveGasPrice:     big.NewInt(0),
			BlockHash:             blockHash,
			BlockNumber:           blockNumber,
			TransactionIndex:      7,
			DepositNonce:          &depNonce1,
			DepositReceiptVersion: nil,
		},
		&Receipt{
			Type:              DepositTxType,
			PostState:         common.Hash{5}.Bytes(),
			CumulativeGasUsed: 60 + 50 + 28,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x33}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[8].Hash(),
					TxIndex:     8,
					BlockHash:   blockHash,
					Index:       6,
				},
				{
					Address: common.BytesToAddress([]byte{0x03, 0x33}),
					Topics:  []common.Hash{common.HexToHash("dead"), common.HexToHash("beef")},
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[8].Hash(),
					TxIndex:     8,
					BlockHash:   blockHash,
					Index:       7,
				},
			},
			TxHash:                txs[8].Hash(),
			ContractAddress:       common.HexToAddress("0x117814af22cb83d8ad6e8489e9477d28265bc105"),
			GasUsed:               60,
			EffectiveGasPrice:     big.NewInt(0),
			BlockHash:             blockHash,
			BlockNumber:           blockNumber,
			TransactionIndex:      8,
			DepositNonce:          &depNonce2,
			DepositReceiptVersion: &canyonDepositReceiptVersion,
		},
	}
)

func TestDecodeEmptyTypedReceipt(t *testing.T) {
	input := []byte{0x80}
	var r Receipt
	err := rlp.DecodeBytes(input, &r)
	if err != errShortTypedReceipt {
		t.Fatal("wrong error:", err)
	}
}

// Tests that receipt data can be correctly derived from the contextual infos
func TestDeriveFields(t *testing.T) {
	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	blobGasPrice := big.NewInt(920)
	derivedReceipts := clearComputedFieldsOnReceipts(receipts)
	err := Receipts(derivedReceipts).DeriveFields(params.TestChainConfig, blockHash, blockNumber.Uint64(), blockTime, baseFee, blobGasPrice, txs)
	if err != nil {
		t.Fatalf("DeriveFields(...) = %v, want <nil>", err)
	}

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

// Test that we can marshal/unmarshal receipts to/from json without errors.
// This also confirms that our test receipts contain all the required fields.
func TestReceiptJSON(t *testing.T) {
	for i := range receipts {
		b, err := receipts[i].MarshalJSON()
		if err != nil {
			t.Fatal("error marshaling receipt to json:", err)
		}
		r := Receipt{}
		err = r.UnmarshalJSON(b)
		if err != nil {
			t.Fatal("error unmarshaling receipt from json:", err)
		}

		// Make sure marshal/unmarshal doesn't affect receipt hash root computation by comparing
		// the output of EncodeIndex
		rsBefore := Receipts([]*Receipt{receipts[i]})
		rsAfter := Receipts([]*Receipt{&r})

		encBefore, encAfter := bytes.Buffer{}, bytes.Buffer{}
		rsBefore.EncodeIndex(0, &encBefore)
		rsAfter.EncodeIndex(0, &encAfter)
		if !bytes.Equal(encBefore.Bytes(), encAfter.Bytes()) {
			t.Errorf("%v: EncodeIndex differs after JSON marshal/unmarshal", i)
		}
	}
}

// Test we can still parse receipt without EffectiveGasPrice for backwards compatibility, even
// though it is required per the spec.
func TestEffectiveGasPriceNotRequired(t *testing.T) {
	r := *receipts[0]
	r.EffectiveGasPrice = nil
	b, err := r.MarshalJSON()
	if err != nil {
		t.Fatal("error marshaling receipt to json:", err)
	}
	r2 := Receipt{}
	err = r2.UnmarshalJSON(b)
	if err != nil {
		t.Fatal("error unmarshaling receipt from json:", err)
	}
}

// TestTypedReceiptEncodingDecoding reproduces a flaw that existed in the receipt
// rlp decoder, which failed due to a shadowing error.
func TestTypedReceiptEncodingDecoding(t *testing.T) {
	var payload = common.FromHex("f9043eb9010c01f90108018262d4b9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0b9010c01f901080182cd14b9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0b9010d01f901090183013754b9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0b9010d01f90109018301a194b9010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c0")
	check := func(bundle []*Receipt) {
		t.Helper()
		for i, receipt := range bundle {
			if got, want := receipt.Type, uint8(1); got != want {
				t.Fatalf("bundle %d: got %x, want %x", i, got, want)
			}
		}
	}
	{
		var bundle []*Receipt
		rlp.DecodeBytes(payload, &bundle)
		check(bundle)
	}
	{
		var bundle []*Receipt
		r := bytes.NewReader(payload)
		s := rlp.NewStream(r, uint64(len(payload)))
		if err := s.Decode(&bundle); err != nil {
			t.Fatal(err)
		}
		check(bundle)
	}
}

func TestReceiptMarshalBinary(t *testing.T) {
	// Legacy Receipt
	legacyReceipt.Bloom = CreateBloom(Receipts{legacyReceipt})
	have, err := legacyReceipt.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary error: %v", err)
	}
	legacyReceipts := Receipts{legacyReceipt}
	buf := new(bytes.Buffer)
	legacyReceipts.EncodeIndex(0, buf)
	haveEncodeIndex := buf.Bytes()
	if !bytes.Equal(have, haveEncodeIndex) {
		t.Errorf("BinaryMarshal and EncodeIndex mismatch, got %x want %x", have, haveEncodeIndex)
	}
	buf.Reset()
	if err := legacyReceipt.EncodeRLP(buf); err != nil {
		t.Fatalf("encode rlp error: %v", err)
	}
	haveRLPEncode := buf.Bytes()
	if !bytes.Equal(have, haveRLPEncode) {
		t.Errorf("BinaryMarshal and EncodeRLP mismatch for legacy tx, got %x want %x", have, haveRLPEncode)
	}
	legacyWant := common.FromHex("f901c58001b9010000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000010000080000000000000000000004000000000000000000000000000040000000000000000000000000000800000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000f8bef85d940000000000000000000000000000000000000011f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100fff85d940000000000000000000000000000000000000111f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100ff")
	if !bytes.Equal(have, legacyWant) {
		t.Errorf("encoded RLP mismatch, got %x want %x", have, legacyWant)
	}

	// 2930 Receipt
	buf.Reset()
	accessListReceipt.Bloom = CreateBloom(Receipts{accessListReceipt})
	have, err = accessListReceipt.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary error: %v", err)
	}
	accessListReceipts := Receipts{accessListReceipt}
	accessListReceipts.EncodeIndex(0, buf)
	haveEncodeIndex = buf.Bytes()
	if !bytes.Equal(have, haveEncodeIndex) {
		t.Errorf("BinaryMarshal and EncodeIndex mismatch, got %x want %x", have, haveEncodeIndex)
	}
	accessListWant := common.FromHex("01f901c58001b9010000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000010000080000000000000000000004000000000000000000000000000040000000000000000000000000000800000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000f8bef85d940000000000000000000000000000000000000011f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100fff85d940000000000000000000000000000000000000111f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100ff")
	if !bytes.Equal(have, accessListWant) {
		t.Errorf("encoded RLP mismatch, got %x want %x", have, accessListWant)
	}

	// 1559 Receipt
	buf.Reset()
	eip1559Receipt.Bloom = CreateBloom(Receipts{eip1559Receipt})
	have, err = eip1559Receipt.MarshalBinary()
	if err != nil {
		t.Fatalf("marshal binary error: %v", err)
	}
	eip1559Receipts := Receipts{eip1559Receipt}
	eip1559Receipts.EncodeIndex(0, buf)
	haveEncodeIndex = buf.Bytes()
	if !bytes.Equal(have, haveEncodeIndex) {
		t.Errorf("BinaryMarshal and EncodeIndex mismatch, got %x want %x", have, haveEncodeIndex)
	}
	eip1559Want := common.FromHex("02f901c58001b9010000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000010000080000000000000000000004000000000000000000000000000040000000000000000000000000000800000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000f8bef85d940000000000000000000000000000000000000011f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100fff85d940000000000000000000000000000000000000111f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100ff")
	if !bytes.Equal(have, eip1559Want) {
		t.Errorf("encoded RLP mismatch, got %x want %x", have, eip1559Want)
	}
}

func TestReceiptUnmarshalBinary(t *testing.T) {
	// Legacy Receipt
	legacyBinary := common.FromHex("f901c58001b9010000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000010000080000000000000000000004000000000000000000000000000040000000000000000000000000000800000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000f8bef85d940000000000000000000000000000000000000011f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100fff85d940000000000000000000000000000000000000111f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100ff")
	gotLegacyReceipt := new(Receipt)
	if err := gotLegacyReceipt.UnmarshalBinary(legacyBinary); err != nil {
		t.Fatalf("unmarshal binary error: %v", err)
	}
	legacyReceipt.Bloom = CreateBloom(Receipts{legacyReceipt})
	if !reflect.DeepEqual(gotLegacyReceipt, legacyReceipt) {
		t.Errorf("receipt unmarshalled from binary mismatch, got %v want %v", gotLegacyReceipt, legacyReceipt)
	}

	// 2930 Receipt
	accessListBinary := common.FromHex("01f901c58001b9010000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000010000080000000000000000000004000000000000000000000000000040000000000000000000000000000800000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000f8bef85d940000000000000000000000000000000000000011f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100fff85d940000000000000000000000000000000000000111f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100ff")
	gotAccessListReceipt := new(Receipt)
	if err := gotAccessListReceipt.UnmarshalBinary(accessListBinary); err != nil {
		t.Fatalf("unmarshal binary error: %v", err)
	}
	accessListReceipt.Bloom = CreateBloom(Receipts{accessListReceipt})
	if !reflect.DeepEqual(gotAccessListReceipt, accessListReceipt) {
		t.Errorf("receipt unmarshalled from binary mismatch, got %v want %v", gotAccessListReceipt, accessListReceipt)
	}

	// 1559 Receipt
	eip1559RctBinary := common.FromHex("02f901c58001b9010000000000000010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000014000000000000000000000000000000000000000000000000000000000000000000000000000010000080000000000000000000004000000000000000000000000000040000000000000000000000000000800000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000f8bef85d940000000000000000000000000000000000000011f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100fff85d940000000000000000000000000000000000000111f842a0000000000000000000000000000000000000000000000000000000000000deada0000000000000000000000000000000000000000000000000000000000000beef830100ff")
	got1559Receipt := new(Receipt)
	if err := got1559Receipt.UnmarshalBinary(eip1559RctBinary); err != nil {
		t.Fatalf("unmarshal binary error: %v", err)
	}
	eip1559Receipt.Bloom = CreateBloom(Receipts{eip1559Receipt})
	if !reflect.DeepEqual(got1559Receipt, eip1559Receipt) {
		t.Errorf("receipt unmarshalled from binary mismatch, got %v want %v", got1559Receipt, eip1559Receipt)
	}
}

func clearComputedFieldsOnReceipts(receipts []*Receipt) []*Receipt {
	r := make([]*Receipt, len(receipts))
	for i, receipt := range receipts {
		r[i] = clearComputedFieldsOnReceipt(receipt)
	}
	return r
}

func clearComputedFieldsOnReceipt(receipt *Receipt) *Receipt {
	cpy := *receipt
	cpy.TxHash = common.Hash{0xff, 0xff, 0x11}
	cpy.BlockHash = common.Hash{0xff, 0xff, 0x22}
	cpy.BlockNumber = big.NewInt(math.MaxUint32)
	cpy.TransactionIndex = math.MaxUint32
	cpy.ContractAddress = common.Address{0xff, 0xff, 0x33}
	cpy.GasUsed = 0xffffffff
	cpy.Logs = clearComputedFieldsOnLogs(receipt.Logs)
	cpy.EffectiveGasPrice = big.NewInt(0)
	cpy.BlobGasUsed = 0
	cpy.BlobGasPrice = nil
	return &cpy
}

func clearComputedFieldsOnLogs(logs []*Log) []*Log {
	l := make([]*Log, len(logs))
	for i, log := range logs {
		cpy := *log
		cpy.BlockNumber = math.MaxUint32
		cpy.BlockHash = common.Hash{}
		cpy.TxHash = common.Hash{}
		cpy.TxIndex = math.MaxUint32
		cpy.Index = math.MaxUint32
		l[i] = &cpy
	}
	return l
}

func getOptimismTxReceipts(
	t *testing.T, l1AttributesPayload []byte,
	l1GasPrice, l1GasUsed *big.Int, feeScalar *big.Float, l1Fee *big.Int) ([]*Transaction, []*Receipt) {
	//to4 := common.HexToAddress("0x4")
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
	txs, receipts := getOptimismTxReceipts(t, payload, l1GasPrice, l1GasUsed, feeScalar, l1Fee)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnReceipts(receipts)
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	if err != nil {
		t.Fatalf("DeriveFields(...) = %v, want <nil>", err)
	}
	checkBedrockReceipts(t, receipts, derivedReceipts)

	// Should get same result with the Ecotone config because it will assume this is "first ecotone block"
	// if it sees the bedrock style L1 attributes.
	err = Receipts(derivedReceipts).DeriveFields(ecotoneTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	if err != nil {
		t.Fatalf("DeriveFields(...) = %v, want <nil>", err)
	}
	checkBedrockReceipts(t, receipts, derivedReceipts)
}

func TestDeriveOptimismEcotoneTxReceipts(t *testing.T) {
	// Ecotone style l1 attributes with baseFeeScalar=2, blobBaseFeeScalar=3, baseFee=1000*1e6, blobBaseFee=10*1e6
	payload := common.Hex2Bytes("440a5e20000000020000000300000000000004d200000000000004d200000000000004d2000000000000000000000000000000000000000000000000000000003b9aca00000000000000000000000000000000000000000000000000000000000098968000000000000000000000000000000000000000000000000000000000000004d200000000000000000000000000000000000000000000000000000000000004d2")
	// the parameters we use below are defined in rollup_test.go
	l1GasPrice := baseFee
	l1GasUsed := ecotoneGas
	l1Fee := ecotoneFee
	txs, receipts := getOptimismTxReceipts(t, payload, l1GasPrice, l1GasUsed, nil /*feeScalar*/, l1Fee)

	// Re-derive receipts.
	baseFee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnReceipts(receipts)
	// Should error out if we try to process this with a pre-Ecotone config
	err := Receipts(derivedReceipts).DeriveFields(bedrockGenesisTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	if err == nil {
		t.Fatalf("expected error from deriving ecotone receipts with pre-ecotone config, got none")
	}

	err = Receipts(derivedReceipts).DeriveFields(ecotoneTestConfig, blockHash, blockNumber.Uint64(), 0, baseFee, nil, txs)
	if err != nil {
		t.Fatalf("DeriveFields(...) = %v, want <nil>", err)
	}
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
