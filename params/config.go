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
	"errors"
	"fmt"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params/forks"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	HoleskyGenesisHash = common.HexToHash("0xb5f7f912443c940f21fd611f12828d75b534364ed9e95ca4e307729a4661bde4")
	SepoliaGenesisHash = common.HexToHash("0x25a5cc106eea7138acab33231d7160d69cb777ee0c2c553fcddf5138993e6dd9")
	HoodiGenesisHash   = common.HexToHash("0xbbe312868b376a3001692a646dd2d7d1e4406380dfd86b98aa8a34d1557c971b")
)

const (
	OPMainnetChainID   = 10
	BaseMainnetChainID = 8453
	baseSepoliaChainID = 84532
)

func newUint64(val uint64) *uint64 { return &val }

var (
	MainnetTerminalTotalDifficulty, _ = new(big.Int).SetString("58_750_000_000_000_000_000_000", 0)

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(1_150_000),
		DAOForkBlock:            big.NewInt(1_920_000),
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(2_463_000),
		EIP155Block:             big.NewInt(2_675_000),
		EIP158Block:             big.NewInt(2_675_000),
		ByzantiumBlock:          big.NewInt(4_370_000),
		ConstantinopleBlock:     big.NewInt(7_280_000),
		PetersburgBlock:         big.NewInt(7_280_000),
		IstanbulBlock:           big.NewInt(9_069_000),
		MuirGlacierBlock:        big.NewInt(9_200_000),
		BerlinBlock:             big.NewInt(12_244_000),
		LondonBlock:             big.NewInt(12_965_000),
		ArrowGlacierBlock:       big.NewInt(13_773_000),
		GrayGlacierBlock:        big.NewInt(15_050_000),
		TerminalTotalDifficulty: MainnetTerminalTotalDifficulty, // 58_750_000_000_000_000_000_000
		ShanghaiTime:            newUint64(1681338455),
		CancunTime:              newUint64(1710338135),
		PragueTime:              newUint64(1746612311),
		OsakaTime:               newUint64(1764798551),
		BPO1Time:                newUint64(1765290071),
		BPO2Time:                newUint64(1767747671),
		DepositContractAddress:  common.HexToAddress("0x00000000219ab540356cbb839cbe05303d7705fa"),
		Ethash:                  new(EthashConfig),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
			Osaka:  DefaultOsakaBlobConfig,
			BPO1:   DefaultBPO1BlobConfig,
			BPO2:   DefaultBPO2BlobConfig,
		},
	}
	// HoleskyChainConfig contains the chain parameters to run a node on the Holesky test network.
	HoleskyChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(17000),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        nil,
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       nil,
		GrayGlacierBlock:        nil,
		TerminalTotalDifficulty: big.NewInt(0),
		MergeNetsplitBlock:      nil,
		ShanghaiTime:            newUint64(1696000704),
		CancunTime:              newUint64(1707305664),
		PragueTime:              newUint64(1740434112),
		OsakaTime:               newUint64(1759308480),
		BPO1Time:                newUint64(1759800000),
		BPO2Time:                newUint64(1760389824),
		DepositContractAddress:  common.HexToAddress("0x4242424242424242424242424242424242424242"),
		Ethash:                  new(EthashConfig),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
			Osaka:  DefaultOsakaBlobConfig,
			BPO1:   DefaultBPO1BlobConfig,
			BPO2:   DefaultBPO2BlobConfig,
		},
	}
	// SepoliaChainConfig contains the chain parameters to run a node on the Sepolia test network.
	SepoliaChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(11155111),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       nil,
		GrayGlacierBlock:        nil,
		TerminalTotalDifficulty: big.NewInt(17_000_000_000_000_000),
		MergeNetsplitBlock:      big.NewInt(1735371),
		ShanghaiTime:            newUint64(1677557088),
		CancunTime:              newUint64(1706655072),
		PragueTime:              newUint64(1741159776),
		OsakaTime:               newUint64(1760427360),
		BPO1Time:                newUint64(1761017184),
		BPO2Time:                newUint64(1761607008),
		DepositContractAddress:  common.HexToAddress("0x7f02c3e3c98b133055b8b348b2ac625669ed295d"),
		Ethash:                  new(EthashConfig),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
			Osaka:  DefaultOsakaBlobConfig,
			BPO1:   DefaultBPO1BlobConfig,
			BPO2:   DefaultBPO2BlobConfig,
		},
	}
	// HoodiChainConfig contains the chain parameters to run a node on the Hoodi test network.
	HoodiChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(560048),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          true,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       nil,
		GrayGlacierBlock:        nil,
		TerminalTotalDifficulty: big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		PragueTime:              newUint64(1742999832),
		OsakaTime:               newUint64(1761677592),
		BPO1Time:                newUint64(1762365720),
		BPO2Time:                newUint64(1762955544),
		DepositContractAddress:  common.HexToAddress("0x00000000219ab540356cBB839Cbe05303d7705Fa"),
		Ethash:                  new(EthashConfig),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
			Osaka:  DefaultOsakaBlobConfig,
			BPO1:   DefaultBPO1BlobConfig,
			BPO2:   DefaultBPO2BlobConfig,
		},
	}
	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Ethash consensus.
	AllEthashProtocolChanges = &ChainConfig{
		ChainID:                 big.NewInt(1337),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		TerminalTotalDifficulty: big.NewInt(math.MaxInt64),
		MergeNetsplitBlock:      nil,
		ShanghaiTime:            nil,
		CancunTime:              nil,
		PragueTime:              nil,
		OsakaTime:               nil,
		VerkleTime:              nil,
		Ethash:                  new(EthashConfig),
		Clique:                  nil,
	}

	AllDevChainProtocolChanges = &ChainConfig{
		ChainID:                 big.NewInt(1337),
		HomesteadBlock:          big.NewInt(0),
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		TerminalTotalDifficulty: big.NewInt(0),
		PragueTime:              newUint64(0),
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
		},
	}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	AllCliqueProtocolChanges = &ChainConfig{
		ChainID:                 big.NewInt(1337),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       nil,
		GrayGlacierBlock:        nil,
		MergeNetsplitBlock:      nil,
		ShanghaiTime:            nil,
		CancunTime:              nil,
		PragueTime:              nil,
		OsakaTime:               nil,
		VerkleTime:              nil,
		TerminalTotalDifficulty: big.NewInt(math.MaxInt64),
		Ethash:                  nil,
		Clique:                  &CliqueConfig{Period: 0, Epoch: 30000},
	}

	// TestChainConfig contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers for testing purposes.
	TestChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		MergeNetsplitBlock:      nil,
		ShanghaiTime:            nil,
		CancunTime:              nil,
		PragueTime:              nil,
		OsakaTime:               nil,
		VerkleTime:              nil,
		TerminalTotalDifficulty: big.NewInt(math.MaxInt64),
		Ethash:                  new(EthashConfig),
		Clique:                  nil,
	}

	// MergedTestChainConfig contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers for testing purposes.
	MergedTestChainConfig = &ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          big.NewInt(0),
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             big.NewInt(0),
		EIP155Block:             big.NewInt(0),
		EIP158Block:             big.NewInt(0),
		ByzantiumBlock:          big.NewInt(0),
		ConstantinopleBlock:     big.NewInt(0),
		PetersburgBlock:         big.NewInt(0),
		IstanbulBlock:           big.NewInt(0),
		MuirGlacierBlock:        big.NewInt(0),
		BerlinBlock:             big.NewInt(0),
		LondonBlock:             big.NewInt(0),
		ArrowGlacierBlock:       big.NewInt(0),
		GrayGlacierBlock:        big.NewInt(0),
		MergeNetsplitBlock:      big.NewInt(0),
		ShanghaiTime:            newUint64(0),
		CancunTime:              newUint64(0),
		PragueTime:              newUint64(0),
		OsakaTime:               newUint64(0),
		VerkleTime:              nil,
		TerminalTotalDifficulty: big.NewInt(0),
		Ethash:                  new(EthashConfig),
		Clique:                  nil,
		BlobScheduleConfig: &BlobScheduleConfig{
			Cancun: DefaultCancunBlobConfig,
			Prague: DefaultPragueBlobConfig,
			Osaka:  DefaultOsakaBlobConfig,
		},
	}

	// NonActivatedConfig defines the chain configuration without activating
	// any protocol change (EIPs).
	NonActivatedConfig = &ChainConfig{
		ChainID:                 big.NewInt(1),
		HomesteadBlock:          nil,
		DAOForkBlock:            nil,
		DAOForkSupport:          false,
		EIP150Block:             nil,
		EIP155Block:             nil,
		EIP158Block:             nil,
		ByzantiumBlock:          nil,
		ConstantinopleBlock:     nil,
		PetersburgBlock:         nil,
		IstanbulBlock:           nil,
		MuirGlacierBlock:        nil,
		BerlinBlock:             nil,
		LondonBlock:             nil,
		ArrowGlacierBlock:       nil,
		GrayGlacierBlock:        nil,
		MergeNetsplitBlock:      nil,
		ShanghaiTime:            nil,
		CancunTime:              nil,
		PragueTime:              nil,
		OsakaTime:               nil,
		VerkleTime:              nil,
		TerminalTotalDifficulty: big.NewInt(math.MaxInt64),
		Ethash:                  new(EthashConfig),
		Clique:                  nil,
	}
	TestRules = TestChainConfig.Rules(new(big.Int), false, 0)

	// OP-Stack chain config with bedrock starting a block 5, introduced for historical endpoint testing, largely based on the clique config
	OptimismTestCliqueConfig = func() *ChainConfig {
		conf := *AllCliqueProtocolChanges // copy the config
		conf.Clique = nil
		conf.BedrockBlock = big.NewInt(5)
		conf.Optimism = &OptimismConfig{EIP1559Elasticity: 50, EIP1559Denominator: 10}
		return &conf
	}()

	// OP-Stack chain config with all production forks activated, based on the MergedTestChainConfig
	OptimismTestConfig = func() *ChainConfig {
		conf := *MergedTestChainConfig // copy the config
		conf.BlobScheduleConfig = nil
		conf.OsakaTime = nil // needs to be removed when production fork introduces Osaka
		conf.BedrockBlock = big.NewInt(0)
		zero := uint64(0)
		conf.RegolithTime = &zero
		conf.CanyonTime = &zero
		conf.EcotoneTime = &zero
		conf.FjordTime = &zero
		conf.GraniteTime = &zero
		conf.HoloceneTime = &zero
		conf.IsthmusTime = &zero
		conf.JovianTime = nil
		conf.InteropTime = nil
		conf.Optimism = &OptimismConfig{EIP1559Elasticity: 6, EIP1559Denominator: 50, EIP1559DenominatorCanyon: uint64ptr(250)}
		return &conf
	}()
)

