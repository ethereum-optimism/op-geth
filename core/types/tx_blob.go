// Copyright 2023 The go-ethereum Authors
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

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/params"
)

// BlobTx represents an EIP-4844 transaction.
type BlobTx struct {
	ChainID    *big.Int
	Nonce      uint64
	GasTipCap  *big.Int // a.k.a. maxPriorityFeePerGas
	GasFeeCap  *big.Int // a.k.a. maxFeePerGas
	Gas        uint64
	To         common.Address
	Value      *big.Int
	Data       []byte
	AccessList AccessList
	BlobFeeCap *big.Int // a.k.a. maxFeePerDataGas
	BlobHashes []common.Hash

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

type BlobTxWithBlobs struct {
	*BlobTx
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

func (tx *BlobTxWithBlobs) copy() TxData {
	cpy := BlobTxWithBlobs{
		BlobTx:      tx.BlobTx.copy().(*BlobTx),
		Blobs:       make([]kzg4844.Blob, len(tx.Blobs)),
		Commitments: make([]kzg4844.Commitment, len(tx.Commitments)),
		Proofs:      make([]kzg4844.Proof, len(tx.Proofs)),
	}
	copy(cpy.Blobs, tx.Blobs)
	copy(cpy.Commitments, tx.Commitments)
	copy(cpy.Proofs, tx.Proofs)
	return &cpy
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *BlobTx) copy() TxData {
	cpy := &BlobTx{
		Nonce: tx.Nonce,
		To:    tx.To,
		Data:  common.CopyBytes(tx.Data),
		Gas:   tx.Gas,
		// These are copied below.
		AccessList: make(AccessList, len(tx.AccessList)),
		BlobHashes: make([]common.Hash, len(tx.BlobHashes)),
		Value:      new(big.Int),
		ChainID:    new(big.Int),
		GasTipCap:  new(big.Int),
		GasFeeCap:  new(big.Int),
		BlobFeeCap: new(big.Int),
		V:          new(big.Int),
		R:          new(big.Int),
		S:          new(big.Int),
	}
	copy(cpy.AccessList, tx.AccessList)
	copy(cpy.BlobHashes, tx.BlobHashes)

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
	if tx.BlobFeeCap != nil {
		cpy.BlobFeeCap.Set(tx.BlobFeeCap)
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
func (tx *BlobTx) txType() byte              { return BlobTxType }
func (tx *BlobTx) chainID() *big.Int         { return tx.ChainID }
func (tx *BlobTx) accessList() AccessList    { return tx.AccessList }
func (tx *BlobTx) data() []byte              { return tx.Data }
func (tx *BlobTx) gas() uint64               { return tx.Gas }
func (tx *BlobTx) gasFeeCap() *big.Int       { return tx.GasFeeCap }
func (tx *BlobTx) gasTipCap() *big.Int       { return tx.GasTipCap }
func (tx *BlobTx) gasPrice() *big.Int        { return tx.GasFeeCap }
func (tx *BlobTx) value() *big.Int           { return tx.Value }
func (tx *BlobTx) nonce() uint64             { return tx.Nonce }
func (tx *BlobTx) to() *common.Address       { tmp := tx.To; return &tmp }
func (tx *BlobTx) blobGas() uint64           { return params.BlobTxDataGasPerBlob * uint64(len(tx.BlobHashes)) }
func (tx *BlobTx) blobGasFeeCap() *big.Int   { return tx.BlobFeeCap }
func (tx *BlobTx) blobHashes() []common.Hash { return tx.BlobHashes }

func (tx *BlobTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	if baseFee == nil {
		return dst.Set(tx.GasFeeCap)
	}
	tip := dst.Sub(tx.GasFeeCap, baseFee)
	if tip.Cmp(tx.GasTipCap) > 0 {
		tip.Set(tx.GasTipCap)
	}
	return tip.Add(tip, baseFee)
}

func (tx *BlobTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *BlobTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID.Set(chainID)
	tx.V.Set(v)
	tx.R.Set(r)
	tx.S.Set(s)
}
