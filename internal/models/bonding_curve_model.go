package models

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type BondingCurve struct {
	VirtualTokenReserves big.Int
	VirtualSolReserves   big.Int
	RealTokenReserves    big.Int
	RealSolReserves      big.Int
	MaxTokens            big.Int
	IsCompleted          bool
	DevWallet               solana.PublicKey
}
