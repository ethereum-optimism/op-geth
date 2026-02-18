package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/holiman/uint256"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/params"
)

var (
	baseFee  = big.NewInt(1000 * 1e6)
	overhead = big.NewInt(50)
	scalar   = big.NewInt(7 * 1e6)

	blobBaseFee         = big.NewInt(10 * 1e6)
	baseFeeScalar       = big.NewInt(2)
	blobBaseFeeScalar   = big.NewInt(3)
	operatorFeeScalar   = big.NewInt(1439103868)
	operatorFeeConstant = big.NewInt(1256417826609331460)

	// below are the expected cost func outcomes for the above parameter settings on the emptyTx
	// which is defined in transaction_test.go
	bedrockFee  = big.NewInt(11326000000000)
	regolithFee = big.NewInt(3710000000000)
	ecotoneFee  = big.NewInt(960900) // (480/16)*(2*16*1000 + 3*10) == 960900
	// the emptyTx is out of bounds for the linear regression so it uses the minimum size
	fjordFee          = big.NewInt(3203000)                 // 100_000_000 * (2 * 1000 * 1e6 * 16 + 3 * 10 * 1e6) / 1e12
	ithmusOperatorFee = uint256.NewInt(1256417826611659930) // 1618 * 1439103868 / 1e6 + 1256417826609331460
	jovianOperatorFee = uint256.NewInt(1256650673615173860) // 1618 * 1439103868 * 100 + 1256417826609331460

	bedrockGas      = big.NewInt(1618)
	regolithGas     = big.NewInt(530) // 530  = 1618 - (16*68)
	ecotoneGas      = big.NewInt(480)
	minimumFjordGas = big.NewInt(1600) // fastlz size of minimum txn, 100_000_000 * 16 / 1e6
)

func TestBedrockL1CostFunc(t *testing.T) {
	costFunc0 := newL1CostFuncBedrockHelper(baseFee, overhead, scalar, false /*isRegolith*/)
	costFunc1 := newL1CostFuncBedrockHelper(baseFee, overhead, scalar, true)

	c0, g0 := costFunc0(emptyTx.RollupCostData()) // pre-Regolith
	c1, g1 := costFunc1(emptyTx.RollupCostData())

	require.Equal(t, bedrockFee, c0)
	require.Equal(t, bedrockGas, g0) // gas-used

	require.Equal(t, regolithFee, c1)
	require.Equal(t, regolithGas, g1)
}

func TestEcotoneL1CostFunc(t *testing.T) {
	costFunc := newL1CostFuncEcotone(baseFee, blobBaseFee, baseFeeScalar, blobBaseFeeScalar)

	c0, g0 := costFunc(emptyTx.RollupCostData())

	require.Equal(t, ecotoneGas, g0)
	require.Equal(t, ecotoneFee, c0)
}

func TestFjordL1CostFuncMinimumBounds(t *testing.T) {
	costFunc := NewL1CostFuncFjord(
		baseFee,
		blobBaseFee,
		baseFeeScalar,
		blobBaseFeeScalar,
	)

	// Minimum size transactions:
	// -42.5856 + 0.8365*110 = 49.4294
	// -42.5856 + 0.8365*150 = 82.8894
	// -42.5856 + 0.8365*170 = 99.6194
	for _, fastLzsize := range []uint64{100, 150, 170} {
		c, g := costFunc(RollupCostData{
			FastLzSize: fastLzsize,
		})

		require.Equal(t, minimumFjordGas, g)
		require.Equal(t, fjordFee, c)
	}

	// Larger size transactions:
	// -42.5856 + 0.8365*171 = 100.4559
	// -42.5856 + 0.8365*175 = 108.8019
	// -42.5856 + 0.8365*200 = 124.7144
	for _, fastLzsize := range []uint64{171, 175, 200} {
		c, g := costFunc(RollupCostData{
			FastLzSize: fastLzsize,
		})

		require.Greater(t, g.Uint64(), minimumFjordGas.Uint64())
		require.Greater(t, c.Uint64(), fjordFee.Uint64())
	}
}

