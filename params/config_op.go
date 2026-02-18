package params

import (
	"fmt"
	"math/big"
)

func (c *ChainConfig) opCheckCompatible(newcfg *ChainConfig, headNumber *big.Int, headTimestamp uint64, genesisTimestamp *uint64) *ConfigCompatError {
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

func (c *ChainConfig) opDescription() string {
	var banner string
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
	banner += "\nAll op fork specifications can be found at https://specs.optimism.io/\n"
	return banner
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
