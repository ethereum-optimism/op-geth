package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

const Deposit2TxType = 0x7D

type Deposit2Tx struct {
	// SourceHash uniquely identifies the source of the deposit
	SourceHash common.Hash
	// From is exposed through the types.Signer, not through TxData
	From common.Address
	// nil means contract creation
	To *common.Address `rlp:"nil"`
	// Mint is minted on L2, locked on L1, nil if no minting.
	Mint *big.Int `rlp:"nil"`
	// Value is transferred from L2 balance, executed after Mint (if any)
	Value *big.Int
	// gas limit
	Gas uint64
	// Field indicating if this transaction is exempt from the L2 gas limit.
	IsSystemTransaction bool
	// Normal Tx data
	Data []byte
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *Deposit2Tx) copy() TxData {
	cpy := &Deposit2Tx{
		SourceHash:          tx.SourceHash,
		From:                tx.From,
		To:                  copyAddressPtr(tx.To),
		Mint:                nil,
		Value:               new(big.Int),
		Gas:                 tx.Gas,
		IsSystemTransaction: tx.IsSystemTransaction,
		Data:                common.CopyBytes(tx.Data),
	}
	if tx.Mint != nil {
		cpy.Mint = new(big.Int).Set(tx.Mint)
	}
	if tx.Value != nil {
		cpy.Value.Set(tx.Value)
	}
	return cpy
}

// accessors for innerTx.
func (tx *Deposit2Tx) txType() byte           { return Deposit2TxType }
func (tx *Deposit2Tx) chainID() *big.Int      { return common.Big0 }
func (tx *Deposit2Tx) accessList() AccessList { return nil }
func (tx *Deposit2Tx) data() []byte           { return tx.Data }
func (tx *Deposit2Tx) gas() uint64            { return tx.Gas }
func (tx *Deposit2Tx) gasFeeCap() *big.Int    { return new(big.Int) }
func (tx *Deposit2Tx) gasTipCap() *big.Int    { return new(big.Int) }
func (tx *Deposit2Tx) gasPrice() *big.Int     { return new(big.Int) }
func (tx *Deposit2Tx) value() *big.Int        { return tx.Value }
func (tx *Deposit2Tx) nonce() uint64          { return 0 }
func (tx *Deposit2Tx) to() *common.Address    { return tx.To }
func (tx *Deposit2Tx) isSystemTx() bool       { return tx.IsSystemTransaction }

func (tx *Deposit2Tx) rawSignatureValues() (v, r, s *big.Int) {
	return common.Big0, common.Big0, common.Big0
}

func (tx *Deposit2Tx) setSignatureValues(chainID, v, r, s *big.Int) {
	// this is a noop for deposit transactions
}
