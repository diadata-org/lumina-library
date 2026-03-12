// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package diaOracleV3

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

// DiaOracleV3MetaData contains all meta data concerning the DiaOracleV3 contract.
var DiaOracleV3MetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"DEFAULT_ADMIN_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_ALLOWED_HISTORY_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_HISTORY_SIZE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_TIMESTAMP_GAP\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPDATER_ROLE\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"UPGRADE_INTERFACE_VERSION\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"decimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getDecimals\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMaxHistorySize\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRawData\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getRoleAdmin\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValue\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValueAt\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"value\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"volume\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValueCount\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getValueHistory\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple[]\",\"internalType\":\"structIDIAOracleV3.ValueEntry[]\",\"components\":[{\"name\":\"value\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"volume\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"grantRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"hasRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"proxiableUUID\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"rawData\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"renounceRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"callerConfirmation\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"revokeRole\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setDecimals\",\"inputs\":[{\"name\":\"decimalPrecision\",\"type\":\"uint8\",\"internalType\":\"uint8\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMultipleRawValues\",\"inputs\":[{\"name\":\"dataArray\",\"type\":\"bytes[]\",\"internalType\":\"bytes[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setMultipleValues\",\"inputs\":[{\"name\":\"keys\",\"type\":\"string[]\",\"internalType\":\"string[]\"},{\"name\":\"compressedValues\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setRawValue\",\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"setValue\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"supportsInterface\",\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\",\"internalType\":\"bytes4\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"upgradeToAndCall\",\"inputs\":[{\"name\":\"newImplementation\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\",\"internalType\":\"bytes\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"function\",\"name\":\"values\",\"inputs\":[{\"name\":\"\",\"type\":\"string\",\"internalType\":\"string\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"DecimalsUpdate\",\"inputs\":[{\"name\":\"decimals\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"uint8\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Initialized\",\"inputs\":[{\"name\":\"version\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleUpdate\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OracleUpdateRaw\",\"inputs\":[{\"name\":\"key\",\"type\":\"string\",\"indexed\":false,\"internalType\":\"string\"},{\"name\":\"value\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"timestamp\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"volume\",\"type\":\"uint128\",\"indexed\":false,\"internalType\":\"uint128\"},{\"name\":\"data\",\"type\":\"bytes\",\"indexed\":false,\"internalType\":\"bytes\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleAdminChanged\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"previousAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"newAdminRole\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleGranted\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RoleRevoked\",\"inputs\":[{\"name\":\"role\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"account\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"sender\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"UpdaterAddressChange\",\"inputs\":[{\"name\":\"newUpdater\",\"type\":\"address\",\"indexed\":false,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Upgraded\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AccessControlBadConfirmation\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AccessControlUnauthorizedAccount\",\"inputs\":[{\"name\":\"account\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"neededRole\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]},{\"type\":\"error\",\"name\":\"AddressEmptyCode\",\"inputs\":[{\"name\":\"target\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967InvalidImplementation\",\"inputs\":[{\"name\":\"implementation\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"ERC1967NonPayable\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FailedCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidHistoryIndex\",\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"maxIndex\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"InvalidInitialization\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"MismatchedArrayLengths\",\"inputs\":[{\"name\":\"keysLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"valuesLength\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"NotInitializing\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimestampNotIncreasing\",\"inputs\":[{\"name\":\"newTimestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"existingTimestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"}]},{\"type\":\"error\",\"name\":\"TimestampTooFarInFuture\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"blockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"TimestampTooFarInPast\",\"inputs\":[{\"name\":\"timestamp\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"blockTime\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]},{\"type\":\"error\",\"name\":\"UUPSUnauthorizedCallContext\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"UUPSUnsupportedProxiableUUID\",\"inputs\":[{\"name\":\"slot\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}]}]",
	Bin: "0x60c080604052346100f35730608052606460a0525f5160206121575f395f51905f525460ff8160401c166100e4576002600160401b03196001600160401b03821601610091575b60405161205f90816100f882396080518181816108e1015261097f015260a0518181816111d50152818161134b0152818161171f01528181611a6501528181611adb0152611be80152f35b6001600160401b0319166001600160401b039081175f5160206121575f395f51905f525581527fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d290602090a15f80610046565b63f92ee8a960e01b5f5260045ffd5b5f80fdfe6080806040526004361015610012575f80fd5b5f3560e01c90816301ffc9a7146112385750806309daaa95146110d7578063135d90c714610e7a578063248a9ca314610e3b57806324b1db1a146101db5780632c484ae514610c185780632f2ff15d14610bce578063313ce5671461016657806336568abe14610b8a57806347e6338014610b505780634df71096146106b15780634ed696fc14610b345780634f1ef2861461093557806352d1902d146108cf57806359c3852c1461087a5780635a9ade8b146108265780637898e0c21461077b5780637a1395aa1461071b5780637a2fa442146106b15780638129fc1c146105bf5780638d241526146103bb5780638d97ecf21461031b57806391d14854146102c6578063960384a014610261578063a217fddf14610247578063ad3cb1cc146101fc578063be2e0179146101e0578063c1066251146101db578063d547741f1461018a5763f0141d8414610166575f80fd5b34610186575f36600319011261018657602060ff60055416604051908152f35b5f80fd5b34610186576040366003190112610186576101d96004356101a961136e565b906101d46101cf825f525f516020611fca5f395f51905f52602052600160405f20015490565b611854565b611e81565b005b611334565b34610186575f366003190112610186576020604051610e108152f35b34610186575f3660031901126101865761024360405161021d6040826112c0565b60058152640352e302e360dc1b6020820152604051918291602083526020830190611384565b0390f35b34610186575f3660031901126101865760206040515f8152f35b34610186576020366003190112610186576004356001600160401b0381116101865760208061029660409336906004016112e1565b8351928184925191829101835e81015f815203019020546001600160801b038251918060801c8352166020820152f35b34610186576040366003190112610186576102df61136e565b6004355f525f516020611fca5f395f51905f5260205260405f209060018060a01b03165f52602052602060ff60405f2054166040519015158152f35b34610186576020366003190112610186576004356001600160401b0381116101865761034e6103539136906004016112e1565b611638565b6040518091602082016020835281518091526020604084019201905f5b81811061037e575050500390f35b91935091602060606001926001600160801b0360408851828151168452828682015116868501520151166040820152019401910191849392610370565b34610186576040366003190112610186576004356001600160401b03811161018657366023820112156101865780600401356103f681611494565b9161040460405193846112c0565b8183526024602084019260051b820101903682116101865760248101925b82841061059057846024356001600160401b03811161018657366023820112156101865780600401359061045582611494565b9161046360405193846112c0565b8083526024602084019160051b8301019136831161018657602401905b828210610580575050506104926117e5565b815181519081810361056b5750505f5b8251906001600160801b038116918210156101d9577fa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782826104f66104ee6001600160801b039588611624565b519186611624565b519061053d8260801c928681169061050e828561189a565b604051602081865180838901835e81015f815203019020556105318185856119c6565b604051938493846115f6565b0390a1166001600160801b038114610557576001016104a2565b634e487b7160e01b5f52601160045260245ffd5b633adf150f60e21b5f5260045260245260445ffd5b8135815260209182019101610480565b83356001600160401b038111610186576020916105b48392602436918701016112e1565b815201930192610422565b34610186575f366003190112610186575f516020611fea5f395f51905f525460ff8160401c16801561069d575b61068e5768ffffffffffffffffff191668010000000000000001175f516020611fea5f395f51905f52556005805460ff1916600817905561062c33611c62565b5061063633611d11565b5068ff0000000000000000195f516020611fea5f395f51905f5254165f516020611fea5f395f51905f52557fc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2602060405160018152a1005b63f92ee8a960e01b5f5260045ffd5b5060016001600160401b03821610156105ec565b34610186576020366003190112610186576004356001600160401b038111610186576107076020806106ea6102439436906004016112e1565b604051928184925191829101835e810160048152030190206113f4565b604051918291602083526020830190611384565b346101865760203660031901126101865760043560ff81168091036101865760207fc03a161db992c7f007555d2c71b8de00760358129d9547bce9475bff4739bf60916107666117e5565b8060ff196005541617600555604051908152a1005b34610186576060366003190112610186576004356001600160401b038111610186576107ab9036906004016112e1565b6024356001600160801b038116810361018657604435916001600160801b038316808403610186577fa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee22917829361050e610821926108046117e5565b61080e838661189a565b6001600160801b03198660801b166114ab565b0390a1005b34610186576020366003190112610186576004356001600160401b0381116101865760208061085a819336906004016112e1565b604051928184925191829101835e81015f81520301902054604051908152f35b34610186576020366003190112610186576004356001600160401b038111610186576020806108ae819336906004016112e1565b604051928184925191829101835e8101600381520301902054604051908152f35b34610186575f366003190112610186577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031630036109265760206040515f516020611faa5f395f51905f528152f35b63703e46dd60e11b5f5260045ffd5b6040366003190112610186576004356001600160a01b03811690818103610186576024356001600160401b038111610186576109759036906004016112e1565b6001600160a01b037f000000000000000000000000000000000000000000000000000000000000000016308114908115610b12575b5061092657335f9081527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d602052604090205460ff1615610afb576040516352d1902d60e01b8152602081600481875afa5f9181610ac7575b50610a1b5783634c9c8ce360e01b5f5260045260245ffd5b805f516020611faa5f395f51905f52859203610ab55750823b15610aa3575f516020611faa5f395f51905f5280546001600160a01b031916821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2805115610a8b576101d991611f1d565b505034610a9457005b63b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b632a87526960e21b5f5260045260245ffd5b9091506020813d602011610af3575b81610ae3602093836112c0565b8101031261018657519085610a03565b3d9150610ad6565b63e2517d3f60e01b5f52336004525f60245260445ffd5b5f516020611faa5f395f51905f52546001600160a01b031614159050846109aa565b34610186575f3660031901126101865760206040516103e88152f35b34610186575f3660031901126101865760206040517f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab8152f35b3461018657604036600319011261018657610ba361136e565b336001600160a01b03821603610bbf576101d990600435611e81565b63334bd91960e11b5f5260045ffd5b34610186576040366003190112610186576101d9600435610bed61136e565b90610c136101cf825f525f516020611fca5f395f51905f52602052600160405f20015490565b611ddd565b34610186576020366003190112610186576004356001600160401b03811161018657366023820112156101865780600401356001600160401b038111610186578101602401368111610186576024610c7892610c726117e5565b01611542565b90610c858386949661189a565b610ca56001600160801b0386166001600160801b03198660801b166114ab565b60405190845191602081818801948086835e81015f815203019020556020604051809286518091835e81016004815203019020948251956001600160401b038711610e2757610cf481546113bc565b601f8111610dd6575b50602096601f8111600114610d61579081610821959493925f51602061200a5f395f51905f52995f91610d56575b508160011b915f199060031b1c19161790555b610d4a82828888611b54565b604051958695866115b3565b90508501518a610d2b565b601f19811697825f52805f20985f5b818110610dbe5750915f51602061200a5f395f51905f529960019282610821999897969510610da6575b5050811b019055610d3e565b8701515f1960f88460031b161c191690558a80610d9a565b878301518b556001909a019960209283019201610d70565b87811115610cfd57815f5260205f20601f890160051c9060208a10610e1f575b81601f9101920160051c03905f5b828110610e12575050610cfd565b5f82820155600101610e04565b5f9150610df6565b634e487b7160e01b5f52604160045260245ffd5b34610186576020366003190112610186576020610e726004355f525f516020611fca5f395f51905f52602052600160405f20015490565b604051908152f35b34610186576020366003190112610186576004356001600160401b03811161018657366023820112156101865780600401356001600160401b038111610186573660248260051b840101116101865790610ed26117e5565b3681900360421901905f5b838110156101d95760248160051b830101358381121561018657820160248101356001600160401b038111610186576044820190803603821361018657610f28920160440190611542565b91610f358186959661189a565b610f556001600160801b0382166001600160801b03198760801b166114ab565b60405190855191602081818901948086835e81015f815203019020556020604051809287518091835e8101600481520301902083516001600160401b038111610e2757610fa282546113bc565b601f8111611086575b506020601f821160011461100e5792600198979592825f51602061200a5f395f51905f52989693610ffa965f91611003575b505f19600383901b1c1916908b1b179055610d4a82828888611b54565b0390a101610edd565b90508501518f610fdd565b601f19821690835f52805f20915f5b81811061106e575083610ffa969360019c9b9996935f51602061200a5f395f51905f529b99968e9410611056575050811b019055610d3e565b8701515f1960f88460031b161c191690558f80610d9a565b9192602060018192868c01518155019401920161101d565b81811115610fab57825f5260205f20601f830160051c90602084106110cf575b81601f9101920160051c03905f5b8281106110c2575050610fab565b5f828201556001016110b4565b5f91506110a6565b34610186576040366003190112610186576004356001600160401b038111610186576111079036906004016112e1565b602435604051825190602081818601938085835e810160018152030190209260405160208183518086835e8101600381520301902054808410156112225750602090604051928391518091835e810160028152030190205490600181018082116105575782106111cf575f1982019182116105575760609261118f6111959261119b946114b8565b906114e3565b50611510565b6001600160801b03815116906001600160801b03604081602084015116920151169060405192835260208301526040820152f35b916111fb7f000000000000000000000000000000000000000000000000000000000000000080936114ab565b5f198101919082116105575761118f61119b9361121d606096611195956114b8565b6114c5565b8363c9d622b960e01b5f5260045260245260445ffd5b34610186576020366003190112610186576004359063ffffffff60e01b8216809203610186576020916362652b8f60e01b811490811561127a575b5015158152f35b637965db0b60e01b811491508115611294575b5083611273565b6301ffc9a760e01b1490508361128d565b606081019081106001600160401b03821117610e2757604052565b90601f801991011681019081106001600160401b03821117610e2757604052565b81601f82011215610186576020813591016001600160401b038211610e275760405192611318601f8401601f1916602001856112c0565b8284528282011161018657815f92602092838601378301015290565b34610186575f3660031901126101865760206040517f00000000000000000000000000000000000000000000000000000000000000008152f35b602435906001600160a01b038216820361018657565b805180835260209291819084018484015e5f828201840152601f01601f1916010190565b35906001600160801b038216820361018657565b90600182811c921680156113ea575b60208310146113d657565b634e487b7160e01b5f52602260045260245ffd5b91607f16916113cb565b9060405191825f825492611407846113bc565b8084529360018116908115611472575060011461142e575b5061142c925003836112c0565b565b90505f9291925260205f20905f915b81831061145657505090602061142c928201015f61141f565b602091935080600191548385890101520191019091849261143d565b90506020925061142c94915060ff191682840152151560051b8201015f61141f565b6001600160401b038111610e275760051b60200190565b9190820180921161055757565b9190820391821161055757565b81156114cf570690565b634e487b7160e01b5f52601260045260245ffd5b80548210156114fc575f5260205f209060011b01905f90565b634e487b7160e01b5f52603260045260245ffd5b9060405161151d816112a5565b60406001600160801b03600183958054838116865260801c6020860152015416910152565b91909160a0818403126101865780356001600160401b038111610186578361156b9183016112e1565b92611578602083016113a8565b92611585604084016113a8565b92611592606082016113a8565b9260808201356001600160401b038111610186576115b092016112e1565b90565b91936001600160801b036115b09694816115d6819560a0885260a0880190611384565b971660208601521660408401521660608201526080818403910152611384565b9193926001600160801b039081611617604094606087526060870190611384565b9616602085015216910152565b80518210156114fc5760209160051b010190565b604051815190602081818501938085835e8101600181520301902060405160208185518086835e81016003815203019020549182156117985761167a83611494565b9361168860405195866112c0565b838552601f1961169785611494565b015f5b81811061176f575050602090604051928391518091835e81016002815203019020545f5b8381106116cc575050505090565b6001810180821161055757821061171d575f19820190828211610557576117016111956116fb836001956114b8565b866114e3565b61170b8288611624565b526117168187611624565b50016116be565b7f00000000000000000000000000000000000000000000000000000000000000009061174982846114ab565b5f1981019081116105575761119561176a60019461121d85611701956114b8565b6116fb565b60209060405161177e816112a5565b5f81525f838201525f604082015282828a0101520161169a565b505050506040516117aa6020826112c0565b5f81525f805b8181106117bc57505090565b6020906040516117cb816112a5565b5f81525f838201525f6040820152828286010152016117b0565b335f9081527f268477d1c17bf7595397186358fec802c44c92bb39292de205ad03cbfe096e02602052604090205460ff161561181d57565b63e2517d3f60e01b5f52336004527f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab60245260445ffd5b5f8181525f516020611fca5f395f51905f526020908152604080832033845290915290205460ff16156118845750565b63e2517d3f60e01b5f523360045260245260445ffd5b610e10420191824211610557576001600160801b038091169216821161194d57610e1042118061192f575b6119185760208091604051928184925191829101835e81015f81520301902054806118ee575050565b6001600160801b031690818110611903575050565b633e4ebc1960e21b5f5260045260245260445ffd5b5063514a513560e01b5f526004524260245260445ffd5b50610e0f194201428111610557576001600160801b031682106118c5565b506304c3bb5b60e21b5f526004524260245260445ffd5b91906119b3578051602082015160801b6fffffffffffffffffffffffffffffffff199081166001600160801b039283161784556040929092015160019390930180549092169216919091179055565b634e487b7160e01b5f525f60045260245ffd5b60405191815192602081818501958087835e810160018152030190209060405160208185518088835e81016002815203019020549160405160208186518089835e810160038152030190205495815415611ad5575b83611a5593926001600160801b03611a4f938160405196611a3b886112a5565b1686521660208501525f60408501526114e3565b90611964565b6001810180911161055757611a8b7f000000000000000000000000000000000000000000000000000000000000000080926114c5565b60405160208185518088835e81016002815203019020558310611aad57505050565b5f198314610557576020906001604051938492518091845e8201946003865201930301902055565b955091507f00000000000000000000000000000000000000000000000000000000000000005f5b818110611b0f57505f9586939150611a1b565b60405190611b1c826112a5565b5f82525f60208301525f60408301528454600160401b811015610e2757600192611a4f8285611b4e94018955886114e3565b01611afc565b90919260405192825193602081818601968088835e8101600181520301902060405160208186518089835e8101600281520301902054926040516020818751808a835e810160038152030190205496825415611be1575b611a5593926001600160801b03611a4f938188948160405198611bcd8a6112a5565b1688521660208701521660408501526114e3565b96509192507f00000000000000000000000000000000000000000000000000000000000000005f5b818110611c1d57505f968794939150611bab565b60405190611c2a826112a5565b5f82525f60208301525f60408301528354600160401b811015610e2757600192611a4f8285611c5c94018855876114e3565b01611c09565b6001600160a01b0381165f9081527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d602052604090205460ff16611d0c576001600160a01b03165f8181527fb7db2dd08fcb62d0c9e08c51941cae53c267786a0b75803fb7960902fc8ef97d60205260408120805460ff191660011790553391907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d8180a4600190565b505f90565b6001600160a01b0381165f9081527f268477d1c17bf7595397186358fec802c44c92bb39292de205ad03cbfe096e02602052604090205460ff16611d0c576001600160a01b03165f8181527f268477d1c17bf7595397186358fec802c44c92bb39292de205ad03cbfe096e0260205260408120805460ff191660011790553391907f73e573f9566d61418a34d5de3ff49360f9c51fec37f7486551670290f6285dab907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b5f8181525f516020611fca5f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff16611e7b575f8181525f516020611fca5f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19166001179055339291907f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d9080a4600190565b50505f90565b5f8181525f516020611fca5f395f51905f52602090815260408083206001600160a01b038616845290915290205460ff1615611e7b575f8181525f516020611fca5f395f51905f52602090815260408083206001600160a01b0395909516808452949091528120805460ff19169055339291907ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9080a4600190565b905f8091602081519101845af48080611f96575b15611f515750506040513d81523d5f602083013e60203d82010160405290565b15611f7657639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b3d15611f87576040513d5f823e3d90fd5b63d6bda27560e01b5f5260045ffd5b503d151580611f315750813b1515611f3156fe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc02dd7bc7dec4dceedda775e58dd541e08a116c6c53815c0bd028192f7b626800f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a000ec1e0298284e066eddd5e448f165c9337bf3f9447b7159177c72e0cada227d3a2646970667358221220d7e782cdf4ef8118152e941fc26324f4f79956cdbd46637119ddcbea2fc0312564736f6c63430008220033f0c57e16840df040f15088dc2f81fe391c3923bec73e23a9662efc9c229c6a00",
}

