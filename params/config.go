// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package params

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params/forks"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	HoleskyGenesisHash = common.HexToHash("0xb5f7f912443c940f21fd611f12828d75b534364ed9e95ca4e307729a4661bde4")
	SepoliaGenesisHash = common.HexToHash("0x25a5cc106eea7138acab33231d7160d69cb777ee0c2c553fcddf5138993e6dd9")
	GoerliGenesisHash  = common.HexToHash("0xbf7e331f7f7c1dd2e05159666b3bf8bc7a8a3a9eb1d518969eab529dd9b88c1a")
)

const (
	OPMainnetChainID        = 10
	OPGoerliChainID         = 420
	BaseMainnetChainID      = 8453
	BaseGoerliChainID       = 84531
	baseSepoliaChainID      = 84532
	baseGoerliDevnetChainID = 11763071
	pgnSepoliaChainID       = 58008
	devnetChainID           = 997
	chaosnetChainID         = 888
)

// OP Stack chain config
var (
	// March 17, 2023 @ 7:00:00 pm UTC
	OptimismGoerliRegolithTime = uint64(1679079600)
	// May 4, 2023 @ 5:00:00 pm UTC
	BaseGoerliRegolithTime = uint64(1683219600)
	// Apr 21, 2023 @ 6:30:00 pm UTC
	baseGoerliDevnetRegolithTime = uint64(1682101800)
	// March 5, 2023 @ 2:48:00 am UTC
	devnetRegolithTime = uint64(1677984480)
	// August 16, 2023 @ 3:34:22 am UTC
	chaosnetRegolithTime = uint64(1692156862)
)

func newUint64(val uint64) *uint64 { return &val }

var (
	MainnetTerminalTotalDifficulty, _ = new(big.Int).SetString("58_750_000_000_000_000_000_000", 0)

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                big.NewInt(1_150_000),
		DAOForkBlock:                  big.NewInt(1_920_000),
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(2_463_000),
		EIP155Block:                   big.NewInt(2_675_000),
		EIP158Block:                   big.NewInt(2_675_000),
		ByzantiumBlock:                big.NewInt(4_370_000),
		ConstantinopleBlock:           big.NewInt(7_280_000),
		PetersburgBlock:               big.NewInt(7_280_000),
		IstanbulBlock:                 big.NewInt(9_069_000),
		MuirGlacierBlock:              big.NewInt(9_200_000),
		BerlinBlock:                   big.NewInt(12_244_000),
		LondonBlock:                   big.NewInt(12_965_000),
		ArrowGlacierBlock:             big.NewInt(13_773_000),
		GrayGlacierBlock:              big.NewInt(15_050_000),
		TerminalTotalDifficulty:       MainnetTerminalTotalDifficulty, // 58_750_000_000_000_000_000_000
		TerminalTotalDifficultyPassed: true,
		ShanghaiTime:                  newUint64(1681338455),
		CancunTime:                    newUint64(1710338135),
		Ethash:                        new(EthashConfig),
	}
	// HoleskyChainConfig contains the chain parameters to run a node on the Holesky test network.
	HoleskyChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(17000),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              nil,
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             nil,
		GrayGlacierBlock:              nil,
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  newUint64(1696000704),
		CancunTime:                    newUint64(1707305664),
		Ethash:                        new(EthashConfig),
	}
	// SepoliaChainConfig contains the chain parameters to run a node on the Sepolia test network.
	SepoliaChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(11155111),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             nil,
		GrayGlacierBlock:              nil,
		TerminalTotalDifficulty:       big.NewInt(17_000_000_000_000_000),
		TerminalTotalDifficultyPassed: true,
		MergeNetsplitBlock:            big.NewInt(1735371),
		ShanghaiTime:                  newUint64(1677557088),
		CancunTime:                    newUint64(1706655072),
		Ethash:                        new(EthashConfig),
	}
	// GoerliChainConfig contains the chain parameters to run a node on the Görli test network.
	GoerliChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(5),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                true,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(1_561_651),
		MuirGlacierBlock:              nil,
		BerlinBlock:                   big.NewInt(4_460_644),
		LondonBlock:                   big.NewInt(5_062_605),
		ArrowGlacierBlock:             nil,
		TerminalTotalDifficulty:       big.NewInt(10_790_000),
		TerminalTotalDifficultyPassed: true,
		ShanghaiTime:                  newUint64(1678832736),
		CancunTime:                    newUint64(1705473120),
		Clique: &CliqueConfig{
			Period: 15,
			Epoch:  30000,
		},
	}
	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	AllEthashProtocolChanges = &ChainConfig{
		ChainID:                       big.NewInt(1337),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		VerkleTime:                    nil,
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: true,
		Ethash:                        new(EthashConfig),
		Clique:                        nil,
	}

	AllDevChainProtocolChanges = &ChainConfig{
		ChainID:                       big.NewInt(1337),
		HomesteadBlock:                big.NewInt(0),
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		ShanghaiTime:                  newUint64(0),
		CancunTime:                    newUint64(0),
		Cel2Time:                      newUint64(0),
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
	}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	AllCliqueProtocolChanges = &ChainConfig{
		ChainID:                       big.NewInt(1337),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             nil,
		GrayGlacierBlock:              nil,
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		VerkleTime:                    nil,
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: false,
		Ethash:                        nil,
		Clique:                        &CliqueConfig{Period: 0, Epoch: 30000},
	}

	// TestChainConfig contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers for testing purposes.
	TestChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		VerkleTime:                    nil,
		Cel2Time:                      newUint64(0),
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: false,
		Ethash:                        new(EthashConfig),
		Clique:                        nil,
	}

	// MergedTestChainConfig contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers for testing purposes.
	MergedTestChainConfig = &ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                big.NewInt(0),
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   big.NewInt(0),
		EIP155Block:                   big.NewInt(0),
		EIP158Block:                   big.NewInt(0),
		ByzantiumBlock:                big.NewInt(0),
		ConstantinopleBlock:           big.NewInt(0),
		PetersburgBlock:               big.NewInt(0),
		IstanbulBlock:                 big.NewInt(0),
		MuirGlacierBlock:              big.NewInt(0),
		BerlinBlock:                   big.NewInt(0),
		LondonBlock:                   big.NewInt(0),
		ArrowGlacierBlock:             big.NewInt(0),
		GrayGlacierBlock:              big.NewInt(0),
		MergeNetsplitBlock:            big.NewInt(0),
		ShanghaiTime:                  newUint64(0),
		CancunTime:                    newUint64(0),
		PragueTime:                    nil,
		VerkleTime:                    nil,
		TerminalTotalDifficulty:       big.NewInt(0),
		TerminalTotalDifficultyPassed: true,
		Ethash:                        new(EthashConfig),
		Clique:                        nil,
	}

	// NonActivatedConfig defines the chain configuration without activating
	// any protocol change (EIPs).
	NonActivatedConfig = &ChainConfig{
		ChainID:                       big.NewInt(1),
		HomesteadBlock:                nil,
		DAOForkBlock:                  nil,
		DAOForkSupport:                false,
		EIP150Block:                   nil,
		EIP155Block:                   nil,
		EIP158Block:                   nil,
		ByzantiumBlock:                nil,
		ConstantinopleBlock:           nil,
		PetersburgBlock:               nil,
		IstanbulBlock:                 nil,
		MuirGlacierBlock:              nil,
		BerlinBlock:                   nil,
		LondonBlock:                   nil,
		ArrowGlacierBlock:             nil,
		GrayGlacierBlock:              nil,
		MergeNetsplitBlock:            nil,
		ShanghaiTime:                  nil,
		CancunTime:                    nil,
		PragueTime:                    nil,
		VerkleTime:                    nil,
		TerminalTotalDifficulty:       nil,
		TerminalTotalDifficultyPassed: false,
		Ethash:                        new(EthashConfig),
		Clique:                        nil,
	}
	TestRules = TestChainConfig.Rules(new(big.Int), false, 0)

	// This is an Optimism chain config with bedrock starting a block 5, introduced for historical endpoint testing, largely based on the clique config
	OptimismTestConfig = func() *ChainConfig {
		conf := *AllCliqueProtocolChanges // copy the config
		conf.Clique = nil
		conf.TerminalTotalDifficultyPassed = true
		conf.BedrockBlock = big.NewInt(5)
		conf.Optimism = &OptimismConfig{EIP1559Elasticity: 50, EIP1559Denominator: 10}
		return &conf
	}()
)

