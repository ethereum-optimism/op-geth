// Copyright 2022 The go-ethereum Authors
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
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>

package miner

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/clique"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/txpool/legacypool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/params"
)

var (
	// Test chain configurations
	testTxPoolConfig  legacypool.Config
	ethashChainConfig *params.ChainConfig
	cliqueChainConfig *params.ChainConfig

	// Test accounts
	testBankKey, _  = crypto.GenerateKey()
	testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
	testBankFunds   = big.NewInt(1000000000000000000)

	testUserKey, _  = crypto.GenerateKey()
	testUserAddress = crypto.PubkeyToAddress(testUserKey.PublicKey)

	testRecipient = common.HexToAddress("0xdeadbeef")
	testTimestamp = uint64(time.Now().Unix())

	// Test transactions
	pendingTxs []*types.Transaction
	newTxs     []*types.Transaction

	testConfig = Config{
		PendingFeeRecipient: testBankAddress,
		Recommit:            time.Second,
		GasCeil:             params.GenesisGasLimit,
	}
)

func init() {
	testTxPoolConfig = legacypool.DefaultConfig
	testTxPoolConfig.Journal = ""
	ethashChainConfig = new(params.ChainConfig)
	*ethashChainConfig = *params.TestChainConfig
	cliqueChainConfig = new(params.ChainConfig)
	*cliqueChainConfig = *params.TestChainConfig
	cliqueChainConfig.Clique = &params.CliqueConfig{
		Period: 10,
		Epoch:  30000,
	}

	signer := types.LatestSigner(params.TestChainConfig)
	tx1 := types.MustSignNewTx(testBankKey, signer, &types.AccessListTx{
		ChainID:  params.TestChainConfig.ChainID,
		Nonce:    0,
		To:       &testUserAddress,
		Value:    big.NewInt(1000),
		Gas:      params.TxGas,
		GasPrice: big.NewInt(params.InitialBaseFee),
	})
	pendingTxs = append(pendingTxs, tx1)

	tx2 := types.MustSignNewTx(testBankKey, signer, &types.LegacyTx{
		Nonce:    1,
		To:       &testUserAddress,
		Value:    big.NewInt(1000),
		Gas:      params.TxGas,
		GasPrice: big.NewInt(params.InitialBaseFee),
	})
	newTxs = append(newTxs, tx2)
}

// testWorkerBackend implements worker.Backend interfaces and wraps all information needed during the testing.
type testWorkerBackend struct {
	db      ethdb.Database
	txPool  *txpool.TxPool
	chain   *core.BlockChain
	genesis *core.Genesis
}

func newTestWorkerBackend(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, n int) *testWorkerBackend {
	var gspec = &core.Genesis{
		Config: chainConfig,
		Alloc:  types.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
	}
	switch e := engine.(type) {
	case *clique.Clique:
		gspec.ExtraData = make([]byte, 32+common.AddressLength+crypto.SignatureLength)
		copy(gspec.ExtraData[32:32+common.AddressLength], testBankAddress.Bytes())
		e.Authorize(testBankAddress, func(account accounts.Account, s string, data []byte) ([]byte, error) {
			return crypto.Sign(crypto.Keccak256(data), testBankKey)
		})
	case *ethash.Ethash:
	default:
		t.Fatalf("unexpected consensus engine type: %T", engine)
	}
	if chainConfig.HoloceneTime != nil {
		// genesis block extraData needs to be correct format
		gspec.ExtraData = []byte{0, 0, 1, 2, 3, 4, 5, 6, 7}
	}
	chain, err := core.NewBlockChain(db, &core.CacheConfig{TrieDirtyDisabled: true}, gspec, nil, engine, vm.Config{}, nil)
	if err != nil {
		t.Fatalf("core.NewBlockChain failed: %v", err)
	}
	pool := legacypool.New(testTxPoolConfig, chain)
	txpool, _ := txpool.New(testTxPoolConfig.PriceLimit, chain, []txpool.SubPool{pool})

	return &testWorkerBackend{
		db:      db,
		chain:   chain,
		txPool:  txpool,
		genesis: gspec,
	}
}

func (b *testWorkerBackend) BlockChain() *core.BlockChain { return b.chain }
func (b *testWorkerBackend) TxPool() *txpool.TxPool       { return b.txPool }

func newTestWorker(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, db ethdb.Database, blocks int) (*Miner, *testWorkerBackend) {
	backend := newTestWorkerBackend(t, chainConfig, engine, db, blocks)
	backend.txPool.Add(pendingTxs, true, true)
	w := New(backend, testConfig, engine)
	return w, backend
}

