// Copyright 2023 The go-ethereum Authors
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

package ethapi

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestNewRPCTransactionDepositTx(t *testing.T) {
	tx := types.NewTx(&types.DepositTx{
		SourceHash:          common.HexToHash("0x1234"),
		IsSystemTransaction: true,
		Mint:                big.NewInt(34),
	})
	nonce := uint64(7)
	receipt := &types.Receipt{
		DepositNonce: &nonce,
	}
	got := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{}, receipt)
	// Should provide zero values for unused fields that are required in other transactions
	require.Equal(t, got.GasPrice, (*hexutil.Big)(big.NewInt(0)), "newRPCTransaction().GasPrice = %v, want 0x0", got.GasPrice)
	require.Equal(t, got.V, (*hexutil.Big)(big.NewInt(0)), "newRPCTransaction().V = %v, want 0x0", got.V)
	require.Equal(t, got.R, (*hexutil.Big)(big.NewInt(0)), "newRPCTransaction().R = %v, want 0x0", got.R)
	require.Equal(t, got.S, (*hexutil.Big)(big.NewInt(0)), "newRPCTransaction().S = %v, want 0x0", got.S)

	// Should include deposit tx specific fields
	require.Equal(t, *got.SourceHash, tx.SourceHash(), "newRPCTransaction().SourceHash = %v, want %v", got.SourceHash, tx.SourceHash())
	require.Equal(t, *got.IsSystemTx, tx.IsSystemTx(), "newRPCTransaction().IsSystemTx = %v, want %v", got.IsSystemTx, tx.IsSystemTx())
	require.Equal(t, got.Mint, (*hexutil.Big)(tx.Mint()), "newRPCTransaction().Mint = %v, want %v", got.Mint, tx.Mint())
	require.Equal(t, got.Nonce, (hexutil.Uint64)(nonce), "newRPCTransaction().Mint = %v, want %v", got.Nonce, nonce)
}

func TestNewRPCTransactionOmitIsSystemTxFalse(t *testing.T) {
	tx := types.NewTx(&types.DepositTx{
		IsSystemTransaction: false,
	})
	got := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{}, nil)

	require.Nil(t, got.IsSystemTx, "should omit IsSystemTx when false")
}

func TestUnmarshalRpcDepositTx(t *testing.T) {
	tests := []struct {
		name     string
		modifier func(tx *RPCTransaction)
		valid    bool
	}{
		{
			name:     "Unmodified",
			modifier: func(tx *RPCTransaction) {},
			valid:    true,
		},
		{
			name: "Zero Values",
			modifier: func(tx *RPCTransaction) {
				tx.V = (*hexutil.Big)(common.Big0)
				tx.R = (*hexutil.Big)(common.Big0)
				tx.S = (*hexutil.Big)(common.Big0)
				tx.GasPrice = (*hexutil.Big)(common.Big0)
			},
			valid: true,
		},
		{
			name: "Nil Values",
			modifier: func(tx *RPCTransaction) {
				tx.V = nil
				tx.R = nil
				tx.S = nil
				tx.GasPrice = nil
			},
			valid: true,
		},
		{
			name: "Non-Zero GasPrice",
			modifier: func(tx *RPCTransaction) {
				tx.GasPrice = (*hexutil.Big)(big.NewInt(43))
			},
			valid: false,
		},
		{
			name: "Non-Zero V",
			modifier: func(tx *RPCTransaction) {
				tx.V = (*hexutil.Big)(big.NewInt(43))
			},
			valid: false,
		},
		{
			name: "Non-Zero R",
			modifier: func(tx *RPCTransaction) {
				tx.R = (*hexutil.Big)(big.NewInt(43))
			},
			valid: false,
		},
		{
			name: "Non-Zero S",
			modifier: func(tx *RPCTransaction) {
				tx.S = (*hexutil.Big)(big.NewInt(43))
			},
			valid: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tx := types.NewTx(&types.DepositTx{
				SourceHash:          common.HexToHash("0x1234"),
				IsSystemTransaction: true,
				Mint:                big.NewInt(34),
			})
			rpcTx := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{}, nil)
			test.modifier(rpcTx)
			json, err := json.Marshal(rpcTx)
			require.NoError(t, err, "marshalling failed: %w", err)
			parsed := &types.Transaction{}
			err = parsed.UnmarshalJSON(json)
			if test.valid {
				require.NoError(t, err, "unmarshal failed: %w", err)
			} else {
				require.Error(t, err, "unmarshal should have failed but did not")
			}
		})
	}
}

