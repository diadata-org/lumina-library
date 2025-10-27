package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/diadata-org/lumina-library/utils"
	"github.com/tkanos/gonfig"
)

// Pool is the container for liquidity pools on DEXes.
type Pool struct {
	Exchange      Exchange      `json:"Exchange"`
	Blockchain    Blockchain    `json:"Blockchain"`
	Address       string        `json:"Address"`
	Assetvolumes  []AssetVolume `json:"Assetvolumes"`
	Order         int           `json:"Order"`
	Time          time.Time     `json:"Time"`
	WatchDogDelay int64         `json:"WatchDogDelay"`
}

type AssetVolume struct {
	Asset  Asset   `json:"Asset"`
	Volume float64 `json:"Volume"`
	Index  uint8   `json:"Index"`
}

// MakePoolMap maps an exchange name on the underlying slice of pool structs.
func MakePoolMap(pools []Pool) map[string][]Pool {
	poolMap := make(map[string][]Pool)
	for _, pool := range pools {
		poolMap[pool.Exchange.Name] = append(poolMap[pool.Exchange.Name], pool)
	}
	return poolMap
}

// PoolsFromJSON reads a file in the format
//
//	"exchange": { "Pools": [ { "Address": "...", "Order": "2", "WatchDogDelay": 300 }, ... ] }
//
// and returns []Pool, with Exchange.Name set to the input exchange (e.g. "UniswapV2").
func PoolsFromConfigFile(exchange string) ([]Pool, error) {
	type filePool struct {
		Address       string `json:"Address"`
		Order         string `json:"Order"`         // note: the file contains strings
		WatchDogDelay int64  `json:"WatchDogDelay"` // optional, default to 300 if not provided
	}
	type fileSchema struct {
		Pools []filePool `json:"Pools"`
	}

	path := getPath2Config("pools") + strings.TrimSpace(exchange) + ".json"

	// read and deserialize the file
	var cfg fileSchema

	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("decode %s: %w", path, err)
	}

	// construct the output
	out := make([]Pool, 0, len(cfg.Pools))
	for i, p := range cfg.Pools {
		if strings.TrimSpace(p.Address) == "" {
			return out, fmt.Errorf("pools[%d] has empty address", i)
		}
		ord, err := strconv.Atoi(strings.TrimSpace(p.Order))
		if err != nil {
			return out, fmt.Errorf("pools[%d] bad order %q: %w", i, p.Order, err)
		}
		wd := p.WatchDogDelay
		if wd <= 0 {
			wd = 300
		}
		out = append(out, Pool{
			Exchange:      Exchange{Name: exchange},
			Address:       p.Address,
			Order:         ord,
			WatchDogDelay: wd,
		})
	}
	return out, nil
}

// aggregate multiple pools from multiple config files
func PoolsFromConfigFiles(exchanges []string) ([]Pool, error) {
	if len(exchanges) == 0 {
		return []Pool{}, nil
	}
	if len(exchanges) == 1 && len(strings.TrimSpace(exchanges[0])) == 0 {
		return []Pool{}, nil
	}

	var (
		all  []Pool
		errs []error
	)
	for _, ex := range exchanges {
		ex = strings.TrimSpace(ex)
		if ex == "" {
			continue
		}
		pools, err := PoolsFromConfigFile(ex)
		if err != nil {
			errs = append(errs, err)
			continue // do not interrupt, collect results from other exchanges
		}
		all = append(all, pools...)
	}
	// join errors
	if len(errs) > 0 {
		return all, errors.Join(errs...)
	}
	return all, nil
}

func GetPoolsFromConfig(exchange string) ([]Pool, error) {
	path := utils.GetPath("pools/", exchange)
	type Pools struct {
		Pools []Pool
	}
	var p Pools
	err := gonfig.GetConf(path, &p)
	if err != nil {
		return []Pool{}, err
	}
	return p.Pools, nil
}
