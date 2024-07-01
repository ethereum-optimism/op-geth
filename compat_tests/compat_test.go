package compat_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"
)

type blockTransactions struct {
	Transactions []*types.Transaction `json:"transactions"`
}

type blockHash struct {
	Hash common.Hash `json:"hash"`
}

func TestCompatibilityOfChain(t *testing.T) {
	dumpOutput := false
	failOnFirstMismatch := true
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// celo-blockchain
	c1, err := rpc.DialContext(ctx, "https://alfajores-forno.celo-testnet.org")

	// op-geth
	c2, err := rpc.DialContext(ctx, "http://localhost:8545")
	require.NoError(t, err)
	ec := ethclient.NewClient(c2)
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	id2, err := ec.ChainID(ctx)
	require.NoError(t, err)
	require.Greater(t, id2.Uint64(), uint64(0))

	latestBlock, err := ec.BlockNumber(ctx)
	require.NoError(t, err)
	var res json.RawMessage
	startBlock := uint64(0)

	// We subtract 128 from the amount to avoid handlig blocks where state is
	// present since when state is present baseFeePerGas is set on the celo
	// block with a value, and we can't access that state from the op-geth side
	amount := latestBlock - 128
	incrementalLogs := make([]*types.Log, 0)
	for i := startBlock; i <= startBlock+amount; i++ {
		ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		b, err := ec.BlockByNumber(ctx, big.NewInt(int64(i)))
		require.NoError(t, err)
		h, err := ec.HeaderByNumber(ctx, big.NewInt(int64(i)))
		require.NoError(t, err)

		// Comparing the 2 headers directly doesn't work, they differ on unset
		// and nil big ints, which represent the same thing but are technically
		// not eaual. So instead we compare the marshalled output.
		bhMarshalled, err := json.Marshal(b.Header())
		require.NoError(t, err)
		hMarshalled, err := json.Marshal(h)
		require.NoError(t, err)
		require.Equal(t, bhMarshalled, hMarshalled)

		// Now get the header by hash and compare
		h2, err := ec.HeaderByHash(ctx, h.Hash())
		require.NoError(t, err)
		h2Marshalled, err := json.Marshal(h2)
		require.NoError(t, err)
		require.Equal(t, hMarshalled, h2Marshalled)

		// Now get the block by hash and compare
		b2, err := ec.BlockByHash(ctx, b.Hash())
		require.NoError(t, err)
		b2hMarshalled, err := json.Marshal(b2.Header())
		require.NoError(t, err)
		require.Equal(t, bhMarshalled, b2hMarshalled)

		bbMarshalled, err := json.Marshal(b.Body())
		require.NoError(t, err)
		b2bMarshalled, err := json.Marshal(b2.Body())
		require.NoError(t, err)
		require.Equal(t, bbMarshalled, b2bMarshalled)

		res = rpcCallCompare(t, c1, c2, dumpOutput, failOnFirstMismatch, "eth_getBlockByNumber", hexutil.EncodeUint64(i), true)
		// Check we got a block
		require.NotEqual(t, "null", string(res), "block %d should not be null", i)
		blockHash := blockHash{}
		err = json.Unmarshal(res, &blockHash)
		require.NoError(t, err)
		txs := blockTransactions{}
		err = json.Unmarshal(res, &txs)
		require.NoError(t, err)

		incrementalBlockReceipts := types.Receipts{}
		for i, tx := range txs.Transactions {
			// Compare transactions decoded from rpcCall with transactions from
			// the the block returned by ethclient and those directly retrieved
			// by ethclient. Comparing transactions directly does not work
			// because of differing private fields. The hash is calculated over
			// the data of the transaction and so serves as a good proxy for
			// equality.
			require.Equal(t, tx.Hash(), b.Transactions()[i].Hash())
			tx2, _, err := ec.TransactionByHash(ctx, tx.Hash())
			require.NoError(t, err)
			require.Equal(t, tx.Hash(), tx2.Hash())

			_ = rpcCallCompare(t, c1, c2, dumpOutput, failOnFirstMismatch, "eth_getTransactionByHash", tx.Hash())
			require.NoError(t, err)
			res = rpcCallCompare(t, c1, c2, dumpOutput, failOnFirstMismatch, "eth_getTransactionReceipt", tx.Hash())
			require.NoError(t, err)
			r := types.Receipt{}
			err = json.Unmarshal(res, &r)
			require.NoError(t, err)

			// Check receipt decoded from RPC call matches that returned directly from ethclient.
			r2, err := ec.TransactionReceipt(ctx, tx.Hash())
			require.NoError(t, err)
			require.Equal(t, r, *r2)

			incrementalBlockReceipts = append(incrementalBlockReceipts, &r)
			incrementalLogs = append(incrementalLogs, r.Logs...)
		}
		// Get the Celo block receipt. See https://docs.celo.org/developer/migrate/from-ethereum#core-contract-calls
		res = rpcCallCompare(t, c1, c2, dumpOutput, failOnFirstMismatch, "eth_getBlockReceipt", blockHash.Hash)
		if string(res) != "null" {
			r := types.Receipt{}
			err = json.Unmarshal(res, &r)
			require.NoError(t, err)
			if len(r.Logs) > 0 {
				// eth_getBlockReceipt generates an empty receipt when there
				// are no logs, we want to avoid adding these here since the
				// same is not done in eth_gethBlockReceipts, the output of
				// which we will later compare against.
				incrementalBlockReceipts = append(incrementalBlockReceipts, &r)
			}
			incrementalLogs = append(incrementalLogs, r.Logs...)
		}

		blockReceipts := types.Receipts{}
		res = rpcCallCompare(t, c1, c2, dumpOutput, failOnFirstMismatch, "eth_getBlockReceipts", hexutil.EncodeUint64(i))
		err = json.Unmarshal(res, &blockReceipts)
		require.NoError(t, err)
		require.Equal(t, incrementalBlockReceipts, blockReceipts)
	}

	// Get all logs for the range and compare with the logs extracted from receipts.
	from := rpc.BlockNumber(startBlock)
	to := rpc.BlockNumber(amount + startBlock)
	res = rpcCallCompare(t, c1, c2, dumpOutput, failOnFirstMismatch, "eth_getLogs", filterQuery{
		FromBlock: &from,
		ToBlock:   &to,
	})
	var logs []*types.Log
	err = json.Unmarshal(res, &logs)
	require.NoError(t, err)
	require.Equal(t, len(incrementalLogs), len(logs))
	require.Equal(t, incrementalLogs, logs)

	// Compare logs with those retrived via ethclient
	logs2, err := ec.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: big.NewInt(int64(startBlock)),
		ToBlock:   big.NewInt(int64(startBlock + amount)),
	})
	// ethclient returns non pointers to logs, convert to pointer form
	logPointers := make([]*types.Log, len(logs2))
	for i := range logs2 {
		logPointers[i] = &logs2[i]
	}
	require.Equal(t, logPointers, logs)
}

