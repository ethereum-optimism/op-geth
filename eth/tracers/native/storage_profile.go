package native

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/tracing"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/params"
)

func init() {
	tracers.DefaultDirectory.Register("storageProfileTracer", newStorageProfileTracer, false)
}

// storageProfileResult is the per-transaction result returned by the tracer.
type storageProfileResult struct {
	TxHash          common.Hash `json:"txHash"`
	From            common.Address `json:"from"`
	To              *common.Address `json:"to"`
	GasUsed         uint64 `json:"gasUsed"`
	OPGasRefund     uint64 `json:"opGasRefund"`
	EffectiveGas    uint64 `json:"effectiveGas"`
	RefundRatio     float64 `json:"refundRatio"`
	SstoreCount     uint64 `json:"sstoreCount"`
	SstoreGas       uint64 `json:"sstoreGas"`
	SstoreRatio     float64 `json:"sstoreRatio"`
	StorageHeavy    bool `json:"storageHeavy"`
	WallClockMicros int64 `json:"wallClockMicros"`
	CalldataLen     int `json:"calldataLen"`
	Status          uint64 `json:"status"`
}

// storageProfileTracer captures SDM/opgas profiling data from transaction execution.
// It reads the non-consensus fields populated by innerExecute() via the receipt.
type storageProfileTracer struct {
	// current tx metadata (set in OnTxStart)
	currentFrom common.Address
	currentTo   *common.Address
	calldataLen int

	// collected results
	results []storageProfileResult
}

func newStorageProfileTracer(ctx *tracers.Context, cfg json.RawMessage, chainConfig *params.ChainConfig) (*tracers.Tracer, error) {
	t := &storageProfileTracer{}
	return &tracers.Tracer{
		Hooks: &tracing.Hooks{
			OnTxStart: t.OnTxStart,
			OnTxEnd:   t.OnTxEnd,
		},
		GetResult: t.GetResult,
		Stop:      t.Stop,
	}, nil
}

func (t *storageProfileTracer) OnTxStart(env *tracing.VMContext, tx *types.Transaction, from common.Address) {
	t.currentFrom = from
	t.currentTo = tx.To()
	t.calldataLen = len(tx.Data())
}

func (t *storageProfileTracer) OnTxEnd(receipt *types.Receipt, err error) {
	if err != nil || receipt == nil {
		return
	}

	gasUsed := receipt.GasUsed
	var opGasRefund uint64
	if receipt.OPGasRefund != nil {
		opGasRefund = *receipt.OPGasRefund
	}

	// canonical gas = gasUsed + opGasRefund (receipt.GasUsed already has refund subtracted)
	canonicalGas := gasUsed + opGasRefund
	var refundRatio float64
	if canonicalGas > 0 {
		refundRatio = float64(opGasRefund) / float64(canonicalGas)
	}

	var sstoreRatio float64
	if canonicalGas > 0 {
		sstoreRatio = float64(receipt.SstoreGas) / float64(canonicalGas)
	}

	t.results = append(t.results, storageProfileResult{
		TxHash:          receipt.TxHash,
		From:            t.currentFrom,
		To:              t.currentTo,
		GasUsed:         gasUsed,
		OPGasRefund:     opGasRefund,
		EffectiveGas:    gasUsed,
		RefundRatio:     refundRatio,
		SstoreCount:     receipt.SstoreCount,
		SstoreGas:       receipt.SstoreGas,
		SstoreRatio:     sstoreRatio,
		StorageHeavy:    receipt.StorageHeavy,
		WallClockMicros: receipt.WallClockMicros,
		CalldataLen:     t.calldataLen,
		Status:          receipt.Status,
	})
}

func (t *storageProfileTracer) GetResult() (json.RawMessage, error) {
	return json.Marshal(t.results)
}

func (t *storageProfileTracer) Stop(err error) {
}
