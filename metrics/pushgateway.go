package metrics

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math"
	"math/big"
	"runtime"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/shirou/gopsutil/cpu"
	log "github.com/sirupsen/logrus"
)

func PushMetricsToPushgateway(m *Metrics, startTime time.Time, conn *ethclient.Client, privateKey *ecdsa.PrivateKey, deployedContract string) {
	const sampleWindowSize = 5                         // Number of samples to calculate the rolling average
	cpuSamples := make([]float64, 0, sampleWindowSize) // Circular buffer for CPU usage samples

	for {
		uptime := time.Since(startTime).Hours()
		m.uptime.Set(uptime)

		// Update memory usage
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		memoryUsageMB := float64(memStats.Alloc) / 1024 / 1024 // Convert bytes to megabytes
		m.memoryUsage.Set(memoryUsageMB)

		// Update CPU usage using gopsutil with smoothing
		percent, err := cpu.Percent(0, false)
		if err != nil {
			log.Errorf("Error gathering CPU usage: %v", err)
		} else if len(percent) > 0 {
			// Add the new sample to the buffer
			if len(cpuSamples) == sampleWindowSize {
				cpuSamples = cpuSamples[1:] // Remove the oldest sample if buffer is full
			}
			cpuSamples = append(cpuSamples, percent[0])

			// Calculate the rolling average
			var sum float64
			for _, v := range cpuSamples {
				sum += v
			}
			avgCPUUsage := sum / float64(len(cpuSamples))
			m.cpuUsage.Set(avgCPUUsage) // Update the metric with the smoothed value
		}

		// Get the gas wallet balance
		gasBalance, err := getAddressBalance(conn, privateKey)
		if err != nil {
			log.Errorf("Failed to fetch address balance: %v", err)
		}
		m.gasBalance.Set(gasBalance)

		// Get the latest event timestamp
		lastUpdateTime, err := getLatestEventTimestamp(conn, deployedContract)
		if err != nil {
			log.Errorf("Error fetching latest event timestamp: %v", err)
		}
		m.lastUpdateTime.Set(lastUpdateTime)

		// Push metrics to the Pushgateway
		pushCollector := push.New(m.pushGatewayURL, m.jobName).
			Collector(m.uptime).
			Collector(m.cpuUsage).
			Collector(m.memoryUsage).
			Collector(m.Contract).
			Collector(m.ExchangePairs).
			Collector(m.gasBalance).
			Collector(m.lastUpdateTime).
			Collector(m.chainID).
			Collector(m.imageVersion).
			Collector(TradesTotal).
			Collector(WatchdogFailoversTotal)

		if err := pushCollector.
			BasicAuth(m.authUser, m.authPassword).
			Push(); err != nil {
			log.Errorf("Could not push metrics to Pushgateway: %v", err)
		} else {
			log.Printf("Metrics pushed successfully to Pushgateway")
		}

		time.Sleep(30 * time.Second) // update metrics every 30 seconds
	}
}

func getAddressBalance(client *ethclient.Client, privateKey *ecdsa.PrivateKey) (float64, error) {

	publicKeyECDSA, ok := privateKey.Public().(*ecdsa.PublicKey)
	if !ok {
		return math.NaN(), fmt.Errorf("failed to cast public key to ECDSA")
	}

	balance, err := client.BalanceAt(context.Background(), crypto.PubkeyToAddress(*publicKeyECDSA), nil)
	if err != nil {
		return math.NaN(), fmt.Errorf("failed to get balance: %w", err)
	}

	balanceInDIA, _ := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18)).Float64()
	return balanceInDIA, nil
}

func getLatestEventTimestamp(client *ethclient.Client, contractAddress string) (float64, error) {
	// Get the latest block number
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return math.NaN(), fmt.Errorf("failed to fetch latest block header: %v", err)
	}
	latestBlock := header.Number.Int64()

	// Calculate the start block for the query
	startBlock := latestBlock - 1000
	if startBlock < 0 {
		startBlock = 0 // Ensure the start block is not negative
	}

	// Define filter query for the last 'blockRange' blocks
	query := ethereum.FilterQuery{
		Addresses: []common.Address{common.HexToAddress(contractAddress)},
		FromBlock: big.NewInt(startBlock),
		ToBlock:   big.NewInt(latestBlock),
	}

	// Fetch logs for the specified block range
	logs, err := client.FilterLogs(context.Background(), query)
	if err != nil {
		return math.NaN(), fmt.Errorf("failed to fetch logs: %v", err)
	}

	// Check if logs are empty
	if len(logs) == 0 {
		return math.NaN(), fmt.Errorf("no events found in the last 1000 blocks")
	}

	// Get the latest timestamp from the last log
	lastLog := logs[len(logs)-1]
	blockHeader, err := client.HeaderByHash(context.Background(), lastLog.BlockHash)
	if err != nil {
		return math.NaN(), fmt.Errorf("failed to fetch block header for log: %v", err)
	}

	return float64(blockHeader.Time), nil
}
