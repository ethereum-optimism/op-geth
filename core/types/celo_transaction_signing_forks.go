package types

import "github.com/ethereum/go-ethereum/params"

var (
	// celoForks is the list of celo forks that are supported by the
	// celoSigner. This list is ordered with more recent forks appearing
	// earlier. It is assumed that if a more recent fork is active then all
	// previous forks are also active.
	celoForks = forks{&cel2{}, &celoLegacy{}}
)

type forks []fork

// activeForks returns the active forks for the given block time and chain config.
func (f forks) activeForks(blockTime uint64, config *params.ChainConfig) []fork {
	for i, fork := range f {
		if fork.active(blockTime, config) {
			return f[i:]
		}
	}
	return nil
}

// findTxFuncs returns the txFuncs for the given tx if there is a fork that supports it.
func (f forks) findTxFuncs(tx *Transaction) *txFuncs {
	for _, fork := range f {
		if funcs := fork.txFuncs(tx); funcs != nil {
			return funcs
		}
	}
	return nil
}

// fork contains functionality to determine if it is active for a given block
// time and chain config. It also acts as a container for functionality related
// to transactions enabled or deprecated in that fork.
type fork interface {
	// active returns true if the fork is active at the given block time.
	active(blockTime uint64, config *params.ChainConfig) bool
	// equal returns true if the given fork is the same underlying type as this fork.
	equal(fork) bool
	// txFuncs returns the txFuncs for the given tx if it is supported by the
	// fork. If a fork deprecates a tx type then this function should return
	// deprecatedTxFuncs for that tx type.
	txFuncs(tx *Transaction) *txFuncs
}

// Cel2 is the fork marking the transition point from an L1 to an L2.
// It deprecates CeloDynamicFeeTxType and LegacyTxTypes with CeloLegacy set to true.
type cel2 struct{}

func (c *cel2) active(blockTime uint64, config *params.ChainConfig) bool {
	return config.IsCel2(blockTime)
}

func (c *cel2) equal(other fork) bool {
	_, ok := other.(*cel2)
	return ok
}

func (c *cel2) txFuncs(tx *Transaction) *txFuncs {
	t := tx.Type()
	switch {
	case t == LegacyTxType && tx.IsCeloLegacy():
		return deprecatedTxFuncs
	case t == CeloDynamicFeeTxType:
		return deprecatedTxFuncs
	}
	return nil
}

// celoLegacy isn't actually a fork, but a placeholder for all historical celo
// related forks occurring on the celo L1. We don't need to construct the full
// signer chain from the celo legacy project because we won't support
// historical transaction execution, so we just need to be able to derive the
// senders for historical transactions and since we assume that the historical
// data is correct we just need one blanket signer that can cover all legacy
// celo transactions, before the L2 transition point.
type celoLegacy struct{}

func (c *celoLegacy) active(blockTime uint64, config *params.ChainConfig) bool {
	// The celo legacy fork is always active in a celo context
	return config.Cel2Time != nil
}

func (c *celoLegacy) equal(other fork) bool {
	_, ok := other.(*cel2)
	return ok
}

func (c *celoLegacy) txFuncs(tx *Transaction) *txFuncs {
	t := tx.Type()
	switch {
	case t == uint8(LegacyTxType) && tx.IsCeloLegacy():
		return celoLegacyTxFuncs
	case t == DynamicFeeTxType:
		// We handle the dynamic fee tx type here because we need to handle
		// migrated dynamic fee txs. These were enabeled in celo in the Espresso
		// hardfork, which doesn't have any analogue in op-geth. Even though
		// op-geth does enable support for dynamic fee txs in the London
		// hardfork (which we set to the cel2 block) that fork contains a lot of
		// changes that were not part of Espresso. So instead we ned to handle
		// DynamicFeeTxTypes here.
		return dynamicFeeTxFuncs
	case t == AccessListTxType:
		// Similar to the dynamic fee tx type, we need to handle the access list tx type that was also enabled by the
		// espresso hardfork.
		return accessListTxFuncs
	case t == CeloDynamicFeeTxV2Type:
		return celoDynamicFeeTxV2Funcs
	case t == CeloDynamicFeeTxType:
		return celoDynamicFeeTxFuncs
	}
	return nil
}
