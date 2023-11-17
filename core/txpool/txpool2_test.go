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
package txpool

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rlp"

	"github.com/stretchr/testify/require"
)

func pricedValuedTransaction(nonce uint64, value int64, gaslimit uint64, gasprice *big.Int, key *ecdsa.PrivateKey) *types.Transaction {
	tx, _ := types.SignTx(types.NewTransaction(nonce, common.Address{}, big.NewInt(value), gaslimit, gasprice, nil), types.HomesteadSigner{}, key)
	return tx
}

func count(t *testing.T, pool *TxPool) (pending int, queued int) {
	t.Helper()
	pending, queued = pool.stats()
	if err := validatePoolInternals(pool); err != nil {
		t.Fatalf("pool internal state corrupted: %v", err)
	}
	return pending, queued
}

func fillPool(t *testing.T, pool *TxPool) {
	t.Helper()
	// Create a number of test accounts, fund them and make transactions
	executableTxs := types.Transactions{}
	nonExecutableTxs := types.Transactions{}
	for i := 0; i < 384; i++ {
		key, _ := crypto.GenerateKey()
		pool.currentState.AddBalance(crypto.PubkeyToAddress(key.PublicKey), big.NewInt(10000000000))
		// Add executable ones
		for j := 0; j < int(pool.config.AccountSlots); j++ {
			executableTxs = append(executableTxs, pricedTransaction(uint64(j), 100000, big.NewInt(300), key))
		}
	}
	// Import the batch and verify that limits have been enforced
	pool.AddRemotesSync(executableTxs)
	pool.AddRemotesSync(nonExecutableTxs)
	pending, queued := pool.Stats()
	slots := pool.all.Slots()
	// sanity-check that the test prerequisites are ok (pending full)
	if have, want := pending, slots; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}
	if have, want := queued, 0; have != want {
		t.Fatalf("have %d, want %d", have, want)
	}

	t.Logf("pool.config: GlobalSlots=%d, GlobalQueue=%d\n", pool.config.GlobalSlots, pool.config.GlobalQueue)
	t.Logf("pending: %d queued: %d, all: %d\n", pending, queued, slots)
}

// Tests that if a batch high-priced of non-executables arrive, they do not kick out
// executable transactions
func TestTransactionFutureAttack(t *testing.T) {
	t.Parallel()

	// Create the pool to test the limit enforcement with
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(1000000, statedb, new(event.Feed))
	config := testTxPoolConfig
	config.GlobalQueue = 100
	config.GlobalSlots = 100
	pool := NewTxPool(config, eip1559Config, blockchain)
	defer pool.Stop()
	fillPool(t, pool)
	pending, _ := pool.Stats()
	// Now, future transaction attack starts, let's add a bunch of expensive non-executables, and see if the pending-count drops
	{
		key, _ := crypto.GenerateKey()
		pool.currentState.AddBalance(crypto.PubkeyToAddress(key.PublicKey), big.NewInt(100000000000))
		futureTxs := types.Transactions{}
		for j := 0; j < int(pool.config.GlobalSlots+pool.config.GlobalQueue); j++ {
			futureTxs = append(futureTxs, pricedTransaction(1000+uint64(j), 100000, big.NewInt(500), key))
		}
		for i := 0; i < 5; i++ {
			pool.AddRemotesSync(futureTxs)
			newPending, newQueued := count(t, pool)
			t.Logf("pending: %d queued: %d, all: %d\n", newPending, newQueued, pool.all.Slots())
		}
	}
	newPending, _ := pool.Stats()
	// Pending should not have been touched
	if have, want := newPending, pending; have < want {
		t.Errorf("wrong pending-count, have %d, want %d (GlobalSlots: %d)",
			have, want, pool.config.GlobalSlots)
	}
}