// DiaOracleV3ABI is the input ABI used to generate the binding from.
// Deprecated: Use DiaOracleV3MetaData.ABI instead.
var DiaOracleV3ABI = DiaOracleV3MetaData.ABI

// DiaOracleV3Bin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use DiaOracleV3MetaData.Bin instead.
var DiaOracleV3Bin = DiaOracleV3MetaData.Bin

// DeployDiaOracleV3 deploys a new Ethereum contract, binding an instance of DiaOracleV3 to it.
func DeployDiaOracleV3(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *DiaOracleV3, error) {
	parsed, err := DiaOracleV3MetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(DiaOracleV3Bin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &DiaOracleV3{DiaOracleV3Caller: DiaOracleV3Caller{contract: contract}, DiaOracleV3Transactor: DiaOracleV3Transactor{contract: contract}, DiaOracleV3Filterer: DiaOracleV3Filterer{contract: contract}}, nil
}

// DiaOracleV3 is an auto generated Go binding around an Ethereum contract.
type DiaOracleV3 struct {
	DiaOracleV3Caller     // Read-only binding to the contract
	DiaOracleV3Transactor // Write-only binding to the contract
	DiaOracleV3Filterer   // Log filterer for contract events
}

// DiaOracleV3Caller is an auto generated read-only Go binding around an Ethereum contract.
type DiaOracleV3Caller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DiaOracleV3Transactor is an auto generated write-only Go binding around an Ethereum contract.
type DiaOracleV3Transactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DiaOracleV3Filterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DiaOracleV3Filterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DiaOracleV3Session is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DiaOracleV3Session struct {
	Contract     *DiaOracleV3      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DiaOracleV3CallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DiaOracleV3CallerSession struct {
	Contract *DiaOracleV3Caller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// DiaOracleV3TransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DiaOracleV3TransactorSession struct {
	Contract     *DiaOracleV3Transactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// DiaOracleV3Raw is an auto generated low-level Go binding around an Ethereum contract.
type DiaOracleV3Raw struct {
	Contract *DiaOracleV3 // Generic contract binding to access the raw methods on
}

// DiaOracleV3CallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DiaOracleV3CallerRaw struct {
	Contract *DiaOracleV3Caller // Generic read-only contract binding to access the raw methods on
}

// DiaOracleV3TransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DiaOracleV3TransactorRaw struct {
	Contract *DiaOracleV3Transactor // Generic write-only contract binding to access the raw methods on
}

// NewDiaOracleV3 creates a new instance of DiaOracleV3, bound to a specific deployed contract.
func NewDiaOracleV3(address common.Address, backend bind.ContractBackend) (*DiaOracleV3, error) {
	contract, err := bindDiaOracleV3(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3{DiaOracleV3Caller: DiaOracleV3Caller{contract: contract}, DiaOracleV3Transactor: DiaOracleV3Transactor{contract: contract}, DiaOracleV3Filterer: DiaOracleV3Filterer{contract: contract}}, nil
}

// NewDiaOracleV3Caller creates a new read-only instance of DiaOracleV3, bound to a specific deployed contract.
func NewDiaOracleV3Caller(address common.Address, caller bind.ContractCaller) (*DiaOracleV3Caller, error) {
	contract, err := bindDiaOracleV3(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3Caller{contract: contract}, nil
}

// NewDiaOracleV3Transactor creates a new write-only instance of DiaOracleV3, bound to a specific deployed contract.
func NewDiaOracleV3Transactor(address common.Address, transactor bind.ContractTransactor) (*DiaOracleV3Transactor, error) {
	contract, err := bindDiaOracleV3(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3Transactor{contract: contract}, nil
}

// NewDiaOracleV3Filterer creates a new log filterer instance of DiaOracleV3, bound to a specific deployed contract.
func NewDiaOracleV3Filterer(address common.Address, filterer bind.ContractFilterer) (*DiaOracleV3Filterer, error) {
	contract, err := bindDiaOracleV3(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3Filterer{contract: contract}, nil
}

// bindDiaOracleV3 binds a generic wrapper to an already deployed contract.
func bindDiaOracleV3(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := DiaOracleV3MetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DiaOracleV3 *DiaOracleV3Raw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DiaOracleV3.Contract.DiaOracleV3Caller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DiaOracleV3 *DiaOracleV3Raw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.DiaOracleV3Transactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DiaOracleV3 *DiaOracleV3Raw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.DiaOracleV3Transactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DiaOracleV3 *DiaOracleV3CallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DiaOracleV3.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DiaOracleV3 *DiaOracleV3TransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DiaOracleV3 *DiaOracleV3TransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Caller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Session) DEFAULTADMINROLE() ([32]byte, error) {
	return _DiaOracleV3.Contract.DEFAULTADMINROLE(&_DiaOracleV3.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3CallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _DiaOracleV3.Contract.DEFAULTADMINROLE(&_DiaOracleV3.CallOpts)
}

// MAXALLOWEDHISTORYSIZE is a free data retrieval call binding the contract method 0x4ed696fc.
//
// Solidity: function MAX_ALLOWED_HISTORY_SIZE() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Caller) MAXALLOWEDHISTORYSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "MAX_ALLOWED_HISTORY_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXALLOWEDHISTORYSIZE is a free data retrieval call binding the contract method 0x4ed696fc.
//
// Solidity: function MAX_ALLOWED_HISTORY_SIZE() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Session) MAXALLOWEDHISTORYSIZE() (*big.Int, error) {
	return _DiaOracleV3.Contract.MAXALLOWEDHISTORYSIZE(&_DiaOracleV3.CallOpts)
}

// MAXALLOWEDHISTORYSIZE is a free data retrieval call binding the contract method 0x4ed696fc.
//
// Solidity: function MAX_ALLOWED_HISTORY_SIZE() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3CallerSession) MAXALLOWEDHISTORYSIZE() (*big.Int, error) {
	return _DiaOracleV3.Contract.MAXALLOWEDHISTORYSIZE(&_DiaOracleV3.CallOpts)
}

// MAXHISTORYSIZE is a free data retrieval call binding the contract method 0xc1066251.
//
// Solidity: function MAX_HISTORY_SIZE() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Caller) MAXHISTORYSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "MAX_HISTORY_SIZE")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXHISTORYSIZE is a free data retrieval call binding the contract method 0xc1066251.
//
// Solidity: function MAX_HISTORY_SIZE() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Session) MAXHISTORYSIZE() (*big.Int, error) {
	return _DiaOracleV3.Contract.MAXHISTORYSIZE(&_DiaOracleV3.CallOpts)
}

// MAXHISTORYSIZE is a free data retrieval call binding the contract method 0xc1066251.
//
// Solidity: function MAX_HISTORY_SIZE() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3CallerSession) MAXHISTORYSIZE() (*big.Int, error) {
	return _DiaOracleV3.Contract.MAXHISTORYSIZE(&_DiaOracleV3.CallOpts)
}

// MAXTIMESTAMPGAP is a free data retrieval call binding the contract method 0xbe2e0179.
//
// Solidity: function MAX_TIMESTAMP_GAP() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Caller) MAXTIMESTAMPGAP(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "MAX_TIMESTAMP_GAP")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// MAXTIMESTAMPGAP is a free data retrieval call binding the contract method 0xbe2e0179.
//
// Solidity: function MAX_TIMESTAMP_GAP() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Session) MAXTIMESTAMPGAP() (*big.Int, error) {
	return _DiaOracleV3.Contract.MAXTIMESTAMPGAP(&_DiaOracleV3.CallOpts)
}

// MAXTIMESTAMPGAP is a free data retrieval call binding the contract method 0xbe2e0179.
//
// Solidity: function MAX_TIMESTAMP_GAP() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3CallerSession) MAXTIMESTAMPGAP() (*big.Int, error) {
	return _DiaOracleV3.Contract.MAXTIMESTAMPGAP(&_DiaOracleV3.CallOpts)
}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Caller) UPDATERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "UPDATER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Session) UPDATERROLE() ([32]byte, error) {
	return _DiaOracleV3.Contract.UPDATERROLE(&_DiaOracleV3.CallOpts)
}

