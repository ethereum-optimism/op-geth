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
	currencySet[getCurrencyKey(nil)] = struct{}{} // Always include native token to handle potential value transfers
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

func (c *list) totalCostVar(feeCurrency *common.Address) *uint256.Int {
	key := getCurrencyKey(feeCurrency)
	if tc, ok := c.totalCost[key]; ok {
		return tc
	}
	newTc := new(uint256.Int)
	c.totalCost[key] = newTc
	return newTc
}

func (c *list) TotalCostFor(feeCurrency *common.Address) *uint256.Int {
	if tc, ok := c.totalCost[getCurrencyKey(feeCurrency)]; ok {
		return new(uint256.Int).Set(tc)
	}
	return new(uint256.Int)
}

func (c *list) costCapFor(feeCurrency *common.Address) *uint256.Int {
	if tc, ok := c.costCap[getCurrencyKey(feeCurrency)]; ok {
		return tc
	}
	return new(uint256.Int)
}

func (c *list) updateCostCapFor(feeCurrency *common.Address, possibleCap *uint256.Int) {
	currentCap := c.costCapFor(feeCurrency)
	if possibleCap.Cmp(currentCap) > 0 {
		c.costCap[getCurrencyKey(feeCurrency)] = possibleCap
	}
}

func (c *list) costCapsLowerThan(costLimits map[common.Address]*uint256.Int) bool {
	for curr, cap := range c.costCap {
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

func (c *list) setCapsTo(caps map[common.Address]*uint256.Int) {
	c.costCap = make(map[common.Address]*uint256.Int)
	for curr, cap := range caps {
		if cap == nil || cap.IsZero() {
			c.costCap[curr] = new(uint256.Int)
		} else {
			c.costCap[curr] = new(uint256.Int).Set(cap)
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
