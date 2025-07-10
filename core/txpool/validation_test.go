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

package txpool

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"
)

func TestValidateTransactionMaxTxGasLimit(t *testing.T) {
	// Create a test private key and signer
	key, _ := crypto.GenerateKey()
	signer := types.NewEIP155Signer(big.NewInt(1))

	tests := []struct {
		name          string
		maxTxGasLimit uint64
		txGasLimit    uint64
		expectError   bool
		expectedError error
	}{
		{
			name:          "No limit set",
			maxTxGasLimit: 0,
			txGasLimit:    1000000,
			expectError:   false,
		},
		{
			name:          "Under limit",
			maxTxGasLimit: 100000,
			txGasLimit:    50000,
			expectError:   false,
		},
		{
			name:          "At limit",
			maxTxGasLimit: 100000,
			txGasLimit:    100000,
			expectError:   false,
		},
		{
			name:          "Over limit",
			maxTxGasLimit: 100000,
			txGasLimit:    150000,
			expectError:   true,
			expectedError: ErrTxGasLimitExceeded,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create test transaction with specified gas limit
			tx := types.NewTransaction(0, common.Address{}, big.NewInt(0), test.txGasLimit, big.NewInt(1000000000), nil)

			// Sign the transaction
			signedTx, err := types.SignTx(tx, signer, key)
			if err != nil {
				t.Fatalf("Failed to sign transaction: %v", err)
			}

			// Create minimal validation options
			opts := &ValidationOptions{
				Config:        params.TestChainConfig,
				Accept:        1 << types.LegacyTxType,
				MaxSize:       32 * 1024,
				MinTip:        big.NewInt(0),
				MaxTxGasLimit: test.maxTxGasLimit,
			}

			// Create test header with high gas limit to not interfere
			header := &types.Header{
				Number:     big.NewInt(1),
				GasLimit:   10000000,
				Time:       0,
				Difficulty: big.NewInt(0),
			}

			err = ValidateTransaction(signedTx, header, signer, opts)

			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				} else if !errors.Is(err, test.expectedError) {
					t.Errorf("Expected error %v, got %v", test.expectedError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}
