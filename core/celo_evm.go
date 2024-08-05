package core

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// XXYXX
func setCeloFieldsInBlockContext(blockContext *vm.BlockContext, header *types.Header, config *params.ChainConfig, statedb vm.StateDB) {
	if !config.IsCel2(header.Time) {
		return
	}

	caller := &contracts.CeloBackend{ChainConfig: config, State: statedb}

	feeCurrencyContext, err := contracts.GetFeeCurrencyContext(caller)
	if err != nil {
		log.Error("Error fetching exchange rates!", "err", err)
	}
	blockContext.FeeCurrencyContext = feeCurrencyContext
}

func GetExchangeRates(header *types.Header, config *params.ChainConfig, statedb vm.StateDB) common.ExchangeRates {
	if !config.IsCel2(header.Time) {
		return nil
	}

	caller := &contracts.CeloBackend{ChainConfig: config, State: statedb}

	feeCurrencyContext, err := contracts.GetFeeCurrencyContext(caller)
	if err != nil {
		log.Error("Error fetching exchange rates!", "err", err)
	}
	return feeCurrencyContext.ExchangeRates
}
