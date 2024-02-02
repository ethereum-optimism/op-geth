package downloader

import (
	"bytes"
	"embed"
	"encoding/gob"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
)

type depositNonces struct {
	ChainID uint64              `json:"chainID"`
	Start   uint64              `json:"start"`
	End     uint64              `json:"end"`
	Results map[uint64][]uint64 `json:"results"`
}

var receiptReferences = make(map[uint64]*depositNonces)

func initReceiptReferences() {
	// lazy load the receipt references
	if len(receiptReferences) == 0 {
		//go:embed receipt_reference/*.gob
		var receiptReference embed.FS
		files, _ := receiptReference.ReadDir(".")
		for _, file := range files {
			// load to map
			bs, _ := receiptReference.ReadFile(file.Name())
			ref := new(depositNonces)
			gob.NewDecoder(bytes.NewReader(bs)).Decode(&ref)
			receiptReferences[ref.ChainID] = ref
		}
	}
}

func correctReceipts(receipts types.Receipts, transactions types.Transactions, blockNumber uint64) types.Receipts {
	initReceiptReferences()
	if len(transactions) == 0 {
		return receipts
	}
	cid := transactions[0].ChainId().Uint64()
	refs, ok := receiptReferences[cid]
	if !ok {
		log.Warn("No data source for chain", "chainID", cid)
	}
	// only correct if the block is within the range
	if blockNumber < refs.Start || blockNumber > refs.End {
		log.Debug("Block is out of range for receipt reference", "blockNumber", blockNumber, "start", refs.Start, "end", refs.End)
		return receipts
	}
	// get the block nonces
	blockNonces, ok := refs.Results[blockNumber]
	if !ok {
		log.Warn("User deposit found, but no block data found", "blockNumber", blockNumber)
		return receipts
	}
	touched := 0
	for i := 0; i < len(receipts); i++ {
		r := receipts[i]
		tx := transactions[i]
		from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
		if err != nil {
			log.Warn("Failed to determine sender", "err", err)
			continue
		}
		if r.Type == 126 && from != common.HexToAddress("0xDeaDDEaDDeAdDeAdDEAdDEaddeAddEAdDEAd0001") {
			nonce := blockNonces[touched]
			touched++
			if nonce != *r.DepositNonce {
				log.Warn("Corrected deposit nonce", "nonce", *r.DepositNonce, "corrected", nonce)
				r.DepositNonce = &nonce
			}
		}
	}
	return receipts
}
