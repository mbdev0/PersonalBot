CREATE TABLE [trading_tasks] (
    id INTEGER PRIMARY KEY, 
    trading_type TEXT NOT NULL,
    wallet_id TEXT NOT NULL,
    slippage FLOAT NOT NULL,
    compute_units FLOAT NOT NULL,
    config TEXT,
    FOREIGN KEY(wallet_id) REFERENCES crypto_wallets(id) 
)