package types

import (
	"encoding/binary"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/require"
)

var (
	basefee  = big.NewInt(1000 * 1e6)
	overhead = big.NewInt(50)
	scalar   = big.NewInt(7 * 1e6)

	blobBasefee       = big.NewInt(10 * 1e6)
	basefeeScalar     = big.NewInt(2)
	blobBasefeeScalar = big.NewInt(3)

	// below are the expected cost func outcomes for the above parameter settings on the emptyTx
	// which is defined in transaction_test.go
	bedrockFee  = big.NewInt(11326000000000)
	regolithFee = big.NewInt(3710000000000)
	ecotoneFee  = big.NewInt(960900) // (480/16)*(2*16*1000 + 3*10) == 960900

	bedrockGas  = big.NewInt(1618)
	regolithGas = big.NewInt(530) // 530  = 1618 - (16*68)
	ecotoneGas  = big.NewInt(480)
)

func TestBedrockL1CostFunc(t *testing.T) {
	costFunc0 := newL1CostFuncBedrockHelper(basefee, overhead, scalar, false /*isRegolith*/)
	costFunc1 := newL1CostFuncBedrockHelper(basefee, overhead, scalar, true)

	c0, g0 := costFunc0(emptyTx.RollupCostData()) // pre-Regolith
	c1, g1 := costFunc1(emptyTx.RollupCostData())

	require.Equal(t, bedrockFee, c0)
	require.Equal(t, bedrockGas, g0) // gas-used

	require.Equal(t, regolithFee, c1)
	require.Equal(t, regolithGas, g1)
}

func TestEcotoneL1CostFunc(t *testing.T) {
	costFunc := newL1CostFuncEcotone(basefee, blobBasefee, basefeeScalar, blobBasefeeScalar)
	c, g := costFunc(emptyTx.RollupCostData())
	require.Equal(t, ecotoneGas, g)
	require.Equal(t, ecotoneFee, c)
}

func TestExtractBedrockGasParams(t *testing.T) {
	regolithTime := uint64(1)
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &regolithTime,
	}

	data := getBedrockL1Attributes(basefee, overhead, scalar)

	_, costFuncPreRegolith, _, err := extractL1GasParams(config, regolithTime-1, data)
	require.NoError(t, err)

	// Function should continue to succeed even with extra data (that just gets ignored) since we
	// have been testing the data size is at least the expected number of bytes instead of exactly
	// the expected number of bytes. It's unclear if this flexibility was intentional, but since
	// it's been in production we shouldn't change this behavior.
	data = append(data, []byte{0xBE, 0xEE, 0xEE, 0xFF}...) // tack on garbage data
	_, costFuncRegolith, _, err := extractL1GasParams(config, regolithTime, data)
	require.NoError(t, err)

	c, _ := costFuncPreRegolith(emptyTx.RollupCostData())
	require.Equal(t, bedrockFee, c)

	c, _ = costFuncRegolith(emptyTx.RollupCostData())
	require.Equal(t, regolithFee, c)

	// try to extract from data which has not enough params, should get error.
	data = data[:len(data)-4-32]
	_, _, _, err = extractL1GasParams(config, regolithTime, data)
	require.Error(t, err)
}

func TestExtractEcotoneGasParams(t *testing.T) {
	zeroTime := uint64(0)
	// create a config where ecotone upgrade is active
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &zeroTime,
		EcotoneTime:  &zeroTime,
	}
	require.True(t, config.IsOptimismEcotone(0))

	data := getEcotoneL1Attributes(basefee, blobBasefee, basefeeScalar, blobBasefeeScalar)

	_, costFunc, _, err := extractL1GasParams(config, 0, data)
	require.NoError(t, err)

	c, g := costFunc(emptyTx.RollupCostData())

	require.Equal(t, ecotoneGas, g)
	require.Equal(t, ecotoneFee, c)

	// make sure wrong amont of data results in error
	data = append(data, 0x00) // tack on garbage byte
	_, _, err = extractL1GasParamsEcotone(data)
	require.Error(t, err)
}

// make sure the first block of the ecotone upgrade is properly detected, and invokes the bedrock
// cost function appropriately
func TestFirstBlockEcotoneGasParams(t *testing.T) {
	zeroTime := uint64(0)
	// create a config where ecotone upgrade is active
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &zeroTime,
		EcotoneTime:  &zeroTime,
	}
	require.True(t, config.IsOptimismEcotone(0))

	data := getBedrockL1Attributes(basefee, overhead, scalar)

	_, oldCostFunc, _, err := extractL1GasParams(config, 0, data)
	require.NoError(t, err)
	c, _ := oldCostFunc(emptyTx.RollupCostData())
	require.Equal(t, regolithFee, c)
}

