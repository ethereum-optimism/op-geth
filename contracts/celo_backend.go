package contracts

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/celo/abigen"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
)

// CeloBackend provide a partial ContractBackend implementation, so that we can
// access core contracts during block processing.
type CeloBackend struct {
	ChainConfig *params.ChainConfig
	State       vm.StateDB
}

// ContractCaller implementation

func (b *CeloBackend) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	return b.State.GetCode(contract), nil
}

func (b *CeloBackend) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	// Ensure message is initialized properly.
	if call.Gas == 0 {
		// Chosen to be the same as ethconfig.Defaults.RPCGasCap
		call.Gas = 50000000
	}
	if call.Value == nil {
		call.Value = new(big.Int)
	}

	// Minimal initialization, might need to be extended when CeloBackend
	// is used in more places. Also initializing blockNumber and time with
	// 0 works now, but will break once we add hardforks at a later time.
	if blockNumber == nil {
		blockNumber = common.Big0
	}
	blockCtx := vm.BlockContext{BlockNumber: blockNumber, Time: 0}
	txCtx := vm.TxContext{}
	vmConfig := vm.Config{}

	readOnlyStateDB := ReadOnlyStateDB{StateDB: b.State}
	evm := vm.NewEVM(blockCtx, txCtx, &readOnlyStateDB, b.ChainConfig, vmConfig)
	ret, _, err := evm.StaticCall(vm.AccountRef(evm.Origin), *call.To, call.Data, call.Gas)

	return ret, err
}

// GetBalanceERC20 returns an account's balance on a given ERC20 currency
func (b *CeloBackend) GetBalanceERC20(accountOwner common.Address, contractAddress common.Address) (result *big.Int, err error) {
	token, err := abigen.NewFeeCurrencyCaller(contractAddress, b)
	if err != nil {
		return nil, fmt.Errorf("failed to access FeeCurrency: %w", err)
	}

	balance, err := token.BalanceOf(&bind.CallOpts{}, accountOwner)
	if err != nil {
		return nil, err
	}

	return balance, nil
}

// GetFeeBalance returns the account's balance from the specified feeCurrency
// (if feeCurrency is nil or ZeroAddress, native currency balance is returned).
func (b *CeloBackend) GetFeeBalance(account common.Address, feeCurrency *common.Address) *big.Int {
	if feeCurrency == nil || *feeCurrency == common.ZeroAddress {
		return b.State.GetBalance(account).ToBig()
	}
	balance, err := b.GetBalanceERC20(account, *feeCurrency)
	if err != nil {
		log.Error("Error while trying to get ERC20 balance:", "cause", err, "contract", feeCurrency.Hex(), "account", account.Hex())
	}
	return balance
}

// GetExchangeRates returns the exchange rates for all gas currencies from CELO
func (b *CeloBackend) GetExchangeRates() (common.ExchangeRates, error) {
	exchangeRates := map[common.Address]*big.Rat{}
	whitelist, err := abigen.NewFeeCurrencyWhitelistCaller(FeeCurrencyWhitelistAddress, b)
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to access FeeCurrencyWhitelist: %w", err)
	}
	oracle, err := abigen.NewSortedOraclesCaller(SortedOraclesAddress, b)
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
			log.Error("Failed to get medianRate for gas currency!", "err", err, "tokenAddress", tokenAddress.Hex())
			continue
		}
		if denominator.Sign() == 0 {
			log.Error("Bad exchange rate for fee currency", "tokenAddress", tokenAddress.Hex(), "numerator", numerator, "denominator", denominator)
			continue
		}
		exchangeRates[tokenAddress] = big.NewRat(numerator.Int64(), denominator.Int64())
	}

	return exchangeRates, nil
}
