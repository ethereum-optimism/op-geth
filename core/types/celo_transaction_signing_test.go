package types

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

// Tests that by default the celo legacy signer will sign transactions in a protected manner.
func TestProtectedCeloLegacyTxSigning(t *testing.T) {
	tx := newCeloTx(t)
	// Configure config and block time to enable the celoLegacy signer, legacy
	// transactions are deprecated after cel2
	cel2Time := uint64(2000)
	config := &params.ChainConfig{
		ChainID:  big.NewInt(10000),
		Cel2Time: &cel2Time,
	}
	number := new(big.Int).SetUint64(100)
	blockTime := uint64(1000)
	s := MakeSigner(config, number, blockTime)

	senderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	signed, err := SignTx(tx, s, senderKey)
	require.NoError(t, err)

	// Check the sender just to be sure that the signing worked correctly
	actualSender, err := Sender(s, signed)
	require.NoError(t, err)
	require.Equal(t, crypto.PubkeyToAddress(senderKey.PublicKey), actualSender)
	// Validate that the transaction is protected
	require.True(t, signed.Protected())
}

// Tests that the celo legacy signer can still derive the sender of an unprotected transaction.
func TestUnprotectedCeloLegacyTxSenderDerivation(t *testing.T) {
	tx := newCeloTx(t)
	// Configure config and block time to enable the celoLegacy signer, legacy
	// transactions are deprecated after cel2
	cel2Time := uint64(2000)
	config := &params.ChainConfig{
		ChainID:  big.NewInt(10000),
		Cel2Time: &cel2Time,
	}
	number := new(big.Int).SetUint64(100)
	blockTime := uint64(1000)
	s := MakeSigner(config, number, blockTime)
	u := &unprotectedSigner{config.ChainID}

	senderKey, err := crypto.GenerateKey()
	require.NoError(t, err)
	// Sign unprotected
	signed, err := SignTx(tx, u, senderKey)
	require.NoError(t, err)

	// Check that the sender can be derived with the signer from MakeSigner
	actualSender, err := Sender(s, signed)
	require.NoError(t, err)
	require.Equal(t, crypto.PubkeyToAddress(senderKey.PublicKey), actualSender)
	// Validate that the transaction is not protected
	require.False(t, signed.Protected())
}

func newCeloTx(t *testing.T) *Transaction {
	return NewTx(&LegacyTx{
		Nonce:    1,
		GasPrice: new(big.Int).SetUint64(10000),
		Gas:      100000,

		FeeCurrency:         randomAddress(t),
		GatewayFee:          new(big.Int).SetUint64(100),
		GatewayFeeRecipient: randomAddress(t),

		To:    randomAddress(t),
		Value: new(big.Int).SetUint64(1000),
		Data:  []byte{},

		CeloLegacy: true,
	})
}

func randomAddress(t *testing.T) *common.Address {
	addr := common.Address{}
	_, err := rand.Read(addr[:])
	require.NoError(t, err)
	return &addr
}

// This signer mimics Homestead signing but for Celo transactions
type unprotectedSigner struct {
	chainID *big.Int
}

// ChainID implements Signer.
func (u *unprotectedSigner) ChainID() *big.Int {
	return u.chainID
}

// Equal implements Signer.
func (u *unprotectedSigner) Equal(Signer) bool {
	panic("unimplemented")
}

// Hash implements Signer.
func (u *unprotectedSigner) Hash(tx *Transaction) common.Hash {
	return rlpHash(baseCeloLegacyTxSigningFields(tx))
}

// Sender implements Signer.
func (u *unprotectedSigner) Sender(tx *Transaction) (common.Address, error) {
	panic("unimplemented")
}

// SignatureValues implements Signer.
func (u *unprotectedSigner) SignatureValues(tx *Transaction, sig []byte) (r *big.Int, s *big.Int, v *big.Int, err error) {
	r, s, v = decodeSignature(sig)
	return r, s, v, nil
}
