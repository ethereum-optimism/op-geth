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
	"bytes"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/kzg4844"
)

func TestGlamsterdamEngineWireFields(t *testing.T) {
	targetGasLimit := uint64(60_000_000)
	attrs := PayloadAttributes{TargetGasLimit: &targetGasLimit}
	encodedAttrs, err := json.Marshal(attrs)
	if err != nil {
		t.Fatal(err)
	}
	var decodedAttrs PayloadAttributes
	if err := json.Unmarshal(encodedAttrs, &decodedAttrs); err != nil {
		t.Fatal(err)
	}
	if decodedAttrs.TargetGasLimit == nil || *decodedAttrs.TargetGasLimit != targetGasLimit {
		t.Fatalf("unexpected target gas limit: %v", decodedAttrs.TargetGasLimit)
	}

	emptyBAL := hexutil.Bytes{0xc0}
	payload := ExecutableData{
		LogsBloom:       make([]byte, 256),
		BaseFeePerGas:   big.NewInt(1),
		Transactions:    make([][]byte, 0),
		BlockAccessList: &emptyBAL,
	}
	encodedPayload, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	var decodedPayload ExecutableData
	if err := json.Unmarshal(encodedPayload, &decodedPayload); err != nil {
		t.Fatal(err)
	}
	if decodedPayload.BlockAccessList == nil || !bytes.Equal(*decodedPayload.BlockAccessList, emptyBAL) {
		t.Fatalf("unexpected block access list: %x", decodedPayload.BlockAccessList)
	}
}

func TestExecutableDataToBlockSetsAmsterdamBlockAccessListHash(t *testing.T) {
	slotNumber := uint64(1)
	blockAccessList := hexutil.Bytes{0xc1, 0x80}
	data := ExecutableData{
		LogsBloom:       make([]byte, 256),
		BaseFeePerGas:   big.NewInt(1),
		Transactions:    make([][]byte, 0),
		SlotNumber:      &slotNumber,
		BlockAccessList: &blockAccessList,
	}
	block, err := ExecutableDataToBlockNoHash(data, []common.Hash{}, nil, nil, types.DefaultBlockConfig)
	if err != nil {
		t.Fatal(err)
	}
	want := crypto.Keccak256Hash(blockAccessList)
	if hash := block.Header().BlockAccessListHash; hash == nil || *hash != want {
		t.Fatalf("unexpected block access list hash: have %v, want %v", hash, want)
	}

	data.BlockAccessList = nil
	if _, err := ExecutableDataToBlockNoHash(data, []common.Hash{}, nil, nil, types.DefaultBlockConfig); err == nil {
		t.Fatal("expected missing block access list to be rejected")
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
