package models

import "github.com/gagliardetto/solana-go"

type CompletePumpAmmMigrationEvent struct {
	User             solana.PublicKey `json:"user"`
	Mint             solana.PublicKey `json:"mint"`
	MintAmount       uint64           `json:"mint_amount"`
	SolAmount        uint64           `json:"sol_amount"`
	PoolMigrationFee uint64           `json:"pool_migration_fee"`
	BondingCurve     solana.PublicKey `json:"bonding_curve"`
	Timestamp        int64            `json:"timestamp"`
	Pool             solana.PublicKey `json:"pool"`
	QuoteMint        solana.PublicKey `json:"quote_mint"`
}
