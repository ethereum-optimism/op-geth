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

// CurrencyConfig is an auto generated low-level Go binding around an user-defined struct.
type CurrencyConfig struct {
	CurrencyIdentifier common.Address
	Oracle             common.Address
	IntrinsicGas       *big.Int
}

// FeeCurrencyDirectoryMetaData contains all meta data concerning the FeeCurrencyDirectory contract.
var FeeCurrencyDirectoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"test\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"currencies\",\"inputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"currencyIdentifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"intrinsicGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrencies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrencyConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structCurrencyConfig\",\"components\":[{\"name\":\"currencyIdentifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"intrinsicGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExchangeRate\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"numerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getVersionNumber\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"pure\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"initialized\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"removeCurrencies\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"renounceOwnership\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setCurrencyConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"currencyIdentifier\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"intrinsicGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"transferOwnership\",\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"OwnershipTransferred\",\"inputs\":[{\"name\":\"previousOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"newOwner\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false}]",
}

// FeeCurrencyDirectoryABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeCurrencyDirectoryMetaData.ABI instead.
var FeeCurrencyDirectoryABI = FeeCurrencyDirectoryMetaData.ABI

// FeeCurrencyDirectory is an auto generated Go binding around an Ethereum contract.
type FeeCurrencyDirectory struct {
	FeeCurrencyDirectoryCaller     // Read-only binding to the contract
	FeeCurrencyDirectoryTransactor // Write-only binding to the contract
	FeeCurrencyDirectoryFilterer   // Log filterer for contract events
}

// FeeCurrencyDirectoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeCurrencyDirectoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeCurrencyDirectoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeCurrencyDirectoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeCurrencyDirectoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeCurrencyDirectoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeCurrencyDirectorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeCurrencyDirectorySession struct {
	Contract     *FeeCurrencyDirectory // Generic contract binding to set the session for
	CallOpts     bind.CallOpts         // Call options to use throughout this session
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// FeeCurrencyDirectoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeCurrencyDirectoryCallerSession struct {
	Contract *FeeCurrencyDirectoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts               // Call options to use throughout this session
}

// FeeCurrencyDirectoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeCurrencyDirectoryTransactorSession struct {
	Contract     *FeeCurrencyDirectoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts               // Transaction auth options to use throughout this session
}

// FeeCurrencyDirectoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeCurrencyDirectoryRaw struct {
	Contract *FeeCurrencyDirectory // Generic contract binding to access the raw methods on
}

// FeeCurrencyDirectoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeCurrencyDirectoryCallerRaw struct {
	Contract *FeeCurrencyDirectoryCaller // Generic read-only contract binding to access the raw methods on
}

// FeeCurrencyDirectoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeCurrencyDirectoryTransactorRaw struct {
	Contract *FeeCurrencyDirectoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeCurrencyDirectory creates a new instance of FeeCurrencyDirectory, bound to a specific deployed contract.
func NewFeeCurrencyDirectory(address common.Address, backend bind.ContractBackend) (*FeeCurrencyDirectory, error) {
	contract, err := bindFeeCurrencyDirectory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyDirectory{FeeCurrencyDirectoryCaller: FeeCurrencyDirectoryCaller{contract: contract}, FeeCurrencyDirectoryTransactor: FeeCurrencyDirectoryTransactor{contract: contract}, FeeCurrencyDirectoryFilterer: FeeCurrencyDirectoryFilterer{contract: contract}}, nil
}

// NewFeeCurrencyDirectoryCaller creates a new read-only instance of FeeCurrencyDirectory, bound to a specific deployed contract.
func NewFeeCurrencyDirectoryCaller(address common.Address, caller bind.ContractCaller) (*FeeCurrencyDirectoryCaller, error) {
	contract, err := bindFeeCurrencyDirectory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyDirectoryCaller{contract: contract}, nil
}

// NewFeeCurrencyDirectoryTransactor creates a new write-only instance of FeeCurrencyDirectory, bound to a specific deployed contract.
func NewFeeCurrencyDirectoryTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeCurrencyDirectoryTransactor, error) {
	contract, err := bindFeeCurrencyDirectory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyDirectoryTransactor{contract: contract}, nil
}

// NewFeeCurrencyDirectoryFilterer creates a new log filterer instance of FeeCurrencyDirectory, bound to a specific deployed contract.
func NewFeeCurrencyDirectoryFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeCurrencyDirectoryFilterer, error) {
	contract, err := bindFeeCurrencyDirectory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyDirectoryFilterer{contract: contract}, nil
}

// bindFeeCurrencyDirectory binds a generic wrapper to an already deployed contract.
func bindFeeCurrencyDirectory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeCurrencyDirectoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeCurrencyDirectory.Contract.FeeCurrencyDirectoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.FeeCurrencyDirectoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.FeeCurrencyDirectoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeCurrencyDirectory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.contract.Transact(opts, method, params...)
}

