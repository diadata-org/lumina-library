package onchain

import (
	"context"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"strings"
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
	conn *ethclient.Client,
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
			log.Warnf("updater - Failed to update Oracle: %v.", err)
			return
		}
	}

}

func OracleUpdateExecutor(
	// publishedPrices map[string]float64,
	// newPrices map[string]float64,
	// deviationPermille int,
	auth *bind.TransactOpts,
	contract *diaOracleV2MultiupdateService.DiaOracleV2MultiupdateService,
	conn *ethclient.Client,
	chainId int64,
	// compatibilityMode bool,
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
		log.Infof("updater - Attempting to update oracle with %d values to endpoint %s (backup: %s)",
			len(values), os.Getenv("BLOCKCHAIN_NODE"), os.Getenv("BACKUP_NODE"))

		err := updateOracleMultiValues(conn, contract, auth, chainId, keys, values, timestamp)
		if err != nil {
			log.Warnf("updater - Failed to update Oracle: %v. ChainID: %v",
				err, chainId)

			// Try to analyze the error for more details
			if strings.Contains(err.Error(), "404 Not Found") {
				log.Warnf("updater - Connection error details: Primary node: %s, Backup node: %s",
					os.Getenv("BLOCKCHAIN_NODE"), os.Getenv("BACKUP_NODE"))
			}
			return
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

	// Add additional context for debugging
	log.Infof("updater - Attempting updateOracleMultiValues with chainId: %d, keys count: %d", chainId, len(keys))

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
			log.Errorf("updater - SuggestGasPrice failed: %v. This often indicates RPC connection issues.", err)
			// Try to further diagnose the connection error
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_, err := client.HeaderByNumber(ctx, nil)
			if err != nil {
				log.Errorf("updater - RPC connection test failed: %v", err)
			}
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
	// check if tx is sendable then fgo backup
	if err != nil {
		// backup in here
		return err
	}

	log.Infof("updater - Gas price: %d.", tx.GasPrice())
	// log.Printf("Data: %x\n", tx.Data())
	log.Infof("updater - Nonce: %d.", tx.Nonce())
	log.Infof("updater - Tx To: %s.", tx.To().String())
	log.Infof("updater - Tx Hash: 0x%x.", tx.Hash())
	return nil
}
