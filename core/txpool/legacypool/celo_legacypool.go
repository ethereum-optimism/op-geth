package legacypool

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/holiman/uint256"
)

// filter Filters transactions from the given list, according to remaining balance (per currency)
// and gasLimit. Returns drops and invalid txs.
func (pool *LegacyPool) filter(list *list, addr common.Address, gasLimit uint64) (types.Transactions, types.Transactions) {
	// CELO: drop all transactions that no longer have a whitelisted currency
	dropsWhitelist, invalidsWhitelist := list.FilterWhitelisted(pool.currentRates)
	// Check from which currencies we need to get balances
	currenciesInList := list.FeeCurrencies()
	drops, invalids := list.Filter(pool.getBalances(addr, currenciesInList), gasLimit)
	totalDrops := append(dropsWhitelist, drops...)
	totalInvalids := append(invalidsWhitelist, invalids...)
	return totalDrops, totalInvalids
}

func (pool *LegacyPool) getBalances(address common.Address, currencies []common.Address) map[common.Address]*uint256.Int {
	balances := make(map[common.Address]*uint256.Int, len(currencies))
	for _, curr := range currencies {
		balances[curr] = uint256.MustFromBig(pool.celoBackend.GetFeeBalance(address, &curr))
	}
	return balances
}

func (pool *LegacyPool) recreateCeloProperties() {
	pool.celoBackend = &core.CeloBackend{
		ChainConfig: pool.chainconfig,
		State:       pool.currentState,
	}
	currentRates, err := pool.celoBackend.GetExchangeRates()
	if err != nil {
		log.Error("Error trying to get exchange rates in txpool.", "cause", err)
	}
	pool.currentRates = currentRates
}
