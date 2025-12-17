package models

import (
	"encoding/json"
	"errors"
	"math"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	luminametacontract "github.com/diadata-org/lumina-library/contracts/lumina/metacontract"
	"github.com/diadata-org/lumina-library/utils"
)

const tokenABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"mintingFinished\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_spender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_from\",\"type\":\"address\"},{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unpause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"mint\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"paused\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finishMinting\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"pause\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"}],\"name\":\"transfer\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_to\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_releaseTime\",\"type\":\"uint256\"}],\"name\":\"mintTimelocked\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_spender\",\"type\":\"address\"}],\"name\":\"allowance\",\"outputs\":[{\"name\":\"remaining\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Mint\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"MintFinished\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Pause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[],\"name\":\"Unpause\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"spender\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"value\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"}]" //nolint:gosec

type TokenCaller struct {
	Contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Asset is the data type for all assets, ranging from fiat to crypto.
type Asset struct {
	Symbol     string `json:"Symbol"`
	Name       string `json:"Name"`
	Address    string `json:"Address"`
	Decimals   uint8  `json:"Decimals"`
	Blockchain string `json:"Blockchain"`
}

func (a *Asset) AssetIdentifier() string {
	return a.Blockchain + "-" + a.Address
}

// GetPrice returns the latest price of asset @a as published by DIA metacontract with @address.
// Fallback in case of error or 0 price is diadata API.
func (a *Asset) GetPrice(
	metacontractAddress common.Address,
	precision int,
	client *ethclient.Client,
) (assetQuotation AssetQuotation, err error) {

	assetQuotation, err = a.GetOnchainPrice(metacontractAddress, precision, client)
	if err != nil || assetQuotation.Price == 0 {
		log.Warnf("On-chain price not available for address-blockchain: %s -- %s", a.Blockchain, a.Address)
		return a.GetPriceFromDiaAPI()
	}

	// TO DO: Should we add a threshold such that onchain price is only taken into account if timestamp not older?
	assetQuotation.Source = "DIA_Lumina_" + metacontractAddress.Hex()

	return
}

// GetOnchainPrice returns a quotation for asset @a by querying the metacontract with @metacontractAddress.
func (a *Asset) GetOnchainPrice(
	metacontractAddress common.Address,
	precision int,
	client *ethclient.Client,
) (aq AssetQuotation, err error) {

	var caller *luminametacontract.LuminametacontractCaller
	caller, err = luminametacontract.NewLuminametacontractCaller(metacontractAddress, client)
	if err != nil {
		return
	}

	priceBig, timeUnixBig, err := caller.GetValue(&bind.CallOpts{}, a.Symbol+"/USD")
	if err != nil {
		return
	}

	if priceBig == nil {
		err = errors.New("price is zero")
		return
	}

	aq.Price, _ = new(big.Float).Mul(new(big.Float).SetInt(priceBig), big.NewFloat(math.Pow10(-precision))).Float64()

	if timeUnixBig != nil {
		aq.Time = time.Unix(timeUnixBig.Int64(), 0)
	} else {
		aq.Time = time.Now()
	}

	aq.Asset = *a
	aq.Source = "DIA_Lumina_" + metacontractAddress.Hex()

	return
}

// GetPriceFromDiaAPI returns the price of the asset given by @address and @blockchain as returned by
// assetQuotation endpoint of diadata API.
func (a *Asset) GetPriceFromDiaAPI() (aq AssetQuotation, err error) {

	baseString := "https://api.diadata.org/v1/assetQuotation/" + a.Blockchain + "/" + a.Address

	var response []byte
	response, _, err = utils.GetRequest(baseString)
	if err != nil {
		return
	}

	err = json.Unmarshal(response, &aq)

	return
}

func NewTokenCaller(address common.Address, caller bind.ContractCaller) (*TokenCaller, error) {
	contract, err := bindToken(address, caller, nil)
	if err != nil {
		return nil, err
	}
	return &TokenCaller{Contract: contract}, nil
}

// bindToken binds a generic wrapper to an already deployed contract.
func bindToken(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(tokenABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, nil), nil
}

func GetAsset(address common.Address, blockchain string, client *ethclient.Client) (asset Asset, err error) {
	tokenCaller, err := NewTokenCaller(address, client)
	if err != nil {
		log.Error("NewIERC20Caller: ", err)
		return
	}

	if address == common.HexToAddress("0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2") && blockchain == utils.ETHEREUM {
		asset.Symbol = "MKR"
	} else {
		var symbol []interface{}
		err = tokenCaller.Contract.Call(&bind.CallOpts{}, &symbol, "symbol")
		if err != nil {
			return Asset{}, err
		}
		asset.Symbol = symbol[0].(string)
	}

	if address == common.HexToAddress("0x9f8F72aA9304c8B593d555F12eF6589cC3A579A2") && blockchain == utils.ETHEREUM {
		asset.Name = "Maker"
	} else {
		var name []interface{}
		err = tokenCaller.Contract.Call(&bind.CallOpts{}, &name, "name")
		if err != nil {
			return Asset{}, err
		}
		asset.Name = name[0].(string)
	}

	var decimals []interface{}
	err = tokenCaller.Contract.Call(&bind.CallOpts{}, &decimals, "decimals")
	if err != nil {
		return Asset{}, err
	}
	aux := decimals[0].(*big.Int)
	asset.Decimals = uint8(aux.Int64())
	asset.Address = address.Hex()
	asset.Blockchain = blockchain

	return
}

// getBalance returns the balance of @a in pool with address @poolAddress.
func (a *Asset) GetBalance(poolAddress common.Address, client *ethclient.Client) (float64, error) {
	balanceBig, err := utils.GetBalanceOf(common.HexToAddress(a.Address), poolAddress, client)
	if err != nil {
		return 0, err
	}
	balance, _ := new(big.Float).Quo(big.NewFloat(0).SetInt(balanceBig), new(big.Float).SetFloat64(math.Pow10(int(a.Decimals)))).Float64()
	return balance, nil
}
