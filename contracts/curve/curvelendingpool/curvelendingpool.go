// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package curvelendingpool

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

// CurvelendingpoolMetaData contains all meta data concerning the Curvelendingpool contract.
var CurvelendingpoolMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"TokenExchange\",\"inputs\":[{\"type\":\"address\",\"name\":\"buyer\",\"indexed\":true},{\"type\":\"int128\",\"name\":\"sold_id\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"tokens_sold\",\"indexed\":false},{\"type\":\"int128\",\"name\":\"bought_id\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"tokens_bought\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"TokenExchangeUnderlying\",\"inputs\":[{\"type\":\"address\",\"name\":\"buyer\",\"indexed\":true},{\"type\":\"int128\",\"name\":\"sold_id\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"tokens_sold\",\"indexed\":false},{\"type\":\"int128\",\"name\":\"bought_id\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"tokens_bought\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"AddLiquidity\",\"inputs\":[{\"type\":\"address\",\"name\":\"provider\",\"indexed\":true},{\"type\":\"uint256[3]\",\"name\":\"token_amounts\",\"indexed\":false},{\"type\":\"uint256[3]\",\"name\":\"fees\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"invariant\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"token_supply\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RemoveLiquidity\",\"inputs\":[{\"type\":\"address\",\"name\":\"provider\",\"indexed\":true},{\"type\":\"uint256[3]\",\"name\":\"token_amounts\",\"indexed\":false},{\"type\":\"uint256[3]\",\"name\":\"fees\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"token_supply\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RemoveLiquidityOne\",\"inputs\":[{\"type\":\"address\",\"name\":\"provider\",\"indexed\":true},{\"type\":\"uint256\",\"name\":\"token_amount\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"coin_amount\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RemoveLiquidityImbalance\",\"inputs\":[{\"type\":\"address\",\"name\":\"provider\",\"indexed\":true},{\"type\":\"uint256[3]\",\"name\":\"token_amounts\",\"indexed\":false},{\"type\":\"uint256[3]\",\"name\":\"fees\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"invariant\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"token_supply\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"CommitNewAdmin\",\"inputs\":[{\"type\":\"uint256\",\"name\":\"deadline\",\"indexed\":true},{\"type\":\"address\",\"name\":\"admin\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"NewAdmin\",\"inputs\":[{\"type\":\"address\",\"name\":\"admin\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"CommitNewFee\",\"inputs\":[{\"type\":\"uint256\",\"name\":\"deadline\",\"indexed\":true},{\"type\":\"uint256\",\"name\":\"fee\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"admin_fee\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"offpeg_fee_multiplier\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"NewFee\",\"inputs\":[{\"type\":\"uint256\",\"name\":\"fee\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"admin_fee\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"offpeg_fee_multiplier\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RampA\",\"inputs\":[{\"type\":\"uint256\",\"name\":\"old_A\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"new_A\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"initial_time\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"future_time\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"StopRampA\",\"inputs\":[{\"type\":\"uint256\",\"name\":\"A\",\"indexed\":false},{\"type\":\"uint256\",\"name\":\"t\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"outputs\":[],\"inputs\":[{\"type\":\"address[3]\",\"name\":\"_coins\"},{\"type\":\"address[3]\",\"name\":\"_underlying_coins\"},{\"type\":\"address\",\"name\":\"_pool_token\"},{\"type\":\"address\",\"name\":\"_aave_lending_pool\"},{\"type\":\"uint256\",\"name\":\"_A\"},{\"type\":\"uint256\",\"name\":\"_fee\"},{\"type\":\"uint256\",\"name\":\"_admin_fee\"},{\"type\":\"uint256\",\"name\":\"_offpeg_fee_multiplier\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"name\":\"A\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":5199},{\"name\":\"A_precise\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":5161},{\"name\":\"dynamic_fee\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"int128\",\"name\":\"j\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":10278},{\"name\":\"balances\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"i\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2731},{\"name\":\"get_virtual_price\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2680120},{\"name\":\"calc_token_amount\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256[3]\",\"name\":\"_amounts\"},{\"type\":\"bool\",\"name\":\"is_deposit\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":5346581},{\"name\":\"add_liquidity\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256[3]\",\"name\":\"_amounts\"},{\"type\":\"uint256\",\"name\":\"_min_mint_amount\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"add_liquidity\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256[3]\",\"name\":\"_amounts\"},{\"type\":\"uint256\",\"name\":\"_min_mint_amount\"},{\"type\":\"bool\",\"name\":\"_use_underlying\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"get_dy\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"int128\",\"name\":\"j\"},{\"type\":\"uint256\",\"name\":\"dx\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":6239547},{\"name\":\"get_dy_underlying\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"int128\",\"name\":\"j\"},{\"type\":\"uint256\",\"name\":\"dx\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":6239577},{\"name\":\"exchange\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"int128\",\"name\":\"j\"},{\"type\":\"uint256\",\"name\":\"dx\"},{\"type\":\"uint256\",\"name\":\"min_dy\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":6361682},{\"name\":\"exchange_underlying\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"int128\",\"name\":\"j\"},{\"type\":\"uint256\",\"name\":\"dx\"},{\"type\":\"uint256\",\"name\":\"min_dy\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":6369753},{\"name\":\"remove_liquidity\",\"outputs\":[{\"type\":\"uint256[3]\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"_amount\"},{\"type\":\"uint256[3]\",\"name\":\"_min_amounts\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"remove_liquidity\",\"outputs\":[{\"type\":\"uint256[3]\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"_amount\"},{\"type\":\"uint256[3]\",\"name\":\"_min_amounts\"},{\"type\":\"bool\",\"name\":\"_use_underlying\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"remove_liquidity_imbalance\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256[3]\",\"name\":\"_amounts\"},{\"type\":\"uint256\",\"name\":\"_max_burn_amount\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"remove_liquidity_imbalance\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256[3]\",\"name\":\"_amounts\"},{\"type\":\"uint256\",\"name\":\"_max_burn_amount\"},{\"type\":\"bool\",\"name\":\"_use_underlying\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"calc_withdraw_one_coin\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"_token_amount\"},{\"type\":\"int128\",\"name\":\"i\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":4449067},{\"name\":\"remove_liquidity_one_coin\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"_token_amount\"},{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"uint256\",\"name\":\"_min_amount\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"remove_liquidity_one_coin\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"_token_amount\"},{\"type\":\"int128\",\"name\":\"i\"},{\"type\":\"uint256\",\"name\":\"_min_amount\"},{\"type\":\"bool\",\"name\":\"_use_underlying\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"name\":\"ramp_A\",\"outputs\":[],\"inputs\":[{\"type\":\"uint256\",\"name\":\"_future_A\"},{\"type\":\"uint256\",\"name\":\"_future_time\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":151954},{\"name\":\"stop_ramp_A\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":148715},{\"name\":\"commit_new_fee\",\"outputs\":[],\"inputs\":[{\"type\":\"uint256\",\"name\":\"new_fee\"},{\"type\":\"uint256\",\"name\":\"new_admin_fee\"},{\"type\":\"uint256\",\"name\":\"new_offpeg_fee_multiplier\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":146482},{\"name\":\"apply_new_fee\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":133744},{\"name\":\"revert_new_parameters\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":21985},{\"name\":\"commit_transfer_ownership\",\"outputs\":[],\"inputs\":[{\"type\":\"address\",\"name\":\"_owner\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":74723},{\"name\":\"apply_transfer_ownership\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":60800},{\"name\":\"revert_transfer_ownership\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":22075},{\"name\":\"withdraw_admin_fees\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":71651},{\"name\":\"donate_admin_fees\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":62276},{\"name\":\"kill_me\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":38058},{\"name\":\"unkill_me\",\"outputs\":[],\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":22195},{\"name\":\"set_aave_referral\",\"outputs\":[],\"inputs\":[{\"type\":\"uint256\",\"name\":\"referral_code\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"gas\":37325},{\"name\":\"coins\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"arg0\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2310},{\"name\":\"underlying_coins\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"arg0\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2340},{\"name\":\"admin_balances\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[{\"type\":\"uint256\",\"name\":\"arg0\"}],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2370},{\"name\":\"fee\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2291},{\"name\":\"offpeg_fee_multiplier\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2321},{\"name\":\"admin_fee\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2351},{\"name\":\"owner\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2381},{\"name\":\"lp_token\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2411},{\"name\":\"initial_A\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2441},{\"name\":\"future_A\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2471},{\"name\":\"initial_A_time\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2501},{\"name\":\"future_A_time\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2531},{\"name\":\"admin_actions_deadline\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2561},{\"name\":\"transfer_ownership_deadline\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2591},{\"name\":\"future_fee\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2621},{\"name\":\"future_admin_fee\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2651},{\"name\":\"future_offpeg_fee_multiplier\",\"outputs\":[{\"type\":\"uint256\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2681},{\"name\":\"future_owner\",\"outputs\":[{\"type\":\"address\",\"name\":\"\"}],\"inputs\":[],\"stateMutability\":\"view\",\"type\":\"function\",\"gas\":2711}]",
}

// CurvelendingpoolABI is the input ABI used to generate the binding from.
// Deprecated: Use CurvelendingpoolMetaData.ABI instead.
var CurvelendingpoolABI = CurvelendingpoolMetaData.ABI

// Curvelendingpool is an auto generated Go binding around an Ethereum contract.
type Curvelendingpool struct {
	CurvelendingpoolCaller     // Read-only binding to the contract
	CurvelendingpoolTransactor // Write-only binding to the contract
	CurvelendingpoolFilterer   // Log filterer for contract events
}

// CurvelendingpoolCaller is an auto generated read-only Go binding around an Ethereum contract.
type CurvelendingpoolCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurvelendingpoolTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CurvelendingpoolTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurvelendingpoolFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CurvelendingpoolFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurvelendingpoolSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CurvelendingpoolSession struct {
	Contract     *Curvelendingpool // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CurvelendingpoolCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CurvelendingpoolCallerSession struct {
	Contract *CurvelendingpoolCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// CurvelendingpoolTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CurvelendingpoolTransactorSession struct {
	Contract     *CurvelendingpoolTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// CurvelendingpoolRaw is an auto generated low-level Go binding around an Ethereum contract.
type CurvelendingpoolRaw struct {
	Contract *Curvelendingpool // Generic contract binding to access the raw methods on
}

// CurvelendingpoolCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CurvelendingpoolCallerRaw struct {
	Contract *CurvelendingpoolCaller // Generic read-only contract binding to access the raw methods on
}

// CurvelendingpoolTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CurvelendingpoolTransactorRaw struct {
	Contract *CurvelendingpoolTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCurvelendingpool creates a new instance of Curvelendingpool, bound to a specific deployed contract.
func NewCurvelendingpool(address common.Address, backend bind.ContractBackend) (*Curvelendingpool, error) {
	contract, err := bindCurvelendingpool(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Curvelendingpool{CurvelendingpoolCaller: CurvelendingpoolCaller{contract: contract}, CurvelendingpoolTransactor: CurvelendingpoolTransactor{contract: contract}, CurvelendingpoolFilterer: CurvelendingpoolFilterer{contract: contract}}, nil
}

// NewCurvelendingpoolCaller creates a new read-only instance of Curvelendingpool, bound to a specific deployed contract.
func NewCurvelendingpoolCaller(address common.Address, caller bind.ContractCaller) (*CurvelendingpoolCaller, error) {
	contract, err := bindCurvelendingpool(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolCaller{contract: contract}, nil
}

// NewCurvelendingpoolTransactor creates a new write-only instance of Curvelendingpool, bound to a specific deployed contract.
func NewCurvelendingpoolTransactor(address common.Address, transactor bind.ContractTransactor) (*CurvelendingpoolTransactor, error) {
	contract, err := bindCurvelendingpool(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolTransactor{contract: contract}, nil
}

// NewCurvelendingpoolFilterer creates a new log filterer instance of Curvelendingpool, bound to a specific deployed contract.
func NewCurvelendingpoolFilterer(address common.Address, filterer bind.ContractFilterer) (*CurvelendingpoolFilterer, error) {
	contract, err := bindCurvelendingpool(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolFilterer{contract: contract}, nil
}

// bindCurvelendingpool binds a generic wrapper to an already deployed contract.
func bindCurvelendingpool(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CurvelendingpoolMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Curvelendingpool *CurvelendingpoolRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Curvelendingpool.Contract.CurvelendingpoolCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Curvelendingpool *CurvelendingpoolRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.CurvelendingpoolTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Curvelendingpool *CurvelendingpoolRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.CurvelendingpoolTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Curvelendingpool *CurvelendingpoolCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Curvelendingpool.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Curvelendingpool *CurvelendingpoolTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Curvelendingpool *CurvelendingpoolTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.contract.Transact(opts, method, params...)
}

// A is a free data retrieval call binding the contract method 0xf446c1d0.
//
// Solidity: function A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) A(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "A")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// A is a free data retrieval call binding the contract method 0xf446c1d0.
//
// Solidity: function A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) A() (*big.Int, error) {
	return _Curvelendingpool.Contract.A(&_Curvelendingpool.CallOpts)
}

// A is a free data retrieval call binding the contract method 0xf446c1d0.
//
// Solidity: function A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) A() (*big.Int, error) {
	return _Curvelendingpool.Contract.A(&_Curvelendingpool.CallOpts)
}

// APrecise is a free data retrieval call binding the contract method 0x76a2f0f0.
//
// Solidity: function A_precise() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) APrecise(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "A_precise")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// APrecise is a free data retrieval call binding the contract method 0x76a2f0f0.
//
// Solidity: function A_precise() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) APrecise() (*big.Int, error) {
	return _Curvelendingpool.Contract.APrecise(&_Curvelendingpool.CallOpts)
}

// APrecise is a free data retrieval call binding the contract method 0x76a2f0f0.
//
// Solidity: function A_precise() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) APrecise() (*big.Int, error) {
	return _Curvelendingpool.Contract.APrecise(&_Curvelendingpool.CallOpts)
}

// AdminActionsDeadline is a free data retrieval call binding the contract method 0x405e28f8.
//
// Solidity: function admin_actions_deadline() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) AdminActionsDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "admin_actions_deadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AdminActionsDeadline is a free data retrieval call binding the contract method 0x405e28f8.
//
// Solidity: function admin_actions_deadline() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) AdminActionsDeadline() (*big.Int, error) {
	return _Curvelendingpool.Contract.AdminActionsDeadline(&_Curvelendingpool.CallOpts)
}

// AdminActionsDeadline is a free data retrieval call binding the contract method 0x405e28f8.
//
// Solidity: function admin_actions_deadline() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) AdminActionsDeadline() (*big.Int, error) {
	return _Curvelendingpool.Contract.AdminActionsDeadline(&_Curvelendingpool.CallOpts)
}

// AdminBalances is a free data retrieval call binding the contract method 0xe2e7d264.
//
// Solidity: function admin_balances(uint256 arg0) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) AdminBalances(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "admin_balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AdminBalances is a free data retrieval call binding the contract method 0xe2e7d264.
//
// Solidity: function admin_balances(uint256 arg0) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) AdminBalances(arg0 *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.AdminBalances(&_Curvelendingpool.CallOpts, arg0)
}

// AdminBalances is a free data retrieval call binding the contract method 0xe2e7d264.
//
// Solidity: function admin_balances(uint256 arg0) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) AdminBalances(arg0 *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.AdminBalances(&_Curvelendingpool.CallOpts, arg0)
}

// AdminFee is a free data retrieval call binding the contract method 0xfee3f7f9.
//
// Solidity: function admin_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) AdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "admin_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AdminFee is a free data retrieval call binding the contract method 0xfee3f7f9.
//
// Solidity: function admin_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) AdminFee() (*big.Int, error) {
	return _Curvelendingpool.Contract.AdminFee(&_Curvelendingpool.CallOpts)
}

// AdminFee is a free data retrieval call binding the contract method 0xfee3f7f9.
//
// Solidity: function admin_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) AdminFee() (*big.Int, error) {
	return _Curvelendingpool.Contract.AdminFee(&_Curvelendingpool.CallOpts)
}

// Balances is a free data retrieval call binding the contract method 0x4903b0d1.
//
// Solidity: function balances(uint256 i) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) Balances(opts *bind.CallOpts, i *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "balances", i)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x4903b0d1.
//
// Solidity: function balances(uint256 i) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) Balances(i *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.Balances(&_Curvelendingpool.CallOpts, i)
}

// Balances is a free data retrieval call binding the contract method 0x4903b0d1.
//
// Solidity: function balances(uint256 i) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) Balances(i *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.Balances(&_Curvelendingpool.CallOpts, i)
}

// CalcTokenAmount is a free data retrieval call binding the contract method 0x3883e119.
//
// Solidity: function calc_token_amount(uint256[3] _amounts, bool is_deposit) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) CalcTokenAmount(opts *bind.CallOpts, _amounts [3]*big.Int, is_deposit bool) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "calc_token_amount", _amounts, is_deposit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcTokenAmount is a free data retrieval call binding the contract method 0x3883e119.
//
// Solidity: function calc_token_amount(uint256[3] _amounts, bool is_deposit) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) CalcTokenAmount(_amounts [3]*big.Int, is_deposit bool) (*big.Int, error) {
	return _Curvelendingpool.Contract.CalcTokenAmount(&_Curvelendingpool.CallOpts, _amounts, is_deposit)
}

// CalcTokenAmount is a free data retrieval call binding the contract method 0x3883e119.
//
// Solidity: function calc_token_amount(uint256[3] _amounts, bool is_deposit) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) CalcTokenAmount(_amounts [3]*big.Int, is_deposit bool) (*big.Int, error) {
	return _Curvelendingpool.Contract.CalcTokenAmount(&_Curvelendingpool.CallOpts, _amounts, is_deposit)
}

// CalcWithdrawOneCoin is a free data retrieval call binding the contract method 0xcc2b27d7.
//
// Solidity: function calc_withdraw_one_coin(uint256 _token_amount, int128 i) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) CalcWithdrawOneCoin(opts *bind.CallOpts, _token_amount *big.Int, i *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "calc_withdraw_one_coin", _token_amount, i)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcWithdrawOneCoin is a free data retrieval call binding the contract method 0xcc2b27d7.
//
// Solidity: function calc_withdraw_one_coin(uint256 _token_amount, int128 i) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) CalcWithdrawOneCoin(_token_amount *big.Int, i *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.CalcWithdrawOneCoin(&_Curvelendingpool.CallOpts, _token_amount, i)
}

// CalcWithdrawOneCoin is a free data retrieval call binding the contract method 0xcc2b27d7.
//
// Solidity: function calc_withdraw_one_coin(uint256 _token_amount, int128 i) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) CalcWithdrawOneCoin(_token_amount *big.Int, i *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.CalcWithdrawOneCoin(&_Curvelendingpool.CallOpts, _token_amount, i)
}

// Coins is a free data retrieval call binding the contract method 0xc6610657.
//
// Solidity: function coins(uint256 arg0) view returns(address)
func (_Curvelendingpool *CurvelendingpoolCaller) Coins(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "coins", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Coins is a free data retrieval call binding the contract method 0xc6610657.
//
// Solidity: function coins(uint256 arg0) view returns(address)
func (_Curvelendingpool *CurvelendingpoolSession) Coins(arg0 *big.Int) (common.Address, error) {
	return _Curvelendingpool.Contract.Coins(&_Curvelendingpool.CallOpts, arg0)
}

// Coins is a free data retrieval call binding the contract method 0xc6610657.
//
// Solidity: function coins(uint256 arg0) view returns(address)
func (_Curvelendingpool *CurvelendingpoolCallerSession) Coins(arg0 *big.Int) (common.Address, error) {
	return _Curvelendingpool.Contract.Coins(&_Curvelendingpool.CallOpts, arg0)
}

// DynamicFee is a free data retrieval call binding the contract method 0x76a9cd3e.
//
// Solidity: function dynamic_fee(int128 i, int128 j) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) DynamicFee(opts *bind.CallOpts, i *big.Int, j *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "dynamic_fee", i, j)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DynamicFee is a free data retrieval call binding the contract method 0x76a9cd3e.
//
// Solidity: function dynamic_fee(int128 i, int128 j) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) DynamicFee(i *big.Int, j *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.DynamicFee(&_Curvelendingpool.CallOpts, i, j)
}

// DynamicFee is a free data retrieval call binding the contract method 0x76a9cd3e.
//
// Solidity: function dynamic_fee(int128 i, int128 j) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) DynamicFee(i *big.Int, j *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.DynamicFee(&_Curvelendingpool.CallOpts, i, j)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) Fee() (*big.Int, error) {
	return _Curvelendingpool.Contract.Fee(&_Curvelendingpool.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) Fee() (*big.Int, error) {
	return _Curvelendingpool.Contract.Fee(&_Curvelendingpool.CallOpts)
}