// UPDATERROLE is a free data retrieval call binding the contract method 0x47e63380.
//
// Solidity: function UPDATER_ROLE() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3CallerSession) UPDATERROLE() ([32]byte, error) {
	return _DiaOracleV3.Contract.UPDATERROLE(&_DiaOracleV3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DiaOracleV3 *DiaOracleV3Caller) UPGRADEINTERFACEVERSION(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "UPGRADE_INTERFACE_VERSION")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DiaOracleV3 *DiaOracleV3Session) UPGRADEINTERFACEVERSION() (string, error) {
	return _DiaOracleV3.Contract.UPGRADEINTERFACEVERSION(&_DiaOracleV3.CallOpts)
}

// UPGRADEINTERFACEVERSION is a free data retrieval call binding the contract method 0xad3cb1cc.
//
// Solidity: function UPGRADE_INTERFACE_VERSION() view returns(string)
func (_DiaOracleV3 *DiaOracleV3CallerSession) UPGRADEINTERFACEVERSION() (string, error) {
	return _DiaOracleV3.Contract.UPGRADEINTERFACEVERSION(&_DiaOracleV3.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DiaOracleV3 *DiaOracleV3Caller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DiaOracleV3 *DiaOracleV3Session) Decimals() (uint8, error) {
	return _DiaOracleV3.Contract.Decimals(&_DiaOracleV3.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() view returns(uint8)
func (_DiaOracleV3 *DiaOracleV3CallerSession) Decimals() (uint8, error) {
	return _DiaOracleV3.Contract.Decimals(&_DiaOracleV3.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_DiaOracleV3 *DiaOracleV3Caller) GetDecimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getDecimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_DiaOracleV3 *DiaOracleV3Session) GetDecimals() (uint8, error) {
	return _DiaOracleV3.Contract.GetDecimals(&_DiaOracleV3.CallOpts)
}

// GetDecimals is a free data retrieval call binding the contract method 0xf0141d84.
//
// Solidity: function getDecimals() view returns(uint8)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetDecimals() (uint8, error) {
	return _DiaOracleV3.Contract.GetDecimals(&_DiaOracleV3.CallOpts)
}

// GetMaxHistorySize is a free data retrieval call binding the contract method 0x24b1db1a.
//
// Solidity: function getMaxHistorySize() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Caller) GetMaxHistorySize(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getMaxHistorySize")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetMaxHistorySize is a free data retrieval call binding the contract method 0x24b1db1a.
//
// Solidity: function getMaxHistorySize() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Session) GetMaxHistorySize() (*big.Int, error) {
	return _DiaOracleV3.Contract.GetMaxHistorySize(&_DiaOracleV3.CallOpts)
}

// GetMaxHistorySize is a free data retrieval call binding the contract method 0x24b1db1a.
//
// Solidity: function getMaxHistorySize() view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetMaxHistorySize() (*big.Int, error) {
	return _DiaOracleV3.Contract.GetMaxHistorySize(&_DiaOracleV3.CallOpts)
}

