package params

const (
	// An inclusive limit on the max cost for the conditional attached to a tx
	TransactionConditionalMaxCost = 1000

	TransactionConditionalRejectedErrCode        = -32003
	TransactionConditionalCostExceededMaxErrCode = -32005
)
