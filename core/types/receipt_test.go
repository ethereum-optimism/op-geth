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
	"github.com/kylelemons/godebug/diff"
	"github.com/stretchr/testify/require"
)

var (
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
		Status:            ReceiptStatusFailed,
		CumulativeGasUsed: 1,
		DepositNonce:      &nonce,
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
	// Create a few transactions to have receipts for
	to2 := common.HexToAddress("0x2")
	to3 := common.HexToAddress("0x3")
	to4 := common.HexToAddress("0x4")
	to5 := common.HexToAddress("0x5")
	txs := Transactions{
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
			GasFeeCap: big.NewInt(1045),
		}),
		NewTx(&DynamicFeeTx{
			To:        &to5,
			Nonce:     5,
			Value:     big.NewInt(5),
			Gas:       5,
			GasTipCap: big.NewInt(56),
			GasFeeCap: big.NewInt(1055),
		}),
		NewTx(&DepositTx{
			To:    nil, // contract creation
			Value: big.NewInt(6),
			Gas:   50,
		}),
	}
	depNonce := uint64(7)
	blockNumber := big.NewInt(1)
	blockHash := common.BytesToHash([]byte{0x03, 0x14})

	// Create the corresponding receipts
	receipts := Receipts{
		&Receipt{
			Status:            ReceiptStatusFailed,
			CumulativeGasUsed: 1,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x11}),
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[0].Hash(),
					TxIndex:     0,
					BlockHash:   blockHash,
					Index:       0,
				},
				{
					Address: common.BytesToAddress([]byte{0x01, 0x11}),
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
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[1].Hash(),
					TxIndex:     1,
					BlockHash:   blockHash,
					Index:       2,
				},
				{
					Address: common.BytesToAddress([]byte{0x02, 0x22}),
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
			Type:              DepositTxType,
			PostState:         common.Hash{5}.Bytes(),
			CumulativeGasUsed: 50 + 15,
			Logs: []*Log{
				{
					Address: common.BytesToAddress([]byte{0x33}),
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[5].Hash(),
					TxIndex:     5,
					BlockHash:   blockHash,
					Index:       4,
				},
				{
					Address: common.BytesToAddress([]byte{0x03, 0x33}),
					// derived fields:
					BlockNumber: blockNumber.Uint64(),
					TxHash:      txs[5].Hash(),
					TxIndex:     5,
					BlockHash:   blockHash,
					Index:       5,
				},
			},
			TxHash:            txs[5].Hash(),
			ContractAddress:   common.HexToAddress("0x3bb898b4bbe24f68a4e9be46cfe72d1787fd74f4"),
			GasUsed:           50,
			EffectiveGasPrice: big.NewInt(0),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  5,
			DepositNonce:      &depNonce,
		},
	}

	// Re-derive receipts.
	basefee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnReceipts(receipts)
	err := Receipts(derivedReceipts).DeriveFields(params.TestChainConfig, blockHash, blockNumber.Uint64(), 0, basefee, txs)
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

func TestDeriveOptimismTxReceipt(t *testing.T) {
	to4 := common.HexToAddress("0x4")
	// Create a few transactions to have receipts for
	txs := Transactions{
		NewTx(&DepositTx{
			To:    nil, // contract creation
			Value: big.NewInt(6),
			Gas:   50,
			// System config with L1Scalar=2_000_000 (becomes 2 after division), L1Overhead=2500, L1BaseFee=5000
			Data: common.Hex2Bytes("015d8eb900000000000000000000000000000000000000000000000026b39534042076f70000000000000000000000000000000000000000000000007e33b7c4995967580000000000000000000000000000000000000000000000000000000000001388547dea8ff339566349ed0ef6384876655d1b9b955e36ac165c6b8ab69b9af5cd0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000123400000000000000000000000000000000000000000000000000000000000009c400000000000000000000000000000000000000000000000000000000001e8480"),
		}),
		NewTx(&DynamicFeeTx{
			To:        &to4,
			Nonce:     4,
			Value:     big.NewInt(4),
			Gas:       4,
			GasTipCap: big.NewInt(44),
			GasFeeCap: big.NewInt(1045),
			Data:      []byte{0, 1, 255, 0},
		}),
	}
	depNonce := uint64(7)
	blockNumber := big.NewInt(1)
	blockHash := common.BytesToHash([]byte{0x03, 0x14})

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
			DepositNonce:      &depNonce,
		},
		&Receipt{
			Type:              DynamicFeeTxType,
			PostState:         common.Hash{4}.Bytes(),
			CumulativeGasUsed: 10,
			Logs:              []*Log{},
			// derived fields:
			TxHash:            txs[1].Hash(),
			GasUsed:           18446744073709551561,
			EffectiveGasPrice: big.NewInt(1044),
			BlockHash:         blockHash,
			BlockNumber:       blockNumber,
			TransactionIndex:  1,
			L1GasPrice:        big.NewInt(5000),
			L1GasUsed:         big.NewInt(3976),
			L1Fee:             big.NewInt(39760000),
			FeeScalar:         big.NewFloat(2),
		},
	}

	// Re-derive receipts.
	basefee := big.NewInt(1000)
	derivedReceipts := clearComputedFieldsOnReceipts(receipts)
	err := Receipts(derivedReceipts).DeriveFields(params.OptimismTestConfig, blockHash, blockNumber.Uint64(), 0, basefee, txs)
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

	// Check that we preserved the invariant: l1Fee = l1GasPrice * l1GasUsed * l1FeeScalar
	// but with more difficult int math...
	l2Rcpt := derivedReceipts[1]
	l1GasCost := new(big.Int).Mul(l2Rcpt.L1GasPrice, l2Rcpt.L1GasUsed)
	l1Fee := new(big.Float).Mul(new(big.Float).SetInt(l1GasCost), l2Rcpt.FeeScalar)
	require.Equal(t, new(big.Float).SetInt(l2Rcpt.L1Fee), l1Fee)
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
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := test.rcpt.MarshalBinary()
			require.NoError(t, err)

			d := &Receipt{}
			err = d.UnmarshalBinary(data)
			require.NoError(t, err)
			require.Equal(t, test.rcpt, d)
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
		})
	}
}
