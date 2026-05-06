package onchain

import (
	"context"
	"io"
	"math/big"
	"net/http"
	"time"

	diaOracleV3 "github.com/diadata-org/lumina-library/contracts/lumina/diaoraclev3"
	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

var (
	log                 *logrus.Logger
	firstUpdateComplete bool = false
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
	contract *diaOracleV3.DIAOracleV3,
	contractBackup *diaOracleV3.DIAOracleV3,
	conn *ethclient.Client,
	connBackup *ethclient.Client,
	chainId int64,
	decimals int,
	filtersChannel <-chan []models.FilterPointPair,
) {

	for filterPoints := range filtersChannel {
		timestamp := time.Now().Unix()
		var keys []string
		var values []*big.Int
		for _, fp := range filterPoints {
			log.Infof(
				"updater - filterPoint received at %v: %v -- %v -- %v.",
				time.Unix(timestamp, 0),
				fp.Pair.QuoteToken.Symbol,
				fp.Value,
				fp.Time,
			)
			log.Infof("updater -- filterPoint received at unix timestamp (now) %v vs fp.Time %v", timestamp, fp.Time.Unix())
			key := fp.Pair.GetOracleKey(fp.SourceType)
			if key == "" {
				log.Warnf("updater - skipping filter point with empty oracle key for asset %s", fp.Pair.QuoteToken.Symbol)
				continue
			}
			keys = append(keys, key)
			values = append(values, utils.ScaleFloat(fp.Value, decimals))
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
	contract *diaOracleV3.DIAOracleV3,
	contractBackup *diaOracleV3.DIAOracleV3,
	conn *ethclient.Client,
	connBackup *ethclient.Client,
	chainId int64,
	decimals int,
	batchSize int,
	filtersChannel <-chan []models.FilterPointPair,
) {

	for filterPoints := range filtersChannel {
		keysMap := make(map[string]struct{})

		timestamp := time.Now().Unix()
		var keys []string
		var values []*big.Int
		for _, fp := range filterPoints {
			log.Infof(
				"updater - filterPoint received at %v: %s: %s-%s -- %v -- %v.",
				time.Unix(timestamp, 0),
				fp.Pair.QuoteToken.Symbol,
				fp.Pair.QuoteToken.Blockchain,
				fp.Pair.QuoteToken.Address,
				fp.Value,
				fp.Time,
			)
			key := fp.Pair.GetOracleKey(fp.SourceType)
			if key == "" {
				log.Warnf("updater - skipping filter point with empty oracle key for asset %s", fp.Pair.QuoteToken.Symbol)
				continue
			}

			// TO DO: amend this check once we switch to blockchain-address identifier!!
			if _, ok := keysMap[key]; !ok {
				keys = append(keys, key)
				values = append(values, utils.ScaleFloat(fp.Value, decimals))
				keysMap[key] = struct{}{}
			} else {
				log.Warnf("symbol %s already existing.", key)
			}

		}

		if !firstUpdateComplete {
			numBatches := (len(keys) + batchSize - 1) / batchSize
			log.Infof("First oracle update run - enabling batch mode (%d key-value pairs in %d batches)", len(keys), numBatches)

			allSuccessful := true
			for i := 0; i < numBatches; i++ {
				start := i * batchSize
				end := start + batchSize
				if end > len(keys) {
					end = len(keys)
				}

				batchKeys := keys[start:end]
				batchValues := values[start:end]

				log.Infof("Processing batch %d of %d (%d pairs)", i+1, numBatches, len(batchKeys))

				err := updateOracleMultiValues(conn, contract, auth, chainId, batchKeys, batchValues, timestamp)
				if err != nil {
					log.Warnf("updater - Failed to update Oracle batch %d: %v. Retry with backup node.", i+1, err)
					err := updateOracleMultiValues(connBackup, contractBackup, auth, chainId, batchKeys, batchValues, timestamp)
					if err != nil {
						log.Errorf("backup updater - Failed to update Oracle batch %d: %v.", i+1, err)
						allSuccessful = false
						break
					}
				}
			}

			if allSuccessful {
				firstUpdateComplete = true
				log.Info("First oracle update run completed successfully")
			}
		} else {
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
}

func updateOracleMultiValues(
	client *ethclient.Client,
	contract *diaOracleV3.DIAOracleV3,
	auth *bind.TransactOpts,
	chainId int64,
	keys []string,
	values []*big.Int,
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
		contents, err := io.ReadAll(response.Body)
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

	for i, value := range values {
		// Create compressed argument with values/timestamps
		cValue := value
		cValue = cValue.Lsh(cValue, 128)
		cValue = cValue.Add(cValue, big.NewInt(timestamp))
		cValues = append(cValues, cValue)
		// Debug logs.
		log.Debugf("key -- value: %s -- %s", keys[i], value.String())
		log.Debug("cValue: ", cValue.String())
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
