// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abigen

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// MockSortedOraclesMetaData contains all meta data concerning the MockSortedOracles contract.
var MockSortedOraclesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"DENOMINATOR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"addOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isOldestReportExpired\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"medianRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"medianTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"numRates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"numTimestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"numerators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"removeExpiredReports\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"removeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// MockSortedOraclesABI is the input ABI used to generate the binding from.
// Deprecated: Use MockSortedOraclesMetaData.ABI instead.
var MockSortedOraclesABI = MockSortedOraclesMetaData.ABI

// MockSortedOracles is an auto generated Go binding around an Ethereum contract.
type MockSortedOracles struct {
	MockSortedOraclesCaller     // Read-only binding to the contract
	MockSortedOraclesTransactor // Write-only binding to the contract
	MockSortedOraclesFilterer   // Log filterer for contract events
}

// MockSortedOraclesCaller is an auto generated read-only Go binding around an Ethereum contract.
type MockSortedOraclesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockSortedOraclesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MockSortedOraclesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockSortedOraclesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MockSortedOraclesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MockSortedOraclesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MockSortedOraclesSession struct {
	Contract     *MockSortedOracles // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// MockSortedOraclesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MockSortedOraclesCallerSession struct {
	Contract *MockSortedOraclesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// MockSortedOraclesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MockSortedOraclesTransactorSession struct {
	Contract     *MockSortedOraclesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// MockSortedOraclesRaw is an auto generated low-level Go binding around an Ethereum contract.
type MockSortedOraclesRaw struct {
	Contract *MockSortedOracles // Generic contract binding to access the raw methods on
}

// MockSortedOraclesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MockSortedOraclesCallerRaw struct {
	Contract *MockSortedOraclesCaller // Generic read-only contract binding to access the raw methods on
}

// MockSortedOraclesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MockSortedOraclesTransactorRaw struct {
	Contract *MockSortedOraclesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMockSortedOracles creates a new instance of MockSortedOracles, bound to a specific deployed contract.