// Currencies is a free data retrieval call binding the contract method 0x6036cba3.
//
// Solidity: function currencies(address ) view returns(address currencyIdentifier, address oracle, uint256 intrinsicGas)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) Currencies(opts *bind.CallOpts, arg0 common.Address) (struct {
	CurrencyIdentifier common.Address
	Oracle             common.Address
	IntrinsicGas       *big.Int
}, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "currencies", arg0)

	outstruct := new(struct {
		CurrencyIdentifier common.Address
		Oracle             common.Address
		IntrinsicGas       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.CurrencyIdentifier = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Oracle = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.IntrinsicGas = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Currencies is a free data retrieval call binding the contract method 0x6036cba3.
//
// Solidity: function currencies(address ) view returns(address currencyIdentifier, address oracle, uint256 intrinsicGas)
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) Currencies(arg0 common.Address) (struct {
	CurrencyIdentifier common.Address
	Oracle             common.Address
	IntrinsicGas       *big.Int
}, error) {
	return _FeeCurrencyDirectory.Contract.Currencies(&_FeeCurrencyDirectory.CallOpts, arg0)
}

// Currencies is a free data retrieval call binding the contract method 0x6036cba3.
//
// Solidity: function currencies(address ) view returns(address currencyIdentifier, address oracle, uint256 intrinsicGas)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) Currencies(arg0 common.Address) (struct {
	CurrencyIdentifier common.Address
	Oracle             common.Address
	IntrinsicGas       *big.Int
}, error) {
	return _FeeCurrencyDirectory.Contract.Currencies(&_FeeCurrencyDirectory.CallOpts, arg0)
}

// GetCurrencies is a free data retrieval call binding the contract method 0x61c661de.
//
// Solidity: function getCurrencies() view returns(address[])
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) GetCurrencies(opts *bind.CallOpts) ([]common.Address, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "getCurrencies")

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetCurrencies is a free data retrieval call binding the contract method 0x61c661de.
//
// Solidity: function getCurrencies() view returns(address[])
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) GetCurrencies() ([]common.Address, error) {
	return _FeeCurrencyDirectory.Contract.GetCurrencies(&_FeeCurrencyDirectory.CallOpts)
}

// GetCurrencies is a free data retrieval call binding the contract method 0x61c661de.
//
// Solidity: function getCurrencies() view returns(address[])
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) GetCurrencies() ([]common.Address, error) {
	return _FeeCurrencyDirectory.Contract.GetCurrencies(&_FeeCurrencyDirectory.CallOpts)
}

// GetCurrencyConfig is a free data retrieval call binding the contract method 0xeab43d97.
//
// Solidity: function getCurrencyConfig(address token) view returns((address,address,uint256))
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) GetCurrencyConfig(opts *bind.CallOpts, token common.Address) (CurrencyConfig, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "getCurrencyConfig", token)

	if err != nil {
		return *new(CurrencyConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(CurrencyConfig)).(*CurrencyConfig)

	return out0, err

}

// GetCurrencyConfig is a free data retrieval call binding the contract method 0xeab43d97.
//
// Solidity: function getCurrencyConfig(address token) view returns((address,address,uint256))
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) GetCurrencyConfig(token common.Address) (CurrencyConfig, error) {
	return _FeeCurrencyDirectory.Contract.GetCurrencyConfig(&_FeeCurrencyDirectory.CallOpts, token)
}

// GetCurrencyConfig is a free data retrieval call binding the contract method 0xeab43d97.
//
// Solidity: function getCurrencyConfig(address token) view returns((address,address,uint256))
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) GetCurrencyConfig(token common.Address) (CurrencyConfig, error) {
	return _FeeCurrencyDirectory.Contract.GetCurrencyConfig(&_FeeCurrencyDirectory.CallOpts, token)
}

