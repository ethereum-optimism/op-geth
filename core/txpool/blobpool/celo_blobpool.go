package blobpool

import (
	"github.com/ethereum/go-ethereum/contracts"
	"github.com/ethereum/go-ethereum/log"
)

func (pool *BlobPool) recreateCeloProperties() {
	pool.celoBackend = &contracts.CeloBackend{
		ChainConfig: pool.chain.Config(),
		State:       pool.state,
	}
	currentRates, err := contracts.GetExchangeRates(pool.celoBackend)
	if err != nil {
		log.Error("Error trying to get exchange rates in txpool.", "cause", err)
	}
	pool.currentRates = currentRates
}
