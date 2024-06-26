package sequencerapi

import (
	"context"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	sendRawTxConditionalCostMeter       = metrics.NewRegisteredMeter("sequencer/sendRawTransactionConditional/cost", nil)
	sendRawTxConditionalRequestsCounter = metrics.NewRegisteredCounter("sequencer/sendRawTransactionConditional/requests", nil)
	sendRawTxConditionalAcceptedCounter = metrics.NewRegisteredCounter("sequencer/sendRawTransactionConditional/accepted", nil)
)

type sendRawTxCond struct {
	b ethapi.Backend
}

func GetSendRawTxConditionalAPI(b ethapi.Backend) rpc.API {
	return rpc.API{
		Namespace: "eth",
		Service:   &sendRawTxCond{b},
	}
}

func (s *sendRawTxCond) SendRawTransactionConditional(ctx context.Context, txBytes hexutil.Bytes, cond types.TransactionConditional) (common.Hash, error) {
	sendRawTxConditionalRequestsCounter.Inc(1)

	cost := cond.Cost()
	sendRawTxConditionalCostMeter.Mark(int64(cost))
	if cost > types.TransactionConditionalMaxCost {
		return common.Hash{}, &jsonRpcError{
			message: fmt.Sprintf("conditional cost, %d, exceeded max: %d", cost, types.TransactionConditionalMaxCost),
			code:    types.TransactionConditionalCostExceededMaxErrCode,
		}
	}

	// Perform sanity validation prior to state lookups
	if err := cond.Validate(); err != nil {
		return common.Hash{}, &jsonRpcError{
			message: fmt.Sprintf("failed conditional validation: %s", err),
			code:    types.TransactionConditionalRejectedErrCode,
		}
	}

	state, header, err := s.b.StateAndHeaderByNumber(context.Background(), rpc.LatestBlockNumber)
	if err != nil {
		return common.Hash{}, err
	}
	if err := header.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &jsonRpcError{
			message: fmt.Sprintf("failed header check: %s", err),
			code:    types.TransactionConditionalRejectedErrCode,
		}
	}
	if err := state.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &jsonRpcError{
			message: fmt.Sprintf("failed state check: %s", err),
			code:    types.TransactionConditionalRejectedErrCode,
		}
	}

	// We also check against the parent block to eliminate the MEV incentive in comparison with sendRawTransaction
	parentBlock := rpc.BlockNumberOrHash{BlockHash: &header.ParentHash}
	parentState, _, err := s.b.StateAndHeaderByNumberOrHash(context.Background(), parentBlock)
	if err != nil {
		return common.Hash{}, err
	}
	if err := parentState.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &jsonRpcError{
			message: fmt.Sprintf("failed parent block %s state check: %s", header.ParentHash, err),
			code:    types.TransactionConditionalRejectedErrCode,
		}
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		return common.Hash{}, err
	}

	// Set internal fields
	cond.SubmissionTime = time.Now()
	cond.Rejected = &atomic.Bool{}

	tx.SetConditional(&cond)
	sendRawTxConditionalAcceptedCounter.Inc(1)

	return ethapi.SubmitTransaction(ctx, s.b, tx)
}
