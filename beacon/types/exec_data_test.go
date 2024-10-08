package types

import (
	"testing"

	"github.com/ethereum/go-ethereum/beacon/engine"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

var (
	withdrawalsHash1 = common.HexToHash("dead")

	isthumusExecData = engine.ExecutableData{
		ParentHash:      common.HexToHash("parent"),
		FeeRecipient:    common.HexToAddress("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		StateRoot:       common.HexToHash("sRoot"),
		ReceiptsRoot:    common.HexToHash("rRoot"),
		LogsBloom:       common.Hex2Bytes("0x376c47978271565f56DEB45495afa69E59c16Ab2"),
		Random:          common.HexToHash("randao"),
		BaseFeePerGas:   hexutil.MustDecodeBig("0x2000000"),
		Transactions:    [][]byte{},
		Withdrawals:     []*types.Withdrawal{},
		WithdrawalsRoot: &withdrawalsHash1,
	}

	executableDataSamples = []engine.ExecutableData{
		isthumusExecData,
	}
)

func TestExecutableDataJSONEncodeDecode(t *testing.T) {
	for i := range executableDataSamples {
		b, err := executableDataSamples[i].MarshalJSON()
		if err != nil {
			t.Fatal("error marshaling executable data to json:", err)
		}
		r := engine.ExecutableData{}
		err = r.UnmarshalJSON(b)
		if err != nil {
			t.Fatal("error unmarshalling executable data from json:", err)
		}
		assert.Equal(t, withdrawalsHash1, *r.WithdrawalsRoot)
	}
}
