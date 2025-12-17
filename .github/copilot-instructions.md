# lumina-library: GitHub Copilot Instructions

## About This Repository

This is a Go library (v1.23.0) for decentralized price data aggregation for DIA (Decentralized Information & Analytics). We collect real-time trade data from centralized exchanges (CEXs) and decentralized exchanges (DEXs), apply filters to aggregate prices, and support on-chain price publishing via smart contracts.

## Architecture Overview

Our codebase follows a **3-stage pipeline pattern**:

```
Scrapers → Collector → Processor → On-Chain Publishing
```

1. **Scrapers**: WebSocket connections to exchanges that produce `Trade` events
2. **Collector**: Aggregates trades into time-based `TradesBlock` batches
3. **Processor**: Applies filters (per-exchange) and metafilters (cross-exchange) to compute final prices
4. **On-Chain**: Publishes validated prices to smart contracts

## Repository Structure

- `scrapers/`: WebSocket data collection from exchanges (APIScraper.go is the main router)
- `models/`: Core data structures (`Trade`, `Pair`, `Asset`, `Pool`, `TradesBlock`, `FilterPointPair`)
- `processor/`: Two-step aggregation (filters per exchange, metafilters across exchanges)
- `filters/`: Per-exchange price filters (e.g., `LastPrice`, `VWAP`)
- `metafilters/`: Cross-exchange aggregators (e.g., `Median`, `Mean`)
- `onchain/`: Smart contract integration for price publishing
- `contracts/`: Contract ABIs (Uniswap, Curve, PancakeSwap, etc.)
- `utils/`: Shared utilities (environment variables, Ethereum helpers, HTTP clients)

## Code Generation Guidelines

### Scraper Implementation

When creating new scrapers:
- Implement the `Scraper` interface with `TradesChannel()` and `Close()` methods
- Use robust websocket reconnection logic with exponential backoff
- Include watchdog timers to detect stale connections
- Respect exchange rate limits
- Handle websocket errors gracefully and log appropriately

### Data Models

Our core types:
- Use `*big.Int` for all price and volume calculations to prevent overflows
- Always validate decimals and handle zero prices
- Check for nil pointers before dereferencing
- Use proper time handling with `time.Time` for timestamps

### Error Handling

Follow Go idiomatic patterns:
- Check errors immediately: `if err != nil { return fmt.Errorf("context: %w", err) }`
- Wrap errors with context using `%w` for error chain visibility
- Use early returns to reduce nesting
- Log errors at appropriate levels (ERROR for critical, WARN for recoverable)

### Concurrency Patterns

- Use `context.Context` for cancellation propagation in goroutines
- Synchronize with `sync.WaitGroup` when needed
- Avoid channel deadlocks by using buffered channels or select with default
- Always clean up resources with `defer` (especially websocket/HTTP connections)

### Configuration & Environment

- Access environment variables via `utils.Getenv(key, default)` instead of `os.Getenv()`
- Document all new environment variables
- Never hardcode secrets, API keys, private keys, or RPC URLs
- Store sensitive data in environment variables or secure vaults

### Security (Critical for Blockchain/DeFi)

- Validate all external inputs (API responses, blockchain data, config files)
- Never expose secrets in logs or error messages
- Verify smart contract addresses before interaction
- Use `big.Int` for all financial calculations to prevent overflow attacks
- Consider price manipulation resistance (flash loan attacks, oracle manipulation)
- Validate transaction parameters (gas limits, nonce, value) for on-chain operations

### Testing Conventions

- Use table-driven tests with slices of test cases
- Name test files `*_test.go`
- Test edge cases: nil pointers, zero values, empty slices, stale data
- Test error paths, not just happy paths
- Use websocket mocks or fixtures for scraper tests

### Logging Standards

- Use the package-level logger pattern with logrus
- Log levels: `ERROR` (critical failures), `WARN` (recoverable issues), `INFO` (lifecycle events), `DEBUG` (verbose)
- Never log sensitive data (private keys, API secrets, user data)
- Include context in log messages (exchange name, pair, block number, etc.)

### Performance Considerations

- Pre-allocate slices/maps when size is known: `make([]T, 0, capacity)`
- Use non-blocking channel operations with select where appropriate
- Batch or cache RPC calls to minimize blockchain network overhead
- Profile and benchmark critical paths in price aggregation

### Naming Conventions

- Exported identifiers: PascalCase (`TradesChannel`, `FilterPoint`)
- Unexported identifiers: camelCase (`tradesBuffer`, `reconnectDelay`)
- Document all exported functions/types with godoc comments
- Use descriptive names that reflect domain concepts (avoid abbreviations unless standard)

### Dependencies

- Maintain strict compatibility with **Go 1.23.0**
- Vet new third-party dependencies for security and maintenance status
- Regenerate contract ABI bindings from official sources
- Keep dependencies minimal and justified

## Code Quality Expectations

We value:
- **Modularity**: Single Responsibility Principle, clean interfaces
- **Robustness**: Comprehensive error handling, edge case coverage
- **Security**: Input validation, no secrets in code/logs, secure blockchain interactions
- **Clarity**: Readable code with clear control flow, early returns, minimal nesting
- **Observability**: Appropriate logging, metrics for monitoring (Prometheus)

CI automation validates basic correctness (`go mod tidy`, `go build`, `go test`, `go fmt`). Focus on architecture, design, security, and complex logic rather than formatting.
