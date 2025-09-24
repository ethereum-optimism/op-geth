package params

// OPStack diff
// This file contains ephemeral feature toggles which should be removed
// after the fork scope is locked.

func (c *ChainConfig) IsMinBaseFee(time uint64) bool {
	return c.IsJovian(time) // Replace with return false to disable
}

func (c *ChainConfig) IsDAFootprintBlockLimit(time uint64) bool {
	return c.IsJovian(time) // Replace with return false to disable
}

func (c *ChainConfig) IsOperatorFeeFix(time uint64) bool {
	return c.IsJovian(time) // Replace with return false to disable
}
