package blobpool

import (
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/log"
)

func (pool *BlobPool) recreateCeloProperties() {
	pool.celoBackend = &core.CeloBackend{
		ChainConfig: pool.chain.Config(),
		State:       pool.state,
	}
	currentRates, err := pool.celoBackend.GetExchangeRates()
	if err != nil {
		log.Error("Error trying to get exchange rates in txpool.", "cause", err)
	}
	pool.currentRates = currentRates
}
