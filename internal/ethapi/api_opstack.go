package ethapi

import (
	"context"
	"fmt"

	"github.com/holiman/uint256"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

type OPStackBackend interface {
	Backend // TODO we can narrow down the backend interface, instead of relying on the full backend
}

// OPStackAPI provides an API to access OP-Stack specific data
type OPStackAPI struct {
	backend OPStackBackend
}

// NewOPStackAPI creates a new Ethereum protocol API for full nodes.
func NewOPStackAPI(backend OPStackBackend) *OPStackAPI {
	return &OPStackAPI{backend: backend}
}

type InteropMessage struct {
	Nonce uint256.Int

	SourceChain common.Hash
	TargetChain common.Hash

	From common.Address
	To   common.Address

	Value uint256.Int

	// TODO limit this to u64 in contract
	GasLimit hexutil.Uint64

	Data hexutil.Bytes
}

func (msg *InteropMessage) Root() common.Hash {
	return common.Hash{}
}

type InteropSourceInfo struct {
	BlockHash              common.Hash
	StateHash              common.Hash
	BlockNumber            hexutil.Uint64
	Timestamp              hexutil.Uint64
	WithdrawalsStorageHash common.Hash
}

// The L2 output root (of different versions) can be computed from the above data.

type InteropMessages struct {
	Info     *InteropSourceInfo
	Messages []*InteropMessage
}

var crossL2OutboxAddr = common.HexToAddress("")

func (api *OPStackAPI) InteropMessages(ctx context.Context, blockNrOrHash rpc.BlockNumberOrHash) (*InteropMessages, error) {
	header, err := api.backend.HeaderByNumberOrHash(ctx, blockNrOrHash)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve header of %s: %v", blockNrOrHash.String(), err)
	}
	info := &InteropSourceInfo{
		BlockHash:              header.Hash(),
		StateHash:              header.Root,
		BlockNumber:            hexutil.Uint64(header.Number.Uint64()),
		Timestamp:              hexutil.Uint64(header.Time),
		WithdrawalsStorageHash: common.Hash{}, // TODO need a backend method for state storage read
	}
	// warning: the retrieved logs raw, not hydrated, i.e. we do not spend time on adding metadata like hashes etc. to them.
	blockLogs, err := api.backend.GetLogs(ctx, header.Hash(), header.Number.Uint64())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs of %s (%d): %v", header.Hash(), header.Number.Uint64(), err)
	}
	var messages []*InteropMessage
	for i, txLogs := range blockLogs {
		for j, log := range txLogs {
			msg, err := logToInteropMessageMaybe(log)
			if err != nil {
				return nil, fmt.Errorf("failed to parse interop log of block %s, tx %d, relative index %d: %w", header.Hash(), i, j, err)
			}
			if msg != nil {
				messages = append(messages, msg)
			}
		}
	}
	return &InteropMessages{
		Info:     info,
		Messages: messages,
	}, nil
}

// logToInteropMessageMaybe parses log data into an interop message.
// Nil, nil is returned if it's not an interop log.
// An error is returned if it's an invalid interop log
func logToInteropMessageMaybe(log *types.Log) (*InteropMessage, error) {
	if log.Address != crossL2OutboxAddr {
		return nil, nil
	}
	// TODO check log topic
	// TODO parse log contents into message
	return nil, nil
}
