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

// SortedOraclesMetaData contains all meta data concerning the SortedOracles contract.
var SortedOraclesMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"test\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"MedianUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"name\":\"OracleAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"name\":\"OracleRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"}],\"name\":\"OracleReportRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oracle\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"OracleReported\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reportExpiry\",\"type\":\"uint256\"}],\"name\":\"ReportExpirySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reportExpiry\",\"type\":\"uint256\"}],\"name\":\"TokenReportExpirySet\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"}],\"name\":\"addOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getOracles\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getRates\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"enumSortedLinkedListWithMedian.MedianRelation[]\",\"name\":\"\",\"type\":\"uint8[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTimestamps\",\"outputs\":[{\"internalType\":\"address[]\",\"name\":\"\",\"type\":\"address[]\"},{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"},{\"internalType\":\"enumSortedLinkedListWithMedian.MedianRelation[]\",\"name\":\"\",\"type\":\"uint8[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"getTokenReportExpirySeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getVersionNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_reportExpirySeconds\",\"type\":\"uint256\"}],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"isOldestReportExpired\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isOracle\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"medianRate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"medianTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"numRates\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"name\":\"numTimestamps\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"oracles\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"n\",\"type\":\"uint256\"}],\"name\":\"removeExpiredReports\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"oracleAddress\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"removeOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"lesserKey\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"greaterKey\",\"type\":\"address\"}],\"name\":\"report\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"reportExpirySeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_reportExpirySeconds\",\"type\":\"uint256\"}],\"name\":\"setReportExpiry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_reportExpirySeconds\",\"type\":\"uint256\"}],\"name\":\"setTokenReportExpiry\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokenReportExpirySeconds\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// SortedOraclesABI is the input ABI used to generate the binding from.
// Deprecated: Use SortedOraclesMetaData.ABI instead.
var SortedOraclesABI = SortedOraclesMetaData.ABI

// SortedOracles is an auto generated Go binding around an Ethereum contract.
type SortedOracles struct {
	SortedOraclesCaller     // Read-only binding to the contract
	SortedOraclesTransactor // Write-only binding to the contract
	SortedOraclesFilterer   // Log filterer for contract events
}

// SortedOraclesCaller is an auto generated read-only Go binding around an Ethereum contract.
type SortedOraclesCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SortedOraclesTransactor is an auto generated write-only Go binding around an Ethereum contract.
type SortedOraclesTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SortedOraclesFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type SortedOraclesFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// SortedOraclesSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type SortedOraclesSession struct {
	Contract     *SortedOracles    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// SortedOraclesCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type SortedOraclesCallerSession struct {
	Contract *SortedOraclesCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// SortedOraclesTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type SortedOraclesTransactorSession struct {
	Contract     *SortedOraclesTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// SortedOraclesRaw is an auto generated low-level Go binding around an Ethereum contract.
type SortedOraclesRaw struct {
	Contract *SortedOracles // Generic contract binding to access the raw methods on
}

// SortedOraclesCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type SortedOraclesCallerRaw struct {
	Contract *SortedOraclesCaller // Generic read-only contract binding to access the raw methods on
}

// SortedOraclesTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type SortedOraclesTransactorRaw struct {
	Contract *SortedOraclesTransactor // Generic write-only contract binding to access the raw methods on
}

