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

// Package rocksdb implements the key-value database layer based on RocksDB.
package rocksdb

import (
	"bytes"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/linxGnu/grocksdb"
)

const (
	// minCache is the minimum amount of memory in megabytes to allocate to rocksdb
	// read and write caching, split half and half.
	minCache = 16

	// minHandles is the minimum number of files handles to allocate to the open
	// database files.
	minHandles = 16

	// metricsGatheringInterval specifies the interval to retrieve rocksdb database
	// compaction, io and pause stats to report to the user.
	metricsGatheringInterval = 3 * time.Second

	// degradationWarnInterval specifies how often warning should be printed if the
	// rocksdb database cannot keep up with requested writes.
	degradationWarnInterval = time.Minute
)

// Database is a persistent key-value store based on the RocksDB storage engine.
// Apart from basic data storage functionality it also supports batch writes and
// iterating over the keyspace in binary-alphabetical order.
type Database struct {
	fn        string       // filename for reporting
	db        *grocksdb.DB // Underlying RocksDB storage engine
	namespace string       // Namespace for metrics

	compTimeMeter          *metrics.Meter   // Meter for measuring the total time spent in database compaction
	compReadMeter          *metrics.Meter   // Meter for measuring the data read during compaction
	compWriteMeter         *metrics.Meter   // Meter for measuring the data written during compaction
	writeDelayNMeter       *metrics.Meter   // Meter for measuring the write delay number due to database compaction
	writeDelayMeter        *metrics.Meter   // Meter for measuring the write delay duration due to database compaction
	diskSizeGauge          *metrics.Gauge   // Gauge for tracking the size of all the levels in the database
	diskReadMeter          *metrics.Meter   // Meter for measuring the effective amount of data read
	diskWriteMeter         *metrics.Meter   // Meter for measuring the effective amount of data written
	memCompGauge           *metrics.Gauge   // Gauge for tracking the number of memory compaction
	level0CompGauge        *metrics.Gauge   // Gauge for tracking the number of table compaction in level0
	nonlevel0CompGauge     *metrics.Gauge   // Gauge for tracking the number of table compaction in non0 level
	seekCompGauge          *metrics.Gauge   // Gauge for tracking the number of table compaction caused by read opt
	manualMemAllocGauge    *metrics.Gauge   // Gauge for tracking amount of non-managed memory currently allocated
	liveMemTablesGauge     *metrics.Gauge   // Gauge for tracking the number of live memory tables
	zombieMemTablesGauge   *metrics.Gauge   // Gauge for tracking the number of zombie memory tables
	blockCacheHitGauge     *metrics.Gauge   // Gauge for tracking the number of total hit in the block cache
	blockCacheMissGauge    *metrics.Gauge   // Gauge for tracking the number of total miss in the block cache
	tableCacheHitGauge     *metrics.Gauge   // Gauge for tracking the number of total hit in the table cache
	tableCacheMissGauge    *metrics.Gauge   // Gauge for tracking the number of total miss in the table cache
	filterHitGauge         *metrics.Gauge   // Gauge for tracking the number of total hit in bloom filter
	filterMissGauge        *metrics.Gauge   // Gauge for tracking the number of total miss in bloom filter
	estimatedCompDebtGauge *metrics.Gauge   // Gauge for tracking the number of bytes that need to be compacted
	liveCompGauge          *metrics.Gauge   // Gauge for tracking the number of in-progress compactions
	liveCompSizeGauge      *metrics.Gauge   // Gauge for tracking the size of in-progress compactions
	liveIterGauge          *metrics.Gauge   // Gauge for tracking the number of live database iterators
	levelsGauge            []*metrics.Gauge // Gauge for tracking the number of tables in levels

	quitLock sync.RWMutex    // Mutex protecting the quit channel and the closed flag
	quitChan chan chan error // Quit channel to stop the metrics collection before closing the database
	closed   bool            // keep track of whether we're Closed

	log log.Logger // Contextual logger tracking the database path

	activeComp    int           // Current number of active compactions
	compStartTime time.Time     // The start time of the earliest currently-active compaction
	compTime      atomic.Int64  // Total time spent in compaction in ns
	level0Comp    atomic.Uint32 // Total number of level-zero compactions
	nonLevel0Comp atomic.Uint32 // Total number of non level-zero compactions

	writeStalled        atomic.Bool  // Flag whether the write is stalled
	writeDelayStartTime time.Time    // The start time of the latest write stall
	writeDelayReason    string       // The reason of the latest write stall
	writeDelayCount     atomic.Int64 // Total number of write stall counts
	writeDelayTime      atomic.Int64 // Total time spent in write stalls

	writeOptions *grocksdb.WriteOptions
}

