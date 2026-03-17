package ethclient_test

import (
	"context"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"testing"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/internal/ethapi"
	"github.com/ethereum/go-ethereum/internal/ethapi/override"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/rpc"
)

func testEstimateGas(t *testing.T, client *rpc.Client) {
	ec := ethclient.NewClient(client)

	// EstimateGas
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	gas, err := ec.EstimateGas(context.Background(), msg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gas != 21000 {
		t.Fatalf("unexpected gas price: %v", gas)
	}
}

func testHistoricalRPC(t *testing.T, client *rpc.Client) {
	ec := ethclient.NewClient(client)

	// Estimate Gas RPC
	msg := ethereum.CallMsg{
		From:  testAddr,
		To:    &common.Address{},
		Gas:   21000,
		Value: big.NewInt(1),
	}
	var res hexutil.Uint64
	callMsg := map[string]interface{}{
		"from":  msg.From,
		"to":    msg.To,
		"gas":   hexutil.Uint64(msg.Gas),
		"value": (*hexutil.Big)(msg.Value),
	}
	err := client.CallContext(context.Background(), &res, "eth_estimateGas", callMsg, rpc.BlockNumberOrHashWithNumber(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res != 12345 {
		t.Fatalf("invalid result: %d", res)
	}

	// Call Contract RPC
	histVal, err := ec.CallContract(context.Background(), msg, big.NewInt(1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(histVal) != "test" {
		t.Fatalf("expected %s to equal test", string(histVal))
	}
}

type mockHistoricalBackend struct{}

func (m *mockHistoricalBackend) Call(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash rpc.BlockNumberOrHash, overrides *override.StateOverride) (hexutil.Bytes, error) {
	num, ok := blockNrOrHash.Number()
	if ok && num == 1 {
		return hexutil.Bytes("test"), nil
	}
	return nil, ethereum.NotFound
}

func (m *mockHistoricalBackend) EstimateGas(ctx context.Context, args ethapi.TransactionArgs, blockNrOrHash *rpc.BlockNumberOrHash) (hexutil.Uint64, error) {
	num, ok := blockNrOrHash.Number()
	if ok && num == 1 {
		return hexutil.Uint64(12345), nil
	}
	return 0, ethereum.NotFound
}

func newMockHistoricalBackend(t *testing.T) string {
	s := rpc.NewServer()
	err := node.RegisterApis([]rpc.API{
		{
			Namespace:     "eth",
			Service:       new(mockHistoricalBackend),
			Public:        true,
			Authenticated: false,
		},
	}, nil, s)
	if err != nil {
		t.Fatalf("error creating mock historical backend: %v", err)
	}

	hdlr := node.NewHTTPHandlerStack(s, []string{"*"}, []string{"*"}, nil)
	mux := http.NewServeMux()
	mux.Handle("/", hdlr)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("error creating mock historical backend listener: %v", err)
	}

	go func() {
		httpS := &http.Server{Handler: mux}
		httpS.Serve(listener)

		t.Cleanup(func() {
			httpS.Shutdown(context.Background())
		})
	}()

	return fmt.Sprintf("http://%s", listener.Addr().String())
}
