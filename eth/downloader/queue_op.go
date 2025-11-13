package downloader

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type OPStackChainConfig interface {
	IsOptimismEcotone(time uint64) bool
	IsOptimismIsthmus(time uint64) bool
	IsOptimismJovian(time uint64) bool
}

func (q *queue) opValidateBody(header *types.Header, txs []*types.Transaction, withdrawalsHash common.Hash) error {
	if q.opConfig == nil {
		return nil
	}

	if len(txs) == 0 {
		return fmt.Errorf("%w: no txs in optimism block", errInvalidBody)
	}
	if !txs[0].IsDepositTx() {
		return fmt.Errorf("%w: first tx in optimism block is not a deposit", errInvalidBody)
	}

	if q.opConfig.IsOptimismIsthmus(header.Time) && header.WithdrawalsHash != nil && withdrawalsHash != types.EmptyWithdrawalsHash {
		// If Isthmus, we expect an empty list of withdrawal operations,
		// but the WithdrawalsHash in the header is used for the withdrawals state storage-root.
		return fmt.Errorf("%w: nonempty Isthmus withdrawals hash", errInvalidBody)
	}

	if q.opConfig.IsOptimismEcotone(header.Time) && header.BlobGasUsed == nil {
		return fmt.Errorf("%w: nil blobGasUsed after Ecotone", errInvalidBody)
	}

	// Jovian changes the interpretation of the BlobGasUsed field.
	if q.opConfig.IsOptimismJovian(header.Time) && !txs[len(txs)-1].IsDepositTx() && *header.BlobGasUsed == 0 {
		return fmt.Errorf("%w: blobGasUsed is zero with at least one non-deposit tx", errInvalidBody)
	}

	return nil
}