// Tests that if a batch high-priced of non-executables arrive, they do not kick out
// executable transactions
func TestTransactionFuture1559(t *testing.T) {
	t.Parallel()
	// Create the pool to test the pricing enforcement with
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(1000000, statedb, new(event.Feed))
	pool := NewTxPool(testTxPoolConfig, eip1559Config, blockchain)
	defer pool.Stop()

	// Create a number of test accounts, fund them and make transactions
	fillPool(t, pool)
	pending, _ := pool.Stats()

	// Now, future transaction attack starts, let's add a bunch of expensive non-executables, and see if the pending-count drops
	{
		key, _ := crypto.GenerateKey()
		pool.currentState.AddBalance(crypto.PubkeyToAddress(key.PublicKey), big.NewInt(100000000000))
		futureTxs := types.Transactions{}
		for j := 0; j < int(pool.config.GlobalSlots+pool.config.GlobalQueue); j++ {
			futureTxs = append(futureTxs, dynamicFeeTx(1000+uint64(j), 100000, big.NewInt(200), big.NewInt(101), key))
		}
		pool.AddRemotesSync(futureTxs)
	}
	newPending, _ := pool.Stats()
	// Pending should not have been touched
	if have, want := newPending, pending; have != want {
		t.Errorf("Wrong pending-count, have %d, want %d (GlobalSlots: %d)",
			have, want, pool.config.GlobalSlots)
	}
}

// Tests that if a batch of balance-overdraft txs arrive, they do not kick out
// executable transactions
func TestTransactionZAttack(t *testing.T) {
	t.Parallel()
	// Create the pool to test the pricing enforcement with
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(1000000, statedb, new(event.Feed))
	pool := NewTxPool(testTxPoolConfig, eip1559Config, blockchain)
	defer pool.Stop()
	// Create a number of test accounts, fund them and make transactions
	fillPool(t, pool)

	countInvalidPending := func() int {
		t.Helper()
		var ivpendingNum int
		pendingtxs, _ := pool.Content()
		for account, txs := range pendingtxs {
			cur_balance := new(big.Int).Set(pool.currentState.GetBalance(account))
			for _, tx := range txs {
				if cur_balance.Cmp(tx.Value()) <= 0 {
					ivpendingNum++
				} else {
					cur_balance.Sub(cur_balance, tx.Value())
				}
			}
		}
		if err := validatePoolInternals(pool); err != nil {
			t.Fatalf("pool internal state corrupted: %v", err)
		}
		return ivpendingNum
	}
	ivPending := countInvalidPending()
	t.Logf("invalid pending: %d\n", ivPending)

	// Now, DETER-Z attack starts, let's add a bunch of expensive non-executables (from N accounts) along with balance-overdraft txs (from one account), and see if the pending-count drops
	for j := 0; j < int(pool.config.GlobalQueue); j++ {
		futureTxs := types.Transactions{}
		key, _ := crypto.GenerateKey()
		pool.currentState.AddBalance(crypto.PubkeyToAddress(key.PublicKey), big.NewInt(100000000000))
		futureTxs = append(futureTxs, pricedTransaction(1000+uint64(j), 21000, big.NewInt(500), key))
		pool.AddRemotesSync(futureTxs)
	}

	overDraftTxs := types.Transactions{}
	{
		key, _ := crypto.GenerateKey()
		pool.currentState.AddBalance(crypto.PubkeyToAddress(key.PublicKey), big.NewInt(100000000000))
		for j := 0; j < int(pool.config.GlobalSlots); j++ {
			overDraftTxs = append(overDraftTxs, pricedValuedTransaction(uint64(j), 60000000000, 21000, big.NewInt(500), key))
		}
	}
	pool.AddRemotesSync(overDraftTxs)
	pool.AddRemotesSync(overDraftTxs)
	pool.AddRemotesSync(overDraftTxs)
	pool.AddRemotesSync(overDraftTxs)
	pool.AddRemotesSync(overDraftTxs)

	newPending, newQueued := count(t, pool)
	newIvPending := countInvalidPending()
	t.Logf("pool.all.Slots(): %d\n", pool.all.Slots())
	t.Logf("pending: %d queued: %d, all: %d\n", newPending, newQueued, pool.all.Slots())
	t.Logf("invalid pending: %d\n", newIvPending)

	// Pending should not have been touched
	if newIvPending != ivPending {
		t.Errorf("Wrong invalid pending-count, have %d, want %d (GlobalSlots: %d, queued: %d)",
			newIvPending, ivPending, pool.config.GlobalSlots, newQueued)
	}
}

