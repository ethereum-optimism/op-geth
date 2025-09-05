package eip1559

import (
	"encoding/binary"
	"errors"
	"fmt"
	gomath "math"
)

const HoloceneExtraDataVersionByte = uint8(0x00)
const JovianExtraDataVersionByte = uint8(0x01)

type ForkChecker interface {
	IsHolocene(time uint64) bool
	IsJovian(time uint64) bool
}

// ValidateOptimismExtraData validates the Optimism extra data.
// It uses the config and parent time to determine how to do the validation.
func ValidateOptimismExtraData(fc ForkChecker, time uint64, extraData []byte) error {
	if fc.IsJovian(time) {
		return ValidateJovianExtraData(extraData)
	} else if fc.IsHolocene(time) {
		return ValidateHoloceneExtraData(extraData)
	} else if len(extraData) > 0 { // pre-Holocene
		return errors.New("extraData must be empty before Holocene")
	}
	return nil
}

// DecodeOptimismExtraData decodes the Optimism extra data.
// It uses the config and parent time to determine how to do the decoding.
// The parent.extraData is expected to be valid (i.e. ValidateOptimismExtraData has been called previously)
func DecodeOptimismExtraData(fc ForkChecker, time uint64, extraData []byte) (uint64, uint64, *uint64) {
	if fc.IsJovian(time) {
		denominator, elasticity, minBaseFee := DecodeJovianExtraData(extraData)
		return denominator, elasticity, minBaseFee
	} else if fc.IsHolocene(time) {
		denominator, elasticity := DecodeHoloceneExtraData(extraData)
		return denominator, elasticity, nil
	}
	return 0, 0, nil
}

// EncodeOptimismExtraData encodes the Optimism extra data.
// It uses the config and parent time to determine how to do the encoding.
func EncodeOptimismExtraData(fc ForkChecker, time uint64, denominator, elasticity uint64, minBaseFee *uint64) []byte {
	if fc.IsJovian(time) {
		if minBaseFee == nil {
			panic("minBaseFee cannot be nil for Jovian")
		}
		return EncodeJovianExtraData(denominator, elasticity, *minBaseFee)
	} else if fc.IsHolocene(time) {
		return EncodeHoloceneExtraData(denominator, elasticity)
	} else {
		return nil
	}
}

// DecodeHolocene1559Params extracts the Holcene 1559 parameters from the encoded form defined here:
// https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/holocene/exec-engine.md#eip-1559-parameters-in-payloadattributesv3
//
// Returns 0,0 if the format is invalid, though ValidateHolocene1559Params should be used instead of this function for
// validity checking.
func DecodeHolocene1559Params(params []byte) (uint64, uint64) {
	if len(params) != 8 {
		return 0, 0
	}
	denominator := binary.BigEndian.Uint32(params[:4])
	elasticity := binary.BigEndian.Uint32(params[4:])
	return uint64(denominator), uint64(elasticity)
}

// DecodeHoloceneExtraData decodes the Holocene 1559 parameters from the encoded form defined here:
// https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/holocene/exec-engine.md#eip-1559-parameters-in-block-header
//
// Returns 0,0 if the format is invalid, though ValidateHoloceneExtraData should be used instead of this function for
// validity checking.
func DecodeHoloceneExtraData(extra []byte) (uint64, uint64) {
	if len(extra) != 9 {
		return 0, 0
	}
	return DecodeHolocene1559Params(extra[1:])
}

// EncodeHolocene1559Params encodes the eip-1559 parameters into 'PayloadAttributes.EIP1559Params' format. Will panic if
// either value is outside uint32 range.
func EncodeHolocene1559Params(denom, elasticity uint64) []byte {
	r := make([]byte, 8)
	if denom > gomath.MaxUint32 || elasticity > gomath.MaxUint32 {
		panic("eip-1559 parameters out of uint32 range")
	}
	binary.BigEndian.PutUint32(r[:4], uint32(denom))
	binary.BigEndian.PutUint32(r[4:], uint32(elasticity))
	return r
}

