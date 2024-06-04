package miner

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

// cStables addresses on mainnet
var (
	cUSD_TOKEN  = common.HexToAddress("0x765DE816845861e75A25fCA122bb6898B8B1282a")
	cEUR_TOKEN  = common.HexToAddress("0xD8763CBa276a3738E6DE85b4b3bF5FDed6D6cA73")
	cREAL_TOKEN = common.HexToAddress("0xe8537a3d056DA446677B9E9d6c5dB704EaAb4787")
	USDC_TOKEN  = common.HexToAddress("0xcebA9300f2b948710d2653dD7B07f33A8B32118C")
	USDT_TOKEN  = common.HexToAddress("0x48065fbBE25f71C9282ddf5e1cD6D6A887483D5e")
)

// default limits default fraction
const DefaultFeeCurrencyLimit = 0.5

// default limits configuration
var DefaultFeeCurrencyLimits = map[uint64]map[common.Address]float64{
	params.CeloMainnetChainID: {
		cUSD_TOKEN:  0.9,
		USDT_TOKEN:  0.9,
		USDC_TOKEN:  0.9,
		cEUR_TOKEN:  0.5,
		cREAL_TOKEN: 0.5,
	},
}
