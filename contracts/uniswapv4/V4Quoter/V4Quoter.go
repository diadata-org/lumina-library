// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package v4quoter

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

// IV4QuoterQuoteExactParams is an auto generated low-level Go binding around an user-defined struct.
type IV4QuoterQuoteExactParams struct {
	ExactCurrency common.Address
	Path          []PathKey
	ExactAmount   *big.Int
}

// IV4QuoterQuoteExactSingleParams is an auto generated low-level Go binding around an user-defined struct.
type IV4QuoterQuoteExactSingleParams struct {
	PoolKey     PoolKey
	ZeroForOne  bool
	ExactAmount *big.Int
	HookData    []byte
}

// PathKey is an auto generated low-level Go binding around an user-defined struct.
type PathKey struct {
	IntermediateCurrency common.Address
	Fee                  *big.Int
	TickSpacing          *big.Int
	Hooks                common.Address
	HookData             []byte
}

// PoolKey is an auto generated low-level Go binding around an user-defined struct.
type PoolKey struct {
	Currency0   common.Address
	Currency1   common.Address
	Fee         *big.Int
	TickSpacing *big.Int
	Hooks       common.Address
}

// V4QuoterMetaData contains all meta data concerning the V4Quoter contract.
var V4QuoterMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_poolManager\",\"type\":\"address\",\"internalType\":\"contractIPoolManager\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_quoteExactInput\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactParams\",\"components\":[{\"name\":\"exactCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structPathKey[]\",\"components\":[{\"name\":\"intermediateCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_quoteExactInputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactSingleParams\",\"components\":[{\"name\":\"poolKey\",\"type\":\"tuple\",\"internalType\":\"structPoolKey\",\"components\":[{\"name\":\"currency0\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"currency1\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"}]},{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"_quoteExactOutput\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactParams\",\"components\":[{\"name\":\"exactCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structPathKey[]\",\"components\":[{\"name\":\"intermediateCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"_quoteExactOutputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactSingleParams\",\"components\":[{\"name\":\"poolKey\",\"type\":\"tuple\",\"internalType\":\"structPoolKey\",\"components\":[{\"name\":\"currency0\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"currency1\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"}]},{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"msgSender\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"poolManager\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"contractIPoolManager\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quoteExactInput\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactParams\",\"components\":[{\"name\":\"exactCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structPathKey[]\",\"components\":[{\"name\":\"intermediateCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteExactInputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactSingleParams\",\"components\":[{\"name\":\"poolKey\",\"type\":\"tuple\",\"internalType\":\"structPoolKey\",\"components\":[{\"name\":\"currency0\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"currency1\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"}]},{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"amountOut\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"quoteExactOutput\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactParams\",\"components\":[{\"name\":\"exactCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"path\",\"type\":\"tuple[]\",\"internalType\":\"structPathKey[]\",\"components\":[{\"name\":\"intermediateCurrency\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"quoteExactOutputSingle\",\"inputs\":[{\"name\":\"params\",\"type\":\"tuple\",\"internalType\":\"structIV4Quoter.QuoteExactSingleParams\",\"components\":[{\"name\":\"poolKey\",\"type\":\"tuple\",\"internalType\":\"structPoolKey\",\"components\":[{\"name\":\"currency0\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"currency1\",\"type\":\"address\",\"internalType\":\"Currency\"},{\"name\":\"fee\",\"type\":\"uint24\",\"internalType\":\"uint24\"},{\"name\":\"tickSpacing\",\"type\":\"int24\",\"internalType\":\"int24\"},{\"name\":\"hooks\",\"type\":\"address\",\"internalType\":\"contractIHooks\"}]},{\"name\":\"zeroForOne\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"exactAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"hookData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]}],\"outputs\":[{\"name\":\"amountIn\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"gasEstimate\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"unlockCallback\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"error\",\"name\":\"NotEnoughLiquidity\",\"inputs\":[{\"name\":\"poolId\",\"type\":\"bytes32\",\"internalType\":\"PoolId\"}]},{\"type\":\"error\",\"name\":\"NotPoolManager\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotSelf\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"QuoteSwap\",\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UnexpectedCallSuccess\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UnexpectedRevertBytes\",\"inputs\":[{\"name\":\"revertData\",\"type\":\"bytes\",\"internalType\":\"bytes\"}]},{\"inputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"address\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"}],\"name\":\"getPoolId\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// V4QuoterABI is the input ABI used to generate the binding from.
// Deprecated: Use V4QuoterMetaData.ABI instead.
var V4QuoterABI = V4QuoterMetaData.ABI

