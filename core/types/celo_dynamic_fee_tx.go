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

const CeloDynamicFeeTxType = 0x7c

type CeloDynamicFeeTx struct {
	ChainID             *big.Int
	Nonce               uint64
	GasTipCap           *big.Int
	GasFeeCap           *big.Int
	Gas                 uint64
	FeeCurrency         *common.Address `rlp:"nil"` // nil means native currency
	GatewayFeeRecipient *common.Address `rlp:"nil"` // nil means no gateway fee is paid
	GatewayFee          *big.Int        `rlp:"nil"`
	To                  *common.Address `rlp:"nil"` // nil means contract creation
	Value               *big.Int
	Data                []byte
	AccessList          AccessList

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *CeloDynamicFeeTx) copy() TxData {
	cpy := &CeloDynamicFeeTx{
		Nonce:               tx.Nonce,
		To:                  copyAddressPtr(tx.To),
		Data:                common.CopyBytes(tx.Data),
		Gas:                 tx.Gas,
		FeeCurrency:         copyAddressPtr(tx.FeeCurrency),
		GatewayFeeRecipient: copyAddressPtr(tx.GatewayFeeRecipient),
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		GatewayFee: new(big.Int),
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
	if tx.GatewayFee != nil {
		cpy.GatewayFee.Set(tx.GatewayFee)
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
func (tx *CeloDynamicFeeTx) txType() byte           { return CeloDynamicFeeTxType }
func (tx *CeloDynamicFeeTx) chainID() *big.Int      { return tx.ChainID }
func (tx *CeloDynamicFeeTx) accessList() AccessList { return tx.AccessList }
func (tx *CeloDynamicFeeTx) data() []byte           { return tx.Data }
func (tx *CeloDynamicFeeTx) gas() uint64            { return tx.Gas }
func (tx *CeloDynamicFeeTx) gasFeeCap() *big.Int    { return tx.GasFeeCap }
func (tx *CeloDynamicFeeTx) gasTipCap() *big.Int    { return tx.GasTipCap }
func (tx *CeloDynamicFeeTx) gasPrice() *big.Int     { return tx.GasFeeCap }
func (tx *CeloDynamicFeeTx) value() *big.Int        { return tx.Value }
func (tx *CeloDynamicFeeTx) nonce() uint64          { return tx.Nonce }
func (tx *CeloDynamicFeeTx) to() *common.Address    { return tx.To }
func (tx *CeloDynamicFeeTx) isSystemTx() bool       { return false }

func (tx *CeloDynamicFeeTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	if baseFee == nil {
		return dst.Set(tx.GasFeeCap)
	}
	tip := dst.Sub(tx.GasFeeCap, baseFee)
	if tip.Cmp(tx.GasTipCap) > 0 {
		tip.Set(tx.GasTipCap)
	}
	return tip.Add(tip, baseFee)
}

func (tx *CeloDynamicFeeTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *CeloDynamicFeeTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}

func (tx *CeloDynamicFeeTx) encode(b *bytes.Buffer) error {
	return rlp.Encode(b, tx)
}

func (tx *CeloDynamicFeeTx) decode(input []byte) error {
	return rlp.DecodeBytes(input, tx)
}