// GetRawData is a free data retrieval call binding the contract method 0x4df71096.
//
// Solidity: function getRawData(string key) view returns(bytes)
func (_DiaOracleV3 *DiaOracleV3Caller) GetRawData(opts *bind.CallOpts, key string) ([]byte, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getRawData", key)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetRawData is a free data retrieval call binding the contract method 0x4df71096.
//
// Solidity: function getRawData(string key) view returns(bytes)
func (_DiaOracleV3 *DiaOracleV3Session) GetRawData(key string) ([]byte, error) {
	return _DiaOracleV3.Contract.GetRawData(&_DiaOracleV3.CallOpts, key)
}

// GetRawData is a free data retrieval call binding the contract method 0x4df71096.
//
// Solidity: function getRawData(string key) view returns(bytes)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetRawData(key string) ([]byte, error) {
	return _DiaOracleV3.Contract.GetRawData(&_DiaOracleV3.CallOpts, key)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Caller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Session) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DiaOracleV3.Contract.GetRoleAdmin(&_DiaOracleV3.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _DiaOracleV3.Contract.GetRoleAdmin(&_DiaOracleV3.CallOpts, role)
}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_DiaOracleV3 *DiaOracleV3Caller) GetValue(opts *bind.CallOpts, key string) (*big.Int, *big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getValue", key)

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
func (_DiaOracleV3 *DiaOracleV3Session) GetValue(key string) (*big.Int, *big.Int, error) {
	return _DiaOracleV3.Contract.GetValue(&_DiaOracleV3.CallOpts, key)
}

// GetValue is a free data retrieval call binding the contract method 0x960384a0.
//
// Solidity: function getValue(string key) view returns(uint128, uint128)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetValue(key string) (*big.Int, *big.Int, error) {
	return _DiaOracleV3.Contract.GetValue(&_DiaOracleV3.CallOpts, key)
}