// V4Quoter is an auto generated Go binding around an Ethereum contract.
type V4Quoter struct {
	V4QuoterCaller     // Read-only binding to the contract
	V4QuoterTransactor // Write-only binding to the contract
	V4QuoterFilterer   // Log filterer for contract events
}

// V4QuoterCaller is an auto generated read-only Go binding around an Ethereum contract.
type V4QuoterCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V4QuoterTransactor is an auto generated write-only Go binding around an Ethereum contract.
type V4QuoterTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V4QuoterFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type V4QuoterFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// V4QuoterSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type V4QuoterSession struct {
	Contract     *V4Quoter         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// V4QuoterCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type V4QuoterCallerSession struct {
	Contract *V4QuoterCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// V4QuoterTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type V4QuoterTransactorSession struct {
	Contract     *V4QuoterTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// V4QuoterRaw is an auto generated low-level Go binding around an Ethereum contract.
type V4QuoterRaw struct {
	Contract *V4Quoter // Generic contract binding to access the raw methods on
}

// V4QuoterCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type V4QuoterCallerRaw struct {
	Contract *V4QuoterCaller // Generic read-only contract binding to access the raw methods on
}

// V4QuoterTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type V4QuoterTransactorRaw struct {
	Contract *V4QuoterTransactor // Generic write-only contract binding to access the raw methods on
}

