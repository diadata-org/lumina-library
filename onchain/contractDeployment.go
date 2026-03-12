package onchain

import (
	"time"

	diaOracleV3 "github.com/diadata-org/lumina-library/contracts/lumina/diaoraclev3"
	diaERC1967Proxy "github.com/diadata-org/lumina-library/contracts/lumina/erc1967proxy"
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
	contract **diaOracleV3.DiaOracleV3,
	contractBackup **diaOracleV3.DiaOracleV3,
) error {
	var err error
	if deployedContract != "" {

		// bind primary and backup to existing proxy address
		*contract, err = diaOracleV3.NewDiaOracleV3(common.HexToAddress(deployedContract), conn)
		if err != nil {
			return err
		}
		*contractBackup, err = diaOracleV3.NewDiaOracleV3(common.HexToAddress(deployedContract), connBackup)
		if err != nil {
			return err
		}

	} else {
		// deploy implementation
		var implAddr common.Address
		var implTx *types.Transaction
		implAddr, implTx, _, err = diaOracleV3.DeployDiaOracleV3(auth, conn)
		if err != nil {
			log.Fatalf("could not deploy contract implementation: %v", err)
			return err
		}
		log.Infof("Implementation pending deploy: 0x%x.", implAddr)
		log.Infof("Transaction waiting to be mined: 0x%x.", implTx.Hash())
		time.Sleep(180000 * time.Millisecond)

		// deploy ERC1967 proxy pointing to implementation
		var proxyAddr common.Address
		var proxyTx *types.Transaction
		proxyAddr, proxyTx, _, err = diaERC1967Proxy.DeployDiaERC1967Proxy(auth, conn, implAddr, []byte{})
		if err != nil {
			log.Fatalf("could not deploy proxy: %v", err)
			return err
		}
		log.Infof("Proxy pending deploy: 0x%x.", proxyAddr)
		log.Infof("Transaction waiting to be mined: 0x%x.", proxyTx.Hash())
		time.Sleep(180000 * time.Millisecond)

		// bind to proxy and initialize
		*contract, err = diaOracleV3.NewDiaOracleV3(proxyAddr, conn)
		if err != nil {
			return err
		}
		_, err = (*contract).Initialize(auth)
		if err != nil {
			log.Fatalf("could not initialize contract: %v", err)
			return err
		}
		log.Infof("Contract initialized at proxy: 0x%x.", proxyAddr)
		time.Sleep(180000 * time.Millisecond)

		// bind backup to proxy
		*contractBackup, err = diaOracleV3.NewDiaOracleV3(proxyAddr, connBackup)
		if err != nil {
			return err
		}
	}
	return nil
}