// New returns a wrapped RocksDB object. The namespace is the prefix that the
// metrics reporting should use for surfacing internal stats.
func New(file string, cache int, handles int, namespace string, readonly bool) (*Database, error) {
	// Ensure we have some minimal caching and file guarantees
	if cache < minCache {
		cache = minCache
	}
	if handles < minHandles {
		handles = minHandles
	}
	logger := log.New("database", file)
	logger.Info("Allocated cache and file handles", "cache", common.StorageSize(cache*1024*1024), "handles", handles)

	return NewCustom(file, namespace, func(opts *grocksdb.Options) {
		// Set default options
		opts.SetMaxOpenFiles(handles)

		// Set cache sizes (split between block cache and write buffer)
		blockCache := grocksdb.NewLRUCache(uint64(cache / 2 * 1024 * 1024))
		bbto := grocksdb.NewDefaultBlockBasedTableOptions()
		bbto.SetBlockCache(blockCache)
		opts.SetBlockBasedTableFactory(bbto)

		// Write buffer size - RocksDB uses this for memtable
		opts.SetWriteBufferSize(uint64(cache / 4 * 1024 * 1024))
		opts.SetMaxWriteBufferNumber(3) // Similar to LevelDB's behavior

		// Performance tuning
		opts.SetMaxBackgroundCompactions(4)
		opts.SetMaxBackgroundFlushes(2)

		// Bloom filter for better read performance
		bbto.SetFilterPolicy(grocksdb.NewBloomFilter(10))

		if readonly {
			// For readonly mode, we don't need write buffers
			opts.SetWriteBufferSize(0)
		}
	})
}

// NewCustom returns a wrapped RocksDB object. The namespace is the prefix that the
// metrics reporting should use for surfacing internal stats.
// The customize function allows the caller to modify the rocksdb options.
func NewCustom(file string, namespace string, customize func(options *grocksdb.Options)) (*Database, error) {
	options := configureOptions(customize)
	logger := log.New("database", file)

	// Log configuration details
	usedCache := options.GetWriteBufferSize() * uint64(options.GetMaxWriteBufferNumber())
	logCtx := []interface{}{"cache", common.StorageSize(usedCache), "handles", options.GetMaxOpenFiles()}
	logger.Info("Allocated cache and file handles", logCtx...)

	// Open the db
	db, err := grocksdb.OpenDb(options, file)
	if err != nil {
		return nil, err
	}

	// Create the database wrapper
	rdb := &Database{
		fn:        file,
		db:        db,
		namespace: namespace,
		log:       logger,
		quitChan:  make(chan chan error),

		// Use default write options
		writeOptions: grocksdb.NewDefaultWriteOptions(),
	}

	// Initialize metrics using GetOrRegisterMeter/Gauge pattern like Pebble
	rdb.compTimeMeter = metrics.GetOrRegisterMeter(namespace+"compact/time", nil)
	rdb.compReadMeter = metrics.GetOrRegisterMeter(namespace+"compact/input", nil)
	rdb.compWriteMeter = metrics.GetOrRegisterMeter(namespace+"compact/output", nil)
	rdb.diskSizeGauge = metrics.GetOrRegisterGauge(namespace+"disk/size", nil)
	rdb.diskReadMeter = metrics.GetOrRegisterMeter(namespace+"disk/read", nil)
	rdb.diskWriteMeter = metrics.GetOrRegisterMeter(namespace+"disk/write", nil)
	rdb.writeDelayMeter = metrics.GetOrRegisterMeter(namespace+"compact/writedelay/duration", nil)
	rdb.writeDelayNMeter = metrics.GetOrRegisterMeter(namespace+"compact/writedelay/counter", nil)
	rdb.memCompGauge = metrics.GetOrRegisterGauge(namespace+"compact/memory", nil)
	rdb.level0CompGauge = metrics.GetOrRegisterGauge(namespace+"compact/level0", nil)
	rdb.nonlevel0CompGauge = metrics.GetOrRegisterGauge(namespace+"compact/nonlevel0", nil)
	rdb.seekCompGauge = metrics.GetOrRegisterGauge(namespace+"compact/seek", nil)
	rdb.manualMemAllocGauge = metrics.GetOrRegisterGauge(namespace+"memory/manualalloc", nil)
	rdb.liveMemTablesGauge = metrics.GetOrRegisterGauge(namespace+"table/live", nil)
	rdb.zombieMemTablesGauge = metrics.GetOrRegisterGauge(namespace+"table/zombie", nil)
	rdb.blockCacheHitGauge = metrics.GetOrRegisterGauge(namespace+"cache/block/hit", nil)
	rdb.blockCacheMissGauge = metrics.GetOrRegisterGauge(namespace+"cache/block/miss", nil)
	rdb.tableCacheHitGauge = metrics.GetOrRegisterGauge(namespace+"cache/table/hit", nil)
	rdb.tableCacheMissGauge = metrics.GetOrRegisterGauge(namespace+"cache/table/miss", nil)
	rdb.filterHitGauge = metrics.GetOrRegisterGauge(namespace+"filter/hit", nil)
	rdb.filterMissGauge = metrics.GetOrRegisterGauge(namespace+"filter/miss", nil)
	rdb.estimatedCompDebtGauge = metrics.GetOrRegisterGauge(namespace+"compact/estimateDebt", nil)
	rdb.liveCompGauge = metrics.GetOrRegisterGauge(namespace+"compact/live/count", nil)
	rdb.liveCompSizeGauge = metrics.GetOrRegisterGauge(namespace+"compact/live/size", nil)
	rdb.liveIterGauge = metrics.GetOrRegisterGauge(namespace+"iter/count", nil)

	// Start up the metrics gathering and return
	go rdb.meter(metricsGatheringInterval, namespace)
	return rdb, nil
}