var (
	// DefaultCancunBlobConfig is the default blob configuration for the Cancun fork.
	DefaultCancunBlobConfig = &BlobConfig{
		Target:         3,
		Max:            6,
		UpdateFraction: 3338477,
	}
	// DefaultPragueBlobConfig is the default blob configuration for the Prague fork.
	DefaultPragueBlobConfig = &BlobConfig{
		Target:         6,
		Max:            9,
		UpdateFraction: 5007716,
	}
	// DefaultOsakaBlobConfig is the default blob configuration for the Osaka fork.
	DefaultOsakaBlobConfig = &BlobConfig{
		Target:         6,
		Max:            9,
		UpdateFraction: 5007716,
	}
	// DefaultBPO1BlobConfig is the default blob configuration for the Osaka fork.
	DefaultBPO1BlobConfig = &BlobConfig{
		Target:         10,
		Max:            15,
		UpdateFraction: 8346193,
	}
	// DefaultBPO1BlobConfig is the default blob configuration for the Osaka fork.
	DefaultBPO2BlobConfig = &BlobConfig{
		Target:         14,
		Max:            21,
		UpdateFraction: 11684671,
	}
	// DefaultBPO1BlobConfig is the default blob configuration for the Osaka fork.
	DefaultBPO3BlobConfig = &BlobConfig{
		Target:         21,
		Max:            32,
		UpdateFraction: 20609697,
	}
	// DefaultBPO1BlobConfig is the default blob configuration for the Osaka fork.
	DefaultBPO4BlobConfig = &BlobConfig{
		Target:         14,
		Max:            21,
		UpdateFraction: 13739630,
	}
	// DefaultBlobSchedule is the latest configured blob schedule for Ethereum mainnet.
	DefaultBlobSchedule = &BlobScheduleConfig{
		Cancun: DefaultCancunBlobConfig,
		Prague: DefaultPragueBlobConfig,
		Osaka:  DefaultOsakaBlobConfig,
	}
)

