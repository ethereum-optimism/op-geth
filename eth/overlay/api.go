package overlay

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/eth/filters"
	"github.com/ethereum/go-ethereum/eth/tracers"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

const defaultTraceReexec = uint64(128)

var (
	errExceedMaxTopics = errors.New("exceed max topics")
)

// The maximum number of topic criteria allowed, vm.LOG4 - vm.LOG0
const maxTopics = 4

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error)
	HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error)
	BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error)
	BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error)
	GetTransaction(ctx context.Context, txHash common.Hash) (bool, *types.Transaction, common.Hash, uint64, uint64, error)
	RPCGasCap() uint64
	ChainConfig() *params.ChainConfig
	Engine() consensus.Engine
	ChainDb() ethdb.Database
	StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (*state.StateDB, tracers.StateReleaseFunc, error)
	StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (*core.Message, vm.BlockContext, *state.StateDB, tracers.StateReleaseFunc, error)
	HistoricalRPCService() *rpc.Client
	GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error)
}

type ContractCreation struct {
	Creator common.Address
	TxIndex int
	Block   *types.Block
}

// API is the collection of tracing APIs exposed over the private debugging endpoint.
type API struct {
	backend Backend
}

// NewAPI creates a new API definition for the tracing methods of the Ethereum service.
func NewAPI(backend Backend) *API {
	return &API{backend: backend}
}

// APIs return the collection of RPC services the tracer package offers.
func APIs(backend Backend) []rpc.API {
	// Append all the local APIs and return
	return []rpc.API{
		{
			Namespace: "overlay",
			Service:   NewAPI(backend),
		},
	}
}

// chainContext constructs the context reader which is used by the evm for reading
// the necessary chain context.
func (api *API) chainContext(ctx context.Context) core.ChainContext {
	return ethapi.NewChainContext(ctx, api.backend)
}

// / OVERLAYS ///
type CreationCode struct {
	Code *hexutil.Bytes `json:"code"`
}

func (api *API) CallConstructor(ctx context.Context, address common.Address, code *hexutil.Bytes) (*CreationCode, error) {
	defer func(start time.Time) { log.Trace("CallConstructor finished", "runtime", time.Since(start)) }(time.Now())

	contractCreation, err := api.GetContractCreation(ctx, address)
	if err != nil {
		return nil, err
	}
	block := contractCreation.Block
	if block == nil {
		return nil, fmt.Errorf("block #%d not found", block.Number())
	}
	blockBefore, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(block.NumberU64()-1))
	if err != nil {
		return nil, err
	}
	parent := block.Header()
	if parent == nil {
		return nil, fmt.Errorf("block %d(%x) not found", block.Number(), block.Hash())
	}
	if api.backend.ChainConfig().IsOptimismPreBedrock(block.Number()) {
		return nil, errors.New("l2geth does not have overlay support yet")
	}
	log.Debug("[overlay_callConstructor] found creationBlock", "address", address, "block", block.Number().Int64())

	// get a fresh statedb
	statedb, release, err := api.backend.StateAtBlock(ctx, blockBefore, defaultTraceReexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()

	// first: apply all txs until the contractCreation.TxIndex (excluding)
	replayTransactions := block.Transactions()[:contractCreation.TxIndex]
	signer := types.MakeSigner(api.backend.ChainConfig(), block.Number(), block.Time())
	blockCtx := core.NewEVMBlockContext(block.Header(), api.chainContext(ctx), nil, api.backend.ChainConfig(), statedb)
	vmEvm := vm.NewEVM(blockCtx, vm.TxContext{GasPrice: big.NewInt(0)}, statedb, api.backend.ChainConfig(), vm.Config{})
	for txIdx, tx := range replayTransactions {
		statedb.SetTxContext(tx.Hash(), txIdx)
		msg, err := core.TransactionToMessage(tx, signer, block.BaseFee())
		if err != nil {
			return nil, err
		}
		txContext := core.NewEVMTxContext(msg)
		vmEvm.TxContext = txContext

		_, err = core.ApplyMessage(vmEvm, msg, new(core.GasPool).AddGas(msg.GasLimit))
		if err != nil {
			return nil, err
		}
	}

	// second: apply the overlay tx with the modified bytecode
	creationTx := block.Transactions()[contractCreation.TxIndex]
	log.Debug("[overlay_callConstructor]", "contractCreation.TxIndex", contractCreation.TxIndex, "creationTx", creationTx.Hash().String())
	msg, err := api.NewOverlayMessage(creationTx, signer, block.BaseFee())
	if err != nil {
		return nil, err
	}

	contractAddr := crypto.CreateAddress(msg.From, msg.Nonce)
	if creationTx.To() == nil && contractAddr == address {
		// replace with new code
		msg.Data = *code
	}

	txContext := core.NewEVMTxContext(msg)
	ct := &OverlayCreateTracer{contractAddress: address, code: *code, gasCap: api.backend.RPCGasCap()}
	vmConfig := vm.Config{NoBaseFee: true, Tracer: ct}
	vmEvm.Config = vmConfig
	vmEvm.TxContext = txContext

	statedb.SetTxContext(creationTx.Hash(), contractCreation.TxIndex)
	_, err = core.ApplyMessage(vmEvm, msg, new(core.GasPool).AddGas(msg.GasLimit))

	deployedCode := ct.resultCode
	log.Debug("[overlay_callConstructor]", "deployedCode", deployedCode)

	resultCode := &CreationCode{}
	if deployedCode != nil {
		log.Debug("deployedCode != nil")
		c := hexutil.Bytes(deployedCode)
		resultCode.Code = &c
	} else {
		log.Debug("deployedCode == nil")
		// err from core.ApplyMessage()
		if err != nil {
			return nil, err
		}
		c := hexutil.Bytes(statedb.GetCode(address))
		resultCode.Code = &c
	}
	return resultCode, nil
}

