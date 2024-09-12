package types

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

var (

	// deprecatedTxFuncs should be returned by forks that have deprecated support for a tx type.
	deprecatedTxFuncs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return tx.Hash()
		},
		signatureValues: func(tx *Transaction, sig []byte, signerChainID *big.Int) (r *big.Int, s *big.Int, v *big.Int, err error) {
			return nil, nil, nil, fmt.Errorf("%w %v", ErrDeprecatedTxType, tx.Type())
		},
		sender: func(tx *Transaction, hashFunc func(tx *Transaction, chainID *big.Int) common.Hash, signerChainID *big.Int) (common.Address, error) {
			return common.Address{}, fmt.Errorf("%w %v", ErrDeprecatedTxType, tx.Type())
		},
	}

	// Although celo allowed unprotected transactions it never supported signing
	// them with signers retrieved by MakeSigner or LatestSigner (if you wanted
	// to make an unprotected transaction you needed to use the HomesteadSigner
	// directly), so both hash and signatureValues functions here provide
	// protected values, but sender can accept unprotected transactions. See
	// https://github.com/celo-org/celo-blockchain/pull/1748/files and
	// https://github.com/celo-org/celo-blockchain/issues/1734 and
	// https://github.com/celo-org/celo-proposals/blob/master/CIPs/cip-0050.md
	celoLegacyTxFuncs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return rlpHash(append(baseCeloLegacyTxSigningFields(tx), chainID, uint(0), uint(0)))
		},
		signatureValues: func(tx *Transaction, sig []byte, signerChainID *big.Int) (r *big.Int, s *big.Int, v *big.Int, err error) {
			r, s, v = decodeSignature(sig)
			if signerChainID.Sign() != 0 {
				v = big.NewInt(int64(sig[64] + 35))
				signerChainMul := new(big.Int).Mul(signerChainID, big.NewInt(2))
				v.Add(v, signerChainMul)
			}
			return r, s, v, nil
		},
		sender: func(tx *Transaction, hashFunc func(tx *Transaction, chainID *big.Int) common.Hash, signerChainID *big.Int) (common.Address, error) {
			if tx.Protected() {
				if tx.ChainId().Cmp(signerChainID) != 0 {
					return common.Address{}, fmt.Errorf("%w: have %d want %d", ErrInvalidChainId, tx.ChainId(), signerChainID)
				}
				v, r, s := tx.RawSignatureValues()
				signerChainMul := new(big.Int).Mul(signerChainID, big.NewInt(2))
				v = new(big.Int).Sub(v, signerChainMul)
				v.Sub(v, big8)
				return recoverPlain(hashFunc(tx, signerChainID), r, s, v, true)
			} else {
				v, r, s := tx.RawSignatureValues()
				return recoverPlain(rlpHash(baseCeloLegacyTxSigningFields(tx)), r, s, v, true)
			}
		},
	}

	accessListTxFuncs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return NewEIP2930Signer(chainID).Hash(tx)
		},
		signatureValues: func(tx *Transaction, sig []byte, signerChainID *big.Int) (r *big.Int, s *big.Int, v *big.Int, err error) {
			return NewEIP2930Signer(signerChainID).SignatureValues(tx, sig)
		},
		sender: func(tx *Transaction, hashFunc func(tx *Transaction, chainID *big.Int) common.Hash, signerChainID *big.Int) (common.Address, error) {
			return NewEIP2930Signer(tx.ChainId()).Sender(tx)
		},
	}

	dynamicFeeTxFuncs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return NewLondonSigner(chainID).Hash(tx)
		},
		signatureValues: func(tx *Transaction, sig []byte, signerChainID *big.Int) (r *big.Int, s *big.Int, v *big.Int, err error) {
			return NewLondonSigner(signerChainID).SignatureValues(tx, sig)
		},
		sender: func(tx *Transaction, hashFunc func(tx *Transaction, chainID *big.Int) common.Hash, signerChainID *big.Int) (common.Address, error) {
			return NewLondonSigner(signerChainID).Sender(tx)
		},
	}

	celoDynamicFeeTxFuncs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return prefixedRlpHash(
				tx.Type(),
				[]interface{}{
					chainID,
					tx.Nonce(),
					tx.GasTipCap(),
					tx.GasFeeCap(),
					tx.Gas(),
					tx.FeeCurrency(),
					tx.GatewayFeeRecipient(),
					tx.GatewayFee(),
					tx.To(),
					tx.Value(),
					tx.Data(),
					tx.AccessList(),
				})
		},
		signatureValues: dynamicAndDenominatedTxSigValues,
		sender:          dynamicAndDenominatedTxSender,
	}

	// Custom signing functionality for CeloDynamicFeeTxV2 txs.
	celoDynamicFeeTxV2Funcs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return prefixedRlpHash(tx.Type(), baseDynomicatedTxSigningFields(tx, chainID))
		},
		signatureValues: dynamicAndDenominatedTxSigValues,
		sender:          dynamicAndDenominatedTxSender,
	}

	// Custom signing functionality for CeloDenominatedTx txs.
	//
	// TODO remove this nolint directive when we do enable support for cip66 transactions.
	//nolint:unused
	celoDenominatedTxFuncs = &txFuncs{
		hash: func(tx *Transaction, chainID *big.Int) common.Hash {
			return prefixedRlpHash(tx.Type(), append(baseDynomicatedTxSigningFields(tx, chainID), tx.MaxFeeInFeeCurrency()))
		},
		signatureValues: dynamicAndDenominatedTxSigValues,
		sender:          dynamicAndDenominatedTxSender,
	}
)

