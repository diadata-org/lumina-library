// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package diaoraclev3

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

// IDIAOracleV3ValueEntry is an auto generated low-level Go binding around an user-defined struct.
type IDIAOracleV3ValueEntry struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}

// DIAOracleV3MetaData contains all meta data concerning the DIAOracleV3 contract.
var DIAOracleV3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_HISTORY_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_TIMESTAMP_GAP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPDATER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxHistorySize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRawData\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValue\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValueAt\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"volume\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValueCount\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValueHistory\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIDIAOracleV3.ValueEntry[]\",\"components\":[{\"name\":\"value\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"volume\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"decimalPrecision\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rawData\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMultipleRawValues\",\"inputs\":[{\"name\":\"dataArray\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMultipleValues\",\"inputs\":[{\"name\":\"keys\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"compressedValues\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRawValue\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setValue\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"values\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleUpdate\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleUpdateRaw\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"volume\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdaterAddressChange\",\"inputs\":[{\"name\":\"newUpdater\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHistoryIndex\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[{\"name\":\"keysLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"valuesLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimestampNotIncreasing\",\"inputs\":[{\"name\":\"newTimestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"existingTimestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"TimestampTooFarInFuture\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"blockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TimestampTooFarInPast\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"blockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60c080604052346100f35730608052606460a0525f51602061213c5f395f51905f525460ff8160401c166100e4576002600160401b03196001600160401b03821601610091575b60405161204490816100f8823960805181818161076f015261080d015260a05181818161114f015281816112c50152818161170301528181611a4a01528181611ac00152611bcd0152f35b6001600160401b0319166001600160401b039081175f51602061213c5f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610046565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a7146111b25750806309daaa951461103d578063135d90c714610dff578063248a9ca314610dc057806324b1db1a146101c55780632c484ae514610ba55780632f2ff15d14610b5b578063313ce5671461015057806336568abe14610b175780634351e6b6146109fc57806347e63380146109c25780634df710961461059c5780634f1ef286146107c357806352d1902d1461075d57806359c3852c146107085780635a9ade8b146106b45780637898e0c2146106065780637a2fa4421461059c5780638d241526146103a55780638d97ecf21461030557806391d14854146102b0578063960384a01461024b578063a217fddf14610231578063ad3cb1cc146101e6578063be2e0179146101ca578063c1066251146101c5578063d547741f146101745763f0141d8414610150575f80fd5b34610170575f36600319011261017057602060ff60055416604051908152f35b5f80fd5b34610170576040366003190112610170576101c36004356101936112e8565b906101be6101b9825f525f516020611faf5f395f51905f52602052600160405f20015490565b611838565b611e66565b005b6112ae565b34610170575f366003190112610170576020604051610e108152f35b34610170575f3660031901126101705761022d60405161020760408261123a565b60058152640352e302e360dc1b60208201526040519182916020835260208301906112fe565b0390f35b34610170575f3660031901126101705760206040515f8152f35b34610170576020366003190112610170576004356001600160401b03811161017057602080610280604093369060040161125b565b8351928184925191829101835e81015f815203019020546001600160801b038251918060801c8352166020820152f35b34610170576040366003190112610170576102c96112e8565b6004355f525f516020611faf5f395f51905f5260205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34610170576020366003190112610170576004356001600160401b0381116101705761033861033d91369060040161125b565b61161c565b6040518091602082016020835281518091526020604084019201905f5b818110610368575050500390f35b91935091602060606001926001600160801b036040885182815116845282868201511686850152015116604082015201940191019184939261035a565b34610170576040366003190112610170576004356001600160401b03811161017057366023820112156101705780600401356103e08161140e565b916103ee604051938461123a565b8183526024602084019260051b820101903682116101705760248101925b82841061056d57602435856001600160401b03821161017057366023830112156101705781600401359161043f8361140e565b9261044d604051948561123a565b8084526024602085019160051b8301019136831161017057602401905b82821061055d5750505061047c6117c9565b80518251908181036105485750505f5b81518110156101c357807fa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee22917826104c360019385611608565b516104ce8387611608565b519061053f8260801c926105286001600160801b038216916104f0838661187e565b60405190855191602081818901948086835e81015f815203019020556020604051809287518091835e8101600481520301902061158b565b6105338185856119ab565b604051938493846115da565b0390a10161048c565b633adf150f60e21b5f5260045260245260445ffd5b813581526020918201910161046a565b83356001600160401b0381116101705760209161059183926024369187010161125b565b81520193019261040c565b34610170576020366003190112610170576004356001600160401b038111610170576105f26020806105d561022d94369060040161125b565b604051928184925191829101835e8101600481520301902061136e565b6040519182916020835260208301906112fe565b34610170576060366003190112610170576004356001600160401b0381116101705761063690369060040161125b565b6024356001600160801b038116810361017057604435916001600160801b038316808403610170577fa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782936105286104f06106af936106926117c9565b61069c848761187e565b6001600160801b03198760801b16611425565b0390a1005b34610170576020366003190112610170576004356001600160401b038111610170576020806106e88193369060040161125b565b604051928184925191829101835e81015f81520301902054604051908152f35b34610170576020366003190112610170576004356001600160401b0381116101705760208061073c8193369060040161125b565b604051928184925191829101835e8101600381520301902054604051908152f35b34610170575f366003190112610170577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031630036107b45760206040515f516020611f8f5f395f51905f528152f35b63703e46dd60e11b5f5260045ffd5b6040366003190112610170576004356001600160a01b03811690818103610170576024356001600160401b0381116101705761080390369060040161125b565b6001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000163081149081156109a0575b506107b457335f9081527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d602052604090205460ff1615610989576040516352d1902d60e01b8152602081600481875afa5f9181610955575b506108a95783634c9c8ce360e01b5f5260045260245ffd5b805f516020611f8f5f395f51905f528592036109435750823b15610931575f516020611f8f5f395f51905f5280546001600160a01b031916821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2805115610919576101c391611f02565b50503461092257005b63b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b632a87526960e21b5f5260045260245ffd5b9091506020813d602011610981575b816109716020938361123a565b8101031261017057519085610891565b3d9150610964565b63e2517d3f60e01b5f52336004525f60245260445ffd5b5f516020611f8f5f395f51905f52546001600160a01b03161415905084610838565b34610170575f3660031901126101705760206040517f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab8152f35b346101705760203660031901126101705760043560ff8116809103610170575f516020611fcf5f395f51905f525460ff8160401c168015610b03575b610af45768ffffffffffffffffff191668010000000000000001175f516020611fcf5f395f51905f5281905560401c60ff1615610ae55760ff196005541617600555610a8333611c47565b50610a8d33611cf6565b5068ff0000000000000000195f516020611fcf5f395f51905f5254165f516020611fcf5f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160018152a1005b631afcd79f60e31b5f5260045ffd5b63f92ee8a960e01b5f5260045ffd5b5060016001600160401b0382161015610a38565b3461017057604036600319011261017057610b306112e8565b336001600160a01b03821603610b4c576101c390600435611e66565b63334bd91960e11b5f5260045ffd5b34610170576040366003190112610170576101c3600435610b7a6112e8565b90610ba06101b9825f525f516020611faf5f395f51905f52602052600160405f20015490565b611dc2565b34610170576020366003190112610170576004356001600160401b03811161017057366023820112156101705780600401356001600160401b038111610170578101602401368111610170576024610c0592610bff6117c9565b016114bc565b90610c128386949661187e565b610c326001600160801b0386166001600160801b03198660801b16611425565b60405190845191602081818801948086835e81015f815203019020556020604051809286518091835e81016004815203019020948251956001600160401b038711610dac57610c818154611336565b601f8111610d69575b50602096601f8111600114610cf45790610cce826106af969594935f516020611fef5f395f51905f529a5f91610ce9575b508160011b915f199060031b1c19161790565b90555b610cdd82828888611b39565b60405195869586611548565b90508601518b610cbb565b601f19811697825f52805f20985f5b818110610d515750915f516020611fef5f395f51905f5299600192826106af999897969510610d39575b5050811b019055610cd1565b8701515f1960f88460031b161c191690558a80610d2d565b878301518b556001909a019960209283019201610d03565b87811115610c8a57610d9e90825f5260205f2090601f8a0160051c9060208b10610da4575b601f82910160051c03910161152d565b87610c8a565b5f9150610d8e565b634e487b7160e01b5f52604160045260245ffd5b34610170576020366003190112610170576020610df76004355f525f516020611faf5f395f51905f52602052600160405f20015490565b604051908152f35b34610170576020366003190112610170576004356001600160401b03811161017057366023820112156101705780600401356001600160401b038111610170573660248260051b840101116101705790610e576117c9565b3681900360421901905f5b838110156101c35760248160051b830101358381121561017057820160248101356001600160401b038111610170576044820190803603821361017057610ead9201604401906114bc565b91610eba8186959661187e565b610eda6001600160801b0382166001600160801b03198760801b16611425565b60405190855191602081818901948086835e81015f815203019020556020604051809287518091835e8101600481520301902083516001600160401b038111610dac57610f278254611336565b601f8111611003575b506020601f8211600114610f8b5792600198979592610cce835f516020611fef5f395f51905f52999794610f77975f91610f8057508160011b915f199060031b1c19161790565b0390a101610e62565b90508601515f610cbb565b601f19821690835f52805f20915f5b818110610feb575083610f77969360019c9b9996935f516020611fef5f395f51905f529b99968e9410610fd3575050811b019055610cd1565b8701515f1960f88460031b161c191690558f80610d2d565b9192602060018192868c015181550194019201610f9a565b81811115610f305761103790835f5260205f2090601f840160051c9060208510610da457601f82910160051c03910161152d565b8b610f30565b34610170576040366003190112610170576004356001600160401b0381116101705761106d90369060040161125b565b602435604051825190602081818601938085835e810160018152030190209260405160208183518086835e81016003815203019020548084101561119c5750602090604051928391518091835e81016002815203019020549060018101808211611135578210611149575f198201918211611135576060926110f56110fb9261110194611432565b9061145d565b5061148a565b6001600160801b03815116906001600160801b03604081602084015116920151169060405192835260208301526040820152f35b634e487b7160e01b5f52601160045260245ffd5b916111757f00000000000000000000000000000000000000000000000000000000000000008093611425565b5f19810191908211611135576110f5611101936111976060966110fb95611432565b61143f565b8363c9d622b960e01b5f5260045260245260445ffd5b34610170576020366003190112610170576004359063ffffffff60e01b821680920361017057602091631876be2560e01b81149081156111f4575b5015158152f35b637965db0b60e01b81149150811561120e575b50836111ed565b6301ffc9a760e01b14905083611207565b606081019081106001600160401b03821117610dac57604052565b90601f801991011681019081106001600160401b03821117610dac57604052565b81601f82011215610170576020813591016001600160401b038211610dac5760405192611292601f8401601f19166020018561123a565b8284528282011161017057815f92602092838601378301015290565b34610170575f3660031901126101705760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b602435906001600160a01b038216820361017057565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b35906001600160801b038216820361017057565b90600182811c92168015611364575b602083101461135057565b634e487b7160e01b5f52602260045260245ffd5b91607f1691611345565b9060405191825f82549261138184611336565b80845293600181169081156113ec57506001146113a8575b506113a69250038361123a565b565b90505f9291925260205f20905f915b8183106113d05750509060206113a6928201015f611399565b60209193508060019154838589010152019101909184926113b7565b9050602092506113a694915060ff191682840152151560051b8201015f611399565b6001600160401b038111610dac5760051b60200190565b9190820180921161113557565b9190820391821161113557565b8115611449570690565b634e487b7160e01b5f52601260045260245ffd5b8054821015611476575f5260205f209060011b01905f90565b634e487b7160e01b5f52603260045260245ffd5b906040516114978161121f565b60406001600160801b03600183958054838116865260801c6020860152015416910152565b91909160a0818403126101705780356001600160401b03811161017057836114e591830161125b565b926114f260208301611322565b926114ff60408401611322565b9261150c60608201611322565b9260808201356001600160401b0381116101705761152a920161125b565b90565b5f5b82811061153b57505050565b5f8282015560010161152f565b91936001600160801b0361152a96948161156b819560a0885260a08801906112fe565b9716602086015216604084015216606082015260808184039101526112fe565b6115958154611336565b908161159f575050565b81601f5f93116001146115b0575055565b818352602083206115cd91601f0160051c8419019060010161152d565b8082528160208120915555565b9193926001600160801b0390816115fb6040946060875260608701906112fe565b9616602085015216910152565b80518210156114765760209160051b010190565b604051815190602081818501938085835e8101600181520301902060405160208185518086835e810160038152030190205491821561177c5761165e8361140e565b9361166c604051958661123a565b838552601f1961167b8561140e565b015f5b818110611753575050602090604051928391518091835e81016002815203019020545f5b8381106116b0575050505090565b60018101808211611135578210611701575f19820190828211611135576116e56110fb6116df83600195611432565b8661145d565b6116ef8288611608565b526116fa8187611608565b50016116a2565b7f00000000000000000000000000000000000000000000000000000000000000009061172d8284611425565b5f198101908111611135576110fb61174e600194611197856116e595611432565b6116df565b6020906040516117628161121f565b5f81525f838201525f604082015282828a0101520161167e565b5050505060405161178e60208261123a565b5f81525f805b8181106117a057505090565b6020906040516117af8161121f565b5f81525f838201525f604082015282828601015201611794565b335f9081527f268477d1c17bf7595397186358fec802c44c92bb39292de205ad03cbfe096e02602052604090205460ff161561180157565b63e2517d3f60e01b5f52336004527f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab60245260445ffd5b5f8181525f516020611faf5f395f51905f526020908152604080832033845290915290205460ff16156118685750565b63e2517d3f60e01b5f523360045260245260445ffd5b610e10420191824211611135576001600160801b038091169216821161193257610e10421180611914575b6118fd5760208091604051928184925191829101835e81015f81520301902054806118d2575050565b6001600160801b031690818111156118e8575050565b633e4ebc1960e21b5f5260045260245260445ffd5b5063514a513560e01b5f526004524260245260445ffd5b50610e0f194201428111611135576001600160801b031682106118a9565b506304c3bb5b60e21b5f526004524260245260445ffd5b9190611998578051602082015160801b6fffffffffffffffffffffffffffffffff199081166001600160801b039283161784556040929092015160019390930180549092169216919091179055565b634e487b7160e01b5f525f60045260245ffd5b60405191815192602081818501958087835e810160018152030190209060405160208185518088835e81016002815203019020549160405160208186518089835e810160038152030190205495815415611aba575b83611a3a93926001600160801b03611a34938160405196611a208861121f565b1686521660208501525f604085015261145d565b90611949565b6001810180911161113557611a707f0000000000000000000000000000000000000000000000000000000000000000809261143f565b60405160208185518088835e81016002815203019020558310611a9257505050565b5f198314611135576020906001604051938492518091845e8201946003865201930301902055565b955091507f00000000000000000000000000000000000000000000000000000000000000005f5b818110611af457505f9586939150611a00565b60405190611b018261121f565b5f82525f60208301525f60408301528454600160401b811015610dac57600192611a348285611b33940189558861145d565b01611ae1565b90919260405192825193602081818601968088835e8101600181520301902060405160208186518089835e8101600281520301902054926040516020818751808a835e810160038152030190205496825415611bc6575b611a3a93926001600160801b03611a34938188948160405198611bb28a61121f565b16885216602087015216604085015261145d565b96509192507f00000000000000000000000000000000000000000000000000000000000000005f5b818110611c0257505f968794939150611b90565b60405190611c0f8261121f565b5f82525f60208301525f60408301528354600160401b811015610dac57600192611a348285611c41940188558761145d565b01611bee565b6001600160a01b0381165f9081527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d602052604090205460ff16611cf1576001600160a01b03165f8181527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d60205260408120805460ff191660011790553391907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d8180a4600190565b505f90565b6001600160a01b0381165f9081527f268477d1c17bf7595397186358fec802c44c92bb39292de205ad03cbfe096e02602052604090205460ff16611cf1576001600160a01b03165f8181527f268477d1c17bf7595397186358fec802c44c92bb39292de205ad03cbfe096e0260205260408120805460ff191660011790553391907f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b5f8181525f516020611faf5f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff16611e60575f8181525f516020611faf5f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b5f8181525f516020611faf5f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff1615611e60575f8181525f516020611faf5f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b905f8091602081519101845af48080611f7b575b15611f365750506040513d81523d5f602083013e60203d82010160405290565b15611f5b57639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b3d15611f6c576040513d5f823e3d90fd5b63d6bda27560e01b5f5260045ffd5b503d151580611f165750813b1515611f1656fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a000ec1e0298284e066eddd5e448f165c9337bf3f9447b7159177c72e0cada227d3a264697066735822122073394c4177146c8ae952e915ea1d03fb7b53b53565e7d0823cb296b6fcf69d0d64736f6c63430008220033f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
}

