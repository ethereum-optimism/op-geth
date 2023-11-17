package types

import (
	"crypto/ecdsa"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/stretchr/testify/require"
)

var (
	userKey, _           = crypto.HexToECDSA("eef77acb6c6a6eebc5b363a475ac583ec7eccdb42b6481424c60f59aa326547f")
	gasFeeSponsorKey1, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	gasFeeSponsorKey2, _ = crypto.HexToECDSA("0288ef00023598499cb6c940146d050d2b1fb914198c327f76aad590bead68b6")
)

func generateMetaTxData(dynamicTx *DynamicFeeTx, expireHeight uint64, sponsorPercent uint64,
	gasFeeSponsorAddr common.Address, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	metaTxSignData := &MetaTxSignData{
		ChainID:        dynamicTx.ChainID,
		Nonce:          dynamicTx.Nonce,
		GasTipCap:      dynamicTx.GasTipCap,
		GasFeeCap:      dynamicTx.GasFeeCap,
		Gas:            dynamicTx.Gas,
		To:             dynamicTx.To,
		Value:          dynamicTx.Value,
		Data:           dynamicTx.Data,
		AccessList:     dynamicTx.AccessList,
		ExpireHeight:   expireHeight,
		SponsorPercent: sponsorPercent,
	}

	sponsorSig, err := crypto.Sign(metaTxSignData.Hash().Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	r, s, v := decodeSignature(sponsorSig)

	metaTxData := &MetaTxParams{
		ExpireHeight:   expireHeight,
		Payload:        metaTxSignData.Data,
		GasFeeSponsor:  gasFeeSponsorAddr,
		SponsorPercent: sponsorPercent,
		R:              r,
		S:              s,
		V:              v,
	}

	metaTxDataBz, err := rlp.EncodeToBytes(metaTxData)
	if err != nil {
		return nil, err
	}

	return append(MetaTxPrefix, metaTxDataBz...), nil
}

func generateMetaTxDataWithMockSig(dynamicTx *DynamicFeeTx, expireHeight uint64, sponsorPercent uint64,
	gasFeeSponsorAddr common.Address, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	metaTxSignData := &MetaTxSignData{
		ChainID:        dynamicTx.ChainID,
		Nonce:          dynamicTx.Nonce,
		GasTipCap:      dynamicTx.GasTipCap,
		GasFeeCap:      dynamicTx.GasFeeCap,
		Gas:            dynamicTx.Gas,
		To:             dynamicTx.To,
		Value:          dynamicTx.Value,
		Data:           dynamicTx.Data,
		AccessList:     dynamicTx.AccessList,
		ExpireHeight:   expireHeight,
		SponsorPercent: sponsorPercent,
	}

	sponsorSig, err := crypto.Sign(metaTxSignData.Hash().Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	sponsorSig[len(sponsorSig)-1] = sponsorSig[len(sponsorSig)-1] + 2

	r, s, v := decodeSignature(sponsorSig)
	metaTxData := &MetaTxParams{
		ExpireHeight:   expireHeight,
		Payload:        metaTxSignData.Data,
		GasFeeSponsor:  gasFeeSponsorAddr,
		SponsorPercent: sponsorPercent,
		R:              r,
		S:              s,
		V:              v,
	}

	metaTxDataBz, err := rlp.EncodeToBytes(metaTxData)
	if err != nil {
		return nil, err
	}

	return append(MetaTxPrefix, metaTxDataBz...), nil
}

func TestDecodeMetaTxParams(t *testing.T) {
	gasFeeSponsorPublicKey := gasFeeSponsorKey1.Public()
	pubKeyECDSA, _ := gasFeeSponsorPublicKey.(*ecdsa.PublicKey)
	gasFeeSponsorAddr := crypto.PubkeyToAddress(*pubKeyECDSA)

	chainId := big.NewInt(1)
	depositABICalldata, _ := hexutil.Decode("0xd0e30db0")
	to := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	expireHeight := uint64(20_000_010)
	dynamicTx := &DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      100,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(1e15),
		Gas:        4700000,
		To:         &to,
		Value:      big.NewInt(1e18),
		Data:       depositABICalldata,
		AccessList: nil,
	}

	metaTxData := &MetaTxParams{
		ExpireHeight:   expireHeight,
		Payload:        depositABICalldata,
		GasFeeSponsor:  gasFeeSponsorAddr,
		SponsorPercent: 50,
	}

	metaTxDataBz, err := rlp.EncodeToBytes(metaTxData)
	require.NoError(t, err)

	dynamicTx.Data = append(MetaTxPrefix, metaTxDataBz...)

	metaTxParams, err := DecodeMetaTxParams(dynamicTx.Data)
	require.NoError(t, err)

	require.Equal(t, gasFeeSponsorAddr.String(), metaTxParams.GasFeeSponsor.String())
	require.Equal(t, hexutil.Encode(depositABICalldata), hexutil.Encode(metaTxParams.Payload))

	metaTxData = &MetaTxParams{
		ExpireHeight:   expireHeight,
		Payload:        depositABICalldata,
		GasFeeSponsor:  gasFeeSponsorAddr,
		SponsorPercent: 101,
	}

	metaTxDataBz, err = rlp.EncodeToBytes(metaTxData)
	require.NoError(t, err)

	dynamicTx.Data = append(MetaTxPrefix, metaTxDataBz...)

	metaTxParams, err = DecodeMetaTxParams(dynamicTx.Data)
	require.Equal(t, ErrInvalidSponsorPercent, err)

}

func TestDecodeAndVerifyMetaTxParams(t *testing.T) {
	gasFeeSponsorPublicKey := gasFeeSponsorKey1.Public()
	pubKeyECDSA, _ := gasFeeSponsorPublicKey.(*ecdsa.PublicKey)
	gasFeeSponsorAddr := crypto.PubkeyToAddress(*pubKeyECDSA)

	chainId := big.NewInt(1)
	depositABICalldata, _ := hexutil.Decode("0xd0e30db0")
	to := common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2")
	expireHeight := uint64(20_000_010)
	dynamicTx := &DynamicFeeTx{
		ChainID:    chainId,
		Nonce:      100,
		GasTipCap:  big.NewInt(1e9),
		GasFeeCap:  big.NewInt(1e15),
		Gas:        4700000,
		To:         &to,
		Value:      big.NewInt(1e18),
		Data:       depositABICalldata,
		AccessList: nil,
	}

	payload, err := generateMetaTxData(dynamicTx, expireHeight, 50, gasFeeSponsorAddr, gasFeeSponsorKey1)
	require.NoError(t, err)

	dynamicTx.Data = payload
	tx := NewTx(dynamicTx)
	signer := LatestSignerForChainID(chainId)
	txSignature, err := crypto.Sign(signer.Hash(tx).Bytes(), userKey)
	require.NoError(t, err)
	signedTx, err := tx.WithSignature(signer, txSignature)
	require.NoError(t, err)

	// test normal metaTx
	metaTxParams, err := DecodeAndVerifyMetaTxParams(signedTx)
	require.NoError(t, err)

	require.Equal(t, gasFeeSponsorAddr.String(), metaTxParams.GasFeeSponsor.String())
	require.Equal(t, hexutil.Encode(depositABICalldata), hexutil.Encode(metaTxParams.Payload))

	// Test ErrInvalidGasFeeSponsorSig
	dynamicTx.Data = depositABICalldata
	payload, err = generateMetaTxDataWithMockSig(dynamicTx, expireHeight, 100, gasFeeSponsorAddr, gasFeeSponsorKey1)
	require.NoError(t, err)

	dynamicTx.Data = payload
	tx = NewTx(dynamicTx)
	txSignature, err = crypto.Sign(signer.Hash(tx).Bytes(), userKey)
	require.NoError(t, err)
	signedTx, err = tx.WithSignature(signer, txSignature)
	require.NoError(t, err)

	_, err = DecodeAndVerifyMetaTxParams(signedTx)
	require.Equal(t, err, ErrInvalidGasFeeSponsorSig)

	// Test ErrGasFeeSponsorMismatch
	dynamicTx.Data = depositABICalldata
	payload, err = generateMetaTxData(dynamicTx, expireHeight, 80, gasFeeSponsorAddr, gasFeeSponsorKey2)
	require.NoError(t, err)

	dynamicTx.Data = payload
	tx = NewTx(dynamicTx)
	txSignature, err = crypto.Sign(signer.Hash(tx).Bytes(), userKey)
	require.NoError(t, err)
	signedTx, err = tx.WithSignature(signer, txSignature)
	require.NoError(t, err)

	_, err = DecodeAndVerifyMetaTxParams(signedTx)
	require.Equal(t, err, ErrGasFeeSponsorMismatch)

	// Test ErrGasFeeSponsorMismatch
	dynamicTx.Data = depositABICalldata
	payload, err = generateMetaTxData(dynamicTx, expireHeight, 101, gasFeeSponsorAddr, gasFeeSponsorKey2)
	require.NoError(t, err)

	dynamicTx.Data = payload
	tx = NewTx(dynamicTx)
	txSignature, err = crypto.Sign(signer.Hash(tx).Bytes(), userKey)
	require.NoError(t, err)
	signedTx, err = tx.WithSignature(signer, txSignature)
	require.NoError(t, err)

	_, err = DecodeAndVerifyMetaTxParams(signedTx)
	require.Equal(t, err, ErrInvalidSponsorPercent)
}
