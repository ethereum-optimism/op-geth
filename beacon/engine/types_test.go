// Copyright 2025 The go-ethereum Authors
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

package engine

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

func TestExecutableDataToBlockSetsAmsterdamBlockAccessListHash(t *testing.T) {
	slotNumber := uint64(1)
	data := ExecutableData{
		LogsBloom:     make([]byte, 256),
		BaseFeePerGas: big.NewInt(1),
		Transactions:  make([][]byte, 0),
		SlotNumber:    &slotNumber,
	}
	block, err := ExecutableDataToBlockNoHash(data, []common.Hash{}, nil, nil, types.DefaultBlockConfig)
	if err != nil {
		t.Fatal(err)
	}
	if hash := block.Header().BlockAccessListHash; hash == nil || *hash != types.EmptyBlockAccessListHash {
		t.Fatalf("unexpected block access list hash: %v", hash)
	}
}

func TestBlobs(t *testing.T) {
	var (
		emptyBlob          = new(kzg4844.Blob)
		emptyBlobCommit, _ = kzg4844.BlobToCommitment(emptyBlob)
		emptyBlobProof, _  = kzg4844.ComputeBlobProof(emptyBlob, emptyBlobCommit)
		emptyCellProof, _  = kzg4844.ComputeCellProofs(emptyBlob)
	)
	header := types.Header{}
	block := types.NewBlock(&header, &types.Body{Withdrawals: []*types.Withdrawal{}}, nil, nil, types.DefaultBlockConfig)

	sidecarWithoutCellProofs := types.NewBlobTxSidecar(types.BlobSidecarVersion0, []kzg4844.Blob{*emptyBlob}, []kzg4844.Commitment{emptyBlobCommit}, []kzg4844.Proof{emptyBlobProof})
	env := BlockToExecutableData(block, common.Big0, []*types.BlobTxSidecar{sidecarWithoutCellProofs}, nil)
	if len(env.BlobsBundle.Proofs) != 1 {
		t.Fatalf("Expect 1 proof in blobs bundle, got %v", len(env.BlobsBundle.Proofs))
	}

	sidecarWithCellProofs := types.NewBlobTxSidecar(types.BlobSidecarVersion0, []kzg4844.Blob{*emptyBlob}, []kzg4844.Commitment{emptyBlobCommit}, emptyCellProof)
	env = BlockToExecutableData(block, common.Big0, []*types.BlobTxSidecar{sidecarWithCellProofs}, nil)
	if len(env.BlobsBundle.Proofs) != 128 {
		t.Fatalf("Expect 128 proofs in blobs bundle, got %v", len(env.BlobsBundle.Proofs))
	}
}
