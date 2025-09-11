package catalyst

import (
	"errors"
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func preCanyon() *params.ChainConfig {
	cfg := new(params.ChainConfig)
	// Mark as an Optimism chain so IsOptimismFoo is true when FooTime is active
	cfg.Optimism = &params.OptimismConfig{}
	return cfg
}

func postCanyon() *params.ChainConfig {
	cfg := preCanyon()
	cfg.CanyonTime = new(uint64)
	return cfg
}

func postHolocene() *params.ChainConfig {
	cfg := postCanyon()
	cfg.HoloceneTime = new(uint64)
	return cfg
}

func postIsthmus() *params.ChainConfig {
	cfg := postHolocene()
	cfg.IsthmusTime = new(uint64)
	return cfg
}

func postJovian() *params.ChainConfig {
	cfg := postIsthmus()
	cfg.JovianTime = new(uint64)
	return cfg
}

var valid1559Params = []byte{0, 1, 2, 3, 4, 5, 6, 7}
var validExtraData = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}
var emptyWithdrawals = make([]*types.Withdrawal, 0)
var validJovianExtraData = append(append([]byte{1}, valid1559Params...), make([]byte, 8)...) // version=1, 8 bytes params, 8 byte minBaseFee

func TestCheckOptimismPayload(t *testing.T) {
	tests := []struct {
		name     string
		params   engine.ExecutableData
		cfg      *params.ChainConfig
		expected error
	}{
		{
			name: "valid payload pre-Canyon",
			params: engine.ExecutableData{
				ExtraData: []byte{},
			},
			cfg:      preCanyon(),
			expected: nil,
		},
		{
			name: "valid payload post-Canyon",
			params: engine.ExecutableData{
				ExtraData:   []byte{},
				Withdrawals: emptyWithdrawals,
			},
			cfg: postCanyon(),
		},
		{
			name: "invalid empty withdrawals post-Canyon",
			params: engine.ExecutableData{
				ExtraData:   []byte{},
				Withdrawals: make([]*types.Withdrawal, 1),
			},
			cfg:      postCanyon(),
			expected: errors.New("non-empty withdrawals post-Canyon"),
		},
		{
			name: "non-nil withdrawalsRoot pre-Isthmus",
			params: engine.ExecutableData{
				ExtraData:       []byte{},
				Withdrawals:     emptyWithdrawals,
				WithdrawalsRoot: &types.EmptyWithdrawalsHash,
			},
			cfg:      postCanyon(),
			expected: errors.New("non-nil withdrawalsRoot pre-Isthmus"),
		},
		{
			name: "invalid non-empty extraData pre-Holocene",
			params: engine.ExecutableData{
				ExtraData:   []byte{1, 2, 3},
				Withdrawals: emptyWithdrawals,
			},
			cfg:      postCanyon(),
			expected: errors.New("extraData must be empty before Holocene"),
		},
		{
			name: "invalid extraData post-Holocene",
			params: engine.ExecutableData{
				ExtraData:   []byte{1, 2, 3},
				Withdrawals: emptyWithdrawals,
			},
			cfg:      postHolocene(),
			expected: errors.New("holocene extraData should be 9 bytes, got 3"),
		},
		{
			name: "valid payload post-Holocene with extraData",
			params: engine.ExecutableData{
				ExtraData:   validExtraData,
				Withdrawals: emptyWithdrawals,
			},
			cfg:      postHolocene(),
			expected: nil,
		},
		{
			name: "invalid non-nil withdrawalsRoot post-Holocene",
			params: engine.ExecutableData{
				ExtraData:       validExtraData,
				Withdrawals:     emptyWithdrawals,
				WithdrawalsRoot: &types.EmptyWithdrawalsHash,
			},
			cfg:      postHolocene(),
			expected: errors.New("non-nil withdrawalsRoot pre-Isthmus"),
		},
		{
			name: "valid payload post-Isthmus",
			params: engine.ExecutableData{
				ExtraData:       validExtraData,
				Withdrawals:     emptyWithdrawals,
				WithdrawalsRoot: &types.EmptyWithdrawalsHash,
			},
			cfg:      postIsthmus(),
			expected: nil,
		},
		{
			name: "invalid nil withdrawals root post-isthmus",
			params: engine.ExecutableData{
				Withdrawals: emptyWithdrawals,
				ExtraData:   validExtraData,
			},
			cfg:      postIsthmus(),
			expected: errors.New("nil withdrawalsRoot post-Isthmus"),
		},
		{
			name: "valid payload post-Jovian with 17-byte extraData",
			params: engine.ExecutableData{
				Timestamp:       0,
				ExtraData:       validJovianExtraData,
				WithdrawalsRoot: &types.EmptyWithdrawalsHash,
			},
			cfg:      postJovian(),
			expected: nil,
		},
		{
			name: "invalid payload post-Jovian with 9-byte (Holocene) extraData",
			params: engine.ExecutableData{
				Timestamp: 0,
				ExtraData: validExtraData,
			},
			cfg:      postJovian(),
			expected: errors.New("MinBaseFee extraData should be 17 bytes, got 9"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := checkOptimismPayload(test.params, test.cfg)
			if test.expected == nil {
				require.NoError(t, err)
			} else {
				require.EqualError(t, err, test.expected.Error())
			}
		})
	}
}

