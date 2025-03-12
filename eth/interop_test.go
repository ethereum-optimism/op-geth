package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestTxToInteropAccessList(t *testing.T) {
	t.Run("Tx has no access list", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{})
		require.Nil(t, txToInteropAccessList(tx))
	})
	t.Run("Tx has access list with no interop address", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			AccessList: types.AccessList{
				{Address: common.Address{0xaa}, StorageKeys: []common.Hash{{0xbb}}},
			},
		})
		require.Nil(t, txToInteropAccessList(tx))
	})
	t.Run("Tx has access list with interop messages", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			AccessList: types.AccessList{
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xaa}}},
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xbb}, {0xcc}}},
			},
		})
		require.Equal(t, []common.Hash{{0xaa}, {0xbb}, {0xcc}}, txToInteropAccessList(tx))
	})
	t.Run("Tx has access list with gaps between interop messages", func(t *testing.T) {
		tx := types.NewTx(&types.DynamicFeeTx{
			AccessList: types.AccessList{
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xaa}}},
				{Address: common.Address{0xcc}, StorageKeys: []common.Hash{{0xcc}}},
				{Address: params.InteropCrossL2InboxAddress, StorageKeys: []common.Hash{{0xbb}}},
			},
		})
		require.Equal(t, []common.Hash{{0xaa}, {0xbb}}, txToInteropAccessList(tx))
	})
}
