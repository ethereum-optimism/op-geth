package interoptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func TestSafetyLevel(t *testing.T) {
	require.True(t, Invalid.wellFormatted())
	require.True(t, Unsafe.wellFormatted())
	require.True(t, CrossUnsafe.wellFormatted())
	require.True(t, LocalSafe.wellFormatted())
	require.True(t, Safe.wellFormatted())
	require.True(t, Finalized.wellFormatted())
	require.False(t, SafetyLevel("hello").wellFormatted())
	require.False(t, SafetyLevel("").wellFormatted())
}

func TestTxToInteropAccessList(t *testing.T) {
	t.Run("Tx has no access list", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{})
		require.Nil(t, TxToInteropAccessList(tx))
	})
	t.Run("Tx has access list with no interop address", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			AccessList: types.AccessList{
				{Address: common.Address{0xaa}, StorageKeys: []common.Hash{{0xbb}}},
			},
		})
		require.Nil(t, TxToInteropAccessList(tx))
	})
	t.Run("Tx has access list with interop messages", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			AccessList: types.AccessList{
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xaa}}},
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xbb}, {0xcc}}},
			},
		})
		require.Equal(t, []common.Hash{{0xaa}, {0xbb}, {0xcc}}, TxToInteropAccessList(tx))
	})
	t.Run("Tx has access list with gaps between interop messages", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			AccessList: types.AccessList{
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xaa}}},
				{Address: common.Address{0xcc}, StorageKeys: []common.Hash{{0xcc}}},
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xbb}}},
			},
		})
		require.Equal(t, []common.Hash{{0xaa}, {0xbb}}, TxToInteropAccessList(tx))
	})
}
