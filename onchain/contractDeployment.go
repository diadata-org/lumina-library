package onchain

import (
	"strings"
	"time"

	diaOracleV3 "github.com/diadata-org/lumina-library/contracts/lumina/diaoraclev3"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

func DeployOrBindContract(
	deployedContract string,
	conn *ethclient.Client,
	connBackup *ethclient.Client,
	auth *bind.TransactOpts,
	contract **diaOracleV3.DIAOracleV3,
	contractBackup **diaOracleV3.DIAOracleV3,
	decimalPrecision uint8,
) error {
	var err error
	if deployedContract != "" {

		// bind primary and backup to existing proxy address
		*contract, err = diaOracleV3.NewDIAOracleV3(common.HexToAddress(deployedContract), conn)
		if err != nil {
			return err
		}
		*contractBackup, err = diaOracleV3.NewDIAOracleV3(common.HexToAddress(deployedContract), connBackup)
		if err != nil {
			return err
		}

	} else {
		// deploy implementation
		var implAddr common.Address
		var implTx *types.Transaction
		implAddr, implTx, _, err = diaOracleV3.DeployDIAOracleV3(auth, conn)
		if err != nil {
			log.Fatalf("could not deploy contract implementation: %v", err)
			return err
		}
		log.Infof("Implementation pending deploy: 0x%x.", implAddr)
		log.Infof("Implementation Transaction waiting to be mined: 0x%x.", implTx.Hash())
		time.Sleep(30 * time.Second)

		// pack initialize(uint8) calldata to pass into proxy constructor —
		// this atomically initializes the contract in the same tx as proxy deployment,
		// eliminating the front-running window that would exist if initialize() were
		// called in a separate transaction.
		diaOracleV3ABI, err := diaOracleV3.DIAOracleV3MetaData.GetAbi()
		if err != nil {
			log.Fatalf("could not parse DIAOracleV3 ABI: %v", err)
			return err
		}

		initData, err := diaOracleV3ABI.Pack("initialize", decimalPrecision)
		if err != nil {
			log.Fatalf("could not pack initialize data: %v", err)
			return err
		}

		log.Info("Initialization data packed successfully")

		// Get ERC1967Proxy ABI
		proxyABI, err := abi.JSON(strings.NewReader("[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_logic\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"_data\",\"type\":\"bytes\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"previousAdmin\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"newAdmin\",\"type\":\"address\"}],\"name\":\"AdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"beacon\",\"type\":\"address\"}],\"name\":\"BeaconUpgraded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"implementation\",\"type\":\"address\"}],\"name\":\"Upgraded\",\"type\":\"event\"},{\"stateMutability\":\"nonpayable\",\"type\":\"fallback\"}]"))
		if err != nil {
			log.Fatalf("could not parse proxy ABI: %v", err)
			return err
		}

		// ERC1967Proxy bytecode from OpenZeppelin v5.5.0
		proxyBytecode := common.FromHex("0x608060405261028c8038038061001481610158565b9283398101604082820312610140578151916001600160a01b03831690818403610140576020810151906001600160401b038211610140570182601f82011215610140578051906001600160401b0382116101445761007c601f8301601f1916602001610158565b938285526020838301011161014057815f9260208093018387015e84010152823b1561012e577f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc80546001600160a01b031916821790557fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b5f80a2805115610117576101079161017d565b505b6040516082908161020a8239f35b505034156101095763b398979f60e01b5f5260045ffd5b634c9c8ce360e01b5f5260045260245ffd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6040519190601f01601f191682016001600160401b0381118382101761014457604052565b905f8091602081519101845af480806101f6575b156101b15750506040513d81523d5f602083013e60203d82010160405290565b156101d657639996b31560e01b5f9081526001600160a01b0391909116600452602490fd5b3d156101e7576040513d5f823e3d90fd5b63d6bda27560e01b5f5260045ffd5b503d1515806101915750813b151561019156fe60806040527f360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc545f9081906001600160a01b0316368280378136915af43d5f803e156048573d5ff35b3d5ffdfea264697066735822122014e90b4a0575eac550ea455374a3550d1857b2a03810b58b0a1ba6c7cc0039cb64736f6c63430008220033")

		// deploy ERC1967 proxy pointing to implementation, atomically initializing it
		proxyAddr, proxyTx, _, err := bind.DeployContract(
			auth,
			proxyABI,
			proxyBytecode,
			conn,
			implAddr,
			initData,
		)
		if err != nil {
			log.Fatalf("could not deploy proxy: %v", err)
			return err
		}

		log.Infof("ERC1967Proxy pending deploy: 0x%x.", proxyAddr)
		log.Infof("Proxy Transaction waiting to be mined: 0x%x.", proxyTx.Hash())
		time.Sleep(30 * time.Second)

		// bind primary and backup to proxy
		*contract, err = diaOracleV3.NewDIAOracleV3(proxyAddr, conn)
		if err != nil {
			log.Fatalf("could not bind to proxy: %v", err)
			return err
		}
		*contractBackup, err = diaOracleV3.NewDIAOracleV3(proxyAddr, connBackup)
		if err != nil {
			log.Fatalf("could not bind backup to proxy: %v", err)
			return err
		}

		log.Info("Deployment successful!")
		log.Infof("DIAOracleV3 Implementation: 0x%x", implAddr)
		log.Infof("ERC1967Proxy (use this address): 0x%x", proxyAddr)
		log.Infof("Decimals: %d", decimalPrecision)
	}
	return nil
}
