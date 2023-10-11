package contracts

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/addresses"
	"github.com/ethereum/go-ethereum/contracts/celo/abigen"
	"github.com/ethereum/go-ethereum/log"
)

// GetExchangeRates returns the exchange rates for all gas currencies from CELO
func GetExchangeRates(caller bind.ContractCaller) (common.ExchangeRates, error) {
	exchangeRates := map[common.Address]*big.Rat{}
	directory, err := abigen.NewFeeCurrencyDirectoryCaller(addresses.FeeCurrencyDirectoryAddress, caller)
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to access FeeCurrencyDirectory: %w", err)
	}

	registeredTokens, err := directory.GetCurrencies(&bind.CallOpts{})
	if err != nil {
		return exchangeRates, fmt.Errorf("Failed to get whitelisted tokens: %w", err)
	}
	for _, tokenAddress := range registeredTokens {
		rate, err := directory.GetExchangeRate(&bind.CallOpts{}, tokenAddress)
		if err != nil {
			log.Error("Failed to get medianRate for gas currency!", "err", err, "tokenAddress", tokenAddress.Hex())
			continue
		}
		if rate.Numerator.Sign() <= 0 || rate.Denominator.Sign() <= 0 {
			log.Error("Bad exchange rate for fee currency", "tokenAddress", tokenAddress.Hex(), "numerator", rate.Numerator, "denominator", rate.Denominator)
			continue
		}
		exchangeRates[tokenAddress] = new(big.Rat).SetFrac(rate.Numerator, rate.Denominator)
	}

	return exchangeRates, nil
}
