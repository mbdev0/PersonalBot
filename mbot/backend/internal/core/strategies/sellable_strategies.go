package strategies

import "personal_bot/backend/internal/core/wallets"

type SellableStrategy interface {
	GetProgram() string
	GetWallet() wallets.SolanaWallet
	GetComputeUnits() uint32
	GetSlippage() float64
	GetSellFee() *float64
	StrategyTaskId() int64
	RPCGroupId() int64
}