func TestBuildPayload(t *testing.T) {
	t.Run("no-tx-pool", func(t *testing.T) { testBuildPayload(t, true, false, nil) })
	// no-tx-pool case with interrupt not interesting because no-tx-pool doesn't run
	// the builder routine
	t.Run("with-tx-pool", func(t *testing.T) { testBuildPayload(t, false, false, nil) })
	t.Run("with-tx-pool-interrupt", func(t *testing.T) { testBuildPayload(t, false, true, nil) })
	params1559 := []byte{0, 1, 2, 3, 4, 5, 6, 7}
	t.Run("with-params", func(t *testing.T) { testBuildPayload(t, false, false, params1559) })
	t.Run("with-params-no-tx-pool", func(t *testing.T) { testBuildPayload(t, true, false, params1559) })
	t.Run("with-params-interrupt", func(t *testing.T) { testBuildPayload(t, false, true, params1559) })

	t.Run("wrong-config-no-params", func(t *testing.T) { testBuildPayloadWrongConfig(t, nil) })
	t.Run("wrong-config-params", func(t *testing.T) { testBuildPayloadWrongConfig(t, params1559) })

	zeroParams := make([]byte, 8)
	t.Run("with-zero-params", func(t *testing.T) { testBuildPayload(t, true, false, zeroParams) })
}

func holoceneConfig() *params.ChainConfig {
	config := *params.TestChainConfig
	config.LondonBlock = big.NewInt(0)
	t := uint64(0)
	config.CanyonTime = &t
	config.HoloceneTime = &t
	canyonDenom := uint64(250)
	config.Optimism = &params.OptimismConfig{
		EIP1559Elasticity:        6,
		EIP1559Denominator:       50,
		EIP1559DenominatorCanyon: &canyonDenom,
	}
	return &config
}

// newPayloadArgs returns a BuildPaylooadArgs with the given parentHash and eip-1559 params,
// testTimestamp for Timestamp, and testRecipient for recipient. NoTxPool is set to true.
func newPayloadArgs(parentHash common.Hash, params1559 []byte) *BuildPayloadArgs {
	return &BuildPayloadArgs{
		Parent:        parentHash,
		Timestamp:     testTimestamp,
		Random:        common.Hash{},
		FeeRecipient:  testRecipient,
		NoTxPool:      true,
		EIP1559Params: params1559,
	}
}

func testBuildPayload(t *testing.T, noTxPool, interrupt bool, params1559 []byte) {
	t.Parallel()
	db := rawdb.NewMemoryDatabase()
	config := params.TestChainConfig
	if len(params1559) != 0 {
		config = holoceneConfig()
	}
	w, b := newTestWorker(t, config, ethash.NewFaker(), db, 0)

	const numInterruptTxs = 256
	if interrupt {
		// when doing interrupt testing, create a large pool so interruption will
		// definitely be visible.
		txs := genTxs(1, numInterruptTxs)
		b.txPool.Add(txs, true, false)
	}

	args := newPayloadArgs(b.chain.CurrentBlock().Hash(), params1559)
	args.NoTxPool = noTxPool

	// payload resolution now interrupts block building, so we have to
	// wait for the payloading building process to build its first block
	payload, err := w.buildPayload(args, false)
	if err != nil {
		t.Fatalf("Failed to build payload %v", err)
	}
	verify := func(outer *engine.ExecutionPayloadEnvelope, txs int) {
		t.Helper()
		if outer == nil {
			t.Fatal("ExecutionPayloadEnvelope is nil")
		}
		payload := outer.ExecutionPayload
		if payload.ParentHash != b.chain.CurrentBlock().Hash() {
			t.Fatal("Unexpected parent hash")
		}
		if payload.Random != (common.Hash{}) {
			t.Fatal("Unexpected random value")
		}
		if payload.Timestamp != testTimestamp {
			t.Fatal("Unexpected timestamp")
		}
		if payload.FeeRecipient != testRecipient {
			t.Fatal("Unexpected fee recipient")
		}
		if !interrupt && len(payload.Transactions) != txs {
			t.Fatalf("Unexpect transaction set: got %d, expected %d", len(payload.Transactions), txs)
		} else if interrupt && len(payload.Transactions) >= txs {
			t.Fatalf("Unexpect transaction set: got %d, expected less than %d", len(payload.Transactions), txs)
		}
	}

	// make sure the 1559 params we've specied (if any) ends up in both the full and empty block headers
	var expected []byte
	if len(params1559) != 0 {
		expected = []byte{0}
		d, _ := eip1559.DecodeHolocene1559Params(params1559)
		if d == 0 {
			expected = append(expected, eip1559.EncodeHolocene1559Params(250, 6)...) // canyon defaults
		} else {
			expected = append(expected, params1559...)
		}
	}
	if payload.full != nil && !bytes.Equal(payload.full.Header().Extra, expected) {
		t.Fatalf("ExtraData doesn't match. want: %x, got %x", expected, payload.full.Header().Extra)
	}
	if payload.empty != nil && !bytes.Equal(payload.empty.Header().Extra, expected) {
		t.Fatalf("ExtraData doesn't match on empty block. want: %x, got %x", expected, payload.empty.Header().Extra)
	}

	if noTxPool {
		// we only build the empty block when ignoring the tx pool
		empty := payload.ResolveEmpty()
		verify(empty, 0)
		full := payload.ResolveFull()
		verify(full, 0)
	} else if interrupt {
		full := payload.ResolveFull()
		verify(full, len(pendingTxs)+numInterruptTxs)
	} else { // tx-pool and no interrupt
		payload.WaitFull()
		full := payload.ResolveFull()
		verify(full, len(pendingTxs))
	}

	// Ensure resolve can be called multiple times and the
	// result should be unchanged
	dataOne := payload.Resolve()
	dataTwo := payload.Resolve()
	if !reflect.DeepEqual(dataOne, dataTwo) {
		t.Fatal("Unexpected payload data")
	}
}