func getBedrockL1Attributes(basefee, overhead, scalar *big.Int) []byte {
	uint256 := make([]byte, 32)
	ignored := big.NewInt(1234)
	data := []byte{}
	data = append(data, BedrockL1AttributesSelector...)
	data = append(data, ignored.FillBytes(uint256)...)  // arg 0
	data = append(data, ignored.FillBytes(uint256)...)  // arg 1
	data = append(data, basefee.FillBytes(uint256)...)  // arg 2
	data = append(data, ignored.FillBytes(uint256)...)  // arg 3
	data = append(data, ignored.FillBytes(uint256)...)  // arg 4
	data = append(data, ignored.FillBytes(uint256)...)  // arg 5
	data = append(data, overhead.FillBytes(uint256)...) // arg 6
	data = append(data, scalar.FillBytes(uint256)...)   // arg 7
	return data
}

func getEcotoneL1Attributes(basefee, blobBasefee, basefeeScalar, blobBasefeeScalar *big.Int) []byte {
	ignored := big.NewInt(1234)
	data := []byte{}
	uint256 := make([]byte, 32)
	uint64 := make([]byte, 8)
	uint32 := make([]byte, 4)
	data = append(data, EcotoneL1AttributesSelector...)
	data = append(data, basefeeScalar.FillBytes(uint32)...)
	data = append(data, blobBasefeeScalar.FillBytes(uint32)...)
	data = append(data, ignored.FillBytes(uint64)...)
	data = append(data, ignored.FillBytes(uint64)...)
	data = append(data, ignored.FillBytes(uint64)...)
	data = append(data, basefee.FillBytes(uint256)...)
	data = append(data, blobBasefee.FillBytes(uint256)...)
	data = append(data, ignored.FillBytes(uint256)...)
	data = append(data, ignored.FillBytes(uint256)...)
	return data
}

type testStateGetter struct {
	basefee, blobBasefee, overhead, scalar *big.Int
	basefeeScalar, blobBasefeeScalar       uint32
}

func (sg *testStateGetter) GetState(addr common.Address, slot common.Hash) common.Hash {
	buf := common.Hash{}
	switch slot {
	case L1BasefeeSlot:
		sg.basefee.FillBytes(buf[:])
	case OverheadSlot:
		sg.overhead.FillBytes(buf[:])
	case ScalarSlot:
		sg.scalar.FillBytes(buf[:])
	case L1BlobBasefeeSlot:
		sg.blobBasefee.FillBytes(buf[:])
	case L1FeeScalarsSlot:
		offset := scalarSectionStart
		binary.BigEndian.PutUint32(buf[offset:offset+4], sg.basefeeScalar)
		binary.BigEndian.PutUint32(buf[offset+4:offset+8], sg.blobBasefeeScalar)
	default:
		panic("unknown slot")
	}
	return buf
}

// TestNewL1CostFunc tests that the appropriate cost function is selected based on the
// configuration and statedb values.
func TestNewL1CostFunc(t *testing.T) {
	time := uint64(1)
	config := &params.ChainConfig{
		Optimism: params.OptimismTestConfig.Optimism,
	}
	statedb := &testStateGetter{
		basefee:           basefee,
		overhead:          overhead,
		scalar:            scalar,
		blobBasefee:       blobBasefee,
		basefeeScalar:     uint32(basefeeScalar.Uint64()),
		blobBasefeeScalar: uint32(blobBasefeeScalar.Uint64()),
	}

	costFunc := NewL1CostFunc(config, statedb)
	require.NotNil(t, costFunc)

	// empty cost data should result in nil fee
	fee := costFunc(RollupCostData{}, time)
	require.Nil(t, fee)

	// emptyTx fee w/ bedrock config should be the bedrock fee
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, bedrockFee, fee)

	// emptyTx fee w/ regolith config should be the regolith fee
	config.RegolithTime = &time
	costFunc = NewL1CostFunc(config, statedb)
	require.NotNil(t, costFunc)
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, regolithFee, fee)

	// emptyTx fee w/ ecotone config should be the ecotone fee
	config.EcotoneTime = &time
	costFunc = NewL1CostFunc(config, statedb)
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, ecotoneFee, fee)

	// emptyTx fee w/ ecotone config, but simulate first ecotone block by blowing away the ecotone
	// params. Should result in regolith fee.
	statedb.basefeeScalar = 0
	statedb.blobBasefeeScalar = 0
	statedb.blobBasefee = new(big.Int)
	costFunc = NewL1CostFunc(config, statedb)
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, regolithFee, fee)
}
