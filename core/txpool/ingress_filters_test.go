package txpool

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/types/interoptypes"
	"github.com/stretchr/testify/require"
)

type mockInteropFilterAPI struct {
	timeFn       func() (uint64, error)
	accessListFn func(tx *types.Transaction) []common.Hash
	checkFn      func(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error
}

func (m *mockInteropFilterAPI) CurrentInteropBlockTime() (uint64, error) {
	if m.timeFn != nil {
		return m.timeFn()
	}
	return 0, nil
}

func (m *mockInteropFilterAPI) TxToInteropAccessList(tx *types.Transaction) []common.Hash {
	if m.accessListFn != nil {
		return m.accessListFn(tx)
	}
	return nil
}

func (m *mockInteropFilterAPI) CheckAccessList(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error {
	if m.checkFn != nil {
		return m.checkFn(ctx, inboxEntries, minSafety, ed)
	}
	return nil
}

func TestInteropFilter(t *testing.T) {
	api := &mockInteropFilterAPI{}
	filter := NewInteropFilter(api)
	tx := types.NewTx(&types.DynamicFeeTx{})

	t.Run("Tx has no access list", func(t *testing.T) {
		api.accessListFn = func(tx *types.Transaction) []common.Hash {
			return nil
		}
		require.True(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx errored when checking current interop block time", func(t *testing.T) {
		api.timeFn = func() (uint64, error) {
			return 0, errors.New("error")
		}
		require.True(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx has valid executing message", func(t *testing.T) {
		api.timeFn = func() (uint64, error) {
			return 0, nil
		}
		api.accessListFn = func(tx *types.Transaction) []common.Hash {
			return []common.Hash{{0xaa}}
		}
		api.checkFn = func(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error {
			require.Equal(t, common.Hash{0xaa}, inboxEntries[0])
			return nil
		}
		require.True(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx has invalid executing message", func(t *testing.T) {
		api.timeFn = func() (uint64, error) {
			return 1, nil
		}
		api.accessListFn = func(tx *types.Transaction) []common.Hash {
			return []common.Hash{{0xaa}}
		}
		api.checkFn = func(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error {
			require.Equal(t, common.Hash{0xaa}, inboxEntries[0])
			return errors.New("error")
		}
		require.False(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx has valid executing message equal to than expiry", func(t *testing.T) {
		api.timeFn = func() (uint64, error) {
			expiredT := tx.Time().Add(86400 * time.Second)
			return uint64(expiredT.Unix()), nil
		}
		api.accessListFn = func(tx *types.Transaction) []common.Hash {
			return []common.Hash{{0xaa}}
		}
		api.checkFn = func(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error {
			require.Equal(t, common.Hash{0xaa}, inboxEntries[0])
			return nil
		}
		require.True(t, filter.FilterTx(context.Background(), tx))
	})
	t.Run("Tx has valid executing message older than expiry", func(t *testing.T) {
		api.timeFn = func() (uint64, error) {
			expiredT := tx.Time().Add(86401 * time.Second)
			return uint64(expiredT.Unix()), nil
		}
		api.accessListFn = func(tx *types.Transaction) []common.Hash {
			return []common.Hash{{0xaa}}
		}
		api.checkFn = func(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error {
			require.Equal(t, common.Hash{0xaa}, inboxEntries[0])
			return nil
		}
		require.False(t, filter.FilterTx(context.Background(), tx))
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
			api := &mockInteropFilterAPI{}
			filter := NewInteropFilter(api)
			api.accessListFn = func(tx *types.Transaction) []common.Hash {
				return []common.Hash{{0xaa}}
			}
			api.checkFn = func(ctx context.Context, inboxEntries []common.Hash, minSafety interoptypes.SafetyLevel, ed interoptypes.ExecutingDescriptor) error {
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

			result := filter.FilterTx(context.Background(), &types.Transaction{})
			require.Equal(t, false, result, "FilterTx result mismatch")
		})
	}
}
