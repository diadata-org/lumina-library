// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package poolManager

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

// IPoolManagerModifyLiquidityParams is an auto generated low-level Go binding around an user-defined struct.
type IPoolManagerModifyLiquidityParams struct {
	TickLower      *big.Int
	TickUpper      *big.Int
	LiquidityDelta *big.Int
	Salt           [32]byte
}

// IPoolManagerSwapParams is an auto generated low-level Go binding around an user-defined struct.
type IPoolManagerSwapParams struct {
	ZeroForOne        bool
	AmountSpecified   *big.Int
	SqrtPriceLimitX96 *big.Int
}

// PoolKey is an auto generated low-level Go binding around an user-defined struct.
type PoolKey struct {
	Currency0   common.Address
	Currency1   common.Address
	Fee         *big.Int
	TickSpacing *big.Int
	Hooks       common.Address
}

// PoolManagerMetaData contains all meta data concerning the PoolManager contract.
var PoolManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AlreadyUnlocked\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"currency1\",\"type\":\"address\"}],\"name\":\"CurrenciesOutOfOrderOrEqual\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"CurrencyNotSettled\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"DelegateCallNotAllowed\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidCaller\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ManagerLocked\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"MustClearExactPositiveDelta\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NonzeroNativeValue\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"PoolNotInitialized\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ProtocolFeeCurrencySynced\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"ProtocolFeeTooLarge\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SwapAmountCannotBeZero\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"name\":\"TickSpacingTooLarge\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"}],\"name\":\"TickSpacingTooSmall\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"UnauthorizedDynamicLPFeeUpdate\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"PoolId\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"}],\"name\":\"Donate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"PoolId\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"name\":\"Initialize\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"PoolId\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"int256\",\"name\":\"liquidityDelta\",\"type\":\"int256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"name\":\"ModifyLiquidity\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"OperatorSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"protocolFeeController\",\"type\":\"address\"}],\"name\":\"ProtocolFeeControllerUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"PoolId\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"protocolFee\",\"type\":\"uint24\"}],\"name\":\"ProtocolFeeUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"PoolId\",\"name\":\"id\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"int128\",\"name\":\"amount0\",\"type\":\"int128\"},{\"indexed\":false,\"internalType\":\"int128\",\"name\":\"amount1\",\"type\":\"int128\"},{\"indexed\":false,\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"},{\"indexed\":false,\"internalType\":\"uint128\",\"name\":\"liquidity\",\"type\":\"uint128\"},{\"indexed\":false,\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"},{\"indexed\":false,\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"}],\"name\":\"Swap\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"caller\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"allowance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"spender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"burn\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"Currency\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"clear\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"collectProtocolFees\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountCollected\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"amount0\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount1\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"hookData\",\"type\":\"bytes\"}],\"name\":\"donate\",\"outputs\":[{\"internalType\":\"BalanceDelta\",\"name\":\"delta\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"extsload\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"startSlot\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"nSlots\",\"type\":\"uint256\"}],\"name\":\"extsload\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"slots\",\"type\":\"bytes32[]\"}],\"name\":\"extsload\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"slots\",\"type\":\"bytes32[]\"}],\"name\":\"exttload\",\"outputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"\",\"type\":\"bytes32[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"slot\",\"type\":\"bytes32\"}],\"name\":\"exttload\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceX96\",\"type\":\"uint160\"}],\"name\":\"initialize\",\"outputs\":[{\"internalType\":\"int24\",\"name\":\"tick\",\"type\":\"int24\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"isOperator\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"int24\",\"name\":\"tickLower\",\"type\":\"int24\"},{\"internalType\":\"int24\",\"name\":\"tickUpper\",\"type\":\"int24\"},{\"internalType\":\"int256\",\"name\":\"liquidityDelta\",\"type\":\"int256\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"}],\"internalType\":\"structIPoolManager.ModifyLiquidityParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"hookData\",\"type\":\"bytes\"}],\"name\":\"modifyLiquidity\",\"outputs\":[{\"internalType\":\"BalanceDelta\",\"name\":\"callerDelta\",\"type\":\"int256\"},{\"internalType\":\"BalanceDelta\",\"name\":\"feesAccrued\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"protocolFeeController\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"Currency\",\"name\":\"currency\",\"type\":\"address\"}],\"name\":\"protocolFeesAccrued\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setOperator\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint24\",\"name\":\"newProtocolFee\",\"type\":\"uint24\"}],\"name\":\"setProtocolFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"controller\",\"type\":\"address\"}],\"name\":\"setProtocolFeeController\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"settle\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"recipient\",\"type\":\"address\"}],\"name\":\"settleFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bool\",\"name\":\"zeroForOne\",\"type\":\"bool\"},{\"internalType\":\"int256\",\"name\":\"amountSpecified\",\"type\":\"int256\"},{\"internalType\":\"uint160\",\"name\":\"sqrtPriceLimitX96\",\"type\":\"uint160\"}],\"internalType\":\"structIPoolManager.SwapParams\",\"name\":\"params\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"hookData\",\"type\":\"bytes\"}],\"name\":\"swap\",\"outputs\":[{\"internalType\":\"BalanceDelta\",\"name\":\"swapDelta\",\"type\":\"int256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"Currency\",\"name\":\"currency\",\"type\":\"address\"}],\"name\":\"sync\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"Currency\",\"name\":\"currency\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"take\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"id\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"unlock\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"result\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"Currency\",\"name\":\"currency0\",\"type\":\"address\"},{\"internalType\":\"Currency\",\"name\":\"currency1\",\"type\":\"address\"},{\"internalType\":\"uint24\",\"name\":\"fee\",\"type\":\"uint24\"},{\"internalType\":\"int24\",\"name\":\"tickSpacing\",\"type\":\"int24\"},{\"internalType\":\"contractIHooks\",\"name\":\"hooks\",\"type\":\"address\"}],\"internalType\":\"structPoolKey\",\"name\":\"key\",\"type\":\"tuple\"},{\"internalType\":\"uint24\",\"name\":\"newDynamicLPFee\",\"type\":\"uint24\"}],\"name\":\"updateDynamicLPFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// PoolManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use PoolManagerMetaData.ABI instead.
var PoolManagerABI = PoolManagerMetaData.ABI

// PoolManager is an auto generated Go binding around an Ethereum contract.
type PoolManager struct {
	PoolManagerCaller     // Read-only binding to the contract
	PoolManagerTransactor // Write-only binding to the contract
	PoolManagerFilterer   // Log filterer for contract events
}

// PoolManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type PoolManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type PoolManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type PoolManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// PoolManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type PoolManagerSession struct {
	Contract     *PoolManager      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// PoolManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type PoolManagerCallerSession struct {
	Contract *PoolManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// PoolManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type PoolManagerTransactorSession struct {
	Contract     *PoolManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// PoolManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type PoolManagerRaw struct {
	Contract *PoolManager // Generic contract binding to access the raw methods on
}

// PoolManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type PoolManagerCallerRaw struct {
	Contract *PoolManagerCaller // Generic read-only contract binding to access the raw methods on
}

// PoolManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type PoolManagerTransactorRaw struct {
	Contract *PoolManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewPoolManager creates a new instance of PoolManager, bound to a specific deployed contract.
func NewPoolManager(address common.Address, backend bind.ContractBackend) (*PoolManager, error) {
	contract, err := bindPoolManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &PoolManager{PoolManagerCaller: PoolManagerCaller{contract: contract}, PoolManagerTransactor: PoolManagerTransactor{contract: contract}, PoolManagerFilterer: PoolManagerFilterer{contract: contract}}, nil
}

// NewPoolManagerCaller creates a new read-only instance of PoolManager, bound to a specific deployed contract.
func NewPoolManagerCaller(address common.Address, caller bind.ContractCaller) (*PoolManagerCaller, error) {
	contract, err := bindPoolManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &PoolManagerCaller{contract: contract}, nil
}

