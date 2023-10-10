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

// FeeCurrencyMetaData contains all meta data concerning the FeeCurrency contract.
var FeeCurrencyMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"name_\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"symbol_\",\"type\":\"string\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"feeRecipient\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"communityFund\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"refund\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"tipTxFee\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"baseTxFee\",\"type\":\"uint256\"}],\"name\":\"creditGasFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"debitGasFees\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"subtractedValue\",\"type\":\"uint256\"}],\"name\":\"decreaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"addedValue\",\"type\":\"uint256\"}],\"name\":\"increaseAllowance\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// FeeCurrencyABI is the input ABI used to generate the binding from.
// Deprecated: Use FeeCurrencyMetaData.ABI instead.
var FeeCurrencyABI = FeeCurrencyMetaData.ABI

// FeeCurrency is an auto generated Go binding around an Ethereum contract.
type FeeCurrency struct {
	FeeCurrencyCaller     // Read-only binding to the contract
	FeeCurrencyTransactor // Write-only binding to the contract
	FeeCurrencyFilterer   // Log filterer for contract events
}

// FeeCurrencyCaller is an auto generated read-only Go binding around an Ethereum contract.
type FeeCurrencyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeCurrencyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type FeeCurrencyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeCurrencyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type FeeCurrencyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// FeeCurrencySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type FeeCurrencySession struct {
	Contract     *FeeCurrency      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// FeeCurrencyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type FeeCurrencyCallerSession struct {
	Contract *FeeCurrencyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// FeeCurrencyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type FeeCurrencyTransactorSession struct {
	Contract     *FeeCurrencyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// FeeCurrencyRaw is an auto generated low-level Go binding around an Ethereum contract.
type FeeCurrencyRaw struct {
	Contract *FeeCurrency // Generic contract binding to access the raw methods on
}

// FeeCurrencyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type FeeCurrencyCallerRaw struct {
	Contract *FeeCurrencyCaller // Generic read-only contract binding to access the raw methods on
}

// FeeCurrencyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type FeeCurrencyTransactorRaw struct {
	Contract *FeeCurrencyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewFeeCurrency creates a new instance of FeeCurrency, bound to a specific deployed contract.
