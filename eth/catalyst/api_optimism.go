package catalyst

import (
	"errors"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/consensus/misc/eip1559"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// checkOptimismPayload performs Optimism-specific checks on the payload data (called during [(*ConsensusAPI).newPayload]).
func checkOptimismPayload(params engine.ExecutableData, cfg *params.ChainConfig) error {
	// Canyon
	if cfg.IsCanyon(params.Timestamp) && !cfg.IsIsthmus(params.Timestamp) {
		if params.WithdrawalsRoot == nil || *params.WithdrawalsRoot != types.EmptyWithdrawalsHash {
			return errors.New("withdrawalsRoot not equal to MPT root of empty list post-Canyon and pre-Isthmus")
		}
	}

	// Holocene
	if cfg.IsHolocene(params.Timestamp) {
		if err := eip1559.ValidateHoloceneExtraData(params.ExtraData); err != nil {
			return err
		}
	} else {
		if len(params.ExtraData) > 0 {
			return errors.New("extraData must be empty before Holocene")
		}
	}

	// Isthmus
	if cfg.IsIsthmus(params.Timestamp) {
		if params.Withdrawals == nil || len(params.Withdrawals) != 0 {
			return errors.New("non-empty or nil withdrawals post-isthmus")
		}
		if params.WithdrawalsRoot == nil {
			return errors.New("nil withdrawalsRoot post-isthmus")
		}
	}

	return nil
}

// checkOptimismPayloadAttributes performs Optimism-specific checks on the payload attributes (called during [(*ConsensusAPI).forkChoiceUpdated].
// Will panic if payloadAttributes is nil.
func checkOptimismPayloadAttributes(payloadAttributes *engine.PayloadAttributes, cfg *params.ChainConfig) error {
	if payloadAttributes.GasLimit == nil {
		return errors.New("gasLimit parameter is required")
	}

	// Holocene
	if cfg.IsHolocene(payloadAttributes.Timestamp) {
		if err := eip1559.ValidateHolocene1559Params(payloadAttributes.EIP1559Params); err != nil {
			return err
		}
	} else if len(payloadAttributes.EIP1559Params) != 0 {
		return errors.New("eip155Params not supported prior to Holocene upgrade")
	}

	// Isthmus
	if cfg.IsIsthmus(payloadAttributes.Timestamp) && payloadAttributes.Withdrawals == nil || len(payloadAttributes.Withdrawals) != 0 {
		return errors.New("non-empty or nil withdrawals post-isthmus")
	}

	return nil
}