// NewPoolManagerTransactor creates a new write-only instance of PoolManager, bound to a specific deployed contract.
func NewPoolManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*PoolManagerTransactor, error) {
	contract, err := bindPoolManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &PoolManagerTransactor{contract: contract}, nil
}

// NewPoolManagerFilterer creates a new log filterer instance of PoolManager, bound to a specific deployed contract.
func NewPoolManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*PoolManagerFilterer, error) {
	contract, err := bindPoolManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &PoolManagerFilterer{contract: contract}, nil
}

// bindPoolManager binds a generic wrapper to an already deployed contract.
func bindPoolManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := PoolManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolManager *PoolManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolManager.Contract.PoolManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolManager *PoolManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManager.Contract.PoolManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolManager *PoolManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolManager.Contract.PoolManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_PoolManager *PoolManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _PoolManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_PoolManager *PoolManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_PoolManager *PoolManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _PoolManager.Contract.contract.Transact(opts, method, params...)
}

// Allowance is a free data retrieval call binding the contract method 0x598af9e7.
//
// Solidity: function allowance(address owner, address spender, uint256 id) view returns(uint256 amount)
func (_PoolManager *PoolManagerCaller) Allowance(opts *bind.CallOpts, owner common.Address, spender common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "allowance", owner, spender, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Allowance is a free data retrieval call binding the contract method 0x598af9e7.
//
// Solidity: function allowance(address owner, address spender, uint256 id) view returns(uint256 amount)
func (_PoolManager *PoolManagerSession) Allowance(owner common.Address, spender common.Address, id *big.Int) (*big.Int, error) {
	return _PoolManager.Contract.Allowance(&_PoolManager.CallOpts, owner, spender, id)
}

// Allowance is a free data retrieval call binding the contract method 0x598af9e7.
//
// Solidity: function allowance(address owner, address spender, uint256 id) view returns(uint256 amount)
func (_PoolManager *PoolManagerCallerSession) Allowance(owner common.Address, spender common.Address, id *big.Int) (*big.Int, error) {
	return _PoolManager.Contract.Allowance(&_PoolManager.CallOpts, owner, spender, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address owner, uint256 id) view returns(uint256 balance)
func (_PoolManager *PoolManagerCaller) BalanceOf(opts *bind.CallOpts, owner common.Address, id *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "balanceOf", owner, id)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address owner, uint256 id) view returns(uint256 balance)
func (_PoolManager *PoolManagerSession) BalanceOf(owner common.Address, id *big.Int) (*big.Int, error) {
	return _PoolManager.Contract.BalanceOf(&_PoolManager.CallOpts, owner, id)
}

// BalanceOf is a free data retrieval call binding the contract method 0x00fdd58e.
//
// Solidity: function balanceOf(address owner, uint256 id) view returns(uint256 balance)
func (_PoolManager *PoolManagerCallerSession) BalanceOf(owner common.Address, id *big.Int) (*big.Int, error) {
	return _PoolManager.Contract.BalanceOf(&_PoolManager.CallOpts, owner, id)
}

// Extsload is a free data retrieval call binding the contract method 0x1e2eaeaf.
//
// Solidity: function extsload(bytes32 slot) view returns(bytes32)
func (_PoolManager *PoolManagerCaller) Extsload(opts *bind.CallOpts, slot [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "extsload", slot)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Extsload is a free data retrieval call binding the contract method 0x1e2eaeaf.
//
// Solidity: function extsload(bytes32 slot) view returns(bytes32)
func (_PoolManager *PoolManagerSession) Extsload(slot [32]byte) ([32]byte, error) {
	return _PoolManager.Contract.Extsload(&_PoolManager.CallOpts, slot)
}

// Extsload is a free data retrieval call binding the contract method 0x1e2eaeaf.
//
// Solidity: function extsload(bytes32 slot) view returns(bytes32)
func (_PoolManager *PoolManagerCallerSession) Extsload(slot [32]byte) ([32]byte, error) {
	return _PoolManager.Contract.Extsload(&_PoolManager.CallOpts, slot)
}

// Extsload0 is a free data retrieval call binding the contract method 0x35fd631a.
//
// Solidity: function extsload(bytes32 startSlot, uint256 nSlots) view returns(bytes32[])
func (_PoolManager *PoolManagerCaller) Extsload0(opts *bind.CallOpts, startSlot [32]byte, nSlots *big.Int) ([][32]byte, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "extsload0", startSlot, nSlots)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// Extsload0 is a free data retrieval call binding the contract method 0x35fd631a.
//
// Solidity: function extsload(bytes32 startSlot, uint256 nSlots) view returns(bytes32[])
func (_PoolManager *PoolManagerSession) Extsload0(startSlot [32]byte, nSlots *big.Int) ([][32]byte, error) {
	return _PoolManager.Contract.Extsload0(&_PoolManager.CallOpts, startSlot, nSlots)
}

// Extsload0 is a free data retrieval call binding the contract method 0x35fd631a.
//
// Solidity: function extsload(bytes32 startSlot, uint256 nSlots) view returns(bytes32[])
func (_PoolManager *PoolManagerCallerSession) Extsload0(startSlot [32]byte, nSlots *big.Int) ([][32]byte, error) {
	return _PoolManager.Contract.Extsload0(&_PoolManager.CallOpts, startSlot, nSlots)
}

// Extsload1 is a free data retrieval call binding the contract method 0xdbd035ff.
//
// Solidity: function extsload(bytes32[] slots) view returns(bytes32[])
func (_PoolManager *PoolManagerCaller) Extsload1(opts *bind.CallOpts, slots [][32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "extsload1", slots)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// Extsload1 is a free data retrieval call binding the contract method 0xdbd035ff.
//
// Solidity: function extsload(bytes32[] slots) view returns(bytes32[])
func (_PoolManager *PoolManagerSession) Extsload1(slots [][32]byte) ([][32]byte, error) {
	return _PoolManager.Contract.Extsload1(&_PoolManager.CallOpts, slots)
}

// Extsload1 is a free data retrieval call binding the contract method 0xdbd035ff.
//
// Solidity: function extsload(bytes32[] slots) view returns(bytes32[])
func (_PoolManager *PoolManagerCallerSession) Extsload1(slots [][32]byte) ([][32]byte, error) {
	return _PoolManager.Contract.Extsload1(&_PoolManager.CallOpts, slots)
}

// Exttload is a free data retrieval call binding the contract method 0x9bf6645f.
//
// Solidity: function exttload(bytes32[] slots) view returns(bytes32[])
func (_PoolManager *PoolManagerCaller) Exttload(opts *bind.CallOpts, slots [][32]byte) ([][32]byte, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "exttload", slots)

	if err != nil {
		return *new([][32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][32]byte)).(*[][32]byte)

	return out0, err

}

// Exttload is a free data retrieval call binding the contract method 0x9bf6645f.
//
// Solidity: function exttload(bytes32[] slots) view returns(bytes32[])
func (_PoolManager *PoolManagerSession) Exttload(slots [][32]byte) ([][32]byte, error) {
	return _PoolManager.Contract.Exttload(&_PoolManager.CallOpts, slots)
}

// Exttload is a free data retrieval call binding the contract method 0x9bf6645f.
//
// Solidity: function exttload(bytes32[] slots) view returns(bytes32[])
func (_PoolManager *PoolManagerCallerSession) Exttload(slots [][32]byte) ([][32]byte, error) {
	return _PoolManager.Contract.Exttload(&_PoolManager.CallOpts, slots)
}

// Exttload0 is a free data retrieval call binding the contract method 0xf135baaa.
//
// Solidity: function exttload(bytes32 slot) view returns(bytes32)
func (_PoolManager *PoolManagerCaller) Exttload0(opts *bind.CallOpts, slot [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "exttload0", slot)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// Exttload0 is a free data retrieval call binding the contract method 0xf135baaa.
//
// Solidity: function exttload(bytes32 slot) view returns(bytes32)
func (_PoolManager *PoolManagerSession) Exttload0(slot [32]byte) ([32]byte, error) {
	return _PoolManager.Contract.Exttload0(&_PoolManager.CallOpts, slot)
}

// Exttload0 is a free data retrieval call binding the contract method 0xf135baaa.
//
// Solidity: function exttload(bytes32 slot) view returns(bytes32)
func (_PoolManager *PoolManagerCallerSession) Exttload0(slot [32]byte) ([32]byte, error) {
	return _PoolManager.Contract.Exttload0(&_PoolManager.CallOpts, slot)
}

// IsOperator is a free data retrieval call binding the contract method 0xb6363cf2.
//
// Solidity: function isOperator(address owner, address operator) view returns(bool isOperator)
func (_PoolManager *PoolManagerCaller) IsOperator(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "isOperator", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsOperator is a free data retrieval call binding the contract method 0xb6363cf2.
//
// Solidity: function isOperator(address owner, address operator) view returns(bool isOperator)
func (_PoolManager *PoolManagerSession) IsOperator(owner common.Address, operator common.Address) (bool, error) {
	return _PoolManager.Contract.IsOperator(&_PoolManager.CallOpts, owner, operator)
}

// IsOperator is a free data retrieval call binding the contract method 0xb6363cf2.
//
// Solidity: function isOperator(address owner, address operator) view returns(bool isOperator)
func (_PoolManager *PoolManagerCallerSession) IsOperator(owner common.Address, operator common.Address) (bool, error) {
	return _PoolManager.Contract.IsOperator(&_PoolManager.CallOpts, owner, operator)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PoolManager *PoolManagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PoolManager *PoolManagerSession) Owner() (common.Address, error) {
	return _PoolManager.Contract.Owner(&_PoolManager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_PoolManager *PoolManagerCallerSession) Owner() (common.Address, error) {
	return _PoolManager.Contract.Owner(&_PoolManager.CallOpts)
}

// ProtocolFeeController is a free data retrieval call binding the contract method 0xf02de3b2.
//
// Solidity: function protocolFeeController() view returns(address)
func (_PoolManager *PoolManagerCaller) ProtocolFeeController(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "protocolFeeController")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ProtocolFeeController is a free data retrieval call binding the contract method 0xf02de3b2.
//
// Solidity: function protocolFeeController() view returns(address)
func (_PoolManager *PoolManagerSession) ProtocolFeeController() (common.Address, error) {
	return _PoolManager.Contract.ProtocolFeeController(&_PoolManager.CallOpts)
}

// ProtocolFeeController is a free data retrieval call binding the contract method 0xf02de3b2.
//
// Solidity: function protocolFeeController() view returns(address)
func (_PoolManager *PoolManagerCallerSession) ProtocolFeeController() (common.Address, error) {
	return _PoolManager.Contract.ProtocolFeeController(&_PoolManager.CallOpts)
}

// ProtocolFeesAccrued is a free data retrieval call binding the contract method 0x97e8cd4e.
//
// Solidity: function protocolFeesAccrued(address currency) view returns(uint256 amount)
func (_PoolManager *PoolManagerCaller) ProtocolFeesAccrued(opts *bind.CallOpts, currency common.Address) (*big.Int, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "protocolFeesAccrued", currency)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ProtocolFeesAccrued is a free data retrieval call binding the contract method 0x97e8cd4e.
//
// Solidity: function protocolFeesAccrued(address currency) view returns(uint256 amount)
func (_PoolManager *PoolManagerSession) ProtocolFeesAccrued(currency common.Address) (*big.Int, error) {
	return _PoolManager.Contract.ProtocolFeesAccrued(&_PoolManager.CallOpts, currency)
}

// ProtocolFeesAccrued is a free data retrieval call binding the contract method 0x97e8cd4e.
//
// Solidity: function protocolFeesAccrued(address currency) view returns(uint256 amount)
func (_PoolManager *PoolManagerCallerSession) ProtocolFeesAccrued(currency common.Address) (*big.Int, error) {
	return _PoolManager.Contract.ProtocolFeesAccrued(&_PoolManager.CallOpts, currency)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PoolManager *PoolManagerCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _PoolManager.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PoolManager *PoolManagerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PoolManager.Contract.SupportsInterface(&_PoolManager.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_PoolManager *PoolManagerCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _PoolManager.Contract.SupportsInterface(&_PoolManager.CallOpts, interfaceId)
}

// Approve is a paid mutator transaction binding the contract method 0x426a8493.
//
// Solidity: function approve(address spender, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerTransactor) Approve(opts *bind.TransactOpts, spender common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "approve", spender, id, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x426a8493.
//
// Solidity: function approve(address spender, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerSession) Approve(spender common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Approve(&_PoolManager.TransactOpts, spender, id, amount)
}

// Approve is a paid mutator transaction binding the contract method 0x426a8493.
//
// Solidity: function approve(address spender, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerTransactorSession) Approve(spender common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Approve(&_PoolManager.TransactOpts, spender, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 id, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactor) Burn(opts *bind.TransactOpts, from common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "burn", from, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 id, uint256 amount) returns()
func (_PoolManager *PoolManagerSession) Burn(from common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Burn(&_PoolManager.TransactOpts, from, id, amount)
}

// Burn is a paid mutator transaction binding the contract method 0xf5298aca.
//
// Solidity: function burn(address from, uint256 id, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactorSession) Burn(from common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Burn(&_PoolManager.TransactOpts, from, id, amount)
}

// Clear is a paid mutator transaction binding the contract method 0x80f0b44c.
//
// Solidity: function clear(address currency, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactor) Clear(opts *bind.TransactOpts, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "clear", currency, amount)
}

// Clear is a paid mutator transaction binding the contract method 0x80f0b44c.
//
// Solidity: function clear(address currency, uint256 amount) returns()
func (_PoolManager *PoolManagerSession) Clear(currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Clear(&_PoolManager.TransactOpts, currency, amount)
}

// Clear is a paid mutator transaction binding the contract method 0x80f0b44c.
//
// Solidity: function clear(address currency, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactorSession) Clear(currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Clear(&_PoolManager.TransactOpts, currency, amount)
}

// CollectProtocolFees is a paid mutator transaction binding the contract method 0x8161b874.
//
// Solidity: function collectProtocolFees(address recipient, address currency, uint256 amount) returns(uint256 amountCollected)
func (_PoolManager *PoolManagerTransactor) CollectProtocolFees(opts *bind.TransactOpts, recipient common.Address, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "collectProtocolFees", recipient, currency, amount)
}

// CollectProtocolFees is a paid mutator transaction binding the contract method 0x8161b874.
//
// Solidity: function collectProtocolFees(address recipient, address currency, uint256 amount) returns(uint256 amountCollected)
func (_PoolManager *PoolManagerSession) CollectProtocolFees(recipient common.Address, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.CollectProtocolFees(&_PoolManager.TransactOpts, recipient, currency, amount)
}

// CollectProtocolFees is a paid mutator transaction binding the contract method 0x8161b874.
//
// Solidity: function collectProtocolFees(address recipient, address currency, uint256 amount) returns(uint256 amountCollected)
func (_PoolManager *PoolManagerTransactorSession) CollectProtocolFees(recipient common.Address, currency common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.CollectProtocolFees(&_PoolManager.TransactOpts, recipient, currency, amount)
}

// Donate is a paid mutator transaction binding the contract method 0x234266d7.
//
// Solidity: function donate((address,address,uint24,int24,address) key, uint256 amount0, uint256 amount1, bytes hookData) returns(int256 delta)
func (_PoolManager *PoolManagerTransactor) Donate(opts *bind.TransactOpts, key PoolKey, amount0 *big.Int, amount1 *big.Int, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "donate", key, amount0, amount1, hookData)
}

