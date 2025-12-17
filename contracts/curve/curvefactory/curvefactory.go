// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package curvefactory

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

// CurvefactoryMetaData contains all meta data concerning the Curvefactory contract.
var CurvefactoryMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"CryptoPoolDeployed\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"indexed\":false},{\"name\":\"coins\",\"type\":\"address[2]\",\"indexed\":false},{\"name\":\"A\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"gamma\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"mid_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"out_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"allowed_extra_profit\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"fee_gamma\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"adjustment_step\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"admin_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"ma_half_time\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"initial_price\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"deployer\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"LiquidityGaugeDeployed\",\"inputs\":[{\"name\":\"pool\",\"type\":\"address\",\"indexed\":false},{\"name\":\"token\",\"type\":\"address\",\"indexed\":false},{\"name\":\"gauge\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"UpdateFeeReceiver\",\"inputs\":[{\"name\":\"_old_fee_receiver\",\"type\":\"address\",\"indexed\":false},{\"name\":\"_new_fee_receiver\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"UpdatePoolImplementation\",\"inputs\":[{\"name\":\"_old_pool_implementation\",\"type\":\"address\",\"indexed\":false},{\"name\":\"_new_pool_implementation\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"UpdateTokenImplementation\",\"inputs\":[{\"name\":\"_old_token_implementation\",\"type\":\"address\",\"indexed\":false},{\"name\":\"_new_token_implementation\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"UpdateGaugeImplementation\",\"inputs\":[{\"name\":\"_old_gauge_implementation\",\"type\":\"address\",\"indexed\":false},{\"name\":\"_new_gauge_implementation\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"TransferOwnership\",\"inputs\":[{\"name\":\"_old_owner\",\"type\":\"address\",\"indexed\":false},{\"name\":\"_new_owner\",\"type\":\"address\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"constructor\",\"inputs\":[{\"name\":\"_fee_receiver\",\"type\":\"address\"},{\"name\":\"_pool_implementation\",\"type\":\"address\"},{\"name\":\"_token_implementation\",\"type\":\"address\"},{\"name\":\"_gauge_implementation\",\"type\":\"address\"},{\"name\":\"_weth\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"deploy_pool\",\"inputs\":[{\"name\":\"_name\",\"type\":\"string\"},{\"name\":\"_symbol\",\"type\":\"string\"},{\"name\":\"_coins\",\"type\":\"address[2]\"},{\"name\":\"A\",\"type\":\"uint256\"},{\"name\":\"gamma\",\"type\":\"uint256\"},{\"name\":\"mid_fee\",\"type\":\"uint256\"},{\"name\":\"out_fee\",\"type\":\"uint256\"},{\"name\":\"allowed_extra_profit\",\"type\":\"uint256\"},{\"name\":\"fee_gamma\",\"type\":\"uint256\"},{\"name\":\"adjustment_step\",\"type\":\"uint256\"},{\"name\":\"admin_fee\",\"type\":\"uint256\"},{\"name\":\"ma_half_time\",\"type\":\"uint256\"},{\"name\":\"initial_price\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"deploy_gauge\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_fee_receiver\",\"inputs\":[{\"name\":\"_fee_receiver\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_pool_implementation\",\"inputs\":[{\"name\":\"_pool_implementation\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_token_implementation\",\"inputs\":[{\"name\":\"_token_implementation\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_gauge_implementation\",\"inputs\":[{\"name\":\"_gauge_implementation\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"commit_transfer_ownership\",\"inputs\":[{\"name\":\"_addr\",\"type\":\"address\"}],\"outputs\":[]},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"accept_transfer_ownership\",\"inputs\":[],\"outputs\":[]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"find_pool_for_coins\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"find_pool_for_coins\",\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"i\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_coins\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address[2]\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_decimals\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_balances\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256[2]\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_coin_indices\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"},{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_gauge\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_eth_index\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_token\",\"inputs\":[{\"name\":\"_pool\",\"type\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_admin\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"fee_receiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"pool_implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"token_implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"gauge_implementation\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"pool_count\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"pool_list\",\"inputs\":[{\"name\":\"arg0\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}]}]",
}

// CurvefactoryABI is the input ABI used to generate the binding from.
// Deprecated: Use CurvefactoryMetaData.ABI instead.
var CurvefactoryABI = CurvefactoryMetaData.ABI

// Curvefactory is an auto generated Go binding around an Ethereum contract.
type Curvefactory struct {
	CurvefactoryCaller     // Read-only binding to the contract
	CurvefactoryTransactor // Write-only binding to the contract
	CurvefactoryFilterer   // Log filterer for contract events
}

