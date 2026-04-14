// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package diaERC1967Proxy

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

// DiaERC1967ProxyMetaData contains all meta data concerning the DiaERC1967Proxy contract.
var DiaERC1967ProxyMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"payable\"},{\"type\":\"fallback\",\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ERC1967ProxyUninitialized\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]}]",
	Bin: "0x60806040526102a2803803806100148161016e565b9283398101604082820312610156578151916001600160a01b03831690818403610156576020810151906001600160401b038211610156570182601f82011215610156578051906001600160401b03821161015a5761007c601f8301601f191660200161016e565b938285526020838301011161015657815f9260208093018387015e8401015281511561014757823b15610135577f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b031916821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a280511561011e5761010e91610193565b505b604051608290816102208239f35b505034156101105763b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b6330a289cf60e21b5f5260045ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6040519190601f01601f191682016001600160401b0381118382101761015a57604052565b905f8091602081519101845af4808061020c575b156101c75750506040513d81523d5f602083013e60203d82010160405290565b156101ec57639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b3d156101fd576040513d5f823e3d90fd5b63d6bda27560e01b5f5260045ffd5b503d1515806101a75750813b15156101a756fe60806040527f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc545f9081906001600160a01b0316368280378136915af43d5f803e156048573d5ff35b3d5ffdfea2646970667358221220244b850a7b9d27bc5726c6bbcd68471b04e720d13f77d0a76628f8c913eb108464736f6c63430008220033",
}

// DiaERC1967ProxyABI is the input ABI used to generate the binding from.
// Deprecated: Use DiaERC1967ProxyMetaData.ABI instead.
var DiaERC1967ProxyABI = DiaERC1967ProxyMetaData.ABI

// DiaERC1967ProxyBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DiaERC1967ProxyMetaData.Bin instead.
var DiaERC1967ProxyBin = DiaERC1967ProxyMetaData.Bin

// DeployDiaERC1967Proxy deploys a new Ethereum contract, binding an instance of DiaERC1967Proxy to it.
func DeployDiaERC1967Proxy(auth *bind.TransactOpts, backend bind.ContractBackend, implementation common.Address, _data []byte) (common.Address, *types.Transaction, *DiaERC1967Proxy, error) {
	parsed, err := DiaERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DiaERC1967ProxyBin), backend, implementation, _data)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DiaERC1967Proxy{DiaERC1967ProxyCaller: DiaERC1967ProxyCaller{contract: contract}, DiaERC1967ProxyTransactor: DiaERC1967ProxyTransactor{contract: contract}, DiaERC1967ProxyFilterer: DiaERC1967ProxyFilterer{contract: contract}}, nil
}

// DiaERC1967Proxy is an auto generated Go binding around an Ethereum contract.
type DiaERC1967Proxy struct {
	DiaERC1967ProxyCaller     // Read-only binding to the contract
	DiaERC1967ProxyTransactor // Write-only binding to the contract
	DiaERC1967ProxyFilterer   // Log filterer for contract events
}