// GetValueAt is a free data retrieval call binding the contract method 0x09daaa95.
//
// Solidity: function getValueAt(string key, uint256 index) view returns(uint128 value, uint128 timestamp, uint128 volume)
func (_DiaOracleV3 *DiaOracleV3Caller) GetValueAt(opts *bind.CallOpts, key string, index *big.Int) (struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getValueAt", key, index)

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
func (_DiaOracleV3 *DiaOracleV3Session) GetValueAt(key string, index *big.Int) (struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}, error) {
	return _DiaOracleV3.Contract.GetValueAt(&_DiaOracleV3.CallOpts, key, index)
}

// GetValueAt is a free data retrieval call binding the contract method 0x09daaa95.
//
// Solidity: function getValueAt(string key, uint256 index) view returns(uint128 value, uint128 timestamp, uint128 volume)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetValueAt(key string, index *big.Int) (struct {
	Value     *big.Int
	Timestamp *big.Int
	Volume    *big.Int
}, error) {
	return _DiaOracleV3.Contract.GetValueAt(&_DiaOracleV3.CallOpts, key, index)
}

// GetValueCount is a free data retrieval call binding the contract method 0x59c3852c.
//
// Solidity: function getValueCount(string key) view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Caller) GetValueCount(opts *bind.CallOpts, key string) (*big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getValueCount", key)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetValueCount is a free data retrieval call binding the contract method 0x59c3852c.
//
// Solidity: function getValueCount(string key) view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Session) GetValueCount(key string) (*big.Int, error) {
	return _DiaOracleV3.Contract.GetValueCount(&_DiaOracleV3.CallOpts, key)
}

// GetValueCount is a free data retrieval call binding the contract method 0x59c3852c.
//
// Solidity: function getValueCount(string key) view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetValueCount(key string) (*big.Int, error) {
	return _DiaOracleV3.Contract.GetValueCount(&_DiaOracleV3.CallOpts, key)
}

// GetValueHistory is a free data retrieval call binding the contract method 0x8d97ecf2.
//
// Solidity: function getValueHistory(string key) view returns((uint128,uint128,uint128)[])
func (_DiaOracleV3 *DiaOracleV3Caller) GetValueHistory(opts *bind.CallOpts, key string) ([]IDIAOracleV3ValueEntry, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "getValueHistory", key)

	if err != nil {
		return *new([]IDIAOracleV3ValueEntry), err
	}

	out0 := *abi.ConvertType(out[0], new([]IDIAOracleV3ValueEntry)).(*[]IDIAOracleV3ValueEntry)

	return out0, err

}

// GetValueHistory is a free data retrieval call binding the contract method 0x8d97ecf2.
//
// Solidity: function getValueHistory(string key) view returns((uint128,uint128,uint128)[])
func (_DiaOracleV3 *DiaOracleV3Session) GetValueHistory(key string) ([]IDIAOracleV3ValueEntry, error) {
	return _DiaOracleV3.Contract.GetValueHistory(&_DiaOracleV3.CallOpts, key)
}

// GetValueHistory is a free data retrieval call binding the contract method 0x8d97ecf2.
//
// Solidity: function getValueHistory(string key) view returns((uint128,uint128,uint128)[])
func (_DiaOracleV3 *DiaOracleV3CallerSession) GetValueHistory(key string) ([]IDIAOracleV3ValueEntry, error) {
	return _DiaOracleV3.Contract.GetValueHistory(&_DiaOracleV3.CallOpts, key)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DiaOracleV3 *DiaOracleV3Caller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DiaOracleV3 *DiaOracleV3Session) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DiaOracleV3.Contract.HasRole(&_DiaOracleV3.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_DiaOracleV3 *DiaOracleV3CallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _DiaOracleV3.Contract.HasRole(&_DiaOracleV3.CallOpts, role, account)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Caller) ProxiableUUID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "proxiableUUID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3Session) ProxiableUUID() ([32]byte, error) {
	return _DiaOracleV3.Contract.ProxiableUUID(&_DiaOracleV3.CallOpts)
}

