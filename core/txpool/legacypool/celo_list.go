package legacypool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
)

func celoFilterWhitelisted(blockNumber *big.Int, list *list, all *lookup, fcv txpool.FeeCurrencyValidator) {
	removed := list.txs.Filter(func(tx *types.Transaction) bool {
		return txpool.FeeCurrencyTx(tx) && fcv.IsWhitelisted(tx.FeeCurrency(), blockNumber)
	})
	for _, tx := range removed {
		hash := tx.Hash()
		all.Remove(hash)
	}
}
