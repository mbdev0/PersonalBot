CREATE TABLE "trading_tasks" (
	"id"	INTEGER,
	"trading_type"	TEXT NOT NULL,
	"wallet_id"	TEXT NOT NULL,
	"slippage"	FLOAT NOT NULL,
	"compute_units"	FLOAT NOT NULL,
	"config"	TEXT,
	"time_created"	INTEGER NOT NULL,
	FOREIGN KEY("wallet_id") REFERENCES "crypto_wallets"("id"),
	PRIMARY KEY("id")
);