package core

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	contracts "github.com/ethereum/go-ethereum/contracts/celo"
	"github.com/ethereum/go-ethereum/contracts/celo/abigen"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// Returns the exchange rates for all gas currencies from CELO
func getExchangeRates(caller *CeloBackend) (map[common.Address]*big.Rat, error) {
	exchangeRates := map[common.Address]*big.Rat{}
	whitelist, err := abigen.NewFeeCurrencyWhitelistCaller(contracts.FeeCurrencyWhitelistAddress, caller)
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to access FeeCurrencyWhitelist: %w", err)
	}
	oracle, err := abigen.NewSortedOraclesCaller(contracts.SortedOraclesAddress, caller)
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to access SortedOracle: %w", err)
	}

	whitelistedTokens, err := whitelist.GetWhitelist(&bind.CallOpts{})
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to get whitelisted tokens: %w", err)
	}
	for _, tokenAddress := range whitelistedTokens {
		numerator, denominator, err := oracle.MedianRate(&bind.CallOpts{}, tokenAddress)
		if err != nil {
			log.Error("Failed to get medianRate for gas currency!", "err", err, "tokenAddress", tokenAddress)
			continue
		}
		if denominator.Sign() == 0 {
			log.Error("Bad exchange rate for fee currency", "tokenAddress", tokenAddress, "numerator", numerator, "denominator", denominator)
			continue
		}
		exchangeRates[tokenAddress] = big.NewRat(numerator.Int64(), denominator.Int64())
	}

	return exchangeRates, nil
}

func setCeloFieldsInBlockContext(blockContext *vm.BlockContext, header *types.Header, config *params.ChainConfig, statedb vm.StateDB) {
	if !config.IsCel2(header.Time) {
		return
	}

	caller := &CeloBackend{config, statedb}

	// Add fee currency exchange rates
	var err error
	blockContext.ExchangeRates, err = getExchangeRates(caller)
	if err != nil {
		log.Error("Error fetching exchange rates!", "err", err)
	}
}
