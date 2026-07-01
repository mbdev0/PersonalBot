package tasks

import (
	"math/big"
	"personal_bot/backend/pkg/logger"
)

type Option func(Configurable)

func WithSlippage(s float64) Option {
	return func(t Configurable) { t.SetSlippage(s) }
}
func WithComputeUnits(cu uint32) Option {
	return func(t Configurable) { t.SetComputeUnit(cu) }
}

func WithUnixTime(time int64) Option {
	return func(t Configurable) { t.SetTime(time) }
}

func WithStrategyId(id int64) Option {
	return func(t Configurable) { t.SetStrategyId(id) }
}

func WithHttpNode(rpcNode string) Option {
	return func(t Configurable) { t.SetHttpNode(rpcNode) }
}

func WithWS(ws string) Option {
	return func(t Configurable) { t.SetWSNode(ws) }
}

func WithRPCGroupId(id int64) Option {
	return func(t Configurable) { t.SetRPCGroupId(id) }
}

func WithProgram(program string) Option {
	return func(t Configurable) { t.SetProgram(program) }
}

func WithRetries(retries uint16) Option {
	return func(t Configurable) { t.SetRetries(retries) }
}

func WithRetriesDelayMs(retriesDelayMs uint32) Option {
	return func(t Configurable) { t.SetRetriesDelayMS(retriesDelayMs) }
}

func WithLogger(logger *logger.TaskLogger) Option {
	return func(t Configurable) { t.SetLogger(logger) }
}

type BuyOption func(*BuyTask)

func WithBuyAmount(amnt *big.Int) BuyOption {
	return func(t *BuyTask) { t.mu.Lock(); defer t.mu.Unlock(); t.BuyAmount = amnt }
}
func WithBuyFee(fee float64) BuyOption {
	return func(t *BuyTask) { t.mu.Lock(); defer t.mu.Unlock(); t.Fee = fee }
}

type SellOption func(*SellTask)

func WithSellAmount(amnt float64) SellOption {
	return func(t *SellTask) { t.mu.Lock(); defer t.mu.Unlock(); t.SellPercentage = amnt }
}
func WithSellFee(fee float64) SellOption {
	return func(t *SellTask) { t.mu.Lock(); defer t.mu.Unlock(); t.Fee = fee }
}

func WithSellPositionId(id *int64) SellOption {
	return func(t *SellTask) { t.mu.Lock(); defer t.mu.Unlock(); t.Position_id = id }
}
