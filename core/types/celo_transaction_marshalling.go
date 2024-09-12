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
	"encoding/json"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

func celoTransactionMarshal(tx *Transaction) ([]byte, bool, error) {
	var enc txJSON
	// These are set for all tx types.
	enc.Hash = tx.Hash()
	enc.Type = hexutil.Uint64(tx.Type())
	switch itx := tx.inner.(type) {
	case *LegacyTx:
		if !itx.CeloLegacy {
			return nil, false, nil
		}
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		if tx.Protected() {
			enc.ChainID = (*hexutil.Big)(tx.ChainId())
		}
		// Celo specific fields
		enc.FeeCurrency = itx.FeeCurrency
		enc.GatewayFee = (*hexutil.Big)(itx.GatewayFee)
		enc.GatewayFeeRecipient = itx.GatewayFeeRecipient
		enc.EthCompatible = new(bool)
	case *CeloDynamicFeeTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		// Celo specific fields
		enc.FeeCurrency = itx.FeeCurrency
		enc.GatewayFee = (*hexutil.Big)(itx.GatewayFee)
		enc.GatewayFeeRecipient = itx.GatewayFeeRecipient
	case *CeloDynamicFeeTxV2:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		// Celo specific fields
		enc.FeeCurrency = itx.FeeCurrency
	case *CeloDenominatedTx:
		enc.ChainID = (*hexutil.Big)(itx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&itx.Nonce)
		enc.To = tx.To()
		enc.Gas = (*hexutil.Uint64)(&itx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(itx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(itx.GasTipCap)
		enc.Value = (*hexutil.Big)(itx.Value)
		enc.Input = (*hexutil.Bytes)(&itx.Data)
		enc.AccessList = &itx.AccessList
		enc.V = (*hexutil.Big)(itx.V)
		enc.R = (*hexutil.Big)(itx.R)
		enc.S = (*hexutil.Big)(itx.S)
		// Celo specific fields
		enc.FeeCurrency = itx.FeeCurrency
		enc.MaxFeeInFeeCurrency = (*hexutil.Big)(itx.MaxFeeInFeeCurrency)
	default:
		return nil, false, nil
	}
	bytes, err := json.Marshal(&enc)
	return bytes, true, err
}

func celoTransactionUnmarshal(dec txJSON, inner *TxData) (bool, error) {
	switch dec.Type {
	case LegacyTxType:
		// EthCompatible is only set to false for celo legacy transactions, otherwise its not set. So not set means it is ethCompatible.
		if dec.EthCompatible == nil {
			return false, nil
		}
		var itx LegacyTx
		*inner = &itx
		if dec.Nonce == nil {
			return true, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return true, errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.GasPrice == nil {
			return true, errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Value == nil {
			return true, errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return true, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input

		// signature R
		if dec.R == nil {
			return true, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		// signature S
		if dec.S == nil {
			return true, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		// signature V
		if dec.V == nil {
			return true, errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0 {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return true, err
			}
		}
		itx.CeloLegacy = true
		itx.FeeCurrency = dec.FeeCurrency
		itx.GatewayFeeRecipient = dec.GatewayFeeRecipient
		itx.GatewayFee = (*big.Int)(dec.GatewayFee)

	case CeloDynamicFeeTxType:
		var itx CeloDynamicFeeTx
		*inner = &itx
		if dec.ChainID == nil {
			return true, errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return true, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return true, errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return true, errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return true, errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Value == nil {
			return true, errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return true, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return true, errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return true, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return true, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return true, err
			}
		}
		// Celo specific fields
		itx.FeeCurrency = dec.FeeCurrency
		itx.GatewayFee = (*big.Int)(dec.GatewayFee)
		itx.GatewayFeeRecipient = dec.GatewayFeeRecipient
	case CeloDynamicFeeTxV2Type:
		var itx CeloDynamicFeeTxV2
		*inner = &itx
		if dec.ChainID == nil {
			return true, errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return true, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return true, errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return true, errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return true, errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Value == nil {
			return true, errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return true, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return true, errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return true, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return true, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return true, err
			}
		}
		// Celo specific fields
		itx.FeeCurrency = dec.FeeCurrency
	case CeloDenominatedTxType:
		var itx CeloDenominatedTx
		*inner = &itx
		if dec.ChainID == nil {
			return true, errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return true, errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Gas == nil {
			return true, errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.MaxPriorityFeePerGas == nil {
			return true, errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return true, errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Value == nil {
			return true, errors.New("missing required field 'value' in transaction")
		}
		itx.FeeCurrency = dec.FeeCurrency
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return true, errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return true, errors.New("missing required field 'v' in transaction")
		}
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return true, errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return true, errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return true, err
			}
		}
		// Celo specific fields
		itx.FeeCurrency = dec.FeeCurrency
		itx.MaxFeeInFeeCurrency = (*big.Int)(dec.MaxFeeInFeeCurrency)
	default:
		return false, nil
	}

	return true, nil
}

func celoDecodeTyped(b []byte) (TxData, bool, error) {
	var inner TxData
	switch b[0] {
	case CeloDenominatedTxType:
		inner = new(CeloDenominatedTx)
	case CeloDynamicFeeTxV2Type:
		inner = new(CeloDynamicFeeTxV2)
	case CeloDynamicFeeTxType:
		inner = new(CeloDynamicFeeTx)
	default:
		return nil, false, nil
	}
	err := inner.decode(b[1:])
	return inner, true, err
}