// ProxiableUUID is a free data retrieval call binding the contract method 0x52d1902d.
//
// Solidity: function proxiableUUID() view returns(bytes32)
func (_DiaOracleV3 *DiaOracleV3CallerSession) ProxiableUUID() ([32]byte, error) {
	return _DiaOracleV3.Contract.ProxiableUUID(&_DiaOracleV3.CallOpts)
}

// RawData is a free data retrieval call binding the contract method 0x7a2fa442.
//
// Solidity: function rawData(string ) view returns(bytes)
func (_DiaOracleV3 *DiaOracleV3Caller) RawData(opts *bind.CallOpts, arg0 string) ([]byte, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "rawData", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// RawData is a free data retrieval call binding the contract method 0x7a2fa442.
//
// Solidity: function rawData(string ) view returns(bytes)
func (_DiaOracleV3 *DiaOracleV3Session) RawData(arg0 string) ([]byte, error) {
	return _DiaOracleV3.Contract.RawData(&_DiaOracleV3.CallOpts, arg0)
}

// RawData is a free data retrieval call binding the contract method 0x7a2fa442.
//
// Solidity: function rawData(string ) view returns(bytes)
func (_DiaOracleV3 *DiaOracleV3CallerSession) RawData(arg0 string) ([]byte, error) {
	return _DiaOracleV3.Contract.RawData(&_DiaOracleV3.CallOpts, arg0)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DiaOracleV3 *DiaOracleV3Caller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DiaOracleV3 *DiaOracleV3Session) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DiaOracleV3.Contract.SupportsInterface(&_DiaOracleV3.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_DiaOracleV3 *DiaOracleV3CallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _DiaOracleV3.Contract.SupportsInterface(&_DiaOracleV3.CallOpts, interfaceId)
}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Caller) Values(opts *bind.CallOpts, arg0 string) (*big.Int, error) {
	var out []interface{}
	err := _DiaOracleV3.contract.Call(opts, &out, "values", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3Session) Values(arg0 string) (*big.Int, error) {
	return _DiaOracleV3.Contract.Values(&_DiaOracleV3.CallOpts, arg0)
}

// Values is a free data retrieval call binding the contract method 0x5a9ade8b.
//
// Solidity: function values(string ) view returns(uint256)
func (_DiaOracleV3 *DiaOracleV3CallerSession) Values(arg0 string) (*big.Int, error) {
	return _DiaOracleV3.Contract.Values(&_DiaOracleV3.CallOpts, arg0)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DiaOracleV3 *DiaOracleV3Session) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.GrantRole(&_DiaOracleV3.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.GrantRole(&_DiaOracleV3.TransactOpts, role, account)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DiaOracleV3 *DiaOracleV3Session) Initialize() (*types.Transaction, error) {
	return _DiaOracleV3.Contract.Initialize(&_DiaOracleV3.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) Initialize() (*types.Transaction, error) {
	return _DiaOracleV3.Contract.Initialize(&_DiaOracleV3.TransactOpts)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DiaOracleV3 *DiaOracleV3Session) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.RenounceRole(&_DiaOracleV3.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.RenounceRole(&_DiaOracleV3.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DiaOracleV3 *DiaOracleV3Session) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.RevokeRole(&_DiaOracleV3.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.RevokeRole(&_DiaOracleV3.TransactOpts, role, account)
}

// SetDecimals is a paid mutator transaction binding the contract method 0x7a1395aa.
//
// Solidity: function setDecimals(uint8 decimalPrecision) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) SetDecimals(opts *bind.TransactOpts, decimalPrecision uint8) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "setDecimals", decimalPrecision)
}

// SetDecimals is a paid mutator transaction binding the contract method 0x7a1395aa.
//
// Solidity: function setDecimals(uint8 decimalPrecision) returns()
func (_DiaOracleV3 *DiaOracleV3Session) SetDecimals(decimalPrecision uint8) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetDecimals(&_DiaOracleV3.TransactOpts, decimalPrecision)
}

// SetDecimals is a paid mutator transaction binding the contract method 0x7a1395aa.
//
// Solidity: function setDecimals(uint8 decimalPrecision) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) SetDecimals(decimalPrecision uint8) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetDecimals(&_DiaOracleV3.TransactOpts, decimalPrecision)
}

// SetMultipleRawValues is a paid mutator transaction binding the contract method 0x135d90c7.
//
// Solidity: function setMultipleRawValues(bytes[] dataArray) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) SetMultipleRawValues(opts *bind.TransactOpts, dataArray [][]byte) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "setMultipleRawValues", dataArray)
}

// SetMultipleRawValues is a paid mutator transaction binding the contract method 0x135d90c7.
//
// Solidity: function setMultipleRawValues(bytes[] dataArray) returns()
func (_DiaOracleV3 *DiaOracleV3Session) SetMultipleRawValues(dataArray [][]byte) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetMultipleRawValues(&_DiaOracleV3.TransactOpts, dataArray)
}

// SetMultipleRawValues is a paid mutator transaction binding the contract method 0x135d90c7.
//
// Solidity: function setMultipleRawValues(bytes[] dataArray) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) SetMultipleRawValues(dataArray [][]byte) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetMultipleRawValues(&_DiaOracleV3.TransactOpts, dataArray)
}

// SetMultipleValues is a paid mutator transaction binding the contract method 0x8d241526.
//
// Solidity: function setMultipleValues(string[] keys, uint256[] compressedValues) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) SetMultipleValues(opts *bind.TransactOpts, keys []string, compressedValues []*big.Int) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "setMultipleValues", keys, compressedValues)
}

// SetMultipleValues is a paid mutator transaction binding the contract method 0x8d241526.
//
// Solidity: function setMultipleValues(string[] keys, uint256[] compressedValues) returns()
func (_DiaOracleV3 *DiaOracleV3Session) SetMultipleValues(keys []string, compressedValues []*big.Int) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetMultipleValues(&_DiaOracleV3.TransactOpts, keys, compressedValues)
}

// SetMultipleValues is a paid mutator transaction binding the contract method 0x8d241526.
//
// Solidity: function setMultipleValues(string[] keys, uint256[] compressedValues) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) SetMultipleValues(keys []string, compressedValues []*big.Int) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetMultipleValues(&_DiaOracleV3.TransactOpts, keys, compressedValues)
}

// SetRawValue is a paid mutator transaction binding the contract method 0x2c484ae5.
//
// Solidity: function setRawValue(bytes data) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) SetRawValue(opts *bind.TransactOpts, data []byte) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "setRawValue", data)
}

// SetRawValue is a paid mutator transaction binding the contract method 0x2c484ae5.
//
// Solidity: function setRawValue(bytes data) returns()
func (_DiaOracleV3 *DiaOracleV3Session) SetRawValue(data []byte) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetRawValue(&_DiaOracleV3.TransactOpts, data)
}

// SetRawValue is a paid mutator transaction binding the contract method 0x2c484ae5.
//
// Solidity: function setRawValue(bytes data) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) SetRawValue(data []byte) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetRawValue(&_DiaOracleV3.TransactOpts, data)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) SetValue(opts *bind.TransactOpts, key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "setValue", key, value, timestamp)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_DiaOracleV3 *DiaOracleV3Session) SetValue(key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetValue(&_DiaOracleV3.TransactOpts, key, value, timestamp)
}