func (api *API) NewOverlayMessage(overlayTx *types.Transaction, signer types.Signer, blockBaseFee *big.Int) (*core.Message, error) {
	msg, err := core.TransactionToMessage(overlayTx, signer, blockBaseFee)
	if err != nil {
		return nil, err
	}
	msg.GasLimit = api.backend.RPCGasCap()
	msg.GasPrice = big.NewInt(0)
	msg.GasFeeCap = big.NewInt(0)
	msg.GasTipCap = big.NewInt(0)
	msg.RollupCostData = types.RollupCostData{}

	return msg, nil
}

type blockReplayTask struct {
	idx         int
	BlockNumber int64
}

type blockReplayResult struct {
	BlockNumber int64        `json:"block_number"`
	Logs        []*types.Log `json:"logs,omitempty"`
	Error       string       `json:"error,omitempty"`
}

// GetLogs returns logs matching the given argument that are stored within the state.
func (api *API) GetLogs(ctx context.Context, crit filters.FilterCriteria, stateOverride *ethapi.StateOverride) ([]*types.Log, error) {
	if len(crit.Topics) > maxTopics {
		return nil, errExceedMaxTopics
	}
	begin := crit.FromBlock.Int64()
	end := crit.ToBlock.Int64()
	if begin > end {
		return nil, errors.New("begin > end")
	}
	numBlocks := end - begin + 1
	var (
		results = make([]*blockReplayResult, numBlocks)
		pend    sync.WaitGroup
	)

	threads := runtime.NumCPU()
	if big.NewInt(int64(threads)).Int64() > numBlocks {
		threads = int(numBlocks)
	}
	jobs := make(chan *blockReplayTask, threads)
	for th := 0; th < threads; th++ {
		pend.Add(1)
		go func() {
			defer pend.Done()
			// Fetch and execute the next transaction trace tasks
			for task := range jobs {
				blockNumber := task.BlockNumber
				if err := ctx.Err(); err != nil {
					results[task.idx] = &blockReplayResult{BlockNumber: task.BlockNumber, Error: err.Error()}
					continue
				}
				block, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNumber))
				if err != nil {
					results[task.idx] = &blockReplayResult{BlockNumber: task.BlockNumber, Error: err.Error()}
					continue
				}
				blockLogs, err := api.replayBlock(ctx, block, -1, common.Address{}, stateOverride)
				if err != nil {
					results[task.idx] = &blockReplayResult{BlockNumber: task.BlockNumber, Error: err.Error()}
					continue
				}
				logs := filterLogs(blockLogs, crit.Addresses, crit.Topics)
				log.Debug("logs", "len(logs)", len(logs), "crit.Addresses", crit.Addresses, "crit.Topics", crit.Topics)
				results[task.idx] = &blockReplayResult{BlockNumber: task.BlockNumber, Logs: logs}
			}
		}()
	}

	var failed error
	idx := 0
