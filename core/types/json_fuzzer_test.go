package types

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
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
		var h Header
		if err := decodeEncodeJSON(input, &h); err != nil {
			t.Fatal(err)
		}
		var b Block
		if err := decodeEncodeJSON(input, &b); err != nil {
			t.Fatal(err)
		}
		var tx Transaction
		if err := decodeEncodeJSON(input, &tx); err != nil {
			t.Fatal(err)
		}
		var txs Transactions
		if err := decodeEncodeJSON(input, &txs); err != nil {
			t.Fatal(err)
		}
		var rs Receipts
		if err := decodeEncodeJSON(input, &rs); err != nil {
			t.Fatal(err)
		}
	}
}
