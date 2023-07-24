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

// GoldTokenMetaData contains all meta data concerning the GoldToken contract.
var GoldTokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bool\",\"name\":\"test\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"registryAddress\",\"type\":\"address\"}],\"name\":\"RegistrySet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"string\",\"name\":\"comment\",\"type\":\"string\"}],\"name\":\"TransferComment\",\"type\":\"event\"},{\"constant\":true,\"inputs\":[],\"name\":\"initialized\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractIRegistry\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"registryAddress\",\"type\":\"address\"}],\"name\":\"setRegistry\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getVersionNumber\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"registryAddress\",\"type\":\"address\"}],\"name\":\"initialize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"string\",\"name\":\"comment\",\"type\":\"string\"}],\"name\":\"transferWithComment\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"circulatingSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"increaseSupply\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getBurnedAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// GoldTokenABI is the input ABI used to generate the binding from.
// Deprecated: Use GoldTokenMetaData.ABI instead.
var GoldTokenABI = GoldTokenMetaData.ABI

// GoldToken is an auto generated Go binding around an Ethereum contract.
type GoldToken struct {
	GoldTokenCaller     // Read-only binding to the contract
	GoldTokenTransactor // Write-only binding to the contract
	GoldTokenFilterer   // Log filterer for contract events
}

// GoldTokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type GoldTokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GoldTokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GoldTokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GoldTokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GoldTokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GoldTokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GoldTokenSession struct {
	Contract     *GoldToken        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GoldTokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GoldTokenCallerSession struct {
	Contract *GoldTokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// GoldTokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GoldTokenTransactorSession struct {
	Contract     *GoldTokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// GoldTokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type GoldTokenRaw struct {
	Contract *GoldToken // Generic contract binding to access the raw methods on
}

// GoldTokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GoldTokenCallerRaw struct {
	Contract *GoldTokenCaller // Generic read-only contract binding to access the raw methods on
}

// GoldTokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GoldTokenTransactorRaw struct {
	Contract *GoldTokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGoldToken creates a new instance of GoldToken, bound to a specific deployed contract.
