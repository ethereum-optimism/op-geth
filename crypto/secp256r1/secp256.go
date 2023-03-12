package secp256r1

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

func VerifySignature(pubKeyByte, dataHashByte, signatureByte []byte) bool {
	pubKey := newPubKey(pubKeyByte)
	if pubKey == nil {
		return false
	}

	return ecdsa.VerifyASN1(pubKey, dataHashByte, signatureByte)
}

func newPubKey(pk []byte) *ecdsa.PublicKey {
	pubKey := new(ecdsa.PublicKey)
	pubKey.Curve = elliptic.P256()
	pubKey.X, pubKey.Y = elliptic.UnmarshalCompressed(pubKey.Curve, pk)

	return pubKey
}
