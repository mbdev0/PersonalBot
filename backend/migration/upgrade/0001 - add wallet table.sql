CREATE TABLE [crypto_wallets] (
	id TEXT not null PRIMARY KEY UNIQUE,
	wallet_name TEXT NOT NULL UNIQUE, 
	chain TEXT not NULL,
	private_key TEXT not NULL
)