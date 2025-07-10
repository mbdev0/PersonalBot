package tasks

import "github.com/gagliardetto/solana-go"

type SellTask struct {
	TokenAddress     solana.PublicKey
	Wallet           solana.PrivateKey
	ComputeUnits     uint32
	SellFee          float64
	Slippage         float64
	PercentageToSell float64
}
