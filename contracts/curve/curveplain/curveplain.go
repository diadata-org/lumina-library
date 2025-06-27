// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package curveplain

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

// CurveplainMetaData contains all meta data concerning the Curveplain contract.
var CurveplainMetaData = &bind.MetaData{
	ABI: "[{\"name\":\"TokenExchange\",\"inputs\":[{\"name\":\"buyer\",\"type\":\"address\",\"indexed\":true},{\"name\":\"sold_id\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"tokens_sold\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"bought_id\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"tokens_bought\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"AddLiquidity\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true},{\"name\":\"token_amounts\",\"type\":\"uint256[3]\",\"indexed\":false},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"token_supply\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RemoveLiquidity\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true},{\"name\":\"token_amounts\",\"type\":\"uint256[3]\",\"indexed\":false},{\"name\":\"token_supply\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RemoveLiquidityOne\",\"inputs\":[{\"name\":\"provider\",\"type\":\"address\",\"indexed\":true},{\"name\":\"token_amount\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"coin_index\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"coin_amount\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"CommitNewAdmin\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"indexed\":true},{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"NewAdmin\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"CommitNewParameters\",\"inputs\":[{\"name\":\"deadline\",\"type\":\"uint256\",\"indexed\":true},{\"name\":\"admin_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"mid_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"out_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"fee_gamma\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"price_threshold\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"adjustment_step\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"ma_half_time\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"NewParameters\",\"inputs\":[{\"name\":\"admin_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"mid_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"out_fee\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"fee_gamma\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"price_threshold\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"adjustment_step\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"ma_half_time\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"RampAgamma\",\"inputs\":[{\"name\":\"initial_A\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"future_A\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"initial_time\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"future_time\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"StopRampA\",\"inputs\":[{\"name\":\"current_A\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"current_gamma\",\"type\":\"uint256\",\"indexed\":false},{\"name\":\"time\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"name\":\"ClaimAdminFee\",\"inputs\":[{\"name\":\"admin\",\"type\":\"address\",\"indexed\":true},{\"name\":\"tokens\",\"type\":\"uint256\",\"indexed\":false}],\"anonymous\":false,\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"constructor\",\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"A\",\"type\":\"uint256\"},{\"name\":\"gamma\",\"type\":\"uint256\"},{\"name\":\"mid_fee\",\"type\":\"uint256\"},{\"name\":\"out_fee\",\"type\":\"uint256\"},{\"name\":\"price_threshold\",\"type\":\"uint256\"},{\"name\":\"fee_gamma\",\"type\":\"uint256\"},{\"name\":\"adjustment_step\",\"type\":\"uint256\"},{\"name\":\"admin_fee\",\"type\":\"uint256\"},{\"name\":\"ma_half_time\",\"type\":\"uint256\"},{\"name\":\"initial_prices\",\"type\":\"uint256[2]\"}],\"outputs\":[]},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"price_oracle\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3361},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"price_scale\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3391},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"last_prices\",\"inputs\":[{\"name\":\"k\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3421},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"token\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"gas\":468},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"coins\",\"inputs\":[{\"name\":\"i\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"gas\":582},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"A\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":937},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"gamma\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":12484},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"A_precise\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":929},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":21733},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"fee_calc\",\"inputs\":[{\"name\":\"xp\",\"type\":\"uint256[3]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":11156},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_virtual_price\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":11642},{\"stateMutability\":\"payable\",\"type\":\"function\",\"name\":\"exchange\",\"inputs\":[{\"name\":\"i\",\"type\":\"uint256\"},{\"name\":\"j\",\"type\":\"uint256\"},{\"name\":\"dx\",\"type\":\"uint256\"},{\"name\":\"min_dy\",\"type\":\"uint256\"}],\"outputs\":[]},{\"stateMutability\":\"payable\",\"type\":\"function\",\"name\":\"exchange\",\"inputs\":[{\"name\":\"i\",\"type\":\"uint256\"},{\"name\":\"j\",\"type\":\"uint256\"},{\"name\":\"dx\",\"type\":\"uint256\"},{\"name\":\"min_dy\",\"type\":\"uint256\"},{\"name\":\"use_eth\",\"type\":\"bool\"}],\"outputs\":[]},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"get_dy\",\"inputs\":[{\"name\":\"i\",\"type\":\"uint256\"},{\"name\":\"j\",\"type\":\"uint256\"},{\"name\":\"dx\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3152},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"calc_token_fee\",\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[3]\"},{\"name\":\"xp\",\"type\":\"uint256[3]\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":26702},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"add_liquidity\",\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[3]\"},{\"name\":\"min_mint_amount\",\"type\":\"uint256\"}],\"outputs\":[],\"gas\":682649},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"remove_liquidity\",\"inputs\":[{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"min_amounts\",\"type\":\"uint256[3]\"}],\"outputs\":[],\"gas\":234127},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"calc_token_amount\",\"inputs\":[{\"name\":\"amounts\",\"type\":\"uint256[3]\"},{\"name\":\"deposit\",\"type\":\"bool\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3459},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"calc_withdraw_one_coin\",\"inputs\":[{\"name\":\"token_amount\",\"type\":\"uint256\"},{\"name\":\"i\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":13939},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"remove_liquidity_one_coin\",\"inputs\":[{\"name\":\"token_amount\",\"type\":\"uint256\"},{\"name\":\"i\",\"type\":\"uint256\"},{\"name\":\"min_amount\",\"type\":\"uint256\"}],\"outputs\":[],\"gas\":592064},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"claim_admin_fees\",\"inputs\":[],\"outputs\":[],\"gas\":355740},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"ramp_A_gamma\",\"inputs\":[{\"name\":\"future_A\",\"type\":\"uint256\"},{\"name\":\"future_gamma\",\"type\":\"uint256\"},{\"name\":\"future_time\",\"type\":\"uint256\"}],\"outputs\":[],\"gas\":162031},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"stop_ramp_A_gamma\",\"inputs\":[],\"outputs\":[],\"gas\":157781},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"commit_new_parameters\",\"inputs\":[{\"name\":\"_new_mid_fee\",\"type\":\"uint256\"},{\"name\":\"_new_out_fee\",\"type\":\"uint256\"},{\"name\":\"_new_admin_fee\",\"type\":\"uint256\"},{\"name\":\"_new_fee_gamma\",\"type\":\"uint256\"},{\"name\":\"_new_price_threshold\",\"type\":\"uint256\"},{\"name\":\"_new_adjustment_step\",\"type\":\"uint256\"},{\"name\":\"_new_ma_half_time\",\"type\":\"uint256\"}],\"outputs\":[],\"gas\":306299},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"apply_new_parameters\",\"inputs\":[],\"outputs\":[],\"gas\":649370},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"revert_new_parameters\",\"inputs\":[],\"outputs\":[],\"gas\":23252},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"commit_transfer_ownership\",\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"outputs\":[],\"gas\":77290},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"apply_transfer_ownership\",\"inputs\":[],\"outputs\":[],\"gas\":65967},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"revert_transfer_ownership\",\"inputs\":[],\"outputs\":[],\"gas\":23342},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"kill_me\",\"inputs\":[],\"outputs\":[],\"gas\":40565},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"unkill_me\",\"inputs\":[],\"outputs\":[],\"gas\":23402},{\"stateMutability\":\"nonpayable\",\"type\":\"function\",\"name\":\"set_admin_fee_receiver\",\"inputs\":[{\"name\":\"_admin_fee_receiver\",\"type\":\"address\"}],\"outputs\":[],\"gas\":38535},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"last_prices_timestamp\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3408},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"initial_A_gamma\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3438},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_A_gamma\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3468},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"initial_A_gamma_time\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3498},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_A_gamma_time\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3528},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"price_threshold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3558},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_price_threshoold\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3588},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"fee_gamma\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3618},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_fee_gamma\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3648},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"adjustment_step\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3678},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_adjustment_step\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3708},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"ma_half_time\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3738},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_ma_half_time\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3768},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"mid_fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3798},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"out_fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3828},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"admin_fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3858},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_mid_fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3888},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_out_fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3918},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_admin_fee\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":3948},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"balances\",\"inputs\":[{\"name\":\"arg0\",\"type\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4087},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"D\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4008},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"gas\":4038},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"future_owner\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"gas\":4068},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"xcp_profit\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4098},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"xcp_profit_a\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4128},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"virtual_price\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4158},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"is_killed\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"gas\":4188},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"kill_deadline\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4218},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"transfer_ownership_deadline\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4248},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"admin_actions_deadline\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"gas\":4278},{\"stateMutability\":\"view\",\"type\":\"function\",\"name\":\"admin_fee_receiver\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"gas\":4308}]",
}

// CurveplainABI is the input ABI used to generate the binding from.
// Deprecated: Use CurveplainMetaData.ABI instead.
var CurveplainABI = CurveplainMetaData.ABI

// Curveplain is an auto generated Go binding around an Ethereum contract.
type Curveplain struct {
	CurveplainCaller     // Read-only binding to the contract
	CurveplainTransactor // Write-only binding to the contract
	CurveplainFilterer   // Log filterer for contract events
}

// CurveplainCaller is an auto generated read-only Go binding around an Ethereum contract.
type CurveplainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveplainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type CurveplainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveplainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type CurveplainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// CurveplainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type CurveplainSession struct {
	Contract     *Curveplain       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// CurveplainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type CurveplainCallerSession struct {
	Contract *CurveplainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// CurveplainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type CurveplainTransactorSession struct {
	Contract     *CurveplainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// CurveplainRaw is an auto generated low-level Go binding around an Ethereum contract.
type CurveplainRaw struct {
	Contract *Curveplain // Generic contract binding to access the raw methods on
}

// CurveplainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type CurveplainCallerRaw struct {
	Contract *CurveplainCaller // Generic read-only contract binding to access the raw methods on
}

// CurveplainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type CurveplainTransactorRaw struct {
	Contract *CurveplainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCurveplain creates a new instance of Curveplain, bound to a specific deployed contract.
func NewCurveplain(address common.Address, backend bind.ContractBackend) (*Curveplain, error) {
	contract, err := bindCurveplain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Curveplain{CurveplainCaller: CurveplainCaller{contract: contract}, CurveplainTransactor: CurveplainTransactor{contract: contract}, CurveplainFilterer: CurveplainFilterer{contract: contract}}, nil
}

// NewCurveplainCaller creates a new read-only instance of Curveplain, bound to a specific deployed contract.
func NewCurveplainCaller(address common.Address, caller bind.ContractCaller) (*CurveplainCaller, error) {
	contract, err := bindCurveplain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &CurveplainCaller{contract: contract}, nil
}

// NewCurveplainTransactor creates a new write-only instance of Curveplain, bound to a specific deployed contract.
func NewCurveplainTransactor(address common.Address, transactor bind.ContractTransactor) (*CurveplainTransactor, error) {
	contract, err := bindCurveplain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &CurveplainTransactor{contract: contract}, nil
}

// NewCurveplainFilterer creates a new log filterer instance of Curveplain, bound to a specific deployed contract.
func NewCurveplainFilterer(address common.Address, filterer bind.ContractFilterer) (*CurveplainFilterer, error) {
	contract, err := bindCurveplain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &CurveplainFilterer{contract: contract}, nil
}

// bindCurveplain binds a generic wrapper to an already deployed contract.
func bindCurveplain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := CurveplainMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Curveplain *CurveplainRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Curveplain.Contract.CurveplainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Curveplain *CurveplainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.Contract.CurveplainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Curveplain *CurveplainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Curveplain.Contract.CurveplainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Curveplain *CurveplainCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Curveplain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Curveplain *CurveplainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Curveplain *CurveplainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Curveplain.Contract.contract.Transact(opts, method, params...)
}

// A is a free data retrieval call binding the contract method 0xf446c1d0.
//
// Solidity: function A() view returns(uint256)
func (_Curveplain *CurveplainCaller) A(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "A")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// A is a free data retrieval call binding the contract method 0xf446c1d0.
//
// Solidity: function A() view returns(uint256)
func (_Curveplain *CurveplainSession) A() (*big.Int, error) {
	return _Curveplain.Contract.A(&_Curveplain.CallOpts)
}

