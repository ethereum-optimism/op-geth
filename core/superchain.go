package core

import (
	"fmt"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

func LoadOPStackGenesis(chainID uint64) (*Genesis, error) {
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return nil, fmt.Errorf("unknown chain ID: %d", chainID)
	}

	cfg, err := params.LoadOPStackChainConfig(chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to load params.ChainConfig for chain %d: %w", chainID, err)
	}

	genesis := &Genesis{
		Config:     cfg,
		Nonce:      0,
		Timestamp:  chConfig.Genesis.L2Time,
		ExtraData:  []byte("BEDROCK"),
		GasLimit:   30_000_000,
		Difficulty: nil,
		Mixhash:    common.Hash{},
		Coinbase:   common.Address{},
		Alloc:      nil,
		Number:     0,
		GasUsed:    0,
		ParentHash: common.Hash{},
		BaseFee:    nil,
	}

	// TODO: load state allocations

	// TODO: exceptions for OP-Mainnet and OP-Goerli to handle pre-Bedrock history

	if chConfig.Genesis.ExtraData != nil {
		genesis.ExtraData = *chConfig.Genesis.ExtraData
		if len(genesis.ExtraData) > 32 {
			return nil, fmt.Errorf("chain must have 32 bytes or less extra-data in genesis, got %d", len(genesis.ExtraData))
		}
	}
	// TODO: apply all genesis block values

	// Verify we correctly produced the genesis config by recomputing the genesis-block-hash
	genesisBlock := genesis.ToBlock()
	genesisBlockHash := genesisBlock.Hash()
	if [32]byte(chConfig.Genesis.L2.Hash) != genesisBlockHash {
		return nil, fmt.Errorf("produced genesis with hash %s but expected %s", genesisBlockHash, chConfig.Genesis.L2.Hash)
	}
	return genesis, nil
}

func SystemConfigAddr(chainID uint64) (common.Address, error) {
	// TODO(proto): when we move to CREATE-2 proxy addresses
	// for SystemConfig contracts we can deterministically compute the system config addr,
	// and do not have to load it from the superchain configs.
	chConfig, ok := superchain.OPChains[chainID]
	if !ok {
		return common.Address{}, fmt.Errorf("unknown chain ID: %d", chainID)
	}
	return common.Address(chConfig.SystemConfigAddr), nil
}
