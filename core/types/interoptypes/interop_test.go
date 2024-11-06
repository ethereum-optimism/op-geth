package interoptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

func FuzzMessage_DecodeEvent(f *testing.F) {
	f.Fuzz(func(t *testing.T, validEvTopic bool, numTopics uint8, data []byte) {
		if len(data) < 32 {
			return
		}
		if len(data) > 100_000 {
			return
		}
		if validEvTopic { // valid even signature topic implies a topic to be there
			numTopics += 1
		}
		if numTopics > 4 { // There can be no more than 4 topics per log event
			return
		}
		if int(numTopics)*32 > len(data) {
			return
		}
		var topics []common.Hash
		if validEvTopic {
			topics = append(topics, ExecutingMessageEventTopic)
		}
		for i := 0; i < int(numTopics); i++ {
			var topic common.Hash
			copy(topic[:], data[:])
			data = data[32:]
		}
		require.NotPanics(t, func() {
			var m Message
			_ = m.DecodeEvent(topics, data)
		})
	})
}

func TestSafetyLevel(t *testing.T) {
	require.True(t, Invalid.wellFormatted())
	require.True(t, Unsafe.wellFormatted())
	require.True(t, CrossUnsafe.wellFormatted())
	require.True(t, LocalSafe.wellFormatted())
	require.True(t, Safe.wellFormatted())
	require.True(t, Finalized.wellFormatted())
	require.False(t, SafetyLevel("hello").wellFormatted())
	require.False(t, SafetyLevel("").wellFormatted())
}

func TestInteropMessageFormatEdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		log           *types.Log
		expectedError string
	}{
		{
			name: "Empty Topics",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics:  []common.Hash{},
				Data:    make([]byte, 32*5),
			},
			expectedError: "unexpected number of event topics: 0",
		},
		{
			name: "Wrong Event Topic",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash([]byte("wrong topic")),
					common.BytesToHash([]byte("payloadHash")),
				},
				Data: make([]byte, 32*5),
			},
			expectedError: "unexpected event topic",
		},
		{
			name: "Missing PayloadHash Topic",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash(ExecutingMessageEventTopic[:]),
				},
				Data: make([]byte, 32*5),
			},
			expectedError: "unexpected number of event topics: 1",
		},
		{
			name: "Too Many Topics",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash(ExecutingMessageEventTopic[:]),
					common.BytesToHash([]byte("payloadHash")),
					common.BytesToHash([]byte("extra")),
				},
				Data: make([]byte, 32*5),
			},
			expectedError: "unexpected number of event topics: 3",
		},
		{
			name: "Data Too Short",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash(ExecutingMessageEventTopic[:]),
					common.BytesToHash([]byte("payloadHash")),
				},
				Data: make([]byte, 32*4), // One word too short
			},
			expectedError: "unexpected identifier data length: 128",
		},
		{
			name: "Data Too Long",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash(ExecutingMessageEventTopic[:]),
					common.BytesToHash([]byte("payloadHash")),
				},
				Data: make([]byte, 32*6), // One word too long
			},
			expectedError: "unexpected identifier data length: 192",
		},
		{
			name: "Invalid Address Padding",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash(ExecutingMessageEventTopic[:]),
					common.BytesToHash([]byte("payloadHash")),
				},
				Data: func() []byte {
					data := make([]byte, 32*5)
					data[0] = 1 // Add non-zero byte in address padding
					return data
				}(),
			},
			expectedError: "invalid address padding",
		},
		{
			name: "Invalid Block Number Padding",
			log: &types.Log{
				Address: params.InteropCrossL2InboxAddress,
				Topics: []common.Hash{
					common.BytesToHash(ExecutingMessageEventTopic[:]),
					common.BytesToHash([]byte("payloadHash")),
				},
				Data: func() []byte {
					data := make([]byte, 32*5)
					data[32+23] = 1 // Add non-zero byte in block number padding
					return data
				}(),
			},
			expectedError: "invalid block number padding",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var msg Message
			err := msg.DecodeEvent(tt.log.Topics, tt.log.Data)
			if tt.expectedError != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
