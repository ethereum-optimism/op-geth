package core

import (
	"bytes"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
	"golang.org/x/exp/slog"
)

func TestRewindOnConfigChange(t *testing.T) {
	genesisTime := uint64(12)

	type testCase struct {
		name              string
		override1         func(*params.ChainConfig)
		override2         func(*params.ChainConfig)
		expectChainRewind bool
	}

	tcs := []testCase{
		{
			name:              fmt.Sprintf("CanyonTime changes from 10 to 0 (genesis time is %d)", genesisTime),
			override1:         func(c *params.ChainConfig) { c.CanyonTime = uint64ptr(10) },
			override2:         func(c *params.ChainConfig) { c.CanyonTime = uint64ptr(0) },
			expectChainRewind: false,
		},
		{
			name:              fmt.Sprintf("RegolithTime changes from 10 to 0 (genesis time is %d)", genesisTime),
			override1:         func(c *params.ChainConfig) { c.RegolithTime = uint64ptr(10) },
			override2:         func(c *params.ChainConfig) { c.RegolithTime = uint64ptr(0) },
			expectChainRewind: false,
		},
		{
			name:              fmt.Sprintf("ShanghaiTime changes from 10 to 0 (genesis time is %d)", genesisTime),
			override1:         func(c *params.ChainConfig) { c.ShanghaiTime = uint64ptr(10) },
			override2:         func(c *params.ChainConfig) { c.ShanghaiTime = uint64ptr(0) },
			expectChainRewind: false,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			var genesis = &Genesis{
				BaseFee:   big.NewInt(params.InitialBaseFee),
				Config:    params.AllEthashProtocolChanges,
				Timestamp: genesisTime,
			}

			// Prepare initial config for chain
			tc.override1(genesis.Config)
			db := rawdb.NewMemoryDatabase()
			cc := DefaultCacheConfigWithScheme(rawdb.PathScheme)

			// Start blockchain once to store config in DB
			blockchain, _ := NewBlockChain(db, cc, genesis, nil, ethash.NewFaker(), vm.Config{}, nil, nil)

			// Stop chain after 1 second
			<-time.After(1 * time.Second)
			blockchain.Stop()

			// Setup a buffer to capture logs
			logBuffer := bytes.Buffer{}
			log.SetDefault(log.NewLogger(slog.NewTextHandler(&logBuffer, nil)))

			// Restart chain with modified genesis config
			tc.override2(genesis.Config)
			blockchain, _ = NewBlockChain(db, cc, genesis, nil, ethash.NewFaker(), vm.Config{}, nil, nil)
			<-time.After(1 * time.Second)
			blockchain.Stop()

			// Inspect logs and assert on contents
			rewindTriggered := strings.Contains(logBuffer.String(), "Rewinding chain to upgrade configuration")
			if tc.expectChainRewind {
				require.True(t, rewindTriggered, "Required log line indicating chain rewind, but did not find one")
			} else {
				require.False(t, rewindTriggered, "Required NO log line indicating chain rewind, but found one")
			}
		})
	}
}
