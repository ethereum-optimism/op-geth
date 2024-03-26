package policies

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

type messageSafetyLabel int

const (
	invalid messageSafetyLabel = iota
	safe
	finalized
)

var (
	inboxExecuteMessageSignature = "executeMessage(address,bytes,(address,uint256,uint256,uint256,uint256))"
	inboxExecuteMessageBytes4    = crypto.Keccak256([]byte(inboxExecuteMessageSignature))[:4]
)

type messageIdentifier struct {
	Origin      common.Address
	BlockNumber *big.Int
	LogIndex    uint64
	Timestamp   uint64
	ChainId     *big.Int
}

// Manually parse the tx data according to the function signature. [TODO] introduce a soliabi
// utility OR pass the entire calldata to the backend to centralize where deser logic happens
// with minimal validation here (i.e min tx data byte length)
func unpackInboxExecutionMessageTxData(txData []byte) (*messageIdentifier, []byte, error) {
	// Function Selector
	if len(txData) <= 4 {
		return nil, nil, fmt.Errorf("invalid calldata: function selector")
	} else if !bytes.Equal(txData[:4], inboxExecuteMessageBytes4) {
		return nil, nil, fmt.Errorf("invalid function selector")
	}
	txData = txData[4:]

	// the argument calldata must include at least 8 words (including message byte length)
	if len(txData) < 8*32 {
		return nil, nil, fmt.Errorf("invalid calldata: executeMessage function args")
	}

	// Message Target
	txData = txData[32:]

	// Message Bytes Calldata Location -- The 7th word
	msgDataLoc := common.BytesToHash(txData[:32])
	if msgDataLoc != common.HexToHash("0xe0") { // 7*32-4 (removing function selector)
		return nil, nil, fmt.Errorf("invalid calldata: msg bytes data loc: %x", msgDataLoc)
	}
	txData = txData[32:]

	// Message Identifier
	var uint64Padding [24]byte
	var addressPadding [12]byte
	msgId := messageIdentifier{}

	msgId.Origin = common.BytesToAddress(txData[:32])
	if !bytes.Equal(txData[:12], addressPadding[:]) {
		return nil, nil, fmt.Errorf("origin address padding is non-zero")
	}
	txData = txData[32:]

	msgId.BlockNumber = new(big.Int).SetBytes(txData[:32])
	txData = txData[32:]

	msgId.LogIndex = new(big.Int).SetBytes(txData[:32]).Uint64()
	if !bytes.Equal(txData[:24], uint64Padding[:]) {
		return nil, nil, fmt.Errorf("log index padding is non-zero")
	}
	txData = txData[32:]

	msgId.Timestamp = new(big.Int).SetBytes(txData[:32]).Uint64()
	if !bytes.Equal(txData[:24], uint64Padding[:]) {
		return nil, nil, fmt.Errorf("timestamp padding is non-zero")
	}
	txData = txData[32:]

	msgId.ChainId = new(big.Int).SetBytes(txData[:32])
	txData = txData[32:]

	// Message Bytes
	byteLen := new(big.Int).SetBytes(txData[:32]).Uint64()
	if !bytes.Equal(txData[:24], uint64Padding[:]) {
		return nil, nil, fmt.Errorf("log index padding is non-zero")
	}
	paddedByteLength := byteLen + (32 - (byteLen % 32))
	if uint64(len(txData)) != 32+paddedByteLength {
		return nil, nil, fmt.Errorf("invalid calldata: too many bytes")
	}
	txData = txData[32:]

	msgBytes := txData[:byteLen]
	return &msgId, msgBytes, nil
}