func NewFeeCurrency(address common.Address, backend bind.ContractBackend) (*FeeCurrency, error) {
	contract, err := bindFeeCurrency(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &FeeCurrency{FeeCurrencyCaller: FeeCurrencyCaller{contract: contract}, FeeCurrencyTransactor: FeeCurrencyTransactor{contract: contract}, FeeCurrencyFilterer: FeeCurrencyFilterer{contract: contract}}, nil
}

// NewFeeCurrencyCaller creates a new read-only instance of FeeCurrency, bound to a specific deployed contract.
func NewFeeCurrencyCaller(address common.Address, caller bind.ContractCaller) (*FeeCurrencyCaller, error) {
	contract, err := bindFeeCurrency(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyCaller{contract: contract}, nil
}

// NewFeeCurrencyTransactor creates a new write-only instance of FeeCurrency, bound to a specific deployed contract.
func NewFeeCurrencyTransactor(address common.Address, transactor bind.ContractTransactor) (*FeeCurrencyTransactor, error) {
	contract, err := bindFeeCurrency(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyTransactor{contract: contract}, nil
}

// NewFeeCurrencyFilterer creates a new log filterer instance of FeeCurrency, bound to a specific deployed contract.
func NewFeeCurrencyFilterer(address common.Address, filterer bind.ContractFilterer) (*FeeCurrencyFilterer, error) {
	contract, err := bindFeeCurrency(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyFilterer{contract: contract}, nil
}

// bindFeeCurrency binds a generic wrapper to an already deployed contract.
func bindFeeCurrency(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := FeeCurrencyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeCurrency *FeeCurrencyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeCurrency.Contract.FeeCurrencyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeCurrency *FeeCurrencyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeCurrency.Contract.FeeCurrencyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeCurrency *FeeCurrencyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeCurrency.Contract.FeeCurrencyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_FeeCurrency *FeeCurrencyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _FeeCurrency.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_FeeCurrency *FeeCurrencyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _FeeCurrency.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_FeeCurrency *FeeCurrencyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _FeeCurrency.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FeeCurrency *FeeCurrencyCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeCurrency.contract.Call(opts, &out, "allowance", owner, spender)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FeeCurrency *FeeCurrencySession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FeeCurrency.Contract.Allowance(&_FeeCurrency.CallOpts, owner, spender)
}

// Allowance is a free data retrieval call binding the contract method 0xdd62ed3e.
//
// Solidity: function allowance(address owner, address spender) view returns(uint256)
func (_FeeCurrency *FeeCurrencyCallerSession) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	return _FeeCurrency.Contract.Allowance(&_FeeCurrency.CallOpts, owner, spender)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FeeCurrency *FeeCurrencyCaller) BalanceOf(opts *bind.CallOpts, account common.Address) (*big.Int, error) {
	var out []interface{}
	err := _FeeCurrency.contract.Call(opts, &out, "balanceOf", account)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FeeCurrency *FeeCurrencySession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FeeCurrency.Contract.BalanceOf(&_FeeCurrency.CallOpts, account)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address account) view returns(uint256)
func (_FeeCurrency *FeeCurrencyCallerSession) BalanceOf(account common.Address) (*big.Int, error) {
	return _FeeCurrency.Contract.BalanceOf(&_FeeCurrency.CallOpts, account)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FeeCurrency *FeeCurrencyCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _FeeCurrency.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FeeCurrency *FeeCurrencySession) Decimals() (uint8, error) {
	return _FeeCurrency.Contract.Decimals(&_FeeCurrency.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_FeeCurrency *FeeCurrencyCallerSession) Decimals() (uint8, error) {
	return _FeeCurrency.Contract.Decimals(&_FeeCurrency.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FeeCurrency *FeeCurrencyCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeCurrency.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FeeCurrency *FeeCurrencySession) Name() (string, error) {
	return _FeeCurrency.Contract.Name(&_FeeCurrency.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_FeeCurrency *FeeCurrencyCallerSession) Name() (string, error) {
	return _FeeCurrency.Contract.Name(&_FeeCurrency.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FeeCurrency *FeeCurrencyCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _FeeCurrency.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FeeCurrency *FeeCurrencySession) Symbol() (string, error) {
	return _FeeCurrency.Contract.Symbol(&_FeeCurrency.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_FeeCurrency *FeeCurrencyCallerSession) Symbol() (string, error) {
	return _FeeCurrency.Contract.Symbol(&_FeeCurrency.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FeeCurrency *FeeCurrencyCaller) TotalSupply(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _FeeCurrency.contract.Call(opts, &out, "totalSupply")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FeeCurrency *FeeCurrencySession) TotalSupply() (*big.Int, error) {
	return _FeeCurrency.Contract.TotalSupply(&_FeeCurrency.CallOpts)
}

// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
//
// Solidity: function totalSupply() view returns(uint256)
func (_FeeCurrency *FeeCurrencyCallerSession) TotalSupply() (*big.Int, error) {
	return _FeeCurrency.Contract.TotalSupply(&_FeeCurrency.CallOpts)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactor) Approve(opts *bind.TransactOpts, spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "approve", spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencySession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.Approve(&_FeeCurrency.TransactOpts, spender, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address spender, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactorSession) Approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.Approve(&_FeeCurrency.TransactOpts, spender, amount)
}

// CreditGasFees is a paid mutator transaction binding the contract method 0x6a30b253.
//
// Solidity: function creditGasFees(address from, address feeRecipient, address , address communityFund, uint256 refund, uint256 tipTxFee, uint256 , uint256 baseTxFee) returns()
func (_FeeCurrency *FeeCurrencyTransactor) CreditGasFees(opts *bind.TransactOpts, from common.Address, feeRecipient common.Address, arg2 common.Address, communityFund common.Address, refund *big.Int, tipTxFee *big.Int, arg6 *big.Int, baseTxFee *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "creditGasFees", from, feeRecipient, arg2, communityFund, refund, tipTxFee, arg6, baseTxFee)
}

// CreditGasFees is a paid mutator transaction binding the contract method 0x6a30b253.
//
// Solidity: function creditGasFees(address from, address feeRecipient, address , address communityFund, uint256 refund, uint256 tipTxFee, uint256 , uint256 baseTxFee) returns()
func (_FeeCurrency *FeeCurrencySession) CreditGasFees(from common.Address, feeRecipient common.Address, arg2 common.Address, communityFund common.Address, refund *big.Int, tipTxFee *big.Int, arg6 *big.Int, baseTxFee *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.CreditGasFees(&_FeeCurrency.TransactOpts, from, feeRecipient, arg2, communityFund, refund, tipTxFee, arg6, baseTxFee)
}

// CreditGasFees is a paid mutator transaction binding the contract method 0x6a30b253.
//
// Solidity: function creditGasFees(address from, address feeRecipient, address , address communityFund, uint256 refund, uint256 tipTxFee, uint256 , uint256 baseTxFee) returns()
func (_FeeCurrency *FeeCurrencyTransactorSession) CreditGasFees(from common.Address, feeRecipient common.Address, arg2 common.Address, communityFund common.Address, refund *big.Int, tipTxFee *big.Int, arg6 *big.Int, baseTxFee *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.CreditGasFees(&_FeeCurrency.TransactOpts, from, feeRecipient, arg2, communityFund, refund, tipTxFee, arg6, baseTxFee)
}

// DebitGasFees is a paid mutator transaction binding the contract method 0x58cf9672.
//
// Solidity: function debitGasFees(address from, uint256 value) returns()
func (_FeeCurrency *FeeCurrencyTransactor) DebitGasFees(opts *bind.TransactOpts, from common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "debitGasFees", from, value)
}

// DebitGasFees is a paid mutator transaction binding the contract method 0x58cf9672.
//
// Solidity: function debitGasFees(address from, uint256 value) returns()
func (_FeeCurrency *FeeCurrencySession) DebitGasFees(from common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.DebitGasFees(&_FeeCurrency.TransactOpts, from, value)
}

// DebitGasFees is a paid mutator transaction binding the contract method 0x58cf9672.
//
// Solidity: function debitGasFees(address from, uint256 value) returns()
func (_FeeCurrency *FeeCurrencyTransactorSession) DebitGasFees(from common.Address, value *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.DebitGasFees(&_FeeCurrency.TransactOpts, from, value)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactor) DecreaseAllowance(opts *bind.TransactOpts, spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "decreaseAllowance", spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_FeeCurrency *FeeCurrencySession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.DecreaseAllowance(&_FeeCurrency.TransactOpts, spender, subtractedValue)
}

// DecreaseAllowance is a paid mutator transaction binding the contract method 0xa457c2d7.
//
// Solidity: function decreaseAllowance(address spender, uint256 subtractedValue) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactorSession) DecreaseAllowance(spender common.Address, subtractedValue *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.DecreaseAllowance(&_FeeCurrency.TransactOpts, spender, subtractedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactor) IncreaseAllowance(opts *bind.TransactOpts, spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "increaseAllowance", spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_FeeCurrency *FeeCurrencySession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.IncreaseAllowance(&_FeeCurrency.TransactOpts, spender, addedValue)
}

// IncreaseAllowance is a paid mutator transaction binding the contract method 0x39509351.
//
// Solidity: function increaseAllowance(address spender, uint256 addedValue) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactorSession) IncreaseAllowance(spender common.Address, addedValue *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.IncreaseAllowance(&_FeeCurrency.TransactOpts, spender, addedValue)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactor) Transfer(opts *bind.TransactOpts, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "transfer", to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencySession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.Transfer(&_FeeCurrency.TransactOpts, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0xa9059cbb.
//
// Solidity: function transfer(address to, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactorSession) Transfer(to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.Transfer(&_FeeCurrency.TransactOpts, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.contract.Transact(opts, "transferFrom", from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencySession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.TransferFrom(&_FeeCurrency.TransactOpts, from, to, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 amount) returns(bool)
func (_FeeCurrency *FeeCurrencyTransactorSession) TransferFrom(from common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _FeeCurrency.Contract.TransferFrom(&_FeeCurrency.TransactOpts, from, to, amount)
}

// FeeCurrencyApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the FeeCurrency contract.
type FeeCurrencyApprovalIterator struct {
	Event *FeeCurrencyApproval // Event containing the contract specifics and raw log

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
func (it *FeeCurrencyApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeCurrencyApproval)
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
		it.Event = new(FeeCurrencyApproval)
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
func (it *FeeCurrencyApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeCurrencyApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeCurrencyApproval represents a Approval event raised by the FeeCurrency contract.
type FeeCurrencyApproval struct {
	Owner   common.Address
	Spender common.Address
	Value   *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FeeCurrency *FeeCurrencyFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address) (*FeeCurrencyApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FeeCurrency.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyApprovalIterator{contract: _FeeCurrency.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 value)
func (_FeeCurrency *FeeCurrencyFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *FeeCurrencyApproval, owner []common.Address, spender []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}

	logs, sub, err := _FeeCurrency.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeCurrencyApproval)
				if err := _FeeCurrency.contract.UnpackLog(event, "Approval", log); err != nil {
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
func (_FeeCurrency *FeeCurrencyFilterer) ParseApproval(log types.Log) (*FeeCurrencyApproval, error) {
	event := new(FeeCurrencyApproval)
	if err := _FeeCurrency.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// FeeCurrencyTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the FeeCurrency contract.
type FeeCurrencyTransferIterator struct {
	Event *FeeCurrencyTransfer // Event containing the contract specifics and raw log

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
func (it *FeeCurrencyTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(FeeCurrencyTransfer)
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
		it.Event = new(FeeCurrencyTransfer)
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
func (it *FeeCurrencyTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *FeeCurrencyTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// FeeCurrencyTransfer represents a Transfer event raised by the FeeCurrency contract.
type FeeCurrencyTransfer struct {
	From  common.Address
	To    common.Address
	Value *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FeeCurrency *FeeCurrencyFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address) (*FeeCurrencyTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeCurrency.contract.FilterLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return &FeeCurrencyTransferIterator{contract: _FeeCurrency.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 value)
func (_FeeCurrency *FeeCurrencyFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *FeeCurrencyTransfer, from []common.Address, to []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}

	logs, sub, err := _FeeCurrency.contract.WatchLogs(opts, "Transfer", fromRule, toRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(FeeCurrencyTransfer)
				if err := _FeeCurrency.contract.UnpackLog(event, "Transfer", log); err != nil {
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
func (_FeeCurrency *FeeCurrencyFilterer) ParseTransfer(log types.Log) (*FeeCurrencyTransfer, error) {
	event := new(FeeCurrencyTransfer)
	if err := _FeeCurrency.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
