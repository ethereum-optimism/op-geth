package vm

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
)

var (
	ErrCheckPoint = fmt.Errorf("no existing checkpoints to remove")
)

var (
	NilAddress = common.BytesToHash(make([]byte, 32))
)

type JournalEntry struct {
	Address common.Address
	Key     common.Hash
	Prev    common.Hash
}

type Journal = []*JournalEntry

type CheckPoints = []int

type TransientStorage struct {
	current map[common.Address]map[common.Hash]common.Hash

	journal     Journal
	checkPoints CheckPoints
}

func NewTransientStorage() *TransientStorage {
	return &TransientStorage{
		current:     make(map[common.Address]map[common.Hash]common.Hash),
		checkPoints: make([]int, 0),
	}
}

func (ts *TransientStorage) Set(address common.Address, key common.Hash, value common.Hash) {
	ts.journal = append(ts.journal,
		&JournalEntry{
			Address: address,
			Key:     key,
			Prev:    ts.Get(address, key),
		},
	)

	if _, exists := ts.current[address]; !exists {
		ts.current[address] = make(map[common.Hash]common.Hash)
	}

	ts.current[address][key] = value
}

func (ts *TransientStorage) Get(address common.Address, key common.Hash) common.Hash {
	_, exists := ts.current[address]

	if !exists {
		ts.current[address] = make(map[common.Hash]common.Hash)
	}

	value, exists := ts.current[address][key]

	if !exists {
		return common.BytesToHash(make([]byte, 32))
	}

	return value
}

func (ts *TransientStorage) CheckPoint() {
	ts.checkPoints = append(ts.checkPoints, 0)
	copy(ts.checkPoints[1:], ts.checkPoints)
	ts.checkPoints[0] = len(ts.checkPoints)

}

func (ts *TransientStorage) Commit() error {
	if len(ts.checkPoints) == 0 {
		return ErrCheckPoint
	}
	// Pop Committment
	ts.checkPoints = ts.checkPoints[1:len(ts.checkPoints)]

	return nil
}

func (ts *TransientStorage) Revert() error {
	if len(ts.checkPoints) == 0 {
		return ErrCheckPoint
	}

	if len(ts.journal) == 0 {
		return nil
	}

	recentCheckPoint := ts.checkPoints[0]
	ts.checkPoints = ts.checkPoints[1:len(ts.checkPoints)]

	for i := len(ts.journal) - 1; i > recentCheckPoint; i-- {
		entry := ts.journal[i]

		ts.current[entry.Address][entry.Key] = entry.Prev

	}

	ts.journal = ts.journal[recentCheckPoint : len(ts.journal)-recentCheckPoint]

	return nil
}
