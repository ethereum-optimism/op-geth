package legacypool

func (q *queue) withRollupCostFnProvider(p rollupCostFuncProvider) {
	q.rollupCostFnProvider = p
}