// TestFjordL1CostSolidityParity tests that the cost function for the fjord upgrade matches a Solidity
// test to ensure the outputs are the same.
func TestFjordL1CostSolidityParity(t *testing.T) {
	costFunc := NewL1CostFuncFjord(
		big.NewInt(2*1e6),
		big.NewInt(3*1e6),
		big.NewInt(20),
		big.NewInt(15),
	)

	c0, g0 := costFunc(RollupCostData{
		FastLzSize: 235,
	})

	require.Equal(t, big.NewInt(2463), g0)
	require.Equal(t, big.NewInt(105484), c0)
}

func TestExtractBedrockGasParams(t *testing.T) {
	regolithTime := uint64(1)
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &regolithTime,
	}

	data := getBedrockL1Attributes(baseFee, overhead, scalar)

	gasparams, err := extractL1GasParams(config, regolithTime-1, data)
	costFuncPreRegolith := gasparams.costFunc
	require.NoError(t, err)

	// Function should continue to succeed even with extra data (that just gets ignored) since we
	// have been testing the data size is at least the expected number of bytes instead of exactly
	// the expected number of bytes. It's unclear if this flexibility was intentional, but since
	// it's been in production we shouldn't change this behavior.
	data = append(data, []byte{0xBE, 0xEE, 0xEE, 0xFF}...) // tack on garbage data
	gasparams, err = extractL1GasParams(config, regolithTime, data)
	costFuncRegolith := gasparams.costFunc
	require.NoError(t, err)

	c, _ := costFuncPreRegolith(emptyTx.RollupCostData())
	require.Equal(t, bedrockFee, c)

	c, _ = costFuncRegolith(emptyTx.RollupCostData())
	require.Equal(t, regolithFee, c)

	// try to extract from data which has not enough params, should get error.
	data = data[:len(data)-4-32]
	_, err = extractL1GasParams(config, regolithTime, data)
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
	require.True(t, config.IsOptimismEcotone(zeroTime))

	data := getEcotoneL1Attributes(
		baseFee,
		blobBaseFee,
		baseFeeScalar,
		blobBaseFeeScalar,
	)

	gasparams, err := extractL1GasParams(config, zeroTime, data)
	require.NoError(t, err)
	costFunc := gasparams.costFunc

	c, g := costFunc(emptyTx.RollupCostData())

	require.Equal(t, ecotoneGas, g)
	require.Equal(t, ecotoneFee, c)

	// make sure wrong amont of data results in error
	data = append(data, 0x00) // tack on garbage byte
	_, err = extractL1GasParamsPostEcotone(data)
	require.Error(t, err)
}

func TestExtractFjordGasParams(t *testing.T) {
	zeroTime := uint64(0)
	// create a config where fjord is active
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &zeroTime,
		EcotoneTime:  &zeroTime,
		FjordTime:    &zeroTime,
	}
	require.True(t, config.IsOptimismFjord(zeroTime))

	data := getEcotoneL1Attributes(
		baseFee,
		blobBaseFee,
		baseFeeScalar,
		blobBaseFeeScalar,
	)

	gasparams, err := extractL1GasParams(config, zeroTime, data)
	require.NoError(t, err)
	costFunc := gasparams.costFunc

	c, g := costFunc(emptyTx.RollupCostData())

	require.Equal(t, minimumFjordGas, g)
	require.Equal(t, fjordFee, c)
}

func TestExtractIsthmusGasParams(t *testing.T) {
	zeroTime := uint64(0)
	// create a config where isthmus is active
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &zeroTime,
		EcotoneTime:  &zeroTime,
		FjordTime:    &zeroTime,
		HoloceneTime: &zeroTime,
		IsthmusTime:  &zeroTime,
	}
	require.True(t, config.IsOptimismIsthmus(zeroTime))

	data := getIsthmusL1Attributes(
		baseFee,
		blobBaseFee,
		baseFeeScalar,
		blobBaseFeeScalar,
		operatorFeeScalar,
		operatorFeeConstant,
	)

	gasparams, err := extractL1GasParams(config, zeroTime, data)
	require.NoError(t, err)
	costFunc := gasparams.costFunc

	c, g := costFunc(emptyTx.RollupCostData())

	require.Equal(t, minimumFjordGas, g)
	require.Equal(t, fjordFee, c)
	require.Equal(t, operatorFeeScalar.Uint64(), uint64(*gasparams.operatorFeeScalar))
	require.Equal(t, operatorFeeConstant.Uint64(), *gasparams.operatorFeeConstant)
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

	data := getBedrockL1Attributes(baseFee, overhead, scalar)

	gasparams, err := extractL1GasParams(config, zeroTime, data)
	require.NoError(t, err)
	oldCostFunc := gasparams.costFunc
	c, g := oldCostFunc(emptyTx.RollupCostData())
	require.Equal(t, regolithGas, g)
	require.Equal(t, regolithFee, c)
}