// DIAOracleV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use DIAOracleV3MetaData.ABI instead.
var DIAOracleV3ABI = DIAOracleV3MetaData.ABI

// DIAOracleV3Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DIAOracleV3MetaData.Bin instead.
var DIAOracleV3Bin = DIAOracleV3MetaData.Bin

// DeployDIAOracleV3 deploys a new Ethereum contract, binding an instance of DIAOracleV3 to it.
func DeployDIAOracleV3(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DIAOracleV3, error) {
	parsed, err := DIAOracleV3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DIAOracleV3Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DIAOracleV3{DIAOracleV3Caller: DIAOracleV3Caller{contract: contract}, DIAOracleV3Transactor: DIAOracleV3Transactor{contract: contract}, DIAOracleV3Filterer: DIAOracleV3Filterer{contract: contract}}, nil
}

// DIAOracleV3 is an auto generated Go binding around an Ethereum contract.
type DIAOracleV3 struct {
	DIAOracleV3Caller     // Read-only binding to the contract
	DIAOracleV3Transactor // Write-only binding to the contract
	DIAOracleV3Filterer   // Log filterer for contract events
}

// DIAOracleV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type DIAOracleV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DIAOracleV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type DIAOracleV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DIAOracleV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DIAOracleV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DIAOracleV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DIAOracleV3Session struct {
	Contract     *DIAOracleV3      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DIAOracleV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DIAOracleV3CallerSession struct {
	Contract *DIAOracleV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DIAOracleV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DIAOracleV3TransactorSession struct {
	Contract     *DIAOracleV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DIAOracleV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type DIAOracleV3Raw struct {
	Contract *DIAOracleV3 // Generic contract binding to access the raw methods on
}

// DIAOracleV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DIAOracleV3CallerRaw struct {
	Contract *DIAOracleV3Caller // Generic read-only contract binding to access the raw methods on
}

// DIAOracleV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DIAOracleV3TransactorRaw struct {
	Contract *DIAOracleV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewDIAOracleV3 creates a new instance of DIAOracleV3, bound to a specific deployed contract.
func NewDIAOracleV3(address common.Address, backend bind.ContractBackend) (*DIAOracleV3, error) {
	contract, err := bindDIAOracleV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3{DIAOracleV3Caller: DIAOracleV3Caller{contract: contract}, DIAOracleV3Transactor: DIAOracleV3Transactor{contract: contract}, DIAOracleV3Filterer: DIAOracleV3Filterer{contract: contract}}, nil
}

// NewDIAOracleV3Caller creates a new read-only instance of DIAOracleV3, bound to a specific deployed contract.
func NewDIAOracleV3Caller(address common.Address, caller bind.ContractCaller) (*DIAOracleV3Caller, error) {
	contract, err := bindDIAOracleV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3Caller{contract: contract}, nil
}

// NewDIAOracleV3Transactor creates a new write-only instance of DIAOracleV3, bound to a specific deployed contract.
func NewDIAOracleV3Transactor(address common.Address, transactor bind.ContractTransactor) (*DIAOracleV3Transactor, error) {
	contract, err := bindDIAOracleV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3Transactor{contract: contract}, nil
}

// NewDIAOracleV3Filterer creates a new log filterer instance of DIAOracleV3, bound to a specific deployed contract.
func NewDIAOracleV3Filterer(address common.Address, filterer bind.ContractFilterer) (*DIAOracleV3Filterer, error) {
	contract, err := bindDIAOracleV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3Filterer{contract: contract}, nil
}

// bindDIAOracleV3 binds a generic wrapper to an already deployed contract.
func bindDIAOracleV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DIAOracleV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DIAOracleV3 *DIAOracleV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DIAOracleV3.Contract.DIAOracleV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DIAOracleV3 *DIAOracleV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.DIAOracleV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DIAOracleV3 *DIAOracleV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.DIAOracleV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DIAOracleV3 *DIAOracleV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DIAOracleV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DIAOracleV3 *DIAOracleV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DIAOracleV3 *DIAOracleV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _DIAOracleV3.Contract.DEFAULTADMINROLE(&_DIAOracleV3.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DIAOracleV3.Contract.DEFAULTADMINROLE(&_DIAOracleV3.CallOpts)
}

// MAXHISTORYSIZE is a free data retrieval call binding the contract method 0xc1066251.
//
// Solidity: function MAX_HISTORY_SIZE() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Caller) MAXHISTORYSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "MAX_HISTORY_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXHISTORYSIZE is a free data retrieval call binding the contract method 0xc1066251.
//
// Solidity: function MAX_HISTORY_SIZE() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Session) MAXHISTORYSIZE() (*big.Int, error) {
	return _DIAOracleV3.Contract.MAXHISTORYSIZE(&_DIAOracleV3.CallOpts)
}