// FutureA is a free data retrieval call binding the contract method 0xb4b577ad.
//
// Solidity: function future_A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) FutureA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "future_A")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureA is a free data retrieval call binding the contract method 0xb4b577ad.
//
// Solidity: function future_A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) FutureA() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureA(&_Curvelendingpool.CallOpts)
}

// FutureA is a free data retrieval call binding the contract method 0xb4b577ad.
//
// Solidity: function future_A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) FutureA() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureA(&_Curvelendingpool.CallOpts)
}

// FutureATime is a free data retrieval call binding the contract method 0x14052288.
//
// Solidity: function future_A_time() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) FutureATime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "future_A_time")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureATime is a free data retrieval call binding the contract method 0x14052288.
//
// Solidity: function future_A_time() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) FutureATime() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureATime(&_Curvelendingpool.CallOpts)
}

// FutureATime is a free data retrieval call binding the contract method 0x14052288.
//
// Solidity: function future_A_time() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) FutureATime() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureATime(&_Curvelendingpool.CallOpts)
}

// FutureAdminFee is a free data retrieval call binding the contract method 0xe3824462.
//
// Solidity: function future_admin_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) FutureAdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "future_admin_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureAdminFee is a free data retrieval call binding the contract method 0xe3824462.
//
// Solidity: function future_admin_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) FutureAdminFee() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureAdminFee(&_Curvelendingpool.CallOpts)
}

// FutureAdminFee is a free data retrieval call binding the contract method 0xe3824462.
//
// Solidity: function future_admin_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) FutureAdminFee() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureAdminFee(&_Curvelendingpool.CallOpts)
}

// FutureFee is a free data retrieval call binding the contract method 0x58680d0b.
//
// Solidity: function future_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) FutureFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "future_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureFee is a free data retrieval call binding the contract method 0x58680d0b.
//
// Solidity: function future_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) FutureFee() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureFee(&_Curvelendingpool.CallOpts)
}

// FutureFee is a free data retrieval call binding the contract method 0x58680d0b.
//
// Solidity: function future_fee() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) FutureFee() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureFee(&_Curvelendingpool.CallOpts)
}

// FutureOffpegFeeMultiplier is a free data retrieval call binding the contract method 0x1e4c4ef8.
//
// Solidity: function future_offpeg_fee_multiplier() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) FutureOffpegFeeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "future_offpeg_fee_multiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureOffpegFeeMultiplier is a free data retrieval call binding the contract method 0x1e4c4ef8.
//
// Solidity: function future_offpeg_fee_multiplier() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) FutureOffpegFeeMultiplier() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureOffpegFeeMultiplier(&_Curvelendingpool.CallOpts)
}

// FutureOffpegFeeMultiplier is a free data retrieval call binding the contract method 0x1e4c4ef8.
//
// Solidity: function future_offpeg_fee_multiplier() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) FutureOffpegFeeMultiplier() (*big.Int, error) {
	return _Curvelendingpool.Contract.FutureOffpegFeeMultiplier(&_Curvelendingpool.CallOpts)
}

// FutureOwner is a free data retrieval call binding the contract method 0x1ec0cdc1.
//
// Solidity: function future_owner() view returns(address)
func (_Curvelendingpool *CurvelendingpoolCaller) FutureOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "future_owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FutureOwner is a free data retrieval call binding the contract method 0x1ec0cdc1.
//
// Solidity: function future_owner() view returns(address)
func (_Curvelendingpool *CurvelendingpoolSession) FutureOwner() (common.Address, error) {
	return _Curvelendingpool.Contract.FutureOwner(&_Curvelendingpool.CallOpts)
}

// FutureOwner is a free data retrieval call binding the contract method 0x1ec0cdc1.
//
// Solidity: function future_owner() view returns(address)
func (_Curvelendingpool *CurvelendingpoolCallerSession) FutureOwner() (common.Address, error) {
	return _Curvelendingpool.Contract.FutureOwner(&_Curvelendingpool.CallOpts)
}

// GetDy is a free data retrieval call binding the contract method 0x5e0d443f.
//
// Solidity: function get_dy(int128 i, int128 j, uint256 dx) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) GetDy(opts *bind.CallOpts, i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "get_dy", i, j, dx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDy is a free data retrieval call binding the contract method 0x5e0d443f.
//
// Solidity: function get_dy(int128 i, int128 j, uint256 dx) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) GetDy(i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.GetDy(&_Curvelendingpool.CallOpts, i, j, dx)
}

// GetDy is a free data retrieval call binding the contract method 0x5e0d443f.
//
// Solidity: function get_dy(int128 i, int128 j, uint256 dx) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) GetDy(i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.GetDy(&_Curvelendingpool.CallOpts, i, j, dx)
}

// GetDyUnderlying is a free data retrieval call binding the contract method 0x07211ef7.
//
// Solidity: function get_dy_underlying(int128 i, int128 j, uint256 dx) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) GetDyUnderlying(opts *bind.CallOpts, i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "get_dy_underlying", i, j, dx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDyUnderlying is a free data retrieval call binding the contract method 0x07211ef7.
//
// Solidity: function get_dy_underlying(int128 i, int128 j, uint256 dx) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) GetDyUnderlying(i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.GetDyUnderlying(&_Curvelendingpool.CallOpts, i, j, dx)
}

// GetDyUnderlying is a free data retrieval call binding the contract method 0x07211ef7.
//
// Solidity: function get_dy_underlying(int128 i, int128 j, uint256 dx) view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) GetDyUnderlying(i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	return _Curvelendingpool.Contract.GetDyUnderlying(&_Curvelendingpool.CallOpts, i, j, dx)
}

// GetVirtualPrice is a free data retrieval call binding the contract method 0xbb7b8b80.
//
// Solidity: function get_virtual_price() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) GetVirtualPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "get_virtual_price")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVirtualPrice is a free data retrieval call binding the contract method 0xbb7b8b80.
//
// Solidity: function get_virtual_price() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) GetVirtualPrice() (*big.Int, error) {
	return _Curvelendingpool.Contract.GetVirtualPrice(&_Curvelendingpool.CallOpts)
}

// GetVirtualPrice is a free data retrieval call binding the contract method 0xbb7b8b80.
//
// Solidity: function get_virtual_price() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) GetVirtualPrice() (*big.Int, error) {
	return _Curvelendingpool.Contract.GetVirtualPrice(&_Curvelendingpool.CallOpts)
}

// InitialA is a free data retrieval call binding the contract method 0x5409491a.
//
// Solidity: function initial_A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) InitialA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "initial_A")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialA is a free data retrieval call binding the contract method 0x5409491a.
//
// Solidity: function initial_A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) InitialA() (*big.Int, error) {
	return _Curvelendingpool.Contract.InitialA(&_Curvelendingpool.CallOpts)
}

// InitialA is a free data retrieval call binding the contract method 0x5409491a.
//
// Solidity: function initial_A() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) InitialA() (*big.Int, error) {
	return _Curvelendingpool.Contract.InitialA(&_Curvelendingpool.CallOpts)
}

// InitialATime is a free data retrieval call binding the contract method 0x2081066c.
//
// Solidity: function initial_A_time() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) InitialATime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "initial_A_time")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialATime is a free data retrieval call binding the contract method 0x2081066c.
//
// Solidity: function initial_A_time() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) InitialATime() (*big.Int, error) {
	return _Curvelendingpool.Contract.InitialATime(&_Curvelendingpool.CallOpts)
}

// InitialATime is a free data retrieval call binding the contract method 0x2081066c.
//
// Solidity: function initial_A_time() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) InitialATime() (*big.Int, error) {
	return _Curvelendingpool.Contract.InitialATime(&_Curvelendingpool.CallOpts)
}

// LpToken is a free data retrieval call binding the contract method 0x82c63066.
//
// Solidity: function lp_token() view returns(address)
func (_Curvelendingpool *CurvelendingpoolCaller) LpToken(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "lp_token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// LpToken is a free data retrieval call binding the contract method 0x82c63066.
//
// Solidity: function lp_token() view returns(address)
func (_Curvelendingpool *CurvelendingpoolSession) LpToken() (common.Address, error) {
	return _Curvelendingpool.Contract.LpToken(&_Curvelendingpool.CallOpts)
}

// LpToken is a free data retrieval call binding the contract method 0x82c63066.
//
// Solidity: function lp_token() view returns(address)
func (_Curvelendingpool *CurvelendingpoolCallerSession) LpToken() (common.Address, error) {
	return _Curvelendingpool.Contract.LpToken(&_Curvelendingpool.CallOpts)
}

// OffpegFeeMultiplier is a free data retrieval call binding the contract method 0x8edfdd5f.
//
// Solidity: function offpeg_fee_multiplier() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) OffpegFeeMultiplier(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "offpeg_fee_multiplier")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OffpegFeeMultiplier is a free data retrieval call binding the contract method 0x8edfdd5f.
//
// Solidity: function offpeg_fee_multiplier() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) OffpegFeeMultiplier() (*big.Int, error) {
	return _Curvelendingpool.Contract.OffpegFeeMultiplier(&_Curvelendingpool.CallOpts)
}

// OffpegFeeMultiplier is a free data retrieval call binding the contract method 0x8edfdd5f.
//
// Solidity: function offpeg_fee_multiplier() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) OffpegFeeMultiplier() (*big.Int, error) {
	return _Curvelendingpool.Contract.OffpegFeeMultiplier(&_Curvelendingpool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Curvelendingpool *CurvelendingpoolCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Curvelendingpool *CurvelendingpoolSession) Owner() (common.Address, error) {
	return _Curvelendingpool.Contract.Owner(&_Curvelendingpool.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Curvelendingpool *CurvelendingpoolCallerSession) Owner() (common.Address, error) {
	return _Curvelendingpool.Contract.Owner(&_Curvelendingpool.CallOpts)
}

// TransferOwnershipDeadline is a free data retrieval call binding the contract method 0xe0a0b586.
//
// Solidity: function transfer_ownership_deadline() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCaller) TransferOwnershipDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "transfer_ownership_deadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferOwnershipDeadline is a free data retrieval call binding the contract method 0xe0a0b586.
//
// Solidity: function transfer_ownership_deadline() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) TransferOwnershipDeadline() (*big.Int, error) {
	return _Curvelendingpool.Contract.TransferOwnershipDeadline(&_Curvelendingpool.CallOpts)
}