func getBedrockL1Attributes(baseFee, overhead, scalar *big.Int) []byte {
	uint256 := make([]byte, 32)
	ignored := big.NewInt(1234)
	data := []byte{}
	data = append(data, BedrockL1AttributesSelector...)
	data = append(data, ignored.FillBytes(uint256)...)  // arg 0
	data = append(data, ignored.FillBytes(uint256)...)  // arg 1
	data = append(data, baseFee.FillBytes(uint256)...)  // arg 2
	data = append(data, ignored.FillBytes(uint256)...)  // arg 3
	data = append(data, ignored.FillBytes(uint256)...)  // arg 4
	data = append(data, ignored.FillBytes(uint256)...)  // arg 5
	data = append(data, overhead.FillBytes(uint256)...) // arg 6
	data = append(data, scalar.FillBytes(uint256)...)   // arg 7
	return data
}

func getEcotoneL1Attributes(baseFee, blobBaseFee, baseFeeScalar, blobBaseFeeScalar *big.Int) []byte {
	ignored := big.NewInt(1234)
	data := []byte{}
	uint256Slice := make([]byte, 32)
	uint64Slice := make([]byte, 8)
	uint32Slice := make([]byte, 4)
	data = append(data, EcotoneL1AttributesSelector...)
	data = append(data, baseFeeScalar.FillBytes(uint32Slice)...)
	data = append(data, blobBaseFeeScalar.FillBytes(uint32Slice)...)
	data = append(data, ignored.FillBytes(uint64Slice)...)
	data = append(data, ignored.FillBytes(uint64Slice)...)
	data = append(data, ignored.FillBytes(uint64Slice)...)
	data = append(data, baseFee.FillBytes(uint256Slice)...)
	data = append(data, blobBaseFee.FillBytes(uint256Slice)...)
	data = append(data, ignored.FillBytes(uint256Slice)...)
	data = append(data, ignored.FillBytes(uint256Slice)...)
	return data
}

func getIsthmusL1Attributes(baseFee, blobBaseFee, baseFeeScalar, blobBaseFeeScalar, operatorFeeScalar, operatorFeeConstant *big.Int) []byte {
	ignored := big.NewInt(1234)
	data := []byte{}
	uint256Slice := make([]byte, 32)
	uint64Slice := make([]byte, 8)
	uint32Slice := make([]byte, 4)
	data = append(data, IsthmusL1AttributesSelector...)
	data = append(data, baseFeeScalar.FillBytes(uint32Slice)...)
	data = append(data, blobBaseFeeScalar.FillBytes(uint32Slice)...)
	data = append(data, ignored.FillBytes(uint64Slice)...)
	data = append(data, ignored.FillBytes(uint64Slice)...)
	data = append(data, ignored.FillBytes(uint64Slice)...)
	data = append(data, baseFee.FillBytes(uint256Slice)...)
	data = append(data, blobBaseFee.FillBytes(uint256Slice)...)
	data = append(data, ignored.FillBytes(uint256Slice)...)
	data = append(data, ignored.FillBytes(uint256Slice)...)
	data = append(data, operatorFeeScalar.FillBytes(uint32Slice)...)
	data = append(data, operatorFeeConstant.FillBytes(uint64Slice)...)
	return data
}

type testStateGetter struct {
	baseFee, blobBaseFee, overhead, scalar *big.Int
	baseFeeScalar, blobBaseFeeScalar       uint32
	operatorFeeScalar                      uint32
	operatorFeeConstant                    uint64
}