// Donate is a paid mutator transaction binding the contract method 0x234266d7.
//
// Solidity: function donate((address,address,uint24,int24,address) key, uint256 amount0, uint256 amount1, bytes hookData) returns(int256 delta)
func (_PoolManager *PoolManagerSession) Donate(key PoolKey, amount0 *big.Int, amount1 *big.Int, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.Donate(&_PoolManager.TransactOpts, key, amount0, amount1, hookData)
}

// Donate is a paid mutator transaction binding the contract method 0x234266d7.
//
// Solidity: function donate((address,address,uint24,int24,address) key, uint256 amount0, uint256 amount1, bytes hookData) returns(int256 delta)
func (_PoolManager *PoolManagerTransactorSession) Donate(key PoolKey, amount0 *big.Int, amount1 *big.Int, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.Donate(&_PoolManager.TransactOpts, key, amount0, amount1, hookData)
}

// Initialize is a paid mutator transaction binding the contract method 0x6276cbbe.
//
// Solidity: function initialize((address,address,uint24,int24,address) key, uint160 sqrtPriceX96) returns(int24 tick)
func (_PoolManager *PoolManagerTransactor) Initialize(opts *bind.TransactOpts, key PoolKey, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "initialize", key, sqrtPriceX96)
}

// Initialize is a paid mutator transaction binding the contract method 0x6276cbbe.
//
// Solidity: function initialize((address,address,uint24,int24,address) key, uint160 sqrtPriceX96) returns(int24 tick)
func (_PoolManager *PoolManagerSession) Initialize(key PoolKey, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Initialize(&_PoolManager.TransactOpts, key, sqrtPriceX96)
}

