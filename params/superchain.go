package params

import (
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/superchain-registry/superchain"
	"github.com/ethereum/go-ethereum/common"
)

var OPStackSupport = ToProtocolVersion(0, 3, 1, 0, 1)

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

// ProtocolVersion encodes the OP-Stack protocol version. See OP-Stack superchain-upgrade specification.
type ProtocolVersion [32]byte

func (p ProtocolVersion) MarshalText() ([]byte, error) {
	return common.Hash(p).MarshalText()
}

func (p *ProtocolVersion) UnmarshalText(input []byte) error {
	return (*common.Hash)(p).UnmarshalText(input)
}

func (p ProtocolVersion) Parse() (versionType uint8, build uint64, major, minor, patch, preRelease uint32) {
	versionType = p[0]
	if versionType != 0 {
		return
	}
	// bytes 1:8 reserved for future use
	build = binary.BigEndian.Uint64(p[8:16])       // differentiates forks and custom-builds of standard protocol
	major = binary.BigEndian.Uint32(p[16:20])      // incompatible API changes
	minor = binary.BigEndian.Uint32(p[20:24])      // identifies additional functionality in backwards compatible manner
	patch = binary.BigEndian.Uint32(p[24:28])      // identifies backward-compatible bug-fixes
	preRelease = binary.BigEndian.Uint32(p[28:32]) // identifies unstable versions that may not satisfy the above
	return
}

func (p ProtocolVersion) String() string {
	versionType, build, major, minor, patch, preRelease := p.Parse()
	if versionType != 0 {
		return "v0.0.0-unknown." + common.Hash(p).String()
	}
	ver := fmt.Sprintf("v%d.%d.%d", major, minor, patch)
	if preRelease != 0 {
		ver += fmt.Sprintf("-%d", preRelease)
	}
	if build != 0 {
		ver += fmt.Sprintf("+%d", build)
	}
	return ver
}

// ProtocolVersionComparison is used to identify how far ahead/outdated a protocol version is relative to another.
// This value is used in metrics and switch comparisons, to easily identify each type of version difference.
// Negative values mean the version is outdated.
// Positive values mean the version is up-to-date.
// Matching versions have a 0.
type ProtocolVersionComparison int

const (
	AheadMajor         ProtocolVersionComparison = 4
	OutdatedMajor      ProtocolVersionComparison = -4
	AheadMinor         ProtocolVersionComparison = 3
	OutdatedMinor      ProtocolVersionComparison = -3
	AheadPatch         ProtocolVersionComparison = 2
	OutdatedPatch      ProtocolVersionComparison = -2
	AheadPrerelease    ProtocolVersionComparison = 1
	OutdatedPrerelease ProtocolVersionComparison = -1
	Matching           ProtocolVersionComparison = 0
	DiffVersionType    ProtocolVersionComparison = 100
	DiffBuild          ProtocolVersionComparison = 101
	EmptyVersion       ProtocolVersionComparison = 102
)

func (p ProtocolVersion) Compare(other ProtocolVersion) (cmp ProtocolVersionComparison) {
	if p == (ProtocolVersion{}) || (other == (ProtocolVersion{})) {
		return EmptyVersion
	}
	aVersionType, aBuild, aMajor, aMinor, aPatch, aPreRelease := p.Parse()
	bVersionType, bBuild, bMajor, bMinor, bPatch, bPreRelease := other.Parse()
	if aVersionType != bVersionType {
		return DiffVersionType
	}
	if aBuild != bBuild {
		return DiffBuild
	}
	fn := func(a, b uint32, ahead, outdated ProtocolVersionComparison) ProtocolVersionComparison {
		if a == b {
			return Matching
		}
		if a > b {
			return ahead
		}
		return outdated
	}
	if c := fn(aMajor, bMajor, AheadMajor, OutdatedMajor); c != Matching {
		return c
	}
	if c := fn(aMinor, bMinor, AheadMinor, OutdatedMinor); c != Matching {
		return c
	}
	if c := fn(aPatch, bPatch, AheadPatch, OutdatedPatch); c != Matching {
		return c
	}
	return fn(aPreRelease, bPreRelease, AheadPrerelease, OutdatedPrerelease)
}

func ToProtocolVersion(build uint64, major, minor, patch, preRelease uint32) (out ProtocolVersion) {
	binary.BigEndian.PutUint64(out[8:16], build)
	binary.BigEndian.PutUint32(out[16:20], major)
	binary.BigEndian.PutUint32(out[20:24], minor)
	binary.BigEndian.PutUint32(out[24:28], patch)
	binary.BigEndian.PutUint32(out[28:32], preRelease)
	return
}
