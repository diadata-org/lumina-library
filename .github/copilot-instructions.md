# lumina-library: GitHub Copilot Code Review Instructions

This is a Go library (v1.23.0) for decentralized price data aggregation for DIA (Decentralized Information & Analytics). It collects real-time trade data from centralized exchanges (CEXs) and decentralized exchanges (DEXs), applies filters to aggregate prices, and supports on-chain price publishing via smart contracts.

## Review Philosophy

Focus on substantive improvements over perfection. Assume CI automation (go mod tidy, go build, go test, go fmt) has validated basic correctness. Your review should prioritize architecture, design, security, and complex logic over style and formatting.

## Code Review Checklist

### 1. Architectural Design

- Does the change align with the 3-stage pipeline pattern (Scrapers → Collector → Processor)?
- For scrapers: Does it implement the `Scraper` interface with `TradesChannel()` and `Close()` methods?
- Is the code modular and following Single Responsibility Principle?
- Are changes atomic (single purpose) or bundling unrelated modifications?

### 2. Functionality & Correctness

- Does the code correctly implement the intended price aggregation/filtering logic?
- Are edge cases handled: zero prices, missing decimals, nil pointers, empty trade blocks, stale data?
- For scrapers: Is websocket reconnection logic robust? Are rate limits respected?
- For blockchain interactions: Are contract calls properly validated?
- Are there potential race conditions in goroutines or channel deadlocks?

### 3. Security (Critical for Blockchain/DeFi)

- Are all external inputs (API responses, blockchain data, config files) validated?
- Are secrets, private keys, API keys, or RPC URLs hardcoded or exposed in logs?
- For smart contract interactions: Are contract addresses verified? Are ABI bindings current?
- Is price manipulation resistance considered (flash loan attacks, oracle manipulation)?
- Are integer overflows/underflows prevented in price calculations using `big.Int`?
- For on-chain operations: Are transaction parameters (gas, nonce, value) validated?

### 4. Go Best Practices

- Is error handling idiomatic with immediate `if err != nil` checks and context wrapping (`fmt.Errorf("context: %w", err)`)?
- Are goroutines properly managed with `sync.WaitGroup` or `context.Context` for cancellation?
- Do names follow Go conventions (PascalCase for exported, camelCase for unexported)?
- Are exported functions/types documented with godoc comments?
- Is control flow clear with early returns instead of deep nesting?
- Are environment variables accessed via `utils.Getenv(key, default)` instead of `os.Getenv()`?

### 5. Testing

- Do tests use table-driven patterns with slices of test cases?
- Are edge cases, error paths, and failure modes tested, not just happy paths?
- Are test files named `*_test.go` and using the standard `testing` package?
- For scrapers: Are websocket mocks or fixtures used appropriately?

### 6. Performance & Scalability

- Are goroutines synchronized properly? Are channel operations non-blocking where appropriate?
- Are large slices/maps pre-allocated with `make([]T, 0, capacity)`?
- Are resources cleaned up with `defer`?
- For blockchain code: Are RPC calls batched or cached to minimize network overhead?
- For scrapers: Are websocket connections managed with watchdog timers and retry logic?

### 7. Logging & Observability

- Are log levels appropriate (ERROR for critical failures, WARN for recoverable issues, INFO for lifecycle events)?
- Is the package-level logger pattern followed with logrus?
- Is sensitive data (private keys, API secrets) excluded from logs?
- Are Prometheus metrics updated for new components?

### 8. Dependencies

- Are new third-party dependencies necessary and vetted for security/maintenance?
- Is compatibility with Go 1.23.0 maintained (strict requirement)?
- Are environment variables documented?
- For contracts: Are ABI bindings regenerated from official sources?

## Repository Structure (for Context)

- `scrapers/`: WebSocket data collection from exchanges (APIScraper.go is the router)
- `models/`: Core data structures (Trade, Pair, Asset, Pool, TradesBlock, FilterPointPair)
- `processor/`: Two-step aggregation (filters per exchange, metafilters across exchanges)
- `filters/`: Per-exchange price filters (e.g., LastPrice)
- `metafilters/`: Cross-exchange aggregators (e.g., Median)
- `onchain/`: Smart contract integration for price publishing
- `contracts/`: Contract ABIs (Uniswap, Curve, PancakeSwap)
- `utils/`: Shared utilities (environment variables, Ethereum helpers, HTTP)

## Review Output Guidelines

- Focus on substantive issues that impact correctness, security, maintainability, or scalability
- Explain the engineering principle behind each suggestion
- Prefix minor/optional suggestions with "Nit:"
- Be constructive and specific with actionable feedback
- Reference specific files and line numbers when possible