// Initialize is a paid mutator transaction binding the contract method 0x6276cbbe.
//
// Solidity: function initialize((address,address,uint24,int24,address) key, uint160 sqrtPriceX96) returns(int24 tick)
func (_PoolManager *PoolManagerTransactorSession) Initialize(key PoolKey, sqrtPriceX96 *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Initialize(&_PoolManager.TransactOpts, key, sqrtPriceX96)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactor) Mint(opts *bind.TransactOpts, to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "mint", to, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (_PoolManager *PoolManagerSession) Mint(to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Mint(&_PoolManager.TransactOpts, to, id, amount)
}

// Mint is a paid mutator transaction binding the contract method 0x156e29f6.
//
// Solidity: function mint(address to, uint256 id, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactorSession) Mint(to common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Mint(&_PoolManager.TransactOpts, to, id, amount)
}

// ModifyLiquidity is a paid mutator transaction binding the contract method 0x5a6bcfda.
//
// Solidity: function modifyLiquidity((address,address,uint24,int24,address) key, (int24,int24,int256,bytes32) params, bytes hookData) returns(int256 callerDelta, int256 feesAccrued)
func (_PoolManager *PoolManagerTransactor) ModifyLiquidity(opts *bind.TransactOpts, key PoolKey, params IPoolManagerModifyLiquidityParams, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "modifyLiquidity", key, params, hookData)
}

// ModifyLiquidity is a paid mutator transaction binding the contract method 0x5a6bcfda.
//
// Solidity: function modifyLiquidity((address,address,uint24,int24,address) key, (int24,int24,int256,bytes32) params, bytes hookData) returns(int256 callerDelta, int256 feesAccrued)
func (_PoolManager *PoolManagerSession) ModifyLiquidity(key PoolKey, params IPoolManagerModifyLiquidityParams, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.ModifyLiquidity(&_PoolManager.TransactOpts, key, params, hookData)
}

// ModifyLiquidity is a paid mutator transaction binding the contract method 0x5a6bcfda.
//
// Solidity: function modifyLiquidity((address,address,uint24,int24,address) key, (int24,int24,int256,bytes32) params, bytes hookData) returns(int256 callerDelta, int256 feesAccrued)
func (_PoolManager *PoolManagerTransactorSession) ModifyLiquidity(key PoolKey, params IPoolManagerModifyLiquidityParams, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.ModifyLiquidity(&_PoolManager.TransactOpts, key, params, hookData)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(address operator, bool approved) returns(bool)
func (_PoolManager *PoolManagerTransactor) SetOperator(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "setOperator", operator, approved)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(address operator, bool approved) returns(bool)
func (_PoolManager *PoolManagerSession) SetOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _PoolManager.Contract.SetOperator(&_PoolManager.TransactOpts, operator, approved)
}

// SetOperator is a paid mutator transaction binding the contract method 0x558a7297.
//
// Solidity: function setOperator(address operator, bool approved) returns(bool)
func (_PoolManager *PoolManagerTransactorSession) SetOperator(operator common.Address, approved bool) (*types.Transaction, error) {
	return _PoolManager.Contract.SetOperator(&_PoolManager.TransactOpts, operator, approved)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x7e87ce7d.
//
// Solidity: function setProtocolFee((address,address,uint24,int24,address) key, uint24 newProtocolFee) returns()
func (_PoolManager *PoolManagerTransactor) SetProtocolFee(opts *bind.TransactOpts, key PoolKey, newProtocolFee *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "setProtocolFee", key, newProtocolFee)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x7e87ce7d.
//
// Solidity: function setProtocolFee((address,address,uint24,int24,address) key, uint24 newProtocolFee) returns()
func (_PoolManager *PoolManagerSession) SetProtocolFee(key PoolKey, newProtocolFee *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.SetProtocolFee(&_PoolManager.TransactOpts, key, newProtocolFee)
}

// SetProtocolFee is a paid mutator transaction binding the contract method 0x7e87ce7d.
//
// Solidity: function setProtocolFee((address,address,uint24,int24,address) key, uint24 newProtocolFee) returns()
func (_PoolManager *PoolManagerTransactorSession) SetProtocolFee(key PoolKey, newProtocolFee *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.SetProtocolFee(&_PoolManager.TransactOpts, key, newProtocolFee)
}

// SetProtocolFeeController is a paid mutator transaction binding the contract method 0x2d771389.
//
// Solidity: function setProtocolFeeController(address controller) returns()
func (_PoolManager *PoolManagerTransactor) SetProtocolFeeController(opts *bind.TransactOpts, controller common.Address) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "setProtocolFeeController", controller)
}

// SetProtocolFeeController is a paid mutator transaction binding the contract method 0x2d771389.
//
// Solidity: function setProtocolFeeController(address controller) returns()
func (_PoolManager *PoolManagerSession) SetProtocolFeeController(controller common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.SetProtocolFeeController(&_PoolManager.TransactOpts, controller)
}

// SetProtocolFeeController is a paid mutator transaction binding the contract method 0x2d771389.
//
// Solidity: function setProtocolFeeController(address controller) returns()
func (_PoolManager *PoolManagerTransactorSession) SetProtocolFeeController(controller common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.SetProtocolFeeController(&_PoolManager.TransactOpts, controller)
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() payable returns(uint256)
func (_PoolManager *PoolManagerTransactor) Settle(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "settle")
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() payable returns(uint256)
func (_PoolManager *PoolManagerSession) Settle() (*types.Transaction, error) {
	return _PoolManager.Contract.Settle(&_PoolManager.TransactOpts)
}

// Settle is a paid mutator transaction binding the contract method 0x11da60b4.
//
// Solidity: function settle() payable returns(uint256)
func (_PoolManager *PoolManagerTransactorSession) Settle() (*types.Transaction, error) {
	return _PoolManager.Contract.Settle(&_PoolManager.TransactOpts)
}

// SettleFor is a paid mutator transaction binding the contract method 0x3dd45adb.
//
// Solidity: function settleFor(address recipient) payable returns(uint256)
func (_PoolManager *PoolManagerTransactor) SettleFor(opts *bind.TransactOpts, recipient common.Address) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "settleFor", recipient)
}

// SettleFor is a paid mutator transaction binding the contract method 0x3dd45adb.
//
// Solidity: function settleFor(address recipient) payable returns(uint256)
func (_PoolManager *PoolManagerSession) SettleFor(recipient common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.SettleFor(&_PoolManager.TransactOpts, recipient)
}