func NewGoldToken(address common.Address, backend bind.ContractBackend) (*GoldToken, error) {
	contract, err := bindGoldToken(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &GoldToken{GoldTokenCaller: GoldTokenCaller{contract: contract}, GoldTokenTransactor: GoldTokenTransactor{contract: contract}, GoldTokenFilterer: GoldTokenFilterer{contract: contract}}, nil
}

// NewGoldTokenCaller creates a new read-only instance of GoldToken, bound to a specific deployed contract.
func NewGoldTokenCaller(address common.Address, caller bind.ContractCaller) (*GoldTokenCaller, error) {
	contract, err := bindGoldToken(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GoldTokenCaller{contract: contract}, nil
}

// NewGoldTokenTransactor creates a new write-only instance of GoldToken, bound to a specific deployed contract.
func NewGoldTokenTransactor(address common.Address, transactor bind.ContractTransactor) (*GoldTokenTransactor, error) {
	contract, err := bindGoldToken(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GoldTokenTransactor{contract: contract}, nil
}

// NewGoldTokenFilterer creates a new log filterer instance of GoldToken, bound to a specific deployed contract.
func NewGoldTokenFilterer(address common.Address, filterer bind.ContractFilterer) (*GoldTokenFilterer, error) {
	contract, err := bindGoldToken(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GoldTokenFilterer{contract: contract}, nil
}

// bindGoldToken binds a generic wrapper to an already deployed contract.
func bindGoldToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GoldTokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GoldToken *GoldTokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GoldToken.Contract.GoldTokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GoldToken *GoldTokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GoldToken.Contract.GoldTokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GoldToken *GoldTokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GoldToken.Contract.GoldTokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_GoldToken *GoldTokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _GoldToken.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_GoldToken *GoldTokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GoldToken.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_GoldToken *GoldTokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _GoldToken.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_GoldToken *GoldTokenCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_GoldToken *GoldTokenSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _GoldToken.Contract.Allowance(&_GoldToken.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_GoldToken *GoldTokenCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _GoldToken.Contract.Allowance(&_GoldToken.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GoldToken *GoldTokenCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GoldToken *GoldTokenSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GoldToken.Contract.BalanceOf(&_GoldToken.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_GoldToken *GoldTokenCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _GoldToken.Contract.BalanceOf(&_GoldToken.CallOpts, owner)
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_GoldToken *GoldTokenCaller) CirculatingSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "circulatingSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_GoldToken *GoldTokenSession) CirculatingSupply() (*big.Int, error) {
	return _GoldToken.Contract.CirculatingSupply(&_GoldToken.CallOpts)
}

// CirculatingSupply is a free data retrieval call binding the contract method 0x9358928b.
//
// Solidity: function circulatingSupply() view returns(uint256)
func (_GoldToken *GoldTokenCallerSession) CirculatingSupply() (*big.Int, error) {
	return _GoldToken.Contract.CirculatingSupply(&_GoldToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_GoldToken *GoldTokenCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_GoldToken *GoldTokenSession) Decimals() (uint8, error) {
	return _GoldToken.Contract.Decimals(&_GoldToken.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_GoldToken *GoldTokenCallerSession) Decimals() (uint8, error) {
	return _GoldToken.Contract.Decimals(&_GoldToken.CallOpts)
}

// GetBurnedAmount is a free data retrieval call binding the contract method 0x265126bd.
//
// Solidity: function getBurnedAmount() view returns(uint256)
func (_GoldToken *GoldTokenCaller) GetBurnedAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "getBurnedAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBurnedAmount is a free data retrieval call binding the contract method 0x265126bd.
//
// Solidity: function getBurnedAmount() view returns(uint256)
func (_GoldToken *GoldTokenSession) GetBurnedAmount() (*big.Int, error) {
	return _GoldToken.Contract.GetBurnedAmount(&_GoldToken.CallOpts)
}

// GetBurnedAmount is a free data retrieval call binding the contract method 0x265126bd.
//
// Solidity: function getBurnedAmount() view returns(uint256)
func (_GoldToken *GoldTokenCallerSession) GetBurnedAmount() (*big.Int, error) {
	return _GoldToken.Contract.GetBurnedAmount(&_GoldToken.CallOpts)
}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_GoldToken *GoldTokenCaller) GetVersionNumber(opts *bind.CallOpts) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "getVersionNumber")

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
func (_GoldToken *GoldTokenSession) GetVersionNumber() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _GoldToken.Contract.GetVersionNumber(&_GoldToken.CallOpts)
}

// GetVersionNumber is a free data retrieval call binding the contract method 0x54255be0.
//
// Solidity: function getVersionNumber() pure returns(uint256, uint256, uint256, uint256)
func (_GoldToken *GoldTokenCallerSession) GetVersionNumber() (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	return _GoldToken.Contract.GetVersionNumber(&_GoldToken.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_GoldToken *GoldTokenCaller) Initialized(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "initialized")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_GoldToken *GoldTokenSession) Initialized() (bool, error) {
	return _GoldToken.Contract.Initialized(&_GoldToken.CallOpts)
}

// Initialized is a free data retrieval call binding the contract method 0x158ef93e.
//
// Solidity: function initialized() view returns(bool)
func (_GoldToken *GoldTokenCallerSession) Initialized() (bool, error) {
	return _GoldToken.Contract.Initialized(&_GoldToken.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_GoldToken *GoldTokenCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "isOwner")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_GoldToken *GoldTokenSession) IsOwner() (bool, error) {
	return _GoldToken.Contract.IsOwner(&_GoldToken.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() view returns(bool)
func (_GoldToken *GoldTokenCallerSession) IsOwner() (bool, error) {
	return _GoldToken.Contract.IsOwner(&_GoldToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GoldToken *GoldTokenCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GoldToken *GoldTokenSession) Name() (string, error) {
	return _GoldToken.Contract.Name(&_GoldToken.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_GoldToken *GoldTokenCallerSession) Name() (string, error) {
	return _GoldToken.Contract.Name(&_GoldToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GoldToken *GoldTokenCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GoldToken *GoldTokenSession) Owner() (common.Address, error) {
	return _GoldToken.Contract.Owner(&_GoldToken.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_GoldToken *GoldTokenCallerSession) Owner() (common.Address, error) {
	return _GoldToken.Contract.Owner(&_GoldToken.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_GoldToken *GoldTokenCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "registry")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_GoldToken *GoldTokenSession) Registry() (common.Address, error) {
	return _GoldToken.Contract.Registry(&_GoldToken.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() view returns(address)
func (_GoldToken *GoldTokenCallerSession) Registry() (common.Address, error) {
	return _GoldToken.Contract.Registry(&_GoldToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GoldToken *GoldTokenCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GoldToken *GoldTokenSession) Symbol() (string, error) {
	return _GoldToken.Contract.Symbol(&_GoldToken.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_GoldToken *GoldTokenCallerSession) Symbol() (string, error) {
	return _GoldToken.Contract.Symbol(&_GoldToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GoldToken *GoldTokenCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _GoldToken.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GoldToken *GoldTokenSession) TotalSupply() (*big.Int, error) {
	return _GoldToken.Contract.TotalSupply(&_GoldToken.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_GoldToken *GoldTokenCallerSession) TotalSupply() (*big.Int, error) {
	return _GoldToken.Contract.TotalSupply(&_GoldToken.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) Approve(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "approve", spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Approve(&_GoldToken.TransactOpts, spender, value)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) Approve(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Approve(&_GoldToken.TransactOpts, spender, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) Burn(opts *bind.TransactOpts, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "burn", value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Burn(&_GoldToken.TransactOpts, value)
}

// Burn is a paid mutator transaction binding the contract method 0x42966c68.
//
// Solidity: function burn(uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) Burn(value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Burn(&_GoldToken.TransactOpts, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "decreaseAllowance", spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) DecreaseAllowance(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.DecreaseAllowance(&_GoldToken.TransactOpts, spender, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) DecreaseAllowance(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.DecreaseAllowance(&_GoldToken.TransactOpts, spender, value)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "increaseAllowance", spender, value)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) IncreaseAllowance(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.IncreaseAllowance(&_GoldToken.TransactOpts, spender, value)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) IncreaseAllowance(spender common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.IncreaseAllowance(&_GoldToken.TransactOpts, spender, value)
}

// IncreaseSupply is a paid mutator transaction binding the contract method 0xb921e163.
//
// Solidity: function increaseSupply(uint256 amount) returns()
func (_GoldToken *GoldTokenTransactor) IncreaseSupply(opts *bind.TransactOpts, amount *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "increaseSupply", amount)
}

// IncreaseSupply is a paid mutator transaction binding the contract method 0xb921e163.
//
// Solidity: function increaseSupply(uint256 amount) returns()
func (_GoldToken *GoldTokenSession) IncreaseSupply(amount *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.IncreaseSupply(&_GoldToken.TransactOpts, amount)
}

// IncreaseSupply is a paid mutator transaction binding the contract method 0xb921e163.
//
// Solidity: function increaseSupply(uint256 amount) returns()
func (_GoldToken *GoldTokenTransactorSession) IncreaseSupply(amount *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.IncreaseSupply(&_GoldToken.TransactOpts, amount)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address registryAddress) returns()
func (_GoldToken *GoldTokenTransactor) Initialize(opts *bind.TransactOpts, registryAddress common.Address) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "initialize", registryAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address registryAddress) returns()
func (_GoldToken *GoldTokenSession) Initialize(registryAddress common.Address) (*types.Transaction, error) {
	return _GoldToken.Contract.Initialize(&_GoldToken.TransactOpts, registryAddress)
}

// Initialize is a paid mutator transaction binding the contract method 0xc4d66de8.
//
// Solidity: function initialize(address registryAddress) returns()
func (_GoldToken *GoldTokenTransactorSession) Initialize(registryAddress common.Address) (*types.Transaction, error) {
	return _GoldToken.Contract.Initialize(&_GoldToken.TransactOpts, registryAddress)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) Mint(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "mint", to, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) Mint(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Mint(&_GoldToken.TransactOpts, to, value)
}

// Mint is a paid mutator transaction binding the contract method 0x40c10f19.
//
// Solidity: function mint(address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) Mint(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Mint(&_GoldToken.TransactOpts, to, value)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GoldToken *GoldTokenTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GoldToken *GoldTokenSession) RenounceOwnership() (*types.Transaction, error) {
	return _GoldToken.Contract.RenounceOwnership(&_GoldToken.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_GoldToken *GoldTokenTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _GoldToken.Contract.RenounceOwnership(&_GoldToken.TransactOpts)
}

// SetRegistry is a paid mutator transaction binding the contract method 0xa91ee0dc.
//
// Solidity: function setRegistry(address registryAddress) returns()
func (_GoldToken *GoldTokenTransactor) SetRegistry(opts *bind.TransactOpts, registryAddress common.Address) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "setRegistry", registryAddress)
}

// SetRegistry is a paid mutator transaction binding the contract method 0xa91ee0dc.
//
// Solidity: function setRegistry(address registryAddress) returns()
func (_GoldToken *GoldTokenSession) SetRegistry(registryAddress common.Address) (*types.Transaction, error) {
	return _GoldToken.Contract.SetRegistry(&_GoldToken.TransactOpts, registryAddress)
}

// SetRegistry is a paid mutator transaction binding the contract method 0xa91ee0dc.
//
// Solidity: function setRegistry(address registryAddress) returns()
func (_GoldToken *GoldTokenTransactorSession) SetRegistry(registryAddress common.Address) (*types.Transaction, error) {
	return _GoldToken.Contract.SetRegistry(&_GoldToken.TransactOpts, registryAddress)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) Transfer(opts *bind.TransactOpts, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "transfer", to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Transfer(&_GoldToken.TransactOpts, to, value)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) Transfer(to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.Transfer(&_GoldToken.TransactOpts, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "transferFrom", from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.TransferFrom(&_GoldToken.TransactOpts, from, to, value)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 value) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) TransferFrom(from common.Address, to common.Address, value *big.Int) (*types.Transaction, error) {
	return _GoldToken.Contract.TransferFrom(&_GoldToken.TransactOpts, from, to, value)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GoldToken *GoldTokenTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GoldToken *GoldTokenSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GoldToken.Contract.TransferOwnership(&_GoldToken.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_GoldToken *GoldTokenTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _GoldToken.Contract.TransferOwnership(&_GoldToken.TransactOpts, newOwner)
}

// TransferWithComment is a paid mutator transaction binding the contract method 0xe1d6aceb.
//
// Solidity: function transferWithComment(address to, uint256 value, string comment) returns(bool)
func (_GoldToken *GoldTokenTransactor) TransferWithComment(opts *bind.TransactOpts, to common.Address, value *big.Int, comment string) (*types.Transaction, error) {
	return _GoldToken.contract.Transact(opts, "transferWithComment", to, value, comment)
}

// TransferWithComment is a paid mutator transaction binding the contract method 0xe1d6aceb.
//
// Solidity: function transferWithComment(address to, uint256 value, string comment) returns(bool)
func (_GoldToken *GoldTokenSession) TransferWithComment(to common.Address, value *big.Int, comment string) (*types.Transaction, error) {
	return _GoldToken.Contract.TransferWithComment(&_GoldToken.TransactOpts, to, value, comment)
}

// TransferWithComment is a paid mutator transaction binding the contract method 0xe1d6aceb.
//
// Solidity: function transferWithComment(address to, uint256 value, string comment) returns(bool)
func (_GoldToken *GoldTokenTransactorSession) TransferWithComment(to common.Address, value *big.Int, comment string) (*types.Transaction, error) {
	return _GoldToken.Contract.TransferWithComment(&_GoldToken.TransactOpts, to, value, comment)
}

// GoldTokenApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the GoldToken contract.
type GoldTokenApprovalIterator struct {
	Event *GoldTokenApproval // Event containing the contract specifics and raw log

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
func (it *GoldTokenApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoldTokenApproval)
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
		it.Event = new(GoldTokenApproval)
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
func (it *GoldTokenApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoldTokenApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoldTokenApproval represents a Approval event raised by the GoldToken contract.
type GoldTokenApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_GoldToken *GoldTokenFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*GoldTokenApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _GoldToken.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &GoldTokenApprovalIterator{contract: _GoldToken.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_GoldToken *GoldTokenFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *GoldTokenApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _GoldToken.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoldTokenApproval)
				if err := _GoldToken.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_GoldToken *GoldTokenFilterer) ParseApproval(log types.Log) (*GoldTokenApproval, error) {
	event := new(GoldTokenApproval)
	if err := _GoldToken.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoldTokenOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the GoldToken contract.
type GoldTokenOwnershipTransferredIterator struct {
	Event *GoldTokenOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *GoldTokenOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoldTokenOwnershipTransferred)
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
		it.Event = new(GoldTokenOwnershipTransferred)
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
func (it *GoldTokenOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoldTokenOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoldTokenOwnershipTransferred represents a OwnershipTransferred event raised by the GoldToken contract.
type GoldTokenOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GoldToken *GoldTokenFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*GoldTokenOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GoldToken.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &GoldTokenOwnershipTransferredIterator{contract: _GoldToken.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_GoldToken *GoldTokenFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *GoldTokenOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _GoldToken.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoldTokenOwnershipTransferred)
				if err := _GoldToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_GoldToken *GoldTokenFilterer) ParseOwnershipTransferred(log types.Log) (*GoldTokenOwnershipTransferred, error) {
	event := new(GoldTokenOwnershipTransferred)
	if err := _GoldToken.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoldTokenRegistrySetIterator is returned from FilterRegistrySet and is used to iterate over the raw logs and unpacked data for RegistrySet events raised by the GoldToken contract.
type GoldTokenRegistrySetIterator struct {
	Event *GoldTokenRegistrySet // Event containing the contract specifics and raw log

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
func (it *GoldTokenRegistrySetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoldTokenRegistrySet)
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
		it.Event = new(GoldTokenRegistrySet)
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
func (it *GoldTokenRegistrySetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoldTokenRegistrySetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoldTokenRegistrySet represents a RegistrySet event raised by the GoldToken contract.
type GoldTokenRegistrySet struct {
	RegistryAddress common.Address
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterRegistrySet is a free log retrieval operation binding the contract event 0x27fe5f0c1c3b1ed427cc63d0f05759ffdecf9aec9e18d31ef366fc8a6cb5dc3b.
//
// Solidity: event RegistrySet(address indexed registryAddress)
func (_GoldToken *GoldTokenFilterer) FilterRegistrySet(opts *bind.FilterOpts, registryAddress []common.Address) (*GoldTokenRegistrySetIterator, error) {

	var registryAddressRule []interface{}
	for _, registryAddressItem := range registryAddress {
		registryAddressRule = append(registryAddressRule, registryAddressItem)
	}

	logs, sub, err := _GoldToken.contract.FilterLogs(opts, "RegistrySet", registryAddressRule)
	if err != nil {
		return nil, err
	}
	return &GoldTokenRegistrySetIterator{contract: _GoldToken.contract, event: "RegistrySet", logs: logs, sub: sub}, nil
}

// WatchRegistrySet is a free log subscription operation binding the contract event 0x27fe5f0c1c3b1ed427cc63d0f05759ffdecf9aec9e18d31ef366fc8a6cb5dc3b.
//
// Solidity: event RegistrySet(address indexed registryAddress)
func (_GoldToken *GoldTokenFilterer) WatchRegistrySet(opts *bind.WatchOpts, sink chan<- *GoldTokenRegistrySet, registryAddress []common.Address) (event.Subscription, error) {

	var registryAddressRule []interface{}
	for _, registryAddressItem := range registryAddress {
		registryAddressRule = append(registryAddressRule, registryAddressItem)
	}

	logs, sub, err := _GoldToken.contract.WatchLogs(opts, "RegistrySet", registryAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoldTokenRegistrySet)
				if err := _GoldToken.contract.UnpackLog(event, "RegistrySet", log); err != nil {
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

// ParseRegistrySet is a log parse operation binding the contract event 0x27fe5f0c1c3b1ed427cc63d0f05759ffdecf9aec9e18d31ef366fc8a6cb5dc3b.
//
// Solidity: event RegistrySet(address indexed registryAddress)
func (_GoldToken *GoldTokenFilterer) ParseRegistrySet(log types.Log) (*GoldTokenRegistrySet, error) {
	event := new(GoldTokenRegistrySet)
	if err := _GoldToken.contract.UnpackLog(event, "RegistrySet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoldTokenTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the GoldToken contract.
type GoldTokenTransferIterator struct {
	Event *GoldTokenTransfer // Event containing the contract specifics and raw log

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
func (it *GoldTokenTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoldTokenTransfer)
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
		it.Event = new(GoldTokenTransfer)
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
func (it *GoldTokenTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoldTokenTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoldTokenTransfer represents a Transfer event raised by the GoldToken contract.
type GoldTokenTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_GoldToken *GoldTokenFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*GoldTokenTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _GoldToken.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &GoldTokenTransferIterator{contract: _GoldToken.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_GoldToken *GoldTokenFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *GoldTokenTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _GoldToken.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoldTokenTransfer)
				if err := _GoldToken.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_GoldToken *GoldTokenFilterer) ParseTransfer(log types.Log) (*GoldTokenTransfer, error) {
	event := new(GoldTokenTransfer)
	if err := _GoldToken.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GoldTokenTransferCommentIterator is returned from FilterTransferComment and is used to iterate over the raw logs and unpacked data for TransferComment events raised by the GoldToken contract.
type GoldTokenTransferCommentIterator struct {
	Event *GoldTokenTransferComment // Event containing the contract specifics and raw log

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
func (it *GoldTokenTransferCommentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GoldTokenTransferComment)
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
		it.Event = new(GoldTokenTransferComment)
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
func (it *GoldTokenTransferCommentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GoldTokenTransferCommentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GoldTokenTransferComment represents a TransferComment event raised by the GoldToken contract.
type GoldTokenTransferComment struct {
	Comment string
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransferComment is a free log retrieval operation binding the contract event 0xe5d4e30fb8364e57bc4d662a07d0cf36f4c34552004c4c3624620a2c1d1c03dc.
//
// Solidity: event TransferComment(string comment)
func (_GoldToken *GoldTokenFilterer) FilterTransferComment(opts *bind.FilterOpts) (*GoldTokenTransferCommentIterator, error) {

	logs, sub, err := _GoldToken.contract.FilterLogs(opts, "TransferComment")
	if err != nil {
		return nil, err
	}
	return &GoldTokenTransferCommentIterator{contract: _GoldToken.contract, event: "TransferComment", logs: logs, sub: sub}, nil
}

// WatchTransferComment is a free log subscription operation binding the contract event 0xe5d4e30fb8364e57bc4d662a07d0cf36f4c34552004c4c3624620a2c1d1c03dc.
//
// Solidity: event TransferComment(string comment)
func (_GoldToken *GoldTokenFilterer) WatchTransferComment(opts *bind.WatchOpts, sink chan<- *GoldTokenTransferComment) (event.Subscription, error) {

	logs, sub, err := _GoldToken.contract.WatchLogs(opts, "TransferComment")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GoldTokenTransferComment)
				if err := _GoldToken.contract.UnpackLog(event, "TransferComment", log); err != nil {
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

// ParseTransferComment is a log parse operation binding the contract event 0xe5d4e30fb8364e57bc4d662a07d0cf36f4c34552004c4c3624620a2c1d1c03dc.
//
// Solidity: event TransferComment(string comment)
func (_GoldToken *GoldTokenFilterer) ParseTransferComment(log types.Log) (*GoldTokenTransferComment, error) {
	event := new(GoldTokenTransferComment)
	if err := _GoldToken.contract.UnpackLog(event, "TransferComment", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
