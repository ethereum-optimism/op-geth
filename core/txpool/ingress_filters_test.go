package txpool

import (
	"context"
	"errors"
	"math/big"
	"net"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func TestInteropFilter(t *testing.T) {
	// some placeholder transaction to test with
	tx := types.NewTx(&types.DynamicFeeTx{
		ChainID: big.NewInt(1),
		Nonce:   1,
		To:      &common.Address{},
		Value:   big.NewInt(1),
		Data:    []byte{},
	})
	t.Run("Tx has no logs", func(t *testing.T) {
		logFn := func(tx *types.Transaction) ([]*types.Log, uint64, error) {
			return []*types.Log{}, 0, nil
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error {
			// make this return error, but it won't be called because logs are empty
			return errors.New("error")
		}
		// when there are no logs to process, the transaction should be allowed
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx errored when getting logs", func(t *testing.T) {
		logFn := func(tx *types.Transaction) ([]*types.Log, uint64, error) {
			return []*types.Log{}, 0, errors.New("error")
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error {
			// make this return error, but it won't be called because logs retrieval errored
			return errors.New("error")
		}
		// when log retrieval errors, the transaction should be denied
		filter := NewInteropFilter(logFn, checkFn)
		require.False(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx has no executing messages", func(t *testing.T) {
		logFn := func(tx *types.Transaction) ([]*types.Log, uint64, error) {
			l1 := &types.Log{
				Topics: []common.Hash{common.BytesToHash([]byte("topic1"))},
			}
			return []*types.Log{l1}, 0, nil
		}
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error {
			// make this return error, but it won't be called because logs retrieval doesn't have executing messages
			return errors.New("error")
		}
		// when no executing messages are included, the transaction should be allowed
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(context.Background(), tx))
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
		logFn := func(tx *types.Transaction) ([]*types.Log, uint64, error) {
			return []*types.Log{l1}, 0, nil
		}
		var spyEMs []interoptypes.Message
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error {
			spyEMs = ems
			return nil
		}
		// when there is one executing message, the transaction should be allowed
		// if the checkFn returns nil
		filter := NewInteropFilter(logFn, checkFn)
		require.True(t, filter.FilterTx(context.Background(), tx))
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
		logFn := func(tx *types.Transaction) ([]*types.Log, uint64, error) {
			return []*types.Log{l1}, 0, nil
		}
		var spyEMs []interoptypes.Message
		checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error {
			spyEMs = ems
			return errors.New("error")
		}
		// when there is one executing message, and the checkFn returns an error,
		// (ie the supervisor rejects the transaction) the transaction should be denied
		filter := NewInteropFilter(logFn, checkFn)
		require.False(t, filter.FilterTx(context.Background(), tx))
		// confirm that one executing message was passed to the checkFn
		require.Equal(t, 1, len(spyEMs))
	})
}

func TestInteropFilterRPCFailures(t *testing.T) {
	tests := []struct {
		name        string
		networkErr  bool
		timeout     bool
		invalidResp bool
	}{
		{
			name:       "Network Error",
			networkErr: true,
		},
		{
			name:    "Timeout",
			timeout: true,
		},
		{
			name:        "Invalid Response",
			invalidResp: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock log function that always returns our test log
			logFn := func(tx *types.Transaction) ([]*types.Log, uint64, error) {
				log := &types.Log{
					Address: params.InteropCrossL2InboxAddress,
					Topics: []common.Hash{
						common.BytesToHash(interoptypes.ExecutingMessageEventTopic[:]),
						common.BytesToHash([]byte("payloadHash")),
					},
					Data: make([]byte, 32*5),
				}
				return []*types.Log{log}, 0, nil
			}

			// Create mock check function that simulates RPC failures
			checkFn := func(ctx context.Context, ems []interoptypes.Message, safety interoptypes.SafetyLevel, emsTimestamp uint64) error {
				if tt.networkErr {
					return &net.OpError{Op: "dial", Err: errors.New("connection refused")}
				}

				if tt.timeout {
					return context.DeadlineExceeded
				}

				if tt.invalidResp {
					return errors.New("invalid response format")
				}

				return nil
			}

			// Create and test filter
			filter := NewInteropFilter(logFn, checkFn)
			result := filter.FilterTx(context.Background(), &types.Transaction{})
			require.Equal(t, false, result, "FilterTx result mismatch")
		})
	}
}