// MAXHISTORYSIZE is a free data retrieval call binding the contract method 0xc1066251.
//
// Solidity: function MAX_HISTORY_SIZE() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3CallerSession) MAXHISTORYSIZE() (*big.Int, error) {
	return _DIAOracleV3.Contract.MAXHISTORYSIZE(&_DIAOracleV3.CallOpts)
}

// MAXTIMESTAMPGAP is a free data retrieval call binding the contract method 0xbe2e0179.
//
// Solidity: function MAX_TIMESTAMP_GAP() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Caller) MAXTIMESTAMPGAP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "MAX_TIMESTAMP_GAP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTIMESTAMPGAP is a free data retrieval call binding the contract method 0xbe2e0179.
//
// Solidity: function MAX_TIMESTAMP_GAP() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Session) MAXTIMESTAMPGAP() (*big.Int, error) {
	return _DIAOracleV3.Contract.MAXTIMESTAMPGAP(&_DIAOracleV3.CallOpts)
}

// MAXTIMESTAMPGAP is a free data retrieval call binding the contract method 0xbe2e0179.
//
// Solidity: function MAX_TIMESTAMP_GAP() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3CallerSession) MAXTIMESTAMPGAP() (*big.Int, error) {
	return _DIAOracleV3.Contract.MAXTIMESTAMPGAP(&_DIAOracleV3.CallOpts)
}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Caller) UPDATERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "UPDATER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Session) UPDATERROLE() ([32]byte, error) {
	return _DIAOracleV3.Contract.UPDATERROLE(&_DIAOracleV3.CallOpts)
}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3CallerSession) UPDATERROLE() ([32]byte, error) {
	return _DIAOracleV3.Contract.UPDATERROLE(&_DIAOracleV3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DIAOracleV3 *DIAOracleV3Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DIAOracleV3 *DIAOracleV3Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _DIAOracleV3.Contract.UPGRADEINTERFACEVERSION(&_DIAOracleV3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DIAOracleV3 *DIAOracleV3CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _DIAOracleV3.Contract.UPGRADEINTERFACEVERSION(&_DIAOracleV3.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DIAOracleV3 *DIAOracleV3Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DIAOracleV3 *DIAOracleV3Session) Decimals() (uint8, error) {
	return _DIAOracleV3.Contract.Decimals(&_DIAOracleV3.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DIAOracleV3 *DIAOracleV3CallerSession) Decimals() (uint8, error) {
	return _DIAOracleV3.Contract.Decimals(&_DIAOracleV3.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_DIAOracleV3 *DIAOracleV3Caller) GetDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_DIAOracleV3 *DIAOracleV3Session) GetDecimals() (uint8, error) {
	return _DIAOracleV3.Contract.GetDecimals(&_DIAOracleV3.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetDecimals() (uint8, error) {
	return _DIAOracleV3.Contract.GetDecimals(&_DIAOracleV3.CallOpts)
}

// GetMaxHistorySize is a free data retrieval call binding the contract method 0x24b1db1a.
//
// Solidity: function getMaxHistorySize() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Caller) GetMaxHistorySize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getMaxHistorySize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxHistorySize is a free data retrieval call binding the contract method 0x24b1db1a.
//
// Solidity: function getMaxHistorySize() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Session) GetMaxHistorySize() (*big.Int, error) {
	return _DIAOracleV3.Contract.GetMaxHistorySize(&_DIAOracleV3.CallOpts)
}

// GetMaxHistorySize is a free data retrieval call binding the contract method 0x24b1db1a.
//
// Solidity: function getMaxHistorySize() view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetMaxHistorySize() (*big.Int, error) {
	return _DIAOracleV3.Contract.GetMaxHistorySize(&_DIAOracleV3.CallOpts)
}

// GetRawData is a free data retrieval call binding the contract method 0x4df71096.
//
// Solidity: function getRawData(string key) view returns(bytes)
func (_DIAOracleV3 *DIAOracleV3Caller) GetRawData(opts *bind.CallOpts, key string) ([]byte, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getRawData", key)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRawData is a free data retrieval call binding the contract method 0x4df71096.
//
// Solidity: function getRawData(string key) view returns(bytes)
func (_DIAOracleV3 *DIAOracleV3Session) GetRawData(key string) ([]byte, error) {
	return _DIAOracleV3.Contract.GetRawData(&_DIAOracleV3.CallOpts, key)
}

// GetRawData is a free data retrieval call binding the contract method 0x4df71096.
//
// Solidity: function getRawData(string key) view returns(bytes)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetRawData(key string) ([]byte, error) {
	return _DIAOracleV3.Contract.GetRawData(&_DIAOracleV3.CallOpts, key)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DIAOracleV3.Contract.GetRoleAdmin(&_DIAOracleV3.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DIAOracleV3.Contract.GetRoleAdmin(&_DIAOracleV3.CallOpts, role)
}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_DIAOracleV3 *DIAOracleV3Caller) GetValue(opts *bind.CallOpts, key string) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getValue", key)

	if err != nil {
		return *new(*big.Int), *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	out1 := *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)

	return out0, out1, err

}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_DIAOracleV3 *DIAOracleV3Session) GetValue(key string) (*big.Int, *big.Int, error) {
	return _DIAOracleV3.Contract.GetValue(&_DIAOracleV3.CallOpts, key)
}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetValue(key string) (*big.Int, *big.Int, error) {
	return _DIAOracleV3.Contract.GetValue(&_DIAOracleV3.CallOpts, key)
}

// GetValueAt is a free data retrieval call binding the contract method 0x09daaa95.
//
// Solidity: function getValueAt(string key, uint256 index) view returns(uint128 value, uint128 timestamp, uint128 volume)
func (_DIAOracleV3 *DIAOracleV3Caller) GetValueAt(opts *bind.CallOpts, key string, index *big.Int) (struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getValueAt", key, index)

	outstruct := new(struct {
		Value     *big.Int
		Timestamp *big.Int
		Volume    *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Value = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.Timestamp = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Volume = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GetValueAt is a free data retrieval call binding the contract method 0x09daaa95.
//
// Solidity: function getValueAt(string key, uint256 index) view returns(uint128 value, uint128 timestamp, uint128 volume)
func (_DIAOracleV3 *DIAOracleV3Session) GetValueAt(key string, index *big.Int) (struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}, error) {
	return _DIAOracleV3.Contract.GetValueAt(&_DIAOracleV3.CallOpts, key, index)
}

// GetValueAt is a free data retrieval call binding the contract method 0x09daaa95.
//
// Solidity: function getValueAt(string key, uint256 index) view returns(uint128 value, uint128 timestamp, uint128 volume)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetValueAt(key string, index *big.Int) (struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}, error) {
	return _DIAOracleV3.Contract.GetValueAt(&_DIAOracleV3.CallOpts, key, index)
}

// GetValueCount is a free data retrieval call binding the contract method 0x59c3852c.
//
// Solidity: function getValueCount(string key) view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Caller) GetValueCount(opts *bind.CallOpts, key string) (*big.Int, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getValueCount", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValueCount is a free data retrieval call binding the contract method 0x59c3852c.
//
// Solidity: function getValueCount(string key) view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Session) GetValueCount(key string) (*big.Int, error) {
	return _DIAOracleV3.Contract.GetValueCount(&_DIAOracleV3.CallOpts, key)
}

// GetValueCount is a free data retrieval call binding the contract method 0x59c3852c.
//
// Solidity: function getValueCount(string key) view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetValueCount(key string) (*big.Int, error) {
	return _DIAOracleV3.Contract.GetValueCount(&_DIAOracleV3.CallOpts, key)
}

// GetValueHistory is a free data retrieval call binding the contract method 0x8d97ecf2.
//
// Solidity: function getValueHistory(string key) view returns((uint128,uint128,uint128)[])
func (_DIAOracleV3 *DIAOracleV3Caller) GetValueHistory(opts *bind.CallOpts, key string) ([]IDIAOracleV3ValueEntry, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "getValueHistory", key)

	if err != nil {
		return *new([]IDIAOracleV3ValueEntry), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDIAOracleV3ValueEntry)).(*[]IDIAOracleV3ValueEntry)

	return out0, err

}

