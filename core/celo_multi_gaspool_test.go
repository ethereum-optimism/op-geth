package core

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestMultiCurrencyGasPool(t *testing.T) {
	blockGasLimit := uint64(1_000)
	subGasAmount := 100

	cUSDToken := common.HexToAddress("0x765DE816845861e75A25fCA122bb6898B8B1282a")
	cEURToken := common.HexToAddress("0xD8763CBa276a3738E6DE85b4b3bF5FDed6D6cA73")

	testCases := []struct {
		name                string
		feeCurrency         *FeeCurrency
		allowlist           []FeeCurrency
		defaultLimit        float64
		limits              FeeCurrencyLimitMapping
		defaultPoolExpected bool
		expectedValue       uint64
	}{
		{
			name:                "Empty allowlist, empty mapping, CELO uses default pool",
			feeCurrency:         nil,
			allowlist:           []FeeCurrency{},
			defaultLimit:        0.9,
			limits:              map[FeeCurrency]float64{},
			defaultPoolExpected: true,
			expectedValue:       900, // blockGasLimit - subGasAmount
		},
		{
			name:        "Non-empty allowlist, non-empty mapping, CELO uses default pool",
			feeCurrency: nil,
			allowlist: []FeeCurrency{
				cUSDToken,
			},
			defaultLimit: 0.9,
			limits: map[FeeCurrency]float64{
				cUSDToken: 0.5,
			},
			defaultPoolExpected: true,
			expectedValue:       900, // blockGasLimit - subGasAmount
		},
		{
			name:                "Empty allowlist, empty mapping, non-registered currency fallbacks to the default pool",
			feeCurrency:         &cUSDToken,
			allowlist:           []FeeCurrency{},
			defaultLimit:        0.9,
			limits:              map[FeeCurrency]float64{},
			defaultPoolExpected: true,
			expectedValue:       900, // blockGasLimit - subGasAmount
		},
		{
			name:        "Non-empty allowlist, non-empty mapping, non-registered currency uses default pool",
			feeCurrency: &cEURToken,
			allowlist: []FeeCurrency{
				cUSDToken,
			},
			defaultLimit: 0.9,
			limits: map[FeeCurrency]float64{
				cUSDToken: 0.5,
			},
			defaultPoolExpected: true,
			expectedValue:       900, // blockGasLimit - subGasAmount
		},
		{
			name:        "Non-empty allowlist, empty mapping, registered currency uses default limit",
			feeCurrency: &cUSDToken,
			allowlist: []FeeCurrency{
				cUSDToken,
			},
			defaultLimit:        0.9,
			limits:              map[FeeCurrency]float64{},
			defaultPoolExpected: false,
			expectedValue:       800, // blockGasLimit * defaultLimit - subGasAmount
		},
		{
			name:        "Non-empty allowlist, non-empty mapping, configured registered currency uses configured limits",
			feeCurrency: &cUSDToken,
			allowlist: []FeeCurrency{
				cUSDToken,
			},
			defaultLimit: 0.9,
			limits: map[FeeCurrency]float64{
				cUSDToken: 0.5,
			},
			defaultPoolExpected: false,
			expectedValue:       400, // blockGasLimit * 0.5 - subGasAmount
		},
		{
			name:        "Non-empty allowlist, non-empty mapping, unconfigured registered currency uses default limit",
			feeCurrency: &cEURToken,
			allowlist: []FeeCurrency{
				cUSDToken,
				cEURToken,
			},
			defaultLimit: 0.9,
			limits: map[FeeCurrency]float64{
				cUSDToken: 0.5,
			},
			defaultPoolExpected: false,
			expectedValue:       800, // blockGasLimit * 0.5 - subGasAmount
		},
	}

	for _, c := range testCases {
		t.Run(c.name, func(t *testing.T) {
			mgp := NewMultiGasPool(
				blockGasLimit,
				c.allowlist,
				c.defaultLimit,
				c.limits,
			)

			pool := mgp.PoolFor(c.feeCurrency)
			pool.SubGas(uint64(subGasAmount))

			if c.defaultPoolExpected {
				result := mgp.PoolFor(nil).Gas()
				if result != c.expectedValue {
					t.Error("Default pool expected", c.expectedValue, "got", result)
				}
			} else {
				result := mgp.PoolFor(c.feeCurrency).Gas()

				if result != c.expectedValue {
					t.Error(
						"Expected pool", c.feeCurrency, "value", c.expectedValue,
						"got", result,
					)
				}
			}
		})
	}
}