func decodeSignature(sig []byte) (r, s, v *big.Int) {
	if len(sig) != crypto.SignatureLength {
		panic(fmt.Sprintf("wrong size for signature: got %d, want %d", len(sig), crypto.SignatureLength))
	}
	r = new(big.Int).SetBytes(sig[:32])
	s = new(big.Int).SetBytes(sig[32:64])
	v = new(big.Int).SetBytes([]byte{sig[64] + 27})
	return r, s, v
}

func generateMetaTxData(dynamicTx *types.DynamicFeeTx, expireHeight uint64, sponsorPercent uint64,
	gasFeeSponsorAddr common.Address, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	metaTxSignData := &types.MetaTxSignData{
		ChainID:        dynamicTx.ChainID,
		Nonce:          dynamicTx.Nonce,
		GasTipCap:      dynamicTx.GasTipCap,
		GasFeeCap:      dynamicTx.GasFeeCap,
		Gas:            dynamicTx.Gas,
		To:             dynamicTx.To,
		Value:          dynamicTx.Value,
		Data:           dynamicTx.Data,
		AccessList:     dynamicTx.AccessList,
		ExpireHeight:   expireHeight,
		SponsorPercent: sponsorPercent,
	}

	sponsorSig, err := crypto.Sign(metaTxSignData.Hash().Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	r, s, v := decodeSignature(sponsorSig)

	metaTxData := &types.MetaTxParams{
		ExpireHeight:   expireHeight,
		Payload:        metaTxSignData.Data,
		GasFeeSponsor:  gasFeeSponsorAddr,
		SponsorPercent: sponsorPercent,
		R:              r,
		S:              s,
		V:              v,
	}

	metaTxDataBz, err := rlp.EncodeToBytes(metaTxData)
	if err != nil {
		return nil, err
	}

	return append(types.MetaTxPrefix, metaTxDataBz...), nil
}

// Tests that if invalid meta txs can be picked out and valid meta txs can be executed
func TestMetaTx(t *testing.T) {
	t.Parallel()
	// Create the pool to test the pricing enforcement with
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(10000000, statedb, new(event.Feed))
	pool := NewTxPool(testTxPoolConfig, eip1559Config, blockchain)
	defer pool.Stop()

	sponsorNum := 10
	gasFeeSponsorPrivateKeys := make([]*ecdsa.PrivateKey, 0, sponsorNum)
	gasFeeSponsorAddrs := make([]common.Address, 0, sponsorNum)
	for i := 0; i < sponsorNum; i++ {
		gasFeeSponsor, _ := crypto.GenerateKey()
		gasFeeSponsorPrivateKeys = append(gasFeeSponsorPrivateKeys, gasFeeSponsor)
		gasFeeSponsorAddrs = append(gasFeeSponsorAddrs, crypto.PubkeyToAddress(gasFeeSponsor.PublicKey))
		pool.currentState.AddBalance(crypto.PubkeyToAddress(gasFeeSponsor.PublicKey), new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1e4)))
	}
	userNum1 := 100
	usersAccountsWithoutBalancePrivateKey := make([]*ecdsa.PrivateKey, 0, userNum1)
	usersAccountsWithoutBalanceAddrs := make([]common.Address, 0, userNum1)
	for i := 0; i < userNum1; i++ {
		userAcc, _ := crypto.GenerateKey()
		usersAccountsWithoutBalancePrivateKey = append(usersAccountsWithoutBalancePrivateKey, userAcc)
		usersAccountsWithoutBalanceAddrs = append(usersAccountsWithoutBalanceAddrs, crypto.PubkeyToAddress(userAcc.PublicKey))
	}

	userNum2 := 10
	usersAccountsWithBalancePrivateKey := make([]*ecdsa.PrivateKey, 0, userNum2)
	usersAccountsWithBalanceAddrs := make([]common.Address, 0, userNum2)
	for i := 0; i < userNum2; i++ {
		userAcc, _ := crypto.GenerateKey()
		usersAccountsWithBalancePrivateKey = append(usersAccountsWithBalancePrivateKey, userAcc)
		usersAccountsWithBalanceAddrs = append(usersAccountsWithBalanceAddrs, crypto.PubkeyToAddress(userAcc.PublicKey))
		pool.currentState.AddBalance(crypto.PubkeyToAddress(userAcc.PublicKey), new(big.Int).Mul(big.NewInt(1e18), big.NewInt(1)))
	}

	chainId := params.TestChainConfig.ChainID
	signer := types.LatestSignerForChainID(chainId)
	approveABICallData, _ := hexutil.Decode("0x095ea7b30000000000000000000000001f9090aae28b8a3dceadf281b0f12828e676c3260000000000000000000000000000000000000000000000000de0b6b3a7640000")
	to := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	expireHeight := uint64(20_000_010)

	userNonceMap := make(map[common.Address]uint64)
	for i := 0; i < sponsorNum; i++ {
		for j := 0; j < userNum1; j++ {
			nonce := userNonceMap[usersAccountsWithoutBalanceAddrs[j]]
			dynamicTx := &types.DynamicFeeTx{
				ChainID:    chainId,
				Nonce:      nonce,
				GasTipCap:  big.NewInt(1e9),
				GasFeeCap:  big.NewInt(2e9),
				Gas:        4700000,
				To:         &to,
				Value:      big.NewInt(0),
				Data:       approveABICallData,
				AccessList: nil,
			}
			payload, err := generateMetaTxData(dynamicTx, expireHeight, 100, gasFeeSponsorAddrs[i], gasFeeSponsorPrivateKeys[i])
			require.NoError(t, err)

			dynamicTx.Data = payload
			tx := types.NewTx(dynamicTx)

			txSignature, err := crypto.Sign(signer.Hash(tx).Bytes(), usersAccountsWithoutBalancePrivateKey[j])
			require.NoError(t, err)
			signedTx, err := tx.WithSignature(signer, txSignature)
			require.NoError(t, err)

			err = pool.AddLocal(signedTx)
			require.NoError(t, err)

			userNonceMap[usersAccountsWithoutBalanceAddrs[j]] = nonce + 1
		}
	}

	all := pool.all.Count()
	pending, queued := pool.Stats()
	if pending != userNum1*sponsorNum {
		t.Errorf("Wrong pending-count, want %d, have %d",
			userNum1*sponsorNum, pending)
	}

	// increase nonce so that later tx can only be included into txpool queue
	for j := 0; j < userNum1; j++ {
		userNonceMap[usersAccountsWithoutBalanceAddrs[j]] = userNonceMap[usersAccountsWithoutBalanceAddrs[j]] + 1
	}

	for i := 0; i < sponsorNum; i++ {
		for j := 0; j < userNum1; j++ {
			nonce := userNonceMap[usersAccountsWithoutBalanceAddrs[j]]
			dynamicTx := &types.DynamicFeeTx{
				ChainID:    chainId,
				Nonce:      nonce,
				GasTipCap:  big.NewInt(1e9),
				GasFeeCap:  big.NewInt(2e9),
				Gas:        4700000,
				To:         &to,
				Value:      big.NewInt(0),
				Data:       approveABICallData,
				AccessList: nil,
			}
			payload, err := generateMetaTxData(dynamicTx, expireHeight, 100, gasFeeSponsorAddrs[i], gasFeeSponsorPrivateKeys[i])
			require.NoError(t, err)

			dynamicTx.Data = payload
			tx := types.NewTx(dynamicTx)

			txSignature, err := crypto.Sign(signer.Hash(tx).Bytes(), usersAccountsWithoutBalancePrivateKey[j])
			require.NoError(t, err)
			signedTx, err := tx.WithSignature(signer, txSignature)
			require.NoError(t, err)

			err = pool.AddLocal(signedTx)
			require.NoError(t, err)

			userNonceMap[usersAccountsWithoutBalanceAddrs[j]] = nonce + 1
		}
	}

	all = pool.all.Count()
	pending, queued = pool.Stats()
	if pending != userNum1*sponsorNum {
		t.Errorf("Wrong pending-count, want %d, have %d",
			userNum1*sponsorNum, pending)
	}
	if queued != userNum1*sponsorNum {
		t.Errorf("Wrong queued-count, want %d, have %d",
			userNum1*sponsorNum, pending)
	}
	if all != 2*userNum1*sponsorNum {
		t.Errorf("Wrong queued-count, want %d, have %d",
			2*userNum1*sponsorNum, all)
	}

	{
		nonce := userNonceMap[usersAccountsWithoutBalanceAddrs[0]]
		dynamicTx := &types.DynamicFeeTx{
			ChainID:    chainId,
			Nonce:      nonce,
			GasTipCap:  big.NewInt(1e9),
			GasFeeCap:  big.NewInt(2e9),
			Gas:        4700000,
			To:         &to,
			Value:      big.NewInt(0),
			Data:       approveABICallData,
			AccessList: nil,
		}
		payload, err := generateMetaTxData(dynamicTx, expireHeight, 90, gasFeeSponsorAddrs[0], gasFeeSponsorPrivateKeys[0])
		require.NoError(t, err)

		dynamicTx.Data = payload
		tx := types.NewTx(dynamicTx)

		txSignature, err := crypto.Sign(signer.Hash(tx).Bytes(), usersAccountsWithoutBalancePrivateKey[0])
		require.NoError(t, err)
		signedTx, err := tx.WithSignature(signer, txSignature)
		require.NoError(t, err)

		err = pool.AddLocal(signedTx)
		require.Equal(t, err, core.ErrInsufficientFunds)
	}

	for i := 0; i < sponsorNum; i++ {
		for j := 0; j < userNum2; j++ {
			nonce := userNonceMap[usersAccountsWithBalanceAddrs[j]]
			dynamicTx := &types.DynamicFeeTx{
				ChainID:    chainId,
				Nonce:      nonce,
				GasTipCap:  big.NewInt(1e9),
				GasFeeCap:  big.NewInt(2e9),
				Gas:        4700000,
				To:         &to,
				Value:      big.NewInt(0),
				Data:       approveABICallData,
				AccessList: nil,
			}
			payload, err := generateMetaTxData(dynamicTx, expireHeight, 90, gasFeeSponsorAddrs[i], gasFeeSponsorPrivateKeys[i])
			require.NoError(t, err)

			dynamicTx.Data = payload
			tx := types.NewTx(dynamicTx)

			txSignature, err := crypto.Sign(signer.Hash(tx).Bytes(), usersAccountsWithBalancePrivateKey[j])
			require.NoError(t, err)
			signedTx, err := tx.WithSignature(signer, txSignature)
			require.NoError(t, err)

			err = pool.AddLocal(signedTx)
			require.NoError(t, err)

			userNonceMap[usersAccountsWithBalanceAddrs[j]] = nonce + 1
		}
	}
	all1 := pool.all.Count()
	pending1, queued1 := pool.Stats()
	if pending1-pending != userNum2*sponsorNum {
		t.Errorf("Wrong pending-count delta, want %d, have %d",
			userNum2*sponsorNum, pending1-pending)
	}
	if queued1-queued != 0 {
		t.Errorf("Wrong queued-count delta, want %d, have %d",
			0, queued1-queued)
	}
	if all1-all != userNum2*sponsorNum {
		t.Errorf("Wrong queued-count delta, want %d, have %d",
			userNum2*sponsorNum, all1-all)
	}
}