// EncodeHoloceneExtraData encodes the eip-1559 parameters into the header 'ExtraData' format. Will panic if either
// value is outside uint32 range.
func EncodeHoloceneExtraData(denom, elasticity uint64) []byte {
	r := make([]byte, 9)
	if denom > gomath.MaxUint32 || elasticity > gomath.MaxUint32 {
		panic("eip-1559 parameters out of uint32 range")
	}
	// leave version byte 0
	binary.BigEndian.PutUint32(r[1:5], uint32(denom))
	binary.BigEndian.PutUint32(r[5:], uint32(elasticity))
	return r
}

// ValidateHolocene1559Params checks if the encoded parameters are valid according to the Holocene
// upgrade.
func ValidateHolocene1559Params(params []byte) error {
	if len(params) != 8 {
		return fmt.Errorf("holocene eip-1559 params should be 8 bytes, got %d", len(params))
	}
	d, e := DecodeHolocene1559Params(params)
	if e != 0 && d == 0 {
		return errors.New("holocene params cannot have a 0 denominator unless elasticity is also 0")
	}
	return nil
}

// ValidateHoloceneExtraData checks if the header extraData is valid according to the Holocene
// upgrade.
func ValidateHoloceneExtraData(extra []byte) error {
	if len(extra) != 9 {
		return fmt.Errorf("holocene extraData should be 9 bytes, got %d", len(extra))
	}
	if extra[0] != HoloceneExtraDataVersionByte {
		return fmt.Errorf("holocene extraData version byte should be %d, got %d", HoloceneExtraDataVersionByte, extra[0])
	}
	return ValidateHolocene1559Params(extra[1:])
}

// DecodeJovianExtraData decodes the extraData parameters from the encoded form defined here:
// https://specs.optimism.io/protocol/jovian/exec-engine.html
//
// Returns 0,0,nil if the format is invalid, and d, e, nil for the Holocene length, to provide best effort behavior for non-Jovian extradata, though ValidateMinBaseFeeExtraData should be used instead of this function for
// validity checking.
func DecodeJovianExtraData(extra []byte) (uint64, uint64, *uint64) {
	// Best effort to decode the extraData for every block in the chain's history,
	// including blocks before the minimum base fee feature was enabled.
	if len(extra) == 9 {
		// This is Holocene extraData
		denominator, elasticity := DecodeHolocene1559Params(extra[1:9])
		return denominator, elasticity, nil
	} else if len(extra) == 17 {
		// Decode extraData when the minimum base fee fork is enabled
		denominator, elasticity := DecodeHolocene1559Params(extra[1:9])
		minBaseFee := binary.BigEndian.Uint64(extra[9:])
		return denominator, elasticity, &minBaseFee
	}
	return 0, 0, nil
}

// EncodeJovianExtraData encodes the EIP-1559 and minBaseFee parameters into the header 'ExtraData' format.
// Will panic if EIP-1559 parameters are outside uint32 range.
func EncodeJovianExtraData(denom, elasticity, minBaseFee uint64) []byte {
	r := make([]byte, 17)
	if denom > gomath.MaxUint32 || elasticity > gomath.MaxUint32 {
		panic("eip-1559 parameters out of uint32 range")
	}
	r[0] = JovianExtraDataVersionByte
	binary.BigEndian.PutUint32(r[1:5], uint32(denom))
	binary.BigEndian.PutUint32(r[5:9], uint32(elasticity))
	binary.BigEndian.PutUint64(r[9:], minBaseFee)
	return r
}

// ValidateJovianExtraData checks if the header extraData is valid according to the minimum base fee feature.
func ValidateJovianExtraData(extra []byte) error {
	if len(extra) != 17 {
		return fmt.Errorf("jovian extraData should be 17 bytes, got %d", len(extra))
	}
	if extra[0] != JovianExtraDataVersionByte {
		return fmt.Errorf("jovian extraData version byte should be %d, got %d", JovianExtraDataVersionByte, extra[0])
	}
	return ValidateHolocene1559Params(extra[1:9])
}
