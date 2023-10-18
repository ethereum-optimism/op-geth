package legacypool

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
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

func balanceMinusL1Cost(feeCurrency *common.Address, balance *big.Int, l1Cost *big.Int,
	fvc txpool.FeeCurrencyValidator) *big.Int {
	// TODO: will need currency convertion from native (l1Cost) to feeCurrency.
	return nil
}

func celoFilterBalance(l1cost *big.Int, gasLimit uint64, list *list,
	fcv txpool.FeeCurrencyValidator) (types.Transactions, types.Transactions) {

	// TODO: needs to filter out txs by gas limit, and by balance-l1cost for txs
	// disregarding currency.

	// drops, invalids := list.Filter(balance, gasLimit)
	return nil, nil
}