// NetworkNames are user friendly names to use in the chain spec banner.
var NetworkNames = map[string]string{
	MainnetChainConfig.ChainID.String(): "mainnet",
	SepoliaChainConfig.ChainID.String(): "sepolia",
	HoleskyChainConfig.ChainID.String(): "holesky",
	HoodiChainConfig.ChainID.String():   "hoodi",
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
	OsakaTime    *uint64 `json:"osakaTime,omitempty"`    // Osaka switch time (nil = no fork, 0 = already on osaka)
	VerkleTime   *uint64 `json:"verkleTime,omitempty"`   // Verkle switch time (nil = no fork, 0 = already on verkle)
	BPO1Time     *uint64 `json:"bpo1Time,omitempty"`     // BPO1 switch time (nil = no fork, 0 = already on bpo1)
	BPO2Time     *uint64 `json:"bpo2Time,omitempty"`     // BPO2 switch time (nil = no fork, 0 = already on bpo2)
	BPO3Time     *uint64 `json:"bpo3Time,omitempty"`     // BPO3 switch time (nil = no fork, 0 = already on bpo3)
	BPO4Time     *uint64 `json:"bpo4Time,omitempty"`     // BPO4 switch time (nil = no fork, 0 = already on bpo4)
	BPO5Time     *uint64 `json:"bpo5Time,omitempty"`     // BPO5 switch time (nil = no fork, 0 = already on bpo5)

	BedrockBlock *big.Int `json:"bedrockBlock,omitempty"` // Bedrock switch block (nil = no fork, 0 = already on optimism bedrock)
	RegolithTime *uint64  `json:"regolithTime,omitempty"` // Regolith switch time (nil = no fork, 0 = already on optimism regolith)
	CanyonTime   *uint64  `json:"canyonTime,omitempty"`   // Canyon switch time (nil = no fork, 0 = already on optimism canyon)
	// Delta: the Delta upgrade does not affect the execution-layer, and is thus not configurable in the chain config.
	EcotoneTime  *uint64 `json:"ecotoneTime,omitempty"`  // Ecotone switch time (nil = no fork, 0 = already on optimism ecotone)
	FjordTime    *uint64 `json:"fjordTime,omitempty"`    // Fjord switch time (nil = no fork, 0 = already on Optimism Fjord)
	GraniteTime  *uint64 `json:"graniteTime,omitempty"`  // Granite switch time (nil = no fork, 0 = already on Optimism Granite)
	HoloceneTime *uint64 `json:"holoceneTime,omitempty"` // Holocene switch time (nil = no fork, 0 = already on Optimism Holocene)
	IsthmusTime  *uint64 `json:"isthmusTime,omitempty"`  // Isthmus switch time (nil = no fork, 0 = already on Optimism Isthmus)
	JovianTime   *uint64 `json:"jovianTime,omitempty"`   // Jovian switch time (nil = no fork, 0 = already on Optimism Jovian)

	InteropTime *uint64 `json:"interopTime,omitempty"` // Interop switch time (nil = no fork, 0 = already on optimism interop)

	// TerminalTotalDifficulty is the amount of total difficulty reached by
	// the network that triggers the consensus upgrade.
	TerminalTotalDifficulty *big.Int `json:"terminalTotalDifficulty,omitempty"`

	DepositContractAddress common.Address `json:"depositContractAddress,omitempty"`

	// EnableVerkleAtGenesis is a flag that specifies whether the network uses
	// the Verkle tree starting from the genesis block. If set to true, the
	// genesis state will be committed using the Verkle tree, eliminating the
	// need for any Verkle transition later.
	//
	// This is a temporary flag only for verkle devnet testing, where verkle is
	// activated at genesis, and the configured activation date has already passed.
	//
	// In production networks (mainnet and public testnets), verkle activation
	// always occurs after the genesis block, making this flag irrelevant in
	// those cases.
	EnableVerkleAtGenesis bool `json:"enableVerkleAtGenesis,omitempty"`

	// Various consensus engines
	Ethash             *EthashConfig       `json:"ethash,omitempty"`
	Clique             *CliqueConfig       `json:"clique,omitempty"`
	BlobScheduleConfig *BlobScheduleConfig `json:"blobSchedule,omitempty"`

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
		banner += "Consensus: Beacon (proof-of-stake), merged from Ethash (proof-of-work)\n"
	case c.Clique != nil:
		banner += "Consensus: Beacon (proof-of-stake), merged from Clique (proof-of-authority)\n"
	default:
		banner += "Consensus: unknown\n"
	}
	banner += "\n"

	// Create a list of forks with a short description of them. Forks that only
	// makes sense for mainnet should be optional at printing to avoid bloating
	// the output for testnets and private networks.
	banner += "Pre-Merge hard forks (block based):\n"
	banner += fmt.Sprintf(" - Homestead:                   #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/homestead/__init__.py.html)\n", c.HomesteadBlock)
	if c.DAOForkBlock != nil {
		banner += fmt.Sprintf(" - DAO Fork:                    #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/dao_fork/__init__.py.html)\n", c.DAOForkBlock)
	}
	banner += fmt.Sprintf(" - Tangerine Whistle (EIP 150): #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/tangerine_whistle/__init__.py.html)\n", c.EIP150Block)
	banner += fmt.Sprintf(" - Spurious Dragon/1 (EIP 155): #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/spurious_dragon/__init__.py.html)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Spurious Dragon/2 (EIP 158): #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/spurious_dragon/__init__.py.html)\n", c.EIP155Block)
	banner += fmt.Sprintf(" - Byzantium:                   #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/byzantium/__init__.py.html)\n", c.ByzantiumBlock)
	banner += fmt.Sprintf(" - Constantinople:              #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/constantinople/__init__.py.html)\n", c.ConstantinopleBlock)
	banner += fmt.Sprintf(" - Petersburg:                  #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/constantinople/__init__.py.html)\n", c.PetersburgBlock)
	banner += fmt.Sprintf(" - Istanbul:                    #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/istanbul/__init__.py.html)\n", c.IstanbulBlock)
	if c.MuirGlacierBlock != nil {
		banner += fmt.Sprintf(" - Muir Glacier:                #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/muir_glacier/__init__.py.html)\n", c.MuirGlacierBlock)
	}
	banner += fmt.Sprintf(" - Berlin:                      #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/berlin/__init__.py.html)\n", c.BerlinBlock)
	banner += fmt.Sprintf(" - London:                      #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/london/__init__.py.html)\n", c.LondonBlock)
	if c.ArrowGlacierBlock != nil {
		banner += fmt.Sprintf(" - Arrow Glacier:               #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/arrow_glacier/__init__.py.html)\n", c.ArrowGlacierBlock)
	}
	if c.GrayGlacierBlock != nil {
		banner += fmt.Sprintf(" - Gray Glacier:                #%-8v (https://ethereum.github.io/execution-specs/src/ethereum/forks/gray_glacier/__init__.py.html)\n", c.GrayGlacierBlock)
	}
	banner += "\n"

	// Add a special section for the merge as it's non-obvious
	banner += "Merge configured:\n"
	banner += " - Hard-fork specification:    https://ethereum.github.io/execution-specs/src/ethereum/forks/paris/__init__.py.html\n"
	banner += " - Network known to be merged\n"
	banner += fmt.Sprintf(" - Total terminal difficulty:  %v\n", c.TerminalTotalDifficulty)
	if c.MergeNetsplitBlock != nil {
		banner += fmt.Sprintf(" - Merge netsplit block:       #%-8v\n", c.MergeNetsplitBlock)
	}
	banner += "\n"

	// Create a list of forks post-merge
	banner += "Post-Merge hard forks (timestamp based):\n"
	if c.ShanghaiTime != nil {
		banner += fmt.Sprintf(" - Shanghai:                    @%-10v (https://ethereum.github.io/execution-specs/src/ethereum/forks/shanghai/__init__.py.html)\n", *c.ShanghaiTime)
	}
	if c.CancunTime != nil {
		banner += fmt.Sprintf(" - Cancun:                      @%-10v (https://ethereum.github.io/execution-specs/src/ethereum/forks/cancun/__init__.py.html)\n", *c.CancunTime)
	}
	if c.PragueTime != nil {
		banner += fmt.Sprintf(" - Prague:                      @%-10v (https://ethereum.github.io/execution-specs/src/ethereum/forks/prague/__init__.py.html)\n", *c.PragueTime)
	}
	if c.OsakaTime != nil {
		banner += fmt.Sprintf(" - Osaka:                       @%-10v (https://ethereum.github.io/execution-specs/src/ethereum/forks/osaka/__init__.py.html)\n", *c.OsakaTime)
	}
	if c.VerkleTime != nil {
		banner += fmt.Sprintf(" - Verkle:                      @%-10v\n", *c.VerkleTime)
	}
	if c.BPO1Time != nil {
		banner += fmt.Sprintf(" - BPO1:                        @%-10v\n", *c.BPO1Time)
	}
	if c.BPO2Time != nil {
		banner += fmt.Sprintf(" - BPO2:                        @%-10v\n", *c.BPO2Time)
	}
	if c.BPO3Time != nil {
		banner += fmt.Sprintf(" - BPO3:                        @%-10v\n", *c.BPO3Time)
	}
	if c.BPO4Time != nil {
		banner += fmt.Sprintf(" - BPO4:                        @%-10v\n", *c.BPO4Time)
	}
	if c.BPO5Time != nil {
		banner += fmt.Sprintf(" - BPO5:                        @%-10v\n", *c.BPO5Time)
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
		banner += fmt.Sprintf(" - Holocene:                    @%-10v\n", *c.HoloceneTime)
	}
	if c.IsthmusTime != nil {
		banner += fmt.Sprintf(" - Isthmus:                     @%-10v\n", *c.IsthmusTime)
	}
	if c.JovianTime != nil {
		banner += fmt.Sprintf(" - Jovian:                      @%-10v\n", *c.JovianTime)
	}
	if c.InteropTime != nil {
		banner += fmt.Sprintf(" - Interop:                     @%-10v\n", *c.InteropTime)
	}
	return banner
}

