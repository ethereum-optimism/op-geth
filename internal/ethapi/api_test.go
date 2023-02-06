package ethapi

import (
	"encoding/json"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestNewRPCTransactionDepositTx(t *testing.T) {
	tx := types.NewTx(&types.DepositTx{
		SourceHash:          common.HexToHash("0x1234"),
		IsSystemTransaction: true,
		Mint:                big.NewInt(34),
	})
	got := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{})
	// Should provide zero values for unused fields that are required in other transactions
	if !reflect.DeepEqual(got.GasPrice, (*hexutil.Big)(big.NewInt(0))) {
		t.Errorf("newRPCTransaction().GasPrice = %v, want 0x0", got.GasPrice)
	}
	if !reflect.DeepEqual(got.V, (*hexutil.Big)(big.NewInt(0))) {
		t.Errorf("newRPCTransaction().V = %v, want 0x0", got.V)
	}
	if !reflect.DeepEqual(got.R, (*hexutil.Big)(big.NewInt(0))) {
		t.Errorf("newRPCTransaction().R = %v, want 0x0", got.R)
	}
	if !reflect.DeepEqual(got.S, (*hexutil.Big)(big.NewInt(0))) {
		t.Errorf("newRPCTransaction().S = %v, want 0x0", got.S)
	}

	// Should include deposit tx specific fields
	if *got.SourceHash != tx.SourceHash() {
		t.Errorf("newRPCTransaction().SourceHash = %v, want %v", got.SourceHash, tx.SourceHash())
	}
	if *got.IsSystemTx != tx.IsSystemTx() {
		t.Errorf("newRPCTransaction().IsSystemTx = %v, want %v", got.IsSystemTx, tx.IsSystemTx())
	}
	if !reflect.DeepEqual(got.Mint, (*hexutil.Big)(tx.Mint())) {
		t.Errorf("newRPCTransaction().Mint = %v, want %v", got.Mint, tx.Mint())
	}
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
			rpcTx := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{})
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
