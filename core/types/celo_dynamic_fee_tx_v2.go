// Copyright 2024 The Celo Authors
// This file is part of the celo library.
//
// The celo library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The celo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the celo library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

const CeloDynamicFeeTxV2Type = 0x7b

// CeloDynamicFeeTxV2 represents a CIP-64 transaction.
type CeloDynamicFeeTxV2 struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int
	GasFeeCap  *big.Int
	Gas        uint64
	To         *common.Address `rlp:"nil"` // nil means contract creation
	Value      *big.Int
	Data       []byte
	AccessList AccessList

	FeeCurrency *common.Address `rlp:"nil"` // nil means native currency

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *CeloDynamicFeeTxV2) copy() TxData {
	cpy := &CeloDynamicFeeTxV2{
		Nonce:       tx.Nonce,
		To:          copyAddressPtr(tx.To),
		Data:        common.CopyBytes(tx.Data),
		Gas:         tx.Gas,
		FeeCurrency: copyAddressPtr(tx.FeeCurrency),
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasTipCap:  new(big.Int),
		GasFeeCap:  new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasTipCap != nil {
		cpy.GasTipCap.Set(tx.GasTipCap)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *CeloDynamicFeeTxV2) txType() byte           { return CeloDynamicFeeTxV2Type }
func (tx *CeloDynamicFeeTxV2) chainID() *big.Int      { return tx.ChainID }
func (tx *CeloDynamicFeeTxV2) accessList() AccessList { return tx.AccessList }
func (tx *CeloDynamicFeeTxV2) data() []byte           { return tx.Data }
func (tx *CeloDynamicFeeTxV2) gas() uint64            { return tx.Gas }
func (tx *CeloDynamicFeeTxV2) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *CeloDynamicFeeTxV2) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *CeloDynamicFeeTxV2) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *CeloDynamicFeeTxV2) value() *big.Int        { return tx.Value }
func (tx *CeloDynamicFeeTxV2) nonce() uint64          { return tx.Nonce }
func (tx *CeloDynamicFeeTxV2) to() *common.Address    { return tx.To }
func (tx *CeloDynamicFeeTxV2) isSystemTx() bool       { return false }

func (tx *CeloDynamicFeeTxV2) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	if baseFee == nil {
		return dst.Set(tx.GasFeeCap)
	}
	tip := dst.Sub(tx.GasFeeCap, baseFee)
	if tip.Cmp(tx.GasTipCap) > 0 {
		tip.Set(tx.GasTipCap)
	}
	return tip.Add(tip, baseFee)
}

func (tx *CeloDynamicFeeTxV2) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *CeloDynamicFeeTxV2) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}

func (tx *CeloDynamicFeeTxV2) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *CeloDynamicFeeTxV2) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