// txFuncs serves as a container to hold custom signing functionality for a transaction.
//
// TODO consider changing this to an interface, it might make things easier
// because then I could store custom bits of data relevant to each tx type /
// signer such as the signerChainMul. It would also solve the problem of having
// to pass the hash function into the sender function.
type txFuncs struct {
	hash            func(tx *Transaction, chainID *big.Int) common.Hash
	signatureValues func(tx *Transaction, sig []byte, signerChainID *big.Int) (r *big.Int, s *big.Int, v *big.Int, err error)
	sender          func(tx *Transaction, hashFunc func(tx *Transaction, chainID *big.Int) common.Hash, signerChainID *big.Int) (common.Address, error)
}

// Returns the signature values for CeloDynamicFeeTxV2 and CeloDenominatedTx
// transactions.
func dynamicAndDenominatedTxSigValues(tx *Transaction, sig []byte, signerChainID *big.Int) (r *big.Int, s *big.Int, v *big.Int, err error) {
	// Check that chain ID of tx matches the signer. We also accept ID zero here,
	// because it indicates that the chain ID was not specified in the tx.
	chainID := tx.inner.chainID()
	if chainID.Sign() != 0 && chainID.Cmp(signerChainID) != 0 {
		return nil, nil, nil, ErrInvalidChainId
	}
	r, s, _ = decodeSignature(sig)
	v = big.NewInt(int64(sig[64]))
	return r, s, v, nil
}

// Returns the sender for CeloDynamicFeeTxV2 and CeloDenominatedTx
// transactions.
func dynamicAndDenominatedTxSender(tx *Transaction, hashFunc func(tx *Transaction, chainID *big.Int) common.Hash, signerChainID *big.Int) (common.Address, error) {
	if tx.ChainId().Cmp(signerChainID) != 0 {
		return common.Address{}, ErrInvalidChainId
	}
	V, R, S := tx.RawSignatureValues()
	// DynamicFee txs are defined to use 0 and 1 as their recovery
	// id, add 27 to become equivalent to unprotected Homestead signatures.
	V = new(big.Int).Add(V, big.NewInt(27))
	return recoverPlain(hashFunc(tx, signerChainID), R, S, V, true)
}

// Extracts the common signing fields for CeloLegacy and CeloDynamicFeeTx
// transactions.
func baseCeloLegacyTxSigningFields(tx *Transaction) []interface{} {
	return []interface{}{
		tx.Nonce(),
		tx.GasPrice(),
		tx.Gas(),
		tx.FeeCurrency(),
		tx.GatewayFeeRecipient(),
		tx.GatewayFee(),
		tx.To(),
		tx.Value(),
		tx.Data(),
	}
}

// Extracts the common signing fields for CeloDynamicFeeTxV2 and
// CeloDenominatedTx transactions.
func baseDynomicatedTxSigningFields(tx *Transaction, chainID *big.Int) []interface{} {
	return []interface{}{
		chainID,
		tx.Nonce(),
		tx.GasTipCap(),
		tx.GasFeeCap(),
		tx.Gas(),
		tx.To(),
		tx.Value(),
		tx.Data(),
		tx.AccessList(),
		tx.FeeCurrency(),
	}
}