// configureOptions sets some default options, then runs the provided setter.
func configureOptions(customizeFn func(*grocksdb.Options)) *grocksdb.Options {
	// Set default options
	options := grocksdb.NewDefaultOptions()

	// Create directories if needed
	options.SetCreateIfMissing(true)
	options.SetCreateIfMissingColumnFamilies(true)

	// Allow caller to make custom modifications to the options
	if customizeFn != nil {
		customizeFn(options)
	}
	return options
}

// Close stops the metrics collection, flushes any pending data to disk and closes
// all io accesses to the underlying key-value store.
func (d *Database) Close() error {
	d.quitLock.Lock()
	defer d.quitLock.Unlock()
	// Allow double closing, simplifies things
	if d.closed {
		return nil
	}
	d.closed = true
	if d.quitChan != nil {
		errc := make(chan error)
		d.quitChan <- errc
		if err := <-errc; err != nil {
			d.log.Error("Metrics collection failed", "err", err)
		}
		d.quitChan = nil
	}

	// Close options and database
	if d.writeOptions != nil {
		d.writeOptions.Destroy()
		d.writeOptions = nil
	}
	if d.db != nil {
		d.db.Close()
		d.db = nil
	}
	return nil
}

// Has retrieves if a key is present in the key-value store.
func (d *Database) Has(key []byte) (bool, error) {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return false, errors.New("database closed")
	}

	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	data, err := d.db.Get(ro, key)
	if err != nil {
		return false, err
	}
	defer data.Free()
	return data.Data() != nil, nil
}

// Get retrieves the given key if it's present in the key-value store.
func (d *Database) Get(key []byte) ([]byte, error) {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return nil, errors.New("database closed")
	}

	ro := grocksdb.NewDefaultReadOptions()
	defer ro.Destroy()

	data, err := d.db.Get(ro, key)
	if err != nil {
		return nil, err
	}
	defer data.Free()

	if data.Data() == nil {
		return nil, errors.New("key not found")
	}

	// Copy the data since we're freeing the slice
	result := make([]byte, len(data.Data()))
	copy(result, data.Data())
	return result, nil
}

// Put inserts the given value into the key-value store.
func (d *Database) Put(key []byte, value []byte) error {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return errors.New("database closed")
	}
	return d.db.Put(d.writeOptions, key, value)
}

// Delete removes the key from the key-value store.
func (d *Database) Delete(key []byte) error {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()
	if d.closed {
		return errors.New("database closed")
	}
	return d.db.Delete(d.writeOptions, key)
}

