package celoapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/exchange"
	"github.com/ethereum/go-ethereum/contracts"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/rpc"
)

func NewCeloAPIBackend(b ethapi.Backend) *CeloAPIBackend {
	return &CeloAPIBackend{
		Backend: b,
	}
}

// CeloAPIBackend is a wrapper for the ethapi.Backend, that provides additional Celo specific
// functionality. CeloAPIBackend is mainly passed to the JSON RPC services and provides
// an easy way to make extra functionality available in the service internal methods without
// having to change their call signature significantly.
type CeloAPIBackend struct {
	ethapi.Backend
}

func (b *CeloAPIBackend) getContractCaller(ctx context.Context, blockNumOrHash rpc.BlockNumberOrHash) (*contracts.CeloBackend, error) {
	state, _, err := b.Backend.StateAndHeaderByNumberOrHash(
		ctx,
		blockNumOrHash,
	)
	if err != nil {
		return nil, fmt.Errorf("retrieve state for block hash %s: %w", blockNumOrHash.String(), err)
	}
	return &contracts.CeloBackend{
		ChainConfig: b.Backend.ChainConfig(),
		State:       state,
	}, nil
}

func (b *CeloAPIBackend) GetFeeBalance(ctx context.Context, blockNumOrHash rpc.BlockNumberOrHash, account common.Address, feeCurrency *common.Address) (*big.Int, error) {
	cb, err := b.getContractCaller(ctx, blockNumOrHash)
	if err != nil {
		return nil, err
	}
	return contracts.GetFeeBalance(cb, account, feeCurrency), nil
}

func (b *CeloAPIBackend) GetExchangeRates(ctx context.Context, blockNumOrHash rpc.BlockNumberOrHash) (common.ExchangeRates, error) {
	contractBackend, err := b.getContractCaller(ctx, blockNumOrHash)
	if err != nil {
		return nil, err
	}
	er, err := contracts.GetExchangeRates(contractBackend)
	if err != nil {
		return nil, err
	}
	return er, nil
}

func (b *CeloAPIBackend) ConvertToCurrency(ctx context.Context, blockNumOrHash rpc.BlockNumberOrHash, goldAmount *big.Int, toFeeCurrency *common.Address) (*big.Int, error) {
	er, err := b.GetExchangeRates(ctx, blockNumOrHash)
	if err != nil {
		return nil, err
	}
	return exchange.ConvertGoldToCurrency(er, toFeeCurrency, goldAmount)
}

func (b *CeloAPIBackend) ConvertToGold(ctx context.Context, blockNumOrHash rpc.BlockNumberOrHash, currencyAmount *big.Int, fromFeeCurrency *common.Address) (*big.Int, error) {
	er, err := b.GetExchangeRates(ctx, blockNumOrHash)
	if err != nil {
		return nil, err
	}
	return exchange.ConvertCurrencyToGold(er, currencyAmount, fromFeeCurrency)
}
