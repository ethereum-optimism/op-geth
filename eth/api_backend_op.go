package eth

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

func (b *EthAPIBackend) HistoricalRPCService() *rpc.Client {
	return b.eth.historicalRPCService
}

func (b *EthAPIBackend) Genesis() *types.Block {
	return b.eth.blockchain.Genesis()
}
