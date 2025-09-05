package catalyst

import (
	"errors"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/params"
)

// checkOptimismPayload performs Optimism-specific checks on the payload data (called during [(*ConsensusAPI).newPayload]).
func checkOptimismPayload(params engine.ExecutableData, cfg *params.ChainConfig) error {
	// (non)-nil withdrawals is already checked by Shanghai rules.
	// Canyon - empty withdrawals
	if cfg.IsCanyon(params.Timestamp) {
		if len(params.Withdrawals) != 0 {
			return errors.New("non-empty withdrawals post-Canyon")
		}
	}

	// ExtraData validation
	if err := eip1559.ValidateOptimismExtraData(cfg, params.Timestamp, params.ExtraData); err != nil {
		return err
	}

	// Isthmus - withdrawalsRoot
	if cfg.IsIsthmus(params.Timestamp) {
		if params.WithdrawalsRoot == nil {
			return errors.New("nil withdrawalsRoot post-Isthmus")
		}
	} else if params.WithdrawalsRoot != nil { // pre-Isthmus
		return errors.New("non-nil withdrawalsRoot pre-Isthmus")
	}

	return nil
}

// checkOptimismPayloadAttributes performs Optimism-specific checks on the payload attributes (called during [(*ConsensusAPI).forkChoiceUpdated].
// Will panic if payloadAttributes is nil.
func checkOptimismPayloadAttributes(payloadAttributes *engine.PayloadAttributes, cfg *params.ChainConfig) error {
	if payloadAttributes.GasLimit == nil {
		return errors.New("gasLimit parameter is required")
	}

	// (non)-nil withdrawals is already checked by Shanghai rules.
	// Canyon - empty withdrawals
	if cfg.IsCanyon(payloadAttributes.Timestamp) {
		if len(payloadAttributes.Withdrawals) != 0 {
			return errors.New("non-empty withdrawals post-Canyon")
		}
	}

	// Holocene - extraData
	if cfg.IsHolocene(payloadAttributes.Timestamp) {
		if err := eip1559.ValidateHolocene1559Params(payloadAttributes.EIP1559Params); err != nil {
			return err
		}
	} else if len(payloadAttributes.EIP1559Params) != 0 { // pre-Holocene
		return errors.New("non-empty eip155Params pre-Holocene")
	}

	// Note: PayloadAttributes don't contain the Isthmus withdrawalsRoot, it's set during block assembly.

	return nil
}
