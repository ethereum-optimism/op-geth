package types

import (
	"io"

	"github.com/ethereum/go-ethereum/crypto/kzg4844"
	"github.com/ethereum/go-ethereum/rlp"
)

type BlobTxWithBlobs struct {
	Transaction
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

func NewBlobTxWithBlobs(tx *Transaction, blobs []kzg4844.Blob, commitments []kzg4844.Commitment, proofs []kzg4844.Proof) *BlobTxWithBlobs {
	if tx == nil {
		return nil
	}
	return &BlobTxWithBlobs{
		Transaction: *tx,
		Blobs:       blobs,
		Commitments: commitments,
		Proofs:      proofs,
	}
}

type innerType struct {
	BlobTx      *BlobTx
	Blobs       []kzg4844.Blob
	Commitments []kzg4844.Commitment
	Proofs      []kzg4844.Proof
}

func (tx *BlobTxWithBlobs) DecodeRLP(s *rlp.Stream) error {
	var typedTx Transaction
	err := s.Decode(&typedTx)
	if err == nil {
		tx.Transaction = typedTx
		return nil
	}

	var blobTypedTx innerType
	if s.Decode(&blobTypedTx) == nil {
		tx.Transaction = *NewTx(blobTypedTx.BlobTx)
		tx.Blobs = blobTypedTx.Blobs
		tx.Commitments = blobTypedTx.Commitments
		tx.Proofs = blobTypedTx.Proofs
		return nil
	}

	return err
}

func (tx *BlobTxWithBlobs) EncodeRLP(w io.Writer) error {
	if tx.Transaction.Type() != BlobTxType {
		return tx.Transaction.EncodeRLP(w)
	}
	inner := innerType{
		BlobTx:      tx.inner.(*BlobTx),
		Blobs:       tx.Blobs,
		Commitments: tx.Commitments,
		Proofs:      tx.Proofs,
	}
	w.Write([]byte{tx.Type()})
	return rlp.Encode(w, inner)
}