// SettleFor is a paid mutator transaction binding the contract method 0x3dd45adb.
//
// Solidity: function settleFor(address recipient) payable returns(uint256)
func (_PoolManager *PoolManagerTransactorSession) SettleFor(recipient common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.SettleFor(&_PoolManager.TransactOpts, recipient)
}

// Swap is a paid mutator transaction binding the contract method 0xf3cd914c.
//
// Solidity: function swap((address,address,uint24,int24,address) key, (bool,int256,uint160) params, bytes hookData) returns(int256 swapDelta)
func (_PoolManager *PoolManagerTransactor) Swap(opts *bind.TransactOpts, key PoolKey, params IPoolManagerSwapParams, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "swap", key, params, hookData)
}

// Swap is a paid mutator transaction binding the contract method 0xf3cd914c.
//
// Solidity: function swap((address,address,uint24,int24,address) key, (bool,int256,uint160) params, bytes hookData) returns(int256 swapDelta)
func (_PoolManager *PoolManagerSession) Swap(key PoolKey, params IPoolManagerSwapParams, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.Swap(&_PoolManager.TransactOpts, key, params, hookData)
}

// Swap is a paid mutator transaction binding the contract method 0xf3cd914c.
//
// Solidity: function swap((address,address,uint24,int24,address) key, (bool,int256,uint160) params, bytes hookData) returns(int256 swapDelta)
func (_PoolManager *PoolManagerTransactorSession) Swap(key PoolKey, params IPoolManagerSwapParams, hookData []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.Swap(&_PoolManager.TransactOpts, key, params, hookData)
}

// Sync is a paid mutator transaction binding the contract method 0xa5841194.
//
// Solidity: function sync(address currency) returns()
func (_PoolManager *PoolManagerTransactor) Sync(opts *bind.TransactOpts, currency common.Address) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "sync", currency)
}

// Sync is a paid mutator transaction binding the contract method 0xa5841194.
//
// Solidity: function sync(address currency) returns()
func (_PoolManager *PoolManagerSession) Sync(currency common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.Sync(&_PoolManager.TransactOpts, currency)
}

// Sync is a paid mutator transaction binding the contract method 0xa5841194.
//
// Solidity: function sync(address currency) returns()
func (_PoolManager *PoolManagerTransactorSession) Sync(currency common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.Sync(&_PoolManager.TransactOpts, currency)
}

// Take is a paid mutator transaction binding the contract method 0x0b0d9c09.
//
// Solidity: function take(address currency, address to, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactor) Take(opts *bind.TransactOpts, currency common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "take", currency, to, amount)
}

// Take is a paid mutator transaction binding the contract method 0x0b0d9c09.
//
// Solidity: function take(address currency, address to, uint256 amount) returns()
func (_PoolManager *PoolManagerSession) Take(currency common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Take(&_PoolManager.TransactOpts, currency, to, amount)
}

// Take is a paid mutator transaction binding the contract method 0x0b0d9c09.
//
// Solidity: function take(address currency, address to, uint256 amount) returns()
func (_PoolManager *PoolManagerTransactorSession) Take(currency common.Address, to common.Address, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Take(&_PoolManager.TransactOpts, currency, to, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x095bcdb6.
//
// Solidity: function transfer(address receiver, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerTransactor) Transfer(opts *bind.TransactOpts, receiver common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "transfer", receiver, id, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x095bcdb6.
//
// Solidity: function transfer(address receiver, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerSession) Transfer(receiver common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Transfer(&_PoolManager.TransactOpts, receiver, id, amount)
}

// Transfer is a paid mutator transaction binding the contract method 0x095bcdb6.
//
// Solidity: function transfer(address receiver, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerTransactorSession) Transfer(receiver common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.Transfer(&_PoolManager.TransactOpts, receiver, id, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xfe99049a.
//
// Solidity: function transferFrom(address sender, address receiver, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerTransactor) TransferFrom(opts *bind.TransactOpts, sender common.Address, receiver common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "transferFrom", sender, receiver, id, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xfe99049a.
//
// Solidity: function transferFrom(address sender, address receiver, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerSession) TransferFrom(sender common.Address, receiver common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.TransferFrom(&_PoolManager.TransactOpts, sender, receiver, id, amount)
}