// GetExchangeRate is a free data retrieval call binding the contract method 0xefb7601d.
//
// Solidity: function getExchangeRate(address token) view returns(uint256 numerator, uint256 denominator)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) GetExchangeRate(opts *bind.CallOpts, token common.Address) (struct {
	Numerator   *big.Int
	Denominator *big.Int
}, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "getExchangeRate", token)

	outstruct := new(struct {
		Numerator   *big.Int
		Denominator *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Numerator = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Denominator = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetExchangeRate is a free data retrieval call binding the contract method 0xefb7601d.
//
// Solidity: function getExchangeRate(address token) view returns(uint256 numerator, uint256 denominator)
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) GetExchangeRate(token common.Address) (struct {
	Numerator   *big.Int
	Denominator *big.Int
}, error) {
	return _FeeCurrencyDirectory.Contract.GetExchangeRate(&_FeeCurrencyDirectory.CallOpts, token)
}

// GetExchangeRate is a free data retrieval call binding the contract method 0xefb7601d.
//
// Solidity: function getExchangeRate(address token) view returns(uint256 numerator, uint256 denominator)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) GetExchangeRate(token common.Address) (struct {
	Numerator   *big.Int
	Denominator *big.Int
}, error) {
	return _FeeCurrencyDirectory.Contract.GetExchangeRate(&_FeeCurrencyDirectory.CallOpts, token)
}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) GetVersionNumber(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "getVersionNumber")

	if err != nil {
		return *new(*big.Int), *new(*big.Int), *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	out2 := *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	out3 := *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return out0, out1, out2, out3, err

}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) GetVersionNumber() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _FeeCurrencyDirectory.Contract.GetVersionNumber(&_FeeCurrencyDirectory.CallOpts)
}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) GetVersionNumber() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _FeeCurrencyDirectory.Contract.GetVersionNumber(&_FeeCurrencyDirectory.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) Initialized() (bool, error) {
	return _FeeCurrencyDirectory.Contract.Initialized(&_FeeCurrencyDirectory.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) Initialized() (bool, error) {
	return _FeeCurrencyDirectory.Contract.Initialized(&_FeeCurrencyDirectory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) Owner() (common.Address, error) {
	return _FeeCurrencyDirectory.Contract.Owner(&_FeeCurrencyDirectory.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) Owner() (common.Address, error) {
	return _FeeCurrencyDirectory.Contract.Owner(&_FeeCurrencyDirectory.CallOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) Initialize() (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.Initialize(&_FeeCurrencyDirectory.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorSession) Initialize() (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.Initialize(&_FeeCurrencyDirectory.TransactOpts)
}

// RemoveCurrencies is a paid mutator transaction binding the contract method 0x16be73a8.
//
// Solidity: function removeCurrencies(address token, uint256 index) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactor) RemoveCurrencies(opts *bind.TransactOpts, token common.Address, index *big.Int) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.contract.Transact(opts, "removeCurrencies", token, index)
}

// RemoveCurrencies is a paid mutator transaction binding the contract method 0x16be73a8.
//
// Solidity: function removeCurrencies(address token, uint256 index) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) RemoveCurrencies(token common.Address, index *big.Int) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.RemoveCurrencies(&_FeeCurrencyDirectory.TransactOpts, token, index)
}

// RemoveCurrencies is a paid mutator transaction binding the contract method 0x16be73a8.
//
// Solidity: function removeCurrencies(address token, uint256 index) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorSession) RemoveCurrencies(token common.Address, index *big.Int) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.RemoveCurrencies(&_FeeCurrencyDirectory.TransactOpts, token, index)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) RenounceOwnership() (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.RenounceOwnership(&_FeeCurrencyDirectory.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.RenounceOwnership(&_FeeCurrencyDirectory.TransactOpts)
}

// SetCurrencyConfig is a paid mutator transaction binding the contract method 0x9046e34a.
//
// Solidity: function setCurrencyConfig(address token, address currencyIdentifier, address oracle, uint256 intrinsicGas) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactor) SetCurrencyConfig(opts *bind.TransactOpts, token common.Address, currencyIdentifier common.Address, oracle common.Address, intrinsicGas *big.Int) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.contract.Transact(opts, "setCurrencyConfig", token, currencyIdentifier, oracle, intrinsicGas)
}

// SetCurrencyConfig is a paid mutator transaction binding the contract method 0x9046e34a.
//
// Solidity: function setCurrencyConfig(address token, address currencyIdentifier, address oracle, uint256 intrinsicGas) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) SetCurrencyConfig(token common.Address, currencyIdentifier common.Address, oracle common.Address, intrinsicGas *big.Int) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.SetCurrencyConfig(&_FeeCurrencyDirectory.TransactOpts, token, currencyIdentifier, oracle, intrinsicGas)
}

// SetCurrencyConfig is a paid mutator transaction binding the contract method 0x9046e34a.
//
// Solidity: function setCurrencyConfig(address token, address currencyIdentifier, address oracle, uint256 intrinsicGas) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorSession) SetCurrencyConfig(token common.Address, currencyIdentifier common.Address, oracle common.Address, intrinsicGas *big.Int) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.SetCurrencyConfig(&_FeeCurrencyDirectory.TransactOpts, token, currencyIdentifier, oracle, intrinsicGas)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.TransferOwnership(&_FeeCurrencyDirectory.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _FeeCurrencyDirectory.Contract.TransferOwnership(&_FeeCurrencyDirectory.TransactOpts, newOwner)
}

// FeeCurrencyDirectoryOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the FeeCurrencyDirectory contract.
type FeeCurrencyDirectoryOwnershipTransferredIterator struct {
	Event *FeeCurrencyDirectoryOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *FeeCurrencyDirectoryOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeCurrencyDirectoryOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(FeeCurrencyDirectoryOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *FeeCurrencyDirectoryOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeCurrencyDirectoryOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeCurrencyDirectoryOwnershipTransferred represents a OwnershipTransferred event raised by the FeeCurrencyDirectory contract.
type FeeCurrencyDirectoryOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*FeeCurrencyDirectoryOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FeeCurrencyDirectory.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyDirectoryOwnershipTransferredIterator{contract: _FeeCurrencyDirectory.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *FeeCurrencyDirectoryOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _FeeCurrencyDirectory.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeCurrencyDirectoryOwnershipTransferred)
				if err := _FeeCurrencyDirectory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryFilterer) ParseOwnershipTransferred(log types.Log) (*FeeCurrencyDirectoryOwnershipTransferred, error) {
	event := new(FeeCurrencyDirectoryOwnershipTransferred)
	if err := _FeeCurrencyDirectory.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
