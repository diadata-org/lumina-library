package processor

import (
	"flag"
	"strconv"

	"github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/sirupsen/logrus"
)

const (
	FILTER_LAST_PRICE = models.FilterType("LastPrice")
	FILTER_VWAP       = models.FilterType("VWAP")
	METAFILTER_MEDIAN = models.MetafilterType("Median")
	METAFILTER_VWAP   = models.MetafilterType("VWAP")
)

// For processing, all filters with timestamp older than time.Now()-toleranceSeconds are discarded.
var (
	toleranceSeconds       int64
	watchFeedConfigSeconds int64
	log                    *logrus.Logger
	// These can be removed with the new filter layout.
	filterType     = utils.Getenv("FILTER_TYPE", string(FILTER_LAST_PRICE))
	metaFilterType = utils.Getenv("METAFILTER_TYPE", string(METAFILTER_MEDIAN))
)

func init() {
	var err error
	flag.Parse()
	log = logrus.New()
	loglevel, err := logrus.ParseLevel(utils.Getenv("LOG_LEVEL_PROCESSOR", "info"))
	if err != nil {
		log.Errorf("Parse log level: %v.", err)
	}
	log.SetLevel(loglevel)

	toleranceSeconds, err = strconv.ParseInt(utils.Getenv("TOLERANCE_SECONDS", "2"), 10, 64)
	if err != nil {
		log.Errorf("Parse TOLERANCE_SECONDS environment variable: %v.", err)
	}

	watchFeedConfigSeconds, err = strconv.ParseInt(utils.Getenv("WATCH_FEED_CONFIG_INTERVAL", "60"), 10, 64)
	if err != nil {
		log.Errorf("Parse TOLERANCE_SECONDS environment variable: %v.", err)
	}

	// FILTER_TYPE=LastPrice does not produce volume data, so pairing it with
	// METAFILTER_TYPE=VWAP would silently degrade to an equal-weight average
	// while still reporting Name="vwap" in the output.
	if filterType == string(FILTER_LAST_PRICE) && metaFilterType == string(METAFILTER_VWAP) {
		log.Fatalf(
			"invalid configuration: FILTER_TYPE=%s is incompatible with METAFILTER_TYPE=%s — "+
				"LastPrice produces no volume data so VWAP weighting cannot be applied. "+
				"Use METAFILTER_TYPE=%s or change FILTER_TYPE to %s.",
			filterType, metaFilterType, METAFILTER_MEDIAN, FILTER_VWAP,
		)
	}
}
