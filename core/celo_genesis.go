package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	contracts "github.com/ethereum/go-ethereum/contracts/celo"
	contracts_config "github.com/ethereum/go-ethereum/contracts/config"
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
	proxyBytecode, err := DecodeHex(contracts.ProxyBytecodeRaw)
	if err != nil {
		panic(err)
	}
	registry_owner := common.HexToHash("0x42cf1bbc38BaAA3c4898ce8790e21eD2738c6A4a")
	return map[common.Address]GenesisAccount{
		// Celo Contracts
		contracts_config.RegistrySmartContractAddress: { // Registry Proxy
			Code: proxyBytecode,
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0"):   registry_owner, // `_owner` slot in Registry contract
				proxy_implementation_slot: common.HexToHash("0xce11"),
				proxy_owner_slot:          registry_owner,
			},
			Balance: big.NewInt(0),
		},
		common.HexToAddress("0xce11"): { // Registry Implementation
			Code:    registryBytecode,
			Balance: big.NewInt(0),
		},
	}
}
