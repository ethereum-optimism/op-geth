package legacypool

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
	"github.com/stretchr/testify/assert"
)

func txC(nonce int, feeCap int, tipCap int, gas int, currency *common.Address) *types.Transaction {
	return types.NewTx(&types.CeloDynamicFeeTxV2{
		GasFeeCap:   big.NewInt(int64(feeCap)),
		GasTipCap:   big.NewInt(int64(tipCap)),
		FeeCurrency: currency,
		Gas:         uint64(gas),
		Nonce:       uint64(nonce),
	})
}

func TestListFeeCost(t *testing.T) {
	curr1 := common.HexToAddress("0002")
	curr2 := common.HexToAddress("0004")
	curr3 := common.HexToAddress("0006")
	rates := common.ExchangeRates{
		curr1: big.NewRat(2, 1),
		curr2: big.NewRat(4, 1),
		curr3: big.NewRat(6, 1),
	}
	// Insert the transactions in a random order
	list := newList(false)

	list.Add(txC(7, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	assert.Equal(t, uint64(10000), list.TotalCostFor(&curr1).Uint64())

	toBeRemoved := txC(8, 2, 1, 15000, &curr2)
	list.Add(toBeRemoved, DefaultConfig.PriceBump, nil, rates)
	assert.Equal(t, uint64(30000), list.TotalCostFor(&curr2).Uint64())
	assert.Equal(t, uint64(10000), list.TotalCostFor(&curr1).Uint64())

	list.Add(txC(9, 3, 2, 5000, &curr3), DefaultConfig.PriceBump, nil, rates)
	assert.Equal(t, uint64(15000), list.TotalCostFor(&curr3).Uint64())
	assert.Equal(t, uint64(30000), list.TotalCostFor(&curr2).Uint64())
	assert.Equal(t, uint64(10000), list.TotalCostFor(&curr1).Uint64())

	// Add another tx from curr1, check it adds properly
	list.Add(txC(10, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	assert.Equal(t, uint64(15000), list.TotalCostFor(&curr3).Uint64())
	assert.Equal(t, uint64(30000), list.TotalCostFor(&curr2).Uint64())
	assert.Equal(t, uint64(20000), list.TotalCostFor(&curr1).Uint64())

	// Remove a tx from curr2, check it subtracts properly
	removed, _ := list.Remove(toBeRemoved)
	assert.True(t, removed)

	assert.Equal(t, uint64(15000), list.TotalCostFor(&curr3).Uint64())
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr2).Uint64())
	assert.Equal(t, uint64(20000), list.TotalCostFor(&curr1).Uint64())
}

func TestFilterWhitelisted(t *testing.T) {
	curr1 := common.HexToAddress("0002")
	curr2 := common.HexToAddress("0004")
	curr3 := common.HexToAddress("0006")
	rates := common.ExchangeRates{
		curr1: big.NewRat(2, 1),
		curr2: big.NewRat(4, 1),
		curr3: big.NewRat(6, 1),
	}

	list := newList(false)
	list.Add(txC(7, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	toBeRemoved := txC(8, 2, 1, 15000, &curr2)
	list.Add(toBeRemoved, DefaultConfig.PriceBump, nil, rates)
	list.Add(txC(9, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	assert.Equal(t, uint64(30000), list.TotalCostFor(&curr2).Uint64())

	removed, invalids := list.FilterWhitelisted(common.ExchangeRates{curr1: nil, curr3: nil})
	assert.Len(t, removed, 1)
	assert.Len(t, invalids, 0)
	assert.Equal(t, removed[0], toBeRemoved)
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr2).Uint64())
}

func TestFilterWhitelistedStrict(t *testing.T) {
	curr1 := common.HexToAddress("0002")
	curr2 := common.HexToAddress("0004")
	curr3 := common.HexToAddress("0006")
	rates := common.ExchangeRates{
		curr1: big.NewRat(2, 1),
		curr2: big.NewRat(4, 1),
		curr3: big.NewRat(6, 1),
	}

	list := newList(true)
	list.Add(txC(7, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	toBeRemoved := txC(8, 2, 1, 15000, &curr2)
	list.Add(toBeRemoved, DefaultConfig.PriceBump, nil, rates)
	toBeInvalid := txC(9, 1, 1, 10000, &curr3)
	list.Add(toBeInvalid, DefaultConfig.PriceBump, nil, rates)

	removed, invalids := list.FilterWhitelisted(common.ExchangeRates{curr1: nil, curr3: nil})
	assert.Len(t, removed, 1)
	assert.Len(t, invalids, 1)
	assert.Equal(t, removed[0], toBeRemoved)
	assert.Equal(t, invalids[0], toBeInvalid)
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr2).Uint64())
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr3).Uint64())
	assert.Equal(t, uint64(10000), list.TotalCostFor(&curr1).Uint64())
}

func TestFilterBalance(t *testing.T) {
	curr1 := common.HexToAddress("0002")
	curr2 := common.HexToAddress("0004")
	curr3 := common.HexToAddress("0006")
	rates := common.ExchangeRates{
		curr1: big.NewRat(2, 1),
		curr2: big.NewRat(4, 1),
		curr3: big.NewRat(6, 1),
	}

	list := newList(false)
	// each tx costs 10000 in each currency
	list.Add(txC(7, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	toBeRemoved := txC(8, 1, 1, 10000, &curr2)
	list.Add(toBeRemoved, DefaultConfig.PriceBump, nil, rates)
	list.Add(txC(9, 1, 1, 10000, &curr3), DefaultConfig.PriceBump, nil, rates)

	removed, invalids := list.Filter(map[common.Address]*uint256.Int{
		curr1: uint256.NewInt(10000),
		curr2: uint256.NewInt(9999),
		curr3: uint256.NewInt(10000),
	}, 15000)
	assert.Len(t, removed, 1)
	assert.Len(t, invalids, 0)
	assert.Equal(t, removed[0], toBeRemoved)
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr2).Uint64())
}

func TestFilterBalanceStrict(t *testing.T) {
	curr1 := common.HexToAddress("0002")
	curr2 := common.HexToAddress("0004")
	curr3 := common.HexToAddress("0006")
	rates := common.ExchangeRates{
		curr1: big.NewRat(2, 1),
		curr2: big.NewRat(4, 1),
		curr3: big.NewRat(6, 1),
	}

	list := newList(true)
	// each tx costs 10000 in each currency
	list.Add(txC(7, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	toBeRemoved := txC(8, 1, 1, 10000, &curr2)
	list.Add(toBeRemoved, DefaultConfig.PriceBump, nil, rates)
	toBeInvalid := txC(9, 1, 1, 10000, &curr3)
	list.Add(toBeInvalid, DefaultConfig.PriceBump, nil, rates)

	removed, invalids := list.Filter(map[common.Address]*uint256.Int{
		curr1: uint256.NewInt(10001),
		curr2: uint256.NewInt(9999),
		curr3: uint256.NewInt(10001),
	}, 15000)
	assert.Len(t, removed, 1)
	assert.Len(t, invalids, 1)
	assert.Equal(t, removed[0], toBeRemoved)
	assert.Equal(t, invalids[0], toBeInvalid)
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr2).Uint64())
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr3).Uint64())
}