// A is a free data retrieval call binding the contract method 0xf446c1d0.
//
// Solidity: function A() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) A() (*big.Int, error) {
	return _Curveplain.Contract.A(&_Curveplain.CallOpts)
}

// APrecise is a free data retrieval call binding the contract method 0x76a2f0f0.
//
// Solidity: function A_precise() view returns(uint256)
func (_Curveplain *CurveplainCaller) APrecise(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "A_precise")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// APrecise is a free data retrieval call binding the contract method 0x76a2f0f0.
//
// Solidity: function A_precise() view returns(uint256)
func (_Curveplain *CurveplainSession) APrecise() (*big.Int, error) {
	return _Curveplain.Contract.APrecise(&_Curveplain.CallOpts)
}

// APrecise is a free data retrieval call binding the contract method 0x76a2f0f0.
//
// Solidity: function A_precise() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) APrecise() (*big.Int, error) {
	return _Curveplain.Contract.APrecise(&_Curveplain.CallOpts)
}

// D is a free data retrieval call binding the contract method 0x0f529ba2.
//
// Solidity: function D() view returns(uint256)
func (_Curveplain *CurveplainCaller) D(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "D")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// D is a free data retrieval call binding the contract method 0x0f529ba2.
//
// Solidity: function D() view returns(uint256)
func (_Curveplain *CurveplainSession) D() (*big.Int, error) {
	return _Curveplain.Contract.D(&_Curveplain.CallOpts)
}

// D is a free data retrieval call binding the contract method 0x0f529ba2.
//
// Solidity: function D() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) D() (*big.Int, error) {
	return _Curveplain.Contract.D(&_Curveplain.CallOpts)
}

// AdjustmentStep is a free data retrieval call binding the contract method 0x083812e5.
//
// Solidity: function adjustment_step() view returns(uint256)
func (_Curveplain *CurveplainCaller) AdjustmentStep(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "adjustment_step")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AdjustmentStep is a free data retrieval call binding the contract method 0x083812e5.
//
// Solidity: function adjustment_step() view returns(uint256)
func (_Curveplain *CurveplainSession) AdjustmentStep() (*big.Int, error) {
	return _Curveplain.Contract.AdjustmentStep(&_Curveplain.CallOpts)
}

// AdjustmentStep is a free data retrieval call binding the contract method 0x083812e5.
//
// Solidity: function adjustment_step() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) AdjustmentStep() (*big.Int, error) {
	return _Curveplain.Contract.AdjustmentStep(&_Curveplain.CallOpts)
}

// AdminActionsDeadline is a free data retrieval call binding the contract method 0x405e28f8.
//
// Solidity: function admin_actions_deadline() view returns(uint256)
func (_Curveplain *CurveplainCaller) AdminActionsDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "admin_actions_deadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AdminActionsDeadline is a free data retrieval call binding the contract method 0x405e28f8.
//
// Solidity: function admin_actions_deadline() view returns(uint256)
func (_Curveplain *CurveplainSession) AdminActionsDeadline() (*big.Int, error) {
	return _Curveplain.Contract.AdminActionsDeadline(&_Curveplain.CallOpts)
}

// AdminActionsDeadline is a free data retrieval call binding the contract method 0x405e28f8.
//
// Solidity: function admin_actions_deadline() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) AdminActionsDeadline() (*big.Int, error) {
	return _Curveplain.Contract.AdminActionsDeadline(&_Curveplain.CallOpts)
}

// AdminFee is a free data retrieval call binding the contract method 0xfee3f7f9.
//
// Solidity: function admin_fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) AdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "admin_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// AdminFee is a free data retrieval call binding the contract method 0xfee3f7f9.
//
// Solidity: function admin_fee() view returns(uint256)
func (_Curveplain *CurveplainSession) AdminFee() (*big.Int, error) {
	return _Curveplain.Contract.AdminFee(&_Curveplain.CallOpts)
}

// AdminFee is a free data retrieval call binding the contract method 0xfee3f7f9.
//
// Solidity: function admin_fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) AdminFee() (*big.Int, error) {
	return _Curveplain.Contract.AdminFee(&_Curveplain.CallOpts)
}

// AdminFeeReceiver is a free data retrieval call binding the contract method 0x6e42e4d2.
//
// Solidity: function admin_fee_receiver() view returns(address)
func (_Curveplain *CurveplainCaller) AdminFeeReceiver(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "admin_fee_receiver")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AdminFeeReceiver is a free data retrieval call binding the contract method 0x6e42e4d2.
//
// Solidity: function admin_fee_receiver() view returns(address)
func (_Curveplain *CurveplainSession) AdminFeeReceiver() (common.Address, error) {
	return _Curveplain.Contract.AdminFeeReceiver(&_Curveplain.CallOpts)
}

// AdminFeeReceiver is a free data retrieval call binding the contract method 0x6e42e4d2.
//
// Solidity: function admin_fee_receiver() view returns(address)
func (_Curveplain *CurveplainCallerSession) AdminFeeReceiver() (common.Address, error) {
	return _Curveplain.Contract.AdminFeeReceiver(&_Curveplain.CallOpts)
}

// Balances is a free data retrieval call binding the contract method 0x4903b0d1.
//
// Solidity: function balances(uint256 arg0) view returns(uint256)
func (_Curveplain *CurveplainCaller) Balances(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "balances", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Balances is a free data retrieval call binding the contract method 0x4903b0d1.
//
// Solidity: function balances(uint256 arg0) view returns(uint256)
func (_Curveplain *CurveplainSession) Balances(arg0 *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.Balances(&_Curveplain.CallOpts, arg0)
}

// Balances is a free data retrieval call binding the contract method 0x4903b0d1.
//
// Solidity: function balances(uint256 arg0) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) Balances(arg0 *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.Balances(&_Curveplain.CallOpts, arg0)
}

// CalcTokenAmount is a free data retrieval call binding the contract method 0x3883e119.
//
// Solidity: function calc_token_amount(uint256[3] amounts, bool deposit) view returns(uint256)
func (_Curveplain *CurveplainCaller) CalcTokenAmount(opts *bind.CallOpts, amounts [3]*big.Int, deposit bool) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "calc_token_amount", amounts, deposit)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcTokenAmount is a free data retrieval call binding the contract method 0x3883e119.
//
// Solidity: function calc_token_amount(uint256[3] amounts, bool deposit) view returns(uint256)
func (_Curveplain *CurveplainSession) CalcTokenAmount(amounts [3]*big.Int, deposit bool) (*big.Int, error) {
	return _Curveplain.Contract.CalcTokenAmount(&_Curveplain.CallOpts, amounts, deposit)
}

// CalcTokenAmount is a free data retrieval call binding the contract method 0x3883e119.
//
// Solidity: function calc_token_amount(uint256[3] amounts, bool deposit) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) CalcTokenAmount(amounts [3]*big.Int, deposit bool) (*big.Int, error) {
	return _Curveplain.Contract.CalcTokenAmount(&_Curveplain.CallOpts, amounts, deposit)
}

// CalcTokenFee is a free data retrieval call binding the contract method 0xcde699fa.
//
// Solidity: function calc_token_fee(uint256[3] amounts, uint256[3] xp) view returns(uint256)
func (_Curveplain *CurveplainCaller) CalcTokenFee(opts *bind.CallOpts, amounts [3]*big.Int, xp [3]*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "calc_token_fee", amounts, xp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcTokenFee is a free data retrieval call binding the contract method 0xcde699fa.
//
// Solidity: function calc_token_fee(uint256[3] amounts, uint256[3] xp) view returns(uint256)
func (_Curveplain *CurveplainSession) CalcTokenFee(amounts [3]*big.Int, xp [3]*big.Int) (*big.Int, error) {
	return _Curveplain.Contract.CalcTokenFee(&_Curveplain.CallOpts, amounts, xp)
}

// CalcTokenFee is a free data retrieval call binding the contract method 0xcde699fa.
//
// Solidity: function calc_token_fee(uint256[3] amounts, uint256[3] xp) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) CalcTokenFee(amounts [3]*big.Int, xp [3]*big.Int) (*big.Int, error) {
	return _Curveplain.Contract.CalcTokenFee(&_Curveplain.CallOpts, amounts, xp)
}

// CalcWithdrawOneCoin is a free data retrieval call binding the contract method 0x4fb08c5e.
//
// Solidity: function calc_withdraw_one_coin(uint256 token_amount, uint256 i) view returns(uint256)
func (_Curveplain *CurveplainCaller) CalcWithdrawOneCoin(opts *bind.CallOpts, token_amount *big.Int, i *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "calc_withdraw_one_coin", token_amount, i)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalcWithdrawOneCoin is a free data retrieval call binding the contract method 0x4fb08c5e.
//
// Solidity: function calc_withdraw_one_coin(uint256 token_amount, uint256 i) view returns(uint256)
func (_Curveplain *CurveplainSession) CalcWithdrawOneCoin(token_amount *big.Int, i *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.CalcWithdrawOneCoin(&_Curveplain.CallOpts, token_amount, i)
}

// CalcWithdrawOneCoin is a free data retrieval call binding the contract method 0x4fb08c5e.
//
// Solidity: function calc_withdraw_one_coin(uint256 token_amount, uint256 i) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) CalcWithdrawOneCoin(token_amount *big.Int, i *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.CalcWithdrawOneCoin(&_Curveplain.CallOpts, token_amount, i)
}

// Coins is a free data retrieval call binding the contract method 0xc6610657.
//
// Solidity: function coins(uint256 i) view returns(address)
func (_Curveplain *CurveplainCaller) Coins(opts *bind.CallOpts, i *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "coins", i)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Coins is a free data retrieval call binding the contract method 0xc6610657.
//
// Solidity: function coins(uint256 i) view returns(address)
func (_Curveplain *CurveplainSession) Coins(i *big.Int) (common.Address, error) {
	return _Curveplain.Contract.Coins(&_Curveplain.CallOpts, i)
}

// Coins is a free data retrieval call binding the contract method 0xc6610657.
//
// Solidity: function coins(uint256 i) view returns(address)
func (_Curveplain *CurveplainCallerSession) Coins(i *big.Int) (common.Address, error) {
	return _Curveplain.Contract.Coins(&_Curveplain.CallOpts, i)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) Fee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Curveplain *CurveplainSession) Fee() (*big.Int, error) {
	return _Curveplain.Contract.Fee(&_Curveplain.CallOpts)
}

// Fee is a free data retrieval call binding the contract method 0xddca3f43.
//
// Solidity: function fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) Fee() (*big.Int, error) {
	return _Curveplain.Contract.Fee(&_Curveplain.CallOpts)
}

// FeeCalc is a free data retrieval call binding the contract method 0x572e5625.
//
// Solidity: function fee_calc(uint256[3] xp) view returns(uint256)
func (_Curveplain *CurveplainCaller) FeeCalc(opts *bind.CallOpts, xp [3]*big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "fee_calc", xp)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeCalc is a free data retrieval call binding the contract method 0x572e5625.
//
// Solidity: function fee_calc(uint256[3] xp) view returns(uint256)
func (_Curveplain *CurveplainSession) FeeCalc(xp [3]*big.Int) (*big.Int, error) {
	return _Curveplain.Contract.FeeCalc(&_Curveplain.CallOpts, xp)
}

// FeeCalc is a free data retrieval call binding the contract method 0x572e5625.
//
// Solidity: function fee_calc(uint256[3] xp) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FeeCalc(xp [3]*big.Int) (*big.Int, error) {
	return _Curveplain.Contract.FeeCalc(&_Curveplain.CallOpts, xp)
}

// FeeGamma is a free data retrieval call binding the contract method 0x72d4f0e2.
//
// Solidity: function fee_gamma() view returns(uint256)
func (_Curveplain *CurveplainCaller) FeeGamma(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "fee_gamma")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FeeGamma is a free data retrieval call binding the contract method 0x72d4f0e2.
//
// Solidity: function fee_gamma() view returns(uint256)
func (_Curveplain *CurveplainSession) FeeGamma() (*big.Int, error) {
	return _Curveplain.Contract.FeeGamma(&_Curveplain.CallOpts)
}

