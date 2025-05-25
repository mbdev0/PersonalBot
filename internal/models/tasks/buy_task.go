package tasks

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type BuyTask struct {
	Wallet       solana.PrivateKey
	TokenAddress solana.PublicKey
	BuyAmount    big.Int
	Slippage     float64
	BuyFee       float64
	ComputeUnits uint32
}
