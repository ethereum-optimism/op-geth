package misc

import (
	"github.com/ethereum-optimism/superchain-registry/superchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
)

var create2DeployerAddress = common.HexToAddress("0x13b0D85CcB8bf860b6b79AF3029fCA081AE9beF2")
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
	if !c.IsOptimism() ||
		c.ChainID.Uint64() != params.BaseMainnetChainID ||
		c.CanyonTime == nil || *c.CanyonTime != timestamp ||
		db.GetCodeSize(create2DeployerAddress) > 0 {
		return
	}
	db.SetCode(create2DeployerAddress, create2DeployerCode)
}
