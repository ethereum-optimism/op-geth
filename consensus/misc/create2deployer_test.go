package misc

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

func TestEnsureCreate2Deployer(t *testing.T) {
	canyonTime := uint64(1000)
	tests := []struct {
		name       string
		override   func(cfg *params.ChainConfig)
		timestamp  uint64
		codeExists bool
		applied    bool
	}{
		{
			name:      "at hardfork",
			timestamp: canyonTime,
			applied:   true,
		},
		{
			name: "another chain ID",
			override: func(cfg *params.ChainConfig) {
				cfg.ChainID = big.NewInt(params.OPMainnetChainID)
			},
			timestamp: canyonTime,
			applied:   true,
		},
		{
			name:       "code already exists",
			timestamp:  canyonTime,
			codeExists: true,
			applied:    true,
		},
		{
			name:      "pre canyon",
			timestamp: canyonTime - 1,
			applied:   false,
		},
		{
			name:      "post hardfork",
			timestamp: canyonTime + 1,
			applied:   false,
		},
		{
			name: "canyon not configured",
			override: func(cfg *params.ChainConfig) {
				cfg.CanyonTime = nil
			},
			timestamp: canyonTime,
			applied:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := params.ChainConfig{
				ChainID:    big.NewInt(params.BaseMainnetChainID),
				Optimism:   &params.OptimismConfig{},
				CanyonTime: &canyonTime,
			}
			if tt.override != nil {
				tt.override(&cfg)
			}
			state := &stateDb{
				codeExists: tt.codeExists,
			}
			EnsureCreate2Deployer(&cfg, tt.timestamp, state)
			assert.Equal(t, tt.applied, state.codeSet)
		})
	}
}

type stateDb struct {
	vm.StateDB
	codeExists bool
	codeSet    bool
}

var _ vm.StateDB = (*stateDb)(nil)

func (s *stateDb) GetCodeSize(_ common.Address) int {
	if s.codeExists {
		return 1
	}
	return 0
}

func (s *stateDb) SetCode(_ common.Address, _ []byte, _ tracing.CodeChangeReason) []byte {
	s.codeSet = true
	return nil
}