// BlobConfig specifies the target and max blobs per block for the associated fork.
type BlobConfig struct {
	Target         int    `json:"target"`
	Max            int    `json:"max"`
	UpdateFraction uint64 `json:"baseFeeUpdateFraction"`
}

// BlobScheduleConfig determines target and max number of blobs allow per fork.
type BlobScheduleConfig struct {
	Cancun *BlobConfig `json:"cancun,omitempty"`
	Prague *BlobConfig `json:"prague,omitempty"`
	Osaka  *BlobConfig `json:"osaka,omitempty"`
	Verkle *BlobConfig `json:"verkle,omitempty"`
	BPO1   *BlobConfig `json:"bpo1,omitempty"`
	BPO2   *BlobConfig `json:"bpo2,omitempty"`
	BPO3   *BlobConfig `json:"bpo3,omitempty"`
	BPO4   *BlobConfig `json:"bpo4,omitempty"`
	BPO5   *BlobConfig `json:"bpo5,omitempty"`
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

// IsPostMerge reports whether the given block number is assumed to be post-merge.
// Here we check the MergeNetsplitBlock to allow configuring networks with a PoW or
// PoA chain for unit testing purposes.
func (c *ChainConfig) IsPostMerge(blockNum uint64, timestamp uint64) bool {
	mergedAtGenesis := c.TerminalTotalDifficulty != nil && c.TerminalTotalDifficulty.Sign() == 0
	return mergedAtGenesis ||
		c.MergeNetsplitBlock != nil && blockNum >= c.MergeNetsplitBlock.Uint64() ||
		c.ShanghaiTime != nil && timestamp >= *c.ShanghaiTime ||
		// If OP-Stack then bedrock activation number determines when TTD (eth Merge) has been reached.
		c.IsOptimismBedrock(new(big.Int).SetUint64(blockNum))
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

// IsOsaka returns whether time is either equal to the Osaka fork time or greater.
func (c *ChainConfig) IsOsaka(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.OsakaTime, time)
}

// IsVerkle returns whether time is either equal to the Verkle fork time or greater.
func (c *ChainConfig) IsVerkle(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.VerkleTime, time)
}

