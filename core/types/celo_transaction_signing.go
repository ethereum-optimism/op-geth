// Copyright 2024 The Celo Authors
// This file is part of the celo library.
//
// The celo library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The celo library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the celo library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

var (
	ErrDeprecatedTxType = errors.New("deprecated transaction type")
)

// celoSigner acts as an overlay signer that handles celo specific signing
// functionality and hands off to an upstream signer for any other transaction
// types. Unlike the signers in the go-ethereum library, the celoSigner is
// configured with a list of forks that determine it's signing capabilities so
// there should not be a need to create any further signers to handle celo
// specific transaction types.
type celoSigner struct {
	upstreamSigner Signer
	chainID        *big.Int
	activatedForks forks
}

// makeCeloSigner creates a new celoSigner that is configured to handle all
// celo forks that are active at the given block time. If there are no active
// celo forks the upstream signer will be returned.
func makeCeloSigner(chainConfig *params.ChainConfig, blockTime uint64, upstreamSigner Signer) Signer {
	s := &celoSigner{
		chainID:        chainConfig.ChainID,
		upstreamSigner: upstreamSigner,
		activatedForks: celoForks.activeForks(blockTime, chainConfig),
	}

	// If there are no active celo forks, return the upstream signer
	if len(s.activatedForks) == 0 {
		return upstreamSigner
	}
	return s
}

// latestCeloSigner creates a new celoSigner that is configured to handle all
// celo forks for non celo transaction types it will delegate to the given
// upstream signer.
func latestCeloSigner(chainID *big.Int, upstreamSigner Signer) Signer {
	return &celoSigner{
		chainID:        chainID,
		upstreamSigner: upstreamSigner,
		activatedForks: celoForks,
	}
}

// Sender implements Signer.
func (c *celoSigner) Sender(tx *Transaction) (common.Address, error) {
	if funcs := c.findTxFuncs(tx); funcs != nil {
		return funcs.sender(tx, funcs.hash, c.ChainID())
	}
	return c.upstreamSigner.Sender(tx)
}

// SignatureValues implements Signer.
func (c *celoSigner) SignatureValues(tx *Transaction, sig []byte) (r *big.Int, s *big.Int, v *big.Int, err error) {
	if funcs := c.findTxFuncs(tx); funcs != nil {
		return funcs.signatureValues(tx, sig, c.ChainID())
	}
	return c.upstreamSigner.SignatureValues(tx, sig)
}

// Hash implements Signer.
func (c *celoSigner) Hash(tx *Transaction) common.Hash {
	if funcs := c.findTxFuncs(tx); funcs != nil {
		return funcs.hash(tx, c.ChainID())
	}
	return c.upstreamSigner.Hash(tx)
}

// findTxFuncs returns the txFuncs for the given tx if it is supported by one
// of the active forks. Note that this mechanism can be used to deprecate
// support for tx types by having forks return deprecatedTxFuncs for a tx type.
func (c *celoSigner) findTxFuncs(tx *Transaction) *txFuncs {
	return c.activatedForks.findTxFuncs(tx)
}

// ChainID implements Signer.
func (c *celoSigner) ChainID() *big.Int {
	return c.chainID
}

// Equal implements Signer.
func (c *celoSigner) Equal(s Signer) bool {
	// Normally signers just check to see if the chainID and type are equal,
	// because their logic is hardcoded to a specific fork. In our case we need
	// to also know that the two signers have matching latest forks.
	other, ok := s.(*celoSigner)
	return ok && c.ChainID() == other.ChainID() && c.latestFork().equal(other.latestFork())
}

func (c *celoSigner) latestFork() fork {
	return c.activatedForks[len(c.activatedForks)-1]
}