// SetValue is a paid mutator transaction binding the contract method 0x7898e0c2.
//
// Solidity: function setValue(string key, uint128 value, uint128 timestamp) returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) SetValue(key string, value *big.Int, timestamp *big.Int) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.SetValue(&_DiaOracleV3.TransactOpts, key, value, timestamp)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DiaOracleV3 *DiaOracleV3Transactor) UpgradeToAndCall(opts *bind.TransactOpts, newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DiaOracleV3.contract.Transact(opts, "upgradeToAndCall", newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DiaOracleV3 *DiaOracleV3Session) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.UpgradeToAndCall(&_DiaOracleV3.TransactOpts, newImplementation, data)
}

// UpgradeToAndCall is a paid mutator transaction binding the contract method 0x4f1ef286.
//
// Solidity: function upgradeToAndCall(address newImplementation, bytes data) payable returns()
func (_DiaOracleV3 *DiaOracleV3TransactorSession) UpgradeToAndCall(newImplementation common.Address, data []byte) (*types.Transaction, error) {
	return _DiaOracleV3.Contract.UpgradeToAndCall(&_DiaOracleV3.TransactOpts, newImplementation, data)
}

// DiaOracleV3DecimalsUpdateIterator is returned from FilterDecimalsUpdate and is used to iterate over the raw logs and unpacked data for DecimalsUpdate events raised by the DiaOracleV3 contract.
type DiaOracleV3DecimalsUpdateIterator struct {
	Event *DiaOracleV3DecimalsUpdate // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3DecimalsUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3DecimalsUpdate)
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
		it.Event = new(DiaOracleV3DecimalsUpdate)
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
func (it *DiaOracleV3DecimalsUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3DecimalsUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3DecimalsUpdate represents a DecimalsUpdate event raised by the DiaOracleV3 contract.
type DiaOracleV3DecimalsUpdate struct {
	Decimals uint8
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterDecimalsUpdate is a free log retrieval operation binding the contract event 0xc03a161db992c7f007555d2c71b8de00760358129d9547bce9475bff4739bf60.
//
// Solidity: event DecimalsUpdate(uint8 decimals)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterDecimalsUpdate(opts *bind.FilterOpts) (*DiaOracleV3DecimalsUpdateIterator, error) {

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "DecimalsUpdate")
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3DecimalsUpdateIterator{contract: _DiaOracleV3.contract, event: "DecimalsUpdate", logs: logs, sub: sub}, nil
}

// WatchDecimalsUpdate is a free log subscription operation binding the contract event 0xc03a161db992c7f007555d2c71b8de00760358129d9547bce9475bff4739bf60.
//
// Solidity: event DecimalsUpdate(uint8 decimals)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchDecimalsUpdate(opts *bind.WatchOpts, sink chan<- *DiaOracleV3DecimalsUpdate) (event.Subscription, error) {

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "DecimalsUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3DecimalsUpdate)
				if err := _DiaOracleV3.contract.UnpackLog(event, "DecimalsUpdate", log); err != nil {
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

// ParseDecimalsUpdate is a log parse operation binding the contract event 0xc03a161db992c7f007555d2c71b8de00760358129d9547bce9475bff4739bf60.
//
// Solidity: event DecimalsUpdate(uint8 decimals)
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseDecimalsUpdate(log types.Log) (*DiaOracleV3DecimalsUpdate, error) {
	event := new(DiaOracleV3DecimalsUpdate)
	if err := _DiaOracleV3.contract.UnpackLog(event, "DecimalsUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3InitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the DiaOracleV3 contract.
type DiaOracleV3InitializedIterator struct {
	Event *DiaOracleV3Initialized // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3InitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3Initialized)
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
		it.Event = new(DiaOracleV3Initialized)
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
func (it *DiaOracleV3InitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3InitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3Initialized represents a Initialized event raised by the DiaOracleV3 contract.
type DiaOracleV3Initialized struct {
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterInitialized(opts *bind.FilterOpts) (*DiaOracleV3InitializedIterator, error) {

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3InitializedIterator{contract: _DiaOracleV3.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *DiaOracleV3Initialized) (event.Subscription, error) {

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3Initialized)
				if err := _DiaOracleV3.contract.UnpackLog(event, "Initialized", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseInitialized(log types.Log) (*DiaOracleV3Initialized, error) {
	event := new(DiaOracleV3Initialized)
	if err := _DiaOracleV3.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3OracleUpdateIterator is returned from FilterOracleUpdate and is used to iterate over the raw logs and unpacked data for OracleUpdate events raised by the DiaOracleV3 contract.
type DiaOracleV3OracleUpdateIterator struct {
	Event *DiaOracleV3OracleUpdate // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3OracleUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3OracleUpdate)
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
		it.Event = new(DiaOracleV3OracleUpdate)
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
func (it *DiaOracleV3OracleUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3OracleUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3OracleUpdate represents a OracleUpdate event raised by the DiaOracleV3 contract.
type DiaOracleV3OracleUpdate struct {
	Key       string
	Value     *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterOracleUpdate is a free log retrieval operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterOracleUpdate(opts *bind.FilterOpts) (*DiaOracleV3OracleUpdateIterator, error) {

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "OracleUpdate")
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3OracleUpdateIterator{contract: _DiaOracleV3.contract, event: "OracleUpdate", logs: logs, sub: sub}, nil
}

// WatchOracleUpdate is a free log subscription operation binding the contract event 0xa7fc99ed7617309ee23f63ae90196a1e490d362e6f6a547a59bc809ee2291782.
//
// Solidity: event OracleUpdate(string key, uint128 value, uint128 timestamp)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchOracleUpdate(opts *bind.WatchOpts, sink chan<- *DiaOracleV3OracleUpdate) (event.Subscription, error) {

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "OracleUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3OracleUpdate)
				if err := _DiaOracleV3.contract.UnpackLog(event, "OracleUpdate", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseOracleUpdate(log types.Log) (*DiaOracleV3OracleUpdate, error) {
	event := new(DiaOracleV3OracleUpdate)
	if err := _DiaOracleV3.contract.UnpackLog(event, "OracleUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3OracleUpdateRawIterator is returned from FilterOracleUpdateRaw and is used to iterate over the raw logs and unpacked data for OracleUpdateRaw events raised by the DiaOracleV3 contract.
type DiaOracleV3OracleUpdateRawIterator struct {
	Event *DiaOracleV3OracleUpdateRaw // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3OracleUpdateRawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3OracleUpdateRaw)
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
		it.Event = new(DiaOracleV3OracleUpdateRaw)
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
func (it *DiaOracleV3OracleUpdateRawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3OracleUpdateRawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3OracleUpdateRaw represents a OracleUpdateRaw event raised by the DiaOracleV3 contract.
type DiaOracleV3OracleUpdateRaw struct {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterOracleUpdateRaw(opts *bind.FilterOpts) (*DiaOracleV3OracleUpdateRawIterator, error) {

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "OracleUpdateRaw")
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3OracleUpdateRawIterator{contract: _DiaOracleV3.contract, event: "OracleUpdateRaw", logs: logs, sub: sub}, nil
}

// WatchOracleUpdateRaw is a free log subscription operation binding the contract event 0x0ec1e0298284e066eddd5e448f165c9337bf3f9447b7159177c72e0cada227d3.
//
// Solidity: event OracleUpdateRaw(string key, uint128 value, uint128 timestamp, uint128 volume, bytes data)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchOracleUpdateRaw(opts *bind.WatchOpts, sink chan<- *DiaOracleV3OracleUpdateRaw) (event.Subscription, error) {

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "OracleUpdateRaw")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3OracleUpdateRaw)
				if err := _DiaOracleV3.contract.UnpackLog(event, "OracleUpdateRaw", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseOracleUpdateRaw(log types.Log) (*DiaOracleV3OracleUpdateRaw, error) {
	event := new(DiaOracleV3OracleUpdateRaw)
	if err := _DiaOracleV3.contract.UnpackLog(event, "OracleUpdateRaw", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3RoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the DiaOracleV3 contract.
type DiaOracleV3RoleAdminChangedIterator struct {
	Event *DiaOracleV3RoleAdminChanged // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3RoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3RoleAdminChanged)
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
		it.Event = new(DiaOracleV3RoleAdminChanged)
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
func (it *DiaOracleV3RoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3RoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3RoleAdminChanged represents a RoleAdminChanged event raised by the DiaOracleV3 contract.
type DiaOracleV3RoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*DiaOracleV3RoleAdminChangedIterator, error) {

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

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3RoleAdminChangedIterator{contract: _DiaOracleV3.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *DiaOracleV3RoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

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

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3RoleAdminChanged)
				if err := _DiaOracleV3.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseRoleAdminChanged(log types.Log) (*DiaOracleV3RoleAdminChanged, error) {
	event := new(DiaOracleV3RoleAdminChanged)
	if err := _DiaOracleV3.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3RoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the DiaOracleV3 contract.
type DiaOracleV3RoleGrantedIterator struct {
	Event *DiaOracleV3RoleGranted // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3RoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3RoleGranted)
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
		it.Event = new(DiaOracleV3RoleGranted)
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
func (it *DiaOracleV3RoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3RoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3RoleGranted represents a RoleGranted event raised by the DiaOracleV3 contract.
type DiaOracleV3RoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DiaOracleV3RoleGrantedIterator, error) {

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

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3RoleGrantedIterator{contract: _DiaOracleV3.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *DiaOracleV3RoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3RoleGranted)
				if err := _DiaOracleV3.contract.UnpackLog(event, "RoleGranted", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseRoleGranted(log types.Log) (*DiaOracleV3RoleGranted, error) {
	event := new(DiaOracleV3RoleGranted)
	if err := _DiaOracleV3.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3RoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the DiaOracleV3 contract.
type DiaOracleV3RoleRevokedIterator struct {
	Event *DiaOracleV3RoleRevoked // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3RoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3RoleRevoked)
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
		it.Event = new(DiaOracleV3RoleRevoked)
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
func (it *DiaOracleV3RoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3RoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3RoleRevoked represents a RoleRevoked event raised by the DiaOracleV3 contract.
type DiaOracleV3RoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*DiaOracleV3RoleRevokedIterator, error) {

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

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3RoleRevokedIterator{contract: _DiaOracleV3.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *DiaOracleV3RoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3RoleRevoked)
				if err := _DiaOracleV3.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseRoleRevoked(log types.Log) (*DiaOracleV3RoleRevoked, error) {
	event := new(DiaOracleV3RoleRevoked)
	if err := _DiaOracleV3.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3UpdaterAddressChangeIterator is returned from FilterUpdaterAddressChange and is used to iterate over the raw logs and unpacked data for UpdaterAddressChange events raised by the DiaOracleV3 contract.
type DiaOracleV3UpdaterAddressChangeIterator struct {
	Event *DiaOracleV3UpdaterAddressChange // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3UpdaterAddressChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3UpdaterAddressChange)
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
		it.Event = new(DiaOracleV3UpdaterAddressChange)
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
func (it *DiaOracleV3UpdaterAddressChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3UpdaterAddressChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3UpdaterAddressChange represents a UpdaterAddressChange event raised by the DiaOracleV3 contract.
type DiaOracleV3UpdaterAddressChange struct {
	NewUpdater common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterUpdaterAddressChange is a free log retrieval operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterUpdaterAddressChange(opts *bind.FilterOpts) (*DiaOracleV3UpdaterAddressChangeIterator, error) {

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "UpdaterAddressChange")
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3UpdaterAddressChangeIterator{contract: _DiaOracleV3.contract, event: "UpdaterAddressChange", logs: logs, sub: sub}, nil
}

// WatchUpdaterAddressChange is a free log subscription operation binding the contract event 0x121e958a4cadf7f8dadefa22cc019700365240223668418faebed197da07089f.
//
// Solidity: event UpdaterAddressChange(address newUpdater)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchUpdaterAddressChange(opts *bind.WatchOpts, sink chan<- *DiaOracleV3UpdaterAddressChange) (event.Subscription, error) {

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "UpdaterAddressChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3UpdaterAddressChange)
				if err := _DiaOracleV3.contract.UnpackLog(event, "UpdaterAddressChange", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseUpdaterAddressChange(log types.Log) (*DiaOracleV3UpdaterAddressChange, error) {
	event := new(DiaOracleV3UpdaterAddressChange)
	if err := _DiaOracleV3.contract.UnpackLog(event, "UpdaterAddressChange", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// DiaOracleV3UpgradedIterator is returned from FilterUpgraded and is used to iterate over the raw logs and unpacked data for Upgraded events raised by the DiaOracleV3 contract.
type DiaOracleV3UpgradedIterator struct {
	Event *DiaOracleV3Upgraded // Event containing the contract specifics and raw log

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
func (it *DiaOracleV3UpgradedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DiaOracleV3Upgraded)
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
		it.Event = new(DiaOracleV3Upgraded)
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
func (it *DiaOracleV3UpgradedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DiaOracleV3UpgradedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DiaOracleV3Upgraded represents a Upgraded event raised by the DiaOracleV3 contract.
type DiaOracleV3Upgraded struct {
	Implementation common.Address
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterUpgraded is a free log retrieval operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DiaOracleV3 *DiaOracleV3Filterer) FilterUpgraded(opts *bind.FilterOpts, implementation []common.Address) (*DiaOracleV3UpgradedIterator, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DiaOracleV3.contract.FilterLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return &DiaOracleV3UpgradedIterator{contract: _DiaOracleV3.contract, event: "Upgraded", logs: logs, sub: sub}, nil
}

// WatchUpgraded is a free log subscription operation binding the contract event 0xbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b.
//
// Solidity: event Upgraded(address indexed implementation)
func (_DiaOracleV3 *DiaOracleV3Filterer) WatchUpgraded(opts *bind.WatchOpts, sink chan<- *DiaOracleV3Upgraded, implementation []common.Address) (event.Subscription, error) {

	var implementationRule []interface{}
	for _, implementationItem := range implementation {
		implementationRule = append(implementationRule, implementationItem)
	}

	logs, sub, err := _DiaOracleV3.contract.WatchLogs(opts, "Upgraded", implementationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DiaOracleV3Upgraded)
				if err := _DiaOracleV3.contract.UnpackLog(event, "Upgraded", log); err != nil {
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
func (_DiaOracleV3 *DiaOracleV3Filterer) ParseUpgraded(log types.Log) (*DiaOracleV3Upgraded, error) {
	event := new(DiaOracleV3Upgraded)
	if err := _DiaOracleV3.contract.UnpackLog(event, "Upgraded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
