package txpool

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestInteropFilter(t *testing.T) {
	t.Run("Tx has no logs", func(t *testing.T) {
		logFn := func(tx *types.Transaction) ([]*types.Log, error) {
			return []*types.Log{}, nil
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel) error {
			// make this return error, but it won't be called because logs are empty
			return errors.New("error")
		}
		// when there are no logs to process, the transaction should be allowed
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(&types.Transaction{}))
	})
	t.Run("Tx errored when getting logs", func(t *testing.T) {
		logFn := func(tx *types.Transaction) ([]*types.Log, error) {
			return []*types.Log{}, errors.New("error")
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel) error {
			// make this return error, but it won't be called because logs retrieval errored
			return errors.New("error")
		}
		// when log retrieval errors, the transaction should be allowed
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(&types.Transaction{}))
	})
	t.Run("Tx has no executing messages", func(t *testing.T) {
		logFn := func(tx *types.Transaction) ([]*types.Log, error) {
			l1 := &types.Log{
				Topics: []common.Hash{common.BytesToHash([]byte("topic1"))},
			}
			return []*types.Log{l1}, errors.New("error")
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel) error {
			// make this return error, but it won't be called because logs retrieval doesn't have executing messages
			return errors.New("error")
		}
		// when no executing messages are included, the transaction should be allowed
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(&types.Transaction{}))
	})
	t.Run("Tx has valid executing message", func(t *testing.T) {
		// build a basic executing message
		// the executing message must pass basic decode validation,
		// but the validity check is done by the checkFn
		l1 := &types.Log{
			Address: params.InteropCrossL2InboxAddress,
			Topics: []common.Hash{
				common.BytesToHash(interoptypes.ExecutingMessageEventTopic[:]),
				common.BytesToHash([]byte("payloadHash")),
			},
			Data: []byte{},
		}
		// using all 0s for data allows all takeZeros to pass
		for i := 0; i < 32*5; i++ {
			l1.Data = append(l1.Data, 0)
		}
		logFn := func(tx *types.Transaction) ([]*types.Log, error) {
			return []*types.Log{l1}, nil
		}
		var spyEMs []interoptypes.Message
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel) error {
			spyEMs = ems
			return nil
		}
		// when there is one executing message, the transaction should be allowed
		// if the checkFn returns nil
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(&types.Transaction{}))
		// confirm that one executing message was passed to the checkFn
		require.Equal(t, 1, len(spyEMs))
	})
	t.Run("Tx has invalid executing message", func(t *testing.T) {
		// build a basic executing message
		// the executing message must pass basic decode validation,
		// but the validity check is done by the checkFn
		l1 := &types.Log{
			Address: params.InteropCrossL2InboxAddress,
			Topics: []common.Hash{
				common.BytesToHash(interoptypes.ExecutingMessageEventTopic[:]),
				common.BytesToHash([]byte("payloadHash")),
			},
			Data: []byte{},
		}
		// using all 0s for data allows all takeZeros to pass
		for i := 0; i < 32*5; i++ {
			l1.Data = append(l1.Data, 0)
		}
		logFn := func(tx *types.Transaction) ([]*types.Log, error) {
			return []*types.Log{l1}, nil
		}
		var spyEMs []interoptypes.Message
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel) error {
			spyEMs = ems
			return errors.New("error")
		}
		// when there is one executing message, and the checkFn returns an error,
		// (ie the supervisor rejects the transaction) the transaction should be denied
		filter := NewInteropFilter(logFn, checkFn)
		require.False(t, filter.FilterTx(&types.Transaction{}))
		// confirm that one executing message was passed to the checkFn
		require.Equal(t, 1, len(spyEMs))
	})
}
