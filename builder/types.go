package builder

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"golang.org/x/exp/slices"
)

type PayloadRequestV1 struct {
	Slot       uint64      `json:"slot"`
	ParentHash common.Hash `json:"parentHash"`
}

type BuilderPayloadAttributes struct {
	Timestamp             hexutil.Uint64    `json:"timestamp"`
	Random                common.Hash       `json:"prevRandao"`
	SuggestedFeeRecipient common.Address    `json:"suggestedFeeRecipient,omitempty"`
	Slot                  uint64            `json:"slot"`
	HeadHash              common.Hash       `json:"blockHash"`
	Withdrawals           types.Withdrawals `json:"withdrawals"`
	ParentBeaconBlockRoot *common.Hash      `json:"parentBeaconBlockRoot"`
	GasLimit              uint64            `json:"gasLimit"`

	NoTxPool     bool                 `json:"noTxPool,omitempty"` // Optimism addition: option to disable tx pool contents from being included
	Transactions []*types.Transaction `json:"transactions"`       // Optimism addition: txs forced into the block via engine API
}

func (attrs *BuilderPayloadAttributes) Equal(other *BuilderPayloadAttributes) bool {
	if attrs.Timestamp != other.Timestamp ||
		attrs.Random != other.Random ||
		attrs.SuggestedFeeRecipient != other.SuggestedFeeRecipient ||
		attrs.Slot != other.Slot ||
		attrs.HeadHash != other.HeadHash ||
		attrs.GasLimit != other.GasLimit ||
		attrs.ParentBeaconBlockRoot != other.ParentBeaconBlockRoot ||
		attrs.NoTxPool != other.NoTxPool {
		return false
	}

	if !slices.Equal(attrs.Withdrawals, other.Withdrawals) {
		return false
	}

	if !slices.Equal(attrs.Transactions, other.Transactions) {
		return false
	}

	return true
}