// DiaERC1967ProxyCaller is an auto generated read-only Go binding around an Ethereum contract.
type DiaERC1967ProxyCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DiaERC1967ProxyTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DiaERC1967ProxyTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DiaERC1967ProxyFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DiaERC1967ProxyFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DiaERC1967ProxySession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DiaERC1967ProxySession struct {
	Contract     *DiaERC1967Proxy  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DiaERC1967ProxyCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DiaERC1967ProxyCallerSession struct {
	Contract *DiaERC1967ProxyCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// DiaERC1967ProxyTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DiaERC1967ProxyTransactorSession struct {
	Contract     *DiaERC1967ProxyTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// DiaERC1967ProxyRaw is an auto generated low-level Go binding around an Ethereum contract.
type DiaERC1967ProxyRaw struct {
	Contract *DiaERC1967Proxy // Generic contract binding to access the raw methods on
}

// DiaERC1967ProxyCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DiaERC1967ProxyCallerRaw struct {
	Contract *DiaERC1967ProxyCaller // Generic read-only contract binding to access the raw methods on
}

// DiaERC1967ProxyTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DiaERC1967ProxyTransactorRaw struct {
	Contract *DiaERC1967ProxyTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDiaERC1967Proxy creates a new instance of DiaERC1967Proxy, bound to a specific deployed contract.
func NewDiaERC1967Proxy(address common.Address, backend bind.ContractBackend) (*DiaERC1967Proxy, error) {
	contract, err := bindDiaERC1967Proxy(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DiaERC1967Proxy{DiaERC1967ProxyCaller: DiaERC1967ProxyCaller{contract: contract}, DiaERC1967ProxyTransactor: DiaERC1967ProxyTransactor{contract: contract}, DiaERC1967ProxyFilterer: DiaERC1967ProxyFilterer{contract: contract}}, nil
}

// NewDiaERC1967ProxyCaller creates a new read-only instance of DiaERC1967Proxy, bound to a specific deployed contract.
func NewDiaERC1967ProxyCaller(address common.Address, caller bind.ContractCaller) (*DiaERC1967ProxyCaller, error) {
	contract, err := bindDiaERC1967Proxy(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DiaERC1967ProxyCaller{contract: contract}, nil
}

// NewDiaERC1967ProxyTransactor creates a new write-only instance of DiaERC1967Proxy, bound to a specific deployed contract.
func NewDiaERC1967ProxyTransactor(address common.Address, transactor bind.ContractTransactor) (*DiaERC1967ProxyTransactor, error) {
	contract, err := bindDiaERC1967Proxy(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DiaERC1967ProxyTransactor{contract: contract}, nil
}

// NewDiaERC1967ProxyFilterer creates a new log filterer instance of DiaERC1967Proxy, bound to a specific deployed contract.
func NewDiaERC1967ProxyFilterer(address common.Address, filterer bind.ContractFilterer) (*DiaERC1967ProxyFilterer, error) {
	contract, err := bindDiaERC1967Proxy(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DiaERC1967ProxyFilterer{contract: contract}, nil
}

// bindDiaERC1967Proxy binds a generic wrapper to an already deployed contract.
func bindDiaERC1967Proxy(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DiaERC1967ProxyMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DiaERC1967Proxy *DiaERC1967ProxyRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DiaERC1967Proxy.Contract.DiaERC1967ProxyCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DiaERC1967Proxy *DiaERC1967ProxyRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DiaERC1967Proxy.Contract.DiaERC1967ProxyTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DiaERC1967Proxy *DiaERC1967ProxyRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DiaERC1967Proxy.Contract.DiaERC1967ProxyTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DiaERC1967Proxy *DiaERC1967ProxyCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DiaERC1967Proxy.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DiaERC1967Proxy *DiaERC1967ProxyTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DiaERC1967Proxy.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DiaERC1967Proxy *DiaERC1967ProxyTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DiaERC1967Proxy.Contract.contract.Transact(opts, method, params...)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_DiaERC1967Proxy *DiaERC1967ProxyTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _DiaERC1967Proxy.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_DiaERC1967Proxy *DiaERC1967ProxySession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _DiaERC1967Proxy.Contract.Fallback(&_DiaERC1967Proxy.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_DiaERC1967Proxy *DiaERC1967ProxyTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _DiaERC1967Proxy.Contract.Fallback(&_DiaERC1967Proxy.TransactOpts, calldata)
}

// DiaERC1967ProxyUpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the DiaERC1967Proxy contract.
type DiaERC1967ProxyUpgradedIterator struct {
	Event *DiaERC1967ProxyUpgraded // Event containing the contract specifics and raw log

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
func (it *DiaERC1967ProxyUpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaERC1967ProxyUpgraded)
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
		it.Event = new(DiaERC1967ProxyUpgraded)
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
func (it *DiaERC1967ProxyUpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaERC1967ProxyUpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaERC1967ProxyUpgraded represents a Upgraded event raised by the DiaERC1967Proxy contract.
type DiaERC1967ProxyUpgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DiaERC1967Proxy *DiaERC1967ProxyFilterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DiaERC1967ProxyUpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DiaERC1967Proxy.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DiaERC1967ProxyUpgradedIterator{contract: _DiaERC1967Proxy.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DiaERC1967Proxy *DiaERC1967ProxyFilterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DiaERC1967ProxyUpgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DiaERC1967Proxy.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaERC1967ProxyUpgraded)
				if err := _DiaERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
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

// ParseUpgraded is a log parse operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DiaERC1967Proxy *DiaERC1967ProxyFilterer) ParseUpgraded(log types.Log) (*DiaERC1967ProxyUpgraded, error) {
	event := new(DiaERC1967ProxyUpgraded)
	if err := _DiaERC1967Proxy.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
