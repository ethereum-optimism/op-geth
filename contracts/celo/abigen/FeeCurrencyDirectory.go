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

// IFeeCurrencyDirectoryCurrencyConfig is an auto generated low-level Go binding around an user-defined struct.
type IFeeCurrencyDirectoryCurrencyConfig struct {
	Oracle       common.Address
	IntrinsicGas *big.Int
}

// FeeCurrencyDirectoryMetaData contains all meta data concerning the FeeCurrencyDirectory contract.
var FeeCurrencyDirectoryMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"function\",\"name\":\"getCurrencies\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address[]\",\"internalType\":\"address[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getCurrencyConfig\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIFeeCurrencyDirectory.CurrencyConfig\",\"components\":[{\"name\":\"oracle\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"intrinsicGas\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getExchangeRate\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"numerator\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"denominator\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"}]",
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
// Solidity: function getCurrencyConfig(address token) view returns((address,uint256))
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCaller) GetCurrencyConfig(opts *bind.CallOpts, token common.Address) (IFeeCurrencyDirectoryCurrencyConfig, error) {
	var out []interface{}
	err := _FeeCurrencyDirectory.contract.Call(opts, &out, "getCurrencyConfig", token)

	if err != nil {
		return *new(IFeeCurrencyDirectoryCurrencyConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IFeeCurrencyDirectoryCurrencyConfig)).(*IFeeCurrencyDirectoryCurrencyConfig)

	return out0, err

}

// GetCurrencyConfig is a free data retrieval call binding the contract method 0xeab43d97.
//
// Solidity: function getCurrencyConfig(address token) view returns((address,uint256))
func (_FeeCurrencyDirectory *FeeCurrencyDirectorySession) GetCurrencyConfig(token common.Address) (IFeeCurrencyDirectoryCurrencyConfig, error) {
	return _FeeCurrencyDirectory.Contract.GetCurrencyConfig(&_FeeCurrencyDirectory.CallOpts, token)
}

// GetCurrencyConfig is a free data retrieval call binding the contract method 0xeab43d97.
//
// Solidity: function getCurrencyConfig(address token) view returns((address,uint256))
func (_FeeCurrencyDirectory *FeeCurrencyDirectoryCallerSession) GetCurrencyConfig(token common.Address) (IFeeCurrencyDirectoryCurrencyConfig, error) {
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
