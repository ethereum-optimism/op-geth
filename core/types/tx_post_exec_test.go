package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/require"
)

func TestTxPostExecUnmarshalJSON(t *testing.T) {
	var tx Transaction
	err := json.Unmarshal([]byte(`{
		"type":"0x7d",
		"gas":"0x0",
		"value":"0x0",
		"input":"0xc201c0"
	}`), &tx)
	require.NoError(t, err)
	require.Equal(t, uint8(TxPostExecType), tx.Type())

	inner, ok := tx.inner.(*TxPostExec)
	require.True(t, ok)
	require.Equal(t, hexutil.MustDecode("0xc201c0"), inner.Data)
}

func TestTxPostExecUnmarshalJSONWithRPCMetadata(t *testing.T) {
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
	require.Equal(t, uint8(TxPostExecType), tx.Type())
}

func TestTxPostExecRoundTrips(t *testing.T) {
	original := NewTx(&TxPostExec{Data: hexutil.MustDecode("0xc201c0")})

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

func TestTxPostExecSenderIsZeroAddress(t *testing.T) {
	tx := NewTx(&TxPostExec{Data: hexutil.MustDecode("0xc201c0")})

	signer := NewLondonSigner(big.NewInt(123))
	sender, err := signer.Sender(tx)
	require.NoError(t, err)
	require.Equal(t, common.Address{}, sender)
}
