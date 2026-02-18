package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/params"
)

// deriveOPStackFields derives the OP Stack specific fields for each receipt.
// It must only be called for blocks with at least one transaction (the L1 attributes deposit).
func (rs Receipts) deriveOPStackFields(config *params.ChainConfig, blockTime uint64, txs []*Transaction) error {
	// Exit early if there are only deposit transactions, for which no fields are derived.
	if txs[len(txs)-1].IsDepositTx() {
		return nil
	}

	l1AttributesData := txs[0].Data()
	gasParams, err := extractL1GasParams(config, blockTime, l1AttributesData)
	if err != nil {
		return fmt.Errorf("failed to extract L1 gas params: %w", err)
	}

	var daFootprintGasScalar uint64
	isJovian := config.IsJovian(blockTime)
	if isJovian {
		scalar, err := ExtractDAFootprintGasScalar(l1AttributesData)
		if err != nil {
			return fmt.Errorf("failed to extract DA footprint gas scalar: %w", err)
		}
		daFootprintGasScalar = uint64(scalar)
	}

	for i := range rs {
		if txs[i].IsDepositTx() {
			continue
		}
		rs[i].L1GasPrice = gasParams.l1BaseFee
		rs[i].L1BlobBaseFee = gasParams.l1BlobBaseFee
		rcd := txs[i].RollupCostData()
		rs[i].L1Fee, rs[i].L1GasUsed = gasParams.costFunc(rcd)
		rs[i].FeeScalar = gasParams.feeScalar
		rs[i].L1BaseFeeScalar = u32ptrTou64ptr(gasParams.l1BaseFeeScalar)
		rs[i].L1BlobBaseFeeScalar = u32ptrTou64ptr(gasParams.l1BlobBaseFeeScalar)
		if gasParams.operatorFeeScalar != nil && gasParams.operatorFeeConstant != nil && (*gasParams.operatorFeeScalar != 0 || *gasParams.operatorFeeConstant != 0) {
			rs[i].OperatorFeeScalar = u32ptrTou64ptr(gasParams.operatorFeeScalar)
			rs[i].OperatorFeeConstant = gasParams.operatorFeeConstant
		}
		if isJovian {
			rs[i].DAFootprintGasScalar = &daFootprintGasScalar
			rs[i].BlobGasUsed = daFootprintGasScalar * rcd.EstimatedDASize().Uint64()
		}
	}
	return nil
}
