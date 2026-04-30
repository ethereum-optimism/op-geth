package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestPostExecTxUnmarshalJSON(t *testing.T) {
	var tx Transaction
	err := json.Unmarshal([]byte(`{
		"type":"0x7d",
		"gas":"0x0",
		"value":"0x0",
		"input":"0xc201c0"
	}`), &tx)
	require.NoError(t, err)
	require.Equal(t, uint8(PostExecTxType), tx.Type())

	inner, ok := tx.inner.(*PostExecTx)
	require.True(t, ok)
	require.Equal(t, hexutil.MustDecode("0xc201c0"), inner.Data)
}

func TestPostExecTxUnmarshalJSONWithRPCMetadata(t *testing.T) {
	var tx Transaction
	err := json.Unmarshal([]byte(`{
		"type":"0x7d",
		"from":"0x0000000000000000000000000000000000000000",
		"nonce":"0x0",
		"gasPrice":"0x1",
		"maxFeePerGas":"0x2",
		"maxPriorityFeePerGas":"0x3",
		"input":"0xc201c0",
		"v":"0x0",
		"r":"0x0",
		"s":"0x0"
	}`), &tx)
	require.NoError(t, err)
	require.Equal(t, uint8(PostExecTxType), tx.Type())
}

func TestPostExecTxUnmarshalJSONErrors(t *testing.T) {
	tests := []struct {
		name          string
		json          string
		expectedError string
	}{
		{
			name:          "non-empty accessList",
			json:          `{"type":"0x7d","input":"0xc201c0","accessList":[{"address":"0x0000000000000000000000000000000000000001","storageKeys":[]}]}`,
			expectedError: "unexpected field(s) in post-exec transaction",
		},
		{
			name:          "unepxpected isSystemTx field",
			json:          `{"type":"0x7d","input":"0xc201c0","isSystemTx":true}`,
			expectedError: "unexpected field(s) in post-exec transaction",
		},
		{
			name:          "unexpected to field",
			json:          `{"type":"0x7d","input":"0xc201c0","to":"0x0000000000000000000000000000000000000001"}`,
			expectedError: "unexpected field(s) in post-exec transaction",
		},
		{
			name:          "non-zero from",
			json:          `{"type":"0x7d","input":"0xc201c0","from":"0x0000000000000000000000000000000000000001"}`,
			expectedError: "post-exec transaction from must be zero address or unset",
		},
		{
			name:          "non-zero nonce",
			json:          `{"type":"0x7d","input":"0xc201c0","nonce":"0x1"}`,
			expectedError: "post-exec transaction nonce must be 0 or unset",
		},
		{
			name:          "non-zero value",
			json:          `{"type":"0x7d","input":"0xc201c0","value":"0x1"}`,
			expectedError: "post-exec transaction value must be 0",
		},
		{
			name:          "non-zero gas",
			json:          `{"type":"0x7d","input":"0xc201c0","gas":"0x1"}`,
			expectedError: "post-exec transaction gas must be 0",
		},
		{
			name:          "non-zero v",
			json:          `{"type":"0x7d","input":"0xc201c0","v":"0x1","r":"0x0","s":"0x0"}`,
			expectedError: "post-exec transaction signature must be 0 or unset",
		},
		{
			name:          "non-zero r",
			json:          `{"type":"0x7d","input":"0xc201c0","v":"0x0","r":"0x1","s":"0x0"}`,
			expectedError: "post-exec transaction signature must be 0 or unset",
		},
		{
			name:          "non-zero s",
			json:          `{"type":"0x7d","input":"0xc201c0","v":"0x0","r":"0x0","s":"0x1"}`,
			expectedError: "post-exec transaction signature must be 0 or unset",
		},
		{
			name:          "missing input",
			json:          `{"type":"0x7d"}`,
			expectedError: "missing required field 'input'",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tx Transaction
			err := json.Unmarshal([]byte(test.json), &tx)
			require.ErrorContains(t, err, test.expectedError)
		})
	}
}

func TestPostExecTxRoundTrips(t *testing.T) {
	original := NewTx(&PostExecTx{Data: hexutil.MustDecode("0xc201c0")})

	jsonBytes, err := original.MarshalJSON()
	require.NoError(t, err)

	var fromJSON Transaction
	require.NoError(t, fromJSON.UnmarshalJSON(jsonBytes))
	require.Equal(t, original.Type(), fromJSON.Type())
	require.Equal(t, original.Hash(), fromJSON.Hash())
	require.Equal(t, original.Data(), fromJSON.Data())

	bin, err := original.MarshalBinary()
	require.NoError(t, err)
	require.Equal(t, hexutil.MustDecode("0x7dc201c0"), bin)

	var fromBinary Transaction
	require.NoError(t, fromBinary.UnmarshalBinary(bin))
	require.Equal(t, original.Type(), fromBinary.Type())
	require.Equal(t, original.Hash(), fromBinary.Hash())
	require.Equal(t, original.Data(), fromBinary.Data())
}

func TestPostExecTxSenderIsZeroAddress(t *testing.T) {
	tx := NewTx(&PostExecTx{Data: hexutil.MustDecode("0xc201c0")})

	signer := NewLondonSigner(big.NewInt(123))
	sender, err := signer.Sender(tx)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, sender)
}
