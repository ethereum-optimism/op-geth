package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
)

const maxConcurrent = 1_000

func main() {
	var (
		rpcURL     = flag.String("rpc", "", "RPC URL")
		startBlock = flag.Uint64("start", 0, "Start block number")
		endBlock   = flag.Uint64("end", 105_235_063, "End block number")
	)
	flag.Parse()

	if *rpcURL == "" {
		log.Fatal("rpc url is required")
	}
	if *endBlock < *startBlock {
		log.Fatal("end block must be >= start block")
	}

	client, err := rpc.Dial(*rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to rpc: %v", err)
	}
	defer client.Close()

	var wg sync.WaitGroup
	sem := make(chan struct{}, maxConcurrent)

	for blockNum := *startBlock; blockNum <= *endBlock; blockNum++ {
		wg.Add(1)
		sem <- struct{}{} // acquire semaphore

		go func(blockNum uint64) {
			defer wg.Done()
			defer func() { <-sem }() // release semaphore

			fmt.Fprintf(os.Stderr, "attempting to decode block %d\n", blockNum)

			var result hexutil.Bytes
			err := client.CallContext(context.Background(), &result, "debug_dbAncient", "receipts", blockNum)
			if err != nil {
				fmt.Fprintf(os.Stderr, "block %d: error fetching receipt: %v\n", blockNum, err)
				return
			}

			// Decode the receipts for this block
			receipts := make([]*types.ReceiptForStorage, 0)
			if err := rlp.DecodeBytes(result, &receipts); err != nil {
				x := new(types.LegacyReceiptError)
				if errors.As(err, &x) {
					fmt.Printf("found one! %d\n", blockNum)
				} else {
					// Try decoding as a single receipt
					var receipt types.ReceiptForStorage
					if err := rlp.Decode(bytes.NewReader(result), &receipt); err != nil {
						fmt.Fprintf(os.Stderr, "block %d: error decoding receipt: %v\n", blockNum, err)
						return
					}
				}
				return
			}
		}(blockNum)
	}

	wg.Wait()
}
