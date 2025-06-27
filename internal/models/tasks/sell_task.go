package tasks

import "github.com/gagliardetto/solana-go"

type SellTask struct {
	TokenAddress solana.PublicKey
	Wallet       solana.PrivateKey
	Amount       int
	ComputeUnits uint32
	SellFee      float64
	Slippage     float64
}
