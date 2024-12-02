package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/core/types"
)

func decodeEncodeJSON(input []byte, val interface{}) error {
	if err := json.Unmarshal(input, &val); err != nil {
		// not valid JSON, nothing to do
		return nil
	}
	output, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if !bytes.Equal(input, output) {
		return fmt.Errorf("encode-decode is not equal, \ninput : %x\noutput: %x", input, output)
	}
	return nil
}

func FuzzJSON(f *testing.F) {
	f.Fuzz(fuzzJSON)
}

func fuzzJSON(t *testing.T, input []byte) {
	if len(input) == 0 {
		return
	}
	{
		var h types.Header
		if err := decodeEncodeJSON(input, &h); err != nil {
			t.Fatal(err)
		}
		var b types.Block
		if err := decodeEncodeJSON(input, &b); err != nil {
			t.Fatal(err)
		}
		var tx types.Transaction
		if err := decodeEncodeJSON(input, &tx); err != nil {
			t.Fatal(err)
		}
		var txs types.Transactions
		if err := decodeEncodeJSON(input, &txs); err != nil {
			t.Fatal(err)
		}
		var rs types.Receipts
		if err := decodeEncodeJSON(input, &rs); err != nil {
			t.Fatal(err)
		}
		var e engine.ExecutableData
		if err := decodeEncodeJSON(input, &e); err != nil {
			t.Fatal(err)
		}
	}
}