// TransferOwnershipDeadline is a free data retrieval call binding the contract method 0xe0a0b586.
//
// Solidity: function transfer_ownership_deadline() view returns(uint256)
func (_Curvelendingpool *CurvelendingpoolCallerSession) TransferOwnershipDeadline() (*big.Int, error) {
	return _Curvelendingpool.Contract.TransferOwnershipDeadline(&_Curvelendingpool.CallOpts)
}

// UnderlyingCoins is a free data retrieval call binding the contract method 0xb9947eb0.
//
// Solidity: function underlying_coins(uint256 arg0) view returns(address)
func (_Curvelendingpool *CurvelendingpoolCaller) UnderlyingCoins(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Curvelendingpool.contract.Call(opts, &out, "underlying_coins", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// UnderlyingCoins is a free data retrieval call binding the contract method 0xb9947eb0.
//
// Solidity: function underlying_coins(uint256 arg0) view returns(address)
func (_Curvelendingpool *CurvelendingpoolSession) UnderlyingCoins(arg0 *big.Int) (common.Address, error) {
	return _Curvelendingpool.Contract.UnderlyingCoins(&_Curvelendingpool.CallOpts, arg0)
}

// UnderlyingCoins is a free data retrieval call binding the contract method 0xb9947eb0.
//
// Solidity: function underlying_coins(uint256 arg0) view returns(address)
func (_Curvelendingpool *CurvelendingpoolCallerSession) UnderlyingCoins(arg0 *big.Int) (common.Address, error) {
	return _Curvelendingpool.Contract.UnderlyingCoins(&_Curvelendingpool.CallOpts, arg0)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x4515cef3.
//
// Solidity: function add_liquidity(uint256[3] _amounts, uint256 _min_mint_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) AddLiquidity(opts *bind.TransactOpts, _amounts [3]*big.Int, _min_mint_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "add_liquidity", _amounts, _min_mint_amount)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x4515cef3.
//
// Solidity: function add_liquidity(uint256[3] _amounts, uint256 _min_mint_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) AddLiquidity(_amounts [3]*big.Int, _min_mint_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.AddLiquidity(&_Curvelendingpool.TransactOpts, _amounts, _min_mint_amount)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x4515cef3.
//
// Solidity: function add_liquidity(uint256[3] _amounts, uint256 _min_mint_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) AddLiquidity(_amounts [3]*big.Int, _min_mint_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.AddLiquidity(&_Curvelendingpool.TransactOpts, _amounts, _min_mint_amount)
}

// AddLiquidity0 is a paid mutator transaction binding the contract method 0x2b6e993a.
//
// Solidity: function add_liquidity(uint256[3] _amounts, uint256 _min_mint_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) AddLiquidity0(opts *bind.TransactOpts, _amounts [3]*big.Int, _min_mint_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "add_liquidity0", _amounts, _min_mint_amount, _use_underlying)
}

// AddLiquidity0 is a paid mutator transaction binding the contract method 0x2b6e993a.
//
// Solidity: function add_liquidity(uint256[3] _amounts, uint256 _min_mint_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) AddLiquidity0(_amounts [3]*big.Int, _min_mint_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.AddLiquidity0(&_Curvelendingpool.TransactOpts, _amounts, _min_mint_amount, _use_underlying)
}

// AddLiquidity0 is a paid mutator transaction binding the contract method 0x2b6e993a.
//
// Solidity: function add_liquidity(uint256[3] _amounts, uint256 _min_mint_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) AddLiquidity0(_amounts [3]*big.Int, _min_mint_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.AddLiquidity0(&_Curvelendingpool.TransactOpts, _amounts, _min_mint_amount, _use_underlying)
}

// ApplyNewFee is a paid mutator transaction binding the contract method 0x4f12fe97.
//
// Solidity: function apply_new_fee() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) ApplyNewFee(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "apply_new_fee")
}

// ApplyNewFee is a paid mutator transaction binding the contract method 0x4f12fe97.
//
// Solidity: function apply_new_fee() returns()
func (_Curvelendingpool *CurvelendingpoolSession) ApplyNewFee() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.ApplyNewFee(&_Curvelendingpool.TransactOpts)
}

// ApplyNewFee is a paid mutator transaction binding the contract method 0x4f12fe97.
//
// Solidity: function apply_new_fee() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) ApplyNewFee() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.ApplyNewFee(&_Curvelendingpool.TransactOpts)
}

// ApplyTransferOwnership is a paid mutator transaction binding the contract method 0x6a1c05ae.
//
// Solidity: function apply_transfer_ownership() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) ApplyTransferOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "apply_transfer_ownership")
}

// ApplyTransferOwnership is a paid mutator transaction binding the contract method 0x6a1c05ae.
//
// Solidity: function apply_transfer_ownership() returns()
func (_Curvelendingpool *CurvelendingpoolSession) ApplyTransferOwnership() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.ApplyTransferOwnership(&_Curvelendingpool.TransactOpts)
}

// ApplyTransferOwnership is a paid mutator transaction binding the contract method 0x6a1c05ae.
//
// Solidity: function apply_transfer_ownership() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) ApplyTransferOwnership() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.ApplyTransferOwnership(&_Curvelendingpool.TransactOpts)
}

// CommitNewFee is a paid mutator transaction binding the contract method 0x0746dd5a.
//
// Solidity: function commit_new_fee(uint256 new_fee, uint256 new_admin_fee, uint256 new_offpeg_fee_multiplier) returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) CommitNewFee(opts *bind.TransactOpts, new_fee *big.Int, new_admin_fee *big.Int, new_offpeg_fee_multiplier *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "commit_new_fee", new_fee, new_admin_fee, new_offpeg_fee_multiplier)
}

// CommitNewFee is a paid mutator transaction binding the contract method 0x0746dd5a.
//
// Solidity: function commit_new_fee(uint256 new_fee, uint256 new_admin_fee, uint256 new_offpeg_fee_multiplier) returns()
func (_Curvelendingpool *CurvelendingpoolSession) CommitNewFee(new_fee *big.Int, new_admin_fee *big.Int, new_offpeg_fee_multiplier *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.CommitNewFee(&_Curvelendingpool.TransactOpts, new_fee, new_admin_fee, new_offpeg_fee_multiplier)
}

// CommitNewFee is a paid mutator transaction binding the contract method 0x0746dd5a.
//
// Solidity: function commit_new_fee(uint256 new_fee, uint256 new_admin_fee, uint256 new_offpeg_fee_multiplier) returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) CommitNewFee(new_fee *big.Int, new_admin_fee *big.Int, new_offpeg_fee_multiplier *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.CommitNewFee(&_Curvelendingpool.TransactOpts, new_fee, new_admin_fee, new_offpeg_fee_multiplier)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _owner) returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) CommitTransferOwnership(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "commit_transfer_ownership", _owner)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _owner) returns()
func (_Curvelendingpool *CurvelendingpoolSession) CommitTransferOwnership(_owner common.Address) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.CommitTransferOwnership(&_Curvelendingpool.TransactOpts, _owner)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _owner) returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) CommitTransferOwnership(_owner common.Address) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.CommitTransferOwnership(&_Curvelendingpool.TransactOpts, _owner)
}

// DonateAdminFees is a paid mutator transaction binding the contract method 0x524c3901.
//
// Solidity: function donate_admin_fees() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) DonateAdminFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "donate_admin_fees")
}

// DonateAdminFees is a paid mutator transaction binding the contract method 0x524c3901.
//
// Solidity: function donate_admin_fees() returns()
func (_Curvelendingpool *CurvelendingpoolSession) DonateAdminFees() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.DonateAdminFees(&_Curvelendingpool.TransactOpts)
}

// DonateAdminFees is a paid mutator transaction binding the contract method 0x524c3901.
//
// Solidity: function donate_admin_fees() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) DonateAdminFees() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.DonateAdminFees(&_Curvelendingpool.TransactOpts)
}

// Exchange is a paid mutator transaction binding the contract method 0x3df02124.
//
// Solidity: function exchange(int128 i, int128 j, uint256 dx, uint256 min_dy) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) Exchange(opts *bind.TransactOpts, i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "exchange", i, j, dx, min_dy)
}

// Exchange is a paid mutator transaction binding the contract method 0x3df02124.
//
// Solidity: function exchange(int128 i, int128 j, uint256 dx, uint256 min_dy) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) Exchange(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.Exchange(&_Curvelendingpool.TransactOpts, i, j, dx, min_dy)
}

// Exchange is a paid mutator transaction binding the contract method 0x3df02124.
//
// Solidity: function exchange(int128 i, int128 j, uint256 dx, uint256 min_dy) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) Exchange(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.Exchange(&_Curvelendingpool.TransactOpts, i, j, dx, min_dy)
}

// ExchangeUnderlying is a paid mutator transaction binding the contract method 0xa6417ed6.
//
// Solidity: function exchange_underlying(int128 i, int128 j, uint256 dx, uint256 min_dy) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) ExchangeUnderlying(opts *bind.TransactOpts, i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "exchange_underlying", i, j, dx, min_dy)
}

// ExchangeUnderlying is a paid mutator transaction binding the contract method 0xa6417ed6.
//
// Solidity: function exchange_underlying(int128 i, int128 j, uint256 dx, uint256 min_dy) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) ExchangeUnderlying(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.ExchangeUnderlying(&_Curvelendingpool.TransactOpts, i, j, dx, min_dy)
}

// ExchangeUnderlying is a paid mutator transaction binding the contract method 0xa6417ed6.
//
// Solidity: function exchange_underlying(int128 i, int128 j, uint256 dx, uint256 min_dy) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) ExchangeUnderlying(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.ExchangeUnderlying(&_Curvelendingpool.TransactOpts, i, j, dx, min_dy)
}

// KillMe is a paid mutator transaction binding the contract method 0xe3698853.
//
// Solidity: function kill_me() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) KillMe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "kill_me")
}

// KillMe is a paid mutator transaction binding the contract method 0xe3698853.
//
// Solidity: function kill_me() returns()
func (_Curvelendingpool *CurvelendingpoolSession) KillMe() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.KillMe(&_Curvelendingpool.TransactOpts)
}

// KillMe is a paid mutator transaction binding the contract method 0xe3698853.
//
// Solidity: function kill_me() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) KillMe() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.KillMe(&_Curvelendingpool.TransactOpts)
}