// GetValueHistory is a free data retrieval call binding the contract method 0x8d97ecf2.
//
// Solidity: function getValueHistory(string key) view returns((uint128,uint128,uint128)[])
func (_DIAOracleV3 *DIAOracleV3Session) GetValueHistory(key string) ([]IDIAOracleV3ValueEntry, error) {
	return _DIAOracleV3.Contract.GetValueHistory(&_DIAOracleV3.CallOpts, key)
}

// GetValueHistory is a free data retrieval call binding the contract method 0x8d97ecf2.
//
// Solidity: function getValueHistory(string key) view returns((uint128,uint128,uint128)[])
func (_DIAOracleV3 *DIAOracleV3CallerSession) GetValueHistory(key string) ([]IDIAOracleV3ValueEntry, error) {
	return _DIAOracleV3.Contract.GetValueHistory(&_DIAOracleV3.CallOpts, key)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DIAOracleV3 *DIAOracleV3Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DIAOracleV3 *DIAOracleV3Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DIAOracleV3.Contract.HasRole(&_DIAOracleV3.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DIAOracleV3 *DIAOracleV3CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DIAOracleV3.Contract.HasRole(&_DIAOracleV3.CallOpts, role, account)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3Session) ProxiableUUID() ([32]byte, error) {
	return _DIAOracleV3.Contract.ProxiableUUID(&_DIAOracleV3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DIAOracleV3 *DIAOracleV3CallerSession) ProxiableUUID() ([32]byte, error) {
	return _DIAOracleV3.Contract.ProxiableUUID(&_DIAOracleV3.CallOpts)
}

// RawData is a free data retrieval call binding the contract method 0x7a2fa442.
//
// Solidity: function rawData(string ) view returns(bytes)
func (_DIAOracleV3 *DIAOracleV3Caller) RawData(opts *bind.CallOpts, arg0 string) ([]byte, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "rawData", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RawData is a free data retrieval call binding the contract method 0x7a2fa442.
//
// Solidity: function rawData(string ) view returns(bytes)
func (_DIAOracleV3 *DIAOracleV3Session) RawData(arg0 string) ([]byte, error) {
	return _DIAOracleV3.Contract.RawData(&_DIAOracleV3.CallOpts, arg0)
}

// RawData is a free data retrieval call binding the contract method 0x7a2fa442.
//
// Solidity: function rawData(string ) view returns(bytes)
func (_DIAOracleV3 *DIAOracleV3CallerSession) RawData(arg0 string) ([]byte, error) {
	return _DIAOracleV3.Contract.RawData(&_DIAOracleV3.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DIAOracleV3 *DIAOracleV3Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DIAOracleV3 *DIAOracleV3Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DIAOracleV3.Contract.SupportsInterface(&_DIAOracleV3.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DIAOracleV3 *DIAOracleV3CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DIAOracleV3.Contract.SupportsInterface(&_DIAOracleV3.CallOpts, interfaceId)
}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Caller) Values(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _DIAOracleV3.contract.Call(opts, &out, "values", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3Session) Values(arg0 string) (*big.Int, error) {
	return _DIAOracleV3.Contract.Values(&_DIAOracleV3.CallOpts, arg0)
}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_DIAOracleV3 *DIAOracleV3CallerSession) Values(arg0 string) (*big.Int, error) {
	return _DIAOracleV3.Contract.Values(&_DIAOracleV3.CallOpts, arg0)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DIAOracleV3 *DIAOracleV3Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.GrantRole(&_DIAOracleV3.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.GrantRole(&_DIAOracleV3.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x4351e6b6.
//
// Solidity: function initialize(uint8 decimalPrecision) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) Initialize(opts *bind.TransactOpts, decimalPrecision uint8) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "initialize", decimalPrecision)
}

// Initialize is a paid mutator transaction binding the contract method 0x4351e6b6.
//
// Solidity: function initialize(uint8 decimalPrecision) returns()
func (_DIAOracleV3 *DIAOracleV3Session) Initialize(decimalPrecision uint8) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.Initialize(&_DIAOracleV3.TransactOpts, decimalPrecision)
}

// Initialize is a paid mutator transaction binding the contract method 0x4351e6b6.
//
// Solidity: function initialize(uint8 decimalPrecision) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) Initialize(decimalPrecision uint8) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.Initialize(&_DIAOracleV3.TransactOpts, decimalPrecision)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DIAOracleV3 *DIAOracleV3Session) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.RenounceRole(&_DIAOracleV3.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.RenounceRole(&_DIAOracleV3.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DIAOracleV3 *DIAOracleV3Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.RevokeRole(&_DIAOracleV3.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.RevokeRole(&_DIAOracleV3.TransactOpts, role, account)
}

// SetMultipleRawValues is a paid mutator transaction binding the contract method 0x135d90c7.
//
// Solidity: function setMultipleRawValues(bytes[] dataArray) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) SetMultipleRawValues(opts *bind.TransactOpts, dataArray [][]byte) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "setMultipleRawValues", dataArray)
}

// SetMultipleRawValues is a paid mutator transaction binding the contract method 0x135d90c7.
//
// Solidity: function setMultipleRawValues(bytes[] dataArray) returns()
func (_DIAOracleV3 *DIAOracleV3Session) SetMultipleRawValues(dataArray [][]byte) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetMultipleRawValues(&_DIAOracleV3.TransactOpts, dataArray)
}

// SetMultipleRawValues is a paid mutator transaction binding the contract method 0x135d90c7.
//
// Solidity: function setMultipleRawValues(bytes[] dataArray) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) SetMultipleRawValues(dataArray [][]byte) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetMultipleRawValues(&_DIAOracleV3.TransactOpts, dataArray)
}

// SetMultipleValues is a paid mutator transaction binding the contract method 0x8d241526.
//
// Solidity: function setMultipleValues(string[] keys, uint256[] compressedValues) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) SetMultipleValues(opts *bind.TransactOpts, keys []string, compressedValues []*big.Int) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "setMultipleValues", keys, compressedValues)
}

