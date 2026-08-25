# MarketIQ

> Real-time crypto market quality and smart execution engine built in Go.

MarketIQ is an engineering-focused system for analyzing cryptocurrency market
microstructure and simulating execution strategies using real-time exchange
data.

The project currently connects to CoinDCX's BTC/USDT order-book stream,
maintains an in-memory order book, and computes core market metrics such as
best bid, best ask, mid price, and spread.

The system is being developed incrementally toward a market-quality analysis
and smart-execution engine.

---

## Architecture

```text
                 CoinDCX
                    │
                    ▼
           Real-Time Market Data
                    │
                    ▼
             Market Data Client
                    │
                    ▼
              Order Book Engine
                    │
          ┌─────────┴─────────┐
          ▼                   ▼
      Best Bid             Best Ask
          │                   │
          └─────────┬─────────┘
                    ▼
             Market Quality
               Analysis
                    │
                    ▼
             Execution Engine
                    │
                    ▼
               Backtesting