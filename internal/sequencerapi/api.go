package sequencerapi

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/time/rate"
)

var (
	sendRawTxConditionalCostMeter       = metrics.NewRegisteredMeter("sequencer/sendRawTransactionConditional/cost", nil)
	sendRawTxConditionalRequestsCounter = metrics.NewRegisteredCounter("sequencer/sendRawTransactionConditional/requests", nil)
	sendRawTxConditionalAcceptedCounter = metrics.NewRegisteredCounter("sequencer/sendRawTransactionConditional/accepted", nil)
)

type sendRawTxCond struct {
	b           ethapi.Backend
	seqRPC      *rpc.Client
	costLimiter *rate.Limiter
}

func GetSendRawTxConditionalAPI(b ethapi.Backend, seqRPC *rpc.Client, costRateLimit rate.Limit) rpc.API {
	// Applying a manual bump to the burst to allow conditional txs to queue. Metrics will
	// will inform of adjustments that may need to be made here.
	costLimiter := rate.NewLimiter(costRateLimit, 3*params.TransactionConditionalMaxCost)
	return rpc.API{
		Namespace: "eth",
		Service:   &sendRawTxCond{b, seqRPC, costLimiter},
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

	state, header, err := s.b.StateAndHeaderByNumber(ctx, rpc.LatestBlockNumber)
	if err != nil {
		return common.Hash{}, err
	}
	if err := header.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("failed header check: %s", err),
			Code:    params.TransactionConditionalRejectedErrCode,
		}
	}
	if err := state.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("failed state check: %s", err),
			Code:    params.TransactionConditionalRejectedErrCode,
		}
	}

	// State is checked against an older block to remove the MEV incentive for this endpoint compared with sendRawTransaction
	parentBlock := rpc.BlockNumberOrHash{BlockHash: &header.ParentHash}
	parentState, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, parentBlock)
	if err != nil {
		return common.Hash{}, err
	}
	if err := parentState.CheckTransactionConditional(&cond); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("failed parent block %s state check: %s", header.ParentHash, err),
			Code:    params.TransactionConditionalRejectedErrCode,
		}
	}

	// enforce rate limit on the cost to be observed
	if err := s.costLimiter.WaitN(ctx, cost); err != nil {
		return common.Hash{}, &rpc.JsonError{
			Message: fmt.Sprintf("cost %d rate limited", cost),
			Code:    params.TransactionConditionalCostExceededMaxErrCode,
		}
	}

	tx := new(types.Transaction)
	if err := tx.UnmarshalBinary(txBytes); err != nil {
		return common.Hash{}, err
	}

	// forward if seqRPC is set, otherwise submit the tx
	if s.seqRPC != nil {
		// Some precondition checks done by `ethapi.SubmitTransaction` that are good to also check here
		if err := ethapi.CheckTxFee(tx.GasPrice(), tx.Gas(), s.b.RPCTxFeeCap()); err != nil {
			return common.Hash{}, err
		}
		if !s.b.UnprotectedAllowed() && !tx.Protected() {
			// Ensure only eip155 signed transactions are submitted if EIP155Required is set.
			return common.Hash{}, errors.New("only replay-protected (EIP-155) transactions allowed over RPC")
		}

		var hash common.Hash
		err := s.seqRPC.CallContext(ctx, &hash, "eth_sendRawTransactionConditional", txBytes, cond)
		return hash, err
	} else {
		// Set out-of-consensus internal tx fields
		tx.SetTime(time.Now())
		tx.SetConditional(&cond)

		// `SubmitTransaction` which forwards to `b.SendTx` also checks if its internal `seqRPC` client is
		// set. Since both of these client are constructed when `RollupSequencerHTTP` is supplied, the above
		// block ensures that we're only adding to the txpool for this node.
		sendRawTxConditionalAcceptedCounter.Inc(1)
		return ethapi.SubmitTransaction(ctx, s.b, tx)
	}
}
