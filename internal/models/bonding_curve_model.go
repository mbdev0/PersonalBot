package models

import "math/big"

type BondingCurve struct {
	VirtualTokenReserves big.Int
	VirtualSolReserves   big.Int
	RealTokenReserves    big.Int
	RealSolReserves      big.Int
	MaxTokens            big.Int
	IsCompleted          bool
}
