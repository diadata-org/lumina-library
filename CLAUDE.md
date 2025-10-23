# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

lumina-library is a Go library (v1.23.0) that provides decentralized price data aggregation for DIA (Decentralized Information & Analytics). It collects real-time trade data from multiple centralized exchanges (CEXs) and decentralized exchanges (DEXs), applies filters to aggregate prices, and supports on-chain price publishing via smart contracts.

## Development Commands

### Building and Testing

```bash
# Tidy dependencies
go mod tidy

# Build all packages
go build -v ./...

# Run all tests with verbose output
go test -v ./...

# Run tests for a specific package
go test -v ./scrapers
go test -v ./models
go test -v ./filters

# Run a single test function
go test -v ./scrapers -run TestBinanceScraper
```

### Linting and Formatting

```bash
# Format all Go files (idiomatic formatting)
go fmt ./...

# Run go vet for static analysis
go vet ./...
```

## Core Architecture

### Pipeline Flow

The library implements a **3-stage pipeline architecture** for price aggregation:

1. **Scrapers** → Collect real-time trade data from exchanges via WebSocket connections
2. **Collector** → Aggregates trades into atomic TradesBlocks per exchange-pair
3. **Processor** → Applies filters and metafilters to produce aggregated prices

```
Scrapers → tradesChannelIn → Collector → tradesblockChannel → Processor → filtersChannel
```

### Key Components

#### 1. Scrapers (`scrapers/`)

- **Interface**: `Scraper` with `TradesChannel()` and `Close()` methods
- **Supported Exchanges**:
  - **CEXs**: Binance, Coinbase, Kraken, KuCoin, GateIO, ByBit, MEXC, OKEx, Crypto.com
  - **DEXs**: Uniswap v2/v3, PancakeSwap v3, Curve
- **Pattern**: Each scraper implements:
  - WebSocket connection management with `wsconn` wrapper
  - Watchdog timer with configurable timeout (`{EXCHANGE}_WATCHDOG` env var, default 300s)
  - Automatic failover and restart via `failoverChannel`
  - Context-based cancellation for graceful shutdown
- **Entry Point**: `RunScraper()` in `APIScraper.go` - switch statement routes to exchange-specific scrapers

#### 2. Collector (`scrapers/Collector.go`)

- Starts all scrapers for given exchange pairs and pools
- Collects trades from `tradesChannelIn` into `tradesblockMap`
- Groups trades by `exchangepairIdentifier` (exchange:quote-base)
- Triggered by `triggerChannel` to send atomic TradesBlocks to processor
- Handles scraper restarts via `failoverChannel`

#### 3. Processor (`processor/processor.go`)

- **Two-step aggregation**:
  1. **Per-exchange filtering**: Apply filter (e.g., LastPrice) to each atomic TradesBlock
  2. **Cross-exchange metafiltering**: Aggregate filter values across exchanges (e.g., Median)
- Filters: `filters/` (e.g., `LastPrice`)
- Metafilters: `metafilters/` (e.g., `Median`)
- Configuration via environment variables: `FILTER_TYPE`, `METAFILTER_TYPE`

#### 4. Data Models (`models/`)

Core types representing the data pipeline:

- **`Asset`**: Token with blockchain, address, decimals, symbol
- **`Trade`**: Atomic trade event with quote/base tokens, price, volume, timestamp, exchange
- **`Pair`**: Quote/base asset pair (e.g., ETH-USDC)
- **`ExchangePair`**: Exchange-specific pair mapping with symbols and verification flag
- **`Pool`**: DEX liquidity pool with address, tokens, fee, blockchain
- **`TradesBlock`**: Collection of trades for a single exchange-pair
  - `Atomic: true` indicates all trades are from the same exchange-pair
  - `EndTime`: Timestamp when block was triggered
- **`FilterPointPair`**: Aggregated price point with asset, value, timestamp, source type

#### 5. On-chain Integration (`onchain/`)

- Smart contract deployment and interaction
- Contract ABIs in `contracts/` for Uniswap v2/v3/v4, Curve, PancakeSwap
- Uses `go-ethereum` for Ethereum client operations

### Configuration Management

- **Environment Variables**: Accessed via `utils.Getenv(key, default)` throughout codebase
  - Log levels: `{PACKAGE}_LOG_LEVEL` (TRACE, DEBUG, INFO, WARN, ERROR)
  - Watchdog timeouts: `{EXCHANGE}_WATCHDOG` (seconds)
  - Filter/metafilter types: `FILTER_TYPE`, `METAFILTER_TYPE`
- **JSON Configuration**: Pool/asset whitelists loaded via `gonfig` in `init()` functions

### Logging

- **Library**: `logrus` with structured logging
- **Pattern**: Package-level logger initialized in `init()` function
- **Levels**: ERROR (critical failures), WARN (recoverable issues, stale data), INFO (lifecycle events, aggregated values), DEBUG (detailed diagnostics), TRACE (verbose)
- **Context**: Include exchange name, asset symbol, pair, trade count in log messages

## Adding a New Exchange Scraper

When adding a new CEX scraper, follow this pattern:

1. **Create scraper file** `scrapers/{Exchange}.go`:

   - Define scraper struct with `pairScrapers` map, `tradesChannel`, `wg`
   - Implement `TradesChannel()` and `Close()` methods (Scraper interface)
   - Create `New{Exchange}Scraper()` constructor
   - Implement WebSocket subscription and message parsing
   - Use `wsconn` wrapper for connection management
   - Handle exchange-specific pair normalization