// FeeGamma is a free data retrieval call binding the contract method 0x72d4f0e2.
//
// Solidity: function fee_gamma() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FeeGamma() (*big.Int, error) {
	return _Curveplain.Contract.FeeGamma(&_Curveplain.CallOpts)
}

// FutureAGamma is a free data retrieval call binding the contract method 0xf30cfad5.
//
// Solidity: function future_A_gamma() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureAGamma(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_A_gamma")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureAGamma is a free data retrieval call binding the contract method 0xf30cfad5.
//
// Solidity: function future_A_gamma() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureAGamma() (*big.Int, error) {
	return _Curveplain.Contract.FutureAGamma(&_Curveplain.CallOpts)
}

// FutureAGamma is a free data retrieval call binding the contract method 0xf30cfad5.
//
// Solidity: function future_A_gamma() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureAGamma() (*big.Int, error) {
	return _Curveplain.Contract.FutureAGamma(&_Curveplain.CallOpts)
}

// FutureAGammaTime is a free data retrieval call binding the contract method 0xf9ed9597.
//
// Solidity: function future_A_gamma_time() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureAGammaTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_A_gamma_time")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureAGammaTime is a free data retrieval call binding the contract method 0xf9ed9597.
//
// Solidity: function future_A_gamma_time() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureAGammaTime() (*big.Int, error) {
	return _Curveplain.Contract.FutureAGammaTime(&_Curveplain.CallOpts)
}

// FutureAGammaTime is a free data retrieval call binding the contract method 0xf9ed9597.
//
// Solidity: function future_A_gamma_time() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureAGammaTime() (*big.Int, error) {
	return _Curveplain.Contract.FutureAGammaTime(&_Curveplain.CallOpts)
}

// FutureAdjustmentStep is a free data retrieval call binding the contract method 0x4ea12c7d.
//
// Solidity: function future_adjustment_step() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureAdjustmentStep(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_adjustment_step")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureAdjustmentStep is a free data retrieval call binding the contract method 0x4ea12c7d.
//
// Solidity: function future_adjustment_step() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureAdjustmentStep() (*big.Int, error) {
	return _Curveplain.Contract.FutureAdjustmentStep(&_Curveplain.CallOpts)
}

// FutureAdjustmentStep is a free data retrieval call binding the contract method 0x4ea12c7d.
//
// Solidity: function future_adjustment_step() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureAdjustmentStep() (*big.Int, error) {
	return _Curveplain.Contract.FutureAdjustmentStep(&_Curveplain.CallOpts)
}

// FutureAdminFee is a free data retrieval call binding the contract method 0xe3824462.
//
// Solidity: function future_admin_fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureAdminFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_admin_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureAdminFee is a free data retrieval call binding the contract method 0xe3824462.
//
// Solidity: function future_admin_fee() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureAdminFee() (*big.Int, error) {
	return _Curveplain.Contract.FutureAdminFee(&_Curveplain.CallOpts)
}

// FutureAdminFee is a free data retrieval call binding the contract method 0xe3824462.
//
// Solidity: function future_admin_fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureAdminFee() (*big.Int, error) {
	return _Curveplain.Contract.FutureAdminFee(&_Curveplain.CallOpts)
}

// FutureFeeGamma is a free data retrieval call binding the contract method 0xd7c3dcbe.
//
// Solidity: function future_fee_gamma() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureFeeGamma(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_fee_gamma")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureFeeGamma is a free data retrieval call binding the contract method 0xd7c3dcbe.
//
// Solidity: function future_fee_gamma() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureFeeGamma() (*big.Int, error) {
	return _Curveplain.Contract.FutureFeeGamma(&_Curveplain.CallOpts)
}

// FutureFeeGamma is a free data retrieval call binding the contract method 0xd7c3dcbe.
//
// Solidity: function future_fee_gamma() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureFeeGamma() (*big.Int, error) {
	return _Curveplain.Contract.FutureFeeGamma(&_Curveplain.CallOpts)
}

// FutureMaHalfTime is a free data retrieval call binding the contract method 0x0c5e23d4.
//
// Solidity: function future_ma_half_time() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureMaHalfTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_ma_half_time")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureMaHalfTime is a free data retrieval call binding the contract method 0x0c5e23d4.
//
// Solidity: function future_ma_half_time() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureMaHalfTime() (*big.Int, error) {
	return _Curveplain.Contract.FutureMaHalfTime(&_Curveplain.CallOpts)
}

// FutureMaHalfTime is a free data retrieval call binding the contract method 0x0c5e23d4.
//
// Solidity: function future_ma_half_time() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureMaHalfTime() (*big.Int, error) {
	return _Curveplain.Contract.FutureMaHalfTime(&_Curveplain.CallOpts)
}

// FutureMidFee is a free data retrieval call binding the contract method 0x7cf9aedc.
//
// Solidity: function future_mid_fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureMidFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_mid_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureMidFee is a free data retrieval call binding the contract method 0x7cf9aedc.
//
// Solidity: function future_mid_fee() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureMidFee() (*big.Int, error) {
	return _Curveplain.Contract.FutureMidFee(&_Curveplain.CallOpts)
}

// FutureMidFee is a free data retrieval call binding the contract method 0x7cf9aedc.
//
// Solidity: function future_mid_fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureMidFee() (*big.Int, error) {
	return _Curveplain.Contract.FutureMidFee(&_Curveplain.CallOpts)
}

// FutureOutFee is a free data retrieval call binding the contract method 0x7d1b060c.
//
// Solidity: function future_out_fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) FutureOutFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_out_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FutureOutFee is a free data retrieval call binding the contract method 0x7d1b060c.
//
// Solidity: function future_out_fee() view returns(uint256)
func (_Curveplain *CurveplainSession) FutureOutFee() (*big.Int, error) {
	return _Curveplain.Contract.FutureOutFee(&_Curveplain.CallOpts)
}

// FutureOutFee is a free data retrieval call binding the contract method 0x7d1b060c.
//
// Solidity: function future_out_fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FutureOutFee() (*big.Int, error) {
	return _Curveplain.Contract.FutureOutFee(&_Curveplain.CallOpts)
}

// FutureOwner is a free data retrieval call binding the contract method 0x1ec0cdc1.
//
// Solidity: function future_owner() view returns(address)
func (_Curveplain *CurveplainCaller) FutureOwner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FutureOwner is a free data retrieval call binding the contract method 0x1ec0cdc1.
//
// Solidity: function future_owner() view returns(address)
func (_Curveplain *CurveplainSession) FutureOwner() (common.Address, error) {
	return _Curveplain.Contract.FutureOwner(&_Curveplain.CallOpts)
}

// FutureOwner is a free data retrieval call binding the contract method 0x1ec0cdc1.
//
// Solidity: function future_owner() view returns(address)
func (_Curveplain *CurveplainCallerSession) FutureOwner() (common.Address, error) {
	return _Curveplain.Contract.FutureOwner(&_Curveplain.CallOpts)
}

// FuturePriceThreshoold is a free data retrieval call binding the contract method 0x73fb9711.
//
// Solidity: function future_price_threshoold() view returns(uint256)
func (_Curveplain *CurveplainCaller) FuturePriceThreshoold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "future_price_threshoold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// FuturePriceThreshoold is a free data retrieval call binding the contract method 0x73fb9711.
//
// Solidity: function future_price_threshoold() view returns(uint256)
func (_Curveplain *CurveplainSession) FuturePriceThreshoold() (*big.Int, error) {
	return _Curveplain.Contract.FuturePriceThreshoold(&_Curveplain.CallOpts)
}

// FuturePriceThreshoold is a free data retrieval call binding the contract method 0x73fb9711.
//
// Solidity: function future_price_threshoold() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) FuturePriceThreshoold() (*big.Int, error) {
	return _Curveplain.Contract.FuturePriceThreshoold(&_Curveplain.CallOpts)
}

// Gamma is a free data retrieval call binding the contract method 0xb1373929.
//
// Solidity: function gamma() view returns(uint256)
func (_Curveplain *CurveplainCaller) Gamma(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "gamma")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Gamma is a free data retrieval call binding the contract method 0xb1373929.
//
// Solidity: function gamma() view returns(uint256)
func (_Curveplain *CurveplainSession) Gamma() (*big.Int, error) {
	return _Curveplain.Contract.Gamma(&_Curveplain.CallOpts)
}

// Gamma is a free data retrieval call binding the contract method 0xb1373929.
//
// Solidity: function gamma() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) Gamma() (*big.Int, error) {
	return _Curveplain.Contract.Gamma(&_Curveplain.CallOpts)
}

// GetDy is a free data retrieval call binding the contract method 0x556d6e9f.
//
// Solidity: function get_dy(uint256 i, uint256 j, uint256 dx) view returns(uint256)
func (_Curveplain *CurveplainCaller) GetDy(opts *bind.CallOpts, i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "get_dy", i, j, dx)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetDy is a free data retrieval call binding the contract method 0x556d6e9f.
//
// Solidity: function get_dy(uint256 i, uint256 j, uint256 dx) view returns(uint256)
func (_Curveplain *CurveplainSession) GetDy(i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.GetDy(&_Curveplain.CallOpts, i, j, dx)
}

// GetDy is a free data retrieval call binding the contract method 0x556d6e9f.
//
// Solidity: function get_dy(uint256 i, uint256 j, uint256 dx) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) GetDy(i *big.Int, j *big.Int, dx *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.GetDy(&_Curveplain.CallOpts, i, j, dx)
}

// GetVirtualPrice is a free data retrieval call binding the contract method 0xbb7b8b80.
//
// Solidity: function get_virtual_price() view returns(uint256)
func (_Curveplain *CurveplainCaller) GetVirtualPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "get_virtual_price")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetVirtualPrice is a free data retrieval call binding the contract method 0xbb7b8b80.
//
// Solidity: function get_virtual_price() view returns(uint256)
func (_Curveplain *CurveplainSession) GetVirtualPrice() (*big.Int, error) {
	return _Curveplain.Contract.GetVirtualPrice(&_Curveplain.CallOpts)
}

// GetVirtualPrice is a free data retrieval call binding the contract method 0xbb7b8b80.
//
// Solidity: function get_virtual_price() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) GetVirtualPrice() (*big.Int, error) {
	return _Curveplain.Contract.GetVirtualPrice(&_Curveplain.CallOpts)
}

// InitialAGamma is a free data retrieval call binding the contract method 0x204fe3d5.
//
// Solidity: function initial_A_gamma() view returns(uint256)
func (_Curveplain *CurveplainCaller) InitialAGamma(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "initial_A_gamma")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialAGamma is a free data retrieval call binding the contract method 0x204fe3d5.
//
// Solidity: function initial_A_gamma() view returns(uint256)
func (_Curveplain *CurveplainSession) InitialAGamma() (*big.Int, error) {
	return _Curveplain.Contract.InitialAGamma(&_Curveplain.CallOpts)
}

// InitialAGamma is a free data retrieval call binding the contract method 0x204fe3d5.
//
// Solidity: function initial_A_gamma() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) InitialAGamma() (*big.Int, error) {
	return _Curveplain.Contract.InitialAGamma(&_Curveplain.CallOpts)
}

// InitialAGammaTime is a free data retrieval call binding the contract method 0xe89876ff.
//
// Solidity: function initial_A_gamma_time() view returns(uint256)
func (_Curveplain *CurveplainCaller) InitialAGammaTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "initial_A_gamma_time")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// InitialAGammaTime is a free data retrieval call binding the contract method 0xe89876ff.
//
// Solidity: function initial_A_gamma_time() view returns(uint256)
func (_Curveplain *CurveplainSession) InitialAGammaTime() (*big.Int, error) {
	return _Curveplain.Contract.InitialAGammaTime(&_Curveplain.CallOpts)
}

// InitialAGammaTime is a free data retrieval call binding the contract method 0xe89876ff.
//
// Solidity: function initial_A_gamma_time() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) InitialAGammaTime() (*big.Int, error) {
	return _Curveplain.Contract.InitialAGammaTime(&_Curveplain.CallOpts)
}

