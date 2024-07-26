package ethapi

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// tx fields
	nonce               uint64 = 1
	gasPrice                   = big.NewInt(1000)
	gasLimit            uint64 = 100000
	feeCurrency                = common.HexToAddress("0x0000000000000000000000000000000000000bbb")
	gatewayFee                 = big.NewInt(500)
	gatewayFeeRecipient        = common.HexToAddress("0x0000000000000000000000000000000000000ccc")
	to                         = common.HexToAddress("0x0000000000000000000000000000000000000aaa")
	value                      = big.NewInt(10)
	// block fields
	baseFee                 = big.NewInt(100)
	transactionIndex uint64 = 15
	blockhash               = common.HexToHash("0x6ba4a8c1bfe2619eb498e5296e81b1c393b13cba0198ed63dea0ee3aa619b073")
	blockNumber      uint64 = 100
)

func TestNewRPCTransactionLegacy(t *testing.T) {
	// Enable all block based forks so we get the most recent upstream signer.
	var cel2Time uint64 = 2000
	config := &params.ChainConfig{
		ChainID:             big.NewInt(44787),
		EIP150Block:         big.NewInt(0),
		EIP155Block:         big.NewInt(0),
		EIP158Block:         big.NewInt(0),
		ByzantiumBlock:      big.NewInt(0),
		ConstantinopleBlock: big.NewInt(0),
		PetersburgBlock:     big.NewInt(0),
		IstanbulBlock:       big.NewInt(0),
		MuirGlacierBlock:    big.NewInt(0),
		BerlinBlock:         big.NewInt(0),
		LondonBlock:         big.NewInt(0),
		Cel2Time:            &cel2Time,
	}
	// Block time is set to before Cel2Time so we don't activate the cel2 fork.
	// This gives us the celo legacy signer (legacy transactions are deprecated
	// after cel2).
	blockTime := uint64(1000)
	s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

	key, err := crypto.GenerateKey()
	require.NoError(t, err)

	t.Run("WithCeloFields", func(t *testing.T) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,

			FeeCurrency:         &feeCurrency,
			GatewayFee:          gatewayFee,
			GatewayFeeRecipient: &gatewayFeeRecipient,

			To:    &to,
			Value: value,
			Data:  []byte{},

			CeloLegacy: true,
		})

		signed, err := types.SignTx(tx, s, key)
		require.NoError(t, err)

		rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)

		// check newRPCTransaction has the expected fields
		// Ethereum fields
		checkRPCTransactionFields(
			t,
			rpcTx,
			to,
			value,
			gasLimit,
			gasPrice,
			nonce,
			config.ChainID,
			signed.Hash(),
			blockhash,
			blockNumber,
			transactionIndex,
		)
		// Celo fields
		assert.Equal(t, feeCurrency, *rpcTx.FeeCurrency)
		assert.Equal(t, (*hexutil.Big)(gatewayFee), rpcTx.GatewayFee)
		assert.Equal(t, gatewayFeeRecipient, *rpcTx.GatewayFeeRecipient)
		assert.Equal(t, false, *rpcTx.EthCompatible)
		assert.Nil(t, rpcTx.MaxFeeInFeeCurrency)
	})

	t.Run("WithoutCeloFields", func(t *testing.T) {
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      gasLimit,

			To:    &to,
			Value: value,
			Data:  []byte{},
		})
		signed, err := types.SignTx(tx, s, key)
		require.NoError(t, err)
		rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)

		// check newRPCTransaction has the expected fields
		// Ethereum fields
		checkRPCTransactionFields(
			t,
			rpcTx,
			to,
			value,
			gasLimit,
			gasPrice,
			nonce,
			config.ChainID,
			signed.Hash(),
			blockhash,
			blockNumber,
			transactionIndex,
		)
		// Celo fields
		assert.Nil(t, rpcTx.FeeCurrency)
		assert.Nil(t, rpcTx.GatewayFee)
		assert.Nil(t, rpcTx.GatewayFeeRecipient)
		assert.Nil(t, rpcTx.EthCompatible)
		assert.Nil(t, rpcTx.MaxFeeInFeeCurrency)
	})
}

func checkRPCTransactionFields(
	t *testing.T,
	rpcTx *RPCTransaction,
	to common.Address,
	value *big.Int,
	gasLimit uint64,
	gasPrice *big.Int,
	nonce uint64,
	chainID *big.Int,
	hash common.Hash,
	blockhash common.Hash,
	blockNumber uint64,
	transactionIndex uint64,
) {
	assert.Equal(t, to, *rpcTx.To)
	assert.Equal(t, (*hexutil.Big)(value), rpcTx.Value)
	assert.Equal(t, hexutil.Bytes{}, rpcTx.Input)
	assert.Equal(t, hexutil.Uint64(gasLimit), rpcTx.Gas)
	assert.Equal(t, (*hexutil.Big)(gasPrice), rpcTx.GasPrice)
	assert.Equal(t, hash, rpcTx.Hash)
	assert.Equal(t, hexutil.Uint64(nonce), rpcTx.Nonce)
	assert.Equal(t, (*hexutil.Big)(chainID), rpcTx.ChainID)
	assert.Equal(t, hexutil.Uint64(types.LegacyTxType), rpcTx.Type)
	assert.Nil(t, rpcTx.Accesses)
	assert.Nil(t, rpcTx.GasFeeCap)
	assert.Nil(t, rpcTx.GasTipCap)
	assert.Nil(t, rpcTx.MaxFeePerBlobGas)
	assert.Equal(t, []common.Hash(nil), rpcTx.BlobVersionedHashes)

	// Added fields (not part of the transaction type)
	assert.Equal(t, &blockhash, rpcTx.BlockHash)
	assert.Equal(t, (*hexutil.Big)(big.NewInt(int64(blockNumber))), rpcTx.BlockNumber)
	assert.Equal(t, hexutil.Uint64(transactionIndex), *rpcTx.TransactionIndex)
}
