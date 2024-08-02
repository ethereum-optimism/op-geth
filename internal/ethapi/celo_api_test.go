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
	blockTime        uint64 = 100
)

func TestNewRPCTransactionLegacy(t *testing.T) {
	config := allEnabledChainConfig()
	// Set cel2 time to 2000 so that we don't activate the cel2 fork.
	var cel2Time uint64 = 2000
	config.Cel2Time = &cel2Time
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
		checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, nil)
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
		checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, nil)
	})
}

func TestNewRPCTransactionDynamicFee(t *testing.T) {
	key, err := crypto.GenerateKey()
	require.NoError(t, err)
	feeCap := big.NewInt(1000)
	tipCap := big.NewInt(100)

	t.Run("PendingTransactions", func(t *testing.T) {
		// For pending transactions we expect the gas price to be the gas fee cap.
		gasFeeCap := func(t *testing.T, tx *types.Transaction, rpcTx *RPCTransaction) {
			assert.Equal(t, (*hexutil.Big)(feeCap), rpcTx.GasPrice)
		}
		overrides := map[string]func(*testing.T, *types.Transaction, *RPCTransaction){"gasPrice": gasFeeCap}
		config := allEnabledChainConfig()
		s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

		// An empty bockhash signals pending transactions (I.E no mined block)
		blockhash := common.Hash{}
		t.Run("DynamicFeeTx", func(t *testing.T) {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   config.ChainID,
				Nonce:     nonce,
				Gas:       gasLimit,
				GasFeeCap: feeCap,
				GasTipCap: tipCap,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})

		t.Run("CeloDynamicFeeTxV2", func(t *testing.T) {
			tx := types.NewTx(&types.CeloDynamicFeeTxV2{
				ChainID:     config.ChainID,
				Nonce:       nonce,
				Gas:         gasLimit,
				GasFeeCap:   feeCap,
				GasTipCap:   tipCap,
				FeeCurrency: &feeCurrency,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})
	})

	t.Run("PreGingerbreadMinedDynamicTxs", func(t *testing.T) {
		nilGasPrice := func(t *testing.T, tx *types.Transaction, rpcTx *RPCTransaction) {
			assert.Nil(t, rpcTx.GasPrice)
		}
		overrides := map[string]func(*testing.T, *types.Transaction, *RPCTransaction){"gasPrice": nilGasPrice}
		// For a pre gingerbread mined dynamic txs we expect the gas price to be unset, because without the state we
		// cannot retrieve the base fee, and we currently have no implementation in op-geth to handle retrieving the
		// base fee from state.
		config := allEnabledChainConfig()
		config.GingerbreadBlock = big.NewInt(200) // Setup config so that gingerbread is not active.
		cel2Time := uint64(1000)
		config.Cel2Time = &cel2Time // also deactivate cel2
		s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

		t.Run("DynamicFeeTx", func(t *testing.T) {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   config.ChainID,
				Nonce:     nonce,
				Gas:       gasLimit,
				GasFeeCap: feeCap,
				GasTipCap: tipCap,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})

		t.Run("CeloDynamicFeeTx", func(t *testing.T) {
			tx := types.NewTx(&types.CeloDynamicFeeTx{
				ChainID:             config.ChainID,
				Nonce:               nonce,
				Gas:                 gasLimit,
				GasFeeCap:           feeCap,
				GasTipCap:           tipCap,
				GatewayFee:          gatewayFee,
				GatewayFeeRecipient: &gatewayFeeRecipient,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})
	})

	t.Run("PostGingerbreadMinedDynamicTxsWithNativeFeeCurrency", func(t *testing.T) {
		// For a post gingerbread mined dynamic tx with a native fee currency we expect the gas price to be the
		// effective gas price calculated with the base fee available on the block.
		effectiveGasPrice := func(t *testing.T, tx *types.Transaction, rpcTx *RPCTransaction) {
			assert.Equal(t, (*hexutil.Big)(effectiveGasPrice(tx, baseFee)), rpcTx.GasPrice)
		}
		overrides := map[string]func(*testing.T, *types.Transaction, *RPCTransaction){"gasPrice": effectiveGasPrice}

		config := allEnabledChainConfig()
		s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

		t.Run("DynamicFeeTx", func(t *testing.T) {
			tx := types.NewTx(&types.DynamicFeeTx{
				ChainID:   config.ChainID,
				Nonce:     nonce,
				Gas:       gasLimit,
				GasFeeCap: feeCap,
				GasTipCap: tipCap,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})

		t.Run("CeloDynamicFeeTx", func(t *testing.T) {
			// CeloDynamicFeeTxs are deprecated after cel2 so we need to ensure cel2time is not activated
			config := allEnabledChainConfig()
			cel2Time := uint64(1000)
			config.Cel2Time = &cel2Time
			s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

			tx := types.NewTx(&types.CeloDynamicFeeTx{
				ChainID:             config.ChainID,
				Nonce:               nonce,
				Gas:                 gasLimit,
				GasFeeCap:           feeCap,
				GasTipCap:           tipCap,
				GatewayFee:          gatewayFee,
				GatewayFeeRecipient: &gatewayFeeRecipient,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})

		t.Run("CeloDynamicFeeTxV2", func(t *testing.T) {
			tx := types.NewTx(&types.CeloDynamicFeeTxV2{
				ChainID:   config.ChainID,
				Nonce:     nonce,
				Gas:       gasLimit,
				GasFeeCap: feeCap,
				GasTipCap: tipCap,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})

		// TODO unskip this when cip 66 txs are enabled currently they are not supporeted in the celo signer.
		t.Run("CeloDenominatedTx", func(t *testing.T) {
			t.Skip("CeloDenominatedTx is currently not supported in the celo signer")
			tx := types.NewTx(&types.CeloDenominatedTx{
				ChainID:             config.ChainID,
				Nonce:               nonce,
				Gas:                 gasLimit,
				GasFeeCap:           feeCap,
				GasTipCap:           tipCap,
				FeeCurrency:         &feeCurrency,
				MaxFeeInFeeCurrency: big.NewInt(100000),

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})
	})

	t.Run("PostGingerbreadPreCel2MinedDynamicTxsWithNonNativeFeeCurrency", func(t *testing.T) {
		// For a post gingerbread mined dynamic txs with a non native fee currency we expect the gas price to be unset,
		// because without the state we cannot retrieve the base fee, and we currently have no implementation in op-geth
		// to handle retrieving the base fee from state.

		nilGasPrice := func(t *testing.T, tx *types.Transaction, rpcTx *RPCTransaction) {
			assert.Nil(t, rpcTx.GasPrice)
		}
		overrides := map[string]func(*testing.T, *types.Transaction, *RPCTransaction){"gasPrice": nilGasPrice}

		config := allEnabledChainConfig()
		cel2Time := uint64(1000)
		config.Cel2Time = &cel2Time // Deactivate cel2
		s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

		t.Run("CeloDynamicFeeTx", func(t *testing.T) {
			// CeloDynamicFeeTxs are deprecated after cel2 so we need to ensure cel2time is not activated
			config := allEnabledChainConfig()
			cel2Time := uint64(1000)
			config.Cel2Time = &cel2Time
			s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

			tx := types.NewTx(&types.CeloDynamicFeeTx{
				ChainID:             config.ChainID,
				Nonce:               nonce,
				Gas:                 gasLimit,
				GasFeeCap:           feeCap,
				GasTipCap:           tipCap,
				FeeCurrency:         &feeCurrency,
				GatewayFee:          gatewayFee,
				GatewayFeeRecipient: &gatewayFeeRecipient,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})

		t.Run("CeloDynamicFeeTxV2", func(t *testing.T) {
			tx := types.NewTx(&types.CeloDynamicFeeTxV2{
				ChainID:     config.ChainID,
				Nonce:       nonce,
				Gas:         gasLimit,
				GasFeeCap:   feeCap,
				GasTipCap:   tipCap,
				FeeCurrency: &feeCurrency,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, nil)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})
	})

	t.Run("PostCel2MinedDynamicTxs", func(t *testing.T) {
		receipt := &types.Receipt{}
		receipt.EffectiveGasPrice = big.NewInt(1234)
		effectiveGasPrice := func(t *testing.T, tx *types.Transaction, rpcTx *RPCTransaction) {
			assert.Equal(t, (*hexutil.Big)(receipt.EffectiveGasPrice), rpcTx.GasPrice)
		}
		overrides := map[string]func(*testing.T, *types.Transaction, *RPCTransaction){"gasPrice": effectiveGasPrice}

		config := allEnabledChainConfig()
		s := types.MakeSigner(config, new(big.Int).SetUint64(blockNumber), blockTime)

		t.Run("CeloDynamicFeeTxV2", func(t *testing.T) {
			// For a pre gingerbread mined dynamic fee tx we expect the gas price to be unset.
			tx := types.NewTx(&types.CeloDynamicFeeTxV2{
				ChainID:     config.ChainID,
				Nonce:       nonce,
				Gas:         gasLimit,
				GasFeeCap:   feeCap,
				GasTipCap:   tipCap,
				FeeCurrency: &feeCurrency,

				To:    &to,
				Value: value,
				Data:  []byte{},
			})

			signed, err := types.SignTx(tx, s, key)
			require.NoError(t, err)

			rpcTx := newRPCTransaction(signed, blockhash, blockNumber, blockTime, transactionIndex, baseFee, config, receipt)
			checkTxFields(t, signed, rpcTx, s, blockhash, blockNumber, transactionIndex, overrides)
		})
	})
}

func allEnabledChainConfig() *params.ChainConfig {
	zeroTime := uint64(0)
	return &params.ChainConfig{
		ChainID:             big.NewInt(44787),
		HomesteadBlock:      big.NewInt(0),
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
		ArrowGlacierBlock:   big.NewInt(0),
		GrayGlacierBlock:    big.NewInt(0),
		ShanghaiTime:        &zeroTime,
		CancunTime:          &zeroTime,
		RegolithTime:        &zeroTime,
		CanyonTime:          &zeroTime,
		EcotoneTime:         &zeroTime,
		FjordTime:           &zeroTime,
		Cel2Time:            &zeroTime,
		GingerbreadBlock:    big.NewInt(0),
	}
}

// checkTxFields for the most part checks that the fields of the rpcTx match those of the provided tx, it allows for
// overriding some checks by providing a map of fieldName -> overrideFunc.
func checkTxFields(
	t *testing.T,
	tx *types.Transaction,
	rpcTx *RPCTransaction,
	signer types.Signer,
	blockhash common.Hash,
	blockNumber uint64,
	transactionIndex uint64,
	overrides map[string]func(*testing.T, *types.Transaction, *RPCTransaction),
) {
	// Added fields (not part of the transaction type)
	//
	// If blockhash is empty it signifies a pending tx and for pending txs the block hash, block number and tx index are
	// not set on the rpcTx. on the result.
	if blockhash == (common.Hash{}) {
		assert.Nil(t, rpcTx.BlockHash)
		assert.Nil(t, rpcTx.BlockNumber)
		assert.Nil(t, rpcTx.TransactionIndex)
	} else {
		assert.Equal(t, &blockhash, rpcTx.BlockHash)
		assert.Equal(t, (*hexutil.Big)(big.NewInt(int64(blockNumber))), rpcTx.BlockNumber)
		assert.Equal(t, hexutil.Uint64(transactionIndex), *rpcTx.TransactionIndex)
	}

	from, err := types.Sender(signer, tx)
	require.NoError(t, err)

	assert.Equal(t, from, rpcTx.From)
	assert.Equal(t, hexutil.Uint64(tx.Gas()), rpcTx.Gas)
	assert.Equal(t, tx.To(), rpcTx.To)
	override, ok := overrides["gasPrice"]
	if ok {
		override(t, tx, rpcTx)
	} else {
		assert.Equal(t, (*hexutil.Big)(tx.GasPrice()), rpcTx.GasPrice)
	}
	switch tx.Type() {
	case types.DynamicFeeTxType, types.CeloDynamicFeeTxType, types.CeloDynamicFeeTxV2Type, types.CeloDenominatedTxType:
		assert.Equal(t, (*hexutil.Big)(tx.GasFeeCap()), rpcTx.GasFeeCap)
		assert.Equal(t, (*hexutil.Big)(tx.GasTipCap()), rpcTx.GasTipCap)
	default:
		assert.Nil(t, rpcTx.GasFeeCap)
		assert.Nil(t, rpcTx.GasTipCap)
	}
	assert.Equal(t, (*hexutil.Big)(tx.BlobGasFeeCap()), rpcTx.MaxFeePerBlobGas)
	assert.Equal(t, tx.Hash(), rpcTx.Hash)
	assert.Equal(t, (hexutil.Bytes)(tx.Data()), rpcTx.Input)
	assert.Equal(t, hexutil.Uint64(tx.Nonce()), rpcTx.Nonce)
	assert.Equal(t, tx.To(), rpcTx.To)
	assert.Equal(t, (*hexutil.Big)(tx.Value()), rpcTx.Value)
	assert.Equal(t, hexutil.Uint64(tx.Type()), rpcTx.Type)
	switch tx.Type() {
	case types.AccessListTxType, types.DynamicFeeTxType, types.CeloDynamicFeeTxType, types.CeloDynamicFeeTxV2Type, types.CeloDenominatedTxType, types.BlobTxType:
		assert.Equal(t, tx.AccessList(), *rpcTx.Accesses)
	default:
		assert.Nil(t, rpcTx.Accesses)
	}

	assert.Equal(t, (*hexutil.Big)(tx.ChainId()), rpcTx.ChainID)
	assert.Equal(t, tx.BlobHashes(), rpcTx.BlobVersionedHashes)

	v, r, s := tx.RawSignatureValues()
	assert.Equal(t, (*hexutil.Big)(v), rpcTx.V)
	assert.Equal(t, (*hexutil.Big)(r), rpcTx.R)
	assert.Equal(t, (*hexutil.Big)(s), rpcTx.S)

	switch tx.Type() {
	case types.AccessListTxType, types.DynamicFeeTxType, types.CeloDynamicFeeTxType, types.CeloDynamicFeeTxV2Type, types.CeloDenominatedTxType, types.BlobTxType:
		yparity := (hexutil.Uint64)(v.Sign())
		assert.Equal(t, &yparity, rpcTx.YParity)
	default:
		assert.Nil(t, rpcTx.YParity)
	}

	// optimism fields
	switch tx.Type() {
	case types.DepositTxType:
		assert.Equal(t, tx.SourceHash(), rpcTx.SourceHash)
		assert.Equal(t, tx.Mint(), rpcTx.Mint)
		assert.Equal(t, tx.IsSystemTx(), rpcTx.IsSystemTx)
	default:
		assert.Nil(t, rpcTx.SourceHash)
		assert.Nil(t, rpcTx.Mint)
		assert.Nil(t, rpcTx.IsSystemTx)
	}

	assert.Nil(t, rpcTx.DepositReceiptVersion)

	// celo fields
	assert.Equal(t, tx.FeeCurrency(), rpcTx.FeeCurrency)
	assert.Equal(t, (*hexutil.Big)(tx.MaxFeeInFeeCurrency()), rpcTx.MaxFeeInFeeCurrency)
	if tx.Type() == types.LegacyTxType && tx.IsCeloLegacy() {
		assert.Equal(t, false, *rpcTx.EthCompatible)
	} else {
		assert.Nil(t, rpcTx.EthCompatible)
	}
	assert.Equal(t, (*hexutil.Big)(tx.GatewayFee()), rpcTx.GatewayFee)
	assert.Equal(t, tx.GatewayFeeRecipient(), rpcTx.GatewayFeeRecipient)
}
