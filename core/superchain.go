package core

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/superchain"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

func LoadOPStackGenesis(chainID uint64) (*Genesis, error) {
	chain, err := superchain.GetChain(chainID)
	if err != nil {
		return nil, fmt.Errorf("error getting superchain: %w", err)
	}

	chConfig, err := chain.Config()
	if err != nil {
		return nil, fmt.Errorf("error getting chain config from superchain: %w", err)
	}

	cfg, err := params.LoadOPStackChainConfig(chConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to load params.ChainConfig for chain %d: %w", chainID, err)
	}

	gen, err := readSuperchainGenesis(chain)
	if err != nil {
		return nil, fmt.Errorf("failed to load genesis definition for chain %d: %w", chainID, err)
	}

	genesis := &Genesis{
		Config:        cfg,
		Nonce:         gen.Nonce,
		Timestamp:     gen.Timestamp,
		ExtraData:     gen.ExtraData,
		GasLimit:      gen.GasLimit,
		Difficulty:    gen.Difficulty,
		Mixhash:       gen.Mixhash,
		Coinbase:      gen.Coinbase,
		Alloc:         gen.Alloc,
		Number:        gen.Number,
		GasUsed:       gen.GasUsed,
		ParentHash:    gen.ParentHash,
		BaseFee:       gen.BaseFee,
		ExcessBlobGas: gen.ExcessBlobGas,
		BlobGasUsed:   gen.BlobGasUsed,
	}

	if gen.StateHash != nil {
		if len(gen.Alloc) > 0 {
			return nil, fmt.Errorf("chain definition unexpectedly contains both allocation (%d) and state-hash %s", len(gen.Alloc), *gen.StateHash)
		}
		genesis.StateHash = gen.StateHash
		genesis.Alloc = nil
	}

	genesisBlock := genesis.ToBlock()
	genesisBlockHash := genesisBlock.Hash()
	expectedHash := chConfig.Genesis.L2.Hash

	// Verify we correctly produced the genesis config by recomputing the genesis-block-hash,
	// and check the genesis matches the chain genesis definition.
	if chConfig.Genesis.L2.Number != genesisBlock.NumberU64() {
		switch chainID {
		case params.OPMainnetChainID:
			expectedHash = common.HexToHash("0x7ca38a1916c42007829c55e69d3e9a73265554b586a499015373241b8a3fa48b")
		default:
			return nil, fmt.Errorf("unknown stateless genesis definition for chain %d", chainID)
		}
	}
	if expectedHash != genesisBlockHash {
		return nil, fmt.Errorf("chainID=%d: produced genesis with hash %s but expected %s", chainID, genesisBlockHash, expectedHash)
	}
	return genesis, nil
}

func readSuperchainGenesis(chain *superchain.Chain) (*Genesis, error) {
	genData, err := chain.GenesisData()
	if err != nil {
		return nil, fmt.Errorf("error getting genesis data from superchain: %w", err)
	}
	gen := new(Genesis)
	if err := json.Unmarshal(genData, gen); err != nil {
		return nil, fmt.Errorf("failed to unmarshal genesis data: %w", err)
	}
	return gen, nil
}