type filterQuery struct {
	BlockHash *common.Hash     `json:"blockHash"`
	FromBlock *rpc.BlockNumber `json:"fromBlock"`
	ToBlock   *rpc.BlockNumber `json:"toBlock"`
	Addresses interface{}      `json:"address"`
	Topics    []interface{}    `json:"topics"`
}

func rpcCallCompare(t *testing.T, c1, c2 *rpc.Client, dumpOutput, failOnFirstMismatch bool, method string, args ...interface{}) json.RawMessage {
	res1, err := rpcCall(c1, dumpOutput, method, args...)
	require.NoError(t, err)
	res2, err := rpcCall(c2, dumpOutput, method, args...)
	require.NoError(t, err)

	res1Filtered, res2Filtered, err := filterResponses(res1, res2)
	require.NoError(t, err)

	dst1 := &bytes.Buffer{}
	err = json.Indent(dst1, res1Filtered, "", "  ")
	require.NoError(t, err, "res1: %v\n\nres1filtered: %v\n", string(res1), string(res1Filtered))
	dst2 := &bytes.Buffer{}
	err = json.Indent(dst2, res2Filtered, "", "  ")
	require.NoError(t, err, "res2: %v\n\nres2filtered: %v\n", string(res2), string(res2Filtered))

	if strings.TrimSpace(dst1.String()) != strings.TrimSpace(dst2.String()) {
		fmt.Printf("\nmethod: %v\nexpected (c1):\n%v,\nactual (c2):\n%v\n", method, dst1.String(), dst2.String())
		if failOnFirstMismatch {
			require.JSONEq(t, dst1.String(), dst2.String())
		}
	}

	return res2Filtered
}

var (
	IstanbulExtraVanity = 32 // Fixed number of extra-data bytes reserved for validator vanity
)

// IstanbulAggregatedSeal is the aggregated seal for Istanbul blocks
type IstanbulAggregatedSeal struct {
	// Bitmap is a bitmap having an active bit for each validator that signed this block
	Bitmap *big.Int
	// Signature is an aggregated BLS signature resulting from signatures by each validator that signed this block
	Signature []byte
	// Round is the round in which the signature was created.
	Round *big.Int
}

