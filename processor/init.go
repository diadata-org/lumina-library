package processor

import (
	"flag"
	"strconv"

	models "github.com/diadata-org/lumina-library/models"
	"github.com/diadata-org/lumina-library/utils"
	"github.com/sirupsen/logrus"
)

// For processing, all filters with timestamp older than time.Now()-toleranceSeconds are discarded.
var (
	toleranceSeconds       int64
	watchFeedConfigSeconds int64
	log                    *logrus.Logger
	// These can be removed with the new filter layout.
	filterTypeGlobal     = models.FilterType(utils.Getenv("FILTER_TYPE", string(models.FILTER_LAST_PRICE)))
	metaFilterTypeGlobal = models.MetafilterType(utils.Getenv("METAFILTER_TYPE", string(models.METAFILTER_MEDIAN)))
	usd                  = models.Asset{Blockchain: "Fiat", Address: "840"}
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

	toleranceSeconds, err = strconv.ParseInt(utils.Getenv("TOLERANCE_SECONDS", "20"), 10, 64)
	if err != nil {
		log.Errorf("Parse TOLERANCE_SECONDS environment variable: %v.", err)
	}

	watchFeedConfigSeconds, err = strconv.ParseInt(utils.Getenv("WATCH_FEED_CONFIG_INTERVAL", "60"), 10, 64)
	if err != nil {
		log.Errorf("Parse WATCH_FEED_CONFIG_INTERVAL environment variable: %v.", err)
	}

	// FILTER_TYPE=LastPrice does not produce volume data, so pairing it with
	// METAFILTER_TYPE=VWAP would silently degrade to an equal-weight average
	// while still reporting Name="vwap" in the output.
	if filterTypeGlobal == models.FILTER_LAST_PRICE && metaFilterTypeGlobal == models.METAFILTER_VWAP {
		log.Fatalf(
			"invalid configuration: FILTER_TYPE=%s is incompatible with METAFILTER_TYPE=%s — "+
				"LastPrice produces no volume data so VWAP weighting cannot be applied. "+
				"Use METAFILTER_TYPE=%s or change FILTER_TYPE to %s.",
			filterTypeGlobal, metaFilterTypeGlobal, models.METAFILTER_MEDIAN, models.FILTER_VWAP,
		)
	}
}
