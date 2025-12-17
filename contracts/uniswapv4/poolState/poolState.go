// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package poolState

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

// PoolStateMetaData contains all meta data concerning the PoolState contract.
var PoolStateMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"contractIPoolManager\",\"name\":\"_poolManager\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"}],\"name\":\"getFeeGrowthGlobals\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeGrowthGlobal0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthGlobal1\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"}],\"name\":\"getFeeGrowthInside\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside0X128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside1X128\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"}],\"name\":\"getLiquidity\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"positionId\",\"type\":\"bytes32\"}],\"name\":\"getPositionInfo\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside0LastX128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside1LastX128\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"getPositionInfo\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside0LastX128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthInside1LastX128\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"positionId\",\"type\":\"bytes32\"}],\"name\":\"getPositionLiquidity\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"}],\"name\":\"getSlot0\",\"outputs\":[{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"internalType\":\"uint24\",\"name\":\"protocolFee\",\"type\":\"uint24\"},{\"internalType\":\"uint24\",\"name\":\"lpFee\",\"type\":\"uint24\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"int16\",\"name\":\"tick\",\"type\":\"int16\"}],\"name\":\"getTickBitmap\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tickBitmap\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"name\":\"getTickFeeGrowthOutside\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"feeGrowthOutside0X128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthOutside1X128\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"name\":\"getTickInfo\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidityGross\",\"type\":\"uint128\"},{\"internalType\":\"int128\",\"name\":\"liquidityNet\",\"type\":\"int128\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthOutside0X128\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"feeGrowthOutside1X128\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"PoolId\",\"name\":\"poolId\",\"type\":\"bytes32\"},{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"name\":\"getTickLiquidity\",\"outputs\":[{\"internalType\":\"uint128\",\"name\":\"liquidityGross\",\"type\":\"uint128\"},{\"internalType\":\"int128\",\"name\":\"liquidityNet\",\"type\":\"int128\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolManager\",\"outputs\":[{\"internalType\":\"contractIPoolManager\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// PoolStateABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolStateMetaData.ABI instead.
var PoolStateABI = PoolStateMetaData.ABI

// PoolState is an auto generated Go binding around an Ethereum contract.
type PoolState struct {
	PoolStateCaller     // Read-only binding to the contract
	PoolStateTransactor // Write-only binding to the contract
	PoolStateFilterer   // Log filterer for contract events
}

// PoolStateCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolStateCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolStateTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolStateTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolStateFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolStateFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolStateSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolStateSession struct {
	Contract     *PoolState        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolStateCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolStateCallerSession struct {
	Contract *PoolStateCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// PoolStateTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolStateTransactorSession struct {
	Contract     *PoolStateTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// PoolStateRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolStateRaw struct {
	Contract *PoolState // Generic contract binding to access the raw methods on
}

// PoolStateCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolStateCallerRaw struct {
	Contract *PoolStateCaller // Generic read-only contract binding to access the raw methods on
}

// PoolStateTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolStateTransactorRaw struct {
	Contract *PoolStateTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolState creates a new instance of PoolState, bound to a specific deployed contract.
