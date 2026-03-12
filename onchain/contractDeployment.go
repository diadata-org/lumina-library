package onchain

import (
	"strings"
	"time"

	diaOracleV3MultiupdateService "github.com/diadata-org/lumina-library/contracts/lumina/diaoraclev3"
	diaERC1967Proxy "github.com/diadata-org/lumina-library/contracts/lumina/erc1967proxy"
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
	contract **diaOracleV3MultiupdateService.DiaOracleV3MultiupdateService,
	contractBackup **diaOracleV3MultiupdateService.DiaOracleV3MultiupdateService,
) error {
	var err error
	if deployedContract != "" {

		// bind primary and backup to existing proxy address
		*contract, err = diaOracleV3MultiupdateService.NewDiaOracleV3MultiupdateService(common.HexToAddress(deployedContract), conn)
		if err != nil {
			return err
		}
		*contractBackup, err = diaOracleV3MultiupdateService.NewDiaOracleV3MultiupdateService(common.HexToAddress(deployedContract), connBackup)
		if err != nil {
			return err
		}

	} else {
		// deploy implementation
		var implAddr common.Address
		var implTx *types.Transaction
		implAddr, implTx, _, err = diaOracleV3MultiupdateService.DeployDiaOracleV3MultiupdateService(auth, conn)
		if err != nil {
			log.Fatalf("could not deploy contract implementation: %v", err)
			return err
		}
		log.Infof("Implementation pending deploy: 0x%x.", implAddr)
		log.Infof("Transaction waiting to be mined: 0x%x.", implTx.Hash())
		time.Sleep(180000 * time.Millisecond)

		// pack initialize() calldata to pass into proxy constructor —
		// this atomically initializes the contract in the same tx as proxy deployment,
		// eliminating the front-running window that would exist if initialize() were
		// called in a separate transaction.
		parsedABI, err := abi.JSON(strings.NewReader(diaOracleV3MultiupdateService.DiaOracleV3MultiupdateServiceMetaData.ABI))
		if err != nil {
			return err
		}
		initData, err := parsedABI.Pack("initialize")
		if err != nil {
			return err
		}

		// deploy ERC1967 proxy pointing to implementation, atomically initializing it
		var proxyAddr common.Address
		var proxyTx *types.Transaction
		proxyAddr, proxyTx, _, err = diaERC1967Proxy.DeployDiaERC1967Proxy(auth, conn, implAddr, initData)
		if err != nil {
			log.Fatalf("could not deploy proxy: %v", err)
			return err
		}
		log.Infof("Proxy pending deploy: 0x%x.", proxyAddr)
		log.Infof("Transaction waiting to be mined: 0x%x.", proxyTx.Hash())
		time.Sleep(180000 * time.Millisecond)

		// bind primary and backup to proxy
		*contract, err = diaOracleV3MultiupdateService.NewDiaOracleV3MultiupdateService(proxyAddr, conn)
		if err != nil {
			return err
		}
		*contractBackup, err = diaOracleV3MultiupdateService.NewDiaOracleV3MultiupdateService(proxyAddr, connBackup)
		if err != nil {
			return err
		}
		log.Infof("Contract deployed and initialized at proxy: 0x%x.", proxyAddr)
	}
	return nil
}