// NewSortedOracles creates a new instance of SortedOracles, bound to a specific deployed contract.
func NewSortedOracles(address common.Address, backend bind.ContractBackend) (*SortedOracles, error) {
	contract, err := bindSortedOracles(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &SortedOracles{SortedOraclesCaller: SortedOraclesCaller{contract: contract}, SortedOraclesTransactor: SortedOraclesTransactor{contract: contract}, SortedOraclesFilterer: SortedOraclesFilterer{contract: contract}}, nil
}

// NewSortedOraclesCaller creates a new read-only instance of SortedOracles, bound to a specific deployed contract.
func NewSortedOraclesCaller(address common.Address, caller bind.ContractCaller) (*SortedOraclesCaller, error) {
	contract, err := bindSortedOracles(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesCaller{contract: contract}, nil
}

// NewSortedOraclesTransactor creates a new write-only instance of SortedOracles, bound to a specific deployed contract.
func NewSortedOraclesTransactor(address common.Address, transactor bind.ContractTransactor) (*SortedOraclesTransactor, error) {
	contract, err := bindSortedOracles(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesTransactor{contract: contract}, nil
}

// NewSortedOraclesFilterer creates a new log filterer instance of SortedOracles, bound to a specific deployed contract.
func NewSortedOraclesFilterer(address common.Address, filterer bind.ContractFilterer) (*SortedOraclesFilterer, error) {
	contract, err := bindSortedOracles(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesFilterer{contract: contract}, nil
}

// bindSortedOracles binds a generic wrapper to an already deployed contract.
func bindSortedOracles(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := SortedOraclesMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SortedOracles *SortedOraclesRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SortedOracles.Contract.SortedOraclesCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SortedOracles *SortedOraclesRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortedOracles.Contract.SortedOraclesTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SortedOracles *SortedOraclesRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SortedOracles.Contract.SortedOraclesTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_SortedOracles *SortedOraclesCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _SortedOracles.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_SortedOracles *SortedOraclesTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortedOracles.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_SortedOracles *SortedOraclesTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _SortedOracles.Contract.contract.Transact(opts, method, params...)
}

// GetOracles is a free data retrieval call binding the contract method 0x8e749281.
//
// Solidity: function getOracles(address token) view returns(address[])
func (_SortedOracles *SortedOraclesCaller) GetOracles(opts *bind.CallOpts, token common.Address) ([]common.Address, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "getOracles", token)

	if err != nil {
		return *new([]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)

	return out0, err

}

// GetOracles is a free data retrieval call binding the contract method 0x8e749281.
//
// Solidity: function getOracles(address token) view returns(address[])
func (_SortedOracles *SortedOraclesSession) GetOracles(token common.Address) ([]common.Address, error) {
	return _SortedOracles.Contract.GetOracles(&_SortedOracles.CallOpts, token)
}

// GetOracles is a free data retrieval call binding the contract method 0x8e749281.
//
// Solidity: function getOracles(address token) view returns(address[])
func (_SortedOracles *SortedOraclesCallerSession) GetOracles(token common.Address) ([]common.Address, error) {
	return _SortedOracles.Contract.GetOracles(&_SortedOracles.CallOpts, token)
}

// GetRates is a free data retrieval call binding the contract method 0x02f55b61.
//
// Solidity: function getRates(address token) view returns(address[], uint256[], uint8[])
func (_SortedOracles *SortedOraclesCaller) GetRates(opts *bind.CallOpts, token common.Address) ([]common.Address, []*big.Int, []uint8, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "getRates", token)

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), *new([]uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	out2 := *abi.ConvertType(out[2], new([]uint8)).(*[]uint8)

	return out0, out1, out2, err

}

// GetRates is a free data retrieval call binding the contract method 0x02f55b61.
//
// Solidity: function getRates(address token) view returns(address[], uint256[], uint8[])
func (_SortedOracles *SortedOraclesSession) GetRates(token common.Address) ([]common.Address, []*big.Int, []uint8, error) {
	return _SortedOracles.Contract.GetRates(&_SortedOracles.CallOpts, token)
}

// GetRates is a free data retrieval call binding the contract method 0x02f55b61.
//
// Solidity: function getRates(address token) view returns(address[], uint256[], uint8[])
func (_SortedOracles *SortedOraclesCallerSession) GetRates(token common.Address) ([]common.Address, []*big.Int, []uint8, error) {
	return _SortedOracles.Contract.GetRates(&_SortedOracles.CallOpts, token)
}

// GetTimestamps is a free data retrieval call binding the contract method 0xb9292158.
//
// Solidity: function getTimestamps(address token) view returns(address[], uint256[], uint8[])
func (_SortedOracles *SortedOraclesCaller) GetTimestamps(opts *bind.CallOpts, token common.Address) ([]common.Address, []*big.Int, []uint8, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "getTimestamps", token)

	if err != nil {
		return *new([]common.Address), *new([]*big.Int), *new([]uint8), err
	}

	out0 := *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	out1 := *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	out2 := *abi.ConvertType(out[2], new([]uint8)).(*[]uint8)

	return out0, out1, out2, err

}

// GetTimestamps is a free data retrieval call binding the contract method 0xb9292158.
//
// Solidity: function getTimestamps(address token) view returns(address[], uint256[], uint8[])
func (_SortedOracles *SortedOraclesSession) GetTimestamps(token common.Address) ([]common.Address, []*big.Int, []uint8, error) {
	return _SortedOracles.Contract.GetTimestamps(&_SortedOracles.CallOpts, token)
}

// GetTimestamps is a free data retrieval call binding the contract method 0xb9292158.
//
// Solidity: function getTimestamps(address token) view returns(address[], uint256[], uint8[])
func (_SortedOracles *SortedOraclesCallerSession) GetTimestamps(token common.Address) ([]common.Address, []*big.Int, []uint8, error) {
	return _SortedOracles.Contract.GetTimestamps(&_SortedOracles.CallOpts, token)
}

// GetTokenReportExpirySeconds is a free data retrieval call binding the contract method 0x6deb6799.
//
// Solidity: function getTokenReportExpirySeconds(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCaller) GetTokenReportExpirySeconds(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "getTokenReportExpirySeconds", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTokenReportExpirySeconds is a free data retrieval call binding the contract method 0x6deb6799.
//
// Solidity: function getTokenReportExpirySeconds(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesSession) GetTokenReportExpirySeconds(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.GetTokenReportExpirySeconds(&_SortedOracles.CallOpts, token)
}

// GetTokenReportExpirySeconds is a free data retrieval call binding the contract method 0x6deb6799.
//
// Solidity: function getTokenReportExpirySeconds(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCallerSession) GetTokenReportExpirySeconds(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.GetTokenReportExpirySeconds(&_SortedOracles.CallOpts, token)
}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_SortedOracles *SortedOraclesCaller) GetVersionNumber(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "getVersionNumber")

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
func (_SortedOracles *SortedOraclesSession) GetVersionNumber() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _SortedOracles.Contract.GetVersionNumber(&_SortedOracles.CallOpts)
}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_SortedOracles *SortedOraclesCallerSession) GetVersionNumber() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _SortedOracles.Contract.GetVersionNumber(&_SortedOracles.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_SortedOracles *SortedOraclesCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_SortedOracles *SortedOraclesSession) Initialized() (bool, error) {
	return _SortedOracles.Contract.Initialized(&_SortedOracles.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_SortedOracles *SortedOraclesCallerSession) Initialized() (bool, error) {
	return _SortedOracles.Contract.Initialized(&_SortedOracles.CallOpts)
}

// IsOldestReportExpired is a free data retrieval call binding the contract method 0xffe736bf.
//
// Solidity: function isOldestReportExpired(address token) view returns(bool, address)
func (_SortedOracles *SortedOraclesCaller) IsOldestReportExpired(opts *bind.CallOpts, token common.Address) (bool, common.Address, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "isOldestReportExpired", token)

	if err != nil {
		return *new(bool), *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return out0, out1, err

}

// IsOldestReportExpired is a free data retrieval call binding the contract method 0xffe736bf.
//
// Solidity: function isOldestReportExpired(address token) view returns(bool, address)
func (_SortedOracles *SortedOraclesSession) IsOldestReportExpired(token common.Address) (bool, common.Address, error) {
	return _SortedOracles.Contract.IsOldestReportExpired(&_SortedOracles.CallOpts, token)
}

// IsOldestReportExpired is a free data retrieval call binding the contract method 0xffe736bf.
//
// Solidity: function isOldestReportExpired(address token) view returns(bool, address)
func (_SortedOracles *SortedOraclesCallerSession) IsOldestReportExpired(token common.Address) (bool, common.Address, error) {
	return _SortedOracles.Contract.IsOldestReportExpired(&_SortedOracles.CallOpts, token)
}

// IsOracle is a free data retrieval call binding the contract method 0x370c998e.
//
// Solidity: function isOracle(address , address ) view returns(bool)
func (_SortedOracles *SortedOraclesCaller) IsOracle(opts *bind.CallOpts, arg0 common.Address, arg1 common.Address) (bool, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "isOracle", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOracle is a free data retrieval call binding the contract method 0x370c998e.
//
// Solidity: function isOracle(address , address ) view returns(bool)
func (_SortedOracles *SortedOraclesSession) IsOracle(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _SortedOracles.Contract.IsOracle(&_SortedOracles.CallOpts, arg0, arg1)
}

// IsOracle is a free data retrieval call binding the contract method 0x370c998e.
//
// Solidity: function isOracle(address , address ) view returns(bool)
func (_SortedOracles *SortedOraclesCallerSession) IsOracle(arg0 common.Address, arg1 common.Address) (bool, error) {
	return _SortedOracles.Contract.IsOracle(&_SortedOracles.CallOpts, arg0, arg1)
}

// MedianRate is a free data retrieval call binding the contract method 0xef90e1b0.
//
// Solidity: function medianRate(address token) view returns(uint256, uint256)
func (_SortedOracles *SortedOraclesCaller) MedianRate(opts *bind.CallOpts, token common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "medianRate", token)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// MedianRate is a free data retrieval call binding the contract method 0xef90e1b0.
//
// Solidity: function medianRate(address token) view returns(uint256, uint256)
func (_SortedOracles *SortedOraclesSession) MedianRate(token common.Address) (*big.Int, *big.Int, error) {
	return _SortedOracles.Contract.MedianRate(&_SortedOracles.CallOpts, token)
}

// MedianRate is a free data retrieval call binding the contract method 0xef90e1b0.
//
// Solidity: function medianRate(address token) view returns(uint256, uint256)
func (_SortedOracles *SortedOraclesCallerSession) MedianRate(token common.Address) (*big.Int, *big.Int, error) {
	return _SortedOracles.Contract.MedianRate(&_SortedOracles.CallOpts, token)
}

// MedianTimestamp is a free data retrieval call binding the contract method 0x071b48fc.
//
// Solidity: function medianTimestamp(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCaller) MedianTimestamp(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "medianTimestamp", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MedianTimestamp is a free data retrieval call binding the contract method 0x071b48fc.
//
// Solidity: function medianTimestamp(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesSession) MedianTimestamp(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.MedianTimestamp(&_SortedOracles.CallOpts, token)
}

// MedianTimestamp is a free data retrieval call binding the contract method 0x071b48fc.
//
// Solidity: function medianTimestamp(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCallerSession) MedianTimestamp(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.MedianTimestamp(&_SortedOracles.CallOpts, token)
}