// CurvefactoryCaller is an auto generated read-only Go binding around an Ethereum contract.
type CurvefactoryCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurvefactoryTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CurvefactoryTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurvefactoryFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CurvefactoryFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurvefactorySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CurvefactorySession struct {
	Contract     *Curvefactory     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CurvefactoryCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CurvefactoryCallerSession struct {
	Contract *CurvefactoryCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// CurvefactoryTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CurvefactoryTransactorSession struct {
	Contract     *CurvefactoryTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// CurvefactoryRaw is an auto generated low-level Go binding around an Ethereum contract.
type CurvefactoryRaw struct {
	Contract *Curvefactory // Generic contract binding to access the raw methods on
}

// CurvefactoryCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CurvefactoryCallerRaw struct {
	Contract *CurvefactoryCaller // Generic read-only contract binding to access the raw methods on
}

// CurvefactoryTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CurvefactoryTransactorRaw struct {
	Contract *CurvefactoryTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCurvefactory creates a new instance of Curvefactory, bound to a specific deployed contract.
func NewCurvefactory(address common.Address, backend bind.ContractBackend) (*Curvefactory, error) {
	contract, err := bindCurvefactory(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Curvefactory{CurvefactoryCaller: CurvefactoryCaller{contract: contract}, CurvefactoryTransactor: CurvefactoryTransactor{contract: contract}, CurvefactoryFilterer: CurvefactoryFilterer{contract: contract}}, nil
}

// NewCurvefactoryCaller creates a new read-only instance of Curvefactory, bound to a specific deployed contract.
func NewCurvefactoryCaller(address common.Address, caller bind.ContractCaller) (*CurvefactoryCaller, error) {
	contract, err := bindCurvefactory(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CurvefactoryCaller{contract: contract}, nil
}

// NewCurvefactoryTransactor creates a new write-only instance of Curvefactory, bound to a specific deployed contract.
func NewCurvefactoryTransactor(address common.Address, transactor bind.ContractTransactor) (*CurvefactoryTransactor, error) {
	contract, err := bindCurvefactory(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CurvefactoryTransactor{contract: contract}, nil
}

// NewCurvefactoryFilterer creates a new log filterer instance of Curvefactory, bound to a specific deployed contract.
func NewCurvefactoryFilterer(address common.Address, filterer bind.ContractFilterer) (*CurvefactoryFilterer, error) {
	contract, err := bindCurvefactory(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CurvefactoryFilterer{contract: contract}, nil
}

// bindCurvefactory binds a generic wrapper to an already deployed contract.
func bindCurvefactory(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CurvefactoryMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Curvefactory *CurvefactoryRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Curvefactory.Contract.CurvefactoryCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Curvefactory *CurvefactoryRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvefactory.Contract.CurvefactoryTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Curvefactory *CurvefactoryRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Curvefactory.Contract.CurvefactoryTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Curvefactory *CurvefactoryCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Curvefactory.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Curvefactory *CurvefactoryTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvefactory.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Curvefactory *CurvefactoryTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Curvefactory.Contract.contract.Transact(opts, method, params...)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Curvefactory *CurvefactoryCaller) Admin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Curvefactory *CurvefactorySession) Admin() (common.Address, error) {
	return _Curvefactory.Contract.Admin(&_Curvefactory.CallOpts)
}

// Admin is a free data retrieval call binding the contract method 0xf851a440.
//
// Solidity: function admin() view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) Admin() (common.Address, error) {
	return _Curvefactory.Contract.Admin(&_Curvefactory.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xcab4d3db.
//
// Solidity: function fee_receiver() view returns(address)
func (_Curvefactory *CurvefactoryCaller) FeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "fee_receiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeReceiver is a free data retrieval call binding the contract method 0xcab4d3db.
//
// Solidity: function fee_receiver() view returns(address)
func (_Curvefactory *CurvefactorySession) FeeReceiver() (common.Address, error) {
	return _Curvefactory.Contract.FeeReceiver(&_Curvefactory.CallOpts)
}

// FeeReceiver is a free data retrieval call binding the contract method 0xcab4d3db.
//
// Solidity: function fee_receiver() view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) FeeReceiver() (common.Address, error) {
	return _Curvefactory.Contract.FeeReceiver(&_Curvefactory.CallOpts)
}

// FindPoolForCoins is a free data retrieval call binding the contract method 0xa87df06c.
//
// Solidity: function find_pool_for_coins(address _from, address _to) view returns(address)
func (_Curvefactory *CurvefactoryCaller) FindPoolForCoins(opts *bind.CallOpts, _from common.Address, _to common.Address) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "find_pool_for_coins", _from, _to)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FindPoolForCoins is a free data retrieval call binding the contract method 0xa87df06c.
//
// Solidity: function find_pool_for_coins(address _from, address _to) view returns(address)
func (_Curvefactory *CurvefactorySession) FindPoolForCoins(_from common.Address, _to common.Address) (common.Address, error) {
	return _Curvefactory.Contract.FindPoolForCoins(&_Curvefactory.CallOpts, _from, _to)
}

// FindPoolForCoins is a free data retrieval call binding the contract method 0xa87df06c.
//
// Solidity: function find_pool_for_coins(address _from, address _to) view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) FindPoolForCoins(_from common.Address, _to common.Address) (common.Address, error) {
	return _Curvefactory.Contract.FindPoolForCoins(&_Curvefactory.CallOpts, _from, _to)
}

// FindPoolForCoins0 is a free data retrieval call binding the contract method 0x6982eb0b.
//
// Solidity: function find_pool_for_coins(address _from, address _to, uint256 i) view returns(address)
func (_Curvefactory *CurvefactoryCaller) FindPoolForCoins0(opts *bind.CallOpts, _from common.Address, _to common.Address, i *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "find_pool_for_coins0", _from, _to, i)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FindPoolForCoins0 is a free data retrieval call binding the contract method 0x6982eb0b.
//
// Solidity: function find_pool_for_coins(address _from, address _to, uint256 i) view returns(address)
func (_Curvefactory *CurvefactorySession) FindPoolForCoins0(_from common.Address, _to common.Address, i *big.Int) (common.Address, error) {
	return _Curvefactory.Contract.FindPoolForCoins0(&_Curvefactory.CallOpts, _from, _to, i)
}

// FindPoolForCoins0 is a free data retrieval call binding the contract method 0x6982eb0b.
//
// Solidity: function find_pool_for_coins(address _from, address _to, uint256 i) view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) FindPoolForCoins0(_from common.Address, _to common.Address, i *big.Int) (common.Address, error) {
	return _Curvefactory.Contract.FindPoolForCoins0(&_Curvefactory.CallOpts, _from, _to, i)
}

// FutureAdmin is a free data retrieval call binding the contract method 0x17f7182a.
//
// Solidity: function future_admin() view returns(address)
func (_Curvefactory *CurvefactoryCaller) FutureAdmin(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "future_admin")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FutureAdmin is a free data retrieval call binding the contract method 0x17f7182a.
//
// Solidity: function future_admin() view returns(address)
func (_Curvefactory *CurvefactorySession) FutureAdmin() (common.Address, error) {
	return _Curvefactory.Contract.FutureAdmin(&_Curvefactory.CallOpts)
}

// FutureAdmin is a free data retrieval call binding the contract method 0x17f7182a.
//
// Solidity: function future_admin() view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) FutureAdmin() (common.Address, error) {
	return _Curvefactory.Contract.FutureAdmin(&_Curvefactory.CallOpts)
}

// GaugeImplementation is a free data retrieval call binding the contract method 0x8df24207.
//
// Solidity: function gauge_implementation() view returns(address)
func (_Curvefactory *CurvefactoryCaller) GaugeImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "gauge_implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GaugeImplementation is a free data retrieval call binding the contract method 0x8df24207.
//
// Solidity: function gauge_implementation() view returns(address)
func (_Curvefactory *CurvefactorySession) GaugeImplementation() (common.Address, error) {
	return _Curvefactory.Contract.GaugeImplementation(&_Curvefactory.CallOpts)
}

// GaugeImplementation is a free data retrieval call binding the contract method 0x8df24207.
//
// Solidity: function gauge_implementation() view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) GaugeImplementation() (common.Address, error) {
	return _Curvefactory.Contract.GaugeImplementation(&_Curvefactory.CallOpts)
}