// Tests there are meta txs and normal txs from the same user
func TestMixedMetaTxs(t *testing.T) {
	t.Parallel()
	// Create the pool to test the pricing enforcement with
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	blockchain := newTestBlockChain(10000000, statedb, new(event.Feed))
	pool := NewTxPool(testTxPoolConfig, eip1559Config, blockchain)
	defer pool.Stop()

	gasFeeSponsorPrivateKey, _ := crypto.GenerateKey()
	gasFeeSponsorAddr := crypto.PubkeyToAddress(gasFeeSponsorPrivateKey.PublicKey)

	userAccPrivateKey, _ := crypto.GenerateKey()
	userAccAddr := crypto.PubkeyToAddress(userAccPrivateKey.PublicKey)

	pool.currentState.AddBalance(gasFeeSponsorAddr, big.NewInt(1e18))
	pool.currentState.AddBalance(userAccAddr, big.NewInt(188e14))

	chainId := params.TestChainConfig.ChainID
	signer := types.LatestSignerForChainID(chainId)
	approveABICallData, _ := hexutil.Decode("0x095ea7b30000000000000000000000001f9090aae28b8a3dceadf281b0f12828e676c3260000000000000000000000000000000000000000000000000de0b6b3a7640000")
	to := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	expireHeight := uint64(20_000_010)

	nonce := uint64(0)

	dynamicTx := &types.DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      nonce,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(2e9),
		Gas:        4700000,
		To:         &to,
		Value:      big.NewInt(0),
		Data:       approveABICallData,
		AccessList: nil,
	}
	tx := types.NewTx(dynamicTx)
	txSignature, err := crypto.Sign(signer.Hash(tx).Bytes(), userAccPrivateKey)
	require.NoError(t, err)
	signedTx, err := tx.WithSignature(signer, txSignature)
	require.NoError(t, err)
	err = pool.AddLocal(signedTx)
	require.NoError(t, err)
	nonce++

	dynamicTx = &types.DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      nonce,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(2e9),
		Gas:        470000,
		To:         &to,
		Value:      big.NewInt(0),
		Data:       approveABICallData,
		AccessList: nil,
	}
	payload, err := generateMetaTxData(dynamicTx, expireHeight, 100, gasFeeSponsorAddr, gasFeeSponsorPrivateKey)
	require.NoError(t, err)
	dynamicTx.Data = payload
	tx = types.NewTx(dynamicTx)
	txSignature, err = crypto.Sign(signer.Hash(tx).Bytes(), userAccPrivateKey)
	require.NoError(t, err)
	signedTx, err = tx.WithSignature(signer, txSignature)
	require.NoError(t, err)
	err = pool.AddLocal(signedTx)
	require.NoError(t, err)
	nonce++

	dynamicTx = &types.DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      nonce,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(2e9),
		Gas:        470000,
		To:         &to,
		Value:      big.NewInt(0),
		Data:       approveABICallData,
		AccessList: nil,
	}
	payload, err = generateMetaTxData(dynamicTx, expireHeight, 100, gasFeeSponsorAddr, gasFeeSponsorPrivateKey)
	require.NoError(t, err)
	dynamicTx.Data = payload
	tx = types.NewTx(dynamicTx)
	txSignature, err = crypto.Sign(signer.Hash(tx).Bytes(), userAccPrivateKey)
	require.NoError(t, err)
	signedTx, err = tx.WithSignature(signer, txSignature)
	require.NoError(t, err)
	err = pool.AddLocal(signedTx)
	require.NoError(t, err)
	nonce++

	pending, _ := pool.Stats()
	require.Equal(t, 3, pending)

	dynamicTx = &types.DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      nonce,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(2e9),
		Gas:        4700000,
		To:         &to,
		Value:      big.NewInt(0),
		Data:       approveABICallData,
		AccessList: nil,
	}
	tx = types.NewTx(dynamicTx)
	txSignature, err = crypto.Sign(signer.Hash(tx).Bytes(), userAccPrivateKey)
	require.NoError(t, err)
	signedTx, err = tx.WithSignature(signer, txSignature)
	require.NoError(t, err)
	err = pool.AddLocal(signedTx)
	require.NoError(t, err)
	nonce++

	pending, _ = pool.Stats()
	require.Equal(t, 4, pending)
}