// IsKilled is a free data retrieval call binding the contract method 0x9c868ac0.
//
// Solidity: function is_killed() view returns(bool)
func (_Curveplain *CurveplainCaller) IsKilled(opts *bind.CallOpts) (bool, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "is_killed")

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsKilled is a free data retrieval call binding the contract method 0x9c868ac0.
//
// Solidity: function is_killed() view returns(bool)
func (_Curveplain *CurveplainSession) IsKilled() (bool, error) {
	return _Curveplain.Contract.IsKilled(&_Curveplain.CallOpts)
}

// IsKilled is a free data retrieval call binding the contract method 0x9c868ac0.
//
// Solidity: function is_killed() view returns(bool)
func (_Curveplain *CurveplainCallerSession) IsKilled() (bool, error) {
	return _Curveplain.Contract.IsKilled(&_Curveplain.CallOpts)
}

// KillDeadline is a free data retrieval call binding the contract method 0x2a426896.
//
// Solidity: function kill_deadline() view returns(uint256)
func (_Curveplain *CurveplainCaller) KillDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "kill_deadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// KillDeadline is a free data retrieval call binding the contract method 0x2a426896.
//
// Solidity: function kill_deadline() view returns(uint256)
func (_Curveplain *CurveplainSession) KillDeadline() (*big.Int, error) {
	return _Curveplain.Contract.KillDeadline(&_Curveplain.CallOpts)
}

// KillDeadline is a free data retrieval call binding the contract method 0x2a426896.
//
// Solidity: function kill_deadline() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) KillDeadline() (*big.Int, error) {
	return _Curveplain.Contract.KillDeadline(&_Curveplain.CallOpts)
}

// LastPrices is a free data retrieval call binding the contract method 0x59189017.
//
// Solidity: function last_prices(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainCaller) LastPrices(opts *bind.CallOpts, k *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "last_prices", k)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastPrices is a free data retrieval call binding the contract method 0x59189017.
//
// Solidity: function last_prices(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainSession) LastPrices(k *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.LastPrices(&_Curveplain.CallOpts, k)
}

// LastPrices is a free data retrieval call binding the contract method 0x59189017.
//
// Solidity: function last_prices(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) LastPrices(k *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.LastPrices(&_Curveplain.CallOpts, k)
}

// LastPricesTimestamp is a free data retrieval call binding the contract method 0x6112c747.
//
// Solidity: function last_prices_timestamp() view returns(uint256)
func (_Curveplain *CurveplainCaller) LastPricesTimestamp(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "last_prices_timestamp")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// LastPricesTimestamp is a free data retrieval call binding the contract method 0x6112c747.
//
// Solidity: function last_prices_timestamp() view returns(uint256)
func (_Curveplain *CurveplainSession) LastPricesTimestamp() (*big.Int, error) {
	return _Curveplain.Contract.LastPricesTimestamp(&_Curveplain.CallOpts)
}

// LastPricesTimestamp is a free data retrieval call binding the contract method 0x6112c747.
//
// Solidity: function last_prices_timestamp() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) LastPricesTimestamp() (*big.Int, error) {
	return _Curveplain.Contract.LastPricesTimestamp(&_Curveplain.CallOpts)
}

// MaHalfTime is a free data retrieval call binding the contract method 0x662b6274.
//
// Solidity: function ma_half_time() view returns(uint256)
func (_Curveplain *CurveplainCaller) MaHalfTime(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "ma_half_time")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MaHalfTime is a free data retrieval call binding the contract method 0x662b6274.
//
// Solidity: function ma_half_time() view returns(uint256)
func (_Curveplain *CurveplainSession) MaHalfTime() (*big.Int, error) {
	return _Curveplain.Contract.MaHalfTime(&_Curveplain.CallOpts)
}

// MaHalfTime is a free data retrieval call binding the contract method 0x662b6274.
//
// Solidity: function ma_half_time() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) MaHalfTime() (*big.Int, error) {
	return _Curveplain.Contract.MaHalfTime(&_Curveplain.CallOpts)
}

// MidFee is a free data retrieval call binding the contract method 0x92526c0c.
//
// Solidity: function mid_fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) MidFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "mid_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MidFee is a free data retrieval call binding the contract method 0x92526c0c.
//
// Solidity: function mid_fee() view returns(uint256)
func (_Curveplain *CurveplainSession) MidFee() (*big.Int, error) {
	return _Curveplain.Contract.MidFee(&_Curveplain.CallOpts)
}

// MidFee is a free data retrieval call binding the contract method 0x92526c0c.
//
// Solidity: function mid_fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) MidFee() (*big.Int, error) {
	return _Curveplain.Contract.MidFee(&_Curveplain.CallOpts)
}

// OutFee is a free data retrieval call binding the contract method 0xee8de675.
//
// Solidity: function out_fee() view returns(uint256)
func (_Curveplain *CurveplainCaller) OutFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "out_fee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// OutFee is a free data retrieval call binding the contract method 0xee8de675.
//
// Solidity: function out_fee() view returns(uint256)
func (_Curveplain *CurveplainSession) OutFee() (*big.Int, error) {
	return _Curveplain.Contract.OutFee(&_Curveplain.CallOpts)
}

// OutFee is a free data retrieval call binding the contract method 0xee8de675.
//
// Solidity: function out_fee() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) OutFee() (*big.Int, error) {
	return _Curveplain.Contract.OutFee(&_Curveplain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Curveplain *CurveplainCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Curveplain *CurveplainSession) Owner() (common.Address, error) {
	return _Curveplain.Contract.Owner(&_Curveplain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Curveplain *CurveplainCallerSession) Owner() (common.Address, error) {
	return _Curveplain.Contract.Owner(&_Curveplain.CallOpts)
}

// PriceOracle is a free data retrieval call binding the contract method 0x68727653.
//
// Solidity: function price_oracle(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainCaller) PriceOracle(opts *bind.CallOpts, k *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "price_oracle", k)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PriceOracle is a free data retrieval call binding the contract method 0x68727653.
//
// Solidity: function price_oracle(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainSession) PriceOracle(k *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.PriceOracle(&_Curveplain.CallOpts, k)
}

// PriceOracle is a free data retrieval call binding the contract method 0x68727653.
//
// Solidity: function price_oracle(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) PriceOracle(k *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.PriceOracle(&_Curveplain.CallOpts, k)
}

// PriceScale is a free data retrieval call binding the contract method 0xa3f7cdd5.
//
// Solidity: function price_scale(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainCaller) PriceScale(opts *bind.CallOpts, k *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "price_scale", k)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PriceScale is a free data retrieval call binding the contract method 0xa3f7cdd5.
//
// Solidity: function price_scale(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainSession) PriceScale(k *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.PriceScale(&_Curveplain.CallOpts, k)
}

// PriceScale is a free data retrieval call binding the contract method 0xa3f7cdd5.
//
// Solidity: function price_scale(uint256 k) view returns(uint256)
func (_Curveplain *CurveplainCallerSession) PriceScale(k *big.Int) (*big.Int, error) {
	return _Curveplain.Contract.PriceScale(&_Curveplain.CallOpts, k)
}

// PriceThreshold is a free data retrieval call binding the contract method 0x03b736e0.
//
// Solidity: function price_threshold() view returns(uint256)
func (_Curveplain *CurveplainCaller) PriceThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "price_threshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PriceThreshold is a free data retrieval call binding the contract method 0x03b736e0.
//
// Solidity: function price_threshold() view returns(uint256)
func (_Curveplain *CurveplainSession) PriceThreshold() (*big.Int, error) {
	return _Curveplain.Contract.PriceThreshold(&_Curveplain.CallOpts)
}

// PriceThreshold is a free data retrieval call binding the contract method 0x03b736e0.
//
// Solidity: function price_threshold() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) PriceThreshold() (*big.Int, error) {
	return _Curveplain.Contract.PriceThreshold(&_Curveplain.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Curveplain *CurveplainCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "token")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Curveplain *CurveplainSession) Token() (common.Address, error) {
	return _Curveplain.Contract.Token(&_Curveplain.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() view returns(address)
func (_Curveplain *CurveplainCallerSession) Token() (common.Address, error) {
	return _Curveplain.Contract.Token(&_Curveplain.CallOpts)
}

// TransferOwnershipDeadline is a free data retrieval call binding the contract method 0xe0a0b586.
//
// Solidity: function transfer_ownership_deadline() view returns(uint256)
func (_Curveplain *CurveplainCaller) TransferOwnershipDeadline(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "transfer_ownership_deadline")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TransferOwnershipDeadline is a free data retrieval call binding the contract method 0xe0a0b586.
//
// Solidity: function transfer_ownership_deadline() view returns(uint256)
func (_Curveplain *CurveplainSession) TransferOwnershipDeadline() (*big.Int, error) {
	return _Curveplain.Contract.TransferOwnershipDeadline(&_Curveplain.CallOpts)
}

// TransferOwnershipDeadline is a free data retrieval call binding the contract method 0xe0a0b586.
//
// Solidity: function transfer_ownership_deadline() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) TransferOwnershipDeadline() (*big.Int, error) {
	return _Curveplain.Contract.TransferOwnershipDeadline(&_Curveplain.CallOpts)
}

// VirtualPrice is a free data retrieval call binding the contract method 0x0c46b72a.
//
// Solidity: function virtual_price() view returns(uint256)
func (_Curveplain *CurveplainCaller) VirtualPrice(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "virtual_price")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// VirtualPrice is a free data retrieval call binding the contract method 0x0c46b72a.
//
// Solidity: function virtual_price() view returns(uint256)
func (_Curveplain *CurveplainSession) VirtualPrice() (*big.Int, error) {
	return _Curveplain.Contract.VirtualPrice(&_Curveplain.CallOpts)
}

// VirtualPrice is a free data retrieval call binding the contract method 0x0c46b72a.
//
// Solidity: function virtual_price() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) VirtualPrice() (*big.Int, error) {
	return _Curveplain.Contract.VirtualPrice(&_Curveplain.CallOpts)
}

// XcpProfit is a free data retrieval call binding the contract method 0x7ba1a74d.
//
// Solidity: function xcp_profit() view returns(uint256)
func (_Curveplain *CurveplainCaller) XcpProfit(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "xcp_profit")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// XcpProfit is a free data retrieval call binding the contract method 0x7ba1a74d.
//
// Solidity: function xcp_profit() view returns(uint256)
func (_Curveplain *CurveplainSession) XcpProfit() (*big.Int, error) {
	return _Curveplain.Contract.XcpProfit(&_Curveplain.CallOpts)
}

// XcpProfit is a free data retrieval call binding the contract method 0x7ba1a74d.
//
// Solidity: function xcp_profit() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) XcpProfit() (*big.Int, error) {
	return _Curveplain.Contract.XcpProfit(&_Curveplain.CallOpts)
}

// XcpProfitA is a free data retrieval call binding the contract method 0x0b7b594b.
//
// Solidity: function xcp_profit_a() view returns(uint256)
func (_Curveplain *CurveplainCaller) XcpProfitA(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Curveplain.contract.Call(opts, &out, "xcp_profit_a")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// XcpProfitA is a free data retrieval call binding the contract method 0x0b7b594b.
//
// Solidity: function xcp_profit_a() view returns(uint256)
func (_Curveplain *CurveplainSession) XcpProfitA() (*big.Int, error) {
	return _Curveplain.Contract.XcpProfitA(&_Curveplain.CallOpts)
}

// XcpProfitA is a free data retrieval call binding the contract method 0x0b7b594b.
//
// Solidity: function xcp_profit_a() view returns(uint256)
func (_Curveplain *CurveplainCallerSession) XcpProfitA() (*big.Int, error) {
	return _Curveplain.Contract.XcpProfitA(&_Curveplain.CallOpts)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x4515cef3.
//
// Solidity: function add_liquidity(uint256[3] amounts, uint256 min_mint_amount) returns()
func (_Curveplain *CurveplainTransactor) AddLiquidity(opts *bind.TransactOpts, amounts [3]*big.Int, min_mint_amount *big.Int) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "add_liquidity", amounts, min_mint_amount)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x4515cef3.
//
// Solidity: function add_liquidity(uint256[3] amounts, uint256 min_mint_amount) returns()
func (_Curveplain *CurveplainSession) AddLiquidity(amounts [3]*big.Int, min_mint_amount *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.AddLiquidity(&_Curveplain.TransactOpts, amounts, min_mint_amount)
}

// AddLiquidity is a paid mutator transaction binding the contract method 0x4515cef3.
//
// Solidity: function add_liquidity(uint256[3] amounts, uint256 min_mint_amount) returns()
func (_Curveplain *CurveplainTransactorSession) AddLiquidity(amounts [3]*big.Int, min_mint_amount *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.AddLiquidity(&_Curveplain.TransactOpts, amounts, min_mint_amount)
}

// ApplyNewParameters is a paid mutator transaction binding the contract method 0x2a7dd7cd.
//
// Solidity: function apply_new_parameters() returns()
func (_Curveplain *CurveplainTransactor) ApplyNewParameters(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "apply_new_parameters")
}