// GetBalances is a free data retrieval call binding the contract method 0x92e3cc2d.
//
// Solidity: function get_balances(address _pool) view returns(uint256[2])
func (_Curvefactory *CurvefactoryCaller) GetBalances(opts *bind.CallOpts, _pool common.Address) ([2]*big.Int, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_balances", _pool)

	if err != nil {
		return *new([2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([2]*big.Int)).(*[2]*big.Int)

	return out0, err

}

// GetBalances is a free data retrieval call binding the contract method 0x92e3cc2d.
//
// Solidity: function get_balances(address _pool) view returns(uint256[2])
func (_Curvefactory *CurvefactorySession) GetBalances(_pool common.Address) ([2]*big.Int, error) {
	return _Curvefactory.Contract.GetBalances(&_Curvefactory.CallOpts, _pool)
}

// GetBalances is a free data retrieval call binding the contract method 0x92e3cc2d.
//
// Solidity: function get_balances(address _pool) view returns(uint256[2])
func (_Curvefactory *CurvefactoryCallerSession) GetBalances(_pool common.Address) ([2]*big.Int, error) {
	return _Curvefactory.Contract.GetBalances(&_Curvefactory.CallOpts, _pool)
}

// GetCoinIndices is a free data retrieval call binding the contract method 0xeb85226d.
//
// Solidity: function get_coin_indices(address _pool, address _from, address _to) view returns(uint256, uint256)
func (_Curvefactory *CurvefactoryCaller) GetCoinIndices(opts *bind.CallOpts, _pool common.Address, _from common.Address, _to common.Address) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_coin_indices", _pool, _from, _to)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetCoinIndices is a free data retrieval call binding the contract method 0xeb85226d.
//
// Solidity: function get_coin_indices(address _pool, address _from, address _to) view returns(uint256, uint256)
func (_Curvefactory *CurvefactorySession) GetCoinIndices(_pool common.Address, _from common.Address, _to common.Address) (*big.Int, *big.Int, error) {
	return _Curvefactory.Contract.GetCoinIndices(&_Curvefactory.CallOpts, _pool, _from, _to)
}

// GetCoinIndices is a free data retrieval call binding the contract method 0xeb85226d.
//
// Solidity: function get_coin_indices(address _pool, address _from, address _to) view returns(uint256, uint256)
func (_Curvefactory *CurvefactoryCallerSession) GetCoinIndices(_pool common.Address, _from common.Address, _to common.Address) (*big.Int, *big.Int, error) {
	return _Curvefactory.Contract.GetCoinIndices(&_Curvefactory.CallOpts, _pool, _from, _to)
}

// GetCoins is a free data retrieval call binding the contract method 0x9ac90d3d.
//
// Solidity: function get_coins(address _pool) view returns(address[2])
func (_Curvefactory *CurvefactoryCaller) GetCoins(opts *bind.CallOpts, _pool common.Address) ([2]common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_coins", _pool)

	if err != nil {
		return *new([2]common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new([2]common.Address)).(*[2]common.Address)

	return out0, err

}

// GetCoins is a free data retrieval call binding the contract method 0x9ac90d3d.
//
// Solidity: function get_coins(address _pool) view returns(address[2])
func (_Curvefactory *CurvefactorySession) GetCoins(_pool common.Address) ([2]common.Address, error) {
	return _Curvefactory.Contract.GetCoins(&_Curvefactory.CallOpts, _pool)
}

// GetCoins is a free data retrieval call binding the contract method 0x9ac90d3d.
//
// Solidity: function get_coins(address _pool) view returns(address[2])
func (_Curvefactory *CurvefactoryCallerSession) GetCoins(_pool common.Address) ([2]common.Address, error) {
	return _Curvefactory.Contract.GetCoins(&_Curvefactory.CallOpts, _pool)
}

// GetDecimals is a free data retrieval call binding the contract method 0x52b51555.
//
// Solidity: function get_decimals(address _pool) view returns(uint256[2])
func (_Curvefactory *CurvefactoryCaller) GetDecimals(opts *bind.CallOpts, _pool common.Address) ([2]*big.Int, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_decimals", _pool)

	if err != nil {
		return *new([2]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([2]*big.Int)).(*[2]*big.Int)

	return out0, err

}

// GetDecimals is a free data retrieval call binding the contract method 0x52b51555.
//
// Solidity: function get_decimals(address _pool) view returns(uint256[2])
func (_Curvefactory *CurvefactorySession) GetDecimals(_pool common.Address) ([2]*big.Int, error) {
	return _Curvefactory.Contract.GetDecimals(&_Curvefactory.CallOpts, _pool)
}

// GetDecimals is a free data retrieval call binding the contract method 0x52b51555.
//
// Solidity: function get_decimals(address _pool) view returns(uint256[2])
func (_Curvefactory *CurvefactoryCallerSession) GetDecimals(_pool common.Address) ([2]*big.Int, error) {
	return _Curvefactory.Contract.GetDecimals(&_Curvefactory.CallOpts, _pool)
}

// GetEthIndex is a free data retrieval call binding the contract method 0xccb15605.
//
// Solidity: function get_eth_index(address _pool) view returns(uint256)
func (_Curvefactory *CurvefactoryCaller) GetEthIndex(opts *bind.CallOpts, _pool common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_eth_index", _pool)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetEthIndex is a free data retrieval call binding the contract method 0xccb15605.
//
// Solidity: function get_eth_index(address _pool) view returns(uint256)
func (_Curvefactory *CurvefactorySession) GetEthIndex(_pool common.Address) (*big.Int, error) {
	return _Curvefactory.Contract.GetEthIndex(&_Curvefactory.CallOpts, _pool)
}

// GetEthIndex is a free data retrieval call binding the contract method 0xccb15605.
//
// Solidity: function get_eth_index(address _pool) view returns(uint256)
func (_Curvefactory *CurvefactoryCallerSession) GetEthIndex(_pool common.Address) (*big.Int, error) {
	return _Curvefactory.Contract.GetEthIndex(&_Curvefactory.CallOpts, _pool)
}

// GetGauge is a free data retrieval call binding the contract method 0xdaf297b9.
//
// Solidity: function get_gauge(address _pool) view returns(address)
func (_Curvefactory *CurvefactoryCaller) GetGauge(opts *bind.CallOpts, _pool common.Address) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_gauge", _pool)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetGauge is a free data retrieval call binding the contract method 0xdaf297b9.
//
// Solidity: function get_gauge(address _pool) view returns(address)
func (_Curvefactory *CurvefactorySession) GetGauge(_pool common.Address) (common.Address, error) {
	return _Curvefactory.Contract.GetGauge(&_Curvefactory.CallOpts, _pool)
}

// GetGauge is a free data retrieval call binding the contract method 0xdaf297b9.
//
// Solidity: function get_gauge(address _pool) view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) GetGauge(_pool common.Address) (common.Address, error) {
	return _Curvefactory.Contract.GetGauge(&_Curvefactory.CallOpts, _pool)
}

// GetToken is a free data retrieval call binding the contract method 0x977d9122.
//
// Solidity: function get_token(address _pool) view returns(address)
func (_Curvefactory *CurvefactoryCaller) GetToken(opts *bind.CallOpts, _pool common.Address) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "get_token", _pool)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetToken is a free data retrieval call binding the contract method 0x977d9122.
//
// Solidity: function get_token(address _pool) view returns(address)
func (_Curvefactory *CurvefactorySession) GetToken(_pool common.Address) (common.Address, error) {
	return _Curvefactory.Contract.GetToken(&_Curvefactory.CallOpts, _pool)
}

// GetToken is a free data retrieval call binding the contract method 0x977d9122.
//
// Solidity: function get_token(address _pool) view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) GetToken(_pool common.Address) (common.Address, error) {
	return _Curvefactory.Contract.GetToken(&_Curvefactory.CallOpts, _pool)
}

// PoolCount is a free data retrieval call binding the contract method 0x956aae3a.
//
// Solidity: function pool_count() view returns(uint256)
func (_Curvefactory *CurvefactoryCaller) PoolCount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "pool_count")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolCount is a free data retrieval call binding the contract method 0x956aae3a.
//
// Solidity: function pool_count() view returns(uint256)
func (_Curvefactory *CurvefactorySession) PoolCount() (*big.Int, error) {
	return _Curvefactory.Contract.PoolCount(&_Curvefactory.CallOpts)
}

// PoolCount is a free data retrieval call binding the contract method 0x956aae3a.
//
// Solidity: function pool_count() view returns(uint256)
func (_Curvefactory *CurvefactoryCallerSession) PoolCount() (*big.Int, error) {
	return _Curvefactory.Contract.PoolCount(&_Curvefactory.CallOpts)
}

// PoolImplementation is a free data retrieval call binding the contract method 0x2489a2c3.
//
// Solidity: function pool_implementation() view returns(address)
func (_Curvefactory *CurvefactoryCaller) PoolImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "pool_implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolImplementation is a free data retrieval call binding the contract method 0x2489a2c3.
//
// Solidity: function pool_implementation() view returns(address)
func (_Curvefactory *CurvefactorySession) PoolImplementation() (common.Address, error) {
	return _Curvefactory.Contract.PoolImplementation(&_Curvefactory.CallOpts)
}

// PoolImplementation is a free data retrieval call binding the contract method 0x2489a2c3.
//
// Solidity: function pool_implementation() view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) PoolImplementation() (common.Address, error) {
	return _Curvefactory.Contract.PoolImplementation(&_Curvefactory.CallOpts)
}

