//go:build !rocksdb

// Package rocksdb implements the key-value database layer based on RocksDB.
// This is a stub implementation when RocksDB is not available.
package rocksdb

import (
	"errors"

	"github.com/ethereum/go-ethereum/ethdb"
)

var ErrRocksDBNotSupported = errors.New("rocksdb support is not enabled. Compile with 'rocksdb' build tag and ensure RocksDB C++ library is installed")

// Database is a stub for RocksDB when not available
type Database struct{}

// SyncKeyValue implements ethdb.KeyValueStore.
func (db *Database) SyncKeyValue() error {
	return ErrRocksDBNotSupported
}

// New returns an error indicating RocksDB is not supported
func New(file string, cache int, handles int, namespace string, readonly bool) (*Database, error) {
	return nil, ErrRocksDBNotSupported
}

// NewCustom returns an error indicating RocksDB is not supported
func NewCustom(file string, namespace string, customize func(options interface{})) (*Database, error) {
	return nil, ErrRocksDBNotSupported
}

// Stub implementations for interface compliance - these should never be called
func (db *Database) Close() error                                           { return ErrRocksDBNotSupported }
func (db *Database) Has(key []byte) (bool, error)                           { return false, ErrRocksDBNotSupported }
func (db *Database) Get(key []byte) ([]byte, error)                         { return nil, ErrRocksDBNotSupported }
func (db *Database) Put(key []byte, value []byte) error                     { return ErrRocksDBNotSupported }
func (db *Database) Delete(key []byte) error                                { return ErrRocksDBNotSupported }
func (db *Database) DeleteRange(start, end []byte) error                    { return ErrRocksDBNotSupported }
func (db *Database) NewBatch() ethdb.Batch                                  { return &stubBatch{} }
func (db *Database) NewBatchWithSize(size int) ethdb.Batch                  { return &stubBatch{} }
func (db *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator { return &stubIterator{} }
func (db *Database) Stat() (string, error)                                  { return "", ErrRocksDBNotSupported }
func (db *Database) Compact(start []byte, limit []byte) error               { return ErrRocksDBNotSupported }
func (db *Database) Path() string                                           { return "" }

// stubBatch implements ethdb.Batch for the stub
type stubBatch struct{}

// DeleteRange implements ethdb.Batch.
func (b *stubBatch) DeleteRange(start []byte, end []byte) error {
	return ErrRocksDBNotSupported
}

func (b *stubBatch) Put(key, value []byte) error         { return ErrRocksDBNotSupported }
func (b *stubBatch) Delete(key []byte) error             { return ErrRocksDBNotSupported }
func (b *stubBatch) ValueSize() int                      { return 0 }
func (b *stubBatch) Write() error                        { return ErrRocksDBNotSupported }
func (b *stubBatch) Reset()                              {}
func (b *stubBatch) Replay(w ethdb.KeyValueWriter) error { return ErrRocksDBNotSupported }

// stubIterator implements ethdb.Iterator for the stub
type stubIterator struct{}

func (it *stubIterator) Next() bool    { return false }
func (it *stubIterator) Error() error  { return ErrRocksDBNotSupported }
func (it *stubIterator) Key() []byte   { return nil }
func (it *stubIterator) Value() []byte { return nil }
func (it *stubIterator) Release()      {}
