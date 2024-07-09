package core

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/contracts/addresses"
	"github.com/ethereum/go-ethereum/contracts/celo"
	"github.com/ethereum/go-ethereum/crypto"
)

// Decode 0x prefixed hex string from file (including trailing newline)
func DecodeHex(hexbytes []byte) ([]byte, error) {
	// Strip 0x prefix and trailing newline
	hexbytes = hexbytes[2 : len(hexbytes)-1] // strip 0x prefix

	// Decode hex string
	bytes := make([]byte, hex.DecodedLen(len(hexbytes)))
	_, err := hex.Decode(bytes, hexbytes)
	if err != nil {
		return nil, fmt.Errorf("DecodeHex: %w", err)
	}

	return bytes, nil
}

// Calculate address in evm mapping: keccak(key ++ mapping_slot)
func CalcMapAddr(slot common.Hash, key common.Hash) common.Hash {
	return crypto.Keccak256Hash(append(key.Bytes(), slot.Bytes()...))
}

// Increase a hash value by `i`, used for addresses in 32byte fields
func incHash(addr common.Hash, i int64) common.Hash {
	return common.BigToHash(new(big.Int).Add(addr.Big(), big.NewInt(i)))
}

var (
	DevPrivateKey, _ = crypto.HexToECDSA("2771aff413cac48d9f8c114fabddd9195a2129f3c2c436caa07e27bb7f58ead5")
	DevAddr          = common.BytesToAddress(DevAddr32.Bytes())
	DevAddr32        = common.HexToHash("0x42cf1bbc38BaAA3c4898ce8790e21eD2738c6A4a")

	DevFeeCurrencyAddr  = common.HexToAddress("0xce16") // worth twice as much as native CELO
	DevFeeCurrencyAddr2 = common.HexToAddress("0xce17") // worth half as much as native CELO
	DevBalance, _       = new(big.Int).SetString("100000000000000000000", 10)
	rateNumerator, _    = new(big.Int).SetString("2000000000000000000000000", 10)
	rateNumerator2, _   = new(big.Int).SetString("500000000000000000000000", 10)
	rateDenominator, _  = new(big.Int).SetString("1000000000000000000000000", 10)
	mockOracleAddr      = common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb0001")
	mockOracleAddr2     = common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb0002")
	mockOracleAddr3     = common.HexToAddress("0xbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb0003")
	FaucetAddr          = common.HexToAddress("0xfcf982bb4015852e706100b14e21f947a5bb718e")
)

func celoGenesisAccounts(fundedAddr common.Address) GenesisAlloc {
	// Initialize Bytecodes
	celoTokenBytecode, err := DecodeHex(celo.CeloTokenBytecodeRaw)
	if err != nil {
		panic(err)
	}
	feeCurrencyBytecode, err := DecodeHex(celo.FeeCurrencyBytecodeRaw)
	if err != nil {
		panic(err)
	}
	feeCurrencyDirectoryBytecode, err := DecodeHex(celo.FeeCurrencyDirectoryBytecodeRaw)
	if err != nil {
		panic(err)
	}
	mockOracleBytecode, err := DecodeHex(celo.MockOracleBytecodeRaw)
	if err != nil {
		panic(err)
	}

	var devBalance32 common.Hash
	DevBalance.FillBytes(devBalance32[:])

	arrayAtSlot2 := crypto.Keccak256Hash(common.HexToHash("0x2").Bytes())

	faucetBalance, ok := new(big.Int).SetString("500000000000000000000000000", 10) // 500M
	if !ok {
		panic("Couldn not set faucet balance!")
	}
	genesisAccounts := map[common.Address]GenesisAccount{
		addresses.CeloTokenAddress: {
			Code:    celoTokenBytecode,
			Balance: big.NewInt(0),
		},
		DevFeeCurrencyAddr: {
			Code:    feeCurrencyBytecode,
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				CalcMapAddr(common.HexToHash("0x0"), DevAddr32):                              devBalance32, // _balances[DevAddr]
				CalcMapAddr(common.HexToHash("0x0"), common.BytesToHash(fundedAddr.Bytes())): devBalance32, // _balances[fund]
				common.HexToHash("0x2"): devBalance32, // _totalSupply
			},
		},
		DevFeeCurrencyAddr2: {
			Code:    feeCurrencyBytecode,
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				CalcMapAddr(common.HexToHash("0x0"), DevAddr32):                              devBalance32, // _balances[DevAddr]
				CalcMapAddr(common.HexToHash("0x0"), common.BytesToHash(fundedAddr.Bytes())): devBalance32, // _balances[fund]
				common.HexToHash("0x2"): devBalance32, // _totalSupply
			},
		},
		mockOracleAddr: {
			Code:    mockOracleBytecode,
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0"): common.BigToHash(rateNumerator),
				common.HexToHash("0x1"): common.BigToHash(rateDenominator),
				common.HexToHash("0x3"): common.BytesToHash(DevFeeCurrencyAddr.Bytes()),
			},
		},
		mockOracleAddr2: {
			Code:    mockOracleBytecode,
			Balance: big.NewInt(0),
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x0"): common.BigToHash(rateNumerator2),
				common.HexToHash("0x1"): common.BigToHash(rateDenominator),
				common.HexToHash("0x3"): common.BytesToHash(DevFeeCurrencyAddr2.Bytes()),
			},
		},
		mockOracleAddr3: {
			Code:    mockOracleBytecode,
			Balance: big.NewInt(0),
			// This oracle is available for tests of contracts outside the celo_genesis, so no initialization is done at this point
		},
		DevAddr: {
			Balance: DevBalance,
		},
		FaucetAddr: {
			Balance: faucetBalance,
		},
	}

	// FeeCurrencyDirectory
	devAddrOffset1 := common.Hash{}
	copy(devAddrOffset1[11:], DevAddr.Bytes())
	feeCurrencyDirectoryStorage := map[common.Hash]common.Hash{
		// owner, slot 0 offset 1
		common.HexToHash("0x0"): devAddrOffset1,
		// add entries to currencyList at slot 2
		common.HexToHash("0x2"):  common.HexToHash("0x2"),                         // array length 2
		arrayAtSlot2:             common.BytesToHash(DevFeeCurrencyAddr.Bytes()),  // FeeCurrency
		incHash(arrayAtSlot2, 1): common.BytesToHash(DevFeeCurrencyAddr2.Bytes()), // FeeCurrency2
	}
	// add entries to currencyConfig mapping
	addFeeCurrencyToStorage(DevFeeCurrencyAddr, mockOracleAddr, feeCurrencyDirectoryStorage)
	addFeeCurrencyToStorage(DevFeeCurrencyAddr2, mockOracleAddr2, feeCurrencyDirectoryStorage)
	genesisAccounts[addresses.FeeCurrencyDirectoryAddress] = GenesisAccount{
		Code:    feeCurrencyDirectoryBytecode,
		Balance: big.NewInt(0),
		Storage: feeCurrencyDirectoryStorage,
	}

	return genesisAccounts
}

func addFeeCurrencyToStorage(feeCurrencyAddr common.Address, oracleAddr common.Address, storage map[common.Hash]common.Hash) {
	structStart := CalcMapAddr(common.HexToHash("0x1"), common.BytesToHash(feeCurrencyAddr.Bytes()))
	storage[structStart] = common.BytesToHash(oracleAddr.Bytes())          // oracle
	storage[incHash(structStart, 1)] = common.BigToHash(big.NewInt(50000)) // intrinsicGas
}