// ApplyNewParameters is a paid mutator transaction binding the contract method 0x2a7dd7cd.
//
// Solidity: function apply_new_parameters() returns()
func (_Curveplain *CurveplainSession) ApplyNewParameters() (*types.Transaction, error) {
	return _Curveplain.Contract.ApplyNewParameters(&_Curveplain.TransactOpts)
}

// ApplyNewParameters is a paid mutator transaction binding the contract method 0x2a7dd7cd.
//
// Solidity: function apply_new_parameters() returns()
func (_Curveplain *CurveplainTransactorSession) ApplyNewParameters() (*types.Transaction, error) {
	return _Curveplain.Contract.ApplyNewParameters(&_Curveplain.TransactOpts)
}

// ApplyTransferOwnership is a paid mutator transaction binding the contract method 0x6a1c05ae.
//
// Solidity: function apply_transfer_ownership() returns()
func (_Curveplain *CurveplainTransactor) ApplyTransferOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "apply_transfer_ownership")
}

// ApplyTransferOwnership is a paid mutator transaction binding the contract method 0x6a1c05ae.
//
// Solidity: function apply_transfer_ownership() returns()
func (_Curveplain *CurveplainSession) ApplyTransferOwnership() (*types.Transaction, error) {
	return _Curveplain.Contract.ApplyTransferOwnership(&_Curveplain.TransactOpts)
}

// ApplyTransferOwnership is a paid mutator transaction binding the contract method 0x6a1c05ae.
//
// Solidity: function apply_transfer_ownership() returns()
func (_Curveplain *CurveplainTransactorSession) ApplyTransferOwnership() (*types.Transaction, error) {
	return _Curveplain.Contract.ApplyTransferOwnership(&_Curveplain.TransactOpts)
}

// ClaimAdminFees is a paid mutator transaction binding the contract method 0xc93f49e8.
//
// Solidity: function claim_admin_fees() returns()
func (_Curveplain *CurveplainTransactor) ClaimAdminFees(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "claim_admin_fees")
}

// ClaimAdminFees is a paid mutator transaction binding the contract method 0xc93f49e8.
//
// Solidity: function claim_admin_fees() returns()
func (_Curveplain *CurveplainSession) ClaimAdminFees() (*types.Transaction, error) {
	return _Curveplain.Contract.ClaimAdminFees(&_Curveplain.TransactOpts)
}

// ClaimAdminFees is a paid mutator transaction binding the contract method 0xc93f49e8.
//
// Solidity: function claim_admin_fees() returns()
func (_Curveplain *CurveplainTransactorSession) ClaimAdminFees() (*types.Transaction, error) {
	return _Curveplain.Contract.ClaimAdminFees(&_Curveplain.TransactOpts)
}

// CommitNewParameters is a paid mutator transaction binding the contract method 0xa43c3351.
//
// Solidity: function commit_new_parameters(uint256 _new_mid_fee, uint256 _new_out_fee, uint256 _new_admin_fee, uint256 _new_fee_gamma, uint256 _new_price_threshold, uint256 _new_adjustment_step, uint256 _new_ma_half_time) returns()
func (_Curveplain *CurveplainTransactor) CommitNewParameters(opts *bind.TransactOpts, _new_mid_fee *big.Int, _new_out_fee *big.Int, _new_admin_fee *big.Int, _new_fee_gamma *big.Int, _new_price_threshold *big.Int, _new_adjustment_step *big.Int, _new_ma_half_time *big.Int) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "commit_new_parameters", _new_mid_fee, _new_out_fee, _new_admin_fee, _new_fee_gamma, _new_price_threshold, _new_adjustment_step, _new_ma_half_time)
}

// CommitNewParameters is a paid mutator transaction binding the contract method 0xa43c3351.
//
// Solidity: function commit_new_parameters(uint256 _new_mid_fee, uint256 _new_out_fee, uint256 _new_admin_fee, uint256 _new_fee_gamma, uint256 _new_price_threshold, uint256 _new_adjustment_step, uint256 _new_ma_half_time) returns()
func (_Curveplain *CurveplainSession) CommitNewParameters(_new_mid_fee *big.Int, _new_out_fee *big.Int, _new_admin_fee *big.Int, _new_fee_gamma *big.Int, _new_price_threshold *big.Int, _new_adjustment_step *big.Int, _new_ma_half_time *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.CommitNewParameters(&_Curveplain.TransactOpts, _new_mid_fee, _new_out_fee, _new_admin_fee, _new_fee_gamma, _new_price_threshold, _new_adjustment_step, _new_ma_half_time)
}

// CommitNewParameters is a paid mutator transaction binding the contract method 0xa43c3351.
//
// Solidity: function commit_new_parameters(uint256 _new_mid_fee, uint256 _new_out_fee, uint256 _new_admin_fee, uint256 _new_fee_gamma, uint256 _new_price_threshold, uint256 _new_adjustment_step, uint256 _new_ma_half_time) returns()
func (_Curveplain *CurveplainTransactorSession) CommitNewParameters(_new_mid_fee *big.Int, _new_out_fee *big.Int, _new_admin_fee *big.Int, _new_fee_gamma *big.Int, _new_price_threshold *big.Int, _new_adjustment_step *big.Int, _new_ma_half_time *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.CommitNewParameters(&_Curveplain.TransactOpts, _new_mid_fee, _new_out_fee, _new_admin_fee, _new_fee_gamma, _new_price_threshold, _new_adjustment_step, _new_ma_half_time)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _owner) returns()
func (_Curveplain *CurveplainTransactor) CommitTransferOwnership(opts *bind.TransactOpts, _owner common.Address) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "commit_transfer_ownership", _owner)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _owner) returns()
func (_Curveplain *CurveplainSession) CommitTransferOwnership(_owner common.Address) (*types.Transaction, error) {
	return _Curveplain.Contract.CommitTransferOwnership(&_Curveplain.TransactOpts, _owner)
}

// CommitTransferOwnership is a paid mutator transaction binding the contract method 0x6b441a40.
//
// Solidity: function commit_transfer_ownership(address _owner) returns()
func (_Curveplain *CurveplainTransactorSession) CommitTransferOwnership(_owner common.Address) (*types.Transaction, error) {
	return _Curveplain.Contract.CommitTransferOwnership(&_Curveplain.TransactOpts, _owner)
}

// Exchange is a paid mutator transaction binding the contract method 0x5b41b908.
//
// Solidity: function exchange(uint256 i, uint256 j, uint256 dx, uint256 min_dy) payable returns()
func (_Curveplain *CurveplainTransactor) Exchange(opts *bind.TransactOpts, i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "exchange", i, j, dx, min_dy)
}

// Exchange is a paid mutator transaction binding the contract method 0x5b41b908.
//
// Solidity: function exchange(uint256 i, uint256 j, uint256 dx, uint256 min_dy) payable returns()
func (_Curveplain *CurveplainSession) Exchange(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.Exchange(&_Curveplain.TransactOpts, i, j, dx, min_dy)
}

// Exchange is a paid mutator transaction binding the contract method 0x5b41b908.
//
// Solidity: function exchange(uint256 i, uint256 j, uint256 dx, uint256 min_dy) payable returns()
func (_Curveplain *CurveplainTransactorSession) Exchange(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.Exchange(&_Curveplain.TransactOpts, i, j, dx, min_dy)
}

// Exchange0 is a paid mutator transaction binding the contract method 0x394747c5.
//
// Solidity: function exchange(uint256 i, uint256 j, uint256 dx, uint256 min_dy, bool use_eth) payable returns()
func (_Curveplain *CurveplainTransactor) Exchange0(opts *bind.TransactOpts, i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int, use_eth bool) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "exchange0", i, j, dx, min_dy, use_eth)
}

// Exchange0 is a paid mutator transaction binding the contract method 0x394747c5.
//
// Solidity: function exchange(uint256 i, uint256 j, uint256 dx, uint256 min_dy, bool use_eth) payable returns()
func (_Curveplain *CurveplainSession) Exchange0(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int, use_eth bool) (*types.Transaction, error) {
	return _Curveplain.Contract.Exchange0(&_Curveplain.TransactOpts, i, j, dx, min_dy, use_eth)
}

// Exchange0 is a paid mutator transaction binding the contract method 0x394747c5.
//
// Solidity: function exchange(uint256 i, uint256 j, uint256 dx, uint256 min_dy, bool use_eth) payable returns()
func (_Curveplain *CurveplainTransactorSession) Exchange0(i *big.Int, j *big.Int, dx *big.Int, min_dy *big.Int, use_eth bool) (*types.Transaction, error) {
	return _Curveplain.Contract.Exchange0(&_Curveplain.TransactOpts, i, j, dx, min_dy, use_eth)
}

// KillMe is a paid mutator transaction binding the contract method 0xe3698853.
//
// Solidity: function kill_me() returns()
func (_Curveplain *CurveplainTransactor) KillMe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "kill_me")
}

// KillMe is a paid mutator transaction binding the contract method 0xe3698853.
//
// Solidity: function kill_me() returns()
func (_Curveplain *CurveplainSession) KillMe() (*types.Transaction, error) {
	return _Curveplain.Contract.KillMe(&_Curveplain.TransactOpts)
}

// KillMe is a paid mutator transaction binding the contract method 0xe3698853.
//
// Solidity: function kill_me() returns()
func (_Curveplain *CurveplainTransactorSession) KillMe() (*types.Transaction, error) {
	return _Curveplain.Contract.KillMe(&_Curveplain.TransactOpts)
}

// RampAGamma is a paid mutator transaction binding the contract method 0x5e248072.
//
// Solidity: function ramp_A_gamma(uint256 future_A, uint256 future_gamma, uint256 future_time) returns()
func (_Curveplain *CurveplainTransactor) RampAGamma(opts *bind.TransactOpts, future_A *big.Int, future_gamma *big.Int, future_time *big.Int) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "ramp_A_gamma", future_A, future_gamma, future_time)
}

// RampAGamma is a paid mutator transaction binding the contract method 0x5e248072.
//
// Solidity: function ramp_A_gamma(uint256 future_A, uint256 future_gamma, uint256 future_time) returns()
func (_Curveplain *CurveplainSession) RampAGamma(future_A *big.Int, future_gamma *big.Int, future_time *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.RampAGamma(&_Curveplain.TransactOpts, future_A, future_gamma, future_time)
}

// RampAGamma is a paid mutator transaction binding the contract method 0x5e248072.
//
// Solidity: function ramp_A_gamma(uint256 future_A, uint256 future_gamma, uint256 future_time) returns()
func (_Curveplain *CurveplainTransactorSession) RampAGamma(future_A *big.Int, future_gamma *big.Int, future_time *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.RampAGamma(&_Curveplain.TransactOpts, future_A, future_gamma, future_time)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xecb586a5.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] min_amounts) returns()
func (_Curveplain *CurveplainTransactor) RemoveLiquidity(opts *bind.TransactOpts, _amount *big.Int, min_amounts [3]*big.Int) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "remove_liquidity", _amount, min_amounts)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xecb586a5.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] min_amounts) returns()
func (_Curveplain *CurveplainSession) RemoveLiquidity(_amount *big.Int, min_amounts [3]*big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.RemoveLiquidity(&_Curveplain.TransactOpts, _amount, min_amounts)
}