// PoolList is a free data retrieval call binding the contract method 0x3a1d5d8e.
//
// Solidity: function pool_list(uint256 arg0) view returns(address)
func (_Curvefactory *CurvefactoryCaller) PoolList(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "pool_list", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolList is a free data retrieval call binding the contract method 0x3a1d5d8e.
//
// Solidity: function pool_list(uint256 arg0) view returns(address)
func (_Curvefactory *CurvefactorySession) PoolList(arg0 *big.Int) (common.Address, error) {
	return _Curvefactory.Contract.PoolList(&_Curvefactory.CallOpts, arg0)
}

// PoolList is a free data retrieval call binding the contract method 0x3a1d5d8e.
//
// Solidity: function pool_list(uint256 arg0) view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) PoolList(arg0 *big.Int) (common.Address, error) {
	return _Curvefactory.Contract.PoolList(&_Curvefactory.CallOpts, arg0)
}

// TokenImplementation is a free data retrieval call binding the contract method 0x35214d81.
//
// Solidity: function token_implementation() view returns(address)
func (_Curvefactory *CurvefactoryCaller) TokenImplementation(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvefactory.contract.Call(opts, &out, "token_implementation")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TokenImplementation is a free data retrieval call binding the contract method 0x35214d81.
//
// Solidity: function token_implementation() view returns(address)
func (_Curvefactory *CurvefactorySession) TokenImplementation() (common.Address, error) {
	return _Curvefactory.Contract.TokenImplementation(&_Curvefactory.CallOpts)
}

// TokenImplementation is a free data retrieval call binding the contract method 0x35214d81.
//
// Solidity: function token_implementation() view returns(address)
func (_Curvefactory *CurvefactoryCallerSession) TokenImplementation() (common.Address, error) {
	return _Curvefactory.Contract.TokenImplementation(&_Curvefactory.CallOpts)
}

// AcceptTransferOwnership is a paid mutator transaction binding the contract method 0xe5ea47b8.
//
// Solidity: function accept_transfer_ownership() returns()
func (_Curvefactory *CurvefactoryTransactor) AcceptTransferOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "accept_transfer_ownership")
}

// AcceptTransferOwnership is a paid mutator transaction binding the contract method 0xe5ea47b8.
//
// Solidity: function accept_transfer_ownership() returns()
func (_Curvefactory *CurvefactorySession) AcceptTransferOwnership() (*types.Transaction, error) {
	return _Curvefactory.Contract.AcceptTransferOwnership(&_Curvefactory.TransactOpts)
}

// AcceptTransferOwnership is a paid mutator transaction binding the contract method 0xe5ea47b8.
//
// Solidity: function accept_transfer_ownership() returns()
func (_Curvefactory *CurvefactoryTransactorSession) AcceptTransferOwnership() (*types.Transaction, error) {
	return _Curvefactory.Contract.AcceptTransferOwnership(&_Curvefactory.TransactOpts)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _addr) returns()
func (_Curvefactory *CurvefactoryTransactor) CommitTransferOwnership(opts *bind.TransactOpts, _addr common.Address) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "commit_transfer_ownership", _addr)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _addr) returns()
func (_Curvefactory *CurvefactorySession) CommitTransferOwnership(_addr common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.CommitTransferOwnership(&_Curvefactory.TransactOpts, _addr)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _addr) returns()
func (_Curvefactory *CurvefactoryTransactorSession) CommitTransferOwnership(_addr common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.CommitTransferOwnership(&_Curvefactory.TransactOpts, _addr)
}

// DeployGauge is a paid mutator transaction binding the contract method 0x96bebb34.
//
// Solidity: function deploy_gauge(address _pool) returns(address)
func (_Curvefactory *CurvefactoryTransactor) DeployGauge(opts *bind.TransactOpts, _pool common.Address) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "deploy_gauge", _pool)
}

// DeployGauge is a paid mutator transaction binding the contract method 0x96bebb34.
//
// Solidity: function deploy_gauge(address _pool) returns(address)
func (_Curvefactory *CurvefactorySession) DeployGauge(_pool common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.DeployGauge(&_Curvefactory.TransactOpts, _pool)
}

// DeployGauge is a paid mutator transaction binding the contract method 0x96bebb34.
//
// Solidity: function deploy_gauge(address _pool) returns(address)
func (_Curvefactory *CurvefactoryTransactorSession) DeployGauge(_pool common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.DeployGauge(&_Curvefactory.TransactOpts, _pool)
}

// DeployPool is a paid mutator transaction binding the contract method 0xc955fa04.
//
// Solidity: function deploy_pool(string _name, string _symbol, address[2] _coins, uint256 A, uint256 gamma, uint256 mid_fee, uint256 out_fee, uint256 allowed_extra_profit, uint256 fee_gamma, uint256 adjustment_step, uint256 admin_fee, uint256 ma_half_time, uint256 initial_price) returns(address)
func (_Curvefactory *CurvefactoryTransactor) DeployPool(opts *bind.TransactOpts, _name string, _symbol string, _coins [2]common.Address, A *big.Int, gamma *big.Int, mid_fee *big.Int, out_fee *big.Int, allowed_extra_profit *big.Int, fee_gamma *big.Int, adjustment_step *big.Int, admin_fee *big.Int, ma_half_time *big.Int, initial_price *big.Int) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "deploy_pool", _name, _symbol, _coins, A, gamma, mid_fee, out_fee, allowed_extra_profit, fee_gamma, adjustment_step, admin_fee, ma_half_time, initial_price)
}

