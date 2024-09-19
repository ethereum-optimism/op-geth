package miner

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

const hours uint64 = 60 * 60

var EvictionTimeoutSeconds uint64 = 2 * hours

type AddressBlocklist struct {
	mux        *sync.RWMutex
	currencies map[common.Address]*types.Header
	// fee-currencies blocked at headers with an older timestamp
	// will get evicted when evict() is called
	headerEvictionTimeoutSeconds uint64
	oldestHeader                 *types.Header
}

func NewAddressBlocklist() *AddressBlocklist {
	return &AddressBlocklist{
		mux:                          &sync.RWMutex{},
		currencies:                   map[common.Address]*types.Header{},
		headerEvictionTimeoutSeconds: EvictionTimeoutSeconds,
		oldestHeader:                 nil,
	}
}

func (b *AddressBlocklist) FilterAllowlist(allowlist common.AddressSet, latest *types.Header) common.AddressSet {
	b.mux.RLock()
	defer b.mux.RUnlock()

	filtered := common.AddressSet{}
	for a := range allowlist {
		if !b.isBlocked(a, latest) {
			filtered[a] = struct{}{}
		}
	}
	return filtered
}

func (b *AddressBlocklist) IsBlocked(currency common.Address, latest *types.Header) bool {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.isBlocked(currency, latest)
}

func (b *AddressBlocklist) Remove(currency common.Address) bool {
	b.mux.Lock()
	defer b.mux.Unlock()

	h, ok := b.currencies[currency]
	if !ok {
		return false
	}
	delete(b.currencies, currency)
	if b.oldestHeader.Time >= h.Time {
		b.resetOldestHeader()
	}
	return ok
}

func (b *AddressBlocklist) Add(currency common.Address, head types.Header) {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.oldestHeader == nil || b.oldestHeader.Time > head.Time {
		b.oldestHeader = &head
	}
	b.currencies[currency] = &head
}

func (b *AddressBlocklist) Evict(latest *types.Header) []common.Address {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.evict(latest)
}

func (b *AddressBlocklist) resetOldestHeader() {
	if len(b.currencies) == 0 {
		b.oldestHeader = nil
		return
	}
	for _, v := range b.currencies {
		if b.oldestHeader == nil {
			b.oldestHeader = v
			continue
		}
		if v.Time < b.oldestHeader.Time {
			b.oldestHeader = v
		}
	}
}

func (b *AddressBlocklist) evict(latest *types.Header) []common.Address {
	evicted := []common.Address{}
	if latest == nil {
		return evicted
	}

	if b.oldestHeader == nil || !b.headerEvicted(b.oldestHeader, latest) {
		// nothing set yet
		return evicted
	}
	for feeCurrencyAddress, addedHeader := range b.currencies {
		if b.headerEvicted(addedHeader, latest) {
			delete(b.currencies, feeCurrencyAddress)
			evicted = append(evicted, feeCurrencyAddress)
		}
	}
	b.resetOldestHeader()
	return evicted
}

func (b *AddressBlocklist) headerEvicted(h, latest *types.Header) bool {
	return h.Time+b.headerEvictionTimeoutSeconds < latest.Time
}

func (b *AddressBlocklist) isBlocked(currency common.Address, latest *types.Header) bool {
	h, exists := b.currencies[currency]
	if !exists {
		return false
	}
	if latest == nil {
		// if no latest block provided to check eviction,
		// assume the currency is blocked
		return true
	}
	return !b.headerEvicted(h, latest)
}
