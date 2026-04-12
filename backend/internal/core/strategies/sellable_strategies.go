package strategies

import "personal_bot/internal/core/wallets"

type SellableStrategy interface {
	GetWallet() wallets.SolanaWallet
	GetComputeUnits() uint32
	GetSlippage() float64
	GetSellFee() *float64
	StrategyTaskId() int64
}
