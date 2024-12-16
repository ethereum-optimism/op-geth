package types

type BlockConfig struct {
	CustomWithdrawalsRoot bool
}

func (bc *BlockConfig) HasOptimismWithdrawalsRoot(blockTime uint64) bool {
	return bc.CustomWithdrawalsRoot
}

var (
	DefaultBlockConfig = &BlockConfig{CustomWithdrawalsRoot: false}
	IsthmusBlockConfig = &BlockConfig{CustomWithdrawalsRoot: true}
)