func NewMockSortedOracles(address common.Address, backend bind.ContractBackend) (*MockSortedOracles, error) {
	contract, err := bindMockSortedOracles(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MockSortedOracles{MockSortedOraclesCaller: MockSortedOraclesCaller{contract: contract}, MockSortedOraclesTransactor: MockSortedOraclesTransactor{contract: contract}, MockSortedOraclesFilterer: MockSortedOraclesFilterer{contract: contract}}, nil
}

// NewMockSortedOraclesCaller creates a new read-only instance of MockSortedOracles, bound to a specific deployed contract.
func NewMockSortedOraclesCaller(address common.Address, caller bind.ContractCaller) (*MockSortedOraclesCaller, error) {
	contract, err := bindMockSortedOracles(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MockSortedOraclesCaller{contract: contract}, nil
}

// NewMockSortedOraclesTransactor creates a new write-only instance of MockSortedOracles, bound to a specific deployed contract.
func NewMockSortedOraclesTransactor(address common.Address, transactor bind.ContractTransactor) (*MockSortedOraclesTransactor, error) {
	contract, err := bindMockSortedOracles(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MockSortedOraclesTransactor{contract: contract}, nil
}

// NewMockSortedOraclesFilterer creates a new log filterer instance of MockSortedOracles, bound to a specific deployed contract.
func NewMockSortedOraclesFilterer(address common.Address, filterer bind.ContractFilterer) (*MockSortedOraclesFilterer, error) {
	contract, err := bindMockSortedOracles(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MockSortedOraclesFilterer{contract: contract}, nil
}

// bindMockSortedOracles binds a generic wrapper to an already deployed contract.
func bindMockSortedOracles(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := MockSortedOraclesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockSortedOracles *MockSortedOraclesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockSortedOracles.Contract.MockSortedOraclesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockSortedOracles *MockSortedOraclesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.MockSortedOraclesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockSortedOracles *MockSortedOraclesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.MockSortedOraclesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MockSortedOracles *MockSortedOraclesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MockSortedOracles.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MockSortedOracles *MockSortedOraclesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MockSortedOracles *MockSortedOraclesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.contract.Transact(opts, method, params...)
}

// DENOMINATOR is a free data retrieval call binding the contract method 0x918f8674.
//
// Solidity: function DENOMINATOR() view returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCaller) DENOMINATOR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "DENOMINATOR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DENOMINATOR is a free data retrieval call binding the contract method 0x918f8674.
//
// Solidity: function DENOMINATOR() view returns(uint256)
func (_MockSortedOracles *MockSortedOraclesSession) DENOMINATOR() (*big.Int, error) {
	return _MockSortedOracles.Contract.DENOMINATOR(&_MockSortedOracles.CallOpts)
}

// DENOMINATOR is a free data retrieval call binding the contract method 0x918f8674.
//
// Solidity: function DENOMINATOR() view returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCallerSession) DENOMINATOR() (*big.Int, error) {
	return _MockSortedOracles.Contract.DENOMINATOR(&_MockSortedOracles.CallOpts)
}

// IsOldestReportExpired is a free data retrieval call binding the contract method 0xffe736bf.
//
// Solidity: function isOldestReportExpired(address ) pure returns(bool, address)
func (_MockSortedOracles *MockSortedOraclesCaller) IsOldestReportExpired(opts *bind.CallOpts, arg0 common.Address) (bool, common.Address, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "isOldestReportExpired", arg0)

	if err != nil {
		return *new(bool), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return out0, out1, err

}

// IsOldestReportExpired is a free data retrieval call binding the contract method 0xffe736bf.
//
// Solidity: function isOldestReportExpired(address ) pure returns(bool, address)
func (_MockSortedOracles *MockSortedOraclesSession) IsOldestReportExpired(arg0 common.Address) (bool, common.Address, error) {
	return _MockSortedOracles.Contract.IsOldestReportExpired(&_MockSortedOracles.CallOpts, arg0)
}

// IsOldestReportExpired is a free data retrieval call binding the contract method 0xffe736bf.
//
// Solidity: function isOldestReportExpired(address ) pure returns(bool, address)
func (_MockSortedOracles *MockSortedOraclesCallerSession) IsOldestReportExpired(arg0 common.Address) (bool, common.Address, error) {
	return _MockSortedOracles.Contract.IsOldestReportExpired(&_MockSortedOracles.CallOpts, arg0)
}

// MedianRate is a free data retrieval call binding the contract method 0xef90e1b0.
//
// Solidity: function medianRate(address token) pure returns(uint256, uint256)
func (_MockSortedOracles *MockSortedOraclesCaller) MedianRate(opts *bind.CallOpts, token common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "medianRate", token)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// MedianRate is a free data retrieval call binding the contract method 0xef90e1b0.
//
// Solidity: function medianRate(address token) pure returns(uint256, uint256)
func (_MockSortedOracles *MockSortedOraclesSession) MedianRate(token common.Address) (*big.Int, *big.Int, error) {
	return _MockSortedOracles.Contract.MedianRate(&_MockSortedOracles.CallOpts, token)
}

// MedianRate is a free data retrieval call binding the contract method 0xef90e1b0.
//
// Solidity: function medianRate(address token) pure returns(uint256, uint256)
func (_MockSortedOracles *MockSortedOraclesCallerSession) MedianRate(token common.Address) (*big.Int, *big.Int, error) {
	return _MockSortedOracles.Contract.MedianRate(&_MockSortedOracles.CallOpts, token)
}

// MedianTimestamp is a free data retrieval call binding the contract method 0x071b48fc.
//
// Solidity: function medianTimestamp(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCaller) MedianTimestamp(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "medianTimestamp", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MedianTimestamp is a free data retrieval call binding the contract method 0x071b48fc.
//
// Solidity: function medianTimestamp(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesSession) MedianTimestamp(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.MedianTimestamp(&_MockSortedOracles.CallOpts, arg0)
}

// MedianTimestamp is a free data retrieval call binding the contract method 0x071b48fc.
//
// Solidity: function medianTimestamp(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCallerSession) MedianTimestamp(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.MedianTimestamp(&_MockSortedOracles.CallOpts, arg0)
}

// NumRates is a free data retrieval call binding the contract method 0xbbc66a94.
//
// Solidity: function numRates(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCaller) NumRates(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "numRates", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumRates is a free data retrieval call binding the contract method 0xbbc66a94.
//
// Solidity: function numRates(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesSession) NumRates(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.NumRates(&_MockSortedOracles.CallOpts, arg0)
}

// NumRates is a free data retrieval call binding the contract method 0xbbc66a94.
//
// Solidity: function numRates(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCallerSession) NumRates(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.NumRates(&_MockSortedOracles.CallOpts, arg0)
}

// NumTimestamps is a free data retrieval call binding the contract method 0x6dd6ef0c.
//
// Solidity: function numTimestamps(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCaller) NumTimestamps(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "numTimestamps", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumTimestamps is a free data retrieval call binding the contract method 0x6dd6ef0c.
//
// Solidity: function numTimestamps(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesSession) NumTimestamps(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.NumTimestamps(&_MockSortedOracles.CallOpts, arg0)
}

// NumTimestamps is a free data retrieval call binding the contract method 0x6dd6ef0c.
//
// Solidity: function numTimestamps(address ) pure returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCallerSession) NumTimestamps(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.NumTimestamps(&_MockSortedOracles.CallOpts, arg0)
}

