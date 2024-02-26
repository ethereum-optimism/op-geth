package policies

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/txpool"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	inboxAddress = common.HexToAddress("0x420")

	_ txpool.OptimismTxPoolPolicy = &SuperchainMessagingPolicy{}
)

type SuperchainMessagingPolicy struct {
	cfg     *params.ChainConfig
	chain   txpool.BlockChain
	backend *rpc.Client
}

func NewSuperchainMessagingPolicy(cfg *params.ChainConfig, chain txpool.BlockChain, backend *rpc.Client) *SuperchainMessagingPolicy {
	return &SuperchainMessagingPolicy{cfg, chain, backend}
}

func (m *SuperchainMessagingPolicy) ValidateTx(tx *types.Transaction) (txpool.OptimismTxPolicyStatus, error) {
	time := m.chain.CurrentBlock().Time
	if !m.cfg.IsInterop(time) || tx.To() == nil || *tx.To() != inboxAddress {
		return txpool.OptimismTxPolicyValid, nil
	}

	msgId, msgBytes, err := unpackInboxExecutionMessageTxData(tx.Data())
	if err != nil {
		return txpool.OptimismTxPolicyInvalid, fmt.Errorf("unable to unpack executeMessage tx data: %w", err)
	}
	msgIdBytes, err := json.Marshal(msgId)
	if err != nil {
		return txpool.OptimismTxPolicyInvalid, fmt.Errorf("unable to marshal message identifier: %w", err)
	}

	if msgId.Timestamp > time {
		return txpool.OptimismTxPolicyInvalid, fmt.Errorf("msg identifier has a future time to the chain")
	}

	var safetyLabel messageSafetyLabel
	if err := m.backend.CallContext(context.TODO(), &safetyLabel, "superchain_messageSafety", msgIdBytes, msgBytes); err != nil {
		return txpool.OptimismTxPolicyInvalid, fmt.Errorf("failed to query message safety: %w", err)
	}

	if safetyLabel == finalized {
		return txpool.OptimismTxPolicyValid, nil
	}

	return txpool.OptimismTxPolicyInvalid, nil
}
