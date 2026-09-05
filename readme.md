# AIO Bot

A desktop trading bot for Solana's pump.fun ecosystem, built as a native Go application with a React front end. It watches the chain in real time, executes trades against the pump.fun bonding curve and AMM programs, and manages positions automatically according to user-defined strategies — all through a live dashboard rather than a config file and a prayer.

Currently focused on Solana pump.fun; the architecture (RPC groups, program registry, strategy interfaces) is designed to extend to other EVM and Web2 targets later.

## What it does

- **Watches the chain live.** Subscribes directly to Solana RPC websocket streams for account and transaction data, decoding raw pump.fun bonding-curve and AMM instructions (buys, sells, migrations, coin creation) without relying on a third-party indexer.
- **Builds and signs its own transactions.** A custom transaction builder constructs pump.fun buy/sell instructions from scratch using Borsh serialization, rather than shelling out to an SDK.
- **Runs configurable trading strategies.** Supports quick buy/sell, a spam strategy for aggressive entry, and AFK automation, each implementing a shared strategy interface so new strategies can be dropped in.
- **Automates exits with take-profit / stop-loss.** Exit conditions can be defined by absolute price, percentage gain/loss, or target market cap, and are evaluated continuously against live position data.
- **Tracks positions in real time.** A pub/sub hub pushes position and task updates over a websocket connection straight into the React dashboard as they happen.
- **Distributes load across multiple RPCs.** RPC endpoints are organized into groups so trading load and subscriptions can be spread across providers instead of hammering a single node.
- **Notifies you on Discord.** Rich embed messages report successful buys/sells and failures, complete with links out to Solscan, Axiom, and pump.fun.
- **Logs every task to its own file.** Each trading task writes to a dedicated log file, viewable through a built-in terminal UI (a Bubble Tea TUI) launched straight from the desktop app.
- **Persists state locally.** Wallets, tasks, and RPC configuration are stored in SQLite with versioned migrations, so the bot picks up where it left off after a restart.

## Tech stack

**Backend:** Go, [Wails v3](https://wails.io) (native desktop shell), `solana-go`, Borsh, SQLite, Wire (dependency injection), Bubble Tea (TUI)

**Frontend:** React 19, TypeScript, Vite, Tailwind CSS, Radix UI / shadcn components, TanStack Query & Table, native WebSockets for live updates