func (sg *testStateGetter) GetState(addr common.Address, slot common.Hash) common.Hash {
	buf := common.Hash{}
	switch slot {
	case L1BaseFeeSlot:
		sg.baseFee.FillBytes(buf[:])
	case OverheadSlot:
		sg.overhead.FillBytes(buf[:])
	case ScalarSlot:
		sg.scalar.FillBytes(buf[:])
	case L1BlobBaseFeeSlot:
		sg.blobBaseFee.FillBytes(buf[:])
	case L1FeeScalarsSlot:
		// fetch Ecotone fee scalars
		offset := scalarSectionStart
		binary.BigEndian.PutUint32(buf[offset:offset+4], sg.baseFeeScalar)
		binary.BigEndian.PutUint32(buf[offset+4:offset+8], sg.blobBaseFeeScalar)
	case OperatorFeeParamsSlot:
		// fetch operator fee scalars
		binary.BigEndian.PutUint32(buf[20:24], sg.operatorFeeScalar)
		binary.BigEndian.PutUint64(buf[24:32], sg.operatorFeeConstant)
	default:
		panic("unknown slot")
	}
	return buf
}

// TestNewL1CostFunc tests that the appropriate cost function is selected based on the
// configuration and statedb values.
func TestNewL1CostFunc(t *testing.T) {
	time := uint64(10)
	timeInFuture := uint64(20)
	config := &params.ChainConfig{
		Optimism: params.OptimismTestConfig.Optimism,
	}
	statedb := &testStateGetter{
		baseFee:           baseFee,
		overhead:          overhead,
		scalar:            scalar,
		blobBaseFee:       blobBaseFee,
		baseFeeScalar:     uint32(baseFeeScalar.Uint64()),
		blobBaseFeeScalar: uint32(blobBaseFeeScalar.Uint64()),
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

	// emptyTx fee w/ fjord config should be the fjord fee
	config.FjordTime = &time
	costFunc = NewL1CostFunc(config, statedb)
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, fjordFee, fee)

	// emptyTx fee w/ ecotone config, but simulate first ecotone block by blowing away the ecotone
	// params. Should result in regolith fee.
	config.FjordTime = &timeInFuture
	statedb.baseFeeScalar = 0
	statedb.blobBaseFeeScalar = 0
	statedb.blobBaseFee = new(big.Int)
	costFunc = NewL1CostFunc(config, statedb)
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, regolithFee, fee)

	// emptyTx fee w/ fjord config, but simulate first ecotone block by blowing away the ecotone
	// params. Should result in regolith fee.
	config.EcotoneTime = &time
	config.FjordTime = &time
	statedb.baseFeeScalar = 0
	statedb.blobBaseFeeScalar = 0
	statedb.blobBaseFee = new(big.Int)
	costFunc = NewL1CostFunc(config, statedb)
	fee = costFunc(emptyTx.RollupCostData(), time)
	require.NotNil(t, fee)
	require.Equal(t, regolithFee, fee)
}

// TestNewL1CostFunc tests that the appropriate cost function is selected based on the
// configuration and statedb values.
func TestNewOperatorCostFunc(t *testing.T) {
	time := uint64(10)
	config := &params.ChainConfig{
		Optimism: params.OptimismTestConfig.Optimism,
	}
	statedb := &testStateGetter{
		baseFee:             baseFee,
		overhead:            overhead,
		scalar:              scalar,
		blobBaseFee:         blobBaseFee,
		baseFeeScalar:       uint32(baseFeeScalar.Uint64()),
		blobBaseFeeScalar:   uint32(blobBaseFeeScalar.Uint64()),
		operatorFeeScalar:   uint32(operatorFeeScalar.Uint64()),
		operatorFeeConstant: operatorFeeConstant.Uint64(),
	}

	// emptyTx fee w/ fjord config, operator fee should be 0
	config.FjordTime = &time
	costFunc := NewOperatorCostFunc(config, statedb)
	fee := costFunc(bedrockGas.Uint64(), time)
	require.NotNil(t, fee)
	require.Equal(t, uint256.NewInt(0), fee)

	// emptyTx fee w/ isthmus config should be not 0
	config.IsthmusTime = &time
	costFunc = NewOperatorCostFunc(config, statedb)
	fee = costFunc(bedrockGas.Uint64(), time)
	require.NotNil(t, fee)
	require.Equal(t, ithmusOperatorFee, fee)

	// emptyTx fee w/ jovian config should be not 0
	config.JovianTime = &time
	costFunc = NewOperatorCostFunc(config, statedb)
	fee = costFunc(bedrockGas.Uint64(), time)
	require.NotNil(t, fee)
	require.Equal(t, jovianOperatorFee, fee)
}

