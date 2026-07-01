package strategies

import (
	"personal_bot/backend/internal/core/wallets"
	"personal_bot/backend/pkg/logger"
)

type SellableStrategy interface {
	GetProgram() string
	GetWallet() wallets.SolanaWallet
	GetComputeUnits() uint32
	GetSlippage() float64
	GetSellFee() *float64
	StrategyTaskId() int64
	RPCGroupId() int64
	Logger() *logger.TaskLogger
}
