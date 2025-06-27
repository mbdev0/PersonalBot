package tasks

import (
	"math/big"

	"github.com/gagliardetto/solana-go"
)

type BuyTask struct {
	Wallet       solana.PrivateKey `validate:"required"`
	TokenAddress solana.PublicKey  `validate:"required"`
	BuyAmount    big.Int           `validate:"required,gtZero"`
	Slippage     float64           `validate:"required,gt=0,lt=1"` // Slippage percentage (0.0 to 1.0)
	BuyFee       float64           `validate:"required,gt=0"`
	ComputeUnits uint32            `validate:"required,min=1"`
}
