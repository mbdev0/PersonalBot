package tasks

import "github.com/gagliardetto/solana-go"

type SellTask struct {
	PublicKey    solana.PublicKey
	TokenAddress solana.PublicKey
	Wallet       solana.PrivateKey
	Amount       int
	MinSolAmount int
	ComputeUnits uint32
	BuyFee       float64
}
