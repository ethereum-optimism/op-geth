package misc

import (
	"github.com/ethereum-optimism/superchain-registry/superchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

// create2Deployer is already deployed to Base goerli at 0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2,
// so we deploy it to 0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1 for hardfork testing purposes
var create2DeployerAddresses = map[uint64]common.Address{
	params.BaseGoerliChainID:  common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF1"),
	params.BaseMainnetChainID: common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2"),
}
var create2DeployerCodeHash = common.HexToHash("0xb0550b5b431e30d38000efb7107aaa0ade03d48a7198a140edda9d27134468b2")
var create2DeployerCode []byte

func init() {
	code, err := superchain.LoadContractBytecode(superchain.Hash(create2DeployerCodeHash))
	if err != nil {
		panic(err)
	}
	create2DeployerCode = code
}

func EnsureCreate2Deployer(c *params.ChainConfig, timestamp uint64, db vm.StateDB) {
	if !c.IsOptimism() || c.CanyonTime == nil || *c.CanyonTime != timestamp {
		return
	}
	address, ok := create2DeployerAddresses[c.ChainID.Uint64()]
	if !ok || db.GetCodeSize(address) > 0 {
		return
	}
	db.SetCode(address, create2DeployerCode)
}