// DeleteRange deletes all of the keys (and values) in the range [start,end)
// (inclusive on start, exclusive on end).
func (d *Database) DeleteRange(start, end []byte) error {
	d.quitLock.RLock()
	defer d.quitLock.RUnlock()

	if d.closed {
		return errors.New("database closed")
	}

	if start == nil {
		it := d.db.NewIterator(grocksdb.NewDefaultReadOptions())
		defer it.Close()
		it.SeekToFirst()
		start = it.Key().Data()
	}

	// There is no special flag to represent the end of key range.
	// Last iterator key cannot be used since it should also be deleted.
	// Use an ugly hack to construct a large key to represent it.
	if end == nil {
		end = ethdb.MaximumKey
	}

	if bytes.Compare(start, end) >= 0 {
		return nil
	}

	wopts := grocksdb.NewDefaultWriteOptions()
	defer wopts.Destroy()
	cfh := d.db.GetDefaultColumnFamily()
	return d.db.DeleteRangeCF(wopts, cfh, start, end)
}

// NewBatch creates a write-only key-value store that buffers changes to its host
// database until a final write is called.
func (d *Database) NewBatch() ethdb.Batch {
	return &batch{
		b:  grocksdb.NewWriteBatch(),
		db: d,
	}
}

// NewBatchWithSize creates a write-only database batch with pre-allocated buffer.
func (d *Database) NewBatchWithSize(size int) ethdb.Batch {
	// RocksDB WriteBatch doesn't have size pre-allocation in grocksdb
	wb := grocksdb.NewWriteBatch()
	return &batch{
		b:  wb,
		db: d,
	}
}

// upperBound returns the upper bound for the given prefix
func upperBound(prefix []byte) (limit []byte) {
	for i := len(prefix) - 1; i >= 0; i-- {
		c := prefix[i]
		if c == 0xff {
			continue
		}
		limit = make([]byte, i+1)
		copy(limit, prefix)
		limit[i] = c + 1
		break
	}
	return limit
}

// NewIterator creates a binary-alphabetical iterator over a subset
// of database content with a particular key prefix, starting at a particular
// initial key (or after, if it does not exist).
func (d *Database) NewIterator(prefix []byte, start []byte) ethdb.Iterator {
	ro := grocksdb.NewDefaultReadOptions()

	// Set up iteration bounds
	startKey := append(prefix, start...)
	endKey := upperBound(prefix)

	if endKey != nil {
		ro.SetIterateUpperBound(endKey)
	}

	iter := d.db.NewIterator(ro)
	iter.Seek(startKey)

	return &rocksdbIterator{
		iter:     iter,
		ro:       ro,
		prefix:   prefix,
		moved:    true,
		released: false,
	}
}

// Stat returns the internal metrics of RocksDB in a text format. It's a developer
// method to read everything there is to read, independent of RocksDB version.
func (d *Database) Stat() (string, error) {
	stats := d.db.GetProperty("rocksdb.stats")
	if stats == "" {
		return "No stats available", nil
	}
	return stats, nil
}

// Compact flattens the underlying data store for the given key range. In essence,
// deleted and overwritten versions are discarded, and the data is rearranged to
// reduce the cost of operations needed to access them.
//
// A nil start is treated as a key before all keys in the data store; a nil limit
// is treated as a key after all keys in the data store. If both is nil then it
// will compact entire data store.
func (d *Database) Compact(start []byte, limit []byte) error {
	// There is no special flag to represent the end of key range
	// in rocksdb(nil in leveldb). Use an ugly hack to construct a
	// large key to represent it.
	// Note any prefixed database entry will be smaller than this
	// flag, as for trie nodes we need the 32 byte 0xff because
	// there might be a shared prefix starting with a number of
	// 0xff-s, so 32 ensures than only a hash collision could touch it.
	if limit == nil {
		limit = ethdb.MaximumKey
	}
	// RocksDB CompactRange
	d.db.CompactRange(grocksdb.Range{Start: start, Limit: limit})
	return nil
}

// Path returns the path to the database directory.
func (d *Database) Path() string {
	return d.fn
}

// SyncKeyValue flushes all pending writes in the write-ahead-log to disk,
// ensuring data durability up to that point.
func (d *Database) SyncKeyValue() error {
	// Create a sync write options to force flush
	wo := grocksdb.NewDefaultWriteOptions()
	wo.SetSync(true)
	defer wo.Destroy()

	// Create an empty batch and write it with sync to force WAL flush
	b := grocksdb.NewWriteBatch()
	defer b.Destroy()

	return d.db.Write(wo, b)
}

