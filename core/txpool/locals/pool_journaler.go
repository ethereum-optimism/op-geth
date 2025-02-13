// Copyright 2025 The op-geth Authors
// This file is part of the op-geth library.
//
// The op-geth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The op-geth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the op-geth library. If not, see <http://www.gnu.org/licenses/>.

package locals

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

// A PoolJournaler periodically journales a transaction pool to disk.
// It is meant to be used instead of the TxTracker when the noLocals and
// journalRemotes settings are enabled.
//
// OP-Stack addition.
type PoolJournaler struct {
	journal   *journal       // Journal of local transaction to back up to disk
	rejournal time.Duration  // How often to rotate journal
	pool      *txpool.TxPool // The tx pool to interact with

	shutdownCh chan struct{}
	wg         sync.WaitGroup
}

func NewPoolJournaler(journalPath string, journalTime time.Duration, pool *txpool.TxPool) *PoolJournaler {
	return &PoolJournaler{
		journal:    newTxJournal(journalPath),
		rejournal:  journalTime,
		pool:       pool,
		shutdownCh: make(chan struct{}),
	}
}

// Start implements node.Lifecycle interface
// Start is called after all services have been constructed and the networking
// layer was also initialized to spawn any goroutines required by the service.
func (pj *PoolJournaler) Start() error {
	pj.wg.Add(1)
	go pj.loop()
	return nil
}

// Stop implements node.Lifecycle interface
// Stop terminates all goroutines belonging to the service, blocking until they
// are all terminated.
func (pj *PoolJournaler) Stop() error {
	close(pj.shutdownCh)
	pj.wg.Wait()
	return nil
}

func (pj *PoolJournaler) loop() {
	defer log.Info("PoolJournaler: Stopped")
	defer pj.wg.Done()

	start := time.Now()
	log.Info("PoolJournaler: Start loading transactions from journal...")
	if err := pj.journal.load(func(transactions []*types.Transaction) []error {
		log.Info("PoolJournaler: Start adding transactions to pool", "count", len(transactions), "loading_duration", time.Since(start))
		errs := pj.pool.Add(transactions, true)
		log.Info("PoolJournaler: Done adding transactions to pool", "total_duration", time.Since(start))
		return errs
	}); err != nil {
		log.Error("PoolJournaler: Transaction journal loading failed. Exiting.", "err", err)
		return
	}
	defer pj.journal.close()

	ticker := time.NewTicker(pj.rejournal)
	defer ticker.Stop()

	journal := func() {
		start := time.Now()
		tojournal := pj.pool.ToJournal()
		if err := pj.journal.rotate(tojournal); err != nil {
			log.Error("PoolJournaler: Transaction journal rotation failed", "err", err)
		} else {
			log.Debug("PoolJournaler: Transaction journal rotated", "count", len(tojournal), "duration", time.Since(start))
		}
	}

	for {
		select {
		case <-pj.shutdownCh:
			journal()
			return
		case <-ticker.C:
			journal()
		}
	}
}
