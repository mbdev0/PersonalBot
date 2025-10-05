package tasks

import "math/big"

type Option func(Task)

func WithSlippage(s float64) Option {
	return func(t Task) { t.SetSlippage(s) }
}
func WithComputeUnits(cu uint32) Option {
	return func(t Task) { t.SetComputeUnit(cu) }
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

func WithSellPositionId(id string) SellOption {
	return func(t *SellTask) { t.mu.Lock(); defer t.mu.Unlock(); t.Position_id = id }
}