func TestTransaction_RoundTripRpcJSON(t *testing.T) {
	var (
		config = params.AllEthashProtocolChanges
		signer = types.LatestSigner(config)
		key, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		tests  = allTransactionTypes(common.Address{0xde, 0xad}, config)
	)
	t.Parallel()
	for i, tt := range tests {
		var tx2 types.Transaction
		tx, err := types.SignNewTx(key, signer, tt)
		if err != nil {
			t.Fatalf("test %d: signing failed: %v", i, err)
		}
		// Regular transaction
		if data, err := json.Marshal(tx); err != nil {
			t.Fatalf("test %d: marshalling failed; %v", i, err)
		} else if err = tx2.UnmarshalJSON(data); err != nil {
			t.Fatalf("test %d: sunmarshal failed: %v", i, err)
		} else if want, have := tx.Hash(), tx2.Hash(); want != have {
			t.Fatalf("test %d: stx changed, want %x have %x", i, want, have)
		}

		//  rpcTransaction
		rpcTx := newRPCTransaction(tx, common.Hash{}, 0, 0, nil, config, nil)
		if data, err := json.Marshal(rpcTx); err != nil {
			t.Fatalf("test %d: marshalling failed; %v", i, err)
		} else if err = tx2.UnmarshalJSON(data); err != nil {
			t.Fatalf("test %d: unmarshal failed: %v", i, err)
		} else if want, have := tx.Hash(), tx2.Hash(); want != have {
			t.Fatalf("test %d: tx changed, want %x have %x", i, want, have)
		}
	}
}

func allTransactionTypes(addr common.Address, config *params.ChainConfig) []types.TxData {
	return []types.TxData{
		&types.LegacyTx{
			Nonce:    5,
			GasPrice: big.NewInt(6),
			Gas:      7,
			To:       &addr,
			Value:    big.NewInt(8),
			Data:     []byte{0, 1, 2, 3, 4},
			V:        big.NewInt(9),
			R:        big.NewInt(10),
			S:        big.NewInt(11),
		},
		&types.LegacyTx{
			Nonce:    5,
			GasPrice: big.NewInt(6),
			Gas:      7,
			To:       nil,
			Value:    big.NewInt(8),
			Data:     []byte{0, 1, 2, 3, 4},
			V:        big.NewInt(32),
			R:        big.NewInt(10),
			S:        big.NewInt(11),
		},
		&types.AccessListTx{
			ChainID:  config.ChainID,
			Nonce:    5,
			GasPrice: big.NewInt(6),
			Gas:      7,
			To:       &addr,
			Value:    big.NewInt(8),
			Data:     []byte{0, 1, 2, 3, 4},
			AccessList: types.AccessList{
				types.AccessTuple{
					Address:     common.Address{0x2},
					StorageKeys: []common.Hash{types.EmptyRootHash},
				},
			},
			V: big.NewInt(32),
			R: big.NewInt(10),
			S: big.NewInt(11),
		},
		&types.AccessListTx{
			ChainID:  config.ChainID,
			Nonce:    5,
			GasPrice: big.NewInt(6),
			Gas:      7,
			To:       nil,
			Value:    big.NewInt(8),
			Data:     []byte{0, 1, 2, 3, 4},
			AccessList: types.AccessList{
				types.AccessTuple{
					Address:     common.Address{0x2},
					StorageKeys: []common.Hash{types.EmptyRootHash},
				},
			},
			V: big.NewInt(32),
			R: big.NewInt(10),
			S: big.NewInt(11),
		},
		&types.DynamicFeeTx{
			ChainID:   config.ChainID,
			Nonce:     5,
			GasTipCap: big.NewInt(6),
			GasFeeCap: big.NewInt(9),
			Gas:       7,
			To:        &addr,
			Value:     big.NewInt(8),
			Data:      []byte{0, 1, 2, 3, 4},
			AccessList: types.AccessList{
				types.AccessTuple{
					Address:     common.Address{0x2},
					StorageKeys: []common.Hash{types.EmptyRootHash},
				},
			},
			V: big.NewInt(32),
			R: big.NewInt(10),
			S: big.NewInt(11),
		},
		&types.DynamicFeeTx{
			ChainID:    config.ChainID,
			Nonce:      5,
			GasTipCap:  big.NewInt(6),
			GasFeeCap:  big.NewInt(9),
			Gas:        7,
			To:         nil,
			Value:      big.NewInt(8),
			Data:       []byte{0, 1, 2, 3, 4},
			AccessList: types.AccessList{},
			V:          big.NewInt(32),
			R:          big.NewInt(10),
			S:          big.NewInt(11),
		},
	}
}
