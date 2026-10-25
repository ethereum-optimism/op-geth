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
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// Call is a single call dispatched by the protocol during EIP-8130 transaction
// execution. The wire form is rlp([to, data]); the dispatched call carries no
// value. Calls are grouped into phases on the transaction as [][]Call, which
// encodes as rlp([rlp([Call, ...]), ...]).
type Call struct {
	To   common.Address `json:"to"`
	Data hexutil.Bytes  `json:"data"`
}