// SetMultipleValues is a paid mutator transaction binding the contract method 0x8d241526.
//
// Solidity: function setMultipleValues(string[] keys, uint256[] compressedValues) returns()
func (_DIAOracleV3 *DIAOracleV3Session) SetMultipleValues(keys []string, compressedValues []*big.Int) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetMultipleValues(&_DIAOracleV3.TransactOpts, keys, compressedValues)
}

// SetMultipleValues is a paid mutator transaction binding the contract method 0x8d241526.
//
// Solidity: function setMultipleValues(string[] keys, uint256[] compressedValues) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) SetMultipleValues(keys []string, compressedValues []*big.Int) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetMultipleValues(&_DIAOracleV3.TransactOpts, keys, compressedValues)
}

// SetRawValue is a paid mutator transaction binding the contract method 0x2c484ae5.
//
// Solidity: function setRawValue(bytes data) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) SetRawValue(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "setRawValue", data)
}

// SetRawValue is a paid mutator transaction binding the contract method 0x2c484ae5.
//
// Solidity: function setRawValue(bytes data) returns()
func (_DIAOracleV3 *DIAOracleV3Session) SetRawValue(data []byte) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetRawValue(&_DIAOracleV3.TransactOpts, data)
}

// SetRawValue is a paid mutator transaction binding the contract method 0x2c484ae5.
//
// Solidity: function setRawValue(bytes data) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) SetRawValue(data []byte) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetRawValue(&_DIAOracleV3.TransactOpts, data)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) SetValue(opts *bind.TransactOpts, key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "setValue", key, value, timestamp)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_DIAOracleV3 *DIAOracleV3Session) SetValue(key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetValue(&_DIAOracleV3.TransactOpts, key, value, timestamp)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) SetValue(key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.SetValue(&_DIAOracleV3.TransactOpts, key, value, timestamp)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DIAOracleV3 *DIAOracleV3Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DIAOracleV3.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DIAOracleV3 *DIAOracleV3Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.UpgradeToAndCall(&_DIAOracleV3.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DIAOracleV3 *DIAOracleV3TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DIAOracleV3.Contract.UpgradeToAndCall(&_DIAOracleV3.TransactOpts, newImplementation, data)
}

