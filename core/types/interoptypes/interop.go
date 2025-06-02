package interoptypes

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

type SafetyLevel string

func (lvl SafetyLevel) String() string {
	return string(lvl)
}

// Valid returns if the safety level is a well-formatted safety level.
func (lvl SafetyLevel) wellFormatted() bool {
	switch lvl {
	case Finalized, Safe, LocalSafe, CrossUnsafe, Unsafe, Invalid:
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
	if !x.wellFormatted() {
		return fmt.Errorf("unrecognized safety level: %q", text)
	}
	*lvl = x
	return nil
}

const (
	Finalized   SafetyLevel = "finalized"
	Safe        SafetyLevel = "safe"
	LocalSafe   SafetyLevel = "local-safe"
	CrossUnsafe SafetyLevel = "cross-unsafe"
	Unsafe      SafetyLevel = "unsafe"
	Invalid     SafetyLevel = "invalid"
)

type ExecutingDescriptor struct {
	ChainID   uint256.Int
	Timestamp uint64
	Timeout   uint64
}

type executingDescriptorMarshaling struct {
	ChainID   uint256.Int    `json:"chainID"`
	Timestamp hexutil.Uint64 `json:"timestamp"`
	Timeout   hexutil.Uint64 `json:"timeout"`
}

func (ed ExecutingDescriptor) MarshalJSON() ([]byte, error) {
	var enc executingDescriptorMarshaling
	enc.ChainID = ed.ChainID
	enc.Timestamp = hexutil.Uint64(ed.Timestamp)
	enc.Timeout = hexutil.Uint64(ed.Timeout)
	return json.Marshal(&enc)
}

func (ed *ExecutingDescriptor) UnmarshalJSON(input []byte) error {
	var dec executingDescriptorMarshaling
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	ed.ChainID = dec.ChainID
	ed.Timestamp = uint64(dec.Timestamp)
	ed.Timeout = uint64(dec.Timeout)
	return nil
}

func TxToInteropAccessList(tx *types.Transaction) []common.Hash {
	if tx == nil {
		return nil
	}
	al := tx.AccessList()
	if len(al) == 0 {
		return nil
	}
	var hashes []common.Hash
	for i := range al {
		if al[i].Address == params.InteropCrossL2InboxAddress {
			hashes = append(hashes, al[i].StorageKeys...)
		}
	}
	return hashes
}
