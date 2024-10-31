package interoptypes

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
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