// NumRates is a free data retrieval call binding the contract method 0xbbc66a94.
//
// Solidity: function numRates(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCaller) NumRates(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "numRates", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumRates is a free data retrieval call binding the contract method 0xbbc66a94.
//
// Solidity: function numRates(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesSession) NumRates(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.NumRates(&_SortedOracles.CallOpts, token)
}

// NumRates is a free data retrieval call binding the contract method 0xbbc66a94.
//
// Solidity: function numRates(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCallerSession) NumRates(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.NumRates(&_SortedOracles.CallOpts, token)
}

// NumTimestamps is a free data retrieval call binding the contract method 0x6dd6ef0c.
//
// Solidity: function numTimestamps(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCaller) NumTimestamps(opts *bind.CallOpts, token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "numTimestamps", token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// NumTimestamps is a free data retrieval call binding the contract method 0x6dd6ef0c.
//
// Solidity: function numTimestamps(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesSession) NumTimestamps(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.NumTimestamps(&_SortedOracles.CallOpts, token)
}

// NumTimestamps is a free data retrieval call binding the contract method 0x6dd6ef0c.
//
// Solidity: function numTimestamps(address token) view returns(uint256)
func (_SortedOracles *SortedOraclesCallerSession) NumTimestamps(token common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.NumTimestamps(&_SortedOracles.CallOpts, token)
}

