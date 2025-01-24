package superchain

import (
	"fmt"
	"path"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
)

/*
name = "Mainnet"
protocol_versions_addr = "0x8062AbC286f5e7D9428a0Ccb9AbD71e50d93b935"
superchain_config_addr = "0x95703e0982140D16f8ebA6d158FccEde42f04a4C"
op_contracts_manager_proxy_addr = "0x18CeC91779995AD14c880e4095456B9147160790"

[hardforks]
canyon_time =  1704992401 # Thu 11 Jan 2024 17:00:01 UTC
delta_time =   1708560000 # Thu 22 Feb 2024 00:00:00 UTC
ecotone_time = 1710374401 # Thu 14 Mar 2024 00:00:01 UTC
fjord_time =   1720627201 # Wed 10 Jul 2024 16:00:01 UTC
granite_time = 1726070401 # Wed 11 Sep 2024 16:00:01 UTC
holocene_time = 1736445601 # Thu 09 Jan 2025 18:00:01 UTC

[l1]
  chain_id = 1
  public_rpc = "https://ethereum-rpc.publicnode.com"
  explorer = "https://etherscan.io"

*/

type Superchain struct {
	Name                        string         `toml:"name"`
	ProtocolVersionsAddr        common.Address `toml:"protocol_versions_addr"`
	SuperchainConfigAddr        common.Address `toml:"superchain_config_addr"`
	OpContractsManagerProxyAddr common.Address `toml:"op_contracts_manager_proxy_addr"`
	Hardforks                   HardforkConfig
	L1                          L1Config
}

type L1Config struct {
	ChainID   uint64 `toml:"chain_id"`
	PublicRPC string `toml:"public_rpc"`
	Explorer  string `toml:"explorer"`
}

var (
	superchainsByNetwork = map[string]Superchain{}
	mtx                  sync.Mutex
)

func GetSuperchain(network string) (Superchain, error) {
	mtx.Lock()
	defer mtx.Unlock()

	var sc Superchain
	if sc, ok := superchainsByNetwork[network]; ok {
		return sc, nil
	}

	zr, err := configDataReader.Open(path.Join("configs", network, "superchain.toml"))
	if err != nil {
		return sc, err
	}

	if err := toml.NewDecoder(zr).Decode(&sc); err != nil {
		return sc, fmt.Errorf("error decoding superchain config: %w", err)
	}

	superchainsByNetwork[network] = sc
	return sc, nil
}
