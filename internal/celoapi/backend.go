package celoapi

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/exchange"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/contracts"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/internal/ethapi"
)

type Ethereum interface {
	BlockChain() *core.BlockChain
}

type CeloAPI struct {
	ethAPI *ethapi.EthereumAPI
	eth    Ethereum
}

func NewCeloAPI(e Ethereum, b ethapi.Backend) *CeloAPI {
	return &CeloAPI{
		ethAPI: ethapi.NewEthereumAPI(b),
		eth:    e,
	}
}

func (c *CeloAPI) convertedCurrencyValue(v *hexutil.Big, feeCurrency *common.Address) (*hexutil.Big, error) {
	if feeCurrency != nil {
		convertedTipCap, err := c.convertGoldToCurrency(v.ToInt(), feeCurrency)
		if err != nil {
			return nil, fmt.Errorf("convert to feeCurrency: %w", err)
		}
		v = (*hexutil.Big)(convertedTipCap)
	}
	return v, nil
}

func (c *CeloAPI) celoBackendCurrentState() (*contracts.CeloBackend, error) {
	state, err := c.eth.BlockChain().State()
	if err != nil {
		return nil, fmt.Errorf("retrieve HEAD blockchain state': %w", err)
	}

	cb := &contracts.CeloBackend{
		ChainConfig: c.eth.BlockChain().Config(),
		State:       state,
	}
	return cb, nil
}

func (c *CeloAPI) convertGoldToCurrency(nativePrice *big.Int, feeCurrency *common.Address) (*big.Int, error) {
	cb, err := c.celoBackendCurrentState()
	if err != nil {
		return nil, err
	}
	er, err := contracts.GetExchangeRates(cb)
	if err != nil {
		return nil, fmt.Errorf("retrieve exchange rates from current state: %w", err)
	}
	return exchange.ConvertCeloToCurrency(er, feeCurrency, nativePrice)
}

// GasPrice wraps the original JSON RPC `eth_gasPrice` and adds an additional
// optional parameter `feeCurrency` for fee-currency conversion.
// When `feeCurrency` is not given, then the original JSON RPC method is called without conversion.
func (c *CeloAPI) GasPrice(ctx context.Context, feeCurrency *common.Address) (*hexutil.Big, error) {
	tipcap, err := c.ethAPI.GasPrice(ctx)
	if err != nil {
		return nil, err
	}
	// Between the call to `ethapi.GasPrice` and the call to fetch and convert the rates,
	// there is a chance of a state-change. This means that gas-price suggestion is calculated
	// based on state of block x, while the currency conversion could be calculated based on block
	// x+1.
	// However, a similar race condition is present in the `ethapi.GasPrice` method itself.
	return c.convertedCurrencyValue(tipcap, feeCurrency)
}

// MaxPriorityFeePerGas wraps the original JSON RPC `eth_maxPriorityFeePerGas` and adds an additional
// optional parameter `feeCurrency` for fee-currency conversion.
// When `feeCurrency` is not given, then the original JSON RPC method is called without conversion.
func (c *CeloAPI) MaxPriorityFeePerGas(ctx context.Context, feeCurrency *common.Address) (*hexutil.Big, error) {
	tipcap, err := c.ethAPI.MaxPriorityFeePerGas(ctx)
	if err != nil {
		return nil, err
	}
	// Between the call to `ethapi.MaxPriorityFeePerGas` and the call to fetch and convert the rates,
	// there is a chance of a state-change. This means that gas-price suggestion is calculated
	// based on state of block x, while the currency conversion could be calculated based on block
	// x+1.
	return c.convertedCurrencyValue(tipcap, feeCurrency)
}