// meter periodically retrieves internal rocksdb counters and reports them to
// the metrics subsystem.
func (d *Database) meter(refresh time.Duration, namespace string) {
	var errc chan error
	timer := time.NewTimer(refresh)
	defer timer.Stop()

	// Create storage and warning log tracer for write delay.
	var (
		compTimes  [2]int64
		compWrites [2]int64
		compReads  [2]int64

		nWrites [2]int64

		writeDelayTimes      [2]int64
		writeDelayCounts     [2]int64
		lastWriteStallReport time.Time
	)

	// Iterate ad infinitum and collect the stats
	for i := 1; errc == nil; i++ {
		var (
			compWrite int64
			compRead  int64
			nWrite    int64

			compTime           = d.compTime.Load()
			writeDelayCount    = d.writeDelayCount.Load()
			writeDelayTime     = d.writeDelayTime.Load()
			nonLevel0CompCount = int64(d.nonLevel0Comp.Load())
			level0CompCount    = int64(d.level0Comp.Load())
		)
		writeDelayTimes[i%2] = writeDelayTime
		writeDelayCounts[i%2] = writeDelayCount
		compTimes[i%2] = compTime

		// Get various RocksDB statistics
		if stats := d.db.GetProperty("rocksdb.stats"); stats != "" {
			// Update basic metrics that we can get from RocksDB
			if totalSize := d.db.GetProperty("rocksdb.total-sst-files-size"); totalSize != "" {
				// Parse and update size metrics - simplified for now
				d.diskSizeGauge.Update(0)
			}

			if level0Files := d.db.GetProperty("rocksdb.num-files-at-level0"); level0Files != "" {
				// Parse and update level metrics - simplified for now
				d.level0CompGauge.Update(level0CompCount)
			}

			if memTableSize := d.db.GetProperty("rocksdb.cur-size-all-mem-tables"); memTableSize != "" {
				// Parse and update memory metrics - simplified for now
				d.manualMemAllocGauge.Update(0)
			}
		}

		compWrites[i%2] = compWrite
		compReads[i%2] = compRead
		nWrites[i%2] = nWrite

		d.writeDelayNMeter.Mark(writeDelayCounts[i%2] - writeDelayCounts[(i-1)%2])
		d.writeDelayMeter.Mark(writeDelayTimes[i%2] - writeDelayTimes[(i-1)%2])
		// Print a warning log if writing has been stalled for a while. The log will
		// be printed per minute to avoid overwhelming users.
		if d.writeStalled.Load() && writeDelayCounts[i%2] == writeDelayCounts[(i-1)%2] &&
			time.Now().After(lastWriteStallReport.Add(degradationWarnInterval)) {
			d.log.Warn("Database compacting, degraded performance")
			lastWriteStallReport = time.Now()
		}
		d.compTimeMeter.Mark(compTimes[i%2] - compTimes[(i-1)%2])
		d.compReadMeter.Mark(compReads[i%2] - compReads[(i-1)%2])
		d.compWriteMeter.Mark(compWrites[i%2] - compWrites[(i-1)%2])
		d.diskReadMeter.Mark(0) // rocksdb doesn't track non-compaction reads
		d.diskWriteMeter.Mark(nWrites[i%2] - nWrites[(i-1)%2])

		d.nonlevel0CompGauge.Update(nonLevel0CompCount)
		d.level0CompGauge.Update(level0CompCount)

		// Sleep a bit, then repeat the stats collection
		select {
		case errc = <-d.quitChan:
			// Quit requesting, stop hammering the database
		case <-timer.C:
			timer.Reset(refresh)
			// Timeout, gather a new set of stats
		}
	}
	errc <- nil
}

// batch is a write-only batch that commits changes to its host database
// when Write is called. A batch cannot be used concurrently.
type batch struct {
	b    *grocksdb.WriteBatch
	db   *Database
	size int
}

// Put inserts the given value into the batch for later committing.
func (b *batch) Put(key, value []byte) error {
	b.b.Put(key, value)
	b.size += len(key) + len(value)
	return nil
}

// Delete inserts the key removal into the batch for later committing.
func (b *batch) Delete(key []byte) error {
	b.b.Delete(key)
	b.size += len(key)
	return nil
}

