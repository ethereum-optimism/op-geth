// Copyright 2025 The go-ethereum Authors
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

package beacon

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// mockChainHeaderReader is a minimal consensus.ChainHeaderReader that only
// serves a chain config, which is all the Karst kill switch needs.
type mockChainHeaderReader struct {
	config *params.ChainConfig
}

func (m *mockChainHeaderReader) Config() *params.ChainConfig                 { return m.config }
func (m *mockChainHeaderReader) CurrentHeader() *types.Header                { return nil }
func (m *mockChainHeaderReader) GetHeader(common.Hash, uint64) *types.Header { return nil }
func (m *mockChainHeaderReader) GetHeaderByNumber(uint64) *types.Header      { return nil }
func (m *mockChainHeaderReader) GetHeaderByHash(common.Hash) *types.Header   { return nil }

// karstConfig returns an Optimism test config with Karst activating at the given
// timestamp.
func karstConfig(karstTime uint64) *params.ChainConfig {
	conf := *params.OptimismTestConfig
	conf.KarstTime = &karstTime
	return &conf
}

func TestKarstKillSwitch(t *testing.T) {
	const karstTime = uint64(20000)

	nonOptimism := *params.OptimismTestConfig
	nonOptimism.Optimism = nil
	nonOptimism.KarstTime = &[]uint64{karstTime}[0]

	tests := []struct {
		name    string
		config  *params.ChainConfig
		time    uint64
		blocked bool
	}{
		{"karst-unset", params.OptimismTestConfig, karstTime, false},
		{"before-karst", karstConfig(karstTime), karstTime - 1, false},
		{"activation-block", karstConfig(karstTime), karstTime, true},
		{"after-karst", karstConfig(karstTime), karstTime + 1, true},
		{"non-optimism-chain", &nonOptimism, karstTime + 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := karstKillSwitch(tt.config, tt.time)
			if got := errors.Is(err, errKarstDeprecated); got != tt.blocked {
				t.Fatalf("karstKillSwitch blocked = %v, want %v (err: %v)", got, tt.blocked, err)
			}
		})
	}
}

// TestBeaconKarstKillSwitchEntryPoints verifies that every consensus-engine
// entry point that touches a block refuses a block at or after Karst, so an
// op-geth node can neither import nor build a chain that crosses Karst.
func TestBeaconKarstKillSwitchEntryPoints(t *testing.T) {
	const karstTime = uint64(20000)
	engine := New(nil)

	// At and after the Karst activation timestamp every entry point must refuse.
	for _, time := range []uint64{karstTime, karstTime + 1} {
		chain := &mockChainHeaderReader{config: karstConfig(karstTime)}
		header := &types.Header{Number: big.NewInt(1), Time: time}

		if err := engine.verifyHeader(chain, header, &types.Header{Number: big.NewInt(0)}); !errors.Is(err, errKarstDeprecated) {
			t.Errorf("verifyHeader(time=%d) = %v, want errKarstDeprecated", time, err)
		}
		if err := engine.Prepare(chain, header); !errors.Is(err, errKarstDeprecated) {
			t.Errorf("Prepare(time=%d) = %v, want errKarstDeprecated", time, err)
		}
		if _, err := engine.FinalizeAndAssemble(context.Background(), chain, header, nil, &types.Body{}, nil); !errors.Is(err, errKarstDeprecated) {
			t.Errorf("FinalizeAndAssemble(time=%d) = %v, want errKarstDeprecated", time, err)
		}
	}

	// With Karst unset, the kill switch must not fire. Prepare is a real entry
	// point we can exercise without a full valid header/parent.
	chain := &mockChainHeaderReader{config: params.OptimismTestConfig}
	header := &types.Header{Number: big.NewInt(1), Time: karstTime + 1}
	if err := engine.Prepare(chain, header); errors.Is(err, errKarstDeprecated) {
		t.Errorf("Prepare with Karst unset returned errKarstDeprecated: %v", err)
	}
}