// Oracles is a free data retrieval call binding the contract method 0xa00a8b2c.
//
// Solidity: function oracles(address , uint256 ) view returns(address)
func (_SortedOracles *SortedOraclesCaller) Oracles(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "oracles", arg0, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracles is a free data retrieval call binding the contract method 0xa00a8b2c.
//
// Solidity: function oracles(address , uint256 ) view returns(address)
func (_SortedOracles *SortedOraclesSession) Oracles(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _SortedOracles.Contract.Oracles(&_SortedOracles.CallOpts, arg0, arg1)
}

// Oracles is a free data retrieval call binding the contract method 0xa00a8b2c.
//
// Solidity: function oracles(address , uint256 ) view returns(address)
func (_SortedOracles *SortedOraclesCallerSession) Oracles(arg0 common.Address, arg1 *big.Int) (common.Address, error) {
	return _SortedOracles.Contract.Oracles(&_SortedOracles.CallOpts, arg0, arg1)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SortedOracles *SortedOraclesCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SortedOracles *SortedOraclesSession) Owner() (common.Address, error) {
	return _SortedOracles.Contract.Owner(&_SortedOracles.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_SortedOracles *SortedOraclesCallerSession) Owner() (common.Address, error) {
	return _SortedOracles.Contract.Owner(&_SortedOracles.CallOpts)
}

// ReportExpirySeconds is a free data retrieval call binding the contract method 0x493a353c.
//
// Solidity: function reportExpirySeconds() view returns(uint256)
func (_SortedOracles *SortedOraclesCaller) ReportExpirySeconds(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "reportExpirySeconds")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReportExpirySeconds is a free data retrieval call binding the contract method 0x493a353c.
//
// Solidity: function reportExpirySeconds() view returns(uint256)
func (_SortedOracles *SortedOraclesSession) ReportExpirySeconds() (*big.Int, error) {
	return _SortedOracles.Contract.ReportExpirySeconds(&_SortedOracles.CallOpts)
}

// ReportExpirySeconds is a free data retrieval call binding the contract method 0x493a353c.
//
// Solidity: function reportExpirySeconds() view returns(uint256)
func (_SortedOracles *SortedOraclesCallerSession) ReportExpirySeconds() (*big.Int, error) {
	return _SortedOracles.Contract.ReportExpirySeconds(&_SortedOracles.CallOpts)
}

// TokenReportExpirySeconds is a free data retrieval call binding the contract method 0x2e86bc01.
//
// Solidity: function tokenReportExpirySeconds(address ) view returns(uint256)
func (_SortedOracles *SortedOraclesCaller) TokenReportExpirySeconds(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _SortedOracles.contract.Call(opts, &out, "tokenReportExpirySeconds", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TokenReportExpirySeconds is a free data retrieval call binding the contract method 0x2e86bc01.
//
// Solidity: function tokenReportExpirySeconds(address ) view returns(uint256)
func (_SortedOracles *SortedOraclesSession) TokenReportExpirySeconds(arg0 common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.TokenReportExpirySeconds(&_SortedOracles.CallOpts, arg0)
}

// TokenReportExpirySeconds is a free data retrieval call binding the contract method 0x2e86bc01.
//
// Solidity: function tokenReportExpirySeconds(address ) view returns(uint256)
func (_SortedOracles *SortedOraclesCallerSession) TokenReportExpirySeconds(arg0 common.Address) (*big.Int, error) {
	return _SortedOracles.Contract.TokenReportExpirySeconds(&_SortedOracles.CallOpts, arg0)
}

// AddOracle is a paid mutator transaction binding the contract method 0xf0ca4adb.
//
// Solidity: function addOracle(address token, address oracleAddress) returns()
func (_SortedOracles *SortedOraclesTransactor) AddOracle(opts *bind.TransactOpts, token common.Address, oracleAddress common.Address) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "addOracle", token, oracleAddress)
}

// AddOracle is a paid mutator transaction binding the contract method 0xf0ca4adb.
//
// Solidity: function addOracle(address token, address oracleAddress) returns()
func (_SortedOracles *SortedOraclesSession) AddOracle(token common.Address, oracleAddress common.Address) (*types.Transaction, error) {
	return _SortedOracles.Contract.AddOracle(&_SortedOracles.TransactOpts, token, oracleAddress)
}

// AddOracle is a paid mutator transaction binding the contract method 0xf0ca4adb.
//
// Solidity: function addOracle(address token, address oracleAddress) returns()
func (_SortedOracles *SortedOraclesTransactorSession) AddOracle(token common.Address, oracleAddress common.Address) (*types.Transaction, error) {
	return _SortedOracles.Contract.AddOracle(&_SortedOracles.TransactOpts, token, oracleAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesTransactor) Initialize(opts *bind.TransactOpts, _reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "initialize", _reportExpirySeconds)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesSession) Initialize(_reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.Initialize(&_SortedOracles.TransactOpts, _reportExpirySeconds)
}

// Initialize is a paid mutator transaction binding the contract method 0xfe4b84df.
//
// Solidity: function initialize(uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesTransactorSession) Initialize(_reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.Initialize(&_SortedOracles.TransactOpts, _reportExpirySeconds)
}

// RemoveExpiredReports is a paid mutator transaction binding the contract method 0xdd34ca3b.
//
// Solidity: function removeExpiredReports(address token, uint256 n) returns()
func (_SortedOracles *SortedOraclesTransactor) RemoveExpiredReports(opts *bind.TransactOpts, token common.Address, n *big.Int) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "removeExpiredReports", token, n)
}

// RemoveExpiredReports is a paid mutator transaction binding the contract method 0xdd34ca3b.
//
// Solidity: function removeExpiredReports(address token, uint256 n) returns()
func (_SortedOracles *SortedOraclesSession) RemoveExpiredReports(token common.Address, n *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.RemoveExpiredReports(&_SortedOracles.TransactOpts, token, n)
}

// RemoveExpiredReports is a paid mutator transaction binding the contract method 0xdd34ca3b.
//
// Solidity: function removeExpiredReports(address token, uint256 n) returns()
func (_SortedOracles *SortedOraclesTransactorSession) RemoveExpiredReports(token common.Address, n *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.RemoveExpiredReports(&_SortedOracles.TransactOpts, token, n)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0x53a57297.
//
// Solidity: function removeOracle(address token, address oracleAddress, uint256 index) returns()
func (_SortedOracles *SortedOraclesTransactor) RemoveOracle(opts *bind.TransactOpts, token common.Address, oracleAddress common.Address, index *big.Int) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "removeOracle", token, oracleAddress, index)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0x53a57297.
//
// Solidity: function removeOracle(address token, address oracleAddress, uint256 index) returns()
func (_SortedOracles *SortedOraclesSession) RemoveOracle(token common.Address, oracleAddress common.Address, index *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.RemoveOracle(&_SortedOracles.TransactOpts, token, oracleAddress, index)
}

// RemoveOracle is a paid mutator transaction binding the contract method 0x53a57297.
//
// Solidity: function removeOracle(address token, address oracleAddress, uint256 index) returns()
func (_SortedOracles *SortedOraclesTransactorSession) RemoveOracle(token common.Address, oracleAddress common.Address, index *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.RemoveOracle(&_SortedOracles.TransactOpts, token, oracleAddress, index)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SortedOracles *SortedOraclesTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SortedOracles *SortedOraclesSession) RenounceOwnership() (*types.Transaction, error) {
	return _SortedOracles.Contract.RenounceOwnership(&_SortedOracles.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_SortedOracles *SortedOraclesTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _SortedOracles.Contract.RenounceOwnership(&_SortedOracles.TransactOpts)
}

// Report is a paid mutator transaction binding the contract method 0x80e50744.
//
// Solidity: function report(address token, uint256 value, address lesserKey, address greaterKey) returns()
func (_SortedOracles *SortedOraclesTransactor) Report(opts *bind.TransactOpts, token common.Address, value *big.Int, lesserKey common.Address, greaterKey common.Address) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "report", token, value, lesserKey, greaterKey)
}

