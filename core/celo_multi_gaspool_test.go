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
		whitelist           []FeeCurrency
		defaultLimit        float64
		limits              FeeCurrencyLimitMapping
		defaultPoolExpected bool
		expectedValue       uint64
	}{
		{
			name:                "Empty whitelist, empty mapping, CELO uses default pool",
			feeCurrency:         nil,
			whitelist:           []FeeCurrency{},
			defaultLimit:        0.9,
			limits:              map[FeeCurrency]float64{},
			defaultPoolExpected: true,
			expectedValue:       900, // blockGasLimit - subGasAmount
		},
		{
			name:        "Non-empty whitelist, non-empty mapping, CELO uses default pool",
			feeCurrency: nil,
			whitelist: []FeeCurrency{
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
			name:                "Empty whitelist, empty mapping, non-whitelisted currency fallbacks to the default pool",
			feeCurrency:         &cUSDToken,
			whitelist:           []FeeCurrency{},
			defaultLimit:        0.9,
			limits:              map[FeeCurrency]float64{},
			defaultPoolExpected: true,
			expectedValue:       900, // blockGasLimit - subGasAmount
		},
		{
			name:        "Non-empty whitelist, non-empty mapping, non-whitelisted currency uses default pool",
			feeCurrency: &cEURToken,
			whitelist: []FeeCurrency{
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
			name:        "Non-empty whitelist, empty mapping, whitelisted currency uses default limit",
			feeCurrency: &cUSDToken,
			whitelist: []FeeCurrency{
				cUSDToken,
			},
			defaultLimit:        0.9,
			limits:              map[FeeCurrency]float64{},
			defaultPoolExpected: false,
			expectedValue:       800, // blockGasLimit * defaultLimit - subGasAmount
		},
		{
			name:        "Non-empty whitelist, non-empty mapping, configured whitelisted currency uses configured limits",
			feeCurrency: &cUSDToken,
			whitelist: []FeeCurrency{
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
			name:        "Non-empty whitelist, non-empty mapping, unconfigured whitelisted currency uses default limit",
			feeCurrency: &cEURToken,
			whitelist: []FeeCurrency{
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
				c.whitelist,
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