// RampA is a paid mutator transaction binding the contract method 0x3c157e64.
//
// Solidity: function ramp_A(uint256 _future_A, uint256 _future_time) returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) RampA(opts *bind.TransactOpts, _future_A *big.Int, _future_time *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "ramp_A", _future_A, _future_time)
}

// RampA is a paid mutator transaction binding the contract method 0x3c157e64.
//
// Solidity: function ramp_A(uint256 _future_A, uint256 _future_time) returns()
func (_Curvelendingpool *CurvelendingpoolSession) RampA(_future_A *big.Int, _future_time *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RampA(&_Curvelendingpool.TransactOpts, _future_A, _future_time)
}

// RampA is a paid mutator transaction binding the contract method 0x3c157e64.
//
// Solidity: function ramp_A(uint256 _future_A, uint256 _future_time) returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RampA(_future_A *big.Int, _future_time *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RampA(&_Curvelendingpool.TransactOpts, _future_A, _future_time)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xecb586a5.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] _min_amounts) returns(uint256[3])
func (_Curvelendingpool *CurvelendingpoolTransactor) RemoveLiquidity(opts *bind.TransactOpts, _amount *big.Int, _min_amounts [3]*big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "remove_liquidity", _amount, _min_amounts)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xecb586a5.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] _min_amounts) returns(uint256[3])
func (_Curvelendingpool *CurvelendingpoolSession) RemoveLiquidity(_amount *big.Int, _min_amounts [3]*big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidity(&_Curvelendingpool.TransactOpts, _amount, _min_amounts)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xecb586a5.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] _min_amounts) returns(uint256[3])
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RemoveLiquidity(_amount *big.Int, _min_amounts [3]*big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidity(&_Curvelendingpool.TransactOpts, _amount, _min_amounts)
}

// RemoveLiquidity0 is a paid mutator transaction binding the contract method 0xfce64736.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] _min_amounts, bool _use_underlying) returns(uint256[3])
func (_Curvelendingpool *CurvelendingpoolTransactor) RemoveLiquidity0(opts *bind.TransactOpts, _amount *big.Int, _min_amounts [3]*big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "remove_liquidity0", _amount, _min_amounts, _use_underlying)
}

// RemoveLiquidity0 is a paid mutator transaction binding the contract method 0xfce64736.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] _min_amounts, bool _use_underlying) returns(uint256[3])
func (_Curvelendingpool *CurvelendingpoolSession) RemoveLiquidity0(_amount *big.Int, _min_amounts [3]*big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidity0(&_Curvelendingpool.TransactOpts, _amount, _min_amounts, _use_underlying)
}

// RemoveLiquidity0 is a paid mutator transaction binding the contract method 0xfce64736.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] _min_amounts, bool _use_underlying) returns(uint256[3])
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RemoveLiquidity0(_amount *big.Int, _min_amounts [3]*big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidity0(&_Curvelendingpool.TransactOpts, _amount, _min_amounts, _use_underlying)
}

// RemoveLiquidityImbalance is a paid mutator transaction binding the contract method 0x9fdaea0c.
//
// Solidity: function remove_liquidity_imbalance(uint256[3] _amounts, uint256 _max_burn_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) RemoveLiquidityImbalance(opts *bind.TransactOpts, _amounts [3]*big.Int, _max_burn_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "remove_liquidity_imbalance", _amounts, _max_burn_amount)
}

// RemoveLiquidityImbalance is a paid mutator transaction binding the contract method 0x9fdaea0c.
//
// Solidity: function remove_liquidity_imbalance(uint256[3] _amounts, uint256 _max_burn_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) RemoveLiquidityImbalance(_amounts [3]*big.Int, _max_burn_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityImbalance(&_Curvelendingpool.TransactOpts, _amounts, _max_burn_amount)
}

// RemoveLiquidityImbalance is a paid mutator transaction binding the contract method 0x9fdaea0c.
//
// Solidity: function remove_liquidity_imbalance(uint256[3] _amounts, uint256 _max_burn_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RemoveLiquidityImbalance(_amounts [3]*big.Int, _max_burn_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityImbalance(&_Curvelendingpool.TransactOpts, _amounts, _max_burn_amount)
}

// RemoveLiquidityImbalance0 is a paid mutator transaction binding the contract method 0x5b8369f5.
//
// Solidity: function remove_liquidity_imbalance(uint256[3] _amounts, uint256 _max_burn_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) RemoveLiquidityImbalance0(opts *bind.TransactOpts, _amounts [3]*big.Int, _max_burn_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "remove_liquidity_imbalance0", _amounts, _max_burn_amount, _use_underlying)
}

// RemoveLiquidityImbalance0 is a paid mutator transaction binding the contract method 0x5b8369f5.
//
// Solidity: function remove_liquidity_imbalance(uint256[3] _amounts, uint256 _max_burn_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) RemoveLiquidityImbalance0(_amounts [3]*big.Int, _max_burn_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityImbalance0(&_Curvelendingpool.TransactOpts, _amounts, _max_burn_amount, _use_underlying)
}

// RemoveLiquidityImbalance0 is a paid mutator transaction binding the contract method 0x5b8369f5.
//
// Solidity: function remove_liquidity_imbalance(uint256[3] _amounts, uint256 _max_burn_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RemoveLiquidityImbalance0(_amounts [3]*big.Int, _max_burn_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityImbalance0(&_Curvelendingpool.TransactOpts, _amounts, _max_burn_amount, _use_underlying)
}

// RemoveLiquidityOneCoin is a paid mutator transaction binding the contract method 0x1a4d01d2.
//
// Solidity: function remove_liquidity_one_coin(uint256 _token_amount, int128 i, uint256 _min_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) RemoveLiquidityOneCoin(opts *bind.TransactOpts, _token_amount *big.Int, i *big.Int, _min_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "remove_liquidity_one_coin", _token_amount, i, _min_amount)
}

// RemoveLiquidityOneCoin is a paid mutator transaction binding the contract method 0x1a4d01d2.
//
// Solidity: function remove_liquidity_one_coin(uint256 _token_amount, int128 i, uint256 _min_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) RemoveLiquidityOneCoin(_token_amount *big.Int, i *big.Int, _min_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityOneCoin(&_Curvelendingpool.TransactOpts, _token_amount, i, _min_amount)
}

// RemoveLiquidityOneCoin is a paid mutator transaction binding the contract method 0x1a4d01d2.
//
// Solidity: function remove_liquidity_one_coin(uint256 _token_amount, int128 i, uint256 _min_amount) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RemoveLiquidityOneCoin(_token_amount *big.Int, i *big.Int, _min_amount *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityOneCoin(&_Curvelendingpool.TransactOpts, _token_amount, i, _min_amount)
}

// RemoveLiquidityOneCoin0 is a paid mutator transaction binding the contract method 0x517a55a3.
//
// Solidity: function remove_liquidity_one_coin(uint256 _token_amount, int128 i, uint256 _min_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactor) RemoveLiquidityOneCoin0(opts *bind.TransactOpts, _token_amount *big.Int, i *big.Int, _min_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "remove_liquidity_one_coin0", _token_amount, i, _min_amount, _use_underlying)
}

// RemoveLiquidityOneCoin0 is a paid mutator transaction binding the contract method 0x517a55a3.
//
// Solidity: function remove_liquidity_one_coin(uint256 _token_amount, int128 i, uint256 _min_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolSession) RemoveLiquidityOneCoin0(_token_amount *big.Int, i *big.Int, _min_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityOneCoin0(&_Curvelendingpool.TransactOpts, _token_amount, i, _min_amount, _use_underlying)
}

// RemoveLiquidityOneCoin0 is a paid mutator transaction binding the contract method 0x517a55a3.
//
// Solidity: function remove_liquidity_one_coin(uint256 _token_amount, int128 i, uint256 _min_amount, bool _use_underlying) returns(uint256)
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RemoveLiquidityOneCoin0(_token_amount *big.Int, i *big.Int, _min_amount *big.Int, _use_underlying bool) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RemoveLiquidityOneCoin0(&_Curvelendingpool.TransactOpts, _token_amount, i, _min_amount, _use_underlying)
}

// RevertNewParameters is a paid mutator transaction binding the contract method 0x226840fb.
//
// Solidity: function revert_new_parameters() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) RevertNewParameters(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "revert_new_parameters")
}

// RevertNewParameters is a paid mutator transaction binding the contract method 0x226840fb.
//
// Solidity: function revert_new_parameters() returns()
func (_Curvelendingpool *CurvelendingpoolSession) RevertNewParameters() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RevertNewParameters(&_Curvelendingpool.TransactOpts)
}

// RevertNewParameters is a paid mutator transaction binding the contract method 0x226840fb.
//
// Solidity: function revert_new_parameters() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RevertNewParameters() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RevertNewParameters(&_Curvelendingpool.TransactOpts)
}

// RevertTransferOwnership is a paid mutator transaction binding the contract method 0x86fbf193.
//
// Solidity: function revert_transfer_ownership() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) RevertTransferOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "revert_transfer_ownership")
}

// RevertTransferOwnership is a paid mutator transaction binding the contract method 0x86fbf193.
//
// Solidity: function revert_transfer_ownership() returns()
func (_Curvelendingpool *CurvelendingpoolSession) RevertTransferOwnership() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RevertTransferOwnership(&_Curvelendingpool.TransactOpts)
}

// RevertTransferOwnership is a paid mutator transaction binding the contract method 0x86fbf193.
//
// Solidity: function revert_transfer_ownership() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) RevertTransferOwnership() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.RevertTransferOwnership(&_Curvelendingpool.TransactOpts)
}

// SetAaveReferral is a paid mutator transaction binding the contract method 0xb6aa64c5.
//
// Solidity: function set_aave_referral(uint256 referral_code) returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) SetAaveReferral(opts *bind.TransactOpts, referral_code *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "set_aave_referral", referral_code)
}

// SetAaveReferral is a paid mutator transaction binding the contract method 0xb6aa64c5.
//
// Solidity: function set_aave_referral(uint256 referral_code) returns()
func (_Curvelendingpool *CurvelendingpoolSession) SetAaveReferral(referral_code *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.SetAaveReferral(&_Curvelendingpool.TransactOpts, referral_code)
}

// SetAaveReferral is a paid mutator transaction binding the contract method 0xb6aa64c5.
//
// Solidity: function set_aave_referral(uint256 referral_code) returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) SetAaveReferral(referral_code *big.Int) (*types.Transaction, error) {
	return _Curvelendingpool.Contract.SetAaveReferral(&_Curvelendingpool.TransactOpts, referral_code)
}