// NetworkNames are user friendly names to use in the chain spec banner.
var NetworkNames = map[string]string{
	MainnetChainConfig.ChainID.String(): "mainnet",
	GoerliChainConfig.ChainID.String():  "goerli",
	SepoliaChainConfig.ChainID.String(): "sepolia",
	HoleskyChainConfig.ChainID.String(): "holesky",
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

	HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

	DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
	DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

	// EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
	EIP150Block *big.Int `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

	ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
	ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
	PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
	IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
	MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)
	BerlinBlock         *big.Int `json:"berlinBlock,omitempty"`         // Berlin switch block (nil = no fork, 0 = already on berlin)
	LondonBlock         *big.Int `json:"londonBlock,omitempty"`         // London switch block (nil = no fork, 0 = already on london)
	ArrowGlacierBlock   *big.Int `json:"arrowGlacierBlock,omitempty"`   // Eip-4345 (bomb delay) switch block (nil = no fork, 0 = already activated)
	GrayGlacierBlock    *big.Int `json:"grayGlacierBlock,omitempty"`    // Eip-5133 (bomb delay) switch block (nil = no fork, 0 = already activated)
	MergeNetsplitBlock  *big.Int `json:"mergeNetsplitBlock,omitempty"`  // Virtual fork after The Merge to use as a network splitter

	// Fork scheduling was switched from blocks to timestamps here

	ShanghaiTime *uint64 `json:"shanghaiTime,omitempty"` // Shanghai switch time (nil = no fork, 0 = already on shanghai)
	CancunTime   *uint64 `json:"cancunTime,omitempty"`   // Cancun switch time (nil = no fork, 0 = already on cancun)
	PragueTime   *uint64 `json:"pragueTime,omitempty"`   // Prague switch time (nil = no fork, 0 = already on prague)
	VerkleTime   *uint64 `json:"verkleTime,omitempty"`   // Verkle switch time (nil = no fork, 0 = already on verkle)

	BedrockBlock *big.Int `json:"bedrockBlock,omitempty"` // Bedrock switch block (nil = no fork, 0 = already on optimism bedrock)
	RegolithTime *uint64  `json:"regolithTime,omitempty"` // Regolith switch time (nil = no fork, 0 = already on optimism regolith)
	CanyonTime   *uint64  `json:"canyonTime,omitempty"`   // Canyon switch time (nil = no fork, 0 = already on optimism canyon)
	// Delta: the Delta upgrade does not affect the execution-layer, and is thus not configurable in the chain config.
	EcotoneTime  *uint64 `json:"ecotoneTime,omitempty"`  // Ecotone switch time (nil = no fork, 0 = already on optimism ecotone)
	FjordTime    *uint64 `json:"fjordTime,omitempty"`    // Fjord switch time (nil = no fork, 0 = already on Optimism Fjord)
	GraniteTime  *uint64 `json:"graniteTime,omitempty"`  // Granite switch time (nil = no fork, 0 = already on Optimism Granite)
	HoloceneTime *uint64 `json:"holoceneTime,omitempty"` // Holocene switch time (nil = no fork, 0 = already on Optimism Holocene)

	InteropTime *uint64 `json:"interopTime,omitempty"` // Interop switch time (nil = no fork, 0 = already on optimism interop)

	Cel2Time *uint64 `json:"cel2Time,omitempty"` // Cel2 switch time (nil = no fork, 0 = already on optimism cel2)

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	// TerminalTotalDifficultyPassed is a flag specifying that the network already
	// passed the terminal total difficulty. Its purpose is to disable legacy sync
	// even without having seen the TTD locally (safer long term).
	//
	// TODO(karalabe): Drop this field eventually (always assuming PoS mode)
	TerminalTotalDifficultyPassed bool `json:"terminalTotalDifficultyPassed,omitempty"`

	// Various consensus engines
	Ethash *EthashConfig `json:"ethash,omitempty"`
	Clique *CliqueConfig `json:"clique,omitempty"`

	// Optimism config, nil if not active
	Optimism *OptimismConfig `json:"optimism,omitempty"`
}

// EthashConfig is the consensus engine configs for proof-of-work based sealing.
type EthashConfig struct{}

// String implements the stringer interface, returning the consensus engine details.
func (c EthashConfig) String() string {
	return "ethash"
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c CliqueConfig) String() string {
	return fmt.Sprintf("clique(period: %d, epoch: %d)", c.Period, c.Epoch)
}

// OptimismConfig is the optimism config.
type OptimismConfig struct {
	EIP1559Elasticity        uint64  `json:"eip1559Elasticity"`
	EIP1559Denominator       uint64  `json:"eip1559Denominator"`
	EIP1559DenominatorCanyon *uint64 `json:"eip1559DenominatorCanyon,omitempty"`
}

// String implements the stringer interface, returning the optimism fee config details.
func (o *OptimismConfig) String() string {
	return "optimism"
}

// Description returns a human-readable description of ChainConfig.
func (c *ChainConfig) Description() string {
	var banner string

	// Create some basic network config output
	network := NetworkNames[c.ChainID.String()]
	if network == "" {
		network = "unknown"
	}
	banner += fmt.Sprintf("Chain ID:  %v (%s)\n", c.ChainID, network)
	switch {
	case c.Optimism != nil:
		banner += "Consensus: Optimism\n"
	case c.Ethash != nil:
		if c.TerminalTotalDifficulty == nil {
			banner += "Consensus: Ethash (proof-of-work)\n"
		} else if !c.TerminalTotalDifficultyPassed {
			banner += "Consensus: Beacon (proof-of-stake), merging from Ethash (proof-of-work)\n"
		} else {
			banner += "Consensus: Beacon (proof-of-stake), merged from Ethash (proof-of-work)\n"
		}
	case c.Clique != nil:
		if c.TerminalTotalDifficulty == nil {
			banner += "Consensus: Clique (proof-of-authority)\n"
		} else if !c.TerminalTotalDifficultyPassed {
			banner += "Consensus: Beacon (proof-of-stake), merging from Clique (proof-of-authority)\n"
		} else {
			banner += "Consensus: Beacon (proof-of-stake), merged from Clique (proof-of-authority)\n"
		}
	default:
		banner += "Consensus: unknown\n"
	}
	banner += "\n"

	// Create a list of forks with a short description of them. Forks that only
	// makes sense for mainnet should be optional at printing to avoid bloating
	// the output for testnets and private networks.
	banner += "Pre-Merge hard forks (block based):\n"
	banner += fmt.Sprintf(" - Homestead:                   #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/homestead.md)\n", c.HomesteadBlock)
	if c.DAOForkBlock != nil {
		banner += fmt.Sprintf(" - DAO Fork:                    #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/dao-fork.md)\n", c.DAOForkBlock)
	}
	banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/tangerine-whistle.md)\n", c.EIP150Block)
	banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/spurious-dragon.md)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Byzantium:                   #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/byzantium.md)\n", c.ByzantiumBlock)
	banner += fmt.Sprintf(" - Constantinople:              #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/constantinople.md)\n", c.ConstantinopleBlock)
	banner += fmt.Sprintf(" - Petersburg:                  #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/petersburg.md)\n", c.PetersburgBlock)
	banner += fmt.Sprintf(" - Istanbul:                    #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/istanbul.md)\n", c.IstanbulBlock)
	if c.MuirGlacierBlock != nil {
		banner += fmt.Sprintf(" - Muir Glacier:                #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/muir-glacier.md)\n", c.MuirGlacierBlock)
	}
	banner += fmt.Sprintf(" - Berlin:                      #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/berlin.md)\n", c.BerlinBlock)
	banner += fmt.Sprintf(" - London:                      #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/london.md)\n", c.LondonBlock)
	if c.ArrowGlacierBlock != nil {
		banner += fmt.Sprintf(" - Arrow Glacier:               #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/arrow-glacier.md)\n", c.ArrowGlacierBlock)
	}
	if c.GrayGlacierBlock != nil {
		banner += fmt.Sprintf(" - Gray Glacier:                #%-8v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/gray-glacier.md)\n", c.GrayGlacierBlock)
	}
	banner += "\n"

	// Add a special section for the merge as it's non-obvious
	if c.TerminalTotalDifficulty == nil {
		banner += "The Merge is not yet available for this network!\n"
		banner += " - Hard-fork specification: https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/paris.md\n"
	} else {
		banner += "Merge configured:\n"
		banner += " - Hard-fork specification:    https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/paris.md\n"
		banner += fmt.Sprintf(" - Network known to be merged: %v\n", c.TerminalTotalDifficultyPassed)
		banner += fmt.Sprintf(" - Total terminal difficulty:  %v\n", c.TerminalTotalDifficulty)
		if c.MergeNetsplitBlock != nil {
			banner += fmt.Sprintf(" - Merge netsplit block:       #%-8v\n", c.MergeNetsplitBlock)
		}
	}
	banner += "\n"

	// Create a list of forks post-merge
	banner += "Post-Merge hard forks (timestamp based):\n"
	if c.ShanghaiTime != nil {
		banner += fmt.Sprintf(" - Shanghai:                    @%-10v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/shanghai.md)\n", *c.ShanghaiTime)
	}
	if c.CancunTime != nil {
		banner += fmt.Sprintf(" - Cancun:                      @%-10v (https://github.com/ethereum/execution-specs/blob/master/network-upgrades/mainnet-upgrades/cancun.md)\n", *c.CancunTime)
	}
	if c.PragueTime != nil {
		banner += fmt.Sprintf(" - Prague:                      @%-10v\n", *c.PragueTime)
	}
	if c.VerkleTime != nil {
		banner += fmt.Sprintf(" - Verkle:                      @%-10v\n", *c.VerkleTime)
	}
	if c.RegolithTime != nil {
		banner += fmt.Sprintf(" - Regolith:                    @%-10v\n", *c.RegolithTime)
	}
	if c.CanyonTime != nil {
		banner += fmt.Sprintf(" - Canyon:                      @%-10v\n", *c.CanyonTime)
	}
	if c.EcotoneTime != nil {
		banner += fmt.Sprintf(" - Ecotone:                     @%-10v\n", *c.EcotoneTime)
	}
	if c.FjordTime != nil {
		banner += fmt.Sprintf(" - Fjord:                       @%-10v\n", *c.FjordTime)
	}
	if c.GraniteTime != nil {
		banner += fmt.Sprintf(" - Granite:                     @%-10v\n", *c.GraniteTime)
	}
	if c.HoloceneTime != nil {
		banner += fmt.Sprintf(" - Holocene:                     @%-10v\n", *c.HoloceneTime)
	}
	if c.InteropTime != nil {
		banner += fmt.Sprintf(" - Interop:                     @%-10v\n", *c.InteropTime)
	}
	if c.Cel2Time != nil {
		banner += fmt.Sprintf(" - Cel2:                        @%-10v\n", *c.Cel2Time)
	}
	return banner
}

// IsHomestead returns whether num is either equal to the homestead block or greater.
func (c *ChainConfig) IsHomestead(num *big.Int) bool {
	return isBlockForked(c.HomesteadBlock, num)
}

// IsDAOFork returns whether num is either equal to the DAO fork block or greater.
func (c *ChainConfig) IsDAOFork(num *big.Int) bool {
	return isBlockForked(c.DAOForkBlock, num)
}

// IsEIP150 returns whether num is either equal to the EIP150 fork block or greater.
func (c *ChainConfig) IsEIP150(num *big.Int) bool {
	return isBlockForked(c.EIP150Block, num)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	return isBlockForked(c.EIP155Block, num)
}

// IsEIP158 returns whether num is either equal to the EIP158 fork block or greater.
func (c *ChainConfig) IsEIP158(num *big.Int) bool {
	return isBlockForked(c.EIP158Block, num)
}

// IsByzantium returns whether num is either equal to the Byzantium fork block or greater.
func (c *ChainConfig) IsByzantium(num *big.Int) bool {
	return isBlockForked(c.ByzantiumBlock, num)
}

// IsConstantinople returns whether num is either equal to the Constantinople fork block or greater.
func (c *ChainConfig) IsConstantinople(num *big.Int) bool {
	return isBlockForked(c.ConstantinopleBlock, num)
}

// IsMuirGlacier returns whether num is either equal to the Muir Glacier (EIP-2384) fork block or greater.
func (c *ChainConfig) IsMuirGlacier(num *big.Int) bool {
	return isBlockForked(c.MuirGlacierBlock, num)
}

// IsPetersburg returns whether num is either
// - equal to or greater than the PetersburgBlock fork block,
// - OR is nil, and Constantinople is active
func (c *ChainConfig) IsPetersburg(num *big.Int) bool {
	return isBlockForked(c.PetersburgBlock, num) || c.PetersburgBlock == nil && isBlockForked(c.ConstantinopleBlock, num)
}

// IsIstanbul returns whether num is either equal to the Istanbul fork block or greater.
func (c *ChainConfig) IsIstanbul(num *big.Int) bool {
	return isBlockForked(c.IstanbulBlock, num)
}

// IsBerlin returns whether num is either equal to the Berlin fork block or greater.
func (c *ChainConfig) IsBerlin(num *big.Int) bool {
	return isBlockForked(c.BerlinBlock, num)
}

// IsLondon returns whether num is either equal to the London fork block or greater.
func (c *ChainConfig) IsLondon(num *big.Int) bool {
	return isBlockForked(c.LondonBlock, num)
}

// IsArrowGlacier returns whether num is either equal to the Arrow Glacier (EIP-4345) fork block or greater.
func (c *ChainConfig) IsArrowGlacier(num *big.Int) bool {
	return isBlockForked(c.ArrowGlacierBlock, num)
}

// IsGrayGlacier returns whether num is either equal to the Gray Glacier (EIP-5133) fork block or greater.
func (c *ChainConfig) IsGrayGlacier(num *big.Int) bool {
	return isBlockForked(c.GrayGlacierBlock, num)
}

// IsTerminalPoWBlock returns whether the given block is the last block of PoW stage.
func (c *ChainConfig) IsTerminalPoWBlock(parentTotalDiff *big.Int, totalDiff *big.Int) bool {
	if c.TerminalTotalDifficulty == nil {
		return false
	}
	return parentTotalDiff.Cmp(c.TerminalTotalDifficulty) < 0 && totalDiff.Cmp(c.TerminalTotalDifficulty) >= 0
}

// IsShanghai returns whether time is either equal to the Shanghai fork time or greater.
func (c *ChainConfig) IsShanghai(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.ShanghaiTime, time)
}

// IsCancun returns whether time is either equal to the Cancun fork time or greater.
func (c *ChainConfig) IsCancun(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.CancunTime, time)
}

// IsPrague returns whether time is either equal to the Prague fork time or greater.
func (c *ChainConfig) IsPrague(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.PragueTime, time)
}

// IsVerkle returns whether time is either equal to the Verkle fork time or greater.
func (c *ChainConfig) IsVerkle(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.VerkleTime, time)
}

// IsEIP4762 returns whether eip 4762 has been activated at given block.
func (c *ChainConfig) IsEIP4762(num *big.Int, time uint64) bool {
	return c.IsVerkle(num, time)
}

// IsBedrock returns whether num is either equal to the Bedrock fork block or greater.
func (c *ChainConfig) IsBedrock(num *big.Int) bool {
	return isBlockForked(c.BedrockBlock, num)
}

func (c *ChainConfig) IsRegolith(time uint64) bool {
	return isTimestampForked(c.RegolithTime, time)
}

func (c *ChainConfig) IsCanyon(time uint64) bool {
	return isTimestampForked(c.CanyonTime, time)
}

func (c *ChainConfig) IsEcotone(time uint64) bool {
	return isTimestampForked(c.EcotoneTime, time)
}

func (c *ChainConfig) IsFjord(time uint64) bool {
	return isTimestampForked(c.FjordTime, time)
}

func (c *ChainConfig) IsGranite(time uint64) bool {
	return isTimestampForked(c.GraniteTime, time)
}

func (c *ChainConfig) IsHolocene(time uint64) bool {
	return isTimestampForked(c.HoloceneTime, time)
}

func (c *ChainConfig) IsInterop(time uint64) bool {
	return isTimestampForked(c.InteropTime, time)
}

func (c *ChainConfig) IsCel2(time uint64) bool {
	return isTimestampForked(c.Cel2Time, time)
}

// IsOptimism returns whether the node is an optimism node or not.
func (c *ChainConfig) IsOptimism() bool {
	return c.Optimism != nil
}

// IsOptimismBedrock returns true iff this is an optimism node & bedrock is active
func (c *ChainConfig) IsOptimismBedrock(num *big.Int) bool {
	return c.IsOptimism() && c.IsBedrock(num)
}

func (c *ChainConfig) IsOptimismRegolith(time uint64) bool {
	return c.IsOptimism() && c.IsRegolith(time)
}

func (c *ChainConfig) IsOptimismCanyon(time uint64) bool {
	return c.IsOptimism() && c.IsCanyon(time)
}

func (c *ChainConfig) IsOptimismEcotone(time uint64) bool {
	return c.IsOptimism() && c.IsEcotone(time)
}

func (c *ChainConfig) IsOptimismFjord(time uint64) bool {
	return c.IsOptimism() && c.IsFjord(time)
}

func (c *ChainConfig) IsOptimismGranite(time uint64) bool {
	return c.IsOptimism() && c.IsGranite(time)
}

func (c *ChainConfig) IsOptimismHolocene(time uint64) bool {
	return c.IsOptimism() && c.IsHolocene(time)
}

// IsOptimismPreBedrock returns true iff this is an optimism node & bedrock is not yet active
func (c *ChainConfig) IsOptimismPreBedrock(num *big.Int) bool {
	return c.IsOptimism() && !c.IsBedrock(num)
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height, time uint64, genesisTimestamp *uint64) *ConfigCompatError {
	var (
		bhead = new(big.Int).SetUint64(height)
		btime = time
	)
	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead, btime, genesisTimestamp)
		if err == nil || (lasterr != nil && err.RewindToBlock == lasterr.RewindToBlock && err.RewindToTime == lasterr.RewindToTime) {
			break
		}
		lasterr = err

		if err.RewindToTime > 0 {
			btime = err.RewindToTime
		} else {
			bhead.SetUint64(err.RewindToBlock)
		}
	}
	return lasterr
}

// CheckConfigForkOrder checks that we don't "skip" any forks, geth isn't pluggable enough
// to guarantee that forks can be implemented in a different order than on official networks
func (c *ChainConfig) CheckConfigForkOrder() error {
	type fork struct {
		name      string
		block     *big.Int // forks up to - and including the merge - were defined with block numbers
		timestamp *uint64  // forks after the merge are scheduled using timestamps
		optional  bool     // if true, the fork may be nil and next fork is still allowed
	}
	var lastFork fork
	for _, cur := range []fork{
		{name: "homesteadBlock", block: c.HomesteadBlock},
		{name: "daoForkBlock", block: c.DAOForkBlock, optional: true},
		{name: "eip150Block", block: c.EIP150Block},
		{name: "eip155Block", block: c.EIP155Block},
		{name: "eip158Block", block: c.EIP158Block},
		{name: "byzantiumBlock", block: c.ByzantiumBlock},
		{name: "constantinopleBlock", block: c.ConstantinopleBlock},
		{name: "petersburgBlock", block: c.PetersburgBlock},
		{name: "istanbulBlock", block: c.IstanbulBlock},
		{name: "muirGlacierBlock", block: c.MuirGlacierBlock, optional: true},
		{name: "berlinBlock", block: c.BerlinBlock},
		{name: "londonBlock", block: c.LondonBlock},
		{name: "arrowGlacierBlock", block: c.ArrowGlacierBlock, optional: true},
		{name: "grayGlacierBlock", block: c.GrayGlacierBlock, optional: true},
		{name: "mergeNetsplitBlock", block: c.MergeNetsplitBlock, optional: true},
		{name: "shanghaiTime", timestamp: c.ShanghaiTime},
		{name: "cancunTime", timestamp: c.CancunTime, optional: true},
		{name: "pragueTime", timestamp: c.PragueTime, optional: true},
		{name: "verkleTime", timestamp: c.VerkleTime, optional: true},
	} {
		if lastFork.name != "" {
			switch {
			// Non-optional forks must all be present in the chain config up to the last defined fork
			case lastFork.block == nil && lastFork.timestamp == nil && (cur.block != nil || cur.timestamp != nil):
				if cur.block != nil {
					return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at block %v",
						lastFork.name, cur.name, cur.block)
				} else {
					return fmt.Errorf("unsupported fork ordering: %v not enabled, but %v enabled at timestamp %v",
						lastFork.name, cur.name, *cur.timestamp)
				}

			// Fork (whether defined by block or timestamp) must follow the fork definition sequence
			case (lastFork.block != nil && cur.block != nil) || (lastFork.timestamp != nil && cur.timestamp != nil):
				if lastFork.block != nil && lastFork.block.Cmp(cur.block) > 0 {
					return fmt.Errorf("unsupported fork ordering: %v enabled at block %v, but %v enabled at block %v",
						lastFork.name, lastFork.block, cur.name, cur.block)
				} else if lastFork.timestamp != nil && *lastFork.timestamp > *cur.timestamp {
					return fmt.Errorf("unsupported fork ordering: %v enabled at timestamp %v, but %v enabled at timestamp %v",
						lastFork.name, *lastFork.timestamp, cur.name, *cur.timestamp)
				}

				// Timestamp based forks can follow block based ones, but not the other way around
				if lastFork.timestamp != nil && cur.block != nil {
					return fmt.Errorf("unsupported fork ordering: %v used timestamp ordering, but %v reverted to block ordering",
						lastFork.name, cur.name)
				}
			}
		}
		// If it was optional and not set, then ignore it
		if !cur.optional || (cur.block != nil || cur.timestamp != nil) {
			lastFork = cur
		}
	}
	return nil
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, headNumber *big.Int, headTimestamp uint64, genesisTimestamp *uint64) *ConfigCompatError {
	if isForkBlockIncompatible(c.HomesteadBlock, newcfg.HomesteadBlock, headNumber) {
		return newBlockCompatError("Homestead fork block", c.HomesteadBlock, newcfg.HomesteadBlock)
	}
	if isForkBlockIncompatible(c.DAOForkBlock, newcfg.DAOForkBlock, headNumber) {
		return newBlockCompatError("DAO fork block", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if c.IsDAOFork(headNumber) && c.DAOForkSupport != newcfg.DAOForkSupport {
		return newBlockCompatError("DAO fork support flag", c.DAOForkBlock, newcfg.DAOForkBlock)
	}
	if isForkBlockIncompatible(c.EIP150Block, newcfg.EIP150Block, headNumber) {
		return newBlockCompatError("EIP150 fork block", c.EIP150Block, newcfg.EIP150Block)
	}
	if isForkBlockIncompatible(c.EIP155Block, newcfg.EIP155Block, headNumber) {
		return newBlockCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkBlockIncompatible(c.EIP158Block, newcfg.EIP158Block, headNumber) {
		return newBlockCompatError("EIP158 fork block", c.EIP158Block, newcfg.EIP158Block)
	}
	if c.IsEIP158(headNumber) && !configBlockEqual(c.ChainID, newcfg.ChainID) {
		return newBlockCompatError("EIP158 chain ID", c.EIP158Block, newcfg.EIP158Block)
	}
	if isForkBlockIncompatible(c.ByzantiumBlock, newcfg.ByzantiumBlock, headNumber) {
		return newBlockCompatError("Byzantium fork block", c.ByzantiumBlock, newcfg.ByzantiumBlock)
	}
	if isForkBlockIncompatible(c.ConstantinopleBlock, newcfg.ConstantinopleBlock, headNumber) {
		return newBlockCompatError("Constantinople fork block", c.ConstantinopleBlock, newcfg.ConstantinopleBlock)
	}
	if isForkBlockIncompatible(c.PetersburgBlock, newcfg.PetersburgBlock, headNumber) {
		// the only case where we allow Petersburg to be set in the past is if it is equal to Constantinople
		// mainly to satisfy fork ordering requirements which state that Petersburg fork be set if Constantinople fork is set
		if isForkBlockIncompatible(c.ConstantinopleBlock, newcfg.PetersburgBlock, headNumber) {
			return newBlockCompatError("Petersburg fork block", c.PetersburgBlock, newcfg.PetersburgBlock)
		}
	}
	if isForkBlockIncompatible(c.IstanbulBlock, newcfg.IstanbulBlock, headNumber) {
		return newBlockCompatError("Istanbul fork block", c.IstanbulBlock, newcfg.IstanbulBlock)
	}
	if isForkBlockIncompatible(c.MuirGlacierBlock, newcfg.MuirGlacierBlock, headNumber) {
		return newBlockCompatError("Muir Glacier fork block", c.MuirGlacierBlock, newcfg.MuirGlacierBlock)
	}
	if isForkBlockIncompatible(c.BerlinBlock, newcfg.BerlinBlock, headNumber) {
		return newBlockCompatError("Berlin fork block", c.BerlinBlock, newcfg.BerlinBlock)
	}
	if isForkBlockIncompatible(c.LondonBlock, newcfg.LondonBlock, headNumber) {
		return newBlockCompatError("London fork block", c.LondonBlock, newcfg.LondonBlock)
	}
	if isForkBlockIncompatible(c.ArrowGlacierBlock, newcfg.ArrowGlacierBlock, headNumber) {
		return newBlockCompatError("Arrow Glacier fork block", c.ArrowGlacierBlock, newcfg.ArrowGlacierBlock)
	}
	if isForkBlockIncompatible(c.GrayGlacierBlock, newcfg.GrayGlacierBlock, headNumber) {
		return newBlockCompatError("Gray Glacier fork block", c.GrayGlacierBlock, newcfg.GrayGlacierBlock)
	}
	if isForkBlockIncompatible(c.MergeNetsplitBlock, newcfg.MergeNetsplitBlock, headNumber) {
		return newBlockCompatError("Merge netsplit fork block", c.MergeNetsplitBlock, newcfg.MergeNetsplitBlock)
	}
	if isForkTimestampIncompatible(c.ShanghaiTime, newcfg.ShanghaiTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Shanghai fork timestamp", c.ShanghaiTime, newcfg.ShanghaiTime)
	}
	if isForkTimestampIncompatible(c.CancunTime, newcfg.CancunTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Cancun fork timestamp", c.CancunTime, newcfg.CancunTime)
	}
	if isForkTimestampIncompatible(c.PragueTime, newcfg.PragueTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Prague fork timestamp", c.PragueTime, newcfg.PragueTime)
	}
	if isForkTimestampIncompatible(c.VerkleTime, newcfg.VerkleTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Verkle fork timestamp", c.VerkleTime, newcfg.VerkleTime)
	}
	if isForkBlockIncompatible(c.BedrockBlock, newcfg.BedrockBlock, headNumber) {
		return newBlockCompatError("Bedrock fork block", c.BedrockBlock, newcfg.BedrockBlock)
	}
	if isForkTimestampIncompatible(c.RegolithTime, newcfg.RegolithTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Regolith fork timestamp", c.RegolithTime, newcfg.RegolithTime)
	}
	if isForkTimestampIncompatible(c.CanyonTime, newcfg.CanyonTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Canyon fork timestamp", c.CanyonTime, newcfg.CanyonTime)
	}
	if isForkTimestampIncompatible(c.EcotoneTime, newcfg.EcotoneTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Ecotone fork timestamp", c.EcotoneTime, newcfg.EcotoneTime)
	}
	if isForkTimestampIncompatible(c.FjordTime, newcfg.FjordTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Fjord fork timestamp", c.FjordTime, newcfg.FjordTime)
	}
	if isForkTimestampIncompatible(c.GraniteTime, newcfg.GraniteTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Granite fork timestamp", c.GraniteTime, newcfg.GraniteTime)
	}
	if isForkTimestampIncompatible(c.HoloceneTime, newcfg.HoloceneTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Holocene fork timestamp", c.HoloceneTime, newcfg.HoloceneTime)
	}
	if isForkTimestampIncompatible(c.InteropTime, newcfg.InteropTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Interop fork timestamp", c.InteropTime, newcfg.InteropTime)
	}
	return nil
}

// BaseFeeChangeDenominator bounds the amount the base fee can change between blocks.
// The time parameters is the timestamp of the block to determine if Canyon is active or not
func (c *ChainConfig) BaseFeeChangeDenominator(time uint64) uint64 {
	if c.Optimism != nil {
		if c.IsCanyon(time) {
			if c.Optimism.EIP1559DenominatorCanyon == nil || *c.Optimism.EIP1559DenominatorCanyon == 0 {
				panic("invalid ChainConfig.Optimism.EIP1559DenominatorCanyon value: '0' or 'nil'")
			}
			return *c.Optimism.EIP1559DenominatorCanyon
		}
		return c.Optimism.EIP1559Denominator
	}
	return DefaultBaseFeeChangeDenominator
}

// ElasticityMultiplier bounds the maximum gas limit an EIP-1559 block may have.
func (c *ChainConfig) ElasticityMultiplier() uint64 {
	if c.Optimism != nil {
		return c.Optimism.EIP1559Elasticity
	}
	return DefaultElasticityMultiplier
}

// LatestFork returns the latest time-based fork that would be active for the given time.
func (c *ChainConfig) LatestFork(time uint64) forks.Fork {
	// Assume last non-time-based fork has passed.
	london := c.LondonBlock

	switch {
	case c.IsPrague(london, time):
		return forks.Prague
	case c.IsCancun(london, time):
		return forks.Cancun
	case c.IsShanghai(london, time):
		return forks.Shanghai
	default:
		return forks.Paris
	}
}

// isForkBlockIncompatible returns true if a fork scheduled at block s1 cannot be
// rescheduled to block s2 because head is already past the fork and the fork was scheduled after genesis
func isForkBlockIncompatible(s1, s2, head *big.Int) bool {
	return (isBlockForked(s1, head) || isBlockForked(s2, head)) && !configBlockEqual(s1, s2)
}

// isBlockForked returns whether a fork scheduled at block s is active at the
// given head block. Whilst this method is the same as isTimestampForked, they
// are explicitly separate for clearer reading.
func isBlockForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configBlockEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// isForkTimestampIncompatible returns true if a fork scheduled at timestamp s1
// cannot be rescheduled to timestamp s2 because head is already past the fork.
func isForkTimestampIncompatible(s1, s2 *uint64, head uint64, genesis *uint64) bool {
	return (isTimestampForked(s1, head) || isTimestampForked(s2, head)) && !configTimestampEqual(s1, s2) && !isTimestampPreGenesis(s1, genesis) && !isTimestampPreGenesis(s2, genesis)
}

func isTimestampPreGenesis(s, genesis *uint64) bool {
	if s == nil || genesis == nil {
		return false
	}
	return *s < *genesis
}

// isTimestampForked returns whether a fork scheduled at timestamp s is active
// at the given head timestamp. Whilst this method is the same as isBlockForked,
// they are explicitly separate for clearer reading.
func isTimestampForked(s *uint64, head uint64) bool {
	if s == nil {
		return false
	}
	return *s <= head
}

func configTimestampEqual(x, y *uint64) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return *x == *y
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string

	// block numbers of the stored and new configurations if block based forking
	StoredBlock, NewBlock *big.Int

	// timestamps of the stored and new configurations if time based forking
	StoredTime, NewTime *uint64

	// the block number to which the local chain must be rewound to correct the error
	RewindToBlock uint64

	// the timestamp to which the local chain must be rewound to correct the error
	RewindToTime uint64
}

func newBlockCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{
		What:          what,
		StoredBlock:   storedblock,
		NewBlock:      newblock,
		RewindToBlock: 0,
	}
	if rew != nil && rew.Sign() > 0 {
		err.RewindToBlock = rew.Uint64() - 1
	}
	return err
}

func newTimestampCompatError(what string, storedtime, newtime *uint64) *ConfigCompatError {
	var rew *uint64
	switch {
	case storedtime == nil:
		rew = newtime
	case newtime == nil || *storedtime < *newtime:
		rew = storedtime
	default:
		rew = newtime
	}
	err := &ConfigCompatError{
		What:         what,
		StoredTime:   storedtime,
		NewTime:      newtime,
		RewindToTime: 0,
	}
	if rew != nil && *rew != 0 {
		err.RewindToTime = *rew - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	if err.StoredBlock != nil {
		return fmt.Sprintf("mismatching %s in database (have block %d, want block %d, rewindto block %d)", err.What, err.StoredBlock, err.NewBlock, err.RewindToBlock)
	}

	if err.StoredTime == nil && err.NewTime == nil {
		return ""
	} else if err.StoredTime == nil && err.NewTime != nil {
		return fmt.Sprintf("mismatching %s in database (have timestamp nil, want timestamp %d, rewindto timestamp %d)", err.What, *err.NewTime, err.RewindToTime)
	} else if err.StoredTime != nil && err.NewTime == nil {
		return fmt.Sprintf("mismatching %s in database (have timestamp %d, want timestamp nil, rewindto timestamp %d)", err.What, *err.StoredTime, err.RewindToTime)
	}
	return fmt.Sprintf("mismatching %s in database (have timestamp %d, want timestamp %d, rewindto timestamp %d)", err.What, *err.StoredTime, *err.NewTime, err.RewindToTime)
}

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID                                                 *big.Int
	IsHomestead, IsEIP150, IsEIP155, IsEIP158               bool
	IsEIP2929, IsEIP4762                                    bool
	IsByzantium, IsConstantinople, IsPetersburg, IsIstanbul bool
	IsBerlin, IsLondon                                      bool
	IsMerge, IsShanghai, IsCancun, IsPrague                 bool
	IsVerkle                                                bool
	IsOptimismBedrock, IsOptimismRegolith                   bool
	IsOptimismCanyon, IsOptimismFjord                       bool
	IsOptimismGranite, IsOptimismHolocene                   bool
	IsCel2                                                  bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int, isMerge bool, timestamp uint64) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	// disallow setting Merge out of order
	isMerge = isMerge && c.IsLondon(num)
	isVerkle := isMerge && c.IsVerkle(num, timestamp)
	return Rules{
		ChainID:          new(big.Int).Set(chainID),
		IsHomestead:      c.IsHomestead(num),
		IsEIP150:         c.IsEIP150(num),
		IsEIP155:         c.IsEIP155(num),
		IsEIP158:         c.IsEIP158(num),
		IsByzantium:      c.IsByzantium(num),
		IsConstantinople: c.IsConstantinople(num),
		IsPetersburg:     c.IsPetersburg(num),
		IsIstanbul:       c.IsIstanbul(num),
		IsBerlin:         c.IsBerlin(num),
		IsEIP2929:        c.IsBerlin(num) && !isVerkle,
		IsLondon:         c.IsLondon(num),
		IsMerge:          isMerge,
		IsShanghai:       isMerge && c.IsShanghai(num, timestamp),
		IsCancun:         isMerge && c.IsCancun(num, timestamp),
		IsPrague:         isMerge && c.IsPrague(num, timestamp),
		IsVerkle:         isVerkle,
		IsEIP4762:        isVerkle,
		// Optimism
		IsOptimismBedrock:  isMerge && c.IsOptimismBedrock(num),
		IsOptimismRegolith: isMerge && c.IsOptimismRegolith(timestamp),
		IsOptimismCanyon:   isMerge && c.IsOptimismCanyon(timestamp),
		IsOptimismFjord:    isMerge && c.IsOptimismFjord(timestamp),
		IsOptimismGranite:  isMerge && c.IsOptimismGranite(timestamp),
		IsOptimismHolocene: isMerge && c.IsOptimismHolocene(timestamp),
		// Celo
		IsCel2: c.IsCel2(timestamp),
	}
}