// NewV4Quoter creates a new instance of V4Quoter, bound to a specific deployed contract.
func NewV4Quoter(address common.Address, backend bind.ContractBackend) (*V4Quoter, error) {
	contract, err := bindV4Quoter(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &V4Quoter{V4QuoterCaller: V4QuoterCaller{contract: contract}, V4QuoterTransactor: V4QuoterTransactor{contract: contract}, V4QuoterFilterer: V4QuoterFilterer{contract: contract}}, nil
}

// NewV4QuoterCaller creates a new read-only instance of V4Quoter, bound to a specific deployed contract.
func NewV4QuoterCaller(address common.Address, caller bind.ContractCaller) (*V4QuoterCaller, error) {
	contract, err := bindV4Quoter(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &V4QuoterCaller{contract: contract}, nil
}

// NewV4QuoterTransactor creates a new write-only instance of V4Quoter, bound to a specific deployed contract.
func NewV4QuoterTransactor(address common.Address, transactor bind.ContractTransactor) (*V4QuoterTransactor, error) {
	contract, err := bindV4Quoter(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &V4QuoterTransactor{contract: contract}, nil
}

// NewV4QuoterFilterer creates a new log filterer instance of V4Quoter, bound to a specific deployed contract.
func NewV4QuoterFilterer(address common.Address, filterer bind.ContractFilterer) (*V4QuoterFilterer, error) {
	contract, err := bindV4Quoter(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &V4QuoterFilterer{contract: contract}, nil
}

// bindV4Quoter binds a generic wrapper to an already deployed contract.
func bindV4Quoter(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := V4QuoterMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V4Quoter *V4QuoterRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V4Quoter.Contract.V4QuoterCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V4Quoter *V4QuoterRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V4Quoter.Contract.V4QuoterTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V4Quoter *V4QuoterRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V4Quoter.Contract.V4QuoterTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_V4Quoter *V4QuoterCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _V4Quoter.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_V4Quoter *V4QuoterTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _V4Quoter.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_V4Quoter *V4QuoterTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _V4Quoter.Contract.contract.Transact(opts, method, params...)
}

// QuoteExactInputSingleInternal is a free data retrieval call binding the contract method 0xeebe0c6a.
//
// Solidity: function _quoteExactInputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) view returns(bytes)
func (_V4Quoter *V4QuoterCaller) QuoteExactInputSingleInternal(opts *bind.CallOpts, params IV4QuoterQuoteExactSingleParams) ([]byte, error) {
	var out []interface{}
	err := _V4Quoter.contract.Call(opts, &out, "_quoteExactInputSingle", params)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// QuoteExactInputSingleInternal is a free data retrieval call binding the contract method 0xeebe0c6a.
//
// Solidity: function _quoteExactInputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) view returns(bytes)
func (_V4Quoter *V4QuoterSession) QuoteExactInputSingleInternal(params IV4QuoterQuoteExactSingleParams) ([]byte, error) {
	return _V4Quoter.Contract.QuoteExactInputSingleInternal(&_V4Quoter.CallOpts, params)
}

// QuoteExactInputSingleInternal is a free data retrieval call binding the contract method 0xeebe0c6a.
//
// Solidity: function _quoteExactInputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) view returns(bytes)
func (_V4Quoter *V4QuoterCallerSession) QuoteExactInputSingleInternal(params IV4QuoterQuoteExactSingleParams) ([]byte, error) {
	return _V4Quoter.Contract.QuoteExactInputSingleInternal(&_V4Quoter.CallOpts, params)
}

// GetPoolId is a free data retrieval call binding the contract method 0x9c6b88a8.
//
// Solidity: function getPoolId((address,address,uint24,int24,address) key) view returns(bytes32)
func (_V4Quoter *V4QuoterCaller) GetPoolId(opts *bind.CallOpts, key PoolKey) ([32]byte, error) {
	var out []interface{}
	err := _V4Quoter.contract.Call(opts, &out, "getPoolId", key)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetPoolId is a free data retrieval call binding the contract method 0x9c6b88a8.
//
// Solidity: function getPoolId((address,address,uint24,int24,address) key) view returns(bytes32)
func (_V4Quoter *V4QuoterSession) GetPoolId(key PoolKey) ([32]byte, error) {
	return _V4Quoter.Contract.GetPoolId(&_V4Quoter.CallOpts, key)
}

// GetPoolId is a free data retrieval call binding the contract method 0x9c6b88a8.
//
// Solidity: function getPoolId((address,address,uint24,int24,address) key) view returns(bytes32)
func (_V4Quoter *V4QuoterCallerSession) GetPoolId(key PoolKey) ([32]byte, error) {
	return _V4Quoter.Contract.GetPoolId(&_V4Quoter.CallOpts, key)
}

// MsgSender is a free data retrieval call binding the contract method 0xd737d0c7.
//
// Solidity: function msgSender() view returns(address)
func (_V4Quoter *V4QuoterCaller) MsgSender(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4Quoter.contract.Call(opts, &out, "msgSender")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MsgSender is a free data retrieval call binding the contract method 0xd737d0c7.
//
// Solidity: function msgSender() view returns(address)
func (_V4Quoter *V4QuoterSession) MsgSender() (common.Address, error) {
	return _V4Quoter.Contract.MsgSender(&_V4Quoter.CallOpts)
}

// MsgSender is a free data retrieval call binding the contract method 0xd737d0c7.
//
// Solidity: function msgSender() view returns(address)
func (_V4Quoter *V4QuoterCallerSession) MsgSender() (common.Address, error) {
	return _V4Quoter.Contract.MsgSender(&_V4Quoter.CallOpts)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_V4Quoter *V4QuoterCaller) PoolManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _V4Quoter.contract.Call(opts, &out, "poolManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_V4Quoter *V4QuoterSession) PoolManager() (common.Address, error) {
	return _V4Quoter.Contract.PoolManager(&_V4Quoter.CallOpts)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_V4Quoter *V4QuoterCallerSession) PoolManager() (common.Address, error) {
	return _V4Quoter.Contract.PoolManager(&_V4Quoter.CallOpts)
}

// QuoteExactInputSingle is a free data retrieval call binding the contract method 0xaa9d21cb.
//
// Solidity: function quoteExactInputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) view returns(uint256 amountOut, uint256 gasEstimate)
func (_V4Quoter *V4QuoterCaller) QuoteExactInputSingle(opts *bind.CallOpts, params IV4QuoterQuoteExactSingleParams) (struct {
	AmountOut   *big.Int
	GasEstimate *big.Int
}, error) {
	var out []interface{}
	err := _V4Quoter.contract.Call(opts, &out, "quoteExactInputSingle", params)

	outstruct := new(struct {
		AmountOut   *big.Int
		GasEstimate *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.AmountOut = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.GasEstimate = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// QuoteExactInputSingle is a free data retrieval call binding the contract method 0xaa9d21cb.
//
// Solidity: function quoteExactInputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) view returns(uint256 amountOut, uint256 gasEstimate)
func (_V4Quoter *V4QuoterSession) QuoteExactInputSingle(params IV4QuoterQuoteExactSingleParams) (struct {
	AmountOut   *big.Int
	GasEstimate *big.Int
}, error) {
	return _V4Quoter.Contract.QuoteExactInputSingle(&_V4Quoter.CallOpts, params)
}

// QuoteExactInputSingle is a free data retrieval call binding the contract method 0xaa9d21cb.
//
// Solidity: function quoteExactInputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) view returns(uint256 amountOut, uint256 gasEstimate)
func (_V4Quoter *V4QuoterCallerSession) QuoteExactInputSingle(params IV4QuoterQuoteExactSingleParams) (struct {
	AmountOut   *big.Int
	GasEstimate *big.Int
}, error) {
	return _V4Quoter.Contract.QuoteExactInputSingle(&_V4Quoter.CallOpts, params)
}

// QuoteExactInputInternal is a paid mutator transaction binding the contract method 0x6a36a38c.
//
// Solidity: function _quoteExactInput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(bytes)
func (_V4Quoter *V4QuoterTransactor) QuoteExactInputInternal(opts *bind.TransactOpts, params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "_quoteExactInput", params)
}

// QuoteExactInputInternal is a paid mutator transaction binding the contract method 0x6a36a38c.
//
// Solidity: function _quoteExactInput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(bytes)
func (_V4Quoter *V4QuoterSession) QuoteExactInputInternal(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactInputInternal(&_V4Quoter.TransactOpts, params)
}

// QuoteExactInputInternal is a paid mutator transaction binding the contract method 0x6a36a38c.
//
// Solidity: function _quoteExactInput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(bytes)
func (_V4Quoter *V4QuoterTransactorSession) QuoteExactInputInternal(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactInputInternal(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutputInternal is a paid mutator transaction binding the contract method 0xaa2f1501.
//
// Solidity: function _quoteExactOutput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(bytes)
func (_V4Quoter *V4QuoterTransactor) QuoteExactOutputInternal(opts *bind.TransactOpts, params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "_quoteExactOutput", params)
}

// QuoteExactOutputInternal is a paid mutator transaction binding the contract method 0xaa2f1501.
//
// Solidity: function _quoteExactOutput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(bytes)
func (_V4Quoter *V4QuoterSession) QuoteExactOutputInternal(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutputInternal(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutputInternal is a paid mutator transaction binding the contract method 0xaa2f1501.
//
// Solidity: function _quoteExactOutput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(bytes)
func (_V4Quoter *V4QuoterTransactorSession) QuoteExactOutputInternal(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutputInternal(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutputSingleInternal is a paid mutator transaction binding the contract method 0x595323f5.
//
// Solidity: function _quoteExactOutputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) returns(bytes)
func (_V4Quoter *V4QuoterTransactor) QuoteExactOutputSingleInternal(opts *bind.TransactOpts, params IV4QuoterQuoteExactSingleParams) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "_quoteExactOutputSingle", params)
}

// QuoteExactOutputSingleInternal is a paid mutator transaction binding the contract method 0x595323f5.
//
// Solidity: function _quoteExactOutputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) returns(bytes)
func (_V4Quoter *V4QuoterSession) QuoteExactOutputSingleInternal(params IV4QuoterQuoteExactSingleParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutputSingleInternal(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutputSingleInternal is a paid mutator transaction binding the contract method 0x595323f5.
//
// Solidity: function _quoteExactOutputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) returns(bytes)
func (_V4Quoter *V4QuoterTransactorSession) QuoteExactOutputSingleInternal(params IV4QuoterQuoteExactSingleParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutputSingleInternal(&_V4Quoter.TransactOpts, params)
}

// QuoteExactInputV1 is a paid mutator transaction binding the contract method 0xca253dc9.
//
// Solidity: function quoteExactInput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(uint256 amountOut, uint256 gasEstimate)
func (_V4Quoter *V4QuoterTransactor) QuoteExactInputV1(opts *bind.TransactOpts, params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "quoteExactInput", params)
}

// QuoteExactInputV1 is a paid mutator transaction binding the contract method 0xca253dc9.
//
// Solidity: function quoteExactInput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(uint256 amountOut, uint256 gasEstimate)
func (_V4Quoter *V4QuoterSession) QuoteExactInputV1(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactInputV1(&_V4Quoter.TransactOpts, params)
}

// QuoteExactInputV1 is a paid mutator transaction binding the contract method 0xca253dc9.
//
// Solidity: function quoteExactInput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(uint256 amountOut, uint256 gasEstimate)
func (_V4Quoter *V4QuoterTransactorSession) QuoteExactInputV1(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactInputV1(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x147d2af9.
//
// Solidity: function quoteExactOutput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(uint256 amountIn, uint256 gasEstimate)
func (_V4Quoter *V4QuoterTransactor) QuoteExactOutput(opts *bind.TransactOpts, params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "quoteExactOutput", params)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x147d2af9.
//
// Solidity: function quoteExactOutput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(uint256 amountIn, uint256 gasEstimate)
func (_V4Quoter *V4QuoterSession) QuoteExactOutput(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutput(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutput is a paid mutator transaction binding the contract method 0x147d2af9.
//
// Solidity: function quoteExactOutput((address,(address,uint24,int24,address,bytes)[],uint128) params) returns(uint256 amountIn, uint256 gasEstimate)
func (_V4Quoter *V4QuoterTransactorSession) QuoteExactOutput(params IV4QuoterQuoteExactParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutput(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutputSingleV1 is a paid mutator transaction binding the contract method 0x58733073.
//
// Solidity: function quoteExactOutputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) returns(uint256 amountIn, uint256 gasEstimate)
func (_V4Quoter *V4QuoterTransactor) QuoteExactOutputSingleV1(opts *bind.TransactOpts, params IV4QuoterQuoteExactSingleParams) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "quoteExactOutputSingle", params)
}

// QuoteExactOutputSingleV1 is a paid mutator transaction binding the contract method 0x58733073.
//
// Solidity: function quoteExactOutputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) returns(uint256 amountIn, uint256 gasEstimate)
func (_V4Quoter *V4QuoterSession) QuoteExactOutputSingleV1(params IV4QuoterQuoteExactSingleParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutputSingleV1(&_V4Quoter.TransactOpts, params)
}

// QuoteExactOutputSingleV1 is a paid mutator transaction binding the contract method 0x58733073.
//
// Solidity: function quoteExactOutputSingle(((address,address,uint24,int24,address),bool,uint128,bytes) params) returns(uint256 amountIn, uint256 gasEstimate)
func (_V4Quoter *V4QuoterTransactorSession) QuoteExactOutputSingleV1(params IV4QuoterQuoteExactSingleParams) (*types.Transaction, error) {
	return _V4Quoter.Contract.QuoteExactOutputSingleV1(&_V4Quoter.TransactOpts, params)
}

// UnlockCallback is a paid mutator transaction binding the contract method 0x91dd7346.
//
// Solidity: function unlockCallback(bytes data) returns(bytes)
func (_V4Quoter *V4QuoterTransactor) UnlockCallback(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _V4Quoter.contract.Transact(opts, "unlockCallback", data)
}

// UnlockCallback is a paid mutator transaction binding the contract method 0x91dd7346.
//
// Solidity: function unlockCallback(bytes data) returns(bytes)
func (_V4Quoter *V4QuoterSession) UnlockCallback(data []byte) (*types.Transaction, error) {
	return _V4Quoter.Contract.UnlockCallback(&_V4Quoter.TransactOpts, data)
}

// UnlockCallback is a paid mutator transaction binding the contract method 0x91dd7346.
//
// Solidity: function unlockCallback(bytes data) returns(bytes)
func (_V4Quoter *V4QuoterTransactorSession) UnlockCallback(data []byte) (*types.Transaction, error) {
	return _V4Quoter.Contract.UnlockCallback(&_V4Quoter.TransactOpts, data)
}