blockLoop:
	for blockNumber := begin; blockNumber <= end; blockNumber++ {
		task := &blockReplayTask{idx: idx, BlockNumber: blockNumber}
		select {
		case <-ctx.Done():
			failed = ctx.Err()
			break blockLoop
		case jobs <- task:
		}
		idx++
	}

	close(jobs)
	pend.Wait()

	// If execution failed in between, abort
	if failed != nil {
		return nil, failed
	}

	logs := []*types.Log{}
	for idx := range results {
		res := results[idx]
		log.Debug("logs for", "blockNumber", res.BlockNumber, "len(Logs)", len(res.Logs), "Error", res.Error)
		logs = append(logs, res.Logs...)
	}
	log.Debug("FINAL LOGS", "len(logs)", len(logs))
	return logs, nil
}

// includes returns true if the element is present in the list.
func includes[T comparable](things []T, element T) bool {
	for _, thing := range things {
		if thing == element {
			return true
		}
	}
	return false
}

func filterLogs(logs []*types.Log, addresses []common.Address, topics [][]common.Hash) []*types.Log {
	var check = func(log *types.Log) bool {
		if len(addresses) > 0 && !includes(addresses, log.Address) {
			return false
		}
		// If the to filtered topics is greater than the amount of topics in logs, skip.
		if len(topics) > len(log.Topics) {
			return false
		}
		for i, sub := range topics {
			if len(sub) == 0 {
				continue // empty rule set == wildcard
			}
			if !includes(sub, log.Topics[i]) {
				return false
			}
		}
		return true
	}
	ret := []*types.Log{}
	for _, log := range logs {
		if check(log) {
			ret = append(ret, log)
		}
	}
	return ret
}

func (api *API) replayBlock(ctx context.Context, block *types.Block, maxIndex int, contractAddress common.Address, stateOverride *ethapi.StateOverride) ([]*types.Log, error) {
	log.Debug("[replayBlock]", "block", block.Hash())
	blockLogs := []*types.Log{}

	blockBefore, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(block.NumberU64()-1))
	if err != nil {
		return nil, err
	}
	statedb, release, err := api.backend.StateAtBlock(ctx, blockBefore, defaultTraceReexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()
	if stateOverride != nil {
		err = stateOverride.Apply(statedb)
		if err != nil {
			return nil, err
		}
	}

	signer := types.MakeSigner(api.backend.ChainConfig(), block.Number(), block.Time())
	blockCtx := core.NewEVMBlockContext(block.Header(), api.chainContext(ctx), nil, api.backend.ChainConfig(), statedb)
	replayTransactions := block.Transactions()
	vmEvm := vm.NewEVM(blockCtx, vm.TxContext{GasPrice: big.NewInt(0)}, statedb, api.backend.ChainConfig(), vm.Config{NoBaseFee: true})

	receipts, err := api.backend.GetReceipts(ctx, block.Hash())
	if err != nil {
		return nil, err
	}
	for txIdx, tx := range replayTransactions {
		log.Debug("[replayBlock]", "txIdx", txIdx)
		msg, err := api.NewOverlayMessage(tx, signer, block.BaseFee())
		if err != nil {
			return nil, err
		}
		statedb.SetTxContext(tx.Hash(), txIdx)
		receipt := receipts[uint64(txIdx)]
		if receipt.Status == types.ReceiptStatusFailed {
			log.Debug("[replayBlock] skipping transaction because it has status=failed", "transactionHash", tx.Hash())
			contractCreation := msg.To == nil
			if !contractCreation {
				// bump the nonce of the sender
				sender := vm.AccountRef(msg.From)
				statedb.SetNonce(msg.From, statedb.GetNonce(sender.Address())+1)
				continue
			}
		}

		txContext := core.NewEVMTxContext(msg)
		vmEvm.TxContext = txContext

		res, err := core.ApplyMessage(vmEvm, msg, new(core.GasPool).AddGas(msg.GasLimit))
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		if res.Failed() {
			log.Debug("[replayBlock] res result for transaction", "transactionHash", tx.Hash(), "failed", res.Failed(), "revert", res.Revert(), "error", res.Err)
			log.Debug("[replayBlock] discarding txLogs because tx has status=failed", "transactionHash", tx.Hash())
		} else {
			//append logs only if tx has not reverted
			txLogs := statedb.GetLogs(tx.Hash(), block.NumberU64(), block.Hash())
			log.Debug("[replayBlock]", "len(txLogs)", len(txLogs), "transactionHash", tx.Hash())
			blockLogs = append(blockLogs, txLogs...)
		}
	}
	return blockLogs, nil
}