2. **Add to `APIScraper.go`**:

   - Add constant: `const {EXCHANGE}_EXCHANGE = "{ExchangeName}"`
   - Add case to `RunScraper()` switch statement with watchdog pattern
   - Use environment variable `{EXCHANGE}_WATCHDOG` for timeout

3. **Write tests** `scrapers/{Exchange}_test.go`:

   - Use table-driven tests with test cases slice
   - Mock WebSocket responses with `ws_mock_test.go` patterns
   - Test pair normalization, message parsing, error handling

4. **Environment variable**: Document `{EXCHANGE}_WATCHDOG` configuration

## Testing Conventions

- **Framework**: Standard Go `testing` package
- **File naming**: `*_test.go` in same package as code under test
- **Function naming**: `Test{FunctionName}` (e.g., `TestMedian`)
- **Pattern**: Table-driven tests with slice of test cases
  ```go
  testCases := []struct {
      name     string
      input    []float64
      expected float64
  }{
      {"empty", []float64{}, 0},
      {"single", []float64{5.0}, 5.0},
  }
  ```
- **Assertions**: Use `reflect.DeepEqual()` for struct/slice comparisons
- **Coverage Focus**: Edge cases (empty inputs, nil pointers), error paths, boundary conditions

## Go Idioms and Best Practices

### Error Handling

- Always check `if err != nil` immediately after function calls
- Wrap errors with context: `fmt.Errorf("context: %w", err)`
- Log errors at appropriate level (ERROR for critical, WARN for recoverable)
- Return early on errors to reduce nesting

### Concurrency

- Use `sync.WaitGroup` to coordinate goroutine lifecycle (pass as `*sync.WaitGroup`)
- Use `context.Context` for cancellation signals in scrapers
- Clean up goroutines with `defer wg.Done()` and context cancellation
- Avoid channel deadlocks: ensure senders and receivers are properly coordinated
- Use buffered channels where appropriate to prevent blocking

### Memory Management

- Pre-allocate slices with known capacity: `make([]Trade, 0, capacity)`
- Clean up resources with `defer` (e.g., `defer conn.Close()`)
- Avoid retaining large objects in long-lived maps

### Naming Conventions

- **Exported** (public): PascalCase (e.g., `TradesChannel`, `LastPrice`)
- **Unexported** (private): camelCase (e.g., `tradesChannel`, `pairScrapers`)
- **Constants**: SCREAMING_SNAKE_CASE for exchange names (e.g., `BINANCE_EXCHANGE`)

### Documentation

- Add godoc comments for all exported functions, types, and methods
- Format: `// FunctionName does X and returns Y.`
- Explain "why" in comments, not "what" (code shows "what")

## Blockchain/DeFi Security Considerations

When working with blockchain-related code:

1. **Input Validation**: Validate all external data (API responses, blockchain data, user config)
2. **Secret Management**: Never hardcode private keys, API keys, or RPC URLs - use environment variables
3. **Contract Verification**: Verify contract addresses before interaction
4. **Price Manipulation**: Consider flash loan attacks and oracle manipulation when aggregating prices
5. **Integer Safety**: Use `big.Int` for large numbers and check for overflows/underflows
6. **Gas Optimization**: Batch RPC calls where possible, cache frequently accessed on-chain data
7. **Transaction Safety**: Validate gas, nonce, and value parameters before sending transactions

## Common Patterns

### Watchdog Timer Pattern (CEX scrapers)

All CEX scrapers use a watchdog to detect stale connections:

```go
watchdogTicker := time.NewTicker(time.Duration(watchdogDelay) * time.Second)
lastTradeTime := time.Now()

for {
    select {
    case trade := <-scraper.TradesChannel():
        lastTradeTime = time.Now()
        tradesChannel <- trade
    case <-watchdogTicker.C:
        duration := time.Since(lastTradeTime)
        if duration > time.Duration(watchdogDelay)*time.Second {
            scraper.Close(cancel)
            failoverChannel <- EXCHANGE_NAME
            return
        }
    }
}
```

### Atomic TradesBlock Construction

The Collector builds atomic TradesBlocks per exchange-pair:

```go
exchangepairIdentifier := exchangepair.ExchangePairIdentifier(trade.Exchange.Name)

if _, ok := tradesblockMap[exchangepairIdentifier]; !ok {
    tradesblockMap[exchangepairIdentifier] = models.TradesBlock{
        Trades: []models.Trade{trade},
        Pair:   exchangepair,
    }
} else {
    tradesblock := tradesblockMap[exchangepairIdentifier]
    tradesblock.Trades = append(tradesblock.Trades, trade)
    tradesblockMap[exchangepairIdentifier] = tradesblock
}
```

### Package Initialization

Each package initializes its logger in `init()`:

```go
var log *logrus.Logger

func init() {
    log = logrus.New()
    log.SetLevel(utils.GetLogLevel("PACKAGE_LOG_LEVEL"))
}
```

## Dependencies

Key external dependencies:

- `github.com/ethereum/go-ethereum` - Ethereum client (contract interaction, ABI binding)
- `github.com/sirupsen/logrus` - Structured logging
- `github.com/gorilla/websocket` - WebSocket client for CEX scrapers
- `github.com/daoleno/uniswapv3-sdk` - Uniswap v3 integration
- `github.com/prometheus/client_golang` - Metrics and monitoring
- `github.com/tidwall/gjson` - Fast JSON parsing
- `github.com/tkanos/gonfig` - JSON configuration loading

When adding new dependencies, ensure compatibility with Go 1.23.0 and vet for security/maintenance status.