// Report is a paid mutator transaction binding the contract method 0x80e50744.
//
// Solidity: function report(address token, uint256 value, address lesserKey, address greaterKey) returns()
func (_SortedOracles *SortedOraclesSession) Report(token common.Address, value *big.Int, lesserKey common.Address, greaterKey common.Address) (*types.Transaction, error) {
	return _SortedOracles.Contract.Report(&_SortedOracles.TransactOpts, token, value, lesserKey, greaterKey)
}

// Report is a paid mutator transaction binding the contract method 0x80e50744.
//
// Solidity: function report(address token, uint256 value, address lesserKey, address greaterKey) returns()
func (_SortedOracles *SortedOraclesTransactorSession) Report(token common.Address, value *big.Int, lesserKey common.Address, greaterKey common.Address) (*types.Transaction, error) {
	return _SortedOracles.Contract.Report(&_SortedOracles.TransactOpts, token, value, lesserKey, greaterKey)
}

// SetReportExpiry is a paid mutator transaction binding the contract method 0xebc1d6bb.
//
// Solidity: function setReportExpiry(uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesTransactor) SetReportExpiry(opts *bind.TransactOpts, _reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "setReportExpiry", _reportExpirySeconds)
}

// SetReportExpiry is a paid mutator transaction binding the contract method 0xebc1d6bb.
//
// Solidity: function setReportExpiry(uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesSession) SetReportExpiry(_reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.SetReportExpiry(&_SortedOracles.TransactOpts, _reportExpirySeconds)
}

// SetReportExpiry is a paid mutator transaction binding the contract method 0xebc1d6bb.
//
// Solidity: function setReportExpiry(uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesTransactorSession) SetReportExpiry(_reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.SetReportExpiry(&_SortedOracles.TransactOpts, _reportExpirySeconds)
}

// SetTokenReportExpiry is a paid mutator transaction binding the contract method 0xfc20935d.
//
// Solidity: function setTokenReportExpiry(address _token, uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesTransactor) SetTokenReportExpiry(opts *bind.TransactOpts, _token common.Address, _reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "setTokenReportExpiry", _token, _reportExpirySeconds)
}

// SetTokenReportExpiry is a paid mutator transaction binding the contract method 0xfc20935d.
//
// Solidity: function setTokenReportExpiry(address _token, uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesSession) SetTokenReportExpiry(_token common.Address, _reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.SetTokenReportExpiry(&_SortedOracles.TransactOpts, _token, _reportExpirySeconds)
}