func NewPoolState(address common.Address, backend bind.ContractBackend) (*PoolState, error) {
	contract, err := bindPoolState(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolState{PoolStateCaller: PoolStateCaller{contract: contract}, PoolStateTransactor: PoolStateTransactor{contract: contract}, PoolStateFilterer: PoolStateFilterer{contract: contract}}, nil
}

// NewPoolStateCaller creates a new read-only instance of PoolState, bound to a specific deployed contract.
func NewPoolStateCaller(address common.Address, caller bind.ContractCaller) (*PoolStateCaller, error) {
	contract, err := bindPoolState(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolStateCaller{contract: contract}, nil
}

// NewPoolStateTransactor creates a new write-only instance of PoolState, bound to a specific deployed contract.
func NewPoolStateTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolStateTransactor, error) {
	contract, err := bindPoolState(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolStateTransactor{contract: contract}, nil
}

// NewPoolStateFilterer creates a new log filterer instance of PoolState, bound to a specific deployed contract.
func NewPoolStateFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolStateFilterer, error) {
	contract, err := bindPoolState(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolStateFilterer{contract: contract}, nil
}

// bindPoolState binds a generic wrapper to an already deployed contract.
func bindPoolState(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolStateMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolState *PoolStateRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolState.Contract.PoolStateCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolState *PoolStateRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolState.Contract.PoolStateTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolState *PoolStateRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolState.Contract.PoolStateTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolState *PoolStateCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolState.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolState *PoolStateTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolState.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolState *PoolStateTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolState.Contract.contract.Transact(opts, method, params...)
}

// GetFeeGrowthGlobals is a free data retrieval call binding the contract method 0x9ec538c8.
//
// Solidity: function getFeeGrowthGlobals(bytes32 poolId) view returns(uint256 feeGrowthGlobal0, uint256 feeGrowthGlobal1)
func (_PoolState *PoolStateCaller) GetFeeGrowthGlobals(opts *bind.CallOpts, poolId [32]byte) (struct {
	FeeGrowthGlobal0 *big.Int
	FeeGrowthGlobal1 *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getFeeGrowthGlobals", poolId)

	outstruct := new(struct {
		FeeGrowthGlobal0 *big.Int
		FeeGrowthGlobal1 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FeeGrowthGlobal0 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthGlobal1 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetFeeGrowthGlobals is a free data retrieval call binding the contract method 0x9ec538c8.
//
// Solidity: function getFeeGrowthGlobals(bytes32 poolId) view returns(uint256 feeGrowthGlobal0, uint256 feeGrowthGlobal1)
func (_PoolState *PoolStateSession) GetFeeGrowthGlobals(poolId [32]byte) (struct {
	FeeGrowthGlobal0 *big.Int
	FeeGrowthGlobal1 *big.Int
}, error) {
	return _PoolState.Contract.GetFeeGrowthGlobals(&_PoolState.CallOpts, poolId)
}

// GetFeeGrowthGlobals is a free data retrieval call binding the contract method 0x9ec538c8.
//
// Solidity: function getFeeGrowthGlobals(bytes32 poolId) view returns(uint256 feeGrowthGlobal0, uint256 feeGrowthGlobal1)
func (_PoolState *PoolStateCallerSession) GetFeeGrowthGlobals(poolId [32]byte) (struct {
	FeeGrowthGlobal0 *big.Int
	FeeGrowthGlobal1 *big.Int
}, error) {
	return _PoolState.Contract.GetFeeGrowthGlobals(&_PoolState.CallOpts, poolId)
}

// GetFeeGrowthInside is a free data retrieval call binding the contract method 0x53e9c1fb.
//
// Solidity: function getFeeGrowthInside(bytes32 poolId, int24 tickLower, int24 tickUpper) view returns(uint256 feeGrowthInside0X128, uint256 feeGrowthInside1X128)
func (_PoolState *PoolStateCaller) GetFeeGrowthInside(opts *bind.CallOpts, poolId [32]byte, tickLower *big.Int, tickUpper *big.Int) (struct {
	FeeGrowthInside0X128 *big.Int
	FeeGrowthInside1X128 *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getFeeGrowthInside", poolId, tickLower, tickUpper)

	outstruct := new(struct {
		FeeGrowthInside0X128 *big.Int
		FeeGrowthInside1X128 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FeeGrowthInside0X128 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside1X128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetFeeGrowthInside is a free data retrieval call binding the contract method 0x53e9c1fb.
//
// Solidity: function getFeeGrowthInside(bytes32 poolId, int24 tickLower, int24 tickUpper) view returns(uint256 feeGrowthInside0X128, uint256 feeGrowthInside1X128)
func (_PoolState *PoolStateSession) GetFeeGrowthInside(poolId [32]byte, tickLower *big.Int, tickUpper *big.Int) (struct {
	FeeGrowthInside0X128 *big.Int
	FeeGrowthInside1X128 *big.Int
}, error) {
	return _PoolState.Contract.GetFeeGrowthInside(&_PoolState.CallOpts, poolId, tickLower, tickUpper)
}

// GetFeeGrowthInside is a free data retrieval call binding the contract method 0x53e9c1fb.
//
// Solidity: function getFeeGrowthInside(bytes32 poolId, int24 tickLower, int24 tickUpper) view returns(uint256 feeGrowthInside0X128, uint256 feeGrowthInside1X128)
func (_PoolState *PoolStateCallerSession) GetFeeGrowthInside(poolId [32]byte, tickLower *big.Int, tickUpper *big.Int) (struct {
	FeeGrowthInside0X128 *big.Int
	FeeGrowthInside1X128 *big.Int
}, error) {
	return _PoolState.Contract.GetFeeGrowthInside(&_PoolState.CallOpts, poolId, tickLower, tickUpper)
}

// GetLiquidity is a free data retrieval call binding the contract method 0xfa6793d5.
//
// Solidity: function getLiquidity(bytes32 poolId) view returns(uint128 liquidity)
func (_PoolState *PoolStateCaller) GetLiquidity(opts *bind.CallOpts, poolId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getLiquidity", poolId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetLiquidity is a free data retrieval call binding the contract method 0xfa6793d5.
//
// Solidity: function getLiquidity(bytes32 poolId) view returns(uint128 liquidity)
func (_PoolState *PoolStateSession) GetLiquidity(poolId [32]byte) (*big.Int, error) {
	return _PoolState.Contract.GetLiquidity(&_PoolState.CallOpts, poolId)
}

// GetLiquidity is a free data retrieval call binding the contract method 0xfa6793d5.
//
// Solidity: function getLiquidity(bytes32 poolId) view returns(uint128 liquidity)
func (_PoolState *PoolStateCallerSession) GetLiquidity(poolId [32]byte) (*big.Int, error) {
	return _PoolState.Contract.GetLiquidity(&_PoolState.CallOpts, poolId)
}

// GetPositionInfo is a free data retrieval call binding the contract method 0x97fd7b42.
//
// Solidity: function getPositionInfo(bytes32 poolId, bytes32 positionId) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128)
func (_PoolState *PoolStateCaller) GetPositionInfo(opts *bind.CallOpts, poolId [32]byte, positionId [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getPositionInfo", poolId, positionId)

	outstruct := new(struct {
		Liquidity                *big.Int
		FeeGrowthInside0LastX128 *big.Int
		FeeGrowthInside1LastX128 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Liquidity = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside0LastX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside1LastX128 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPositionInfo is a free data retrieval call binding the contract method 0x97fd7b42.
//
// Solidity: function getPositionInfo(bytes32 poolId, bytes32 positionId) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128)
func (_PoolState *PoolStateSession) GetPositionInfo(poolId [32]byte, positionId [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
}, error) {
	return _PoolState.Contract.GetPositionInfo(&_PoolState.CallOpts, poolId, positionId)
}

// GetPositionInfo is a free data retrieval call binding the contract method 0x97fd7b42.
//
// Solidity: function getPositionInfo(bytes32 poolId, bytes32 positionId) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128)
func (_PoolState *PoolStateCallerSession) GetPositionInfo(poolId [32]byte, positionId [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
}, error) {
	return _PoolState.Contract.GetPositionInfo(&_PoolState.CallOpts, poolId, positionId)
}

// GetPositionInfo0 is a free data retrieval call binding the contract method 0xdacf1d2f.
//
// Solidity: function getPositionInfo(bytes32 poolId, address owner, int24 tickLower, int24 tickUpper, bytes32 salt) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128)
func (_PoolState *PoolStateCaller) GetPositionInfo0(opts *bind.CallOpts, poolId [32]byte, owner common.Address, tickLower *big.Int, tickUpper *big.Int, salt [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getPositionInfo0", poolId, owner, tickLower, tickUpper, salt)

	outstruct := new(struct {
		Liquidity                *big.Int
		FeeGrowthInside0LastX128 *big.Int
		FeeGrowthInside1LastX128 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Liquidity = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside0LastX128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthInside1LastX128 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetPositionInfo0 is a free data retrieval call binding the contract method 0xdacf1d2f.
//
// Solidity: function getPositionInfo(bytes32 poolId, address owner, int24 tickLower, int24 tickUpper, bytes32 salt) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128)
func (_PoolState *PoolStateSession) GetPositionInfo0(poolId [32]byte, owner common.Address, tickLower *big.Int, tickUpper *big.Int, salt [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
}, error) {
	return _PoolState.Contract.GetPositionInfo0(&_PoolState.CallOpts, poolId, owner, tickLower, tickUpper, salt)
}

// GetPositionInfo0 is a free data retrieval call binding the contract method 0xdacf1d2f.
//
// Solidity: function getPositionInfo(bytes32 poolId, address owner, int24 tickLower, int24 tickUpper, bytes32 salt) view returns(uint128 liquidity, uint256 feeGrowthInside0LastX128, uint256 feeGrowthInside1LastX128)
func (_PoolState *PoolStateCallerSession) GetPositionInfo0(poolId [32]byte, owner common.Address, tickLower *big.Int, tickUpper *big.Int, salt [32]byte) (struct {
	Liquidity                *big.Int
	FeeGrowthInside0LastX128 *big.Int
	FeeGrowthInside1LastX128 *big.Int
}, error) {
	return _PoolState.Contract.GetPositionInfo0(&_PoolState.CallOpts, poolId, owner, tickLower, tickUpper, salt)
}

// GetPositionLiquidity is a free data retrieval call binding the contract method 0xf0928f29.
//
// Solidity: function getPositionLiquidity(bytes32 poolId, bytes32 positionId) view returns(uint128 liquidity)
func (_PoolState *PoolStateCaller) GetPositionLiquidity(opts *bind.CallOpts, poolId [32]byte, positionId [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getPositionLiquidity", poolId, positionId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPositionLiquidity is a free data retrieval call binding the contract method 0xf0928f29.
//
// Solidity: function getPositionLiquidity(bytes32 poolId, bytes32 positionId) view returns(uint128 liquidity)
func (_PoolState *PoolStateSession) GetPositionLiquidity(poolId [32]byte, positionId [32]byte) (*big.Int, error) {
	return _PoolState.Contract.GetPositionLiquidity(&_PoolState.CallOpts, poolId, positionId)
}

// GetPositionLiquidity is a free data retrieval call binding the contract method 0xf0928f29.
//
// Solidity: function getPositionLiquidity(bytes32 poolId, bytes32 positionId) view returns(uint128 liquidity)
func (_PoolState *PoolStateCallerSession) GetPositionLiquidity(poolId [32]byte, positionId [32]byte) (*big.Int, error) {
	return _PoolState.Contract.GetPositionLiquidity(&_PoolState.CallOpts, poolId, positionId)
}

// GetSlot0 is a free data retrieval call binding the contract method 0xc815641c.
//
// Solidity: function getSlot0(bytes32 poolId) view returns(uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee)
func (_PoolState *PoolStateCaller) GetSlot0(opts *bind.CallOpts, poolId [32]byte) (struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	ProtocolFee  *big.Int
	LpFee        *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getSlot0", poolId)

	outstruct := new(struct {
		SqrtPriceX96 *big.Int
		Tick         *big.Int
		ProtocolFee  *big.Int
		LpFee        *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.SqrtPriceX96 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Tick = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.ProtocolFee = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.LpFee = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetSlot0 is a free data retrieval call binding the contract method 0xc815641c.
//
// Solidity: function getSlot0(bytes32 poolId) view returns(uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee)
func (_PoolState *PoolStateSession) GetSlot0(poolId [32]byte) (struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	ProtocolFee  *big.Int
	LpFee        *big.Int
}, error) {
	return _PoolState.Contract.GetSlot0(&_PoolState.CallOpts, poolId)
}

// GetSlot0 is a free data retrieval call binding the contract method 0xc815641c.
//
// Solidity: function getSlot0(bytes32 poolId) view returns(uint160 sqrtPriceX96, int24 tick, uint24 protocolFee, uint24 lpFee)
func (_PoolState *PoolStateCallerSession) GetSlot0(poolId [32]byte) (struct {
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	ProtocolFee  *big.Int
	LpFee        *big.Int
}, error) {
	return _PoolState.Contract.GetSlot0(&_PoolState.CallOpts, poolId)
}

// GetTickBitmap is a free data retrieval call binding the contract method 0x1c7ccb4c.
//
// Solidity: function getTickBitmap(bytes32 poolId, int16 tick) view returns(uint256 tickBitmap)
func (_PoolState *PoolStateCaller) GetTickBitmap(opts *bind.CallOpts, poolId [32]byte, tick int16) (*big.Int, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getTickBitmap", poolId, tick)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTickBitmap is a free data retrieval call binding the contract method 0x1c7ccb4c.
//
// Solidity: function getTickBitmap(bytes32 poolId, int16 tick) view returns(uint256 tickBitmap)
func (_PoolState *PoolStateSession) GetTickBitmap(poolId [32]byte, tick int16) (*big.Int, error) {
	return _PoolState.Contract.GetTickBitmap(&_PoolState.CallOpts, poolId, tick)
}

// GetTickBitmap is a free data retrieval call binding the contract method 0x1c7ccb4c.
//
// Solidity: function getTickBitmap(bytes32 poolId, int16 tick) view returns(uint256 tickBitmap)
func (_PoolState *PoolStateCallerSession) GetTickBitmap(poolId [32]byte, tick int16) (*big.Int, error) {
	return _PoolState.Contract.GetTickBitmap(&_PoolState.CallOpts, poolId, tick)
}

// GetTickFeeGrowthOutside is a free data retrieval call binding the contract method 0x8a2bb9e6.
//
// Solidity: function getTickFeeGrowthOutside(bytes32 poolId, int24 tick) view returns(uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128)
func (_PoolState *PoolStateCaller) GetTickFeeGrowthOutside(opts *bind.CallOpts, poolId [32]byte, tick *big.Int) (struct {
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getTickFeeGrowthOutside", poolId, tick)

	outstruct := new(struct {
		FeeGrowthOutside0X128 *big.Int
		FeeGrowthOutside1X128 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.FeeGrowthOutside0X128 = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthOutside1X128 = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTickFeeGrowthOutside is a free data retrieval call binding the contract method 0x8a2bb9e6.
//
// Solidity: function getTickFeeGrowthOutside(bytes32 poolId, int24 tick) view returns(uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128)
func (_PoolState *PoolStateSession) GetTickFeeGrowthOutside(poolId [32]byte, tick *big.Int) (struct {
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
}, error) {
	return _PoolState.Contract.GetTickFeeGrowthOutside(&_PoolState.CallOpts, poolId, tick)
}

// GetTickFeeGrowthOutside is a free data retrieval call binding the contract method 0x8a2bb9e6.
//
// Solidity: function getTickFeeGrowthOutside(bytes32 poolId, int24 tick) view returns(uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128)
func (_PoolState *PoolStateCallerSession) GetTickFeeGrowthOutside(poolId [32]byte, tick *big.Int) (struct {
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
}, error) {
	return _PoolState.Contract.GetTickFeeGrowthOutside(&_PoolState.CallOpts, poolId, tick)
}

// GetTickInfo is a free data retrieval call binding the contract method 0x7c40f1fe.
//
// Solidity: function getTickInfo(bytes32 poolId, int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128)
func (_PoolState *PoolStateCaller) GetTickInfo(opts *bind.CallOpts, poolId [32]byte, tick *big.Int) (struct {
	LiquidityGross        *big.Int
	LiquidityNet          *big.Int
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getTickInfo", poolId, tick)

	outstruct := new(struct {
		LiquidityGross        *big.Int
		LiquidityNet          *big.Int
		FeeGrowthOutside0X128 *big.Int
		FeeGrowthOutside1X128 *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LiquidityGross = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LiquidityNet = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthOutside0X128 = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.FeeGrowthOutside1X128 = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTickInfo is a free data retrieval call binding the contract method 0x7c40f1fe.
//
// Solidity: function getTickInfo(bytes32 poolId, int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128)
func (_PoolState *PoolStateSession) GetTickInfo(poolId [32]byte, tick *big.Int) (struct {
	LiquidityGross        *big.Int
	LiquidityNet          *big.Int
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
}, error) {
	return _PoolState.Contract.GetTickInfo(&_PoolState.CallOpts, poolId, tick)
}

// GetTickInfo is a free data retrieval call binding the contract method 0x7c40f1fe.
//
// Solidity: function getTickInfo(bytes32 poolId, int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet, uint256 feeGrowthOutside0X128, uint256 feeGrowthOutside1X128)
func (_PoolState *PoolStateCallerSession) GetTickInfo(poolId [32]byte, tick *big.Int) (struct {
	LiquidityGross        *big.Int
	LiquidityNet          *big.Int
	FeeGrowthOutside0X128 *big.Int
	FeeGrowthOutside1X128 *big.Int
}, error) {
	return _PoolState.Contract.GetTickInfo(&_PoolState.CallOpts, poolId, tick)
}

// GetTickLiquidity is a free data retrieval call binding the contract method 0xcaedab54.
//
// Solidity: function getTickLiquidity(bytes32 poolId, int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet)
func (_PoolState *PoolStateCaller) GetTickLiquidity(opts *bind.CallOpts, poolId [32]byte, tick *big.Int) (struct {
	LiquidityGross *big.Int
	LiquidityNet   *big.Int
}, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "getTickLiquidity", poolId, tick)

	outstruct := new(struct {
		LiquidityGross *big.Int
		LiquidityNet   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.LiquidityGross = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.LiquidityNet = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetTickLiquidity is a free data retrieval call binding the contract method 0xcaedab54.
//
// Solidity: function getTickLiquidity(bytes32 poolId, int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet)
func (_PoolState *PoolStateSession) GetTickLiquidity(poolId [32]byte, tick *big.Int) (struct {
	LiquidityGross *big.Int
	LiquidityNet   *big.Int
}, error) {
	return _PoolState.Contract.GetTickLiquidity(&_PoolState.CallOpts, poolId, tick)
}

// GetTickLiquidity is a free data retrieval call binding the contract method 0xcaedab54.
//
// Solidity: function getTickLiquidity(bytes32 poolId, int24 tick) view returns(uint128 liquidityGross, int128 liquidityNet)
func (_PoolState *PoolStateCallerSession) GetTickLiquidity(poolId [32]byte, tick *big.Int) (struct {
	LiquidityGross *big.Int
	LiquidityNet   *big.Int
}, error) {
	return _PoolState.Contract.GetTickLiquidity(&_PoolState.CallOpts, poolId, tick)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_PoolState *PoolStateCaller) PoolManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolState.contract.Call(opts, &out, "poolManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_PoolState *PoolStateSession) PoolManager() (common.Address, error) {
	return _PoolState.Contract.PoolManager(&_PoolState.CallOpts)
}

// PoolManager is a free data retrieval call binding the contract method 0xdc4c90d3.
//
// Solidity: function poolManager() view returns(address)
func (_PoolState *PoolStateCallerSession) PoolManager() (common.Address, error) {
	return _PoolState.Contract.PoolManager(&_PoolState.CallOpts)
}
