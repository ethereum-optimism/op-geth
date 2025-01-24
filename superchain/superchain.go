package superchain

import (
	"fmt"
	"path"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/naoina/toml"
)

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