func TestFlzCompressLen(t *testing.T) {
	var (
		emptyTxBytes, _   = emptyTx.MarshalBinary()
		contractCallTxStr = "02f901550a758302df1483be21b88304743f94f8" +
			"0e51afb613d764fa61751affd3313c190a86bb870151bd62fd12adb8" +
			"e41ef24f3f0000000000000000000000000000000000000000000000" +
			"00000000000000006e000000000000000000000000af88d065e77c8c" +
			"c2239327c5edb3a432268e5831000000000000000000000000000000" +
			"000000000000000000000000000003c1e50000000000000000000000" +
			"00000000000000000000000000000000000000000000000000000000" +
			"000000000000000000000000000000000000000000000000a0000000" +
			"00000000000000000000000000000000000000000000000000000000" +
			"148c89ed219d02f1a5be012c689b4f5b731827bebe00000000000000" +
			"0000000000c001a033fd89cb37c31b2cba46b6466e040c61fc9b2a36" +
			"75a7f5f493ebd5ad77c497f8a07cdf65680e238392693019b4092f61" +
			"0222e71b7cec06449cb922b93b6a12744e"
		contractCallTx, _ = hex.DecodeString(contractCallTxStr)
	)

	testCases := []struct {
		input       []byte
		expectedLen uint32
	}{
		// empty input
		{[]byte{}, 0},
		// all 1 inputs
		{bytes.Repeat([]byte{1}, 1000), 21},
		// all 0 inputs
		{make([]byte, 1000), 21},
		// empty tx input
		{emptyTxBytes, 31},
		// contract call tx: https://optimistic.etherscan.io/tx/0x8eb9dd4eb6d33f4dc25fb015919e4b1e9f7542f9b0322bf6622e268cd116b594
		{contractCallTx, 202},
	}

	for _, tc := range testCases {
		output := FlzCompressLen(tc.input)
		require.Equal(t, tc.expectedLen, output)
	}
}

// copy of emptyTx with non-zero gas
var emptyTxWithGas = NewTransaction(
	0,
	common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
	big.NewInt(0), bedrockGas.Uint64(), big.NewInt(0),
	nil,
)

// TestTotalRollupCostFunc tests that the total rollup cost function correctly
// combines the L1 cost and operator cost.
func TestTotalRollupCostFunc(t *testing.T) {
	zero := uint64(0)
	isthmusTime := uint64(10)
	jovianTime := uint64(20)
	config := &params.ChainConfig{
		Optimism:     params.OptimismTestConfig.Optimism,
		RegolithTime: &zero,
		EcotoneTime:  &zero,
		FjordTime:    &zero,
		HoloceneTime: &zero,
		IsthmusTime:  &isthmusTime,
		JovianTime:   &jovianTime,
	}
	statedb := &testStateGetter{
		baseFee:             baseFee,
		overhead:            overhead,
		scalar:              scalar,
		blobBaseFee:         blobBaseFee,
		baseFeeScalar:       uint32(baseFeeScalar.Uint64()),
		blobBaseFeeScalar:   uint32(blobBaseFeeScalar.Uint64()),
		operatorFeeScalar:   uint32(operatorFeeScalar.Uint64()),
		operatorFeeConstant: operatorFeeConstant.Uint64(),
	}

	costFunc := NewTotalRollupCostFunc(config, statedb)

	// Pre-Isthmus: only L1 cost
	cost := costFunc(emptyTxWithGas, isthmusTime-1)
	require.NotNil(t, cost)
	expCost := uint256.MustFromBig(fjordFee)
	require.Equal(t, expCost, cost, "pre-Isthmus total rollup cost should only contain L1 cost")

	// Isthmus: L1 cost + Isthmus operator cost
	cost = costFunc(emptyTxWithGas, isthmusTime+1)
	require.NotNil(t, cost)
	expCost = uint256.MustFromBig(fjordFee)
	expCost.Add(expCost, ithmusOperatorFee)
	require.Equal(t, expCost, cost, "Isthmus total rollup cost should contain L1 cost and Isthmus operator cost")

	// Jovian: L1 cost + fixed operator cost
	cost = costFunc(emptyTxWithGas, jovianTime+1)
	require.NotNil(t, cost)
	expCost = uint256.MustFromBig(fjordFee)
	expCost.Add(expCost, jovianOperatorFee)
	require.Equal(t, expCost, cost, "Jovian total rollup cost should contain L1 cost and Jovian operator cost")
}
