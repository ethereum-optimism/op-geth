package celoapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/exchange"
	"github.com/ethereum/go-ethereum/common/lru"
	"github.com/ethereum/go-ethereum/contracts"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewCeloAPIBackend(b ethapi.Backend) *CeloAPIBackend {
	return &CeloAPIBackend{
		Backend:            b,
		exchangeRatesCache: lru.NewCache[common.Hash, common.ExchangeRates](128),
	}
}

// CeloAPIBackend is a wrapper for the ethapi.Backend, that provides additional Celo specific
// functionality. CeloAPIBackend is mainly passed to the JSON RPC services and provides
// an easy way to make extra functionality available in the service internal methods without
// having to change their call signature significantly.
// CeloAPIBackend keeps a threadsafe LRU cache of block-hash to exchange rates for that block.
// Cache invalidation is only a problem when an already existing blocks' hash
// doesn't change, but the rates change. That shouldn't be possible, since changing the rates
// requires different transaction hashes / state and thus a different block hash.
// If the previous rates change during a reorg, the previous block hash should also change
// and with it the new block's hash.
// Stale branches cache values will get evicted eventually.
type CeloAPIBackend struct {
	ethapi.Backend

	exchangeRatesCache *lru.Cache[common.Hash, common.ExchangeRates]
}

func (b *CeloAPIBackend) getContractCaller(ctx context.Context, atBlock common.Hash) (*contracts.CeloBackend, error) {
	state, _, err := b.Backend.StateAndHeaderByNumberOrHash(
		ctx,
		rpc.BlockNumberOrHashWithHash(atBlock, false),
	)
	if err != nil {
		return nil, fmt.Errorf("retrieve state for block hash %s: %w", atBlock.String(), err)
	}
	return &contracts.CeloBackend{
		ChainConfig: b.Backend.ChainConfig(),
		State:       state,
	}, nil
}

func (b *CeloAPIBackend) GetFeeBalance(ctx context.Context, atBlock common.Hash, account common.Address, feeCurrency *common.Address) (*big.Int, error) {
	cb, err := b.getContractCaller(ctx, atBlock)
	if err != nil {
		return nil, err
	}
	return contracts.GetFeeBalance(cb, account, feeCurrency), nil
}

func (b *CeloAPIBackend) GetExchangeRates(ctx context.Context, atBlock common.Hash) (common.ExchangeRates, error) {
	cachedRates, ok := b.exchangeRatesCache.Get(atBlock)
	if ok {
		return cachedRates, nil
	}
	cb, err := b.getContractCaller(ctx, atBlock)
	if err != nil {
		return nil, err
	}
	er, err := contracts.GetExchangeRates(cb)
	if err != nil {
		return nil, err
	}
	b.exchangeRatesCache.Add(atBlock, er)
	return er, nil
}

func (b *CeloAPIBackend) ConvertToCurrency(ctx context.Context, atBlock common.Hash, value *big.Int, fromFeeCurrency *common.Address) (*big.Int, error) {
	er, err := b.GetExchangeRates(ctx, atBlock)
	if err != nil {
		return nil, err
	}
	return exchange.ConvertCeloToCurrency(er, fromFeeCurrency, value)
}

func (b *CeloAPIBackend) ConvertToGold(ctx context.Context, atBlock common.Hash, value *big.Int, toFeeCurrency *common.Address) (*big.Int, error) {
	er, err := b.GetExchangeRates(ctx, atBlock)
	if err != nil {
		return nil, err
	}
	return exchange.ConvertCurrencyToCelo(er, value, toFeeCurrency)
}