// Numerators is a free data retrieval call binding the contract method 0xf7ca6963.
//
// Solidity: function numerators(address ) view returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCaller) Numerators(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _MockSortedOracles.contract.Call(opts, &out, "numerators", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Numerators is a free data retrieval call binding the contract method 0xf7ca6963.
//
// Solidity: function numerators(address ) view returns(uint256)
func (_MockSortedOracles *MockSortedOraclesSession) Numerators(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.Numerators(&_MockSortedOracles.CallOpts, arg0)
}

// Numerators is a free data retrieval call binding the contract method 0xf7ca6963.
//
// Solidity: function numerators(address ) view returns(uint256)
func (_MockSortedOracles *MockSortedOraclesCallerSession) Numerators(arg0 common.Address) (*big.Int, error) {
	return _MockSortedOracles.Contract.Numerators(&_MockSortedOracles.CallOpts, arg0)
}

// AddOracle is a paid mutator transaction binding the contract method 0xf0ca4adb.
//
// Solidity: function addOracle(address , address ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactor) AddOracle(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address) (*types.Transaction, error) {
	return _MockSortedOracles.contract.Transact(opts, "addOracle", arg0, arg1)
}

// AddOracle is a paid mutator transaction binding the contract method 0xf0ca4adb.
//
// Solidity: function addOracle(address , address ) returns()
func (_MockSortedOracles *MockSortedOraclesSession) AddOracle(arg0 common.Address, arg1 common.Address) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.AddOracle(&_MockSortedOracles.TransactOpts, arg0, arg1)
}

// AddOracle is a paid mutator transaction binding the contract method 0xf0ca4adb.
//
// Solidity: function addOracle(address , address ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactorSession) AddOracle(arg0 common.Address, arg1 common.Address) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.AddOracle(&_MockSortedOracles.TransactOpts, arg0, arg1)
}

// RemoveExpiredReports is a paid mutator transaction binding the contract method 0xdd34ca3b.
//
// Solidity: function removeExpiredReports(address , uint256 ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactor) RemoveExpiredReports(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _MockSortedOracles.contract.Transact(opts, "removeExpiredReports", arg0, arg1)
}

// RemoveExpiredReports is a paid mutator transaction binding the contract method 0xdd34ca3b.
//
// Solidity: function removeExpiredReports(address , uint256 ) returns()
func (_MockSortedOracles *MockSortedOraclesSession) RemoveExpiredReports(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.RemoveExpiredReports(&_MockSortedOracles.TransactOpts, arg0, arg1)
}

// RemoveExpiredReports is a paid mutator transaction binding the contract method 0xdd34ca3b.
//
// Solidity: function removeExpiredReports(address , uint256 ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactorSession) RemoveExpiredReports(arg0 common.Address, arg1 *big.Int) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.RemoveExpiredReports(&_MockSortedOracles.TransactOpts, arg0, arg1)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0x53a57297.
//
// Solidity: function removeOracle(address , address , uint256 ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactor) RemoveOracle(opts *bind.TransactOpts, arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _MockSortedOracles.contract.Transact(opts, "removeOracle", arg0, arg1, arg2)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0x53a57297.
//
// Solidity: function removeOracle(address , address , uint256 ) returns()
func (_MockSortedOracles *MockSortedOraclesSession) RemoveOracle(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.RemoveOracle(&_MockSortedOracles.TransactOpts, arg0, arg1, arg2)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0x53a57297.
//
// Solidity: function removeOracle(address , address , uint256 ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactorSession) RemoveOracle(arg0 common.Address, arg1 common.Address, arg2 *big.Int) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.RemoveOracle(&_MockSortedOracles.TransactOpts, arg0, arg1, arg2)
}

// Report is a paid mutator transaction binding the contract method 0x80e50744.
//
// Solidity: function report(address , uint256 , address , address ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactor) Report(opts *bind.TransactOpts, arg0 common.Address, arg1 *big.Int, arg2 common.Address, arg3 common.Address) (*types.Transaction, error) {
	return _MockSortedOracles.contract.Transact(opts, "report", arg0, arg1, arg2, arg3)
}

// Report is a paid mutator transaction binding the contract method 0x80e50744.
//
// Solidity: function report(address , uint256 , address , address ) returns()
func (_MockSortedOracles *MockSortedOraclesSession) Report(arg0 common.Address, arg1 *big.Int, arg2 common.Address, arg3 common.Address) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.Report(&_MockSortedOracles.TransactOpts, arg0, arg1, arg2, arg3)
}

// Report is a paid mutator transaction binding the contract method 0x80e50744.
//
// Solidity: function report(address , uint256 , address , address ) returns()
func (_MockSortedOracles *MockSortedOraclesTransactorSession) Report(arg0 common.Address, arg1 *big.Int, arg2 common.Address, arg3 common.Address) (*types.Transaction, error) {
	return _MockSortedOracles.Contract.Report(&_MockSortedOracles.TransactOpts, arg0, arg1, arg2, arg3)
}
