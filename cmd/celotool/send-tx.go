// Copyright 2024 The celo Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/urfave/cli/v2"
)

var (
	rpcUrlFlag = &cli.StringFlag{
		Name:        "rpc-url",
		Aliases:     []string{"r"},
		DefaultText: "http://localhost:8545",
		Usage:       "The rpc endpoint",
	}
	privKeyFlag = &cli.StringFlag{
		Name:     "private-key",
		Usage:    "Use the provided private key",
		Required: true,
	}
	valueFlag = &cli.Int64Flag{
		Name:        "value",
		DefaultText: "0",
		Usage:       "The value to send, in wei",
	}
)

var commandSend = &cli.Command{
	Name:      "send",
	Usage:     "send celo tx (cip-64)",
	ArgsUsage: "[to] [feeCurrency]",
	Description: `
Send a CIP-64 transaction.

- to: the address to send the transaction to, in hex format
- feeCurrency: the fee currency address, in hex format

Example:
$ celotool send --rpc-url $RPC_URL --private-key $PRIVATE_KEY $TO $FEECURRENCY --value 1
`,
	Flags: []cli.Flag{
		rpcUrlFlag,
		privKeyFlag,
		valueFlag,
	},
	Action: func(ctx *cli.Context) error {
		privKeyRaw := ctx.String("private-key")
		if len(privKeyRaw) >= 2 && privKeyRaw[0] == '0' && (privKeyRaw[1] == 'x' || privKeyRaw[1] == 'X') {
			privKeyRaw = privKeyRaw[2:]
		}
		privateKey, err := crypto.HexToECDSA(privKeyRaw)
		if err != nil {
			return err
		}

		to := ctx.Args().Get(0)
		if to == "" {
			fmt.Println("missing 'to' address")
			return nil
		}
		toAddress := common.HexToAddress(to)

		feeCurrency := ctx.Args().Get(1)
		if feeCurrency == "" {
			fmt.Println("missing 'feeCurrency' address")
			return nil
		}
		feeCurrencyAddress := common.HexToAddress(feeCurrency)

		value := big.NewInt(ctx.Int64("value"))

		rpcUrl := ctx.String("rpc-url")
		client, err := ethclient.Dial(rpcUrl)
		if err != nil {
			return err
		}

		chainId, err := client.ChainID(context.Background())
		if err != nil {
			return err
		}

		nonce, err := client.PendingNonceAt(context.Background(), crypto.PubkeyToAddress(privateKey.PublicKey))
		if err != nil {
			return err
		}

		feeCap, err := client.SuggestGasPriceForCurrency(context.Background(), &feeCurrencyAddress)
		if err != nil {
			return err
		}

		txdata := &types.CeloDynamicFeeTxV2{
			ChainID:     chainId,
			Nonce:       nonce,
			To:          &toAddress,
			Gas:         100_000,
			GasFeeCap:   feeCap,
			GasTipCap:   big.NewInt(2),
			FeeCurrency: &feeCurrencyAddress,
			Value:       value,
		}

		signer := types.LatestSignerForChainID(chainId)
		tx, err := types.SignNewTx(privateKey, signer, txdata)
		if err != nil {
			return err
		}

		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			return err
		}

		fmt.Printf("tx sent: %s\n", tx.Hash().Hex())

		return nil
	},
}
