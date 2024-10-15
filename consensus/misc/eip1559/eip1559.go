// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package eip1559

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
)

// VerifyEIP1559Header verifies some header attributes which were changed in EIP-1559,
// - gas limit check
// - basefee check
func VerifyEIP1559Header(config *params.ChainConfig, parent, header *types.Header) error {
	// Verify that the gas limit remains within allowed bounds
	parentGasLimit := parent.GasLimit
	if !config.IsLondon(parent.Number) {
		parentGasLimit = parent.GasLimit * config.ElasticityMultiplier()
	}
	if config.Optimism == nil { // gasLimit can adjust instantly in optimism
		if err := misc.VerifyGaslimit(parentGasLimit, header.GasLimit); err != nil {
			return err
		}
	}
	// Verify the header is not malformed
	if header.BaseFee == nil {
		return errors.New("header is missing baseFee")
	}
	// Verify the baseFee is correct based on the parent header.
	expectedBaseFee := CalcBaseFee(config, parent, header.Time)
	if header.BaseFee.Cmp(expectedBaseFee) != 0 {
		return fmt.Errorf("invalid baseFee: have %s, want %s, parentBaseFee %s, parentGasUsed %d",
			header.BaseFee, expectedBaseFee, parent.BaseFee, parent.GasUsed)
	}
	return nil
}

// DecodeHolocene1599Params extracts the Holcene 1599 parameters from the encoded form:
// https://github.com/ethereum-optimism/specs/blob/main/specs/protocol/holocene/exec-engine.md#eip1559params-encoding
func DecodeHolocene1559Params(params types.BlockNonce) (uint64, uint64) {
	elasticity := binary.BigEndian.Uint32(params[4:])
	denominator := binary.BigEndian.Uint32(params[:4])
	return uint64(elasticity), uint64(denominator)
}

func EncodeHolocene1559Params(elasticity, denom uint32) types.BlockNonce {
	var nonce types.BlockNonce
	binary.BigEndian.PutUint32(nonce[4:], elasticity)
	binary.BigEndian.PutUint32(nonce[:4], denom)
	return nonce
}

// ValidateHoloceneParams checks if the encoded parameters are valid according to the Holocene
// upgrade.
func ValidateHoloceneParams(params types.BlockNonce) error {
	e, d := DecodeHolocene1559Params(params)
	if e != 0 && d == 0 {
		return errors.New("holocene params cannot have a 0 denominator unless elasticity is also 0")
	}
	return nil
}

// CalcBaseFee calculates the basefee of the header.
// The time belongs to the new block to check which upgrades are active.
func CalcBaseFee(config *params.ChainConfig, parent *types.Header, time uint64) *big.Int {
	// If the current block is the first EIP-1559 block, return the InitialBaseFee.
	if !config.IsLondon(parent.Number) {
		return new(big.Int).SetUint64(params.InitialBaseFee)
	}
	elasticity := config.ElasticityMultiplier()
	denominator := config.BaseFeeChangeDenominator(time)
	if config.IsHolocene(time) {
		// Holocene requires we get the 1559 parameters from the nonce field of the parent header
		// unless the field is zero, in which case we use the Canyon values.
		if parent.Nonce != types.BlockNonce([8]byte{}) {
			elasticity, denominator = DecodeHolocene1559Params(parent.Nonce)
		}
	}
	parentGasTarget := parent.GasLimit / elasticity
	// If the parent gasUsed is the same as the target, the baseFee remains unchanged.
	if parent.GasUsed == parentGasTarget {
		return new(big.Int).Set(parent.BaseFee)
	}

	var (
		num   = new(big.Int)
		denom = new(big.Int)
	)

	if parent.GasUsed > parentGasTarget {
		// If the parent block used more gas than its target, the baseFee should increase.
		// max(1, parentBaseFee * gasUsedDelta / parentGasTarget / baseFeeChangeDenominator)
		num.SetUint64(parent.GasUsed - parentGasTarget)
		num.Mul(num, parent.BaseFee)
		num.Div(num, denom.SetUint64(parentGasTarget))
		num.Div(num, denom.SetUint64(denominator))
		baseFeeDelta := math.BigMax(num, common.Big1)

		return num.Add(parent.BaseFee, baseFeeDelta)
	} else {
		// Otherwise if the parent block used less gas than its target, the baseFee should decrease.
		// max(0, parentBaseFee * gasUsedDelta / parentGasTarget / baseFeeChangeDenominator)
		num.SetUint64(parentGasTarget - parent.GasUsed)
		num.Mul(num, parent.BaseFee)
		num.Div(num, denom.SetUint64(parentGasTarget))
		num.Div(num, denom.SetUint64(denominator))
		baseFee := num.Sub(parent.BaseFee, num)

		return math.BigMax(baseFee, common.Big0)
	}
}
