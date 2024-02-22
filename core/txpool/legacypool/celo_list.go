package legacypool

import (
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/holiman/uint256"
)

func (l *list) FilterWhitelisted(rates common.ExchangeRates) (types.Transactions, types.Transactions) {
	removed := l.txs.Filter(func(tx *types.Transaction) bool {
		return !common.IsCurrencyWhitelisted(rates, tx.FeeCurrency())
	})

	if len(removed) == 0 {
		return nil, nil
	}

	invalid := l.dropInvalidsAfterRemovalAndReheap(removed)
	l.subTotalCost(removed)
	l.subTotalCost(invalid)
	return removed, invalid
}

func (l *list) dropInvalidsAfterRemovalAndReheap(removed types.Transactions) types.Transactions {
	var invalids types.Transactions
	// If the list was strict, filter anything above the lowest nonce
	// Note that the 'invalid' txs have no intersection with the 'removed' txs
	if l.strict {
		lowest := uint64(math.MaxUint64)
		for _, tx := range removed {
			if nonce := tx.Nonce(); lowest > nonce {
				lowest = nonce
			}
		}
		invalids = l.txs.filter(func(tx *types.Transaction) bool { return tx.Nonce() > lowest })
	}
	l.txs.reheap()
	return invalids
}

func (l *list) FeeCurrencies() []common.Address {
	currencySet := make(map[common.Address]interface{})
	for _, tx := range l.txs.items {
		// native currency (nil) represented as Zero address
		currencySet[getCurrencyKey(tx.FeeCurrency())] = struct{}{}
	}
	currencies := make([]common.Address, 0, len(currencySet))
	for curr := range currencySet {
		currencies = append(currencies, curr)
	}
	return currencies
}

func getCurrencyKey(feeCurrency *common.Address) common.Address {
	if feeCurrency == nil {
		return common.ZeroAddress
	}
	return *feeCurrency
}

func (l *list) totalCostVar(feeCurrency *common.Address) *uint256.Int {
	key := getCurrencyKey(feeCurrency)
	if tc, ok := l.totalCost[key]; ok {
		return tc
	}
	newTc := new(uint256.Int)
	l.totalCost[key] = newTc
	return newTc
}

func (l *list) TotalCostFor(feeCurrency *common.Address) *uint256.Int {
	if tc, ok := l.totalCost[getCurrencyKey(feeCurrency)]; ok {
		return new(uint256.Int).Set(tc)
	}
	return new(uint256.Int)
}

func (l *list) costCapFor(feeCurrency *common.Address) *uint256.Int {
	if tc, ok := l.costCap[getCurrencyKey(feeCurrency)]; ok {
		return tc
	}
	return new(uint256.Int)
}

func (l *list) updateCostCapFor(feeCurrency *common.Address, possibleCap *uint256.Int) {
	currentCap := l.costCapFor(feeCurrency)
	if possibleCap.Cmp(currentCap) > 0 {
		l.costCap[getCurrencyKey(feeCurrency)] = possibleCap
	}
}

func (l *list) costCapsLowerThan(costLimits map[common.Address]*uint256.Int) bool {
	for curr, cap := range l.costCap {
		limit, ok := costLimits[curr]
		if !ok || limit == nil {
			// If there's no limit for the currency we can assume the limit is zero
			return cap.IsZero()
		}
		if cap.Cmp(limit) > 0 {
			return false
		}
	}
	return true
}

func (l *list) setCapsTo(caps map[common.Address]*uint256.Int) {
	l.costCap = make(map[common.Address]*uint256.Int)
	for curr, cap := range caps {
		if cap == nil || cap.IsZero() {
			l.costCap[curr] = new(uint256.Int)
		} else {
			l.costCap[curr] = new(uint256.Int).Set(cap)
		}
	}
}

// GetNativeBaseFee returns the base fee for this priceHeap
func (h *priceHeap) GetNativeBaseFee() *big.Int {
	if h.ratesAndFees == nil {
		return nil
	}
	return h.ratesAndFees.GetNativeBaseFee()
}