// RemoveLiquidity is a paid mutator transaction binding the contract method 0xecb586a5.
//
// Solidity: function remove_liquidity(uint256 _amount, uint256[3] min_amounts) returns()
func (_Curveplain *CurveplainTransactorSession) RemoveLiquidity(_amount *big.Int, min_amounts [3]*big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.RemoveLiquidity(&_Curveplain.TransactOpts, _amount, min_amounts)
}

// RemoveLiquidityOneCoin is a paid mutator transaction binding the contract method 0xf1dc3cc9.
//
// Solidity: function remove_liquidity_one_coin(uint256 token_amount, uint256 i, uint256 min_amount) returns()
func (_Curveplain *CurveplainTransactor) RemoveLiquidityOneCoin(opts *bind.TransactOpts, token_amount *big.Int, i *big.Int, min_amount *big.Int) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "remove_liquidity_one_coin", token_amount, i, min_amount)
}

// RemoveLiquidityOneCoin is a paid mutator transaction binding the contract method 0xf1dc3cc9.
//
// Solidity: function remove_liquidity_one_coin(uint256 token_amount, uint256 i, uint256 min_amount) returns()
func (_Curveplain *CurveplainSession) RemoveLiquidityOneCoin(token_amount *big.Int, i *big.Int, min_amount *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.RemoveLiquidityOneCoin(&_Curveplain.TransactOpts, token_amount, i, min_amount)
}

// RemoveLiquidityOneCoin is a paid mutator transaction binding the contract method 0xf1dc3cc9.
//
// Solidity: function remove_liquidity_one_coin(uint256 token_amount, uint256 i, uint256 min_amount) returns()
func (_Curveplain *CurveplainTransactorSession) RemoveLiquidityOneCoin(token_amount *big.Int, i *big.Int, min_amount *big.Int) (*types.Transaction, error) {
	return _Curveplain.Contract.RemoveLiquidityOneCoin(&_Curveplain.TransactOpts, token_amount, i, min_amount)
}

// RevertNewParameters is a paid mutator transaction binding the contract method 0x226840fb.
//
// Solidity: function revert_new_parameters() returns()
func (_Curveplain *CurveplainTransactor) RevertNewParameters(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "revert_new_parameters")
}

// RevertNewParameters is a paid mutator transaction binding the contract method 0x226840fb.
//
// Solidity: function revert_new_parameters() returns()
func (_Curveplain *CurveplainSession) RevertNewParameters() (*types.Transaction, error) {
	return _Curveplain.Contract.RevertNewParameters(&_Curveplain.TransactOpts)
}

// RevertNewParameters is a paid mutator transaction binding the contract method 0x226840fb.
//
// Solidity: function revert_new_parameters() returns()
func (_Curveplain *CurveplainTransactorSession) RevertNewParameters() (*types.Transaction, error) {
	return _Curveplain.Contract.RevertNewParameters(&_Curveplain.TransactOpts)
}

// RevertTransferOwnership is a paid mutator transaction binding the contract method 0x86fbf193.
//
// Solidity: function revert_transfer_ownership() returns()
func (_Curveplain *CurveplainTransactor) RevertTransferOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "revert_transfer_ownership")
}

// RevertTransferOwnership is a paid mutator transaction binding the contract method 0x86fbf193.
//
// Solidity: function revert_transfer_ownership() returns()
func (_Curveplain *CurveplainSession) RevertTransferOwnership() (*types.Transaction, error) {
	return _Curveplain.Contract.RevertTransferOwnership(&_Curveplain.TransactOpts)
}

// RevertTransferOwnership is a paid mutator transaction binding the contract method 0x86fbf193.
//
// Solidity: function revert_transfer_ownership() returns()
func (_Curveplain *CurveplainTransactorSession) RevertTransferOwnership() (*types.Transaction, error) {
	return _Curveplain.Contract.RevertTransferOwnership(&_Curveplain.TransactOpts)
}

// SetAdminFeeReceiver is a paid mutator transaction binding the contract method 0x7242e524.
//
// Solidity: function set_admin_fee_receiver(address _admin_fee_receiver) returns()
func (_Curveplain *CurveplainTransactor) SetAdminFeeReceiver(opts *bind.TransactOpts, _admin_fee_receiver common.Address) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "set_admin_fee_receiver", _admin_fee_receiver)
}

// SetAdminFeeReceiver is a paid mutator transaction binding the contract method 0x7242e524.
//
// Solidity: function set_admin_fee_receiver(address _admin_fee_receiver) returns()
func (_Curveplain *CurveplainSession) SetAdminFeeReceiver(_admin_fee_receiver common.Address) (*types.Transaction, error) {
	return _Curveplain.Contract.SetAdminFeeReceiver(&_Curveplain.TransactOpts, _admin_fee_receiver)
}

// SetAdminFeeReceiver is a paid mutator transaction binding the contract method 0x7242e524.
//
// Solidity: function set_admin_fee_receiver(address _admin_fee_receiver) returns()
func (_Curveplain *CurveplainTransactorSession) SetAdminFeeReceiver(_admin_fee_receiver common.Address) (*types.Transaction, error) {
	return _Curveplain.Contract.SetAdminFeeReceiver(&_Curveplain.TransactOpts, _admin_fee_receiver)
}

// StopRampAGamma is a paid mutator transaction binding the contract method 0x244c7c2e.
//
// Solidity: function stop_ramp_A_gamma() returns()
func (_Curveplain *CurveplainTransactor) StopRampAGamma(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "stop_ramp_A_gamma")
}

// StopRampAGamma is a paid mutator transaction binding the contract method 0x244c7c2e.
//
// Solidity: function stop_ramp_A_gamma() returns()
func (_Curveplain *CurveplainSession) StopRampAGamma() (*types.Transaction, error) {
	return _Curveplain.Contract.StopRampAGamma(&_Curveplain.TransactOpts)
}

// StopRampAGamma is a paid mutator transaction binding the contract method 0x244c7c2e.
//
// Solidity: function stop_ramp_A_gamma() returns()
func (_Curveplain *CurveplainTransactorSession) StopRampAGamma() (*types.Transaction, error) {
	return _Curveplain.Contract.StopRampAGamma(&_Curveplain.TransactOpts)
}

// UnkillMe is a paid mutator transaction binding the contract method 0x3046f972.
//
// Solidity: function unkill_me() returns()
func (_Curveplain *CurveplainTransactor) UnkillMe(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Curveplain.contract.Transact(opts, "unkill_me")
}

// UnkillMe is a paid mutator transaction binding the contract method 0x3046f972.
//
// Solidity: function unkill_me() returns()
func (_Curveplain *CurveplainSession) UnkillMe() (*types.Transaction, error) {
	return _Curveplain.Contract.UnkillMe(&_Curveplain.TransactOpts)
}

