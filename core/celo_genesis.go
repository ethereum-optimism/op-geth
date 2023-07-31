package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	contracts "github.com/ethereum/go-ethereum/contracts/celo"
	"github.com/ethereum/go-ethereum/crypto"
)

// Decode 0x prefixed hex string from file (including trailing newline)
func DecodeHex(hexbytes []byte) ([]byte, error) {
	// Strip 0x prefix and trailing newline
	hexbytes = hexbytes[2 : len(hexbytes)-1] // strip 0x prefix

	// Decode hex string
	bytes := make([]byte, hex.DecodedLen(len(hexbytes)))
	_, err := hex.Decode(bytes, hexbytes)
	if err != nil {
		return nil, fmt.Errorf("DecodeHex: %w", err)
	}

	return bytes, nil
}

var DevPrivateKey, _ = crypto.HexToECDSA("2771aff413cac48d9f8c114fabddd9195a2129f3c2c436caa07e27bb7f58ead5")
var DevAddr = common.BytesToAddress(DevAddr32.Bytes())
var DevAddr32 = common.HexToHash("0x42cf1bbc38BaAA3c4898ce8790e21eD2738c6A4a")

func celoGenesisAccounts() map[common.Address]GenesisAccount {
	// As defined in ERC-1967: Proxy Storage Slots (https://eips.ethereum.org/EIPS/eip-1967)
	var (
		proxy_owner_slot          = common.HexToHash("0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103")
		proxy_implementation_slot = common.HexToHash("0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc")
	)

	registryBytecode, err := DecodeHex(contracts.RegistryBytecodeRaw)
	if err != nil {
		panic(err)
	}
	goldTokenBytecode, err := DecodeHex(contracts.GoldTokenBytecodeRaw)
	if err != nil {
		panic(err)
	}
	proxyBytecode, err := DecodeHex(contracts.ProxyBytecodeRaw)
	if err != nil {
		panic(err)
	}
	devBalance, ok := new(big.Int).SetString("100000000000000000000", 10)
	if !ok {
		panic("Could not set devBalance!")
	}
	return map[common.Address]GenesisAccount{
		contracts.RegistryAddress: { // Registry Proxy
			Code: proxyBytecode,
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0"):   DevAddr32, // `_owner` slot in Registry contract
				proxy_implementation_slot: common.HexToHash("0xce11"),
				proxy_owner_slot:          DevAddr32,
			},
			Balance: big.NewInt(0),
		},
		common.HexToAddress("0xce11"): { // Registry Implementation
			Code:    registryBytecode,
			Balance: big.NewInt(0),
		},
		contracts.GoldTokenAddress: { // GoldToken Proxy
			Code: proxyBytecode,
			Storage: map[common.Hash]common.Hash{
				proxy_implementation_slot: common.HexToHash("0xce13"),
				proxy_owner_slot:          DevAddr32,
			},
			Balance: big.NewInt(0),
		},
		common.HexToAddress("0xce13"): { // GoldToken Implementation
			Code:    goldTokenBytecode,
			Balance: big.NewInt(0),
		},
		DevAddr: {
			Balance: devBalance,
		},
	}
}