// DeleteRange removes all keys in the range [start, end) from the batch for
// later committing, inclusive on start, exclusive on end.
func (b *batch) DeleteRange(start, end []byte) error {
	// If start is nil, we assume it's the key before all keys
	if start == nil {
		start = []byte{0}
	}
	// If end is nil, we assume it's the key after all keys
	if end == nil {
		end = ethdb.MaximumKey
	}
	b.b.DeleteRange(start, end)
	b.size += len(start) + len(end)
	return nil
}

// ValueSize retrieves the amount of data queued up for writing.
func (b *batch) ValueSize() int {
	return b.size
}

// Write flushes any accumulated data to disk.
func (b *batch) Write() error {
	b.db.quitLock.RLock()
	defer b.db.quitLock.RUnlock()
	if b.db.closed {
		return errors.New("database closed")
	}
	err := b.db.db.Write(b.db.writeOptions, b.b)
	if err != nil {
		return err
	}
	b.b.Clear()
	b.size = 0
	return nil
}

// Reset resets the batch for reuse.
func (b *batch) Reset() {
	b.size = 0
	b.b.Clear()
}

// Replay replays the batch contents.
func (b *batch) Replay(w ethdb.KeyValueWriter) error {
	it := b.b.NewIterator()
	for {
		if !it.Next() {
			return it.Error()
		}
		rec := it.Record()
		if rec == nil {
			return it.Error()
		}
		k := rec.Key
		v := rec.Value
		kind := rec.Type
		switch kind {
		case grocksdb.WriteBatchValueRecord:
			if err := w.Put(k, v); err != nil {
				return err
			}
		case grocksdb.WriteBatchDeletionRecord:
			if err := w.Delete(k); err != nil {
				return err
			}
		case grocksdb.WriteBatchRangeDeletion:
			// For range deletion, k is the start key and v is the end key
			if rangeDeleter, ok := w.(ethdb.KeyValueRangeDeleter); ok {
				if err := rangeDeleter.DeleteRange(k, v); err != nil {
					return err
				}
			} else {
				return errors.New("ethdb.KeyValueWriter does not implement DeleteRange")
			}
		default:
			return fmt.Errorf("unhandled operation, keytype: %v", kind)
		}
	}
}

// rocksdbIterator is a wrapper of underlying iterator in storage engine.
// The purpose of this structure is to implement the missing APIs.
//
// The rocksdb iterator is not thread-safe.
type rocksdbIterator struct {
	iter     *grocksdb.Iterator
	ro       *grocksdb.ReadOptions
	prefix   []byte
	moved    bool
	released bool
}

// Next moves the iterator to the next key/value pair. It returns whether the
// iterator is exhausted.
func (it *rocksdbIterator) Next() bool {
	if it.moved {
		it.moved = false
		return it.Valid()
	}
	if !it.Valid() {
		return false
	}
	it.iter.Next()
	return it.Valid()
}

// Error returns any accumulated error. Exhausting all the key/value pairs
// is not considered to be an error.
func (it *rocksdbIterator) Error() error {
	// if iter is nil, it means the iterator is exhausted and released => no error
	if it.iter == nil {
		return nil
	}
	return it.iter.Err()
}

// Key returns the key of the current key/value pair, or nil if done. The caller
// should not modify the contents of the returned slice, and its contents may
// change on the next call to Next.
func (it *rocksdbIterator) Key() []byte {
	if !it.Valid() {
		return nil
	}
	key := it.iter.Key()
	defer key.Free()
	result := make([]byte, len(key.Data()))
	copy(result, key.Data())
	return result
}

// Value returns the value of the current key/value pair, or nil if done. The
// caller should not modify the contents of the returned slice, and its contents
// may change on the next call to Next.
func (it *rocksdbIterator) Value() []byte {
	if !it.Valid() {
		return nil
	}
	value := it.iter.Value()
	defer value.Free()
	result := make([]byte, len(value.Data()))
	copy(result, value.Data())
	return result
}

// Valid returns whether the iterator is positioned at a valid key/value pair.
func (it *rocksdbIterator) Valid() bool {
	if !it.iter.Valid() {
		return false
	}

	// Check if we're still within the prefix
	key := it.iter.Key()
	defer key.Free()
	return bytes.HasPrefix(key.Data(), it.prefix)
}

// Release releases associated resources. Release should always succeed and can
// be called multiple times without causing error.
func (it *rocksdbIterator) Release() {
	if !it.released {
		if it.iter != nil {
			it.iter.Close()
			it.iter = nil
		}
		if it.ro != nil {
			it.ro.Destroy()
			it.ro = nil
		}
		it.released = true
	}
}
