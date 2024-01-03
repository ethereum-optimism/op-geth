package types

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestL1CostFunc(t *testing.T) {
	basefee := big.NewInt(1)
	overhead := big.NewInt(1)
	scalar := big.NewInt(1_000_000)

	costFunc0 := newL1CostFunc(basefee, overhead, scalar, false /*isRegolith*/)
	costFunc1 := newL1CostFunc(basefee, overhead, scalar, true)

	// emptyTx is a test tx defined in transaction_test.go
	c0, g0 := costFunc0(emptyTx.RollupCostData()) // pre-Regolith
	c1, g1 := costFunc1(emptyTx.RollupCostData())
	require.Equal(t, big.NewInt(1569), c0)
	require.Equal(t, big.NewInt(1569), g0) // gas-used == fee since scalars are all 1
	require.Equal(t, big.NewInt(481), c1)
	require.Equal(t, big.NewInt(481), g1)
}

func TestExtractGasParams(t *testing.T) {
	regolithTime := uint64(1)
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &regolithTime,
	}

	selector := []byte{0x01, 0x5d, 0x8e, 0xb9}
	uint256 := make([]byte, 32)

	ignored := big.NewInt(1234)
	basefee := big.NewInt(1)
	overhead := big.NewInt(1)
	scalar := big.NewInt(1_000_000)

	data := []byte{}
	data = append(data, selector...)                    // selector
	data = append(data, ignored.FillBytes(uint256)...)  // arg 0
	data = append(data, ignored.FillBytes(uint256)...)  // arg 1
	data = append(data, basefee.FillBytes(uint256)...)  // arg 2
	data = append(data, ignored.FillBytes(uint256)...)  // arg 3
	data = append(data, ignored.FillBytes(uint256)...)  // arg 4
	data = append(data, ignored.FillBytes(uint256)...)  // arg 5
	data = append(data, overhead.FillBytes(uint256)...) // arg 6

	// try to extract from data which has not enough params, should get error.
	_, _, _, err := extractL1GasParams(config, regolithTime, data)
	require.Error(t, err)

	data = append(data, scalar.FillBytes(uint256)...) // arg 7

	// now it should succeed
	_, costFuncPreRegolith, _, err := extractL1GasParams(config, regolithTime-1, data)
	require.NoError(t, err)

	// Function should continue to succeed even with extra data (that just gets ignored) since we
	// have been testing the data size is at least the expected number of bytes instead of exactly
	// the expected number of bytes. It's unclear if this flexibility was intentional, but since
	// it's been in production we shouldn't change this behavior.
	data = append(data, ignored.FillBytes(uint256)...) // extra ignored arg
	_, costFuncRegolith, _, err := extractL1GasParams(config, regolithTime, data)
	require.NoError(t, err)

	c, _ := costFuncPreRegolith(emptyTx.RollupCostData())
	require.Equal(t, big.NewInt(1569), c)

	c, _ = costFuncRegolith(emptyTx.RollupCostData())
	require.Equal(t, big.NewInt(481), c)
}