func (api *API) GetContractCreation(ctx context.Context, address common.Address) (*ContractCreation, error) {
	latestBlock, err := api.backend.BlockByNumber(ctx, rpc.LatestBlockNumber)
	if err != nil {
		return nil, err
	}
	// TODO implement for pre-bedrock in l2geth
	var start int64 = 105235063 // bedrock block
	block, err := api.getCreationBlock(ctx, address, start, latestBlock.Number().Int64())
	if err != nil {
		return nil, err
	}

	// try to recompute the state
	blockBefore, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(block.NumberU64()-1))
	if err != nil {
		return nil, err
	}

	// get a fresh statedb
	statedb, release, err := api.backend.StateAtBlock(ctx, blockBefore, defaultTraceReexec, nil, true, false)
	if err != nil {
		return nil, err
	}
	defer release()
	// try to find the creation tx and index inside the block
	// replay all transactions to find the creationIndex
	signer := types.MakeSigner(api.backend.ChainConfig(), block.Number(), block.Time())
	blockCtx := core.NewEVMBlockContext(block.Header(), api.chainContext(ctx), nil, api.backend.ChainConfig(), statedb)
	replayTransactions := block.Transactions()
	ct := &OverlayCreateTracer{contractAddress: address}
	vmConfig := vm.Config{Tracer: ct}
	vmEvm := vm.NewEVM(blockCtx, vm.TxContext{GasPrice: big.NewInt(0)}, statedb, api.backend.ChainConfig(), vmConfig)

	for txIdx, tx := range replayTransactions {
		statedb.SetTxContext(tx.Hash(), txIdx)
		msg, err := core.TransactionToMessage(tx, signer, block.BaseFee())
		if err != nil {
			return nil, err
		}
		txContext := core.NewEVMTxContext(msg)
		vmEvm.TxContext = txContext

		_, err = core.ApplyMessage(vmEvm, msg, new(core.GasPool).AddGas(msg.GasLimit))
		if err != nil {
			return nil, err
		}

		if ct.foundCreator {
			log.Debug("[replayBlock]", "ct.foundCreator", ct.foundCreator, "ct.creationTxIndex", txIdx)
			cc := &ContractCreation{
				Creator: ct.creator,
				Block:   block,
				TxIndex: txIdx,
			}
			log.Debug("ContractCreation", "creator", cc.Creator, "block", cc.Block, "txIndex", cc.TxIndex)
			return cc, nil
		}
	}
	return nil, fmt.Errorf("error! Could not find creator for contract %s", address.String())
}

func (api *API) getCreationBlock(ctx context.Context, contractAddress common.Address, startBlock int64, endBlock int64) (*types.Block, error) {
	if startBlock == endBlock {
		creationBlock, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(startBlock))
		if err != nil {
			return nil, err
		}
		return creationBlock, nil
	}
	midBlock := (startBlock + endBlock) / 2
	codeLength, err := api.codeLength(contractAddress, midBlock)
	if err != nil {
		return nil, err
	}
	if codeLength > 2 {
		return api.getCreationBlock(ctx, contractAddress, startBlock, midBlock)
	} else {
		return api.getCreationBlock(ctx, contractAddress, midBlock+1, endBlock)
	}
}

func (api *API) codeLength(contractAddress common.Address, blockNum int64) (int, error) {
	ctx := context.Background()
	block, err := api.backend.BlockByNumber(ctx, rpc.BlockNumber(blockNum))
	if err != nil {
		return -1, err
	}
	statedb, release, err := api.backend.StateAtBlock(ctx, block, defaultTraceReexec, nil, true, false)
	if err != nil {
		return -1, err
	}
	defer release()
	data := statedb.GetCode(contractAddress)
	return len(data), nil
}
