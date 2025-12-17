package metrics

import (
	"crypto/ecdsa"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Metrics struct {
	uptime         prometheus.Gauge
	cpuUsage       prometheus.Gauge
	memoryUsage    prometheus.Gauge
	Contract       *prometheus.GaugeVec
	ExchangePairs  *prometheus.GaugeVec
	gasBalance     prometheus.Gauge
	lastUpdateTime prometheus.Gauge
	chainID        prometheus.Gauge
	imageVersion   *prometheus.GaugeVec
	pushGatewayURL string
	jobName        string
	authUser       string
	authPassword   string
}

func NewMetrics(reg *prometheus.Registry, pushGatewayURL, jobName, authUser, authPassword string, chainID int64, imageVersion string) *Metrics {
	m := &Metrics{
		uptime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "feeder",
			Name:      "uptime_hours",
			Help:      "Feeder Uptime in hours.",
		}),
		cpuUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "feeder",
			Name:      "cpu_usage_percent",
			Help:      "Feeder CPU usage in percent.",
		}),
		memoryUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "feeder",
			Name:      "memory_usage_megabytes",
			Help:      "Feeder Memory usage in megabytes.",
		}),
		Contract: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "feeder",
				Name:      "contract_info",
				Help:      "Feeder contract information.",
			},
			[]string{"contract"}, // Label to store the contract address
		),
		ExchangePairs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "feeder",
				Name:      "exchange_pairs",
				Help:      "List of exchange pairs to be pushed as labels for each Feeder.",
			},
			[]string{"exchange_pair"}, // Label to store each exchange pair
		),
		gasBalance: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "feeder",
			Name:      "gas_balance",
			Help:      "Gas wallet balance in DIA.",
		}),
		lastUpdateTime: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "feeder",
			Name:      "last_update_time",
			Help:      "Last update time in UTC timestamp.'",
		}),
		chainID: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "feeder",
			Name:      "chain_id",
			Help:      "The chain ID of the blockchain being monitored.",
		}),
		imageVersion: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: "feeder",
				Name:      "image_version",
				Help:      "The Docker image version of the feeder.",
			},
			[]string{"version"},
		),
		pushGatewayURL: pushGatewayURL,
		jobName:        jobName,
		authUser:       authUser,
		authPassword:   authPassword,
	}
	reg.MustRegister(m.uptime)
	reg.MustRegister(m.cpuUsage)
	reg.MustRegister(m.memoryUsage)
	reg.MustRegister(m.Contract)
	reg.MustRegister(m.ExchangePairs)
	reg.MustRegister(m.gasBalance)
	reg.MustRegister(m.lastUpdateTime)
	reg.MustRegister(m.chainID)
	reg.MustRegister(m.imageVersion)

	m.chainID.Set(float64(chainID))
	m.imageVersion.WithLabelValues(imageVersion).Set(1)

	return m
}

func StartPrometheusServer(m *Metrics, port string) {
	if m == nil {
		log.Errorf("Cannot start metrics server: metrics object is nil")
		return
	}

	// Register metrics with the default registry
	prometheus.DefaultRegisterer.MustRegister(m.uptime)
	prometheus.DefaultRegisterer.MustRegister(m.cpuUsage)
	prometheus.DefaultRegisterer.MustRegister(m.memoryUsage)
	prometheus.DefaultRegisterer.MustRegister(m.Contract)
	prometheus.DefaultRegisterer.MustRegister(m.ExchangePairs)
	prometheus.DefaultRegisterer.MustRegister(m.gasBalance)
	prometheus.DefaultRegisterer.MustRegister(m.lastUpdateTime)
	prometheus.DefaultRegisterer.MustRegister(m.chainID)
	prometheus.DefaultRegisterer.MustRegister(m.imageVersion)

	log.Printf("Starting metrics server on :%s", port)
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Printf("Failed to start metrics server: %v", err)
	}
}

func StartMetrics(
	conn *ethclient.Client,
	privateKey *ecdsa.PrivateKey,
	deployedContract string,
	pushgatewayURL string,
	authUser string,
	authPassword string,
	enablePrometheusServer string,
	nodeOperatorName string,
	metricsPort string,
	imageVersion string,
	chainID int64,
	exchangePairsEnv string) {
	// get hostname of the container so that we can display it in monitoring dashboards
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatalf("Failed to get hostname: %v", err)
	}

	// Check if metrics pushing to Pushgateway is enabled
	pushgatewayEnabled := pushgatewayURL != "" && authUser != "" && authPassword != ""

	// Check if Prometheus HTTP server is enabled
	prometheusServerEnabled := strings.ToLower(enablePrometheusServer) == "true"

	// Create the job name for metrics (used for both modes)
	jobName := MakeJobName(hostname, nodeOperatorName)

	if imageVersion == "" {
		imageVersion = "unknown" // fallback if not set
		log.Info("No image version found, using 'unknown'")
	}
	log.Infof("Image version: %s", imageVersion)

	// Set default pushgateway URL if enabled
	if pushgatewayEnabled {
		if pushgatewayURL == "" {
			pushgatewayURL = "https://pushgateway-auth.diadata.org"
		}
		log.Info("Metrics pushing enabled. Pushing to: ", pushgatewayURL)
	} else {
		log.Info("Metrics pushing to Pushgateway disabled")
	}

	// Create metrics object
	m := NewMetrics(
		prometheus.NewRegistry(),
		pushgatewayURL,
		jobName,
		authUser,
		authPassword,
		chainID,
		imageVersion,
	)

	// Start Prometheus HTTP server if enabled
	if prometheusServerEnabled {
		go StartPrometheusServer(m, metricsPort)
		log.Info("Prometheus HTTP server enabled on port:", metricsPort)
	} else {
		log.Info("Prometheus HTTP server disabled")
	}

	// Record start time for uptime calculation
	startTime := time.Now()

	// Move metrics setup here, right before the blocking call
	// Only setup metrics collection if metrics are enabled and metrics object exists
	if pushgatewayEnabled && m != nil {
		// Set the static contract label for Prometheus monitoring
		m.Contract.WithLabelValues(deployedContract).Set(1)

		exchangePairsList := strings.Split(exchangePairsEnv, ",")
		for _, pair := range exchangePairsList {
			pair = strings.TrimSpace(pair) // Clean whitespace
			if pair != "" {
				m.ExchangePairs.WithLabelValues(pair).Set(1)
			}
		}

		// Push metrics to Pushgateway if enabled
		go PushMetricsToPushgateway(m, startTime, conn, privateKey, deployedContract)
	}
}
