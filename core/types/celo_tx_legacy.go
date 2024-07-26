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
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

var ethCompatibleTxNumFields = 9

// ethCompatibleTxRlpList is used for RLP encoding/decoding of eth-compatible transactions.
// As such, it:
// (a) excludes the Celo-only fields,
// (b) doesn't need the Hash or EthCompatible fields, and
// (c) doesn't need the `json` or `gencodec` tags
type ethCompatibleTxRlpList struct {
	Nonce    uint64          // nonce of sender account
	GasPrice *big.Int        // wei per gas
	Gas      uint64          // gas limit
	To       *common.Address `rlp:"nil"` // nil means contract creation
	Value    *big.Int        // wei amount
	Data     []byte          // contract invocation input data
	V, R, S  *big.Int        // signature values
}

// celoTxRlpList is used for RLP encoding/decoding of celo transactions.
type celoTxRlpList struct {
	Nonce               uint64          // nonce of sender account
	GasPrice            *big.Int        // wei per gas
	Gas                 uint64          // gas limit
	FeeCurrency         *common.Address `rlp:"nil"` // nil means native currency
	GatewayFeeRecipient *common.Address `rlp:"nil"` // nil means no gateway fee is paid
	GatewayFee          *big.Int        `rlp:"nil"`
	To                  *common.Address `rlp:"nil"` // nil means contract creation
	Value               *big.Int        // wei amount
	Data                []byte          // contract invocation input data
	V, R, S             *big.Int        // signature values
}

func toEthCompatibleRlpList(tx LegacyTx) ethCompatibleTxRlpList {
	return ethCompatibleTxRlpList{
		Nonce:    tx.Nonce,
		GasPrice: tx.GasPrice,
		Gas:      tx.Gas,
		To:       tx.To,
		Value:    tx.Value,
		Data:     tx.Data,
		V:        tx.V,
		R:        tx.R,
		S:        tx.S,
	}
}

func toCeloRlpList(tx LegacyTx) celoTxRlpList {
	return celoTxRlpList{
		Nonce:    tx.Nonce,
		GasPrice: tx.GasPrice,
		Gas:      tx.Gas,
		To:       tx.To,
		Value:    tx.Value,
		Data:     tx.Data,
		V:        tx.V,
		R:        tx.R,
		S:        tx.S,

		// Celo specific fields
		FeeCurrency:         tx.FeeCurrency,
		GatewayFeeRecipient: tx.GatewayFeeRecipient,
		GatewayFee:          tx.GatewayFee,
	}
}

func setTxFromEthCompatibleRlpList(tx *LegacyTx, rlplist ethCompatibleTxRlpList) {
	tx.Nonce = rlplist.Nonce
	tx.GasPrice = rlplist.GasPrice
	tx.Gas = rlplist.Gas
	tx.To = rlplist.To
	tx.Value = rlplist.Value
	tx.Data = rlplist.Data
	tx.V = rlplist.V
	tx.R = rlplist.R
	tx.S = rlplist.S
	tx.Hash = nil // txdata.Hash is calculated and saved inside tx.Hash()

	// Celo specific fields
	tx.FeeCurrency = nil
	tx.GatewayFeeRecipient = nil
	tx.GatewayFee = nil
	tx.CeloLegacy = false
}

func setTxFromCeloRlpList(tx *LegacyTx, rlplist celoTxRlpList) {
	tx.Nonce = rlplist.Nonce
	tx.GasPrice = rlplist.GasPrice
	tx.Gas = rlplist.Gas
	tx.To = rlplist.To
	tx.Value = rlplist.Value
	tx.Data = rlplist.Data
	tx.V = rlplist.V
	tx.R = rlplist.R
	tx.S = rlplist.S
	tx.Hash = nil // txdata.Hash is calculated and saved inside tx.Hash()

	// Celo specific fields
	tx.FeeCurrency = rlplist.FeeCurrency
	tx.GatewayFeeRecipient = rlplist.GatewayFeeRecipient
	tx.GatewayFee = rlplist.GatewayFee
	tx.CeloLegacy = true
}

// EncodeRLP implements rlp.Encoder
func (tx *LegacyTx) EncodeRLP(w io.Writer) error {
	if tx.CeloLegacy {
		return rlp.Encode(w, toCeloRlpList(*tx))
	} else {
		return rlp.Encode(w, toEthCompatibleRlpList(*tx))
	}
}

// DecodeRLP implements rlp.Decoder
func (tx *LegacyTx) DecodeRLP(s *rlp.Stream) (err error) {
	_, size, _ := s.Kind()
	var raw rlp.RawValue
	err = s.Decode(&raw)
	if err != nil {
		return err
	}
	headerSize := len(raw) - int(size)
	numElems, err := rlp.CountValues(raw[headerSize:])
	if err != nil {
		return err
	}
	if numElems == ethCompatibleTxNumFields {
		rlpList := ethCompatibleTxRlpList{}
		err = rlp.DecodeBytes(raw, &rlpList)
		setTxFromEthCompatibleRlpList(tx, rlpList)
	} else {
		var rlpList celoTxRlpList
		err = rlp.DecodeBytes(raw, &rlpList)
		setTxFromCeloRlpList(tx, rlpList)
	}

	return err
}
