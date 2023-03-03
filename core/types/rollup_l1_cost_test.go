package types

import (
	"math/rand"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestRollupGasData(t *testing.T) {
	for i := 0; i < 100; i++ {
		zeroes := rand.Uint64()
		ones := rand.Uint64()

		r := RollupGasData{
			Zeroes: zeroes,
			Ones:   ones,
		}
		time := uint64(1)
		cfg := &params.ChainConfig{
			RegolithTime: &time,
		}
		gasPreRegolith := r.DataGas(0, cfg)
		gasPostRegolith := r.DataGas(1, cfg)

		require.Equal(t, r.Zeroes*params.TxDataZeroGas+(r.Ones+68)*params.TxDataNonZeroGasEIP2028, gasPreRegolith)
		require.Equal(t, r.Zeroes*params.TxDataZeroGas+r.Ones*params.TxDataNonZeroGasEIP2028, gasPostRegolith)
	}
}
