// TODO: needs copyright header?

package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

// CeloDynamicFeeTx represents a CIP-64 transaction.
type CeloDynamicFeeTx struct {
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
func (tx *CeloDynamicFeeTx) copy() TxData {
	cpy := &CeloDynamicFeeTx{
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

func (tx *CeloDynamicFeeTx) feeCurrency() *common.Address { return tx.FeeCurrency }