// DIAOracleV3InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DIAOracleV3 contract.
type DIAOracleV3InitializedIterator struct {
	Event *DIAOracleV3Initialized // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3Initialized)
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
		it.Event = new(DIAOracleV3Initialized)
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
func (it *DIAOracleV3InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3Initialized represents a Initialized event raised by the DIAOracleV3 contract.
type DIAOracleV3Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterInitialized(opts *bind.FilterOpts) (*DIAOracleV3InitializedIterator, error) {

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3InitializedIterator{contract: _DIAOracleV3.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DIAOracleV3Initialized) (event.Subscription, error) {

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3Initialized)
				if err := _DIAOracleV3.contract.UnpackLog(event, "Initialized", log); err != nil {
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseInitialized(log types.Log) (*DIAOracleV3Initialized, error) {
	event := new(DIAOracleV3Initialized)
	if err := _DIAOracleV3.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3OracleUpdateIterator is returned from FilterOracleUpdate and is used to iterate over the raw logs and unpacked data for OracleUpdate events raised by the DIAOracleV3 contract.
type DIAOracleV3OracleUpdateIterator struct {
	Event *DIAOracleV3OracleUpdate // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3OracleUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3OracleUpdate)
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
		it.Event = new(DIAOracleV3OracleUpdate)
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
func (it *DIAOracleV3OracleUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3OracleUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3OracleUpdate represents a OracleUpdate event raised by the DIAOracleV3 contract.
type DIAOracleV3OracleUpdate struct {
	Key       string
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleUpdate is a free log retrieval operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterOracleUpdate(opts *bind.FilterOpts) (*DIAOracleV3OracleUpdateIterator, error) {

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "OracleUpdate")
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3OracleUpdateIterator{contract: _DIAOracleV3.contract, event: "OracleUpdate", logs: logs, sub: sub}, nil
}

// WatchOracleUpdate is a free log subscription operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchOracleUpdate(opts *bind.WatchOpts, sink chan<- *DIAOracleV3OracleUpdate) (event.Subscription, error) {

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "OracleUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3OracleUpdate)
				if err := _DIAOracleV3.contract.UnpackLog(event, "OracleUpdate", log); err != nil {
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

// ParseOracleUpdate is a log parse operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseOracleUpdate(log types.Log) (*DIAOracleV3OracleUpdate, error) {
	event := new(DIAOracleV3OracleUpdate)
	if err := _DIAOracleV3.contract.UnpackLog(event, "OracleUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3OracleUpdateRawIterator is returned from FilterOracleUpdateRaw and is used to iterate over the raw logs and unpacked data for OracleUpdateRaw events raised by the DIAOracleV3 contract.
type DIAOracleV3OracleUpdateRawIterator struct {
	Event *DIAOracleV3OracleUpdateRaw // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3OracleUpdateRawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3OracleUpdateRaw)
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
		it.Event = new(DIAOracleV3OracleUpdateRaw)
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
func (it *DIAOracleV3OracleUpdateRawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3OracleUpdateRawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3OracleUpdateRaw represents a OracleUpdateRaw event raised by the DIAOracleV3 contract.
type DIAOracleV3OracleUpdateRaw struct {
	Key       string
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
	Data      []byte
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleUpdateRaw is a free log retrieval operation binding the contract event 0x0ec1e0298284e066eddd5e448f165c9337bf3f9447b7159177c72e0cada227d3.
//
// Solidity: event OracleUpdateRaw(string key, uint128 value, uint128 timestamp, uint128 volume, bytes data)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterOracleUpdateRaw(opts *bind.FilterOpts) (*DIAOracleV3OracleUpdateRawIterator, error) {

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "OracleUpdateRaw")
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3OracleUpdateRawIterator{contract: _DIAOracleV3.contract, event: "OracleUpdateRaw", logs: logs, sub: sub}, nil
}

// WatchOracleUpdateRaw is a free log subscription operation binding the contract event 0x0ec1e0298284e066eddd5e448f165c9337bf3f9447b7159177c72e0cada227d3.
//
// Solidity: event OracleUpdateRaw(string key, uint128 value, uint128 timestamp, uint128 volume, bytes data)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchOracleUpdateRaw(opts *bind.WatchOpts, sink chan<- *DIAOracleV3OracleUpdateRaw) (event.Subscription, error) {

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "OracleUpdateRaw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3OracleUpdateRaw)
				if err := _DIAOracleV3.contract.UnpackLog(event, "OracleUpdateRaw", log); err != nil {
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

// ParseOracleUpdateRaw is a log parse operation binding the contract event 0x0ec1e0298284e066eddd5e448f165c9337bf3f9447b7159177c72e0cada227d3.
//
// Solidity: event OracleUpdateRaw(string key, uint128 value, uint128 timestamp, uint128 volume, bytes data)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseOracleUpdateRaw(log types.Log) (*DIAOracleV3OracleUpdateRaw, error) {
	event := new(DIAOracleV3OracleUpdateRaw)
	if err := _DIAOracleV3.contract.UnpackLog(event, "OracleUpdateRaw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the DIAOracleV3 contract.
type DIAOracleV3RoleAdminChangedIterator struct {
	Event *DIAOracleV3RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3RoleAdminChanged)
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
		it.Event = new(DIAOracleV3RoleAdminChanged)
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
func (it *DIAOracleV3RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3RoleAdminChanged represents a RoleAdminChanged event raised by the DIAOracleV3 contract.
type DIAOracleV3RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*DIAOracleV3RoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3RoleAdminChangedIterator{contract: _DIAOracleV3.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *DIAOracleV3RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3RoleAdminChanged)
				if err := _DIAOracleV3.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseRoleAdminChanged(log types.Log) (*DIAOracleV3RoleAdminChanged, error) {
	event := new(DIAOracleV3RoleAdminChanged)
	if err := _DIAOracleV3.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the DIAOracleV3 contract.
type DIAOracleV3RoleGrantedIterator struct {
	Event *DIAOracleV3RoleGranted // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3RoleGranted)
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
		it.Event = new(DIAOracleV3RoleGranted)
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
func (it *DIAOracleV3RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3RoleGranted represents a RoleGranted event raised by the DIAOracleV3 contract.
type DIAOracleV3RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DIAOracleV3RoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3RoleGrantedIterator{contract: _DIAOracleV3.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *DIAOracleV3RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3RoleGranted)
				if err := _DIAOracleV3.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseRoleGranted(log types.Log) (*DIAOracleV3RoleGranted, error) {
	event := new(DIAOracleV3RoleGranted)
	if err := _DIAOracleV3.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the DIAOracleV3 contract.
type DIAOracleV3RoleRevokedIterator struct {
	Event *DIAOracleV3RoleRevoked // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3RoleRevoked)
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
		it.Event = new(DIAOracleV3RoleRevoked)
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
func (it *DIAOracleV3RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3RoleRevoked represents a RoleRevoked event raised by the DIAOracleV3 contract.
type DIAOracleV3RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DIAOracleV3RoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3RoleRevokedIterator{contract: _DIAOracleV3.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *DIAOracleV3RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3RoleRevoked)
				if err := _DIAOracleV3.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseRoleRevoked(log types.Log) (*DIAOracleV3RoleRevoked, error) {
	event := new(DIAOracleV3RoleRevoked)
	if err := _DIAOracleV3.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3UpdaterAddressChangeIterator is returned from FilterUpdaterAddressChange and is used to iterate over the raw logs and unpacked data for UpdaterAddressChange events raised by the DIAOracleV3 contract.
type DIAOracleV3UpdaterAddressChangeIterator struct {
	Event *DIAOracleV3UpdaterAddressChange // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3UpdaterAddressChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3UpdaterAddressChange)
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
		it.Event = new(DIAOracleV3UpdaterAddressChange)
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
func (it *DIAOracleV3UpdaterAddressChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3UpdaterAddressChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3UpdaterAddressChange represents a UpdaterAddressChange event raised by the DIAOracleV3 contract.
type DIAOracleV3UpdaterAddressChange struct {
	NewUpdater common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdaterAddressChange is a free log retrieval operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterUpdaterAddressChange(opts *bind.FilterOpts) (*DIAOracleV3UpdaterAddressChangeIterator, error) {

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "UpdaterAddressChange")
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3UpdaterAddressChangeIterator{contract: _DIAOracleV3.contract, event: "UpdaterAddressChange", logs: logs, sub: sub}, nil
}

// WatchUpdaterAddressChange is a free log subscription operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchUpdaterAddressChange(opts *bind.WatchOpts, sink chan<- *DIAOracleV3UpdaterAddressChange) (event.Subscription, error) {

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "UpdaterAddressChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3UpdaterAddressChange)
				if err := _DIAOracleV3.contract.UnpackLog(event, "UpdaterAddressChange", log); err != nil {
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

// ParseUpdaterAddressChange is a log parse operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseUpdaterAddressChange(log types.Log) (*DIAOracleV3UpdaterAddressChange, error) {
	event := new(DIAOracleV3UpdaterAddressChange)
	if err := _DIAOracleV3.contract.UnpackLog(event, "UpdaterAddressChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DIAOracleV3UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the DIAOracleV3 contract.
type DIAOracleV3UpgradedIterator struct {
	Event *DIAOracleV3Upgraded // Event containing the contract specifics and raw log

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
func (it *DIAOracleV3UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DIAOracleV3Upgraded)
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
		it.Event = new(DIAOracleV3Upgraded)
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
func (it *DIAOracleV3UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DIAOracleV3UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DIAOracleV3Upgraded represents a Upgraded event raised by the DIAOracleV3 contract.
type DIAOracleV3Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DIAOracleV3 *DIAOracleV3Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DIAOracleV3UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DIAOracleV3.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DIAOracleV3UpgradedIterator{contract: _DIAOracleV3.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DIAOracleV3 *DIAOracleV3Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DIAOracleV3Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DIAOracleV3.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DIAOracleV3Upgraded)
				if err := _DIAOracleV3.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_DIAOracleV3 *DIAOracleV3Filterer) ParseUpgraded(log types.Log) (*DIAOracleV3Upgraded, error) {
	event := new(DIAOracleV3Upgraded)
	if err := _DIAOracleV3.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