func TestFilterBalanceGasLimit(t *testing.T) {
	curr1 := common.HexToAddress("0002")
	curr2 := common.HexToAddress("0004")
	curr3 := common.HexToAddress("0006")
	rates := common.ExchangeRates{
		curr1: big.NewRat(2, 1),
		curr2: big.NewRat(4, 1),
		curr3: big.NewRat(6, 1),
	}

	list := newList(false)
	// each tx costs 10000 in each currency
	list.Add(txC(7, 1, 1, 10000, &curr1), DefaultConfig.PriceBump, nil, rates)
	toBeRemoved := txC(8, 1, 1, 10001, &curr2)
	list.Add(toBeRemoved, DefaultConfig.PriceBump, nil, rates)
	list.Add(txC(9, 1, 1, 10000, &curr3), DefaultConfig.PriceBump, nil, rates)

	removed, invalids := list.Filter(map[common.Address]*uint256.Int{
		curr1: uint256.NewInt(20000),
		curr2: uint256.NewInt(20000),
		curr3: uint256.NewInt(20000),
	}, 10000)
	assert.Len(t, removed, 1)
	assert.Len(t, invalids, 0)
	assert.Equal(t, removed[0], toBeRemoved)
	assert.Equal(t, uint64(0), list.TotalCostFor(&curr2).Uint64())
}
