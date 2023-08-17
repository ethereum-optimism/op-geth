package params

import (
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/common"
)

func init() {
	for id, ch := range superchain.OPChains {
		NetworkNames[fmt.Sprintf("%d", id)] = ch.Name
	}
}

func OPStackChainIDByName(name string) (uint64, error) {
	for id, ch := range superchain.OPChains {
		if ch.Chain+"-"+ch.Superchain == name {
			return id, nil
		}
	}
	return 0, fmt.Errorf("unknown chain %q", name)
}

func LoadOPStackChainConfig(chainID uint64) (*ChainConfig, error) {
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}
	superchainConfig, ok := superchain.Superchains[chConfig.Superchain]
	if !ok {
		return nil, fmt.Errorf("unknown superchain %q", chConfig.Superchain)
	}

	genesisActivation := uint64(0)
	out := &ChainConfig{
		ChainID:                       new(big.Int).SetUint64(chainID),
		HomesteadBlock:                common.Big0,
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   common.Big0,
		EIP155Block:                   common.Big0,
		EIP158Block:                   common.Big0,
		ByzantiumBlock:                common.Big0,
		ConstantinopleBlock:           common.Big0,
		PetersburgBlock:               common.Big0,
		IstanbulBlock:                 common.Big0,
		MuirGlacierBlock:              common.Big0,
		BerlinBlock:                   common.Big0,
		LondonBlock:                   common.Big0,
		ArrowGlacierBlock:             common.Big0,
		GrayGlacierBlock:              common.Big0,
		MergeNetsplitBlock:            common.Big0,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		BedrockBlock:                  common.Big0,
		RegolithTime:                  &genesisActivation,
		TerminalTotalDifficulty:       common.Big0,
		TerminalTotalDifficultyPassed: true,
		Ethash:                        nil,
		Clique:                        nil,
		Optimism: &OptimismConfig{
			EIP1559Elasticity:  6,
			EIP1559Denominator: 50,
		},
	}

	// note: no actual parameters are being loaded, yet.
	// Future superchain upgrades are loaded from the superchain chConfig and applied to the geth ChainConfig here.
	_ = superchainConfig.Config

	// special overrides for OP-Stack chains with pre-Regolith upgrade history
	switch chainID {
	case OPGoerliChainID:
		out.LondonBlock = big.NewInt(4061224)
		out.ArrowGlacierBlock = big.NewInt(4061224)
		out.GrayGlacierBlock = big.NewInt(4061224)
		out.MergeNetsplitBlock = big.NewInt(4061224)
		out.BedrockBlock = big.NewInt(4061224)
		out.RegolithTime = &OptimismGoerliRegolithTime
		out.Optimism.EIP1559Elasticity = 10
	case OPMainnetChainID:
		out.BerlinBlock = big.NewInt(3950000)
		out.LondonBlock = big.NewInt(105235063)
		out.ArrowGlacierBlock = big.NewInt(105235063)
		out.GrayGlacierBlock = big.NewInt(105235063)
		out.MergeNetsplitBlock = big.NewInt(105235063)
		out.BedrockBlock = big.NewInt(105235063)
	case BaseGoerliChainID:
		out.RegolithTime = &BaseGoerliRegolithTime
	}

	return out, nil
}
