package common

import (
	"encoding/json"
	"fmt"
	"math/big"
)

var (
	ZeroAddress = BytesToAddress([]byte{})
)

type AddressSet map[Address]struct{}

type ExchangeRates = map[Address]*big.Rat
type IntrinsicGasCosts = map[Address]uint64

type FeeCurrencyContext struct {
	ExchangeRates     ExchangeRates
	IntrinsicGasCosts IntrinsicGasCosts
}

// Only used in tracer tests
func (fc *FeeCurrencyContext) UnmarshalJSON(data []byte) error {
	var raw struct {
		ExchangeRates     map[Address][]json.Number `json:"exchangeRates"`
		IntrinsicGasCosts map[Address]uint64        `json:"intrinsicGasCosts"`
	}

	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	fc.ExchangeRates = make(ExchangeRates)
	for addr, rateArr := range raw.ExchangeRates {
		if len(rateArr) != 2 {
			return fmt.Errorf("invalid exchange rate array for address %s: expected 2 elements, got %d", addr, len(rateArr))
		}
		numerator, ok := new(big.Int).SetString(string(rateArr[0]), 10)
		if !ok {
			return fmt.Errorf("invalid numerator for address %s: %s", addr, rateArr[0])
		}
		denominator, ok := new(big.Int).SetString(string(rateArr[1]), 10)
		if !ok {
			return fmt.Errorf("invalid denominator for address %s: %s", addr, rateArr[1])
		}

		rate := new(big.Rat).SetFrac(numerator, denominator)
		fc.ExchangeRates[addr] = rate
	}
	fc.IntrinsicGasCosts = raw.IntrinsicGasCosts

	return nil
}

func NewAddressSet(addresses ...Address) AddressSet {
	as := AddressSet{}
	for _, address := range addresses {
		as[address] = struct{}{}
	}
	return as
}

func MaxAllowedIntrinsicGasCost(i IntrinsicGasCosts, feeCurrency *Address) (uint64, bool) {
	intrinsicGas, ok := CurrencyIntrinsicGasCost(i, feeCurrency)
	if !ok {
		return 0, false
	}
	// Allow the contract to overshoot 2 times the deducted intrinsic gas
	// during execution.
	// If the feeCurrency is nil, then the max allowed intrinsic gas cost
	// is 0 (i.e. not allowed) for a fee-currency specific EVM call within the STF.
	return intrinsicGas * 3, true
}

func CurrencyIntrinsicGasCost(i IntrinsicGasCosts, feeCurrency *Address) (uint64, bool) {
	// the additional intrinsic gas cost for a non fee-currency
	// transaction is 0
	if feeCurrency == nil {
		return 0, true
	}
	gasCost, ok := i[*feeCurrency]
	if !ok {
		return 0, false
	}
	return gasCost, true
}

func CurrencyAllowlist(exchangeRates ExchangeRates) AddressSet {
	addrs := AddressSet{}
	for k := range exchangeRates {
		addrs[k] = struct{}{}
	}
	return addrs
}

func IsCurrencyAllowed(exchangeRates ExchangeRates, feeCurrency *Address) bool {
	if feeCurrency == nil {
		return true
	}

	// Check if fee currency is registered
	_, ok := exchangeRates[*feeCurrency]
	return ok
}

func AreSameAddress(a, b *Address) bool {
	// both are nil or point to the same address
	if a == b {
		return true
	}
	// if only one is nil
	if a == nil || b == nil {
		return false
	}
	// if they point to the same
	return *a == *b
}
