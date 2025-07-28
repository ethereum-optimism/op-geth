package types

type BlockConfig struct {
	IsIsthmusEnabled bool
}

func (bc *BlockConfig) HasOptimismWithdrawalsRoot(blockTime uint64) bool {
	return bc.IsIsthmusEnabled
}

func (bc *BlockConfig) IsIsthmus(blockTime uint64) bool {
	return bc.IsIsthmusEnabled
}

var (
	DefaultBlockConfig = &BlockConfig{IsIsthmusEnabled: false}
	IsthmusBlockConfig = &BlockConfig{IsIsthmusEnabled: true}
)