// UnkillMe is a paid mutator transaction binding the contract method 0x3046f972.
//
// Solidity: function unkill_me() returns()
func (_Curveplain *CurveplainTransactorSession) UnkillMe() (*types.Transaction, error) {
	return _Curveplain.Contract.UnkillMe(&_Curveplain.TransactOpts)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Curveplain *CurveplainTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _Curveplain.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Curveplain *CurveplainSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Curveplain.Contract.Fallback(&_Curveplain.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_Curveplain *CurveplainTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _Curveplain.Contract.Fallback(&_Curveplain.TransactOpts, calldata)
}

// CurveplainAddLiquidityIterator is returned from FilterAddLiquidity and is used to iterate over the raw logs and unpacked data for AddLiquidity events raised by the Curveplain contract.
type CurveplainAddLiquidityIterator struct {
	Event *CurveplainAddLiquidity // Event containing the contract specifics and raw log

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
func (it *CurveplainAddLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainAddLiquidity)
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
		it.Event = new(CurveplainAddLiquidity)
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
func (it *CurveplainAddLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainAddLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainAddLiquidity represents a AddLiquidity event raised by the Curveplain contract.
type CurveplainAddLiquidity struct {
	Provider     common.Address
	TokenAmounts [3]*big.Int
	Fee          *big.Int
	TokenSupply  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterAddLiquidity is a free log retrieval operation binding the contract event 0x96b486485420b963edd3fdec0b0195730035600feb7de6f544383d7950fa97ee.
//
// Solidity: event AddLiquidity(address indexed provider, uint256[3] token_amounts, uint256 fee, uint256 token_supply)
func (_Curveplain *CurveplainFilterer) FilterAddLiquidity(opts *bind.FilterOpts, provider []common.Address) (*CurveplainAddLiquidityIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "AddLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainAddLiquidityIterator{contract: _Curveplain.contract, event: "AddLiquidity", logs: logs, sub: sub}, nil
}

// WatchAddLiquidity is a free log subscription operation binding the contract event 0x96b486485420b963edd3fdec0b0195730035600feb7de6f544383d7950fa97ee.
//
// Solidity: event AddLiquidity(address indexed provider, uint256[3] token_amounts, uint256 fee, uint256 token_supply)
func (_Curveplain *CurveplainFilterer) WatchAddLiquidity(opts *bind.WatchOpts, sink chan<- *CurveplainAddLiquidity, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "AddLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainAddLiquidity)
				if err := _Curveplain.contract.UnpackLog(event, "AddLiquidity", log); err != nil {
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

// ParseAddLiquidity is a log parse operation binding the contract event 0x96b486485420b963edd3fdec0b0195730035600feb7de6f544383d7950fa97ee.
//
// Solidity: event AddLiquidity(address indexed provider, uint256[3] token_amounts, uint256 fee, uint256 token_supply)
func (_Curveplain *CurveplainFilterer) ParseAddLiquidity(log types.Log) (*CurveplainAddLiquidity, error) {
	event := new(CurveplainAddLiquidity)
	if err := _Curveplain.contract.UnpackLog(event, "AddLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainClaimAdminFeeIterator is returned from FilterClaimAdminFee and is used to iterate over the raw logs and unpacked data for ClaimAdminFee events raised by the Curveplain contract.
type CurveplainClaimAdminFeeIterator struct {
	Event *CurveplainClaimAdminFee // Event containing the contract specifics and raw log

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
func (it *CurveplainClaimAdminFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainClaimAdminFee)
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
		it.Event = new(CurveplainClaimAdminFee)
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
func (it *CurveplainClaimAdminFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainClaimAdminFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainClaimAdminFee represents a ClaimAdminFee event raised by the Curveplain contract.
type CurveplainClaimAdminFee struct {
	Admin  common.Address
	Tokens *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClaimAdminFee is a free log retrieval operation binding the contract event 0x6059a38198b1dc42b3791087d1ff0fbd72b3179553c25f678cd246f52ffaaf59.
//
// Solidity: event ClaimAdminFee(address indexed admin, uint256 tokens)
func (_Curveplain *CurveplainFilterer) FilterClaimAdminFee(opts *bind.FilterOpts, admin []common.Address) (*CurveplainClaimAdminFeeIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "ClaimAdminFee", adminRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainClaimAdminFeeIterator{contract: _Curveplain.contract, event: "ClaimAdminFee", logs: logs, sub: sub}, nil
}

// WatchClaimAdminFee is a free log subscription operation binding the contract event 0x6059a38198b1dc42b3791087d1ff0fbd72b3179553c25f678cd246f52ffaaf59.
//
// Solidity: event ClaimAdminFee(address indexed admin, uint256 tokens)
func (_Curveplain *CurveplainFilterer) WatchClaimAdminFee(opts *bind.WatchOpts, sink chan<- *CurveplainClaimAdminFee, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "ClaimAdminFee", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainClaimAdminFee)
				if err := _Curveplain.contract.UnpackLog(event, "ClaimAdminFee", log); err != nil {
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

// ParseClaimAdminFee is a log parse operation binding the contract event 0x6059a38198b1dc42b3791087d1ff0fbd72b3179553c25f678cd246f52ffaaf59.
//
// Solidity: event ClaimAdminFee(address indexed admin, uint256 tokens)
func (_Curveplain *CurveplainFilterer) ParseClaimAdminFee(log types.Log) (*CurveplainClaimAdminFee, error) {
	event := new(CurveplainClaimAdminFee)
	if err := _Curveplain.contract.UnpackLog(event, "ClaimAdminFee", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainCommitNewAdminIterator is returned from FilterCommitNewAdmin and is used to iterate over the raw logs and unpacked data for CommitNewAdmin events raised by the Curveplain contract.
type CurveplainCommitNewAdminIterator struct {
	Event *CurveplainCommitNewAdmin // Event containing the contract specifics and raw log

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
func (it *CurveplainCommitNewAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainCommitNewAdmin)
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
		it.Event = new(CurveplainCommitNewAdmin)
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
func (it *CurveplainCommitNewAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainCommitNewAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainCommitNewAdmin represents a CommitNewAdmin event raised by the Curveplain contract.
type CurveplainCommitNewAdmin struct {
	Deadline *big.Int
	Admin    common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterCommitNewAdmin is a free log retrieval operation binding the contract event 0x181aa3aa17d4cbf99265dd4443eba009433d3cde79d60164fde1d1a192beb935.
//
// Solidity: event CommitNewAdmin(uint256 indexed deadline, address indexed admin)
func (_Curveplain *CurveplainFilterer) FilterCommitNewAdmin(opts *bind.FilterOpts, deadline []*big.Int, admin []common.Address) (*CurveplainCommitNewAdminIterator, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "CommitNewAdmin", deadlineRule, adminRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainCommitNewAdminIterator{contract: _Curveplain.contract, event: "CommitNewAdmin", logs: logs, sub: sub}, nil
}

// WatchCommitNewAdmin is a free log subscription operation binding the contract event 0x181aa3aa17d4cbf99265dd4443eba009433d3cde79d60164fde1d1a192beb935.
//
// Solidity: event CommitNewAdmin(uint256 indexed deadline, address indexed admin)
func (_Curveplain *CurveplainFilterer) WatchCommitNewAdmin(opts *bind.WatchOpts, sink chan<- *CurveplainCommitNewAdmin, deadline []*big.Int, admin []common.Address) (event.Subscription, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}
	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "CommitNewAdmin", deadlineRule, adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainCommitNewAdmin)
				if err := _Curveplain.contract.UnpackLog(event, "CommitNewAdmin", log); err != nil {
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
func (_Curveplain *CurveplainFilterer) ParseCommitNewAdmin(log types.Log) (*CurveplainCommitNewAdmin, error) {
	event := new(CurveplainCommitNewAdmin)
	if err := _Curveplain.contract.UnpackLog(event, "CommitNewAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainCommitNewParametersIterator is returned from FilterCommitNewParameters and is used to iterate over the raw logs and unpacked data for CommitNewParameters events raised by the Curveplain contract.
type CurveplainCommitNewParametersIterator struct {
	Event *CurveplainCommitNewParameters // Event containing the contract specifics and raw log

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
func (it *CurveplainCommitNewParametersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainCommitNewParameters)
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
		it.Event = new(CurveplainCommitNewParameters)
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
func (it *CurveplainCommitNewParametersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainCommitNewParametersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainCommitNewParameters represents a CommitNewParameters event raised by the Curveplain contract.
type CurveplainCommitNewParameters struct {
	Deadline       *big.Int
	AdminFee       *big.Int
	MidFee         *big.Int
	OutFee         *big.Int
	FeeGamma       *big.Int
	PriceThreshold *big.Int
	AdjustmentStep *big.Int
	MaHalfTime     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterCommitNewParameters is a free log retrieval operation binding the contract event 0x913fde9a37e1f8ab67876a4d0ce80790d764fcfc5692f4529526df9c6bdde553.
//
// Solidity: event CommitNewParameters(uint256 indexed deadline, uint256 admin_fee, uint256 mid_fee, uint256 out_fee, uint256 fee_gamma, uint256 price_threshold, uint256 adjustment_step, uint256 ma_half_time)
func (_Curveplain *CurveplainFilterer) FilterCommitNewParameters(opts *bind.FilterOpts, deadline []*big.Int) (*CurveplainCommitNewParametersIterator, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "CommitNewParameters", deadlineRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainCommitNewParametersIterator{contract: _Curveplain.contract, event: "CommitNewParameters", logs: logs, sub: sub}, nil
}

// WatchCommitNewParameters is a free log subscription operation binding the contract event 0x913fde9a37e1f8ab67876a4d0ce80790d764fcfc5692f4529526df9c6bdde553.
//
// Solidity: event CommitNewParameters(uint256 indexed deadline, uint256 admin_fee, uint256 mid_fee, uint256 out_fee, uint256 fee_gamma, uint256 price_threshold, uint256 adjustment_step, uint256 ma_half_time)
func (_Curveplain *CurveplainFilterer) WatchCommitNewParameters(opts *bind.WatchOpts, sink chan<- *CurveplainCommitNewParameters, deadline []*big.Int) (event.Subscription, error) {

	var deadlineRule []interface{}
	for _, deadlineItem := range deadline {
		deadlineRule = append(deadlineRule, deadlineItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "CommitNewParameters", deadlineRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainCommitNewParameters)
				if err := _Curveplain.contract.UnpackLog(event, "CommitNewParameters", log); err != nil {
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

// ParseCommitNewParameters is a log parse operation binding the contract event 0x913fde9a37e1f8ab67876a4d0ce80790d764fcfc5692f4529526df9c6bdde553.
//
// Solidity: event CommitNewParameters(uint256 indexed deadline, uint256 admin_fee, uint256 mid_fee, uint256 out_fee, uint256 fee_gamma, uint256 price_threshold, uint256 adjustment_step, uint256 ma_half_time)
func (_Curveplain *CurveplainFilterer) ParseCommitNewParameters(log types.Log) (*CurveplainCommitNewParameters, error) {
	event := new(CurveplainCommitNewParameters)
	if err := _Curveplain.contract.UnpackLog(event, "CommitNewParameters", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainNewAdminIterator is returned from FilterNewAdmin and is used to iterate over the raw logs and unpacked data for NewAdmin events raised by the Curveplain contract.
type CurveplainNewAdminIterator struct {
	Event *CurveplainNewAdmin // Event containing the contract specifics and raw log

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
func (it *CurveplainNewAdminIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainNewAdmin)
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
		it.Event = new(CurveplainNewAdmin)
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
func (it *CurveplainNewAdminIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainNewAdminIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainNewAdmin represents a NewAdmin event raised by the Curveplain contract.
type CurveplainNewAdmin struct {
	Admin common.Address
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterNewAdmin is a free log retrieval operation binding the contract event 0x71614071b88dee5e0b2ae578a9dd7b2ebbe9ae832ba419dc0242cd065a290b6c.
//
// Solidity: event NewAdmin(address indexed admin)
func (_Curveplain *CurveplainFilterer) FilterNewAdmin(opts *bind.FilterOpts, admin []common.Address) (*CurveplainNewAdminIterator, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "NewAdmin", adminRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainNewAdminIterator{contract: _Curveplain.contract, event: "NewAdmin", logs: logs, sub: sub}, nil
}

// WatchNewAdmin is a free log subscription operation binding the contract event 0x71614071b88dee5e0b2ae578a9dd7b2ebbe9ae832ba419dc0242cd065a290b6c.
//
// Solidity: event NewAdmin(address indexed admin)
func (_Curveplain *CurveplainFilterer) WatchNewAdmin(opts *bind.WatchOpts, sink chan<- *CurveplainNewAdmin, admin []common.Address) (event.Subscription, error) {

	var adminRule []interface{}
	for _, adminItem := range admin {
		adminRule = append(adminRule, adminItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "NewAdmin", adminRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainNewAdmin)
				if err := _Curveplain.contract.UnpackLog(event, "NewAdmin", log); err != nil {
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
func (_Curveplain *CurveplainFilterer) ParseNewAdmin(log types.Log) (*CurveplainNewAdmin, error) {
	event := new(CurveplainNewAdmin)
	if err := _Curveplain.contract.UnpackLog(event, "NewAdmin", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainNewParametersIterator is returned from FilterNewParameters and is used to iterate over the raw logs and unpacked data for NewParameters events raised by the Curveplain contract.
type CurveplainNewParametersIterator struct {
	Event *CurveplainNewParameters // Event containing the contract specifics and raw log

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
func (it *CurveplainNewParametersIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainNewParameters)
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
		it.Event = new(CurveplainNewParameters)
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
func (it *CurveplainNewParametersIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainNewParametersIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainNewParameters represents a NewParameters event raised by the Curveplain contract.
type CurveplainNewParameters struct {
	AdminFee       *big.Int
	MidFee         *big.Int
	OutFee         *big.Int
	FeeGamma       *big.Int
	PriceThreshold *big.Int
	AdjustmentStep *big.Int
	MaHalfTime     *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterNewParameters is a free log retrieval operation binding the contract event 0x1c65bbdc939f346e5d6f0bde1f072819947438d4fc7b182cc59c2f6dc5504087.
//
// Solidity: event NewParameters(uint256 admin_fee, uint256 mid_fee, uint256 out_fee, uint256 fee_gamma, uint256 price_threshold, uint256 adjustment_step, uint256 ma_half_time)
func (_Curveplain *CurveplainFilterer) FilterNewParameters(opts *bind.FilterOpts) (*CurveplainNewParametersIterator, error) {

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "NewParameters")
	if err != nil {
		return nil, err
	}
	return &CurveplainNewParametersIterator{contract: _Curveplain.contract, event: "NewParameters", logs: logs, sub: sub}, nil
}

// WatchNewParameters is a free log subscription operation binding the contract event 0x1c65bbdc939f346e5d6f0bde1f072819947438d4fc7b182cc59c2f6dc5504087.
//
// Solidity: event NewParameters(uint256 admin_fee, uint256 mid_fee, uint256 out_fee, uint256 fee_gamma, uint256 price_threshold, uint256 adjustment_step, uint256 ma_half_time)
func (_Curveplain *CurveplainFilterer) WatchNewParameters(opts *bind.WatchOpts, sink chan<- *CurveplainNewParameters) (event.Subscription, error) {

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "NewParameters")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainNewParameters)
				if err := _Curveplain.contract.UnpackLog(event, "NewParameters", log); err != nil {
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

// ParseNewParameters is a log parse operation binding the contract event 0x1c65bbdc939f346e5d6f0bde1f072819947438d4fc7b182cc59c2f6dc5504087.
//
// Solidity: event NewParameters(uint256 admin_fee, uint256 mid_fee, uint256 out_fee, uint256 fee_gamma, uint256 price_threshold, uint256 adjustment_step, uint256 ma_half_time)
func (_Curveplain *CurveplainFilterer) ParseNewParameters(log types.Log) (*CurveplainNewParameters, error) {
	event := new(CurveplainNewParameters)
	if err := _Curveplain.contract.UnpackLog(event, "NewParameters", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainRampAgammaIterator is returned from FilterRampAgamma and is used to iterate over the raw logs and unpacked data for RampAgamma events raised by the Curveplain contract.
type CurveplainRampAgammaIterator struct {
	Event *CurveplainRampAgamma // Event containing the contract specifics and raw log

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
func (it *CurveplainRampAgammaIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainRampAgamma)
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
		it.Event = new(CurveplainRampAgamma)
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
func (it *CurveplainRampAgammaIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainRampAgammaIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainRampAgamma represents a RampAgamma event raised by the Curveplain contract.
type CurveplainRampAgamma struct {
	InitialA    *big.Int
	FutureA     *big.Int
	InitialTime *big.Int
	FutureTime  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRampAgamma is a free log retrieval operation binding the contract event 0x7dce008c7221b9d42f26174f7435e50f9cc8f502c57019724c090b546415a029.
//
// Solidity: event RampAgamma(uint256 initial_A, uint256 future_A, uint256 initial_time, uint256 future_time)
func (_Curveplain *CurveplainFilterer) FilterRampAgamma(opts *bind.FilterOpts) (*CurveplainRampAgammaIterator, error) {

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "RampAgamma")
	if err != nil {
		return nil, err
	}
	return &CurveplainRampAgammaIterator{contract: _Curveplain.contract, event: "RampAgamma", logs: logs, sub: sub}, nil
}

// WatchRampAgamma is a free log subscription operation binding the contract event 0x7dce008c7221b9d42f26174f7435e50f9cc8f502c57019724c090b546415a029.
//
// Solidity: event RampAgamma(uint256 initial_A, uint256 future_A, uint256 initial_time, uint256 future_time)
func (_Curveplain *CurveplainFilterer) WatchRampAgamma(opts *bind.WatchOpts, sink chan<- *CurveplainRampAgamma) (event.Subscription, error) {

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "RampAgamma")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainRampAgamma)
				if err := _Curveplain.contract.UnpackLog(event, "RampAgamma", log); err != nil {
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

// ParseRampAgamma is a log parse operation binding the contract event 0x7dce008c7221b9d42f26174f7435e50f9cc8f502c57019724c090b546415a029.
//
// Solidity: event RampAgamma(uint256 initial_A, uint256 future_A, uint256 initial_time, uint256 future_time)
func (_Curveplain *CurveplainFilterer) ParseRampAgamma(log types.Log) (*CurveplainRampAgamma, error) {
	event := new(CurveplainRampAgamma)
	if err := _Curveplain.contract.UnpackLog(event, "RampAgamma", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainRemoveLiquidityIterator is returned from FilterRemoveLiquidity and is used to iterate over the raw logs and unpacked data for RemoveLiquidity events raised by the Curveplain contract.
type CurveplainRemoveLiquidityIterator struct {
	Event *CurveplainRemoveLiquidity // Event containing the contract specifics and raw log

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
func (it *CurveplainRemoveLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainRemoveLiquidity)
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
		it.Event = new(CurveplainRemoveLiquidity)
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
func (it *CurveplainRemoveLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainRemoveLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainRemoveLiquidity represents a RemoveLiquidity event raised by the Curveplain contract.
type CurveplainRemoveLiquidity struct {
	Provider     common.Address
	TokenAmounts [3]*big.Int
	TokenSupply  *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterRemoveLiquidity is a free log retrieval operation binding the contract event 0xd6cc314a0b1e3b2579f8e64248e82434072e8271290eef8ad0886709304195f5.
//
// Solidity: event RemoveLiquidity(address indexed provider, uint256[3] token_amounts, uint256 token_supply)
func (_Curveplain *CurveplainFilterer) FilterRemoveLiquidity(opts *bind.FilterOpts, provider []common.Address) (*CurveplainRemoveLiquidityIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "RemoveLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainRemoveLiquidityIterator{contract: _Curveplain.contract, event: "RemoveLiquidity", logs: logs, sub: sub}, nil
}

// WatchRemoveLiquidity is a free log subscription operation binding the contract event 0xd6cc314a0b1e3b2579f8e64248e82434072e8271290eef8ad0886709304195f5.
//
// Solidity: event RemoveLiquidity(address indexed provider, uint256[3] token_amounts, uint256 token_supply)
func (_Curveplain *CurveplainFilterer) WatchRemoveLiquidity(opts *bind.WatchOpts, sink chan<- *CurveplainRemoveLiquidity, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "RemoveLiquidity", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainRemoveLiquidity)
				if err := _Curveplain.contract.UnpackLog(event, "RemoveLiquidity", log); err != nil {
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

// ParseRemoveLiquidity is a log parse operation binding the contract event 0xd6cc314a0b1e3b2579f8e64248e82434072e8271290eef8ad0886709304195f5.
//
// Solidity: event RemoveLiquidity(address indexed provider, uint256[3] token_amounts, uint256 token_supply)
func (_Curveplain *CurveplainFilterer) ParseRemoveLiquidity(log types.Log) (*CurveplainRemoveLiquidity, error) {
	event := new(CurveplainRemoveLiquidity)
	if err := _Curveplain.contract.UnpackLog(event, "RemoveLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainRemoveLiquidityOneIterator is returned from FilterRemoveLiquidityOne and is used to iterate over the raw logs and unpacked data for RemoveLiquidityOne events raised by the Curveplain contract.
type CurveplainRemoveLiquidityOneIterator struct {
	Event *CurveplainRemoveLiquidityOne // Event containing the contract specifics and raw log

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
func (it *CurveplainRemoveLiquidityOneIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainRemoveLiquidityOne)
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
		it.Event = new(CurveplainRemoveLiquidityOne)
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
func (it *CurveplainRemoveLiquidityOneIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainRemoveLiquidityOneIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainRemoveLiquidityOne represents a RemoveLiquidityOne event raised by the Curveplain contract.
type CurveplainRemoveLiquidityOne struct {
	Provider    common.Address
	TokenAmount *big.Int
	CoinIndex   *big.Int
	CoinAmount  *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterRemoveLiquidityOne is a free log retrieval operation binding the contract event 0x5ad056f2e28a8cec232015406b843668c1e36cda598127ec3b8c59b8c72773a0.
//
// Solidity: event RemoveLiquidityOne(address indexed provider, uint256 token_amount, uint256 coin_index, uint256 coin_amount)
func (_Curveplain *CurveplainFilterer) FilterRemoveLiquidityOne(opts *bind.FilterOpts, provider []common.Address) (*CurveplainRemoveLiquidityOneIterator, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "RemoveLiquidityOne", providerRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainRemoveLiquidityOneIterator{contract: _Curveplain.contract, event: "RemoveLiquidityOne", logs: logs, sub: sub}, nil
}

// WatchRemoveLiquidityOne is a free log subscription operation binding the contract event 0x5ad056f2e28a8cec232015406b843668c1e36cda598127ec3b8c59b8c72773a0.
//
// Solidity: event RemoveLiquidityOne(address indexed provider, uint256 token_amount, uint256 coin_index, uint256 coin_amount)
func (_Curveplain *CurveplainFilterer) WatchRemoveLiquidityOne(opts *bind.WatchOpts, sink chan<- *CurveplainRemoveLiquidityOne, provider []common.Address) (event.Subscription, error) {

	var providerRule []interface{}
	for _, providerItem := range provider {
		providerRule = append(providerRule, providerItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "RemoveLiquidityOne", providerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainRemoveLiquidityOne)
				if err := _Curveplain.contract.UnpackLog(event, "RemoveLiquidityOne", log); err != nil {
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

// ParseRemoveLiquidityOne is a log parse operation binding the contract event 0x5ad056f2e28a8cec232015406b843668c1e36cda598127ec3b8c59b8c72773a0.
//
// Solidity: event RemoveLiquidityOne(address indexed provider, uint256 token_amount, uint256 coin_index, uint256 coin_amount)
func (_Curveplain *CurveplainFilterer) ParseRemoveLiquidityOne(log types.Log) (*CurveplainRemoveLiquidityOne, error) {
	event := new(CurveplainRemoveLiquidityOne)
	if err := _Curveplain.contract.UnpackLog(event, "RemoveLiquidityOne", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainStopRampAIterator is returned from FilterStopRampA and is used to iterate over the raw logs and unpacked data for StopRampA events raised by the Curveplain contract.
type CurveplainStopRampAIterator struct {
	Event *CurveplainStopRampA // Event containing the contract specifics and raw log

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
func (it *CurveplainStopRampAIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainStopRampA)
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
		it.Event = new(CurveplainStopRampA)
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
func (it *CurveplainStopRampAIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainStopRampAIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainStopRampA represents a StopRampA event raised by the Curveplain contract.
type CurveplainStopRampA struct {
	CurrentA     *big.Int
	CurrentGamma *big.Int
	Time         *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterStopRampA is a free log retrieval operation binding the contract event 0x5f0e7fba3d100c9e19446e1c92fe436f0a9a22fe99669360e4fdd6d3de2fc284.
//
// Solidity: event StopRampA(uint256 current_A, uint256 current_gamma, uint256 time)
func (_Curveplain *CurveplainFilterer) FilterStopRampA(opts *bind.FilterOpts) (*CurveplainStopRampAIterator, error) {

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "StopRampA")
	if err != nil {
		return nil, err
	}
	return &CurveplainStopRampAIterator{contract: _Curveplain.contract, event: "StopRampA", logs: logs, sub: sub}, nil
}

// WatchStopRampA is a free log subscription operation binding the contract event 0x5f0e7fba3d100c9e19446e1c92fe436f0a9a22fe99669360e4fdd6d3de2fc284.
//
// Solidity: event StopRampA(uint256 current_A, uint256 current_gamma, uint256 time)
func (_Curveplain *CurveplainFilterer) WatchStopRampA(opts *bind.WatchOpts, sink chan<- *CurveplainStopRampA) (event.Subscription, error) {

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "StopRampA")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainStopRampA)
				if err := _Curveplain.contract.UnpackLog(event, "StopRampA", log); err != nil {
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

// ParseStopRampA is a log parse operation binding the contract event 0x5f0e7fba3d100c9e19446e1c92fe436f0a9a22fe99669360e4fdd6d3de2fc284.
//
// Solidity: event StopRampA(uint256 current_A, uint256 current_gamma, uint256 time)
func (_Curveplain *CurveplainFilterer) ParseStopRampA(log types.Log) (*CurveplainStopRampA, error) {
	event := new(CurveplainStopRampA)
	if err := _Curveplain.contract.UnpackLog(event, "StopRampA", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// CurveplainTokenExchangeIterator is returned from FilterTokenExchange and is used to iterate over the raw logs and unpacked data for TokenExchange events raised by the Curveplain contract.
type CurveplainTokenExchangeIterator struct {
	Event *CurveplainTokenExchange // Event containing the contract specifics and raw log

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
func (it *CurveplainTokenExchangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(CurveplainTokenExchange)
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
		it.Event = new(CurveplainTokenExchange)
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
func (it *CurveplainTokenExchangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *CurveplainTokenExchangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// CurveplainTokenExchange represents a TokenExchange event raised by the Curveplain contract.
type CurveplainTokenExchange struct {
	Buyer        common.Address
	SoldId       *big.Int
	TokensSold   *big.Int
	BoughtId     *big.Int
	TokensBought *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterTokenExchange is a free log retrieval operation binding the contract event 0xb2e76ae99761dc136e598d4a629bb347eccb9532a5f8bbd72e18467c3c34cc98.
//
// Solidity: event TokenExchange(address indexed buyer, uint256 sold_id, uint256 tokens_sold, uint256 bought_id, uint256 tokens_bought)
func (_Curveplain *CurveplainFilterer) FilterTokenExchange(opts *bind.FilterOpts, buyer []common.Address) (*CurveplainTokenExchangeIterator, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Curveplain.contract.FilterLogs(opts, "TokenExchange", buyerRule)
	if err != nil {
		return nil, err
	}
	return &CurveplainTokenExchangeIterator{contract: _Curveplain.contract, event: "TokenExchange", logs: logs, sub: sub}, nil
}

// WatchTokenExchange is a free log subscription operation binding the contract event 0xb2e76ae99761dc136e598d4a629bb347eccb9532a5f8bbd72e18467c3c34cc98.
//
// Solidity: event TokenExchange(address indexed buyer, uint256 sold_id, uint256 tokens_sold, uint256 bought_id, uint256 tokens_bought)
func (_Curveplain *CurveplainFilterer) WatchTokenExchange(opts *bind.WatchOpts, sink chan<- *CurveplainTokenExchange, buyer []common.Address) (event.Subscription, error) {

	var buyerRule []interface{}
	for _, buyerItem := range buyer {
		buyerRule = append(buyerRule, buyerItem)
	}

	logs, sub, err := _Curveplain.contract.WatchLogs(opts, "TokenExchange", buyerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(CurveplainTokenExchange)
				if err := _Curveplain.contract.UnpackLog(event, "TokenExchange", log); err != nil {
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

// ParseTokenExchange is a log parse operation binding the contract event 0xb2e76ae99761dc136e598d4a629bb347eccb9532a5f8bbd72e18467c3c34cc98.
//
// Solidity: event TokenExchange(address indexed buyer, uint256 sold_id, uint256 tokens_sold, uint256 bought_id, uint256 tokens_bought)
func (_Curveplain *CurveplainFilterer) ParseTokenExchange(log types.Log) (*CurveplainTokenExchange, error) {
	event := new(CurveplainTokenExchange)
	if err := _Curveplain.contract.UnpackLog(event, "TokenExchange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
