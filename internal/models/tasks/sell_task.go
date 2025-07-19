package tasks

import "github.com/gagliardetto/solana-go"

type SellTask struct {
	Wallet           solana.PrivateKey
	TokenAddress     solana.PublicKey
	PercentageToSell float64
	SellFee          float64
	Slippage         float64
	ComputeUnits     uint32
}