// IstanbulExtra is the extra-data for Istanbul blocks
type IstanbulExtra struct {
	// AddedValidators are the validators that have been added in the block
	AddedValidators []common.Address
	// AddedValidatorsPublicKeys are the BLS public keys for the validators added in the block
	AddedValidatorsPublicKeys [][96]byte
	// RemovedValidators is a bitmap having an active bit for each removed validator in the block
	RemovedValidators *big.Int
	// Seal is an ECDSA signature by the proposer
	Seal []byte
	// AggregatedSeal contains the aggregated BLS signature created via IBFT consensus.
	AggregatedSeal IstanbulAggregatedSeal
	// ParentAggregatedSeal contains and aggregated BLS signature for the previous block.
	ParentAggregatedSeal IstanbulAggregatedSeal
}

func filterResponses(res1, res2 json.RawMessage) (json.RawMessage, json.RawMessage, error) {
	// // For block responses, filter out null feeCurrency fields in transactions. In CEL2 we will omit them as empty.
	// resFiltered, err := execJQ(res, `if type != "array" and has("transactions") and has("totalDifficulty") and has("parentHash") and (.transactions | length > 0) then .transactions |= map(if .feeCurrency == null then del(.feeCurrency) else . end) else . end`)
	// if err != nil {
	// 	return nil, err
	// }

	// //
	// resFiltered, err := execJQ(res, `if type != "array" and has("transactionIndex") and has("nonce") and has("type") then del(.gatewayFee, .gatewayFeeRecipient, .ethCompatible) | if .feeCurrency == null then del(.feeCurrency) else . end else . end`)
	// if err != nil {
	// 	return nil, err
	// }

	parsed2, err := gabs.ParseJSON(res2)
	if err != nil {
		return nil, nil, err
	}

	// block
	if parsed2.Exists("parentHash") && parsed2.Exists("totalDifficulty") {
		res1Filtered, err := execJQ([]byte(res1), "del(.gasLimit, .randomness, .uncles, .sha3Uncles, .size, .epochSnarkData, .transactions[].chainId)")
		if err != nil {
			return nil, nil, err
		}
		res2Filtered, err := execJQ([]byte(res2), "del(.gasLimit, .randomness, .uncles, .sha3Uncles, .size, .epochSnarkData, .transactions[].chainId, .mixHash, .nonce)")
		if err != nil {
			return nil, nil, err
		}

		extraData1, err := execJQ([]byte(res1), ".extraData")
		if err != nil {
			return nil, nil, err
		}

		if len(extraData1) < IstanbulExtraVanity {
			return nil, nil, fmt.Errorf("invalid istanbul header extra-data length from res1: %d", len(extraData1))
		}

		istanbulExtra := IstanbulExtra{}
		err = rlp.DecodeBytes(hexutil.Bytes(extraData1)[IstanbulExtraVanity:], &istanbulExtra)
		if err != nil {
			return nil, nil, err
		}

		// istanbulExtra := IstanbulExtra{}
		// err = json.Unmarshal(.UnmarshalJSON).UnmarshalFixedJSON()Bytes(extraData1)[IstanbulExtraVanity:], &istanbulExtra)
		// if err != nil {
		// 	return nil, nil, err
		// }

		istanbulExtra.AggregatedSeal = IstanbulAggregatedSeal{}

		payload, err := json.Marshal(&istanbulExtra)
		if err != nil {
			return nil, nil, err
		}

		res1Filtered, err = execJQ([]byte(res1Filtered), ".extraData |=", string(append(extraData1[:IstanbulExtraVanity], payload...)))

		return res1Filtered, res2Filtered, nil
	}

	// transaction
	if parsed2.Exists("transactionIndex") && parsed2.Exists("nonce") && parsed2.Exists("type") {
		res1Filtered, err := execJQ([]byte(res1), "del(.chainId) else . end")
		if err != nil {
			return nil, nil, err
		}
		res2Filtered, err := execJQ([]byte(res2), "del(.chainId) else . end")
		if err != nil {
			return nil, nil, err
		}
		return res1Filtered, res2Filtered, nil
	}

	return res1, res2, nil
}

func rpcCall(c *rpc.Client, dumpOutput bool, method string, args ...interface{}) (json.RawMessage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()
	var m json.RawMessage
	err := c.CallContext(ctx, &m, method, args...)
	if err != nil {
		return nil, err
	}
	if dumpOutput {
		dst := &bytes.Buffer{}
		err = json.Indent(dst, m, "", "  ")
		if err != nil {
			return nil, err
		}
		fmt.Printf("%v\n%v\n", method, dst.String())
	}
	return m, nil
}

