package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/ethereum/go-ethereum/superchain"
)

const (
	// gasLimit() function selector = keccak256("gasLimit()")[:4]
	gasLimitSelector = "f68016b7"
	// Default Ethereum mainnet RPC (can be overridden via ETH_RPC_URL env var)
	defaultRPC = "https://eth.llamarpc.com"
)

type jsonRPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}

type jsonRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  string          `json:"result,omitempty"`
	Error   *jsonRPCError   `json:"error,omitempty"`
}

type jsonRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type chainResult struct {
	Name     string
	ChainID  uint64
	GasLimit uint64
	Error    string
}

func main() {
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		rpcURL = defaultRPC
	}

	fmt.Printf("Using RPC: %s\n", rpcURL)
	fmt.Printf("Fetching gas limits for mainnet chains from SystemConfigProxy...\n\n")

	var results []chainResult

	// Iterate over all chains in the superchain registry
	for chainID, chain := range superchain.Chains {
		// Only process mainnet chains
		if chain.Network != "mainnet" {
			continue
		}

		config, err := chain.Config()
		if err != nil {
			results = append(results, chainResult{
				Name:    chain.Name,
				ChainID: chainID,
				Error:   fmt.Sprintf("failed to get config: %v", err),
			})
			continue
		}

		if config.Addresses.SystemConfigProxy == nil {
			results = append(results, chainResult{
				Name:    chain.Name,
				ChainID: chainID,
				Error:   "no SystemConfigProxy address",
			})
			continue
		}

		gasLimit, err := getGasLimit(rpcURL, config.Addresses.SystemConfigProxy.Hex())
		if err != nil {
			results = append(results, chainResult{
				Name:    chain.Name,
				ChainID: chainID,
				Error:   err.Error(),
			})
			continue
		}

		results = append(results, chainResult{
			Name:     chain.Name,
			ChainID:  chainID,
			GasLimit: gasLimit,
		})
	}

	// Sort results by chain name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	// Print results
	fmt.Printf("%-30s %-12s %s\n", "Chain Name", "Chain ID", "Gas Limit")
	fmt.Printf("%-30s %-12s %s\n", "----------", "--------", "---------")

	for _, r := range results {
		if r.Error != "" {
			fmt.Printf("%-30s %-12d ERROR: %s\n", r.Name, r.ChainID, r.Error)
		} else {
			fmt.Printf("%-30s %-12d %d\n", r.Name, r.ChainID, r.GasLimit)
		}
	}

	fmt.Printf("\nTotal mainnet chains: %d\n", len(results))
}

func getGasLimit(rpcURL, contractAddr string) (uint64, error) {
	callData := map[string]string{
		"to":   contractAddr,
		"data": "0x" + gasLimitSelector,
	}

	req := jsonRPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_call",
		Params:  []interface{}{callData, "latest"},
		ID:      1,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return 0, fmt.Errorf("marshal request: %v", err)
	}

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Post(rpcURL, "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return 0, fmt.Errorf("http post: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("read response: %v", err)
	}

	var rpcResp jsonRPCResponse
	if err := json.Unmarshal(body, &rpcResp); err != nil {
		return 0, fmt.Errorf("unmarshal response: %v", err)
	}

	if rpcResp.Error != nil {
		return 0, fmt.Errorf("rpc error: %s", rpcResp.Error.Message)
	}

	// Remove 0x prefix and decode hex
	hexStr := rpcResp.Result
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// Decode the hex string to bytes
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return 0, fmt.Errorf("decode hex: %v", err)
	}

	// gasLimit is returned as uint64 (right-padded to 32 bytes)
	gasLimit := new(big.Int).SetBytes(data)
	return gasLimit.Uint64(), nil
}
