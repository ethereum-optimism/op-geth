package catalyst

import (
	"errors"
	"math"
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

func postCanyonPreIsthmus() *params.ChainConfig {
	cfg := new(params.ChainConfig)
	cfg.CanyonTime = new(uint64)
	future := uint64(math.MaxUint64)
	cfg.IsthmusTime = &future
	return cfg
}

func preHolocene() *params.ChainConfig {
	cfg := new(params.ChainConfig)
	return cfg
}

func postHolocene() *params.ChainConfig {
	cfg := new(params.ChainConfig)
	cfg.HoloceneTime = new(uint64)
	return cfg
}

func postIsthmus() *params.ChainConfig {
	cfg := new(params.ChainConfig)
	cfg.HoloceneTime = new(uint64)
	cfg.IsthmusTime = new(uint64)
	return cfg
}

var valid1559Params = []byte{0, 1, 2, 3, 4, 5, 6, 7}
var validExtraData = []byte{0, 1, 2, 3, 4, 5, 6, 7, 8}

func TestCheckOptimismPayload(t *testing.T) {
	tests := []struct {
		name     string
		params   engine.ExecutableData
		cfg      *params.ChainConfig
		expected error
	}{
		{
			name: "valid payload post-Canyon/pre-Isthmus",
			params: engine.ExecutableData{
				Timestamp:       0,
				ExtraData:       []byte{},
				WithdrawalsRoot: &types.EmptyWithdrawalsHash,
			},
			cfg: postCanyonPreIsthmus(),
		},
		{
			name: "nil withdrawalsRoot payload post-Canyon/pre-Isthmus",
			params: engine.ExecutableData{
				Timestamp:       0,
				ExtraData:       []byte{},
				WithdrawalsRoot: nil,
			},
			cfg:      postCanyonPreIsthmus(),
			expected: errors.New("withdrawalsRoot not equal to MPT root of empty list post-Canyon and pre-Isthmus"),
		}, {
			name: "incorrect withdrawalsRoot payload post-Canyon/pre-Isthmus",
			params: engine.ExecutableData{
				Timestamp:       0,
				ExtraData:       []byte{},
				WithdrawalsRoot: &(types.EmptyVerkleHash),
			},
			cfg:      postCanyonPreIsthmus(),
			expected: errors.New("withdrawalsRoot not equal to MPT root of empty list post-Canyon and pre-Isthmus"),
		},
		{
			name: "valid payload pre-Holocene",
			params: engine.ExecutableData{
				Timestamp: 0,
				ExtraData: []byte{},
			},
			cfg:      preHolocene(),
			expected: nil,
		},
		{
			name: "invalid payload pre-Holocene with extraData",
			params: engine.ExecutableData{
				Timestamp: 0,
				ExtraData: []byte{1, 2, 3},
			},
			cfg:      preHolocene(),
			expected: errors.New("extraData must be empty before Holocene"),
		},
		{
			name: "invalid payload pre-Holocene with extraData",
			params: engine.ExecutableData{
				Timestamp: 0,
				ExtraData: []byte{1, 2, 3},
			},
			cfg:      postHolocene(),
			expected: errors.New("holocene extraData should be 9 bytes, got 3"),
		},
		{
			name: "valid payload post-Holocene with extraData",
			params: engine.ExecutableData{
				Timestamp: 0,
				ExtraData: validExtraData,
			},
			cfg:      postHolocene(),
			expected: nil,
		},
		{
			name: "nil withdrawals post-isthmus",
			params: engine.ExecutableData{
				Timestamp:   0,
				Withdrawals: nil,
				ExtraData:   validExtraData,
			},
			cfg:      postIsthmus(),
			expected: errors.New("non-empty or nil withdrawals post-isthmus"),
		},
		{
			name: "non-empty withdrawals post-isthmus",
			params: engine.ExecutableData{
				Timestamp:   0,
				Withdrawals: make([]*types.Withdrawal, 1),
				ExtraData:   validExtraData,
			},
			cfg:      postIsthmus(),
			expected: errors.New("non-empty or nil withdrawals post-isthmus"),
		},
		{
			name: "nil withdrawals root post-isthmus",
			params: engine.ExecutableData{
				Timestamp:   0,
				Withdrawals: make([]*types.Withdrawal, 0),
				ExtraData:   validExtraData,
			},
			cfg:      postIsthmus(),
			expected: errors.New("nil withdrawalsRoot post-isthmus"),
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
			name:              "nil payload attributes",
			payloadAttributes: nil,
			cfg:               preHolocene(),
			shouldPanic:       true,
		},
		{
			name: "valid payload attributes pre-Holocene",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp: 0,
				GasLimit:  new(uint64),
			},
			cfg:      preHolocene(),
			expected: nil,
		},
		{
			name: "invalid payload attributes pre-Holocene with gasLimit",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp: 0,
				GasLimit:  nil,
			},
			cfg:      preHolocene(),
			expected: errors.New("gasLimit parameter is required"),
		},
		{
			name: "invalid payload attributes pre-Holocene with eip1559Params",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp:     0,
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
			},
			cfg:      preHolocene(),
			expected: errors.New("eip155Params not supported prior to Holocene upgrade"),
		},
		{
			name: "valid payload attributes post-Holocene with gasLimit",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp:     0,
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
			},
			cfg:      postHolocene(),
			expected: nil,
		},
		{
			name: "non-empty withdrawals post-isthmus",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp:     0,
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
				Withdrawals:   make([]*types.Withdrawal, 1),
			},
			cfg:      postIsthmus(),
			expected: errors.New("non-empty or nil withdrawals post-isthmus"),
		},
		{
			name: "nil withdrawals post-isthmus",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp:     0,
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
				Withdrawals:   nil,
			},
			cfg:      postIsthmus(),
			expected: errors.New("non-empty or nil withdrawals post-isthmus"),
		},
		{
			name: "empty withdrawals post-isthmus",
			payloadAttributes: &engine.PayloadAttributes{
				Timestamp:     0,
				GasLimit:      new(uint64),
				EIP1559Params: valid1559Params,
				Withdrawals:   make([]*types.Withdrawal, 0),
			},
			cfg:      postIsthmus(),
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
