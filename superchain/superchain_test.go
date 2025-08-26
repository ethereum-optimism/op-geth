package superchain

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestGetSuperchain(t *testing.T) {
	mainnet, err := GetSuperchain("mainnet")
	require.NoError(t, err)

	require.Equal(t, "Mainnet", mainnet.Name)
	require.Equal(t, common.HexToAddress("0x8062AbC286f5e7D9428a0Ccb9AbD71e50d93b935"), mainnet.ProtocolVersionsAddr)
	require.EqualValues(t, 1, mainnet.L1.ChainID)

	_, err = GetSuperchain("not a network")
	require.Error(t, err)
}
