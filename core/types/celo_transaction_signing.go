// Copyright 2016 The go-ethereum Authors
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

// TODO(pl)

package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type cel2Signer struct{ cancunSigner }

// NewCel2Signer returns a signer that accepts
// - CIP-64 celo dynamic fee transaction (v2)
// - EIP-4844 blob transactions
// - EIP-1559 dynamic fee transactions
// - EIP-2930 access list transactions,
// - EIP-155 replay protected transactions, and
// - legacy Homestead transactions.
func NewCel2Signer(chainId *big.Int) Signer {
	return cel2Signer{cancunSigner{londonSigner{eip2930Signer{NewEIP155Signer(chainId)}}}}
}

func (s cel2Signer) Sender(tx *Transaction) (common.Address, error) {
	if tx.Type() != CeloDynamicFeeTxType {
		return s.cancunSigner.Sender(tx)
	}
	V, R, S := tx.RawSignatureValues()
	// DynamicFee txs are defined to use 0 and 1 as their recovery
	// id, add 27 to become equivalent to unprotected Homestead signatures.
	V = new(big.Int).Add(V, big.NewInt(27))
	if tx.ChainId().Cmp(s.chainId) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	return recoverPlain(s.Hash(tx), R, S, V, true)
}

func (s cel2Signer) Equal(s2 Signer) bool {
	x, ok := s2.(cel2Signer)
	return ok && x.chainId.Cmp(s.chainId) == 0
}

func (s cel2Signer) SignatureValues(tx *Transaction, sig []byte) (R, S, V *big.Int, err error) {
	txdata, ok := tx.inner.(*CeloDynamicFeeTx)
	if !ok {
		return s.cancunSigner.SignatureValues(tx, sig)
	}
	// Check that chain ID of tx matches the signer. We also accept ID zero here,
	// because it indicates that the chain ID was not specified in the tx.
	if txdata.ChainID.Sign() != 0 && txdata.ChainID.Cmp(s.chainId) != 0 {
		return nil, nil, nil, ErrInvalidChainId
	}
	R, S, _ = decodeSignature(sig)
	V = big.NewInt(int64(sig[64]))
	return R, S, V, nil
}

// Hash returns the hash to be signed by the sender.
// It does not uniquely identify the transaction.
func (s cel2Signer) Hash(tx *Transaction) common.Hash {
	if tx.Type() == CeloDynamicFeeTxType {
		return prefixedRlpHash(
			tx.Type(),
			[]interface{}{
				s.chainId,
				tx.Nonce(),
				tx.GasTipCap(),
				tx.GasFeeCap(),
				tx.Gas(),
				tx.To(),
				tx.Value(),
				tx.Data(),
				tx.AccessList(),
				tx.FeeCurrency(),
			})
	}
	return s.cancunSigner.Hash(tx)
}
