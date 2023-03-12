package crypto

import "github.com/ethereum/go-ethereum/crypto/secp256r1"

func VerifySignatureR1(pubKey, dataHash, signature []byte) []byte {
	if secp256r1.VerifySignature(pubKey, dataHash, signature) {
		return []byte{1}
	}
	return []byte{0}
}