// SetTokenReportExpiry is a paid mutator transaction binding the contract method 0xfc20935d.
//
// Solidity: function setTokenReportExpiry(address _token, uint256 _reportExpirySeconds) returns()
func (_SortedOracles *SortedOraclesTransactorSession) SetTokenReportExpiry(_token common.Address, _reportExpirySeconds *big.Int) (*types.Transaction, error) {
	return _SortedOracles.Contract.SetTokenReportExpiry(&_SortedOracles.TransactOpts, _token, _reportExpirySeconds)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SortedOracles *SortedOraclesTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _SortedOracles.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SortedOracles *SortedOraclesSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SortedOracles.Contract.TransferOwnership(&_SortedOracles.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_SortedOracles *SortedOraclesTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _SortedOracles.Contract.TransferOwnership(&_SortedOracles.TransactOpts, newOwner)
}

// SortedOraclesMedianUpdatedIterator is returned from FilterMedianUpdated and is used to iterate over the raw logs and unpacked data for MedianUpdated events raised by the SortedOracles contract.
type SortedOraclesMedianUpdatedIterator struct {
	Event *SortedOraclesMedianUpdated // Event containing the contract specifics and raw log

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
func (it *SortedOraclesMedianUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesMedianUpdated)
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
		it.Event = new(SortedOraclesMedianUpdated)
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
func (it *SortedOraclesMedianUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesMedianUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesMedianUpdated represents a MedianUpdated event raised by the SortedOracles contract.
type SortedOraclesMedianUpdated struct {
	Token common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterMedianUpdated is a free log retrieval operation binding the contract event 0xa9981ebfc3b766a742486e898f54959b050a66006dbce1a4155c1f84a08bcf41.
//
// Solidity: event MedianUpdated(address indexed token, uint256 value)
func (_SortedOracles *SortedOraclesFilterer) FilterMedianUpdated(opts *bind.FilterOpts, token []common.Address) (*SortedOraclesMedianUpdatedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "MedianUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesMedianUpdatedIterator{contract: _SortedOracles.contract, event: "MedianUpdated", logs: logs, sub: sub}, nil
}

// WatchMedianUpdated is a free log subscription operation binding the contract event 0xa9981ebfc3b766a742486e898f54959b050a66006dbce1a4155c1f84a08bcf41.
//
// Solidity: event MedianUpdated(address indexed token, uint256 value)
func (_SortedOracles *SortedOraclesFilterer) WatchMedianUpdated(opts *bind.WatchOpts, sink chan<- *SortedOraclesMedianUpdated, token []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "MedianUpdated", tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesMedianUpdated)
				if err := _SortedOracles.contract.UnpackLog(event, "MedianUpdated", log); err != nil {
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

// ParseMedianUpdated is a log parse operation binding the contract event 0xa9981ebfc3b766a742486e898f54959b050a66006dbce1a4155c1f84a08bcf41.
//
// Solidity: event MedianUpdated(address indexed token, uint256 value)
func (_SortedOracles *SortedOraclesFilterer) ParseMedianUpdated(log types.Log) (*SortedOraclesMedianUpdated, error) {
	event := new(SortedOraclesMedianUpdated)
	if err := _SortedOracles.contract.UnpackLog(event, "MedianUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesOracleAddedIterator is returned from FilterOracleAdded and is used to iterate over the raw logs and unpacked data for OracleAdded events raised by the SortedOracles contract.
type SortedOraclesOracleAddedIterator struct {
	Event *SortedOraclesOracleAdded // Event containing the contract specifics and raw log

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
func (it *SortedOraclesOracleAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesOracleAdded)
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
		it.Event = new(SortedOraclesOracleAdded)
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
func (it *SortedOraclesOracleAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesOracleAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesOracleAdded represents a OracleAdded event raised by the SortedOracles contract.
type SortedOraclesOracleAdded struct {
	Token         common.Address
	OracleAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOracleAdded is a free log retrieval operation binding the contract event 0x828d2be040dede7698182e08dfa8bfbd663c879aee772509c4a2bd961d0ed43f.
//
// Solidity: event OracleAdded(address indexed token, address indexed oracleAddress)
func (_SortedOracles *SortedOraclesFilterer) FilterOracleAdded(opts *bind.FilterOpts, token []common.Address, oracleAddress []common.Address) (*SortedOraclesOracleAddedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleAddressRule []interface{}
	for _, oracleAddressItem := range oracleAddress {
		oracleAddressRule = append(oracleAddressRule, oracleAddressItem)
	}

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "OracleAdded", tokenRule, oracleAddressRule)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesOracleAddedIterator{contract: _SortedOracles.contract, event: "OracleAdded", logs: logs, sub: sub}, nil
}

// WatchOracleAdded is a free log subscription operation binding the contract event 0x828d2be040dede7698182e08dfa8bfbd663c879aee772509c4a2bd961d0ed43f.
//
// Solidity: event OracleAdded(address indexed token, address indexed oracleAddress)
func (_SortedOracles *SortedOraclesFilterer) WatchOracleAdded(opts *bind.WatchOpts, sink chan<- *SortedOraclesOracleAdded, token []common.Address, oracleAddress []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleAddressRule []interface{}
	for _, oracleAddressItem := range oracleAddress {
		oracleAddressRule = append(oracleAddressRule, oracleAddressItem)
	}

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "OracleAdded", tokenRule, oracleAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesOracleAdded)
				if err := _SortedOracles.contract.UnpackLog(event, "OracleAdded", log); err != nil {
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

// ParseOracleAdded is a log parse operation binding the contract event 0x828d2be040dede7698182e08dfa8bfbd663c879aee772509c4a2bd961d0ed43f.
//
// Solidity: event OracleAdded(address indexed token, address indexed oracleAddress)
func (_SortedOracles *SortedOraclesFilterer) ParseOracleAdded(log types.Log) (*SortedOraclesOracleAdded, error) {
	event := new(SortedOraclesOracleAdded)
	if err := _SortedOracles.contract.UnpackLog(event, "OracleAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesOracleRemovedIterator is returned from FilterOracleRemoved and is used to iterate over the raw logs and unpacked data for OracleRemoved events raised by the SortedOracles contract.
type SortedOraclesOracleRemovedIterator struct {
	Event *SortedOraclesOracleRemoved // Event containing the contract specifics and raw log

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
func (it *SortedOraclesOracleRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesOracleRemoved)
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
		it.Event = new(SortedOraclesOracleRemoved)
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
func (it *SortedOraclesOracleRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesOracleRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesOracleRemoved represents a OracleRemoved event raised by the SortedOracles contract.
type SortedOraclesOracleRemoved struct {
	Token         common.Address
	OracleAddress common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOracleRemoved is a free log retrieval operation binding the contract event 0x6dc84b66cc948d847632b9d829f7cb1cb904fbf2c084554a9bc22ad9d8453340.
//
// Solidity: event OracleRemoved(address indexed token, address indexed oracleAddress)
func (_SortedOracles *SortedOraclesFilterer) FilterOracleRemoved(opts *bind.FilterOpts, token []common.Address, oracleAddress []common.Address) (*SortedOraclesOracleRemovedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleAddressRule []interface{}
	for _, oracleAddressItem := range oracleAddress {
		oracleAddressRule = append(oracleAddressRule, oracleAddressItem)
	}

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "OracleRemoved", tokenRule, oracleAddressRule)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesOracleRemovedIterator{contract: _SortedOracles.contract, event: "OracleRemoved", logs: logs, sub: sub}, nil
}

// WatchOracleRemoved is a free log subscription operation binding the contract event 0x6dc84b66cc948d847632b9d829f7cb1cb904fbf2c084554a9bc22ad9d8453340.
//
// Solidity: event OracleRemoved(address indexed token, address indexed oracleAddress)
func (_SortedOracles *SortedOraclesFilterer) WatchOracleRemoved(opts *bind.WatchOpts, sink chan<- *SortedOraclesOracleRemoved, token []common.Address, oracleAddress []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleAddressRule []interface{}
	for _, oracleAddressItem := range oracleAddress {
		oracleAddressRule = append(oracleAddressRule, oracleAddressItem)
	}

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "OracleRemoved", tokenRule, oracleAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesOracleRemoved)
				if err := _SortedOracles.contract.UnpackLog(event, "OracleRemoved", log); err != nil {
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

// ParseOracleRemoved is a log parse operation binding the contract event 0x6dc84b66cc948d847632b9d829f7cb1cb904fbf2c084554a9bc22ad9d8453340.
//
// Solidity: event OracleRemoved(address indexed token, address indexed oracleAddress)
func (_SortedOracles *SortedOraclesFilterer) ParseOracleRemoved(log types.Log) (*SortedOraclesOracleRemoved, error) {
	event := new(SortedOraclesOracleRemoved)
	if err := _SortedOracles.contract.UnpackLog(event, "OracleRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesOracleReportRemovedIterator is returned from FilterOracleReportRemoved and is used to iterate over the raw logs and unpacked data for OracleReportRemoved events raised by the SortedOracles contract.
type SortedOraclesOracleReportRemovedIterator struct {
	Event *SortedOraclesOracleReportRemoved // Event containing the contract specifics and raw log

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
func (it *SortedOraclesOracleReportRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesOracleReportRemoved)
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
		it.Event = new(SortedOraclesOracleReportRemoved)
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
func (it *SortedOraclesOracleReportRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesOracleReportRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesOracleReportRemoved represents a OracleReportRemoved event raised by the SortedOracles contract.
type SortedOraclesOracleReportRemoved struct {
	Token  common.Address
	Oracle common.Address
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterOracleReportRemoved is a free log retrieval operation binding the contract event 0xe21a44017b6fa1658d84e937d56ff408501facdb4ff7427c479ac460d76f7893.
//
// Solidity: event OracleReportRemoved(address indexed token, address indexed oracle)
func (_SortedOracles *SortedOraclesFilterer) FilterOracleReportRemoved(opts *bind.FilterOpts, token []common.Address, oracle []common.Address) (*SortedOraclesOracleReportRemovedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "OracleReportRemoved", tokenRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesOracleReportRemovedIterator{contract: _SortedOracles.contract, event: "OracleReportRemoved", logs: logs, sub: sub}, nil
}

// WatchOracleReportRemoved is a free log subscription operation binding the contract event 0xe21a44017b6fa1658d84e937d56ff408501facdb4ff7427c479ac460d76f7893.
//
// Solidity: event OracleReportRemoved(address indexed token, address indexed oracle)
func (_SortedOracles *SortedOraclesFilterer) WatchOracleReportRemoved(opts *bind.WatchOpts, sink chan<- *SortedOraclesOracleReportRemoved, token []common.Address, oracle []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "OracleReportRemoved", tokenRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesOracleReportRemoved)
				if err := _SortedOracles.contract.UnpackLog(event, "OracleReportRemoved", log); err != nil {
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

// ParseOracleReportRemoved is a log parse operation binding the contract event 0xe21a44017b6fa1658d84e937d56ff408501facdb4ff7427c479ac460d76f7893.
//
// Solidity: event OracleReportRemoved(address indexed token, address indexed oracle)
func (_SortedOracles *SortedOraclesFilterer) ParseOracleReportRemoved(log types.Log) (*SortedOraclesOracleReportRemoved, error) {
	event := new(SortedOraclesOracleReportRemoved)
	if err := _SortedOracles.contract.UnpackLog(event, "OracleReportRemoved", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesOracleReportedIterator is returned from FilterOracleReported and is used to iterate over the raw logs and unpacked data for OracleReported events raised by the SortedOracles contract.
type SortedOraclesOracleReportedIterator struct {
	Event *SortedOraclesOracleReported // Event containing the contract specifics and raw log

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
func (it *SortedOraclesOracleReportedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesOracleReported)
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
		it.Event = new(SortedOraclesOracleReported)
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
func (it *SortedOraclesOracleReportedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesOracleReportedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesOracleReported represents a OracleReported event raised by the SortedOracles contract.
type SortedOraclesOracleReported struct {
	Token     common.Address
	Oracle    common.Address
	Timestamp *big.Int
	Value     *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleReported is a free log retrieval operation binding the contract event 0x7cebb17173a9ed273d2b7538f64395c0ebf352ff743f1cf8ce66b437a6144213.
//
// Solidity: event OracleReported(address indexed token, address indexed oracle, uint256 timestamp, uint256 value)
func (_SortedOracles *SortedOraclesFilterer) FilterOracleReported(opts *bind.FilterOpts, token []common.Address, oracle []common.Address) (*SortedOraclesOracleReportedIterator, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "OracleReported", tokenRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesOracleReportedIterator{contract: _SortedOracles.contract, event: "OracleReported", logs: logs, sub: sub}, nil
}

// WatchOracleReported is a free log subscription operation binding the contract event 0x7cebb17173a9ed273d2b7538f64395c0ebf352ff743f1cf8ce66b437a6144213.
//
// Solidity: event OracleReported(address indexed token, address indexed oracle, uint256 timestamp, uint256 value)
func (_SortedOracles *SortedOraclesFilterer) WatchOracleReported(opts *bind.WatchOpts, sink chan<- *SortedOraclesOracleReported, token []common.Address, oracle []common.Address) (event.Subscription, error) {

	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}
	var oracleRule []interface{}
	for _, oracleItem := range oracle {
		oracleRule = append(oracleRule, oracleItem)
	}

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "OracleReported", tokenRule, oracleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesOracleReported)
				if err := _SortedOracles.contract.UnpackLog(event, "OracleReported", log); err != nil {
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

// ParseOracleReported is a log parse operation binding the contract event 0x7cebb17173a9ed273d2b7538f64395c0ebf352ff743f1cf8ce66b437a6144213.
//
// Solidity: event OracleReported(address indexed token, address indexed oracle, uint256 timestamp, uint256 value)
func (_SortedOracles *SortedOraclesFilterer) ParseOracleReported(log types.Log) (*SortedOraclesOracleReported, error) {
	event := new(SortedOraclesOracleReported)
	if err := _SortedOracles.contract.UnpackLog(event, "OracleReported", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the SortedOracles contract.
type SortedOraclesOwnershipTransferredIterator struct {
	Event *SortedOraclesOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *SortedOraclesOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesOwnershipTransferred)
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
		it.Event = new(SortedOraclesOwnershipTransferred)
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
func (it *SortedOraclesOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesOwnershipTransferred represents a OwnershipTransferred event raised by the SortedOracles contract.
type SortedOraclesOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SortedOracles *SortedOraclesFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*SortedOraclesOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &SortedOraclesOwnershipTransferredIterator{contract: _SortedOracles.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_SortedOracles *SortedOraclesFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *SortedOraclesOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesOwnershipTransferred)
				if err := _SortedOracles.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_SortedOracles *SortedOraclesFilterer) ParseOwnershipTransferred(log types.Log) (*SortedOraclesOwnershipTransferred, error) {
	event := new(SortedOraclesOwnershipTransferred)
	if err := _SortedOracles.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesReportExpirySetIterator is returned from FilterReportExpirySet and is used to iterate over the raw logs and unpacked data for ReportExpirySet events raised by the SortedOracles contract.
type SortedOraclesReportExpirySetIterator struct {
	Event *SortedOraclesReportExpirySet // Event containing the contract specifics and raw log

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
func (it *SortedOraclesReportExpirySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesReportExpirySet)
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
		it.Event = new(SortedOraclesReportExpirySet)
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
func (it *SortedOraclesReportExpirySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesReportExpirySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesReportExpirySet represents a ReportExpirySet event raised by the SortedOracles contract.
type SortedOraclesReportExpirySet struct {
	ReportExpiry *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterReportExpirySet is a free log retrieval operation binding the contract event 0xc68a9b88effd8a11611ff410efbc83569f0031b7bc70dd455b61344c7f0a042f.
//
// Solidity: event ReportExpirySet(uint256 reportExpiry)
func (_SortedOracles *SortedOraclesFilterer) FilterReportExpirySet(opts *bind.FilterOpts) (*SortedOraclesReportExpirySetIterator, error) {

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "ReportExpirySet")
	if err != nil {
		return nil, err
	}
	return &SortedOraclesReportExpirySetIterator{contract: _SortedOracles.contract, event: "ReportExpirySet", logs: logs, sub: sub}, nil
}

// WatchReportExpirySet is a free log subscription operation binding the contract event 0xc68a9b88effd8a11611ff410efbc83569f0031b7bc70dd455b61344c7f0a042f.
//
// Solidity: event ReportExpirySet(uint256 reportExpiry)
func (_SortedOracles *SortedOraclesFilterer) WatchReportExpirySet(opts *bind.WatchOpts, sink chan<- *SortedOraclesReportExpirySet) (event.Subscription, error) {

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "ReportExpirySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesReportExpirySet)
				if err := _SortedOracles.contract.UnpackLog(event, "ReportExpirySet", log); err != nil {
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

// ParseReportExpirySet is a log parse operation binding the contract event 0xc68a9b88effd8a11611ff410efbc83569f0031b7bc70dd455b61344c7f0a042f.
//
// Solidity: event ReportExpirySet(uint256 reportExpiry)
func (_SortedOracles *SortedOraclesFilterer) ParseReportExpirySet(log types.Log) (*SortedOraclesReportExpirySet, error) {
	event := new(SortedOraclesReportExpirySet)
	if err := _SortedOracles.contract.UnpackLog(event, "ReportExpirySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// SortedOraclesTokenReportExpirySetIterator is returned from FilterTokenReportExpirySet and is used to iterate over the raw logs and unpacked data for TokenReportExpirySet events raised by the SortedOracles contract.
type SortedOraclesTokenReportExpirySetIterator struct {
	Event *SortedOraclesTokenReportExpirySet // Event containing the contract specifics and raw log

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
func (it *SortedOraclesTokenReportExpirySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(SortedOraclesTokenReportExpirySet)
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
		it.Event = new(SortedOraclesTokenReportExpirySet)
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
func (it *SortedOraclesTokenReportExpirySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *SortedOraclesTokenReportExpirySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// SortedOraclesTokenReportExpirySet represents a TokenReportExpirySet event raised by the SortedOracles contract.
type SortedOraclesTokenReportExpirySet struct {
	Token        common.Address
	ReportExpiry *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenReportExpirySet is a free log retrieval operation binding the contract event 0xf8324c8592dfd9991ee3e717351afe0a964605257959e3d99b0eb3d45bff9422.
//
// Solidity: event TokenReportExpirySet(address token, uint256 reportExpiry)
func (_SortedOracles *SortedOraclesFilterer) FilterTokenReportExpirySet(opts *bind.FilterOpts) (*SortedOraclesTokenReportExpirySetIterator, error) {

	logs, sub, err := _SortedOracles.contract.FilterLogs(opts, "TokenReportExpirySet")
	if err != nil {
		return nil, err
	}
	return &SortedOraclesTokenReportExpirySetIterator{contract: _SortedOracles.contract, event: "TokenReportExpirySet", logs: logs, sub: sub}, nil
}

// WatchTokenReportExpirySet is a free log subscription operation binding the contract event 0xf8324c8592dfd9991ee3e717351afe0a964605257959e3d99b0eb3d45bff9422.
//
// Solidity: event TokenReportExpirySet(address token, uint256 reportExpiry)
func (_SortedOracles *SortedOraclesFilterer) WatchTokenReportExpirySet(opts *bind.WatchOpts, sink chan<- *SortedOraclesTokenReportExpirySet) (event.Subscription, error) {

	logs, sub, err := _SortedOracles.contract.WatchLogs(opts, "TokenReportExpirySet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(SortedOraclesTokenReportExpirySet)
				if err := _SortedOracles.contract.UnpackLog(event, "TokenReportExpirySet", log); err != nil {
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

// ParseTokenReportExpirySet is a log parse operation binding the contract event 0xf8324c8592dfd9991ee3e717351afe0a964605257959e3d99b0eb3d45bff9422.
//
// Solidity: event TokenReportExpirySet(address token, uint256 reportExpiry)
func (_SortedOracles *SortedOraclesFilterer) ParseTokenReportExpirySet(log types.Log) (*SortedOraclesTokenReportExpirySet, error) {
	event := new(SortedOraclesTokenReportExpirySet)
	if err := _SortedOracles.contract.UnpackLog(event, "TokenReportExpirySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
