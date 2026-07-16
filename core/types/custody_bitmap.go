// Copyright 2026 The go-ethereum Authors
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

package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CustodyBitmap is the Engine API wire representation of the custody-column set.
type CustodyBitmap [16]byte

// MarshalText implements encoding.TextMarshaler.
func (b CustodyBitmap) MarshalText() ([]byte, error) {
	return []byte(hexutil.Encode(b[:])), nil
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *CustodyBitmap) UnmarshalText(input []byte) error {
	decoded, err := hexutil.Decode(string(input))
	if err != nil {
		return fmt.Errorf("custody bitmap: %v", err)
	}
	if len(decoded) != len(b) {
		return fmt.Errorf("custody bitmap: invalid length %d, want %d", len(decoded), len(b))
	}
	copy(b[:], decoded)
	return nil
}
