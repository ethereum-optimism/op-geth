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
	// ErrUnregisteredFeeCurrency is returned if the currency specified to use for the fees
	// isn't one of the currencies whitelisted for that purpose.
	ErrUnregisteredFeeCurrency = errors.New("unregistered fee-currency address")
)

// ConvertCurrency does an exchange conversion from currencyFrom to currencyTo of the value given.
func ConvertCurrency(exchangeRates common.ExchangeRates, val1 *big.Int, currencyFrom *common.Address, currencyTo *common.Address) *big.Int {
	celoAmount, err := ConvertCurrencyToCelo(exchangeRates, val1, currencyFrom)
	if err != nil {
		log.Error("Error trying to convert from currency to CELO.", "value", val1, "fromCurrency", currencyFrom.Hex())
	}
	toAmount, err := ConvertCeloToCurrency(exchangeRates, currencyTo, celoAmount)
	if err != nil {
		log.Error("Error trying to convert from CELO to currency.", "value", celoAmount, "toCurrency", currencyTo.Hex())
	}
	return toAmount
}

func ConvertCurrencyToCelo(exchangeRates common.ExchangeRates, currencyAmount *big.Int, feeCurrency *common.Address) (*big.Int, error) {
	if feeCurrency == nil {
		return currencyAmount, nil
	}
	if currencyAmount == nil {
		return nil, fmt.Errorf("Can't convert nil amount to CELO.")
	}
	exchangeRate, ok := exchangeRates[*feeCurrency]
	if !ok {
		return nil, fmt.Errorf("could not convert from fee currency to native (fee-currency=%s): %w ", feeCurrency, ErrUnregisteredFeeCurrency)
	}
	return new(big.Int).Div(new(big.Int).Mul(currencyAmount, exchangeRate.Denom()), exchangeRate.Num()), nil
}

func ConvertCeloToCurrency(exchangeRates common.ExchangeRates, feeCurrency *common.Address, celoAmount *big.Int) (*big.Int, error) {
	if feeCurrency == nil {
		return celoAmount, nil
	}
	if celoAmount == nil {
		return nil, fmt.Errorf("Can't convert nil amount to fee currency.")
	}
	exchangeRate, ok := exchangeRates[*feeCurrency]
	if !ok {
		return nil, fmt.Errorf("could not convert from native to fee currency (fee-currency=%s): %w ", feeCurrency, ErrUnregisteredFeeCurrency)
	}
	return new(big.Int).Div(new(big.Int).Mul(celoAmount, exchangeRate.Num()), exchangeRate.Denom()), nil
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

// RatesAndFees holds exchange rates and the basefees expressed in the rates currencies.
type RatesAndFees struct {
	Rates common.ExchangeRates

	nativeBaseFee    *big.Int
	currencyBaseFees map[common.Address]*big.Int
}

// NewRatesAndFees creates a new empty RatesAndFees object.
func NewRatesAndFees(rates common.ExchangeRates, nativeBaseFee *big.Int) *RatesAndFees {
	// While it could be made so that currency basefees are calculated on demand,
	// the low amount of these (usually N < 20)
	return &RatesAndFees{
		Rates:            rates,
		nativeBaseFee:    nativeBaseFee,
		currencyBaseFees: make(map[common.Address]*big.Int, len(rates)),
	}
}

// HasBaseFee returns if the basefee is set.
func (rf *RatesAndFees) HasBaseFee() bool {
	return rf.nativeBaseFee != nil
}

// GetNativeBaseFee returns the basefee in celo currency.
func (rf *RatesAndFees) GetNativeBaseFee() *big.Int {
	return rf.nativeBaseFee
}

// GetBaseFeeIn returns the basefee expressed in the specified currency. Returns nil
// if the currency is not allowlisted.
func (rf *RatesAndFees) GetBaseFeeIn(currency *common.Address) *big.Int {
	// If native currency is being requested, return it
	if currency == nil {
		return rf.nativeBaseFee
	}
	// If a non-native currency is being requested, but it is nil,
	// it means there is no baseFee in this context. Return nil as well.
	if rf.nativeBaseFee == nil {
		return nil
	}
	// Check the cache
	baseFee, ok := rf.currencyBaseFees[*currency]
	if ok {
		return baseFee
	}
	// Not found, calculate
	calculatedBaseFee, err := ConvertCeloToCurrency(rf.Rates, currency, rf.nativeBaseFee)
	if err != nil {
		// Should never happen: error lvl log line
		log.Error("BaseFee requested for unregistered currency",
			"currency", currency.Hex(),
			"exchangeRates", rf.Rates,
			"cause", err)
		return nil
	}
	rf.currencyBaseFees[*currency] = calculatedBaseFee
	return calculatedBaseFee
}
