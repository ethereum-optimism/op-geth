package types

import (
	"github.com/ethereum/go-ethereum/common/exchange"
)

// CompareWithRates compares the effective gas price of two transactions according to the exchange rates and
// the base fees in the transactions currencies.
func CompareWithRates(a, b *Transaction, ratesAndFees *exchange.RatesAndFees) int {
	if ratesAndFees == nil {
		// During node startup the ratesAndFees might not be yet setup, compare nominally
		feeCapCmp := a.GasFeeCapCmp(b)
		if feeCapCmp != 0 {
			return feeCapCmp
		}
		return a.GasTipCapCmp(b)
	}
	rates := ratesAndFees.Rates
	if ratesAndFees.HasBaseFee() {
		tipA := a.EffectiveGasTipValue(ratesAndFees.GetBaseFeeIn(a.inner.feeCurrency()))
		tipB := b.EffectiveGasTipValue(ratesAndFees.GetBaseFeeIn(b.inner.feeCurrency()))
		c, _ := exchange.CompareValue(rates, tipA, a.inner.feeCurrency(), tipB, b.inner.feeCurrency())
		return c
	}

	// Compare fee caps if baseFee is not specified or effective tips are equal
	feeA := a.inner.gasFeeCap()
	feeB := b.inner.gasFeeCap()
	c, _ := exchange.CompareValue(rates, feeA, a.inner.feeCurrency(), feeB, b.inner.feeCurrency())
	if c != 0 {
		return c
	}

	// Compare tips if effective tips and fee caps are equal
	tipCapA := a.inner.gasTipCap()
	tipCapB := b.inner.gasTipCap()
	c, _ = exchange.CompareValue(rates, tipCapA, a.inner.feeCurrency(), tipCapB, b.inner.feeCurrency())
	return c
}
