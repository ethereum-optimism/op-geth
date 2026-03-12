package types

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const PostExecTxType = 0x7D

// PostExecTx is a synthetic OP Stack transaction used to carry post-exec metadata.
type PostExecTx struct {
	Data []byte
}

func (tx *PostExecTx) copy() TxData {
	return &PostExecTx{Data: common.CopyBytes(tx.Data)}
}

func (tx *PostExecTx) txType() byte           { return PostExecTxType }
func (tx *PostExecTx) chainID() *big.Int      { return common.Big0 }
func (tx *PostExecTx) accessList() AccessList { return nil }
func (tx *PostExecTx) data() []byte           { return tx.Data }
func (tx *PostExecTx) gas() uint64            { return 0 }
func (tx *PostExecTx) gasFeeCap() *big.Int    { return new(big.Int) }
func (tx *PostExecTx) gasTipCap() *big.Int    { return new(big.Int) }
func (tx *PostExecTx) gasPrice() *big.Int     { return new(big.Int) }
func (tx *PostExecTx) value() *big.Int        { return new(big.Int) }
func (tx *PostExecTx) nonce() uint64          { return 0 }
func (tx *PostExecTx) to() *common.Address    { return nil }
func (tx *PostExecTx) isSystemTx() bool       { return false }

func (tx *PostExecTx) effectiveGasPrice(dst *big.Int, baseFee *big.Int) *big.Int {
	return dst.Set(new(big.Int))
}

func (tx *PostExecTx) effectiveNonce() *uint64 { return nil }

func (tx *PostExecTx) sigHash(*big.Int) common.Hash {
	panic("post-exec transaction cannot be signed")
}

func (tx *PostExecTx) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *PostExecTx) setSignatureValues(chainID, v, r, s *big.Int) {}

func (tx *PostExecTx) encode(b *bytes.Buffer) error {
	_, err := b.Write(tx.Data)
	return err
}

func (tx *PostExecTx) decode(input []byte) error {
	tx.Data = common.CopyBytes(input)
	return nil
}
