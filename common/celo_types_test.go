package common

import (
	"math/big"
	"testing"
)

var (
	currA         = HexToAddress("0xA")
	currB         = HexToAddress("0xB")
	currX         = HexToAddress("0xF")
	exchangeRates = ExchangeRates{
		currA: big.NewRat(47, 100),
		currB: big.NewRat(45, 100),
	}
)

func TestIsWhitelisted(t *testing.T) {
	tests := []struct {
		name        string
		feeCurrency *Address
		want        bool
	}{
		{
			name:        "no fee currency",
			feeCurrency: nil,
			want:        true,
		},
		{
			name:        "valid fee currency",
			feeCurrency: &currA,
			want:        true,
		},
		{
			name:        "invalid fee currency",
			feeCurrency: &currX,
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCurrencyWhitelisted(exchangeRates, tt.feeCurrency); got != tt.want {
				t.Errorf("IsWhitelisted() = %v, want %v", got, tt.want)
			}
		})
	}
}