func testBuildPayloadWrongConfig(t *testing.T, params1559 []byte) {
	t.Parallel()
	db := rawdb.NewMemoryDatabase()
	config := holoceneConfig()
	if len(params1559) != 0 {
		// deactivate holocene and make sure non-empty params get rejected
		config.HoloceneTime = nil
	}
	w, b := newTestWorker(t, config, ethash.NewFaker(), db, 0)

	args := newPayloadArgs(b.chain.CurrentBlock().Hash(), params1559)
	payload, err := w.buildPayload(args, false)
	if err == nil && (payload == nil || payload.err == nil) {
		t.Fatalf("expected error, got none")
	}
}

func TestBuildPayloadInvalidHoloceneParams(t *testing.T) {
	t.Parallel()
	db := rawdb.NewMemoryDatabase()
	config := holoceneConfig()
	w, b := newTestWorker(t, config, ethash.NewFaker(), db, 0)

	// 0 denominators shouldn't be allowed
	badParams := eip1559.EncodeHolocene1559Params(0, 6)

	args := newPayloadArgs(b.chain.CurrentBlock().Hash(), badParams)
	payload, err := w.buildPayload(args, false)
	if err == nil && (payload == nil || payload.err == nil) {
		t.Fatalf("expected error, got none")
	}
}

func genTxs(startNonce, count uint64) types.Transactions {
	txs := make(types.Transactions, 0, count)
	signer := types.LatestSigner(params.TestChainConfig)
	for nonce := startNonce; nonce < startNonce+count; nonce++ {
		txs = append(txs, types.MustSignNewTx(testBankKey, signer, &types.AccessListTx{
			ChainID:  params.TestChainConfig.ChainID,
			Nonce:    nonce,
			To:       &testUserAddress,
			Value:    big.NewInt(1000),
			Gas:      params.TxGas,
			GasPrice: big.NewInt(params.InitialBaseFee),
		}))
	}
	return txs
}

func TestPayloadId(t *testing.T) {
	t.Parallel()
	ids := make(map[string]int)
	for i, tt := range []*BuildPayloadArgs{
		{
			Parent:       common.Hash{1},
			Timestamp:    1,
			Random:       common.Hash{0x1},
			FeeRecipient: common.Address{0x1},
		},
		// Different parent
		{
			Parent:       common.Hash{2},
			Timestamp:    1,
			Random:       common.Hash{0x1},
			FeeRecipient: common.Address{0x1},
		},
		// Different timestamp
		{
			Parent:       common.Hash{2},
			Timestamp:    2,
			Random:       common.Hash{0x1},
			FeeRecipient: common.Address{0x1},
		},
		// Different Random
		{
			Parent:       common.Hash{2},
			Timestamp:    2,
			Random:       common.Hash{0x2},
			FeeRecipient: common.Address{0x1},
		},
		// Different fee-recipient
		{
			Parent:       common.Hash{2},
			Timestamp:    2,
			Random:       common.Hash{0x2},
			FeeRecipient: common.Address{0x2},
		},
		// Different withdrawals (non-empty)
		{
			Parent:       common.Hash{2},
			Timestamp:    2,
			Random:       common.Hash{0x2},
			FeeRecipient: common.Address{0x2},
			Withdrawals: []*types.Withdrawal{
				{
					Index:     0,
					Validator: 0,
					Address:   common.Address{},
					Amount:    0,
				},
			},
		},
		// Different withdrawals (non-empty)
		{
			Parent:       common.Hash{2},
			Timestamp:    2,
			Random:       common.Hash{0x2},
			FeeRecipient: common.Address{0x2},
			Withdrawals: []*types.Withdrawal{
				{
					Index:     2,
					Validator: 0,
					Address:   common.Address{},
					Amount:    0,
				},
			},
		},
	} {
		id := tt.Id().String()
		if prev, exists := ids[id]; exists {
			t.Errorf("ID collision, case %d and case %d: id %v", prev, i, id)
		}
		ids[id] = i
	}
}
