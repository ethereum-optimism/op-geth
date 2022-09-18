package types

import (
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/tests/fuzzerutils"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

// FuzzTransactionMarshallingRoundTrip executes a fuzz test which constructs arbitrary transactions of various types
// and tests round-trip marshalling of them to ensure original values are kept intact. This tests JSON and binary
// marshalling.
func FuzzTransactionMarshallingRoundTrip(f *testing.F) {
	signingKey, _ := crypto.GenerateKey()
	signer := LatestSignerForChainID(big.NewInt(0))
	f.Fuzz(func(t *testing.T, goFuzzSeed []byte, transactionType uint64) {
		// Create our fuzzer wrapper to generate complex values
		typeProvider := fuzz.NewFromGoFuzz(goFuzzSeed).NilChance(0).MaxDepth(10000).NumElements(0, 0x100).AllowUnexportedFields(true)
		fuzzerutils.AddFuzzerFunctions(typeProvider)

		// Bound our transaction type to determine which type of transaction to construct.
		transactionType %= 4

		// Create the underlying fuzzed transaction.
		var tx *Transaction
		if transactionType == 0 {
			// Create a legacy tx with fields populated by the fuzzer
			innerTx := &LegacyTx{}
			typeProvider.Fuzz(&innerTx)
			tx = MustSignNewTx(signingKey, signer, innerTx)
		} else if transactionType == 1 {
			// Create an access list tx with fields populated by the fuzzer
			innerTx := &AccessListTx{}
			typeProvider.Fuzz(&innerTx)
			innerTx.ChainID = signer.ChainID()
			tx = MustSignNewTx(signingKey, signer, innerTx)
		} else if transactionType == 2 {
			// Create a dynamic fee tx with fields populated by the fuzzer
			innerTx := &DynamicFeeTx{}
			typeProvider.Fuzz(&innerTx)
			innerTx.ChainID = signer.ChainID()
			tx = MustSignNewTx(signingKey, signer, innerTx)
		} else {
			// Create a deposit tx with fields populated by the fuzzer
			innerTx := &DepositTx{}
			typeProvider.Fuzz(&innerTx)
			tx = NewTx(innerTx)
		}

		// Define our decoded tx data we will test
		decodedInnerTxs := make([]TxData, 0)

		// Perform round trip JSON encoding
		jsonTx, err := tx.MarshalJSON()
		require.NoError(t, err)
		decodedJsonTx := &Transaction{}
		err = decodedJsonTx.UnmarshalJSON(jsonTx)
		require.NoError(t, err)
		decodedInnerTxs = append(decodedInnerTxs, decodedJsonTx.inner)

		// Perform round trip binary encoding
		binaryTx, err := tx.MarshalBinary()
		decodedBinaryTx := &Transaction{}
		err = decodedBinaryTx.UnmarshalBinary(binaryTx)
		require.NoError(t, err)
		decodedInnerTxs = append(decodedInnerTxs, decodedBinaryTx.inner)

		// Loop for each tx data to compare.
		// Note: We do this rather than require.EqualValues() as big.Int types can differ when instantiated vs
		// deserialized because their inner big.Nat values may contain nil or not when empty.
		for _, decodedTxDataGeneric := range decodedInnerTxs {
			switch txData := tx.inner.(type) {
			case *LegacyTx:
				// Verify our type and cast appropriately.
				switch decodedTxData := decodedTxDataGeneric.(type) {
				case *LegacyTx:
					require.EqualValues(t, txData.Data, decodedTxData.Data)
					require.EqualValues(t, txData.Value.Bytes(), decodedTxData.Value.Bytes())
					require.EqualValues(t, txData.To, decodedTxData.To)
					require.EqualValues(t, txData.Gas, decodedTxData.Gas)
					require.EqualValues(t, txData.GasPrice.Bytes(), decodedTxData.GasPrice.Bytes())
					require.EqualValues(t, txData.Nonce, decodedTxData.Nonce)
					require.EqualValues(t, txData.V.Bytes(), decodedTxData.V.Bytes())
					require.EqualValues(t, txData.R.Bytes(), decodedTxData.R.Bytes())
					require.EqualValues(t, txData.S.Bytes(), decodedTxData.S.Bytes())
					break
				default:
					require.Fail(t, "decoded tx should have been LegacyTx")
					break
				}
			case *AccessListTx:
				// Verify our type and cast appropriately.
				switch decodedTxData := decodedTxDataGeneric.(type) {
				case *AccessListTx:
					require.EqualValues(t, txData.Data, decodedTxData.Data)
					require.EqualValues(t, txData.Value.Bytes(), decodedTxData.Value.Bytes())
					require.EqualValues(t, txData.To, decodedTxData.To)
					require.EqualValues(t, txData.Gas, decodedTxData.Gas)
					require.EqualValues(t, txData.Nonce, decodedTxData.Nonce)
					require.EqualValues(t, txData.ChainID.Bytes(), decodedTxData.ChainID.Bytes())
					require.EqualValues(t, txData.GasPrice.Bytes(), decodedTxData.GasPrice.Bytes())
					require.EqualValues(t, txData.AccessList, decodedTxData.AccessList)
					require.EqualValues(t, txData.V.Bytes(), decodedTxData.V.Bytes())
					require.EqualValues(t, txData.R.Bytes(), decodedTxData.R.Bytes())
					require.EqualValues(t, txData.S.Bytes(), decodedTxData.S.Bytes())
					break
				default:
					require.Fail(t, "decoded tx should have been AccessListTx")
					break
				}
				break
			case *DynamicFeeTx:
				// Verify our type and cast appropriately.
				switch decodedTxData := decodedTxDataGeneric.(type) {
				case *DynamicFeeTx:
					require.EqualValues(t, txData.Data, decodedTxData.Data)
					require.EqualValues(t, txData.Value.Bytes(), decodedTxData.Value.Bytes())
					require.EqualValues(t, txData.To, decodedTxData.To)
					require.EqualValues(t, txData.Gas, decodedTxData.Gas)
					require.EqualValues(t, txData.Nonce, decodedTxData.Nonce)
					require.EqualValues(t, txData.ChainID.Bytes(), decodedTxData.ChainID.Bytes())
					require.EqualValues(t, txData.GasTipCap.Bytes(), decodedTxData.GasTipCap.Bytes())
					require.EqualValues(t, txData.GasFeeCap.Bytes(), decodedTxData.GasFeeCap.Bytes())
					require.EqualValues(t, txData.AccessList, decodedTxData.AccessList)
					require.EqualValues(t, txData.V.Bytes(), decodedTxData.V.Bytes())
					require.EqualValues(t, txData.R.Bytes(), decodedTxData.R.Bytes())
					require.EqualValues(t, txData.S.Bytes(), decodedTxData.S.Bytes())
					_ = decodedTxData
					break
				default:
					require.Fail(t, "decoded tx should have been DynamicFeeTx")
					break
				}
				break
			case *DepositTx:
				// Verify our type and cast appropriately.
				switch decodedTxData := decodedTxDataGeneric.(type) {
				case *DepositTx:
					require.EqualValues(t, txData.Data, decodedTxData.Data)
					require.EqualValues(t, txData.Value.Bytes(), decodedTxData.Value.Bytes())
					require.EqualValues(t, txData.To, decodedTxData.To)
					require.EqualValues(t, txData.Gas, decodedTxData.Gas)
					require.EqualValues(t, txData.From, decodedTxData.From)
					require.EqualValues(t, txData.SourceHash, decodedTxData.SourceHash)
					require.EqualValues(t, txData.Mint.Bytes(), decodedTxData.Mint.Bytes())
					require.EqualValues(t, txData.IsSystemTransaction, decodedTxData.IsSystemTransaction)
					break
				default:
					require.Fail(t, "decoded tx should have been DepositTx")
					break
				}
				break
			default:
				require.Fail(t, "original tx inner data could not be casted")
			}
		}
	})
}
