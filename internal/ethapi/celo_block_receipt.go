package ethapi

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// GetBlockReceipt returns "system calls" receipt for the block with the given block hash.
func (s *BlockChainAPI) GetBlockReceipt(ctx context.Context, hash common.Hash) (map[string]interface{}, error) {
	block, err := s.b.BlockByHash(ctx, hash)
	if block == nil || err != nil {
		// If no header with that hash is found, err gives "header for hash not found".
		// But we return nil with no error, to match the behavior of eth_getBlockByHash and eth_getTransactionReceipt in these cases.
		return nil, nil
	}
	index := block.Transactions().Len()
	blockNumber := block.NumberU64()
	receipts, err := s.b.GetReceipts(ctx, block.Hash())
	// GetReceipts() doesn't return an error if things go wrong, so we also check len(receipts)
	if err != nil || len(receipts) < index {
		return nil, err
	}

	var receipt *types.Receipt
	if len(receipts) == index {
		// The block didn't have any logs from system calls and no receipt was created.
		// So we create an empty receipt to return, similarly to how system receipts are created.
		receipt = types.NewReceipt(nil, false, 0)
		receipt.Bloom = types.CreateBloom(types.Receipts{receipt})
	} else {
		receipt = receipts[index]
	}
	return marshalBlockReceipt(receipt, hash, blockNumber, index), nil
}

// marshalBlockReceipt marshals a Celo block receipt into a JSON object. See https://docs.celo.org/developer/migrate/from-ethereum#core-contract-calls
func marshalBlockReceipt(receipt *types.Receipt, blockHash common.Hash, blockNumber uint64, index int) map[string]interface{} {
	fields := map[string]interface{}{
		"blockHash":         blockHash,
		"blockNumber":       hexutil.Uint64(blockNumber),
		"transactionHash":   blockHash,
		"transactionIndex":  hexutil.Uint64(index),
		"from":              common.Address{},
		"to":                nil,
		"gasUsed":           hexutil.Uint64(0),
		"cumulativeGasUsed": hexutil.Uint64(0),
		"contractAddress":   nil,
		"logs":              receipt.Logs,
		"logsBloom":         receipt.Bloom,
		"type":              hexutil.Uint(0),
		"status":            hexutil.Uint(types.ReceiptStatusSuccessful),
	}
	if receipt.Logs == nil {
		fields["logs"] = []*types.Log{}
	}
	return fields
}
