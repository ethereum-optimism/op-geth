package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const TxPostExecType = 0x7D

// TxPostExec is a synthetic OP Stack transaction used to carry post-exec metadata.
type TxPostExec struct {
	Data []byte
}

func (tx *TxPostExec) copy() TxData {
	return &TxPostExec{Data: common.CopyBytes(tx.Data)}
}

func (tx *TxPostExec) txType() byte           { return TxPostExecType }
func (tx *TxPostExec) chainID() *big.Int      { return common.Big0 }
func (tx *TxPostExec) accessList() AccessList { return nil }
func (tx *TxPostExec) data() []byte           { return tx.Data }
func (tx *TxPostExec) gas() uint64            { return 0 }
func (tx *TxPostExec) gasFeeCap() *big.Int    { return new(big.Int) }
func (tx *TxPostExec) gasTipCap() *big.Int    { return new(big.Int) }
func (tx *TxPostExec) gasPrice() *big.Int     { return new(big.Int) }
func (tx *TxPostExec) value() *big.Int        { return new(big.Int) }
func (tx *TxPostExec) nonce() uint64          { return 0 }
func (tx *TxPostExec) to() *common.Address    { return nil }
func (tx *TxPostExec) isSystemTx() bool       { return false }

func (tx *TxPostExec) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(new(big.Int))
}

func (tx *TxPostExec) effectiveNonce() *uint64 { return nil }

func (tx *TxPostExec) sigHash(*big.Int) common.Hash {
	panic("post-exec transaction cannot be signed")
}

func (tx *TxPostExec) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *TxPostExec) setSignatureValues(chainID, v, r, s *big.Int) {}

func (tx *TxPostExec) encode(b *bytes.Buffer) error {
	_, err := b.Write(tx.Data)
	return err
}

func (tx *TxPostExec) decode(input []byte) error {
	tx.Data = common.CopyBytes(input)
	return nil
}
