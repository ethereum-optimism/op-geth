package interoptypes

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

var ExecutingMessageEventTopic = crypto.Keccak256Hash([]byte("ExecutingMessage(bytes32,(address,uint256,uint256,uint256,uint256))"))

type Message struct {
	Identifier  Identifier  `json:"identifier"`
	PayloadHash common.Hash `json:"payloadHash"`
}

func (m *Message) DecodeEvent(topics []common.Hash, data []byte) error {
	if len(topics) != 2 { // event hash, indexed payloadHash
		return fmt.Errorf("unexpected number of event topics: %d", len(topics))
	}
	if topics[0] != ExecutingMessageEventTopic {
		return fmt.Errorf("unexpected event topic %q", topics[0])
	}
	if len(data) != 32*5 {
		return fmt.Errorf("unexpected identifier data length: %d", len(data))
	}
	take := func(length uint) []byte {
		taken := data[:length]
		data = data[length:]
		return taken
	}
	takeZeroes := func(length uint) error {
		for _, v := range take(length) {
			if v != 0 {
				return errors.New("expected zero")
			}
		}
		return nil
	}
	if err := takeZeroes(12); err != nil {
		return fmt.Errorf("invalid address padding: %w", err)
	}
	m.Identifier.Origin = common.Address(take(20))
	if err := takeZeroes(32 - 8); err != nil {
		return fmt.Errorf("invalid block number padding: %w", err)
	}
	m.Identifier.BlockNumber = binary.BigEndian.Uint64(take(8))
	if err := takeZeroes(32 - 8); err != nil {
		return fmt.Errorf("invalid log index padding: %w", err)
	}
	m.Identifier.LogIndex = binary.BigEndian.Uint64(take(8))
	if err := takeZeroes(32 - 8); err != nil {
		return fmt.Errorf("invalid timestamp padding: %w", err)
	}
	m.Identifier.Timestamp = binary.BigEndian.Uint64(take(8))
	m.Identifier.ChainID.SetBytes32(take(32))
	m.PayloadHash = topics[1]
	return nil
}

func ExecutingMessagesFromLogs(logs []*types.Log) ([]Message, error) {
	var executingMessages []Message
	for i, l := range logs {
		if l.Address == params.InteropCrossL2InboxAddress {
			var msg Message
			if err := msg.DecodeEvent(l.Topics, l.Data); err != nil {
				return nil, fmt.Errorf("invalid executing message %d, tx-log %d: %w", len(executingMessages), i, err)
			}
			executingMessages = append(executingMessages, msg)
		}
	}
	return executingMessages, nil
}

type Identifier struct {
	Origin      common.Address
	BlockNumber uint64
	LogIndex    uint64
	Timestamp   uint64
	ChainID     uint256.Int // flat, not a pointer, to make Identifier safe as map key
}

type identifierMarshaling struct {
	Origin      common.Address `json:"origin"`
	BlockNumber hexutil.Uint64 `json:"blockNumber"`
	LogIndex    hexutil.Uint64 `json:"logIndex"`
	Timestamp   hexutil.Uint64 `json:"timestamp"`
	ChainID     hexutil.U256   `json:"chainID"`
}

func (id Identifier) MarshalJSON() ([]byte, error) {
	var enc identifierMarshaling
	enc.Origin = id.Origin
	enc.BlockNumber = hexutil.Uint64(id.BlockNumber)
	enc.LogIndex = hexutil.Uint64(id.LogIndex)
	enc.Timestamp = hexutil.Uint64(id.Timestamp)
	enc.ChainID = (hexutil.U256)(id.ChainID)
	return json.Marshal(&enc)
}

func (id *Identifier) UnmarshalJSON(input []byte) error {
	var dec identifierMarshaling
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	id.Origin = dec.Origin
	id.BlockNumber = uint64(dec.BlockNumber)
	id.LogIndex = uint64(dec.LogIndex)
	id.Timestamp = uint64(dec.Timestamp)
	id.ChainID = (uint256.Int)(dec.ChainID)
	return nil
}

type SafetyLevel string

func (lvl SafetyLevel) String() string {
	return string(lvl)
}

func (lvl SafetyLevel) Valid() bool {
	switch lvl {
	case Finalized, Safe, CrossUnsafe, Unsafe:
		return true
	default:
		return false
	}
}

func (lvl SafetyLevel) MarshalText() ([]byte, error) {
	return []byte(lvl), nil
}

func (lvl *SafetyLevel) UnmarshalText(text []byte) error {
	if lvl == nil {
		return errors.New("cannot unmarshal into nil SafetyLevel")
	}
	x := SafetyLevel(text)
	if !x.Valid() {
		return fmt.Errorf("unrecognized safety level: %q", text)
	}
	*lvl = x
	return nil
}

const (
	Finalized   SafetyLevel = "finalized"
	Safe        SafetyLevel = "safe"
	CrossUnsafe SafetyLevel = "cross-unsafe"
	Unsafe      SafetyLevel = "unsafe"
)
