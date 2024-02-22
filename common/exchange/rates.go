package exchange

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

var (
	unitRate = big.NewRat(1, 1)
	// ErrNonWhitelistedFeeCurrency is returned if the currency specified to use for the fees
	// isn't one of the currencies whitelisted for that purpose.
	ErrNonWhitelistedFeeCurrency = errors.New("non-whitelisted fee currency address")
)

// ConvertCurrency does an exchange conversion from currencyFrom to currencyTo of the value given.
func ConvertCurrency(exchangeRates common.ExchangeRates, val1 *big.Int, currencyFrom *common.Address, currencyTo *common.Address) *big.Int {
	goldAmount, err := ConvertCurrencyToGold(exchangeRates, val1, currencyFrom)
	if err != nil {
		log.Error("Error trying to convert from currency to gold.", "value", val1, "fromCurrency", currencyFrom.Hex())
	}
	toAmount, err := ConvertGoldToCurrency(exchangeRates, currencyTo, goldAmount)
	if err != nil {
		log.Error("Error trying to convert from gold to currency.", "value", goldAmount, "toCurrency", currencyTo.Hex())
	}
	return toAmount
}

func ConvertCurrencyToGold(exchangeRates common.ExchangeRates, currencyAmount *big.Int, feeCurrency *common.Address) (*big.Int, error) {
	if feeCurrency == nil {
		return currencyAmount, nil
	}
	exchangeRate, ok := exchangeRates[*feeCurrency]
	if !ok {
		return nil, ErrNonWhitelistedFeeCurrency
	}
	return new(big.Int).Div(new(big.Int).Mul(currencyAmount, exchangeRate.Denom()), exchangeRate.Num()), nil
}

func ConvertGoldToCurrency(exchangeRates common.ExchangeRates, feeCurrency *common.Address, goldAmount *big.Int) (*big.Int, error) {
	if feeCurrency == nil {
		return goldAmount, nil
	}
	exchangeRate, ok := exchangeRates[*feeCurrency]
	if !ok {
		return nil, ErrNonWhitelistedFeeCurrency
	}
	return new(big.Int).Div(new(big.Int).Mul(goldAmount, exchangeRate.Num()), exchangeRate.Denom()), nil
}

func getRate(exchangeRates common.ExchangeRates, feeCurrency *common.Address) (*big.Rat, error) {
	if feeCurrency == nil {
		return unitRate, nil
	}
	rate, ok := exchangeRates[*feeCurrency]
	if !ok {
		return nil, fmt.Errorf("fee currency not registered: %s", feeCurrency.Hex())
	}
	return rate, nil
}

// CompareValue compares values in different currencies (nil currency is native currency)
// returns -1 0 or 1 depending if val1 < val2, val1 == val2, or val1 > val2 respectively.
func CompareValue(exchangeRates common.ExchangeRates, val1 *big.Int, feeCurrency1 *common.Address, val2 *big.Int, feeCurrency2 *common.Address) (int, error) {
	// Short circuit if the fee currency is the same.
	if feeCurrency1 == feeCurrency2 {
		return val1.Cmp(val2), nil
	}

	exchangeRate1, err := getRate(exchangeRates, feeCurrency1)
	if err != nil {
		return 0, err
	}
	exchangeRate2, err := getRate(exchangeRates, feeCurrency2)
	if err != nil {
		return 0, err
	}

	// Below code block is basically evaluating this comparison:
	// val1 * exchangeRate1.denominator / exchangeRate1.numerator < val2 * exchangeRate2.denominator / exchangeRate2.numerator
	// It will transform that comparison to this, to remove having to deal with fractional values.
	// val1 * exchangeRate1.denominator * exchangeRate2.numerator < val2 * exchangeRate2.denominator * exchangeRate1.numerator
	leftSide := new(big.Int).Mul(
		val1,
		new(big.Int).Mul(
			exchangeRate1.Denom(),
			exchangeRate2.Num(),
		),
	)
	rightSide := new(big.Int).Mul(
		val2,
		new(big.Int).Mul(
			exchangeRate2.Denom(),
			exchangeRate1.Num(),
		),
	)

	return leftSide.Cmp(rightSide), nil
}
