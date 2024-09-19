package miner

import (
	"math"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

var (
	feeCurrency1 = common.BigToAddress(big.NewInt(1))
	feeCurrency2 = common.BigToAddress(big.NewInt(2))
	header       = types.Header{Time: 1111111111111}
)

func HeaderAfter(h types.Header, deltaSeconds int64) *types.Header {
	if h.Time > math.MaxInt64 {
		panic("int64 overflow")
	}
	t := int64(h.Time) + deltaSeconds
	if t < 0 {
		panic("uint64 underflow")
	}
	return &types.Header{Time: uint64(t)}
}

func TestBlocklistEviction(t *testing.T) {
	bl := NewAddressBlocklist()
	bl.Add(feeCurrency1, header)

	// latest header is before eviction time
	assert.True(t, bl.IsBlocked(feeCurrency1, HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)))
	// latest header is after eviction time
	assert.False(t, bl.IsBlocked(feeCurrency1, HeaderAfter(header, int64(EvictionTimeoutSeconds)+1)))

	// check filter allowlist removes the currency from the allowlist
	assert.Equal(t, len(bl.FilterAllowlist(
		common.NewAddressSet(feeCurrency1),
		HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)),
	), 0)

	// permanently delete the currency from the blocklist
	bl.Evict(HeaderAfter(header, int64(EvictionTimeoutSeconds)+1))

	// now the currency is removed from the cache, so the currency is not blocked even in earlier headers
	assert.False(t, bl.IsBlocked(feeCurrency1, HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)))

	// check filter allowlist doesn't change the allowlist
	assert.Equal(t, len(bl.FilterAllowlist(
		common.NewAddressSet(feeCurrency1),
		HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)),
	), 1)
}

func TestBlocklistAddAfterEviction(t *testing.T) {
	bl := NewAddressBlocklist()
	bl.Add(feeCurrency1, header)
	bl.Evict(HeaderAfter(header, int64(EvictionTimeoutSeconds)+1))

	header2 := HeaderAfter(header, 10)
	bl.Add(feeCurrency2, *header2)

	// make sure the feeCurrency2 behaves as expected
	assert.True(t, bl.IsBlocked(feeCurrency2, HeaderAfter(*header2, int64(EvictionTimeoutSeconds)-1)))
	assert.False(t, bl.IsBlocked(feeCurrency2, HeaderAfter(*header2, int64(EvictionTimeoutSeconds)+1)))
}

func TestBlocklistRemove(t *testing.T) {
	bl := NewAddressBlocklist()
	bl.Add(feeCurrency1, header)
	bl.Add(feeCurrency2, header)
	bl.Remove(feeCurrency1)

	assert.False(t, bl.IsBlocked(feeCurrency1, HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)))
	assert.True(t, bl.IsBlocked(feeCurrency2, HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)))
}

func TestBlocklistAddAfterRemove(t *testing.T) {
	bl := NewAddressBlocklist()
	bl.Add(feeCurrency1, header)
	bl.Remove(feeCurrency1)
	assert.False(t, bl.IsBlocked(feeCurrency1, HeaderAfter(header, int64(EvictionTimeoutSeconds)-1)))

	header2 := HeaderAfter(header, 10)
	bl.Add(feeCurrency2, *header2)

	// make sure the feeCurrency2 behaves as expected
	assert.True(t, bl.IsBlocked(feeCurrency2, HeaderAfter(*header2, int64(EvictionTimeoutSeconds)-1)))
	assert.False(t, bl.IsBlocked(feeCurrency2, HeaderAfter(*header2, int64(EvictionTimeoutSeconds)+1)))
}