// StopRampA is a paid mutator transaction binding the contract method 0x551a6588.
//
// Solidity: function stop_ramp_A() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) StopRampA(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "stop_ramp_A")
}

// StopRampA is a paid mutator transaction binding the contract method 0x551a6588.
//
// Solidity: function stop_ramp_A() returns()
func (_Curvelendingpool *CurvelendingpoolSession) StopRampA() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.StopRampA(&_Curvelendingpool.TransactOpts)
}

// StopRampA is a paid mutator transaction binding the contract method 0x551a6588.
//
// Solidity: function stop_ramp_A() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) StopRampA() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.StopRampA(&_Curvelendingpool.TransactOpts)
}

// UnkillMe is a paid mutator transaction binding the contract method 0x3046f972.
//
// Solidity: function unkill_me() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) UnkillMe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "unkill_me")
}

// UnkillMe is a paid mutator transaction binding the contract method 0x3046f972.
//
// Solidity: function unkill_me() returns()
func (_Curvelendingpool *CurvelendingpoolSession) UnkillMe() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.UnkillMe(&_Curvelendingpool.TransactOpts)
}

// UnkillMe is a paid mutator transaction binding the contract method 0x3046f972.
//
// Solidity: function unkill_me() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) UnkillMe() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.UnkillMe(&_Curvelendingpool.TransactOpts)
}

// WithdrawAdminFees is a paid mutator transaction binding the contract method 0x30c54085.
//
// Solidity: function withdraw_admin_fees() returns()
func (_Curvelendingpool *CurvelendingpoolTransactor) WithdrawAdminFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curvelendingpool.contract.Transact(opts, "withdraw_admin_fees")
}

// WithdrawAdminFees is a paid mutator transaction binding the contract method 0x30c54085.
//
// Solidity: function withdraw_admin_fees() returns()
func (_Curvelendingpool *CurvelendingpoolSession) WithdrawAdminFees() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.WithdrawAdminFees(&_Curvelendingpool.TransactOpts)
}

// WithdrawAdminFees is a paid mutator transaction binding the contract method 0x30c54085.
//
// Solidity: function withdraw_admin_fees() returns()
func (_Curvelendingpool *CurvelendingpoolTransactorSession) WithdrawAdminFees() (*types.Transaction, error) {
	return _Curvelendingpool.Contract.WithdrawAdminFees(&_Curvelendingpool.TransactOpts)
}

