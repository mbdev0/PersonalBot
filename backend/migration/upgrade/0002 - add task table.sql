CREATE TABLE "tasks" (
	"id"	INTEGER,
	"task_type"	TEXT NOT NULL,
	"wallet_id"	TEXT NOT NULL,
	"slippage_percentage"	INTEGER NOT NULL,
	"compute_units"	INTEGER NOT NULL,
	"config"	TEXT,
	"strategy_id"	INTEGER,
	"state"	TEXT NOT NULL,
	"token"	TEXT NOT NULL,
	PRIMARY KEY("id"),
	FOREIGN KEY("strategy_id") REFERENCES "trading_tasks"("id") ON DELETE CASCADE,
	FOREIGN KEY("wallet_id") REFERENCES "crypto_wallets"("id")
)