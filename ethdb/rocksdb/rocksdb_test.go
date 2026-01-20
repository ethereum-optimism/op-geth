// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package pebble implements the key-value database layer based on pebble.
//go:build !js && !wasip1 && rocksdb

package rocksdb

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/dbtest"
)

func TestRocksDBSuite(t *testing.T) {
	t.Run("DatabaseSuite", func(t *testing.T) {
		dbtest.TestDatabaseSuite(t, func() ethdb.KeyValueStore {
			db, err := New(fmt.Sprintf("/tmp/test-rocksdb-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
			if err != nil {
				t.Fatal(err)
			}
			return db
		})
	})
}

// Before running this benchmark, mount a tmpfs:
// sudo mkdir -p /mnt/tmpfs
// sudo mount -t tmpfs -o size=16G tmpfs /mnt/tmpfs
// Then run:
// go test -benchmem -run=^$ -tags rocksdb -bench ^BenchmarkRocksDB$ github.com/ethereum/go-ethereum/ethdb/rocksdb
func BenchmarkRocksDB(b *testing.B) {
	dbtest.BenchDatabaseSuite(b, func() ethdb.KeyValueStore {
		db, err := New(fmt.Sprintf("/mnt/tmpfs/bench-rocksdb-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
		if err != nil {
			b.Fatal(err)
		}
		return db
	})
}

func BenchmarkRocksDBDisk(b *testing.B) {
	dbtest.BenchDatabaseSuite(b, func() ethdb.KeyValueStore {
		db, err := New(fmt.Sprintf("/tmp/bench-rocksdb-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
		if err != nil {
			b.Fatal(err)
		}
		return db
	})
}

func TestRocksDBBasic(t *testing.T) {
	db, err := New(fmt.Sprintf("/tmp/test-rocksdb-basic-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(db.Path())
	}()

	key := []byte("key123")
	value := []byte("value123")

	if err := db.Put(key, value); err != nil {
		t.Errorf("Put failed: %v", err)
		return
	}

	retrieved, err := db.Get(key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
		return
	}

	if string(retrieved) != string(value) {
		t.Errorf("Value mismatch: got %s, want %s", retrieved, value)
		return
	}
}

func TestRocksDBConcurrent(t *testing.T) {
	db, err := New(fmt.Sprintf("/tmp/test-rocksdb-concurrent-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(db.Path())
	}()

	// Test concurrent access
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				key := []byte(fmt.Sprintf("key-%d-%d", id, j))
				value := []byte(fmt.Sprintf("value-%d-%d", id, j))

				if err := db.Put(key, value); err != nil {
					t.Errorf("Put failed: %v", err)
					return
				}

				retrieved, err := db.Get(key)
				if err != nil {
					t.Errorf("Get failed: %v", err)
					return
				}

				if string(retrieved) != string(value) {
					t.Errorf("Value mismatch: got %s, want %s", retrieved, value)
					return
				}
			}
		}(i)
	}
	wg.Wait()
}

func TestRocksDBIterator(t *testing.T) {
	db, err := New(fmt.Sprintf("/tmp/test-rocksdb-iter-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(db.Path())
	}()

	// Insert test data
	prefix := []byte("test-")
	for i := 0; i < 10; i++ {
		key := append(prefix, []byte(strconv.Itoa(i))...)
		value := []byte(fmt.Sprintf("value-%d", i))
		if err := db.Put(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Test iterator
	iter := db.NewIterator(prefix, nil)
	defer iter.Release()

	count := 0
	for iter.Next() {
		count++
	}

	if err := iter.Error(); err != nil {
		t.Fatal(err)
	}

	if count != 10 {
		t.Fatalf("Expected 10 items, got %d", count)
	}
}

func TestRocksDBBatch(t *testing.T) {
	db, err := New(fmt.Sprintf("/tmp/test-rocksdb-batch-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(db.Path())
	}()

	// Test batch operations
	batch := db.NewBatch()
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("batch-key-%d", i))
		value := []byte(fmt.Sprintf("batch-value-%d", i))
		if err := batch.Put(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Verify data is not yet in database
	if has, _ := db.Has([]byte("batch-key-0")); has {
		t.Fatal("Data should not be in database before batch write")
	}

	// Write batch
	if err := batch.Write(); err != nil {
		t.Fatal(err)
	}

	// Verify data is now in database
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("batch-key-%d", i))
		expected := []byte(fmt.Sprintf("batch-value-%d", i))

		value, err := db.Get(key)
		if err != nil {
			t.Fatal(err)
		}

		if string(value) != string(expected) {
			t.Fatalf("Value mismatch for key %s: got %s, want %s", key, value, expected)
		}
	}
}

func TestRocksDBDeleteRange(t *testing.T) {
	db, err := New(fmt.Sprintf("/tmp/test-rocksdb-delrange-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(db.Path())
	}()

	// Insert test data
	for i := 0; i < 100; i++ {
		key := []byte(fmt.Sprintf("key-%03d", i))
		value := []byte(fmt.Sprintf("value-%d", i))
		if err := db.Put(key, value); err != nil {
			t.Fatal(err)
		}
	}

	// Delete range [key-010, key-020)
	start := []byte("key-010")
	end := []byte("key-020")
	if err := db.DeleteRange(start, end); err != nil {
		t.Fatal(err)
	}

	// Verify keys in range are deleted
	for i := 10; i < 20; i++ {
		key := []byte(fmt.Sprintf("key-%03d", i))
		if has, _ := db.Has(key); has {
			t.Fatalf("Key %s should have been deleted", key)
		}
	}

	// Verify keys outside range still exist
	key := []byte("key-009")
	if has, _ := db.Has(key); !has {
		t.Fatalf("Key %s should still exist", key)
	}

	key = []byte("key-020")
	if has, _ := db.Has(key); !has {
		t.Fatalf("Key %s should still exist", key)
	}
}

// mockWriter implements ethdb.KeyValueWriter for testing batch replay
type mockWriter struct {
	data map[string]string
}

func (m *mockWriter) Put(key, value []byte) error {
	m.data[string(key)] = string(value)
	return nil
}

func (m *mockWriter) Delete(key []byte) error {
	delete(m.data, string(key))
	return nil
}

func TestRocksDBBatchReplay(t *testing.T) {
	db, err := New(fmt.Sprintf("/tmp/test-rocksdb-replay-%d-%d", os.Getpid(), time.Now().UnixNano()), 16, 16, "", false)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		db.Close()
		os.RemoveAll(db.Path())
	}()

	// Create a batch with some operations
	batch := db.NewBatch()
	testData := map[string]string{
		"key1": "value1",
		"key2": "value2",
		"key3": "value3",
	}

	for k, v := range testData {
		if err := batch.Put([]byte(k), []byte(v)); err != nil {
			t.Fatal(err)
		}
	}

	// Create a mock writer to capture replayed operations
	writer := &mockWriter{data: make(map[string]string)}

	// Test replay functionality
	if err := batch.Replay(writer); err != nil {
		t.Fatal(err)
	}

	// Verify replayed data matches original
	for k, v := range testData {
		if writer.data[k] != v {
			t.Fatalf("Replay mismatch for key %s: got %s, want %s", k, writer.data[k], v)
		}
	}
}
