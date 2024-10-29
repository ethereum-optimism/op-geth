package txpool

import (
	"context"
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
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
		logFn := func(tx *types.Transaction) ([]*types.Log, error) {
			// TODO: make executing messages here
			l1 := &types.Log{
				Topics: []common.Hash{common.BytesToHash([]byte("topic1"))},
			}
			return []*types.Log{l1}, errors.New("error")
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel) error {
			// make this return error, but it won't be called because logs retrieval doesn't have executing messages
			return nil
		}
		// when no executing messages are included, the transaction should be allowed
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(&types.Transaction{}))
	})
}
