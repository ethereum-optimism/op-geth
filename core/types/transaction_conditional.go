package types

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
)

// KnownAccounts represents a set of KnownAccounts
type KnownAccounts map[common.Address]KnownAccount

// KnownAccount allows for a user to express their preference of a known
// prestate at a particular account. Only one of the storage root or
// storage slots is allowed to be set. If the storage root is set, then
// the user prefers their transaction to only be included in a block if
// the account's storage root matches. If the storage slots are set,
// then the user prefers their transaction to only be included if the
// particular storage slot values from state match.
type KnownAccount struct {
	StorageRoot  *common.Hash
	StorageSlots map[common.Hash]common.Hash
}

// UnmarshalJSON will parse the JSON bytes into a KnownAccount struct.
func (ka *KnownAccount) UnmarshalJSON(data []byte) error {
	var hash common.Hash
	if err := json.Unmarshal(data, &hash); err == nil {
		ka.StorageRoot = &hash
		ka.StorageSlots = make(map[common.Hash]common.Hash)
		return nil
	}

	var mapping map[common.Hash]common.Hash
	if err := json.Unmarshal(data, &mapping); err != nil {
		return err
	}
	ka.StorageSlots = mapping
	return nil
}

// MarshalJSON will serialize the KnownAccount into JSON bytes.
func (ka KnownAccount) MarshalJSON() ([]byte, error) {
	if ka.StorageRoot != nil {
		return json.Marshal(ka.StorageRoot)
	}
	return json.Marshal(ka.StorageSlots)
}

// Root will return the storage root and true when the user prefers
// execution against an account's storage root, otherwise it will
// return false.
func (ka *KnownAccount) Root() (common.Hash, bool) {
	if ka.StorageRoot == nil {
		return common.Hash{}, false
	}
	return *ka.StorageRoot, true
}

// Slots will return the storage slots and true when the user prefers
// execution against an account's particular storage slots, StorageRoot == nil,
// otherwise it will return false.
func (ka *KnownAccount) Slots() (map[common.Hash]common.Hash, bool) {
	if ka.StorageRoot != nil {
		return ka.StorageSlots, false
	}
	return ka.StorageSlots, true
}

//go:generate go run github.com/fjl/gencodec -type TransactionConditional -field-override transactionConditionalMarshalling -out gen_transaction_conditional_json.go

// TransactionConditional represents the preconditions that determine the
// inclusion of the transaction, enforced out-of-protocol by the sequencer.
type TransactionConditional struct {
	// KnownAccounts represents account prestate conditions
	KnownAccounts KnownAccounts `json:"knownAccounts"`

	// Header state conditionals (inclusive ranges)
	BlockNumberMin *big.Int `json:"blockNumberMin,omitempty"`
	BlockNumberMax *big.Int `json:"blockNumberMax,omitempty"`
	TimestampMin   *uint64  `json:"timestampMin,omitempty"`
	TimestampMax   *uint64  `json:"timestampMax,omitempty"`
}

// field type overrides for gencodec
type transactionConditionalMarshalling struct {
	BlockNumberMax *math.HexOrDecimal256
	BlockNumberMin *math.HexOrDecimal256
	TimestampMin   *math.HexOrDecimal64
	TimestampMax   *math.HexOrDecimal64
}

// Validate will perform sanity checks on the preconditions. This does not check the aggregate cost of the preconditions.
func (cond *TransactionConditional) Validate() error {
	if cond.BlockNumberMin != nil && cond.BlockNumberMax != nil && cond.BlockNumberMin.Cmp(cond.BlockNumberMax) > 0 {
		return fmt.Errorf("block number minimum constraint must be less than the maximum")
	}
	if cond.TimestampMin != nil && cond.TimestampMax != nil && *cond.TimestampMin > *cond.TimestampMax {
		return fmt.Errorf("timestamp minimum constraint must be less than the maximum")
	}
	return nil
}

// Cost computes the aggregate cost of the preconditions; total number of storage lookups required
func (cond *TransactionConditional) Cost() int {
	cost := 0
	for _, account := range cond.KnownAccounts {
		// default cost to handle empty accounts
		cost += 1
		if _, isRoot := account.Root(); isRoot {
			cost += 1
		}
		if slots, isSlots := account.Slots(); isSlots {
			cost += len(slots)
		}
	}
	if cond.BlockNumberMin != nil || cond.BlockNumberMax != nil {
		cost += 1
	}
	if cond.TimestampMin != nil || cond.TimestampMax != nil {
		cost += 1
	}
	return cost
}
