package onchain

import (
	"context"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"time"

	diaOracleV2MultiupdateService "github.com/diadata-org/diadata/pkg/dia/scraper/blockchain-scrapers/blockchains/ethereum/diaOracleV2MultiupdateService"
	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

const (
	DECIMALS_ORACLE_VALUE = 8
)

var (
	log *logrus.Logger
)

func init() {
	log = logrus.New()
	loglevel, err := logrus.ParseLevel(utils.Getenv("LOG_LEVEL_UPDATER", "info"))
	if err != nil {
		log.Errorf("Parse log level: %v.", err)
	}
	log.SetLevel(loglevel)
}

func OracleUpdateExecutorSimulation(
	auth *bind.TransactOpts,
	contract *diaOracleV2MultiupdateService.DiaOracleV2MultiupdateService,
	contractBackup *diaOracleV2MultiupdateService.DiaOracleV2MultiupdateService,
	conn *ethclient.Client,
	connBackup *ethclient.Client,
	chainId int64,
	filtersChannel <-chan []models.FilterPointPair,
) {

	for filterPoints := range filtersChannel {
		timestamp := time.Now().Unix()
		var keys []string
		var values []int64
		for _, fp := range filterPoints {
			log.Infof(
				"updater - filterPoint received at %v: %v -- %v -- %v.",
				time.Unix(timestamp, 0),
				fp.Pair.QuoteToken.Symbol,
				fp.Value,
				fp.Time,
			)

			key := models.GetOracleKeySimulation(fp.Pair)
			keys = append(keys, key)
			values = append(values, int64(fp.Value*math.Pow10(int(DECIMALS_ORACLE_VALUE))))
		}
		err := updateOracleMultiValues(conn, contract, auth, chainId, keys, values, timestamp)
		if err != nil {
			log.Warnf("updater - Failed to update Oracle: %v. Retry with backup node.", err)
			err := updateOracleMultiValues(connBackup, contractBackup, auth, chainId, keys, values, timestamp)
			if err != nil {
				log.Errorf("backup updater - Failed to update Oracle: %v.", err)
				return
			}
		}
	}
}

func OracleUpdateExecutor(
	auth *bind.TransactOpts,
	contract *diaOracleV2MultiupdateService.DiaOracleV2MultiupdateService,
	contractBackup *diaOracleV2MultiupdateService.DiaOracleV2MultiupdateService,
	conn *ethclient.Client,
	connBackup *ethclient.Client,
	chainId int64,
	filtersChannel <-chan []models.FilterPointPair,
) {

	for filterPoints := range filtersChannel {
		timestamp := time.Now().Unix()
		var keys []string
		var values []int64
		for _, fp := range filterPoints {
			log.Infof(
				"updater - filterPoint received at %v: %v -- %v -- %v.",
				time.Unix(timestamp, 0),
				fp.Pair.QuoteToken.Symbol,
				fp.Value,
				fp.Time,
			)

			key := models.GetOracleKey(fp.SourceType, fp.Pair)
			keys = append(keys, key)
			// keys = append(keys, fp.Pair.QuoteToken.Symbol+"/USD")
			values = append(values, int64(fp.Value*math.Pow10(int(DECIMALS_ORACLE_VALUE))))
		}
		err := updateOracleMultiValues(conn, contract, auth, chainId, keys, values, timestamp)
		if err != nil {
			log.Warnf("updater - Failed to update Oracle: %v. Retry with backup node.", err)
			err := updateOracleMultiValues(connBackup, contractBackup, auth, chainId, keys, values, timestamp)
			if err != nil {
				log.Errorf("backup updater - Failed to update Oracle: %v.", err)
				return
			}
		}
	}
}

func updateOracleMultiValues(
	client *ethclient.Client,
	contract *diaOracleV2MultiupdateService.DiaOracleV2MultiupdateService,
	auth *bind.TransactOpts,
	chainId int64,
	keys []string,
	values []int64,
	timestamp int64) error {

	var cValues []*big.Int
	var gasPrice *big.Int
	var err error

	// Get proper gas price depending on chainId
	switch chainId {
	/*case 288: //Bobapkg/scraper/Uniswapv2.go
	gasPrice = big.NewInt(1000000000)*/
	case 592: //Astar
		response, err := http.Get("https://gas.astar.network/api/gasnow?network=astar")
		if err != nil {
			return err
		}

		defer response.Body.Close()
		if 200 != response.StatusCode {
			return err
		}
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return err
		}

		gasSuggestion := gjson.Get(string(contents), "data.fast")
		gasPrice = big.NewInt(gasSuggestion.Int())
	default:
		// Get gas price suggestion
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Errorf("updater - SuggestGasPrice: %v.", err)
			return err
		}

		// Get 110% of the gas price
		fGas := new(big.Float).SetInt(gasPrice)
		fGas.Mul(fGas, big.NewFloat(1.1))
		gasPrice, _ = fGas.Int(nil)
	}

	for _, value := range values {
		// Create compressed argument with values/timestamps
		cValue := big.NewInt(value)
		cValue = cValue.Lsh(cValue, 128)
		cValue = cValue.Add(cValue, big.NewInt(timestamp))
		cValues = append(cValues, cValue)
	}

	// Write values to smart contract
	tx, err := contract.SetMultipleValues(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: gasPrice,
	}, keys, cValues)
	if err != nil {
		return err
	}

	log.Infof("updater - Gas price: %d.", tx.GasPrice())
	// log.Printf("Data: %x\n", tx.Data())
	log.Infof("updater - Nonce: %d.", tx.Nonce())
	log.Infof("updater - Tx To: %s.", tx.To().String())
	log.Infof("updater - Tx Hash: 0x%x.", tx.Hash())
	return nil
}