func execJQ(json []byte, command ...string) ([]byte, error) {
	cmd := exec.Command("jq", command...)
	cmd.Stdin = bytes.NewBuffer(json)
	buf := new(bytes.Buffer)
	cmd.Stdout = buf
	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("jq failed on input: %v\n\n with command: %v\n error: %v", string(json), command[0], err)
	}
	return buf.Bytes(), nil
}

func TestExecJq(t *testing.T) {
	res, err := execJQ([]byte(testblock), "del(.transactions[].gatewayFeeRecipient)")
	require.NoError(t, err)

	fmt.Printf("result %v\n", string(res))

	res, err = execJQ([]byte(testblock), "--exit-status", `has("kparentHash", "ktotalDifficulty")`)
	require.NoError(t, err)
	require.Equal(t, "true", string(res))
}

var testblock = `{
  "difficulty": "0x0",
  "extraData": "0xd983010000846765746889676f312e31332e3130856c696e7578000000000000f8c2c0c080b841b1e6fb24531d1ee3a145773a1c9b943dda5942a1183ae3375f49fc84b7a5622535743f16ad93b85277d2fec8393dd9a3326a2751e2133b3237ff06006ecf478001f83c890184a277df3fd1fffbb0346de4cec20deee720f3d5723c730328764b6f591d1deed14fea11a194919e67776a0d80e7998e9ece4bc35fdf56678180f83c8901efefffff7ff7fffbb0b37d8e6746f976aa8c5c8d363f514d1ca4fc48c133eecbc3c2bdaccd4153f382a08de7d009ead2bd643758b3882f558080",
  "gasUsed": "0x365f1",
  "hash": "0x24303658ca351b0cd03269ec9205e4d9d789f7609b3165054d0aeba86fb7c985",
  "logsBloom": "0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
  "miner": "0x39ec4f2a82f9f0f39929415c65db9ea5df54e41d",
  "number": "0xb05",
  "parentHash": "0x5bac312111c91bd96d0b18f1ddd645e8c64ef605cbdfc591b0bd33101e6c7ad8",
  "receiptsRoot": "0x0367e31620bf2a40cd63abeeb5dc675b1f10881b508c1f696cd8579867f7c5b6",
  "stateRoot": "0x55a2fcfa962c3915a56928bbd5c2ae277994bad019d991e564c68cd822c9499c",
  "timestamp": "0x5ea0a46b",
  "totalDifficulty": "0xb06",
  "transactions": [
    {
      "blockHash": "0x24303658ca351b0cd03269ec9205e4d9d789f7609b3165054d0aeba86fb7c985",
      "blockNumber": "0xb05",
      "ethCompatible": false,
      "feeCurrency": null,
      "from": "0xe23a4c6615669526ab58e9c37088bee4ed2b2dee",
      "gas": "0x1312d00",
      "gasPrice": "0x2540be400",
      "gatewayFeeRecipient": null,
      "hash": "0x721a8fc581d2e699294f10398d3e4ffeb9af5d9c6a08ba1eb76e2ad29617f375",
      "input": "0x608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506102ae806100606000396000f3fe608060405234801561001057600080fd5b506004361061004c5760003560e01c80630900f01014610051578063445df0ac146100955780638da5cb5b146100b3578063fdacd576146100fd575b600080fd5b6100936004803603602081101561006757600080fd5b81019080803573ffffffffffffffffffffffffffffffffffffffff16906020019092919050505061012b565b005b61009d6101f7565b6040518082815260200191505060405180910390f35b6100bb6101fd565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6101296004803603602081101561011357600080fd5b8101908080359060200190929190505050610222565b005b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff1614156101f45760008190508073ffffffffffffffffffffffffffffffffffffffff1663fdacd5766001546040518263ffffffff1660e01b815260040180828152602001915050600060405180830381600087803b1580156101da57600080fd5b505af11580156101ee573d6000803e3d6000fd5b50505050505b50565b60015481565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1681565b6000809054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff163373ffffffffffffffffffffffffffffffffffffffff16141561027f57806001819055505b5056fea165627a7a723058208d331f3378bb7ac07b45c7b734ec6883c68482f8eed0e59739f5b82af926f5810029",
      "nonce": "0x0",
      "r": "0x97d69d5ec9d61c72835acbfa8cfa846fad6584002f0dd503c5f4f67f721025e1",
      "s": "0x7aebe0781721a38da5bce745c6b923306e30154e63156e3742b2cd198e0233a7",
      "to": null,
      "transactionIndex": "0x0",
      "type": "0x0",
      "v": "0x149fb",
      "value": "0x0"
    }
  ],
  "transactionsRoot": "0x82674f07b62b9c6596674894e8cc17c50d444b0f2c57fbe1f9a7be9da14d0718"
}`
