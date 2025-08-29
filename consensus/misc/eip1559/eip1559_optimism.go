package eip1559

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// ValidateOptimismExtraData validates the Optimism extra data.
// It uses the config and parent time to determine how to do the validation.
func ValidateOptimismExtraData(config *params.ChainConfig, parent *types.Header) error {
	if config.IsMinBaseFee(parent.Time) {
		if err := ValidateMinBaseFeeExtraData(parent.Extra); err != nil {
			return err
		}
	} else if config.IsHolocene(parent.Time) {
		if err := ValidateHoloceneExtraData(parent.Extra); err != nil {
			return err
		}
	}
	return nil
}

// DecodeOptimismExtraData decodes the Optimism extra data.
// It uses the config and parent time to determine how to do the decoding.
func DecodeOptimismExtraData(config *params.ChainConfig, parent *types.Header) (uint64, uint64, uint64) {
	if config.IsMinBaseFee(parent.Time) {
		denominator, elasticity, minBaseFee := DecodeMinBaseFeeExtraData(parent.Extra)
		if denominator == 0 {
			// this shouldn't happen as the ExtraData should have been validated prior
			panic("invalid eip-1559 params in extradata")
		}
		return denominator, elasticity, minBaseFee
	} else if config.IsHolocene(parent.Time) {
		denominator, elasticity := DecodeHoloceneExtraData(parent.Extra)
		return denominator, elasticity, 0
	}
	return 0, 0, 0
}
