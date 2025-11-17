package utils

import (
	"crypto/ecdsa"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func SetupOnchain(blockchainNode string, backupNode string, metacontractNode string, privateKeyHex string, chainID int64) (
	conn *ethclient.Client,
	connBackup *ethclient.Client,
	connMetacontract *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	auth *bind.TransactOpts,
) {
	conn, err := MakeEthClient(blockchainNode, backupNode)
	if err != nil {
		log.Fatalf("MakeEthClient for main connection: %v", err)
	}
	connBackup, err = MakeEthClient(backupNode, blockchainNode)
	if err != nil {
		log.Fatalf("MakeEthClient for backup connection: %v", err)
	}
	connMetacontract, err = MakeEthClient(metacontractNode, metacontractNode)
	if err != nil {
		log.Fatalf("MakeEthClient for metacontract connection: %v", err)
	}

	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKey, err = crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}

	auth, err = bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(chainID))
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	return
}