// IsBPO1 returns whether time is either equal to the BPO1 fork time or greater.
func (c *ChainConfig) IsBPO1(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO1Time, time)
}

// IsBPO2 returns whether time is either equal to the BPO2 fork time or greater.
func (c *ChainConfig) IsBPO2(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO2Time, time)
}

// IsBPO3 returns whether time is either equal to the BPO3 fork time or greater.
func (c *ChainConfig) IsBPO3(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO3Time, time)
}

// IsBPO4 returns whether time is either equal to the BPO4 fork time or greater.
func (c *ChainConfig) IsBPO4(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO4Time, time)
}

// IsBPO5 returns whether time is either equal to the BPO5 fork time or greater.
func (c *ChainConfig) IsBPO5(num *big.Int, time uint64) bool {
	return c.IsLondon(num) && isTimestampForked(c.BPO5Time, time)
}

// IsVerkleGenesis checks whether the verkle fork is activated at the genesis block.
//
// Verkle mode is considered enabled if the verkle fork time is configured,
// regardless of whether the local time has surpassed the fork activation time.
// This is a temporary workaround for verkle devnet testing, where verkle is
// activated at genesis, and the configured activation date has already passed.
//
// In production networks (mainnet and public testnets), verkle activation
// always occurs after the genesis block, making this function irrelevant in
// those cases.
func (c *ChainConfig) IsVerkleGenesis() bool {
	return c.EnableVerkleAtGenesis
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

func (c *ChainConfig) IsIsthmus(time uint64) bool {
	return isTimestampForked(c.IsthmusTime, time)
}

func (c *ChainConfig) IsJovian(time uint64) bool {
	return isTimestampForked(c.JovianTime, time)
}

func (c *ChainConfig) IsInterop(time uint64) bool {
	return isTimestampForked(c.InteropTime, time)
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

func (c *ChainConfig) IsOptimismIsthmus(time uint64) bool {
	return c.IsOptimism() && c.IsIsthmus(time)
}

func (c *ChainConfig) IsOptimismJovian(time uint64) bool {
	return c.IsOptimism() && c.IsJovian(time)
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
		log.Info("Checking compatibility", "height", bhead, "time", btime, "error", err)
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
		{name: "osakaTime", timestamp: c.OsakaTime, optional: true},
		{name: "verkleTime", timestamp: c.VerkleTime, optional: true},
		{name: "bpo1", timestamp: c.BPO1Time, optional: true},
		{name: "bpo2", timestamp: c.BPO2Time, optional: true},
		{name: "bpo3", timestamp: c.BPO3Time, optional: true},
		{name: "bpo4", timestamp: c.BPO4Time, optional: true},
		{name: "bpo5", timestamp: c.BPO5Time, optional: true},
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

	// OP-Stack chains don't support blobs, and must have a nil BlobScheduleConfig.
	if c.IsOptimism() {
		if c.BlobScheduleConfig == nil {
			return nil
		} else {
			return errors.New("OP-Stack chains must have empty blob configuration")
		}
	}

	// Check that all forks with blobs explicitly define the blob schedule configuration.
	bsc := c.BlobScheduleConfig
	if bsc == nil {
		bsc = new(BlobScheduleConfig)
	}
	for _, cur := range []struct {
		name      string
		timestamp *uint64
		config    *BlobConfig
	}{
		{name: "cancun", timestamp: c.CancunTime, config: bsc.Cancun},
		{name: "prague", timestamp: c.PragueTime, config: bsc.Prague},
		{name: "osaka", timestamp: c.OsakaTime, config: bsc.Osaka},
		{name: "bpo1", timestamp: c.BPO1Time, config: bsc.BPO1},
		{name: "bpo2", timestamp: c.BPO2Time, config: bsc.BPO2},
		{name: "bpo3", timestamp: c.BPO3Time, config: bsc.BPO3},
		{name: "bpo4", timestamp: c.BPO4Time, config: bsc.BPO4},
		{name: "bpo5", timestamp: c.BPO5Time, config: bsc.BPO5},
	} {
		if cur.config != nil {
			if err := cur.config.validate(); err != nil {
				return fmt.Errorf("invalid chain configuration in blobSchedule for fork %q: %v", cur.name, err)
			}
		}
		if cur.timestamp != nil {
			// If the fork is configured, a blob schedule must be defined for it.
			if cur.config == nil {
				return fmt.Errorf("invalid chain configuration: missing entry for fork %q in blobSchedule", cur.name)
			}
		}
	}
	return nil
}

func (bc *BlobConfig) validate() error {
	if bc.Max < 0 {
		return errors.New("max < 0")
	}
	if bc.Target < 0 {
		return errors.New("target < 0")
	}
	if bc.UpdateFraction == 0 {
		return errors.New("update fraction must be defined and non-zero")
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
	if isForkTimestampIncompatible(c.OsakaTime, newcfg.OsakaTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Osaka fork timestamp", c.OsakaTime, newcfg.OsakaTime)
	}
	if isForkTimestampIncompatible(c.VerkleTime, newcfg.VerkleTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Verkle fork timestamp", c.VerkleTime, newcfg.VerkleTime)
	}
	if isForkTimestampIncompatible(c.BPO1Time, newcfg.BPO1Time, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("BPO1 fork timestamp", c.BPO1Time, newcfg.BPO1Time)
	}
	if isForkTimestampIncompatible(c.BPO2Time, newcfg.BPO2Time, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("BPO2 fork timestamp", c.BPO2Time, newcfg.BPO2Time)
	}
	if isForkTimestampIncompatible(c.BPO3Time, newcfg.BPO3Time, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("BPO3 fork timestamp", c.BPO3Time, newcfg.BPO3Time)
	}
	if isForkTimestampIncompatible(c.BPO4Time, newcfg.BPO4Time, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("BPO4 fork timestamp", c.BPO4Time, newcfg.BPO4Time)
	}
	if isForkTimestampIncompatible(c.BPO5Time, newcfg.BPO5Time, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("BPO5 fork timestamp", c.BPO5Time, newcfg.BPO5Time)
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
	if isForkTimestampIncompatible(c.IsthmusTime, newcfg.IsthmusTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Isthmus fork timestamp", c.IsthmusTime, newcfg.IsthmusTime)
	}
	if isForkTimestampIncompatible(c.JovianTime, newcfg.JovianTime, headTimestamp, genesisTimestamp) {
		return newTimestampCompatError("Jovian fork timestamp", c.JovianTime, newcfg.JovianTime)
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
	case c.IsBPO5(london, time):
		return forks.BPO5
	case c.IsBPO4(london, time):
		return forks.BPO4
	case c.IsBPO3(london, time):
		return forks.BPO3
	case c.IsBPO2(london, time):
		return forks.BPO2
	case c.IsBPO1(london, time):
		return forks.BPO1
	case c.IsOsaka(london, time):
		return forks.Osaka
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

// BlobConfig returns the blob config associated with the provided fork.
func (c *ChainConfig) BlobConfig(fork forks.Fork) *BlobConfig {
	// TODO: https://github.com/ethereum-optimism/op-geth/issues/685
	// This function has a bug.
	switch fork {
	case forks.BPO5:
		return c.BlobScheduleConfig.BPO5
	case forks.BPO4:
		return c.BlobScheduleConfig.BPO4
	case forks.BPO3:
		return c.BlobScheduleConfig.BPO3
	case forks.BPO2:
		return c.BlobScheduleConfig.BPO2
	case forks.BPO1:
		return c.BlobScheduleConfig.BPO1
	case forks.Osaka:
		return c.BlobScheduleConfig.Osaka
	case forks.Prague:
		return c.BlobScheduleConfig.Prague
	case forks.Cancun:
		return c.BlobScheduleConfig.Cancun
	default:
		return nil
	}
}

// ActiveSystemContracts returns the currently active system contracts at the
// given timestamp.
func (c *ChainConfig) ActiveSystemContracts(time uint64) map[string]common.Address {
	fork := c.LatestFork(time)
	active := make(map[string]common.Address)
	if fork >= forks.Osaka {
		// no new system contracts
	}
	if fork >= forks.Prague {
		active["CONSOLIDATION_REQUEST_PREDEPLOY_ADDRESS"] = ConsolidationQueueAddress
		active["DEPOSIT_CONTRACT_ADDRESS"] = c.DepositContractAddress
		active["HISTORY_STORAGE_ADDRESS"] = HistoryStorageAddress
		active["WITHDRAWAL_REQUEST_PREDEPLOY_ADDRESS"] = WithdrawalQueueAddress
	}
	if fork >= forks.Cancun {
		active["BEACON_ROOTS_ADDRESS"] = BeaconRootsAddress
	}
	return active
}

// Timestamp returns the timestamp associated with the fork or returns nil if
// the fork isn't defined or isn't a time-based fork.
func (c *ChainConfig) Timestamp(fork forks.Fork) *uint64 {
	switch {
	case fork == forks.BPO5:
		return c.BPO5Time
	case fork == forks.BPO4:
		return c.BPO4Time
	case fork == forks.BPO3:
		return c.BPO3Time
	case fork == forks.BPO2:
		return c.BPO2Time
	case fork == forks.BPO1:
		return c.BPO1Time
	case fork == forks.Osaka:
		return c.OsakaTime
	case fork == forks.Prague:
		return c.PragueTime
	case fork == forks.Cancun:
		return c.CancunTime
	case fork == forks.Shanghai:
		return c.ShanghaiTime
	default:
		return nil
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
	return (isTimestampForked(s1, head) || isTimestampForked(s2, head)) && !configTimestampEqual(s1, s2) && !(isTimestampPreGenesis(s1, genesis) && isTimestampPreGenesis(s2, genesis))
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
	IsMerge, IsShanghai, IsCancun, IsPrague, IsOsaka        bool
	IsVerkle                                                bool
	IsOptimismBedrock, IsOptimismRegolith                   bool
	IsOptimismCanyon, IsOptimismFjord                       bool
	IsOptimismGranite, IsOptimismHolocene                   bool
	IsOptimismIsthmus, IsOptimismJovian                     bool
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
		IsOsaka:          isMerge && c.IsOsaka(num, timestamp),
		IsVerkle:         isVerkle,
		IsEIP4762:        isVerkle,
		// Optimism
		IsOptimismBedrock:  isMerge && c.IsOptimismBedrock(num),
		IsOptimismRegolith: isMerge && c.IsOptimismRegolith(timestamp),
		IsOptimismCanyon:   isMerge && c.IsOptimismCanyon(timestamp),
		IsOptimismFjord:    isMerge && c.IsOptimismFjord(timestamp),
		IsOptimismGranite:  isMerge && c.IsOptimismGranite(timestamp),
		IsOptimismHolocene: isMerge && c.IsOptimismHolocene(timestamp),
		IsOptimismIsthmus:  isMerge && c.IsOptimismIsthmus(timestamp),
		IsOptimismJovian:   isMerge && c.IsOptimismJovian(timestamp),
	}
}

func (c *ChainConfig) HasOptimismWithdrawalsRoot(blockTime uint64) bool {
	return c.IsOptimismIsthmus(blockTime)
}

// CheckOptimismValidity checks for OP Stack chains:
// - the EIP159 params are set
// - the Ethereum forks are set to the same time as the OP Stack forks that imply them
func (c *ChainConfig) CheckOptimismValidity() error {
	if c.Optimism == nil {
		return nil
	}

	if c.Optimism.EIP1559Denominator == 0 {
		return errors.New("zero EIP1559Denominator")
	}
	if c.Optimism.EIP1559Elasticity == 0 {
		return errors.New("zero EIP1559Elasticity")
	}
	if c.CanyonTime != nil && (c.Optimism.EIP1559DenominatorCanyon == nil || *c.Optimism.EIP1559DenominatorCanyon == 0) {
		return errors.New("missing or zero EIP1559DenominatorCanyon")
	}

	if !equalPtrValues(c.ShanghaiTime, c.CanyonTime) {
		return fmt.Errorf("ShanghaiTime (%s) must equal CanyonTime (%s)", ptrValueString(c.ShanghaiTime), ptrValueString(c.CanyonTime))
	}
	if !equalPtrValues(c.CancunTime, c.EcotoneTime) {
		return fmt.Errorf("CancunTime (%s) must equal EcotoneTime (%s)", ptrValueString(c.CancunTime), ptrValueString(c.EcotoneTime))
	}
	if !equalPtrValues(c.PragueTime, c.IsthmusTime) {
		return fmt.Errorf("PragueTime (%s) must equal IsthmusTime (%s)", ptrValueString(c.PragueTime), ptrValueString(c.IsthmusTime))
	}

	return nil
}

func equalPtrValues[T comparable](a, b *T) bool {
	// also captures nil == nil
	return a == b || (a != nil && b != nil && *a == *b)
}

func ptrValueString[T any](t *T) string {
	if t == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%v", *t)
}
