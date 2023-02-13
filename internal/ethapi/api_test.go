package ethapi

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func Test_newRPCTransaction_UnmarshalDepositTx_RoundTrip(t *testing.T) {
	tx := types.NewTx(&types.DepositTx{
		SourceHash:          common.HexToHash("0x1234"),
		IsSystemTransaction: true,
		Mint:                big.NewInt(34),
	})
	rpcTx := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{})
	json, err := json.Marshal(rpcTx)
	require.NoError(t, err, "marshalling failed")
	parsed := &types.Transaction{}
	err = parsed.UnmarshalJSON(json)
	require.NoError(t, err, "unmarshal failed")
}

func Test_newRPCTransaction_UnmarshalDepositTx_NilValues(t *testing.T) {
	tx := types.NewTx(&types.DepositTx{
		SourceHash:          common.HexToHash("0x1234"),
		IsSystemTransaction: true,
		Mint:                big.NewInt(34),
	})
	rpcTx := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{})
	rpcTx.V = nil
	rpcTx.R = nil
	rpcTx.S = nil
	rpcTx.GasPrice = nil
	json, err := json.Marshal(rpcTx)
	require.NoError(t, err, "marshalling failed")
	parsed := &types.Transaction{}
	err = parsed.UnmarshalJSON(json)
	require.NoError(t, err, "unmarshal failed")
}

func Test_newRPCTransaction_UnmarshalDepositTx_ZeroValues(t *testing.T) {
	tx := types.NewTx(&types.DepositTx{
		SourceHash:          common.HexToHash("0x1234"),
		IsSystemTransaction: true,
		Mint:                big.NewInt(34),
	})
	rpcTx := newRPCTransaction(tx, common.Hash{}, uint64(12), uint64(1), big.NewInt(0), &params.ChainConfig{})
	rpcTx.V = (*hexutil.Big)(common.Big0)
	rpcTx.R = (*hexutil.Big)(common.Big0)
	rpcTx.S = (*hexutil.Big)(common.Big0)
	rpcTx.GasPrice = (*hexutil.Big)(common.Big0)
	json, err := json.Marshal(rpcTx)
	require.NoError(t, err, "marshalling failed")
	parsed := &types.Transaction{}
	err = parsed.UnmarshalJSON(json)
	require.NoError(t, err, "unmarshal failed")
}
