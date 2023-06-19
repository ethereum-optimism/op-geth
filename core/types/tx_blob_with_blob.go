package types

import "github.com/ethereum/go-ethereum/crypto/kzg4844"

type BlobTxWithBlobs struct {
	BlobTx      *BlobTx
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

func NewBlobTxWithBlobs(tx *Transaction, blobs []kzg4844.Blob, commitments []kzg4844.Commitment, proofs []kzg4844.Proof) *BlobTxWithBlobs {
	return &BlobTxWithBlobs{
		BlobTx:      tx.inner.(*BlobTx),
		Blobs:       blobs,
		Commitments: commitments,
		Proofs:      proofs,
	}
}
