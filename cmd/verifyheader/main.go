// Copyright 2025 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// verifyheader is a utility that fetches a block header from an RPC endpoint
// using `cast` (from Foundry) and verifies it against the beacon consensus rules.
//
// Usage:
//
//	go run ./cmd/verifyheader -rpc <rpc-url> -block <block-number-or-hash> [-chain-id <id>]
//
// Examples:
//
//	go run ./cmd/verifyheader -rpc http://localhost:8545 -block 12345678
//	go run ./cmd/verifyheader -rpc http://localhost:8545 -block latest
//	go run ./cmd/verifyheader -rpc https://mainnet.optimism.io -block 148263570 -chain-id 10
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus/beacon"
	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/superchain"
)

func main() {
	rpcURL := flag.String("rpc", "", "RPC endpoint URL (required)")
	blockID := flag.String("block", "latest", "Block number, hash, or tag (latest, finalized, etc.)")
	chainID := flag.Uint64("chain-id", 0, "Chain ID (0 = auto-detect via cast)")
	flag.Parse()

	if *rpcURL == "" {
		fmt.Fprintln(os.Stderr, "Error: -rpc flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Verify cast is available
	if _, err := exec.LookPath("cast"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: 'cast' (Foundry) not found in PATH. Install from https://getfoundry.sh")
		os.Exit(1)
	}

	// Auto-detect chain ID if not specified
	if *chainID == 0 {
		id, err := fetchChainID(*rpcURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error detecting chain ID: %v\n", err)
			os.Exit(1)
		}
		*chainID = id
		fmt.Printf("Detected chain ID: %d\n", *chainID)
	}

	// Resolve block ID to a hex string for the RPC call
	blockParam := resolveBlockParam(*blockID)

	// Fetch the target header
	fmt.Printf("Fetching block %s from %s...\n", *blockID, *rpcURL)
	header, err := fetchHeader(*rpcURL, blockParam)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching header: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Target block: #%s (hash: %s)\n", header.Number, header.Hash().Hex())

	// Fetch the parent header
	parentNum := new(big.Int).Sub(header.Number, big.NewInt(1))
	parentParam := fmt.Sprintf("0x%x", parentNum)
	fmt.Printf("Fetching parent block #%s...\n", parentNum)
	parent, err := fetchHeader(*rpcURL, parentParam)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error fetching parent header: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Parent block: #%s (hash: %s)\n", parent.Number, parent.Hash().Hex())

	// Verify parent hash matches
	if header.ParentHash != parent.Hash() {
		fmt.Fprintf(os.Stderr, "Error: parent hash mismatch!\n  header.ParentHash: %s\n  parent.Hash():     %s\n",
			header.ParentHash.Hex(), parent.Hash().Hex())
		os.Exit(1)
	}

	// Get chain config
	chainConfig, err := getChainConfig(*chainID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Create a mock chain header reader backed by the two headers we fetched
	chain := &mockChainHeaderReader{
		config:  chainConfig,
		headers: map[common.Hash]*types.Header{},
	}
	chain.headers[parent.Hash()] = parent
	chain.headers[header.Hash()] = header

	// Create the beacon consensus engine and verify
	engine := beacon.New(ethash.NewFaker())

	fmt.Println("\nRunning beacon.VerifyHeader...")
	fmt.Println(strings.Repeat("-", 60))

	if err := engine.VerifyHeader(chain, header); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Verification FAILED: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✅ Header verification PASSED")
	fmt.Println(strings.Repeat("-", 60))
	printHeaderSummary(header)
}

// resolveBlockParam converts a user-provided block identifier into the
// hex-encoded format expected by eth_getBlockByNumber.
func resolveBlockParam(blockID string) string {
	// Named tags pass through as-is
	switch blockID {
	case "latest", "finalized", "safe", "pending", "earliest":
		return blockID
	}
	// If it looks like a hex hash (0x + 64 chars), it's a block hash — not supported
	// by eth_getBlockByNumber, so we'd need eth_getBlockByHash. For simplicity,
	// treat everything else as a decimal number and convert to hex.
	if strings.HasPrefix(blockID, "0x") {
		return blockID // already hex
	}
	n := new(big.Int)
	if _, ok := n.SetString(blockID, 10); ok {
		return fmt.Sprintf("0x%x", n)
	}
	return blockID
}

// fetchChainID retrieves the chain ID from the RPC endpoint using cast.
func fetchChainID(rpcURL string) (uint64, error) {
	cmd := exec.Command("cast", "chain-id", "--rpc-url", rpcURL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return 0, fmt.Errorf("cast chain-id failed: %s", strings.TrimSpace(string(out)))
	}
	var id uint64
	if _, err := fmt.Sscanf(strings.TrimSpace(string(out)), "%d", &id); err != nil {
		return 0, fmt.Errorf("failed to parse chain ID: %v", err)
	}
	return id, nil
}

// rpcResponse represents a JSON-RPC response envelope.
type rpcResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// fetchHeader uses `cast rpc eth_getBlockByNumber` to fetch the raw JSON-RPC
// block response and unmarshals it directly into types.Header, which has the
// correct JSON field tags for all header fields.
func fetchHeader(rpcURL, blockParam string) (*types.Header, error) {
	// cast rpc eth_getBlockByNumber <blockParam> false --rpc-url <url>
	// The "false" means don't include full transactions (we only need the header).
	cmd := exec.Command("cast", "rpc", "eth_getBlockByNumber", blockParam, "false", "--rpc-url", rpcURL)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cast rpc failed: %s", strings.TrimSpace(string(out)))
	}

	// cast rpc returns the raw JSON-RPC result directly (not wrapped in {jsonrpc, result, ...})
	raw := strings.TrimSpace(string(out))
	if raw == "null" {
		return nil, fmt.Errorf("block %s not found", blockParam)
	}

	var header types.Header
	if err := json.Unmarshal([]byte(raw), &header); err != nil {
		return nil, fmt.Errorf("failed to unmarshal header: %v", err)
	}
	return &header, nil
}

// printHeaderSummary prints a summary of the verified header.
func printHeaderSummary(h *types.Header) {
	fmt.Printf("\nHeader Summary:\n")
	fmt.Printf("  Block Number:    %s\n", h.Number)
	fmt.Printf("  Hash:            %s\n", h.Hash().Hex())
	fmt.Printf("  Parent Hash:     %s\n", h.ParentHash.Hex())
	fmt.Printf("  Timestamp:       %d\n", h.Time)
	fmt.Printf("  GasLimit:        %d\n", h.GasLimit)
	fmt.Printf("  GasUsed:         %d\n", h.GasUsed)
	fmt.Printf("  Coinbase:        %s\n", h.Coinbase.Hex())
	if h.BaseFee != nil {
		fmt.Printf("  BaseFee:         %s\n", h.BaseFee)
	}
	if h.BlobGasUsed != nil {
		fmt.Printf("  BlobGasUsed:     %d\n", *h.BlobGasUsed)
	}
	if h.ExcessBlobGas != nil {
		fmt.Printf("  ExcessBlobGas:   %d\n", *h.ExcessBlobGas)
	}
	fmt.Printf("  Extra:           %x\n", h.Extra)
}

// getChainConfig returns the chain configuration for a given chain ID.
// It first checks built-in L1 configs, then tries the OP Stack superchain
// registry (which contains correct EIP-1559 parameters and fork times).
func getChainConfig(chainID uint64) (*params.ChainConfig, error) {
	// Check built-in L1 configs first
	switch chainID {
	case params.MainnetChainConfig.ChainID.Uint64():
		return params.MainnetChainConfig, nil
	case params.SepoliaChainConfig.ChainID.Uint64():
		return params.SepoliaChainConfig, nil
	case params.HoodiChainConfig.ChainID.Uint64():
		return params.HoodiChainConfig, nil
	}

	// Try the OP Stack superchain registry
	if ch, ok := superchain.Chains[chainID]; ok {
		fmt.Printf("Using OP Stack superchain config for %s (chain %d)\n", ch.Name, chainID)
		chConfig, err := ch.Config()
		if err != nil {
			return nil, fmt.Errorf("failed to load superchain config: %v", err)
		}
		return params.LoadOPStackChainConfig(chConfig)
	}

	return nil, fmt.Errorf("unsupported chain ID %d: not found in built-in configs or superchain registry", chainID)
}

// mockChainHeaderReader implements consensus.ChainHeaderReader backed by
// a simple in-memory map of headers. This is sufficient for VerifyHeader
// which only needs Config() and GetHeader().
type mockChainHeaderReader struct {
	config  *params.ChainConfig
	headers map[common.Hash]*types.Header
}

func (m *mockChainHeaderReader) Config() *params.ChainConfig {
	return m.config
}

func (m *mockChainHeaderReader) CurrentHeader() *types.Header {
	return nil
}

func (m *mockChainHeaderReader) GetHeader(hash common.Hash, number uint64) *types.Header {
	return m.headers[hash]
}

func (m *mockChainHeaderReader) GetHeaderByNumber(number uint64) *types.Header {
	for _, h := range m.headers {
		if h.Number.Uint64() == number {
			return h
		}
	}
	return nil
}

func (m *mockChainHeaderReader) GetHeaderByHash(hash common.Hash) *types.Header {
	return m.headers[hash]
}