// DeployPool is a paid mutator transaction binding the contract method 0xc955fa04.
//
// Solidity: function deploy_pool(string _name, string _symbol, address[2] _coins, uint256 A, uint256 gamma, uint256 mid_fee, uint256 out_fee, uint256 allowed_extra_profit, uint256 fee_gamma, uint256 adjustment_step, uint256 admin_fee, uint256 ma_half_time, uint256 initial_price) returns(address)
func (_Curvefactory *CurvefactorySession) DeployPool(_name string, _symbol string, _coins [2]common.Address, A *big.Int, gamma *big.Int, mid_fee *big.Int, out_fee *big.Int, allowed_extra_profit *big.Int, fee_gamma *big.Int, adjustment_step *big.Int, admin_fee *big.Int, ma_half_time *big.Int, initial_price *big.Int) (*types.Transaction, error) {
	return _Curvefactory.Contract.DeployPool(&_Curvefactory.TransactOpts, _name, _symbol, _coins, A, gamma, mid_fee, out_fee, allowed_extra_profit, fee_gamma, adjustment_step, admin_fee, ma_half_time, initial_price)
}

// DeployPool is a paid mutator transaction binding the contract method 0xc955fa04.
//
// Solidity: function deploy_pool(string _name, string _symbol, address[2] _coins, uint256 A, uint256 gamma, uint256 mid_fee, uint256 out_fee, uint256 allowed_extra_profit, uint256 fee_gamma, uint256 adjustment_step, uint256 admin_fee, uint256 ma_half_time, uint256 initial_price) returns(address)
func (_Curvefactory *CurvefactoryTransactorSession) DeployPool(_name string, _symbol string, _coins [2]common.Address, A *big.Int, gamma *big.Int, mid_fee *big.Int, out_fee *big.Int, allowed_extra_profit *big.Int, fee_gamma *big.Int, adjustment_step *big.Int, admin_fee *big.Int, ma_half_time *big.Int, initial_price *big.Int) (*types.Transaction, error) {
	return _Curvefactory.Contract.DeployPool(&_Curvefactory.TransactOpts, _name, _symbol, _coins, A, gamma, mid_fee, out_fee, allowed_extra_profit, fee_gamma, adjustment_step, admin_fee, ma_half_time, initial_price)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xe41ab771.
//
// Solidity: function set_fee_receiver(address _fee_receiver) returns()
func (_Curvefactory *CurvefactoryTransactor) SetFeeReceiver(opts *bind.TransactOpts, _fee_receiver common.Address) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "set_fee_receiver", _fee_receiver)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xe41ab771.
//
// Solidity: function set_fee_receiver(address _fee_receiver) returns()
func (_Curvefactory *CurvefactorySession) SetFeeReceiver(_fee_receiver common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetFeeReceiver(&_Curvefactory.TransactOpts, _fee_receiver)
}

// SetFeeReceiver is a paid mutator transaction binding the contract method 0xe41ab771.
//
// Solidity: function set_fee_receiver(address _fee_receiver) returns()
func (_Curvefactory *CurvefactoryTransactorSession) SetFeeReceiver(_fee_receiver common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetFeeReceiver(&_Curvefactory.TransactOpts, _fee_receiver)
}

// SetGaugeImplementation is a paid mutator transaction binding the contract method 0x8f03182c.
//
// Solidity: function set_gauge_implementation(address _gauge_implementation) returns()
func (_Curvefactory *CurvefactoryTransactor) SetGaugeImplementation(opts *bind.TransactOpts, _gauge_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "set_gauge_implementation", _gauge_implementation)
}

// SetGaugeImplementation is a paid mutator transaction binding the contract method 0x8f03182c.
//
// Solidity: function set_gauge_implementation(address _gauge_implementation) returns()
func (_Curvefactory *CurvefactorySession) SetGaugeImplementation(_gauge_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetGaugeImplementation(&_Curvefactory.TransactOpts, _gauge_implementation)
}

// SetGaugeImplementation is a paid mutator transaction binding the contract method 0x8f03182c.
//
// Solidity: function set_gauge_implementation(address _gauge_implementation) returns()
func (_Curvefactory *CurvefactoryTransactorSession) SetGaugeImplementation(_gauge_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetGaugeImplementation(&_Curvefactory.TransactOpts, _gauge_implementation)
}

// SetPoolImplementation is a paid mutator transaction binding the contract method 0x9ed796d0.
//
// Solidity: function set_pool_implementation(address _pool_implementation) returns()
func (_Curvefactory *CurvefactoryTransactor) SetPoolImplementation(opts *bind.TransactOpts, _pool_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "set_pool_implementation", _pool_implementation)
}

// SetPoolImplementation is a paid mutator transaction binding the contract method 0x9ed796d0.
//
// Solidity: function set_pool_implementation(address _pool_implementation) returns()
func (_Curvefactory *CurvefactorySession) SetPoolImplementation(_pool_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetPoolImplementation(&_Curvefactory.TransactOpts, _pool_implementation)
}

// SetPoolImplementation is a paid mutator transaction binding the contract method 0x9ed796d0.
//
// Solidity: function set_pool_implementation(address _pool_implementation) returns()
func (_Curvefactory *CurvefactoryTransactorSession) SetPoolImplementation(_pool_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetPoolImplementation(&_Curvefactory.TransactOpts, _pool_implementation)
}

// SetTokenImplementation is a paid mutator transaction binding the contract method 0x653023c2.
//
// Solidity: function set_token_implementation(address _token_implementation) returns()
func (_Curvefactory *CurvefactoryTransactor) SetTokenImplementation(opts *bind.TransactOpts, _token_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.contract.Transact(opts, "set_token_implementation", _token_implementation)
}

// SetTokenImplementation is a paid mutator transaction binding the contract method 0x653023c2.
//
// Solidity: function set_token_implementation(address _token_implementation) returns()
func (_Curvefactory *CurvefactorySession) SetTokenImplementation(_token_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetTokenImplementation(&_Curvefactory.TransactOpts, _token_implementation)
}

// SetTokenImplementation is a paid mutator transaction binding the contract method 0x653023c2.
//
// Solidity: function set_token_implementation(address _token_implementation) returns()
func (_Curvefactory *CurvefactoryTransactorSession) SetTokenImplementation(_token_implementation common.Address) (*types.Transaction, error) {
	return _Curvefactory.Contract.SetTokenImplementation(&_Curvefactory.TransactOpts, _token_implementation)
}

