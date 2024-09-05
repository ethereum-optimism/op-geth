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
	"github.com/ethereum/go-ethereum/params"
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
	if cost > params.TransactionConditionalMaxCost {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("conditional cost, %d, exceeded max: %d", cost, params.TransactionConditionalMaxCost),
			Code:    params.TransactionConditionalCostExceededMaxErrCode,
		}
	}

	// Perform sanity validation prior to state lookups
	if err := cond.Validate(); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("failed conditional validation: %s", err),
			Code:    params.TransactionConditionalRejectedErrCode,
		}
	}

	header, err := s.b.HeaderByNumber(context.Background(), rpc.LatestBlockNumber)
	if err != nil {
		return common.Hash{}, err
	}
	if err := header.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("failed header check: %s", err),
			Code:    params.TransactionConditionalRejectedErrCode,
		}
	}

	// State is checked against an older block to remove the MEV incentive for this endpoint compared with sendRawTransaction
	parentBlock := rpc.BlockNumberOrHash{BlockHash: &header.ParentHash}
	parentState, _, err := s.b.StateAndHeaderByNumberOrHash(context.Background(), parentBlock)
	if err != nil {
		return common.Hash{}, err
	}
	if err := parentState.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("failed parent block %s state check: %s", header.ParentHash, err),
			Code:    params.TransactionConditionalRejectedErrCode,
		}
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		return common.Hash{}, err
	}

	// Set internal fields
	tx.SetTime(time.Now())
	cond.Rejected = &atomic.Bool{}

	tx.SetConditional(&cond)
	sendRawTxConditionalAcceptedCounter.Inc(1)

	return ethapi.SubmitTransaction(ctx, s.b, tx)
}