// CurvelendingpoolAddLiquidityIterator is returned from FilterAddLiquidity and is used to iterate over the raw logs and unpacked data for AddLiquidity events raised by the Curvelendingpool contract.
type CurvelendingpoolAddLiquidityIterator struct {
	Event *CurvelendingpoolAddLiquidity // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolAddLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolAddLiquidity)
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
		it.Event = new(CurvelendingpoolAddLiquidity)
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
func (it *CurvelendingpoolAddLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolAddLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolAddLiquidity represents a AddLiquidity event raised by the Curvelendingpool contract.
type CurvelendingpoolAddLiquidity struct {
	Provider     common.Address
	TokenAmounts [3]*big.Int
	Fees         [3]*big.Int
	Invariant    *big.Int
	TokenSupply  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAddLiquidity is a free log retrieval operation binding the contract event 0x423f6495a08fc652425cf4ed0d1f9e37e571d9b9529b1c1c23cce780b2e7df0d.
//
// Solidity: event AddLiquidity(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 invariant, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterAddLiquidity(opts *bind.FilterOpts, provider []common.Address) (*CurvelendingpoolAddLiquidityIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "AddLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolAddLiquidityIterator{contract: _Curvelendingpool.contract, event: "AddLiquidity", logs: logs, sub: sub}, nil
}

// WatchAddLiquidity is a free log subscription operation binding the contract event 0x423f6495a08fc652425cf4ed0d1f9e37e571d9b9529b1c1c23cce780b2e7df0d.
//
// Solidity: event AddLiquidity(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 invariant, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchAddLiquidity(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolAddLiquidity, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "AddLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolAddLiquidity)
				if err := _Curvelendingpool.contract.UnpackLog(event, "AddLiquidity", log); err != nil {
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

// ParseAddLiquidity is a log parse operation binding the contract event 0x423f6495a08fc652425cf4ed0d1f9e37e571d9b9529b1c1c23cce780b2e7df0d.
//
// Solidity: event AddLiquidity(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 invariant, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseAddLiquidity(log types.Log) (*CurvelendingpoolAddLiquidity, error) {
	event := new(CurvelendingpoolAddLiquidity)
	if err := _Curvelendingpool.contract.UnpackLog(event, "AddLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolCommitNewAdminIterator is returned from FilterCommitNewAdmin and is used to iterate over the raw logs and unpacked data for CommitNewAdmin events raised by the Curvelendingpool contract.
type CurvelendingpoolCommitNewAdminIterator struct {
	Event *CurvelendingpoolCommitNewAdmin // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolCommitNewAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolCommitNewAdmin)
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
		it.Event = new(CurvelendingpoolCommitNewAdmin)
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
func (it *CurvelendingpoolCommitNewAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolCommitNewAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolCommitNewAdmin represents a CommitNewAdmin event raised by the Curvelendingpool contract.
type CurvelendingpoolCommitNewAdmin struct {
	Deadline *big.Int
	Admin    common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCommitNewAdmin is a free log retrieval operation binding the contract event 0x181aa3aa17d4cbf99265dd4443eba009433d3cde79d60164fde1d1a192beb935.
//
// Solidity: event CommitNewAdmin(uint256 indexed deadline, address indexed admin)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterCommitNewAdmin(opts *bind.FilterOpts, deadline []*big.Int, admin []common.Address) (*CurvelendingpoolCommitNewAdminIterator, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "CommitNewAdmin", deadlineRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolCommitNewAdminIterator{contract: _Curvelendingpool.contract, event: "CommitNewAdmin", logs: logs, sub: sub}, nil
}

// WatchCommitNewAdmin is a free log subscription operation binding the contract event 0x181aa3aa17d4cbf99265dd4443eba009433d3cde79d60164fde1d1a192beb935.
//
// Solidity: event CommitNewAdmin(uint256 indexed deadline, address indexed admin)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchCommitNewAdmin(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolCommitNewAdmin, deadline []*big.Int, admin []common.Address) (event.Subscription, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "CommitNewAdmin", deadlineRule, adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolCommitNewAdmin)
				if err := _Curvelendingpool.contract.UnpackLog(event, "CommitNewAdmin", log); err != nil {
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

// ParseCommitNewAdmin is a log parse operation binding the contract event 0x181aa3aa17d4cbf99265dd4443eba009433d3cde79d60164fde1d1a192beb935.
//
// Solidity: event CommitNewAdmin(uint256 indexed deadline, address indexed admin)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseCommitNewAdmin(log types.Log) (*CurvelendingpoolCommitNewAdmin, error) {
	event := new(CurvelendingpoolCommitNewAdmin)
	if err := _Curvelendingpool.contract.UnpackLog(event, "CommitNewAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolCommitNewFeeIterator is returned from FilterCommitNewFee and is used to iterate over the raw logs and unpacked data for CommitNewFee events raised by the Curvelendingpool contract.
type CurvelendingpoolCommitNewFeeIterator struct {
	Event *CurvelendingpoolCommitNewFee // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolCommitNewFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolCommitNewFee)
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
		it.Event = new(CurvelendingpoolCommitNewFee)
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
func (it *CurvelendingpoolCommitNewFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolCommitNewFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolCommitNewFee represents a CommitNewFee event raised by the Curvelendingpool contract.
type CurvelendingpoolCommitNewFee struct {
	Deadline            *big.Int
	Fee                 *big.Int
	AdminFee            *big.Int
	OffpegFeeMultiplier *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterCommitNewFee is a free log retrieval operation binding the contract event 0xe347cde074ab87e09449fa2b03e8f2cf79094cb1265f4c914365d2247d4147a3.
//
// Solidity: event CommitNewFee(uint256 indexed deadline, uint256 fee, uint256 admin_fee, uint256 offpeg_fee_multiplier)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterCommitNewFee(opts *bind.FilterOpts, deadline []*big.Int) (*CurvelendingpoolCommitNewFeeIterator, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "CommitNewFee", deadlineRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolCommitNewFeeIterator{contract: _Curvelendingpool.contract, event: "CommitNewFee", logs: logs, sub: sub}, nil
}

// WatchCommitNewFee is a free log subscription operation binding the contract event 0xe347cde074ab87e09449fa2b03e8f2cf79094cb1265f4c914365d2247d4147a3.
//
// Solidity: event CommitNewFee(uint256 indexed deadline, uint256 fee, uint256 admin_fee, uint256 offpeg_fee_multiplier)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchCommitNewFee(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolCommitNewFee, deadline []*big.Int) (event.Subscription, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "CommitNewFee", deadlineRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolCommitNewFee)
				if err := _Curvelendingpool.contract.UnpackLog(event, "CommitNewFee", log); err != nil {
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

// ParseCommitNewFee is a log parse operation binding the contract event 0xe347cde074ab87e09449fa2b03e8f2cf79094cb1265f4c914365d2247d4147a3.
//
// Solidity: event CommitNewFee(uint256 indexed deadline, uint256 fee, uint256 admin_fee, uint256 offpeg_fee_multiplier)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseCommitNewFee(log types.Log) (*CurvelendingpoolCommitNewFee, error) {
	event := new(CurvelendingpoolCommitNewFee)
	if err := _Curvelendingpool.contract.UnpackLog(event, "CommitNewFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolNewAdminIterator is returned from FilterNewAdmin and is used to iterate over the raw logs and unpacked data for NewAdmin events raised by the Curvelendingpool contract.
type CurvelendingpoolNewAdminIterator struct {
	Event *CurvelendingpoolNewAdmin // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolNewAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolNewAdmin)
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
		it.Event = new(CurvelendingpoolNewAdmin)
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
func (it *CurvelendingpoolNewAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolNewAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolNewAdmin represents a NewAdmin event raised by the Curvelendingpool contract.
type CurvelendingpoolNewAdmin struct {
	Admin common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNewAdmin is a free log retrieval operation binding the contract event 0x71614071b88dee5e0b2ae578a9dd7b2ebbe9ae832ba419dc0242cd065a290b6c.
//
// Solidity: event NewAdmin(address indexed admin)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterNewAdmin(opts *bind.FilterOpts, admin []common.Address) (*CurvelendingpoolNewAdminIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "NewAdmin", adminRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolNewAdminIterator{contract: _Curvelendingpool.contract, event: "NewAdmin", logs: logs, sub: sub}, nil
}

// WatchNewAdmin is a free log subscription operation binding the contract event 0x71614071b88dee5e0b2ae578a9dd7b2ebbe9ae832ba419dc0242cd065a290b6c.
//
// Solidity: event NewAdmin(address indexed admin)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchNewAdmin(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolNewAdmin, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "NewAdmin", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolNewAdmin)
				if err := _Curvelendingpool.contract.UnpackLog(event, "NewAdmin", log); err != nil {
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

// ParseNewAdmin is a log parse operation binding the contract event 0x71614071b88dee5e0b2ae578a9dd7b2ebbe9ae832ba419dc0242cd065a290b6c.
//
// Solidity: event NewAdmin(address indexed admin)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseNewAdmin(log types.Log) (*CurvelendingpoolNewAdmin, error) {
	event := new(CurvelendingpoolNewAdmin)
	if err := _Curvelendingpool.contract.UnpackLog(event, "NewAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolNewFeeIterator is returned from FilterNewFee and is used to iterate over the raw logs and unpacked data for NewFee events raised by the Curvelendingpool contract.
type CurvelendingpoolNewFeeIterator struct {
	Event *CurvelendingpoolNewFee // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolNewFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolNewFee)
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
		it.Event = new(CurvelendingpoolNewFee)
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
func (it *CurvelendingpoolNewFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolNewFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolNewFee represents a NewFee event raised by the Curvelendingpool contract.
type CurvelendingpoolNewFee struct {
	Fee                 *big.Int
	AdminFee            *big.Int
	OffpegFeeMultiplier *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterNewFee is a free log retrieval operation binding the contract event 0xcfca96e0fef3432146913b2a5a2268a55d3f475fe057e7ffde1082b77693f4f3.
//
// Solidity: event NewFee(uint256 fee, uint256 admin_fee, uint256 offpeg_fee_multiplier)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterNewFee(opts *bind.FilterOpts) (*CurvelendingpoolNewFeeIterator, error) {

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "NewFee")
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolNewFeeIterator{contract: _Curvelendingpool.contract, event: "NewFee", logs: logs, sub: sub}, nil
}

// WatchNewFee is a free log subscription operation binding the contract event 0xcfca96e0fef3432146913b2a5a2268a55d3f475fe057e7ffde1082b77693f4f3.
//
// Solidity: event NewFee(uint256 fee, uint256 admin_fee, uint256 offpeg_fee_multiplier)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchNewFee(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolNewFee) (event.Subscription, error) {

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "NewFee")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolNewFee)
				if err := _Curvelendingpool.contract.UnpackLog(event, "NewFee", log); err != nil {
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

// ParseNewFee is a log parse operation binding the contract event 0xcfca96e0fef3432146913b2a5a2268a55d3f475fe057e7ffde1082b77693f4f3.
//
// Solidity: event NewFee(uint256 fee, uint256 admin_fee, uint256 offpeg_fee_multiplier)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseNewFee(log types.Log) (*CurvelendingpoolNewFee, error) {
	event := new(CurvelendingpoolNewFee)
	if err := _Curvelendingpool.contract.UnpackLog(event, "NewFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolRampAIterator is returned from FilterRampA and is used to iterate over the raw logs and unpacked data for RampA events raised by the Curvelendingpool contract.
type CurvelendingpoolRampAIterator struct {
	Event *CurvelendingpoolRampA // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolRampAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolRampA)
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
		it.Event = new(CurvelendingpoolRampA)
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
func (it *CurvelendingpoolRampAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolRampAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolRampA represents a RampA event raised by the Curvelendingpool contract.
type CurvelendingpoolRampA struct {
	OldA        *big.Int
	NewA        *big.Int
	InitialTime *big.Int
	FutureTime  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRampA is a free log retrieval operation binding the contract event 0xa2b71ec6df949300b59aab36b55e189697b750119dd349fcfa8c0f779e83c254.
//
// Solidity: event RampA(uint256 old_A, uint256 new_A, uint256 initial_time, uint256 future_time)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterRampA(opts *bind.FilterOpts) (*CurvelendingpoolRampAIterator, error) {

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "RampA")
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolRampAIterator{contract: _Curvelendingpool.contract, event: "RampA", logs: logs, sub: sub}, nil
}

// WatchRampA is a free log subscription operation binding the contract event 0xa2b71ec6df949300b59aab36b55e189697b750119dd349fcfa8c0f779e83c254.
//
// Solidity: event RampA(uint256 old_A, uint256 new_A, uint256 initial_time, uint256 future_time)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchRampA(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolRampA) (event.Subscription, error) {

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "RampA")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolRampA)
				if err := _Curvelendingpool.contract.UnpackLog(event, "RampA", log); err != nil {
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

// ParseRampA is a log parse operation binding the contract event 0xa2b71ec6df949300b59aab36b55e189697b750119dd349fcfa8c0f779e83c254.
//
// Solidity: event RampA(uint256 old_A, uint256 new_A, uint256 initial_time, uint256 future_time)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseRampA(log types.Log) (*CurvelendingpoolRampA, error) {
	event := new(CurvelendingpoolRampA)
	if err := _Curvelendingpool.contract.UnpackLog(event, "RampA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolRemoveLiquidityIterator is returned from FilterRemoveLiquidity and is used to iterate over the raw logs and unpacked data for RemoveLiquidity events raised by the Curvelendingpool contract.
type CurvelendingpoolRemoveLiquidityIterator struct {
	Event *CurvelendingpoolRemoveLiquidity // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolRemoveLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolRemoveLiquidity)
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
		it.Event = new(CurvelendingpoolRemoveLiquidity)
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
func (it *CurvelendingpoolRemoveLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolRemoveLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolRemoveLiquidity represents a RemoveLiquidity event raised by the Curvelendingpool contract.
type CurvelendingpoolRemoveLiquidity struct {
	Provider     common.Address
	TokenAmounts [3]*big.Int
	Fees         [3]*big.Int
	TokenSupply  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemoveLiquidity is a free log retrieval operation binding the contract event 0xa49d4cf02656aebf8c771f5a8585638a2a15ee6c97cf7205d4208ed7c1df252d.
//
// Solidity: event RemoveLiquidity(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterRemoveLiquidity(opts *bind.FilterOpts, provider []common.Address) (*CurvelendingpoolRemoveLiquidityIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "RemoveLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolRemoveLiquidityIterator{contract: _Curvelendingpool.contract, event: "RemoveLiquidity", logs: logs, sub: sub}, nil
}

// WatchRemoveLiquidity is a free log subscription operation binding the contract event 0xa49d4cf02656aebf8c771f5a8585638a2a15ee6c97cf7205d4208ed7c1df252d.
//
// Solidity: event RemoveLiquidity(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchRemoveLiquidity(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolRemoveLiquidity, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "RemoveLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolRemoveLiquidity)
				if err := _Curvelendingpool.contract.UnpackLog(event, "RemoveLiquidity", log); err != nil {
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

// ParseRemoveLiquidity is a log parse operation binding the contract event 0xa49d4cf02656aebf8c771f5a8585638a2a15ee6c97cf7205d4208ed7c1df252d.
//
// Solidity: event RemoveLiquidity(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseRemoveLiquidity(log types.Log) (*CurvelendingpoolRemoveLiquidity, error) {
	event := new(CurvelendingpoolRemoveLiquidity)
	if err := _Curvelendingpool.contract.UnpackLog(event, "RemoveLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolRemoveLiquidityImbalanceIterator is returned from FilterRemoveLiquidityImbalance and is used to iterate over the raw logs and unpacked data for RemoveLiquidityImbalance events raised by the Curvelendingpool contract.
type CurvelendingpoolRemoveLiquidityImbalanceIterator struct {
	Event *CurvelendingpoolRemoveLiquidityImbalance // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolRemoveLiquidityImbalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolRemoveLiquidityImbalance)
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
		it.Event = new(CurvelendingpoolRemoveLiquidityImbalance)
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
func (it *CurvelendingpoolRemoveLiquidityImbalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolRemoveLiquidityImbalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolRemoveLiquidityImbalance represents a RemoveLiquidityImbalance event raised by the Curvelendingpool contract.
type CurvelendingpoolRemoveLiquidityImbalance struct {
	Provider     common.Address
	TokenAmounts [3]*big.Int
	Fees         [3]*big.Int
	Invariant    *big.Int
	TokenSupply  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemoveLiquidityImbalance is a free log retrieval operation binding the contract event 0x173599dbf9c6ca6f7c3b590df07ae98a45d74ff54065505141e7de6c46a624c2.
//
// Solidity: event RemoveLiquidityImbalance(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 invariant, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterRemoveLiquidityImbalance(opts *bind.FilterOpts, provider []common.Address) (*CurvelendingpoolRemoveLiquidityImbalanceIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "RemoveLiquidityImbalance", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolRemoveLiquidityImbalanceIterator{contract: _Curvelendingpool.contract, event: "RemoveLiquidityImbalance", logs: logs, sub: sub}, nil
}

// WatchRemoveLiquidityImbalance is a free log subscription operation binding the contract event 0x173599dbf9c6ca6f7c3b590df07ae98a45d74ff54065505141e7de6c46a624c2.
//
// Solidity: event RemoveLiquidityImbalance(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 invariant, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchRemoveLiquidityImbalance(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolRemoveLiquidityImbalance, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "RemoveLiquidityImbalance", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolRemoveLiquidityImbalance)
				if err := _Curvelendingpool.contract.UnpackLog(event, "RemoveLiquidityImbalance", log); err != nil {
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

// ParseRemoveLiquidityImbalance is a log parse operation binding the contract event 0x173599dbf9c6ca6f7c3b590df07ae98a45d74ff54065505141e7de6c46a624c2.
//
// Solidity: event RemoveLiquidityImbalance(address indexed provider, uint256[3] token_amounts, uint256[3] fees, uint256 invariant, uint256 token_supply)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseRemoveLiquidityImbalance(log types.Log) (*CurvelendingpoolRemoveLiquidityImbalance, error) {
	event := new(CurvelendingpoolRemoveLiquidityImbalance)
	if err := _Curvelendingpool.contract.UnpackLog(event, "RemoveLiquidityImbalance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolRemoveLiquidityOneIterator is returned from FilterRemoveLiquidityOne and is used to iterate over the raw logs and unpacked data for RemoveLiquidityOne events raised by the Curvelendingpool contract.
type CurvelendingpoolRemoveLiquidityOneIterator struct {
	Event *CurvelendingpoolRemoveLiquidityOne // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolRemoveLiquidityOneIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolRemoveLiquidityOne)
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
		it.Event = new(CurvelendingpoolRemoveLiquidityOne)
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
func (it *CurvelendingpoolRemoveLiquidityOneIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolRemoveLiquidityOneIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolRemoveLiquidityOne represents a RemoveLiquidityOne event raised by the Curvelendingpool contract.
type CurvelendingpoolRemoveLiquidityOne struct {
	Provider    common.Address
	TokenAmount *big.Int
	CoinAmount  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveLiquidityOne is a free log retrieval operation binding the contract event 0x9e96dd3b997a2a257eec4df9bb6eaf626e206df5f543bd963682d143300be310.
//
// Solidity: event RemoveLiquidityOne(address indexed provider, uint256 token_amount, uint256 coin_amount)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterRemoveLiquidityOne(opts *bind.FilterOpts, provider []common.Address) (*CurvelendingpoolRemoveLiquidityOneIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "RemoveLiquidityOne", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolRemoveLiquidityOneIterator{contract: _Curvelendingpool.contract, event: "RemoveLiquidityOne", logs: logs, sub: sub}, nil
}

// WatchRemoveLiquidityOne is a free log subscription operation binding the contract event 0x9e96dd3b997a2a257eec4df9bb6eaf626e206df5f543bd963682d143300be310.
//
// Solidity: event RemoveLiquidityOne(address indexed provider, uint256 token_amount, uint256 coin_amount)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchRemoveLiquidityOne(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolRemoveLiquidityOne, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "RemoveLiquidityOne", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolRemoveLiquidityOne)
				if err := _Curvelendingpool.contract.UnpackLog(event, "RemoveLiquidityOne", log); err != nil {
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

// ParseRemoveLiquidityOne is a log parse operation binding the contract event 0x9e96dd3b997a2a257eec4df9bb6eaf626e206df5f543bd963682d143300be310.
//
// Solidity: event RemoveLiquidityOne(address indexed provider, uint256 token_amount, uint256 coin_amount)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseRemoveLiquidityOne(log types.Log) (*CurvelendingpoolRemoveLiquidityOne, error) {
	event := new(CurvelendingpoolRemoveLiquidityOne)
	if err := _Curvelendingpool.contract.UnpackLog(event, "RemoveLiquidityOne", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolStopRampAIterator is returned from FilterStopRampA and is used to iterate over the raw logs and unpacked data for StopRampA events raised by the Curvelendingpool contract.
type CurvelendingpoolStopRampAIterator struct {
	Event *CurvelendingpoolStopRampA // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolStopRampAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolStopRampA)
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
		it.Event = new(CurvelendingpoolStopRampA)
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
func (it *CurvelendingpoolStopRampAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolStopRampAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolStopRampA represents a StopRampA event raised by the Curvelendingpool contract.
type CurvelendingpoolStopRampA struct {
	A   *big.Int
	T   *big.Int
	Raw types.Log // Blockchain specific contextual infos
}

// FilterStopRampA is a free log retrieval operation binding the contract event 0x46e22fb3709ad289f62ce63d469248536dbc78d82b84a3d7e74ad606dc201938.
//
// Solidity: event StopRampA(uint256 A, uint256 t)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterStopRampA(opts *bind.FilterOpts) (*CurvelendingpoolStopRampAIterator, error) {

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "StopRampA")
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolStopRampAIterator{contract: _Curvelendingpool.contract, event: "StopRampA", logs: logs, sub: sub}, nil
}

// WatchStopRampA is a free log subscription operation binding the contract event 0x46e22fb3709ad289f62ce63d469248536dbc78d82b84a3d7e74ad606dc201938.
//
// Solidity: event StopRampA(uint256 A, uint256 t)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchStopRampA(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolStopRampA) (event.Subscription, error) {

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "StopRampA")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolStopRampA)
				if err := _Curvelendingpool.contract.UnpackLog(event, "StopRampA", log); err != nil {
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

// ParseStopRampA is a log parse operation binding the contract event 0x46e22fb3709ad289f62ce63d469248536dbc78d82b84a3d7e74ad606dc201938.
//
// Solidity: event StopRampA(uint256 A, uint256 t)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseStopRampA(log types.Log) (*CurvelendingpoolStopRampA, error) {
	event := new(CurvelendingpoolStopRampA)
	if err := _Curvelendingpool.contract.UnpackLog(event, "StopRampA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolTokenExchangeIterator is returned from FilterTokenExchange and is used to iterate over the raw logs and unpacked data for TokenExchange events raised by the Curvelendingpool contract.
type CurvelendingpoolTokenExchangeIterator struct {
	Event *CurvelendingpoolTokenExchange // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolTokenExchangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolTokenExchange)
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
		it.Event = new(CurvelendingpoolTokenExchange)
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
func (it *CurvelendingpoolTokenExchangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolTokenExchangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolTokenExchange represents a TokenExchange event raised by the Curvelendingpool contract.
type CurvelendingpoolTokenExchange struct {
	Buyer        common.Address
	SoldId       *big.Int
	TokensSold   *big.Int
	BoughtId     *big.Int
	TokensBought *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenExchange is a free log retrieval operation binding the contract event 0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140.
//
// Solidity: event TokenExchange(address indexed buyer, int128 sold_id, uint256 tokens_sold, int128 bought_id, uint256 tokens_bought)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterTokenExchange(opts *bind.FilterOpts, buyer []common.Address) (*CurvelendingpoolTokenExchangeIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "TokenExchange", buyerRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolTokenExchangeIterator{contract: _Curvelendingpool.contract, event: "TokenExchange", logs: logs, sub: sub}, nil
}

// WatchTokenExchange is a free log subscription operation binding the contract event 0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140.
//
// Solidity: event TokenExchange(address indexed buyer, int128 sold_id, uint256 tokens_sold, int128 bought_id, uint256 tokens_bought)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchTokenExchange(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolTokenExchange, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "TokenExchange", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolTokenExchange)
				if err := _Curvelendingpool.contract.UnpackLog(event, "TokenExchange", log); err != nil {
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

// ParseTokenExchange is a log parse operation binding the contract event 0x8b3e96f2b889fa771c53c981b40daf005f63f637f1869f707052d15a3dd97140.
//
// Solidity: event TokenExchange(address indexed buyer, int128 sold_id, uint256 tokens_sold, int128 bought_id, uint256 tokens_bought)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseTokenExchange(log types.Log) (*CurvelendingpoolTokenExchange, error) {
	event := new(CurvelendingpoolTokenExchange)
	if err := _Curvelendingpool.contract.UnpackLog(event, "TokenExchange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurvelendingpoolTokenExchangeUnderlyingIterator is returned from FilterTokenExchangeUnderlying and is used to iterate over the raw logs and unpacked data for TokenExchangeUnderlying events raised by the Curvelendingpool contract.
type CurvelendingpoolTokenExchangeUnderlyingIterator struct {
	Event *CurvelendingpoolTokenExchangeUnderlying // Event containing the contract specifics and raw log

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
func (it *CurvelendingpoolTokenExchangeUnderlyingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurvelendingpoolTokenExchangeUnderlying)
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
		it.Event = new(CurvelendingpoolTokenExchangeUnderlying)
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
func (it *CurvelendingpoolTokenExchangeUnderlyingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurvelendingpoolTokenExchangeUnderlyingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurvelendingpoolTokenExchangeUnderlying represents a TokenExchangeUnderlying event raised by the Curvelendingpool contract.
type CurvelendingpoolTokenExchangeUnderlying struct {
	Buyer        common.Address
	SoldId       *big.Int
	TokensSold   *big.Int
	BoughtId     *big.Int
	TokensBought *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenExchangeUnderlying is a free log retrieval operation binding the contract event 0xd013ca23e77a65003c2c659c5442c00c805371b7fc1ebd4c206c41d1536bd90b.
//
// Solidity: event TokenExchangeUnderlying(address indexed buyer, int128 sold_id, uint256 tokens_sold, int128 bought_id, uint256 tokens_bought)
func (_Curvelendingpool *CurvelendingpoolFilterer) FilterTokenExchangeUnderlying(opts *bind.FilterOpts, buyer []common.Address) (*CurvelendingpoolTokenExchangeUnderlyingIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.FilterLogs(opts, "TokenExchangeUnderlying", buyerRule)
	if err != nil {
		return nil, err
	}
	return &CurvelendingpoolTokenExchangeUnderlyingIterator{contract: _Curvelendingpool.contract, event: "TokenExchangeUnderlying", logs: logs, sub: sub}, nil
}

// WatchTokenExchangeUnderlying is a free log subscription operation binding the contract event 0xd013ca23e77a65003c2c659c5442c00c805371b7fc1ebd4c206c41d1536bd90b.
//
// Solidity: event TokenExchangeUnderlying(address indexed buyer, int128 sold_id, uint256 tokens_sold, int128 bought_id, uint256 tokens_bought)
func (_Curvelendingpool *CurvelendingpoolFilterer) WatchTokenExchangeUnderlying(opts *bind.WatchOpts, sink chan<- *CurvelendingpoolTokenExchangeUnderlying, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Curvelendingpool.contract.WatchLogs(opts, "TokenExchangeUnderlying", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurvelendingpoolTokenExchangeUnderlying)
				if err := _Curvelendingpool.contract.UnpackLog(event, "TokenExchangeUnderlying", log); err != nil {
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

// ParseTokenExchangeUnderlying is a log parse operation binding the contract event 0xd013ca23e77a65003c2c659c5442c00c805371b7fc1ebd4c206c41d1536bd90b.
//
// Solidity: event TokenExchangeUnderlying(address indexed buyer, int128 sold_id, uint256 tokens_sold, int128 bought_id, uint256 tokens_bought)
func (_Curvelendingpool *CurvelendingpoolFilterer) ParseTokenExchangeUnderlying(log types.Log) (*CurvelendingpoolTokenExchangeUnderlying, error) {
	event := new(CurvelendingpoolTokenExchangeUnderlying)
	if err := _Curvelendingpool.contract.UnpackLog(event, "TokenExchangeUnderlying", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