// CurvefactoryCryptoPoolDeployedIterator is returned from FilterCryptoPoolDeployed and is used to iterate over the raw logs and unpacked data for CryptoPoolDeployed events raised by the Curvefactory contract.
type CurvefactoryCryptoPoolDeployedIterator struct {
	Event *CurvefactoryCryptoPoolDeployed // Event containing the contract specifics and raw log

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
func (it *CurvefactoryCryptoPoolDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryCryptoPoolDeployed)
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
		it.Event = new(CurvefactoryCryptoPoolDeployed)
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
func (it *CurvefactoryCryptoPoolDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryCryptoPoolDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryCryptoPoolDeployed represents a CryptoPoolDeployed event raised by the Curvefactory contract.
type CurvefactoryCryptoPoolDeployed struct {
	Token              common.Address
	Coins              [2]common.Address
	A                  *big.Int
	Gamma              *big.Int
	MidFee             *big.Int
	OutFee             *big.Int
	AllowedExtraProfit *big.Int
	FeeGamma           *big.Int
	AdjustmentStep     *big.Int
	AdminFee           *big.Int
	MaHalfTime         *big.Int
	InitialPrice       *big.Int
	Deployer           common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterCryptoPoolDeployed is a free log retrieval operation binding the contract event 0x0394cb40d7dbe28dad1d4ee890bdd35bbb0d89e17924a80a542535e83d54ba14.
//
// Solidity: event CryptoPoolDeployed(address token, address[2] coins, uint256 A, uint256 gamma, uint256 mid_fee, uint256 out_fee, uint256 allowed_extra_profit, uint256 fee_gamma, uint256 adjustment_step, uint256 admin_fee, uint256 ma_half_time, uint256 initial_price, address deployer)
func (_Curvefactory *CurvefactoryFilterer) FilterCryptoPoolDeployed(opts *bind.FilterOpts) (*CurvefactoryCryptoPoolDeployedIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "CryptoPoolDeployed")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryCryptoPoolDeployedIterator{contract: _Curvefactory.contract, event: "CryptoPoolDeployed", logs: logs, sub: sub}, nil
}

// WatchCryptoPoolDeployed is a free log subscription operation binding the contract event 0x0394cb40d7dbe28dad1d4ee890bdd35bbb0d89e17924a80a542535e83d54ba14.
//
// Solidity: event CryptoPoolDeployed(address token, address[2] coins, uint256 A, uint256 gamma, uint256 mid_fee, uint256 out_fee, uint256 allowed_extra_profit, uint256 fee_gamma, uint256 adjustment_step, uint256 admin_fee, uint256 ma_half_time, uint256 initial_price, address deployer)
func (_Curvefactory *CurvefactoryFilterer) WatchCryptoPoolDeployed(opts *bind.WatchOpts, sink chan<- *CurvefactoryCryptoPoolDeployed) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "CryptoPoolDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryCryptoPoolDeployed)
				if err := _Curvefactory.contract.UnpackLog(event, "CryptoPoolDeployed", log); err != nil {
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

// ParseCryptoPoolDeployed is a log parse operation binding the contract event 0x0394cb40d7dbe28dad1d4ee890bdd35bbb0d89e17924a80a542535e83d54ba14.
//
// Solidity: event CryptoPoolDeployed(address token, address[2] coins, uint256 A, uint256 gamma, uint256 mid_fee, uint256 out_fee, uint256 allowed_extra_profit, uint256 fee_gamma, uint256 adjustment_step, uint256 admin_fee, uint256 ma_half_time, uint256 initial_price, address deployer)
func (_Curvefactory *CurvefactoryFilterer) ParseCryptoPoolDeployed(log types.Log) (*CurvefactoryCryptoPoolDeployed, error) {
	event := new(CurvefactoryCryptoPoolDeployed)
	if err := _Curvefactory.contract.UnpackLog(event, "CryptoPoolDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvefactoryLiquidityGaugeDeployedIterator is returned from FilterLiquidityGaugeDeployed and is used to iterate over the raw logs and unpacked data for LiquidityGaugeDeployed events raised by the Curvefactory contract.
type CurvefactoryLiquidityGaugeDeployedIterator struct {
	Event *CurvefactoryLiquidityGaugeDeployed // Event containing the contract specifics and raw log

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
func (it *CurvefactoryLiquidityGaugeDeployedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryLiquidityGaugeDeployed)
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
		it.Event = new(CurvefactoryLiquidityGaugeDeployed)
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
func (it *CurvefactoryLiquidityGaugeDeployedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryLiquidityGaugeDeployedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryLiquidityGaugeDeployed represents a LiquidityGaugeDeployed event raised by the Curvefactory contract.
type CurvefactoryLiquidityGaugeDeployed struct {
	Pool  common.Address
	Token common.Address
	Gauge common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterLiquidityGaugeDeployed is a free log retrieval operation binding the contract event 0x1d6247eae69b5feb96b30be78552f35de45f61fdb6d6d7e1b08aae159b6226af.
//
// Solidity: event LiquidityGaugeDeployed(address pool, address token, address gauge)
func (_Curvefactory *CurvefactoryFilterer) FilterLiquidityGaugeDeployed(opts *bind.FilterOpts) (*CurvefactoryLiquidityGaugeDeployedIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "LiquidityGaugeDeployed")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryLiquidityGaugeDeployedIterator{contract: _Curvefactory.contract, event: "LiquidityGaugeDeployed", logs: logs, sub: sub}, nil
}

// WatchLiquidityGaugeDeployed is a free log subscription operation binding the contract event 0x1d6247eae69b5feb96b30be78552f35de45f61fdb6d6d7e1b08aae159b6226af.
//
// Solidity: event LiquidityGaugeDeployed(address pool, address token, address gauge)
func (_Curvefactory *CurvefactoryFilterer) WatchLiquidityGaugeDeployed(opts *bind.WatchOpts, sink chan<- *CurvefactoryLiquidityGaugeDeployed) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "LiquidityGaugeDeployed")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryLiquidityGaugeDeployed)
				if err := _Curvefactory.contract.UnpackLog(event, "LiquidityGaugeDeployed", log); err != nil {
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

// ParseLiquidityGaugeDeployed is a log parse operation binding the contract event 0x1d6247eae69b5feb96b30be78552f35de45f61fdb6d6d7e1b08aae159b6226af.
//
// Solidity: event LiquidityGaugeDeployed(address pool, address token, address gauge)
func (_Curvefactory *CurvefactoryFilterer) ParseLiquidityGaugeDeployed(log types.Log) (*CurvefactoryLiquidityGaugeDeployed, error) {
	event := new(CurvefactoryLiquidityGaugeDeployed)
	if err := _Curvefactory.contract.UnpackLog(event, "LiquidityGaugeDeployed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvefactoryTransferOwnershipIterator is returned from FilterTransferOwnership and is used to iterate over the raw logs and unpacked data for TransferOwnership events raised by the Curvefactory contract.
type CurvefactoryTransferOwnershipIterator struct {
	Event *CurvefactoryTransferOwnership // Event containing the contract specifics and raw log

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
func (it *CurvefactoryTransferOwnershipIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryTransferOwnership)
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
		it.Event = new(CurvefactoryTransferOwnership)
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
func (it *CurvefactoryTransferOwnershipIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryTransferOwnershipIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryTransferOwnership represents a TransferOwnership event raised by the Curvefactory contract.
type CurvefactoryTransferOwnership struct {
	OldOwner common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterTransferOwnership is a free log retrieval operation binding the contract event 0x5c486528ec3e3f0ea91181cff8116f02bfa350e03b8b6f12e00765adbb5af85c.
//
// Solidity: event TransferOwnership(address _old_owner, address _new_owner)
func (_Curvefactory *CurvefactoryFilterer) FilterTransferOwnership(opts *bind.FilterOpts) (*CurvefactoryTransferOwnershipIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "TransferOwnership")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryTransferOwnershipIterator{contract: _Curvefactory.contract, event: "TransferOwnership", logs: logs, sub: sub}, nil
}

// WatchTransferOwnership is a free log subscription operation binding the contract event 0x5c486528ec3e3f0ea91181cff8116f02bfa350e03b8b6f12e00765adbb5af85c.
//
// Solidity: event TransferOwnership(address _old_owner, address _new_owner)
func (_Curvefactory *CurvefactoryFilterer) WatchTransferOwnership(opts *bind.WatchOpts, sink chan<- *CurvefactoryTransferOwnership) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "TransferOwnership")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryTransferOwnership)
				if err := _Curvefactory.contract.UnpackLog(event, "TransferOwnership", log); err != nil {
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

// ParseTransferOwnership is a log parse operation binding the contract event 0x5c486528ec3e3f0ea91181cff8116f02bfa350e03b8b6f12e00765adbb5af85c.
//
// Solidity: event TransferOwnership(address _old_owner, address _new_owner)
func (_Curvefactory *CurvefactoryFilterer) ParseTransferOwnership(log types.Log) (*CurvefactoryTransferOwnership, error) {
	event := new(CurvefactoryTransferOwnership)
	if err := _Curvefactory.contract.UnpackLog(event, "TransferOwnership", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvefactoryUpdateFeeReceiverIterator is returned from FilterUpdateFeeReceiver and is used to iterate over the raw logs and unpacked data for UpdateFeeReceiver events raised by the Curvefactory contract.
type CurvefactoryUpdateFeeReceiverIterator struct {
	Event *CurvefactoryUpdateFeeReceiver // Event containing the contract specifics and raw log

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
func (it *CurvefactoryUpdateFeeReceiverIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryUpdateFeeReceiver)
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
		it.Event = new(CurvefactoryUpdateFeeReceiver)
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
func (it *CurvefactoryUpdateFeeReceiverIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryUpdateFeeReceiverIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryUpdateFeeReceiver represents a UpdateFeeReceiver event raised by the Curvefactory contract.
type CurvefactoryUpdateFeeReceiver struct {
	OldFeeReceiver common.Address
	NewFeeReceiver common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpdateFeeReceiver is a free log retrieval operation binding the contract event 0x2861448678f0be67f11bfb5481b3e3b4cfeb3acc6126ad60a05f95bfc6530666.
//
// Solidity: event UpdateFeeReceiver(address _old_fee_receiver, address _new_fee_receiver)
func (_Curvefactory *CurvefactoryFilterer) FilterUpdateFeeReceiver(opts *bind.FilterOpts) (*CurvefactoryUpdateFeeReceiverIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "UpdateFeeReceiver")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryUpdateFeeReceiverIterator{contract: _Curvefactory.contract, event: "UpdateFeeReceiver", logs: logs, sub: sub}, nil
}

// WatchUpdateFeeReceiver is a free log subscription operation binding the contract event 0x2861448678f0be67f11bfb5481b3e3b4cfeb3acc6126ad60a05f95bfc6530666.
//
// Solidity: event UpdateFeeReceiver(address _old_fee_receiver, address _new_fee_receiver)
func (_Curvefactory *CurvefactoryFilterer) WatchUpdateFeeReceiver(opts *bind.WatchOpts, sink chan<- *CurvefactoryUpdateFeeReceiver) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "UpdateFeeReceiver")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryUpdateFeeReceiver)
				if err := _Curvefactory.contract.UnpackLog(event, "UpdateFeeReceiver", log); err != nil {
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

// ParseUpdateFeeReceiver is a log parse operation binding the contract event 0x2861448678f0be67f11bfb5481b3e3b4cfeb3acc6126ad60a05f95bfc6530666.
//
// Solidity: event UpdateFeeReceiver(address _old_fee_receiver, address _new_fee_receiver)
func (_Curvefactory *CurvefactoryFilterer) ParseUpdateFeeReceiver(log types.Log) (*CurvefactoryUpdateFeeReceiver, error) {
	event := new(CurvefactoryUpdateFeeReceiver)
	if err := _Curvefactory.contract.UnpackLog(event, "UpdateFeeReceiver", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvefactoryUpdateGaugeImplementationIterator is returned from FilterUpdateGaugeImplementation and is used to iterate over the raw logs and unpacked data for UpdateGaugeImplementation events raised by the Curvefactory contract.
type CurvefactoryUpdateGaugeImplementationIterator struct {
	Event *CurvefactoryUpdateGaugeImplementation // Event containing the contract specifics and raw log

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
func (it *CurvefactoryUpdateGaugeImplementationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryUpdateGaugeImplementation)
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
		it.Event = new(CurvefactoryUpdateGaugeImplementation)
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
func (it *CurvefactoryUpdateGaugeImplementationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryUpdateGaugeImplementationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryUpdateGaugeImplementation represents a UpdateGaugeImplementation event raised by the Curvefactory contract.
type CurvefactoryUpdateGaugeImplementation struct {
	OldGaugeImplementation common.Address
	NewGaugeImplementation common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterUpdateGaugeImplementation is a free log retrieval operation binding the contract event 0x1fd705f9c77053962a503f2f2f57f0862b4c3af687c25615c13817a86946c359.
//
// Solidity: event UpdateGaugeImplementation(address _old_gauge_implementation, address _new_gauge_implementation)
func (_Curvefactory *CurvefactoryFilterer) FilterUpdateGaugeImplementation(opts *bind.FilterOpts) (*CurvefactoryUpdateGaugeImplementationIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "UpdateGaugeImplementation")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryUpdateGaugeImplementationIterator{contract: _Curvefactory.contract, event: "UpdateGaugeImplementation", logs: logs, sub: sub}, nil
}

// WatchUpdateGaugeImplementation is a free log subscription operation binding the contract event 0x1fd705f9c77053962a503f2f2f57f0862b4c3af687c25615c13817a86946c359.
//
// Solidity: event UpdateGaugeImplementation(address _old_gauge_implementation, address _new_gauge_implementation)
func (_Curvefactory *CurvefactoryFilterer) WatchUpdateGaugeImplementation(opts *bind.WatchOpts, sink chan<- *CurvefactoryUpdateGaugeImplementation) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "UpdateGaugeImplementation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryUpdateGaugeImplementation)
				if err := _Curvefactory.contract.UnpackLog(event, "UpdateGaugeImplementation", log); err != nil {
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

// ParseUpdateGaugeImplementation is a log parse operation binding the contract event 0x1fd705f9c77053962a503f2f2f57f0862b4c3af687c25615c13817a86946c359.
//
// Solidity: event UpdateGaugeImplementation(address _old_gauge_implementation, address _new_gauge_implementation)
func (_Curvefactory *CurvefactoryFilterer) ParseUpdateGaugeImplementation(log types.Log) (*CurvefactoryUpdateGaugeImplementation, error) {
	event := new(CurvefactoryUpdateGaugeImplementation)
	if err := _Curvefactory.contract.UnpackLog(event, "UpdateGaugeImplementation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvefactoryUpdatePoolImplementationIterator is returned from FilterUpdatePoolImplementation and is used to iterate over the raw logs and unpacked data for UpdatePoolImplementation events raised by the Curvefactory contract.
type CurvefactoryUpdatePoolImplementationIterator struct {
	Event *CurvefactoryUpdatePoolImplementation // Event containing the contract specifics and raw log

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
func (it *CurvefactoryUpdatePoolImplementationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryUpdatePoolImplementation)
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
		it.Event = new(CurvefactoryUpdatePoolImplementation)
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
func (it *CurvefactoryUpdatePoolImplementationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryUpdatePoolImplementationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryUpdatePoolImplementation represents a UpdatePoolImplementation event raised by the Curvefactory contract.
type CurvefactoryUpdatePoolImplementation struct {
	OldPoolImplementation common.Address
	NewPoolImplementation common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterUpdatePoolImplementation is a free log retrieval operation binding the contract event 0x0617fd31aa5ab95ec80eefc1eb61a2c477aa419d1d761b4e46f5f077e47852aa.
//
// Solidity: event UpdatePoolImplementation(address _old_pool_implementation, address _new_pool_implementation)
func (_Curvefactory *CurvefactoryFilterer) FilterUpdatePoolImplementation(opts *bind.FilterOpts) (*CurvefactoryUpdatePoolImplementationIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "UpdatePoolImplementation")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryUpdatePoolImplementationIterator{contract: _Curvefactory.contract, event: "UpdatePoolImplementation", logs: logs, sub: sub}, nil
}

// WatchUpdatePoolImplementation is a free log subscription operation binding the contract event 0x0617fd31aa5ab95ec80eefc1eb61a2c477aa419d1d761b4e46f5f077e47852aa.
//
// Solidity: event UpdatePoolImplementation(address _old_pool_implementation, address _new_pool_implementation)
func (_Curvefactory *CurvefactoryFilterer) WatchUpdatePoolImplementation(opts *bind.WatchOpts, sink chan<- *CurvefactoryUpdatePoolImplementation) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "UpdatePoolImplementation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryUpdatePoolImplementation)
				if err := _Curvefactory.contract.UnpackLog(event, "UpdatePoolImplementation", log); err != nil {
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

// ParseUpdatePoolImplementation is a log parse operation binding the contract event 0x0617fd31aa5ab95ec80eefc1eb61a2c477aa419d1d761b4e46f5f077e47852aa.
//
// Solidity: event UpdatePoolImplementation(address _old_pool_implementation, address _new_pool_implementation)
func (_Curvefactory *CurvefactoryFilterer) ParseUpdatePoolImplementation(log types.Log) (*CurvefactoryUpdatePoolImplementation, error) {
	event := new(CurvefactoryUpdatePoolImplementation)
	if err := _Curvefactory.contract.UnpackLog(event, "UpdatePoolImplementation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvefactoryUpdateTokenImplementationIterator is returned from FilterUpdateTokenImplementation and is used to iterate over the raw logs and unpacked data for UpdateTokenImplementation events raised by the Curvefactory contract.
type CurvefactoryUpdateTokenImplementationIterator struct {
	Event *CurvefactoryUpdateTokenImplementation // Event containing the contract specifics and raw log

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
func (it *CurvefactoryUpdateTokenImplementationIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvefactoryUpdateTokenImplementation)
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
		it.Event = new(CurvefactoryUpdateTokenImplementation)
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
func (it *CurvefactoryUpdateTokenImplementationIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvefactoryUpdateTokenImplementationIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvefactoryUpdateTokenImplementation represents a UpdateTokenImplementation event raised by the Curvefactory contract.
type CurvefactoryUpdateTokenImplementation struct {
	OldTokenImplementation common.Address
	NewTokenImplementation common.Address
	Raw                    types.Log // Blockchain specific contextual infos
}

// FilterUpdateTokenImplementation is a free log retrieval operation binding the contract event 0x1cc4f8e20b0cd3e5109eb156cadcfd3a5496ac0794c6085ca02afd7d710dd566.
//
// Solidity: event UpdateTokenImplementation(address _old_token_implementation, address _new_token_implementation)
func (_Curvefactory *CurvefactoryFilterer) FilterUpdateTokenImplementation(opts *bind.FilterOpts) (*CurvefactoryUpdateTokenImplementationIterator, error) {

	logs, sub, err := _Curvefactory.contract.FilterLogs(opts, "UpdateTokenImplementation")
	if err != nil {
		return nil, err
	}
	return &CurvefactoryUpdateTokenImplementationIterator{contract: _Curvefactory.contract, event: "UpdateTokenImplementation", logs: logs, sub: sub}, nil
}

// WatchUpdateTokenImplementation is a free log subscription operation binding the contract event 0x1cc4f8e20b0cd3e5109eb156cadcfd3a5496ac0794c6085ca02afd7d710dd566.
//
// Solidity: event UpdateTokenImplementation(address _old_token_implementation, address _new_token_implementation)
func (_Curvefactory *CurvefactoryFilterer) WatchUpdateTokenImplementation(opts *bind.WatchOpts, sink chan<- *CurvefactoryUpdateTokenImplementation) (event.Subscription, error) {

	logs, sub, err := _Curvefactory.contract.WatchLogs(opts, "UpdateTokenImplementation")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvefactoryUpdateTokenImplementation)
				if err := _Curvefactory.contract.UnpackLog(event, "UpdateTokenImplementation", log); err != nil {
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

// ParseUpdateTokenImplementation is a log parse operation binding the contract event 0x1cc4f8e20b0cd3e5109eb156cadcfd3a5496ac0794c6085ca02afd7d710dd566.
//
// Solidity: event UpdateTokenImplementation(address _old_token_implementation, address _new_token_implementation)
func (_Curvefactory *CurvefactoryFilterer) ParseUpdateTokenImplementation(log types.Log) (*CurvefactoryUpdateTokenImplementation, error) {
	event := new(CurvefactoryUpdateTokenImplementation)
	if err := _Curvefactory.contract.UnpackLog(event, "UpdateTokenImplementation", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