func TestCheckOptimismPayloadAttributes(t *testing.T) {
	tests := []struct {
		name              string
		payloadAttributes *engine.PayloadAttributes
		cfg               *params.ChainConfig
		expected          error
		shouldPanic       bool
	}{
		{
			name:              "nil payload attributes panic",
			payloadAttributes: nil,
			cfg:               preCanyon(),
			shouldPanic:       true,
		},
		{
			name: "valid payload attributes pre-Canyon",
			payloadAttributes: &engine.PayloadAttributes{
				GasLimit: new(uint64),
			},
			cfg:      preCanyon(),
			expected: nil,
		},
		{
			name: "invalid nil gasLimit",
			payloadAttributes: &engine.PayloadAttributes{
				GasLimit: nil,
			},
			cfg:      preCanyon(),
			expected: errors.New("gasLimit parameter is required"),
		},
		{
			name: "invalid non-empty withdrawals post-Canyon",
			payloadAttributes: &engine.PayloadAttributes{
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
				Withdrawals:   make([]*types.Withdrawal, 1),
			},
			cfg:      postCanyon(),
			expected: errors.New("non-empty withdrawals post-Canyon"),
		},
		{
			name: "invalid non-empty eip1559Params pre-Holocene",
			payloadAttributes: &engine.PayloadAttributes{
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
			},
			cfg:      postCanyon(),
			expected: errors.New("non-empty eip155Params pre-Holocene"),
		},
		{
			name: "invalid eip1559Params post-Holocene",
			payloadAttributes: &engine.PayloadAttributes{
				GasLimit:      new(uint64),
				EIP1559Params: append(valid1559Params, 77),
			},
			cfg:      postHolocene(),
			expected: errors.New("holocene eip-1559 params should be 8 bytes, got 9"),
		},
		{
			name: "valid payload attributes post-Holocene",
			payloadAttributes: &engine.PayloadAttributes{
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
			},
			cfg:      postHolocene(),
			expected: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.shouldPanic {
				require.Panics(t, func() {
					checkOptimismPayloadAttributes(test.payloadAttributes, test.cfg)
				})
			} else {
				err := checkOptimismPayloadAttributes(test.payloadAttributes, test.cfg)
				if test.expected == nil {
					require.NoError(t, err)
				} else {
					require.EqualError(t, err, test.expected.Error())
				}
			}
		})
	}
}

// OP Stack test diff: verify that nil payload attributes does not cause a panic
func TestForkChoiceUpdatedNilPayloadAttributes(t *testing.T) {
	genesis, blocks := generateMergeChain(10, true)
	n, ethservice := startEthService(t, genesis, blocks)
	defer n.Close()
	api := NewConsensusAPI(ethservice)
	cfg := api.eth.BlockChain().Config()
	cfg.Optimism = &params.OptimismConfig{}
	if !cfg.IsOptimism() {
		t.Fatalf("expected optimism config")
	}
	fcState := engine.ForkchoiceStateV1{
		HeadBlockHash: common.Hash{42},
	}
	_, _ = api.forkchoiceUpdated(fcState, nil, engine.PayloadV3, false)
}