// TransferFrom is a paid mutator transaction binding the contract method 0xfe99049a.
//
// Solidity: function transferFrom(address sender, address receiver, uint256 id, uint256 amount) returns(bool)
func (_PoolManager *PoolManagerTransactorSession) TransferFrom(sender common.Address, receiver common.Address, id *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.TransferFrom(&_PoolManager.TransactOpts, sender, receiver, id, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PoolManager *PoolManagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PoolManager *PoolManagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.TransferOwnership(&_PoolManager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_PoolManager *PoolManagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _PoolManager.Contract.TransferOwnership(&_PoolManager.TransactOpts, newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0x48c89491.
//
// Solidity: function unlock(bytes data) returns(bytes result)
func (_PoolManager *PoolManagerTransactor) Unlock(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "unlock", data)
}

// Unlock is a paid mutator transaction binding the contract method 0x48c89491.
//
// Solidity: function unlock(bytes data) returns(bytes result)
func (_PoolManager *PoolManagerSession) Unlock(data []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.Unlock(&_PoolManager.TransactOpts, data)
}

// Unlock is a paid mutator transaction binding the contract method 0x48c89491.
//
// Solidity: function unlock(bytes data) returns(bytes result)
func (_PoolManager *PoolManagerTransactorSession) Unlock(data []byte) (*types.Transaction, error) {
	return _PoolManager.Contract.Unlock(&_PoolManager.TransactOpts, data)
}

// UpdateDynamicLPFee is a paid mutator transaction binding the contract method 0x52759651.
//
// Solidity: function updateDynamicLPFee((address,address,uint24,int24,address) key, uint24 newDynamicLPFee) returns()
func (_PoolManager *PoolManagerTransactor) UpdateDynamicLPFee(opts *bind.TransactOpts, key PoolKey, newDynamicLPFee *big.Int) (*types.Transaction, error) {
	return _PoolManager.contract.Transact(opts, "updateDynamicLPFee", key, newDynamicLPFee)
}

// UpdateDynamicLPFee is a paid mutator transaction binding the contract method 0x52759651.
//
// Solidity: function updateDynamicLPFee((address,address,uint24,int24,address) key, uint24 newDynamicLPFee) returns()
func (_PoolManager *PoolManagerSession) UpdateDynamicLPFee(key PoolKey, newDynamicLPFee *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.UpdateDynamicLPFee(&_PoolManager.TransactOpts, key, newDynamicLPFee)
}

// UpdateDynamicLPFee is a paid mutator transaction binding the contract method 0x52759651.
//
// Solidity: function updateDynamicLPFee((address,address,uint24,int24,address) key, uint24 newDynamicLPFee) returns()
func (_PoolManager *PoolManagerTransactorSession) UpdateDynamicLPFee(key PoolKey, newDynamicLPFee *big.Int) (*types.Transaction, error) {
	return _PoolManager.Contract.UpdateDynamicLPFee(&_PoolManager.TransactOpts, key, newDynamicLPFee)
}

// PoolManagerApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the PoolManager contract.
type PoolManagerApprovalIterator struct {
	Event *PoolManagerApproval // Event containing the contract specifics and raw log

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
func (it *PoolManagerApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerApproval)
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
		it.Event = new(PoolManagerApproval)
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
func (it *PoolManagerApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerApproval represents a Approval event raised by the PoolManager contract.
type PoolManagerApproval struct {
	Owner   common.Address
	Spender common.Address
	Id      *big.Int
	Amount  *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0xb3fd5071835887567a0671151121894ddccc2842f1d10bedad13e0d17cace9a7.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 indexed id, uint256 amount)
func (_PoolManager *PoolManagerFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, spender []common.Address, id []*big.Int) (*PoolManagerApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "Approval", ownerRule, spenderRule, idRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerApprovalIterator{contract: _PoolManager.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0xb3fd5071835887567a0671151121894ddccc2842f1d10bedad13e0d17cace9a7.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 indexed id, uint256 amount)
func (_PoolManager *PoolManagerFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *PoolManagerApproval, owner []common.Address, spender []common.Address, id []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var spenderRule []interface{}
	for _, spenderItem := range spender {
		spenderRule = append(spenderRule, spenderItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "Approval", ownerRule, spenderRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerApproval)
				if err := _PoolManager.contract.UnpackLog(event, "Approval", log); err != nil {
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

// ParseApproval is a log parse operation binding the contract event 0xb3fd5071835887567a0671151121894ddccc2842f1d10bedad13e0d17cace9a7.
//
// Solidity: event Approval(address indexed owner, address indexed spender, uint256 indexed id, uint256 amount)
func (_PoolManager *PoolManagerFilterer) ParseApproval(log types.Log) (*PoolManagerApproval, error) {
	event := new(PoolManagerApproval)
	if err := _PoolManager.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerDonateIterator is returned from FilterDonate and is used to iterate over the raw logs and unpacked data for Donate events raised by the PoolManager contract.
type PoolManagerDonateIterator struct {
	Event *PoolManagerDonate // Event containing the contract specifics and raw log

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
func (it *PoolManagerDonateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerDonate)
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
		it.Event = new(PoolManagerDonate)
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
func (it *PoolManagerDonateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerDonateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerDonate represents a Donate event raised by the PoolManager contract.
type PoolManagerDonate struct {
	Id      [32]byte
	Sender  common.Address
	Amount0 *big.Int
	Amount1 *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDonate is a free log retrieval operation binding the contract event 0x29ef05caaff9404b7cb6d1c0e9bbae9eaa7ab2541feba1a9c4248594c08156cb.
//
// Solidity: event Donate(bytes32 indexed id, address indexed sender, uint256 amount0, uint256 amount1)
func (_PoolManager *PoolManagerFilterer) FilterDonate(opts *bind.FilterOpts, id [][32]byte, sender []common.Address) (*PoolManagerDonateIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "Donate", idRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerDonateIterator{contract: _PoolManager.contract, event: "Donate", logs: logs, sub: sub}, nil
}

// WatchDonate is a free log subscription operation binding the contract event 0x29ef05caaff9404b7cb6d1c0e9bbae9eaa7ab2541feba1a9c4248594c08156cb.
//
// Solidity: event Donate(bytes32 indexed id, address indexed sender, uint256 amount0, uint256 amount1)
func (_PoolManager *PoolManagerFilterer) WatchDonate(opts *bind.WatchOpts, sink chan<- *PoolManagerDonate, id [][32]byte, sender []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "Donate", idRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerDonate)
				if err := _PoolManager.contract.UnpackLog(event, "Donate", log); err != nil {
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

// ParseDonate is a log parse operation binding the contract event 0x29ef05caaff9404b7cb6d1c0e9bbae9eaa7ab2541feba1a9c4248594c08156cb.
//
// Solidity: event Donate(bytes32 indexed id, address indexed sender, uint256 amount0, uint256 amount1)
func (_PoolManager *PoolManagerFilterer) ParseDonate(log types.Log) (*PoolManagerDonate, error) {
	event := new(PoolManagerDonate)
	if err := _PoolManager.contract.UnpackLog(event, "Donate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerInitializeIterator is returned from FilterInitialize and is used to iterate over the raw logs and unpacked data for Initialize events raised by the PoolManager contract.
type PoolManagerInitializeIterator struct {
	Event *PoolManagerInitialize // Event containing the contract specifics and raw log

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
func (it *PoolManagerInitializeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerInitialize)
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
		it.Event = new(PoolManagerInitialize)
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
func (it *PoolManagerInitializeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerInitializeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerInitialize represents a Initialize event raised by the PoolManager contract.
type PoolManagerInitialize struct {
	Id           [32]byte
	Currency0    common.Address
	Currency1    common.Address
	Fee          *big.Int
	TickSpacing  *big.Int
	Hooks        common.Address
	SqrtPriceX96 *big.Int
	Tick         *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterInitialize is a free log retrieval operation binding the contract event 0xdd466e674ea557f56295e2d0218a125ea4b4f0f6f3307b95f85e6110838d6438.
//
// Solidity: event Initialize(bytes32 indexed id, address indexed currency0, address indexed currency1, uint24 fee, int24 tickSpacing, address hooks, uint160 sqrtPriceX96, int24 tick)
func (_PoolManager *PoolManagerFilterer) FilterInitialize(opts *bind.FilterOpts, id [][32]byte, currency0 []common.Address, currency1 []common.Address) (*PoolManagerInitializeIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var currency0Rule []interface{}
	for _, currency0Item := range currency0 {
		currency0Rule = append(currency0Rule, currency0Item)
	}
	var currency1Rule []interface{}
	for _, currency1Item := range currency1 {
		currency1Rule = append(currency1Rule, currency1Item)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "Initialize", idRule, currency0Rule, currency1Rule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerInitializeIterator{contract: _PoolManager.contract, event: "Initialize", logs: logs, sub: sub}, nil
}

// WatchInitialize is a free log subscription operation binding the contract event 0xdd466e674ea557f56295e2d0218a125ea4b4f0f6f3307b95f85e6110838d6438.
//
// Solidity: event Initialize(bytes32 indexed id, address indexed currency0, address indexed currency1, uint24 fee, int24 tickSpacing, address hooks, uint160 sqrtPriceX96, int24 tick)
func (_PoolManager *PoolManagerFilterer) WatchInitialize(opts *bind.WatchOpts, sink chan<- *PoolManagerInitialize, id [][32]byte, currency0 []common.Address, currency1 []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var currency0Rule []interface{}
	for _, currency0Item := range currency0 {
		currency0Rule = append(currency0Rule, currency0Item)
	}
	var currency1Rule []interface{}
	for _, currency1Item := range currency1 {
		currency1Rule = append(currency1Rule, currency1Item)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "Initialize", idRule, currency0Rule, currency1Rule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerInitialize)
				if err := _PoolManager.contract.UnpackLog(event, "Initialize", log); err != nil {
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

// ParseInitialize is a log parse operation binding the contract event 0xdd466e674ea557f56295e2d0218a125ea4b4f0f6f3307b95f85e6110838d6438.
//
// Solidity: event Initialize(bytes32 indexed id, address indexed currency0, address indexed currency1, uint24 fee, int24 tickSpacing, address hooks, uint160 sqrtPriceX96, int24 tick)
func (_PoolManager *PoolManagerFilterer) ParseInitialize(log types.Log) (*PoolManagerInitialize, error) {
	event := new(PoolManagerInitialize)
	if err := _PoolManager.contract.UnpackLog(event, "Initialize", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerModifyLiquidityIterator is returned from FilterModifyLiquidity and is used to iterate over the raw logs and unpacked data for ModifyLiquidity events raised by the PoolManager contract.
type PoolManagerModifyLiquidityIterator struct {
	Event *PoolManagerModifyLiquidity // Event containing the contract specifics and raw log

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
func (it *PoolManagerModifyLiquidityIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerModifyLiquidity)
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
		it.Event = new(PoolManagerModifyLiquidity)
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
func (it *PoolManagerModifyLiquidityIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerModifyLiquidityIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerModifyLiquidity represents a ModifyLiquidity event raised by the PoolManager contract.
type PoolManagerModifyLiquidity struct {
	Id             [32]byte
	Sender         common.Address
	TickLower      *big.Int
	TickUpper      *big.Int
	LiquidityDelta *big.Int
	Salt           [32]byte
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterModifyLiquidity is a free log retrieval operation binding the contract event 0xf208f4912782fd25c7f114ca3723a2d5dd6f3bcc3ac8db5af63baa85f711d5ec.
//
// Solidity: event ModifyLiquidity(bytes32 indexed id, address indexed sender, int24 tickLower, int24 tickUpper, int256 liquidityDelta, bytes32 salt)
func (_PoolManager *PoolManagerFilterer) FilterModifyLiquidity(opts *bind.FilterOpts, id [][32]byte, sender []common.Address) (*PoolManagerModifyLiquidityIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "ModifyLiquidity", idRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerModifyLiquidityIterator{contract: _PoolManager.contract, event: "ModifyLiquidity", logs: logs, sub: sub}, nil
}

// WatchModifyLiquidity is a free log subscription operation binding the contract event 0xf208f4912782fd25c7f114ca3723a2d5dd6f3bcc3ac8db5af63baa85f711d5ec.
//
// Solidity: event ModifyLiquidity(bytes32 indexed id, address indexed sender, int24 tickLower, int24 tickUpper, int256 liquidityDelta, bytes32 salt)
func (_PoolManager *PoolManagerFilterer) WatchModifyLiquidity(opts *bind.WatchOpts, sink chan<- *PoolManagerModifyLiquidity, id [][32]byte, sender []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "ModifyLiquidity", idRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerModifyLiquidity)
				if err := _PoolManager.contract.UnpackLog(event, "ModifyLiquidity", log); err != nil {
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

// ParseModifyLiquidity is a log parse operation binding the contract event 0xf208f4912782fd25c7f114ca3723a2d5dd6f3bcc3ac8db5af63baa85f711d5ec.
//
// Solidity: event ModifyLiquidity(bytes32 indexed id, address indexed sender, int24 tickLower, int24 tickUpper, int256 liquidityDelta, bytes32 salt)
func (_PoolManager *PoolManagerFilterer) ParseModifyLiquidity(log types.Log) (*PoolManagerModifyLiquidity, error) {
	event := new(PoolManagerModifyLiquidity)
	if err := _PoolManager.contract.UnpackLog(event, "ModifyLiquidity", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerOperatorSetIterator is returned from FilterOperatorSet and is used to iterate over the raw logs and unpacked data for OperatorSet events raised by the PoolManager contract.
type PoolManagerOperatorSetIterator struct {
	Event *PoolManagerOperatorSet // Event containing the contract specifics and raw log

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
func (it *PoolManagerOperatorSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerOperatorSet)
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
		it.Event = new(PoolManagerOperatorSet)
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
func (it *PoolManagerOperatorSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerOperatorSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerOperatorSet represents a OperatorSet event raised by the PoolManager contract.
type PoolManagerOperatorSet struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOperatorSet is a free log retrieval operation binding the contract event 0xceb576d9f15e4e200fdb5096d64d5dfd667e16def20c1eefd14256d8e3faa267.
//
// Solidity: event OperatorSet(address indexed owner, address indexed operator, bool approved)
func (_PoolManager *PoolManagerFilterer) FilterOperatorSet(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*PoolManagerOperatorSetIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "OperatorSet", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerOperatorSetIterator{contract: _PoolManager.contract, event: "OperatorSet", logs: logs, sub: sub}, nil
}

// WatchOperatorSet is a free log subscription operation binding the contract event 0xceb576d9f15e4e200fdb5096d64d5dfd667e16def20c1eefd14256d8e3faa267.
//
// Solidity: event OperatorSet(address indexed owner, address indexed operator, bool approved)
func (_PoolManager *PoolManagerFilterer) WatchOperatorSet(opts *bind.WatchOpts, sink chan<- *PoolManagerOperatorSet, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "OperatorSet", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerOperatorSet)
				if err := _PoolManager.contract.UnpackLog(event, "OperatorSet", log); err != nil {
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

// ParseOperatorSet is a log parse operation binding the contract event 0xceb576d9f15e4e200fdb5096d64d5dfd667e16def20c1eefd14256d8e3faa267.
//
// Solidity: event OperatorSet(address indexed owner, address indexed operator, bool approved)
func (_PoolManager *PoolManagerFilterer) ParseOperatorSet(log types.Log) (*PoolManagerOperatorSet, error) {
	event := new(PoolManagerOperatorSet)
	if err := _PoolManager.contract.UnpackLog(event, "OperatorSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the PoolManager contract.
type PoolManagerOwnershipTransferredIterator struct {
	Event *PoolManagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *PoolManagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerOwnershipTransferred)
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
		it.Event = new(PoolManagerOwnershipTransferred)
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
func (it *PoolManagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerOwnershipTransferred represents a OwnershipTransferred event raised by the PoolManager contract.
type PoolManagerOwnershipTransferred struct {
	User     common.Address
	NewOwner common.Address
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed user, address indexed newOwner)
func (_PoolManager *PoolManagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, user []common.Address, newOwner []common.Address) (*PoolManagerOwnershipTransferredIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "OwnershipTransferred", userRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerOwnershipTransferredIterator{contract: _PoolManager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed user, address indexed newOwner)
func (_PoolManager *PoolManagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *PoolManagerOwnershipTransferred, user []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "OwnershipTransferred", userRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerOwnershipTransferred)
				if err := _PoolManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
// Solidity: event OwnershipTransferred(address indexed user, address indexed newOwner)
func (_PoolManager *PoolManagerFilterer) ParseOwnershipTransferred(log types.Log) (*PoolManagerOwnershipTransferred, error) {
	event := new(PoolManagerOwnershipTransferred)
	if err := _PoolManager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerProtocolFeeControllerUpdatedIterator is returned from FilterProtocolFeeControllerUpdated and is used to iterate over the raw logs and unpacked data for ProtocolFeeControllerUpdated events raised by the PoolManager contract.
type PoolManagerProtocolFeeControllerUpdatedIterator struct {
	Event *PoolManagerProtocolFeeControllerUpdated // Event containing the contract specifics and raw log

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
func (it *PoolManagerProtocolFeeControllerUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerProtocolFeeControllerUpdated)
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
		it.Event = new(PoolManagerProtocolFeeControllerUpdated)
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
func (it *PoolManagerProtocolFeeControllerUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerProtocolFeeControllerUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerProtocolFeeControllerUpdated represents a ProtocolFeeControllerUpdated event raised by the PoolManager contract.
type PoolManagerProtocolFeeControllerUpdated struct {
	ProtocolFeeController common.Address
	Raw                   types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeControllerUpdated is a free log retrieval operation binding the contract event 0xb4bd8ef53df690b9943d3318996006dbb82a25f54719d8c8035b516a2a5b8acc.
//
// Solidity: event ProtocolFeeControllerUpdated(address indexed protocolFeeController)
func (_PoolManager *PoolManagerFilterer) FilterProtocolFeeControllerUpdated(opts *bind.FilterOpts, protocolFeeController []common.Address) (*PoolManagerProtocolFeeControllerUpdatedIterator, error) {

	var protocolFeeControllerRule []interface{}
	for _, protocolFeeControllerItem := range protocolFeeController {
		protocolFeeControllerRule = append(protocolFeeControllerRule, protocolFeeControllerItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "ProtocolFeeControllerUpdated", protocolFeeControllerRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerProtocolFeeControllerUpdatedIterator{contract: _PoolManager.contract, event: "ProtocolFeeControllerUpdated", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeControllerUpdated is a free log subscription operation binding the contract event 0xb4bd8ef53df690b9943d3318996006dbb82a25f54719d8c8035b516a2a5b8acc.
//
// Solidity: event ProtocolFeeControllerUpdated(address indexed protocolFeeController)
func (_PoolManager *PoolManagerFilterer) WatchProtocolFeeControllerUpdated(opts *bind.WatchOpts, sink chan<- *PoolManagerProtocolFeeControllerUpdated, protocolFeeController []common.Address) (event.Subscription, error) {

	var protocolFeeControllerRule []interface{}
	for _, protocolFeeControllerItem := range protocolFeeController {
		protocolFeeControllerRule = append(protocolFeeControllerRule, protocolFeeControllerItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "ProtocolFeeControllerUpdated", protocolFeeControllerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerProtocolFeeControllerUpdated)
				if err := _PoolManager.contract.UnpackLog(event, "ProtocolFeeControllerUpdated", log); err != nil {
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

// ParseProtocolFeeControllerUpdated is a log parse operation binding the contract event 0xb4bd8ef53df690b9943d3318996006dbb82a25f54719d8c8035b516a2a5b8acc.
//
// Solidity: event ProtocolFeeControllerUpdated(address indexed protocolFeeController)
func (_PoolManager *PoolManagerFilterer) ParseProtocolFeeControllerUpdated(log types.Log) (*PoolManagerProtocolFeeControllerUpdated, error) {
	event := new(PoolManagerProtocolFeeControllerUpdated)
	if err := _PoolManager.contract.UnpackLog(event, "ProtocolFeeControllerUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerProtocolFeeUpdatedIterator is returned from FilterProtocolFeeUpdated and is used to iterate over the raw logs and unpacked data for ProtocolFeeUpdated events raised by the PoolManager contract.
type PoolManagerProtocolFeeUpdatedIterator struct {
	Event *PoolManagerProtocolFeeUpdated // Event containing the contract specifics and raw log

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
func (it *PoolManagerProtocolFeeUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerProtocolFeeUpdated)
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
		it.Event = new(PoolManagerProtocolFeeUpdated)
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
func (it *PoolManagerProtocolFeeUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerProtocolFeeUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerProtocolFeeUpdated represents a ProtocolFeeUpdated event raised by the PoolManager contract.
type PoolManagerProtocolFeeUpdated struct {
	Id          [32]byte
	ProtocolFee *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterProtocolFeeUpdated is a free log retrieval operation binding the contract event 0xe9c42593e71f84403b84352cd168d693e2c9fcd1fdbcc3feb21d92b43e6696f9.
//
// Solidity: event ProtocolFeeUpdated(bytes32 indexed id, uint24 protocolFee)
func (_PoolManager *PoolManagerFilterer) FilterProtocolFeeUpdated(opts *bind.FilterOpts, id [][32]byte) (*PoolManagerProtocolFeeUpdatedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "ProtocolFeeUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerProtocolFeeUpdatedIterator{contract: _PoolManager.contract, event: "ProtocolFeeUpdated", logs: logs, sub: sub}, nil
}

// WatchProtocolFeeUpdated is a free log subscription operation binding the contract event 0xe9c42593e71f84403b84352cd168d693e2c9fcd1fdbcc3feb21d92b43e6696f9.
//
// Solidity: event ProtocolFeeUpdated(bytes32 indexed id, uint24 protocolFee)
func (_PoolManager *PoolManagerFilterer) WatchProtocolFeeUpdated(opts *bind.WatchOpts, sink chan<- *PoolManagerProtocolFeeUpdated, id [][32]byte) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "ProtocolFeeUpdated", idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerProtocolFeeUpdated)
				if err := _PoolManager.contract.UnpackLog(event, "ProtocolFeeUpdated", log); err != nil {
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

// ParseProtocolFeeUpdated is a log parse operation binding the contract event 0xe9c42593e71f84403b84352cd168d693e2c9fcd1fdbcc3feb21d92b43e6696f9.
//
// Solidity: event ProtocolFeeUpdated(bytes32 indexed id, uint24 protocolFee)
func (_PoolManager *PoolManagerFilterer) ParseProtocolFeeUpdated(log types.Log) (*PoolManagerProtocolFeeUpdated, error) {
	event := new(PoolManagerProtocolFeeUpdated)
	if err := _PoolManager.contract.UnpackLog(event, "ProtocolFeeUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerSwapIterator is returned from FilterSwap and is used to iterate over the raw logs and unpacked data for Swap events raised by the PoolManager contract.
type PoolManagerSwapIterator struct {
	Event *PoolManagerSwap // Event containing the contract specifics and raw log

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
func (it *PoolManagerSwapIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerSwap)
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
		it.Event = new(PoolManagerSwap)
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
func (it *PoolManagerSwapIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerSwapIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerSwap represents a Swap event raised by the PoolManager contract.
type PoolManagerSwap struct {
	Id           [32]byte
	Sender       common.Address
	Amount0      *big.Int
	Amount1      *big.Int
	SqrtPriceX96 *big.Int
	Liquidity    *big.Int
	Tick         *big.Int
	Fee          *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterSwap is a free log retrieval operation binding the contract event 0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f.
//
// Solidity: event Swap(bytes32 indexed id, address indexed sender, int128 amount0, int128 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick, uint24 fee)
func (_PoolManager *PoolManagerFilterer) FilterSwap(opts *bind.FilterOpts, id [][32]byte, sender []common.Address) (*PoolManagerSwapIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "Swap", idRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerSwapIterator{contract: _PoolManager.contract, event: "Swap", logs: logs, sub: sub}, nil
}

// WatchSwap is a free log subscription operation binding the contract event 0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f.
//
// Solidity: event Swap(bytes32 indexed id, address indexed sender, int128 amount0, int128 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick, uint24 fee)
func (_PoolManager *PoolManagerFilterer) WatchSwap(opts *bind.WatchOpts, sink chan<- *PoolManagerSwap, id [][32]byte, sender []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "Swap", idRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerSwap)
				if err := _PoolManager.contract.UnpackLog(event, "Swap", log); err != nil {
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

// ParseSwap is a log parse operation binding the contract event 0x40e9cecb9f5f1f1c5b9c97dec2917b7ee92e57ba5563708daca94dd84ad7112f.
//
// Solidity: event Swap(bytes32 indexed id, address indexed sender, int128 amount0, int128 amount1, uint160 sqrtPriceX96, uint128 liquidity, int24 tick, uint24 fee)
func (_PoolManager *PoolManagerFilterer) ParseSwap(log types.Log) (*PoolManagerSwap, error) {
	event := new(PoolManagerSwap)
	if err := _PoolManager.contract.UnpackLog(event, "Swap", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// PoolManagerTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the PoolManager contract.
type PoolManagerTransferIterator struct {
	Event *PoolManagerTransfer // Event containing the contract specifics and raw log

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
func (it *PoolManagerTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(PoolManagerTransfer)
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
		it.Event = new(PoolManagerTransfer)
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
func (it *PoolManagerTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *PoolManagerTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// PoolManagerTransfer represents a Transfer event raised by the PoolManager contract.
type PoolManagerTransfer struct {
	Caller common.Address
	From   common.Address
	To     common.Address
	Id     *big.Int
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0x1b3d7edb2e9c0b0e7c525b20aaaef0f5940d2ed71663c7d39266ecafac728859.
//
// Solidity: event Transfer(address caller, address indexed from, address indexed to, uint256 indexed id, uint256 amount)
func (_PoolManager *PoolManagerFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, id []*big.Int) (*PoolManagerTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _PoolManager.contract.FilterLogs(opts, "Transfer", fromRule, toRule, idRule)
	if err != nil {
		return nil, err
	}
	return &PoolManagerTransferIterator{contract: _PoolManager.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0x1b3d7edb2e9c0b0e7c525b20aaaef0f5940d2ed71663c7d39266ecafac728859.
//
// Solidity: event Transfer(address caller, address indexed from, address indexed to, uint256 indexed id, uint256 amount)
func (_PoolManager *PoolManagerFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *PoolManagerTransfer, from []common.Address, to []common.Address, id []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}

	logs, sub, err := _PoolManager.contract.WatchLogs(opts, "Transfer", fromRule, toRule, idRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(PoolManagerTransfer)
				if err := _PoolManager.contract.UnpackLog(event, "Transfer", log); err != nil {
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

// ParseTransfer is a log parse operation binding the contract event 0x1b3d7edb2e9c0b0e7c525b20aaaef0f5940d2ed71663c7d39266ecafac728859.
//
// Solidity: event Transfer(address caller, address indexed from, address indexed to, uint256 indexed id, uint256 amount)
func (_PoolManager *PoolManagerFilterer) ParseTransfer(log types.Log) (*PoolManagerTransfer, error) {
	event := new(PoolManagerTransfer)
	if err := _PoolManager.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
