package core

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/metrics"
)

// OPStack additions
var (
	headBaseFeeGauge     = metrics.NewRegisteredGauge("chain/head/basefee", nil)
	headGasUsedGauge     = metrics.NewRegisteredGauge("chain/head/gas_used", nil)
	headBlobGasUsedGauge = metrics.NewRegisteredGauge("chain/head/blob_gas_used", nil)

	headGasUsedHist     = metrics.NewRegisteredHistogram("chain/head/gas_used_hist", nil, metrics.NewExpDecaySample(1028, 0.015))
	headBlobGasUsedHist = metrics.NewRegisteredHistogram("chain/head/blob_gas_used_hist", nil, metrics.NewExpDecaySample(1028, 0.015))
)

func updateOptimismBlockMetrics(header *types.Header) error {
	headBaseFeeGauge.TryUpdate(header.BaseFee)
	headGasUsedGauge.Update(int64(header.GasUsed))
	headBlobGasUsedGauge.TryUpdateUint64(header.BlobGasUsed)
	headGasUsedHist.Update(int64(header.GasUsed))
	if header.BlobGasUsed != nil {
		headBlobGasUsedHist.Update(int64(*header.BlobGasUsed))
	}
	return nil
}
