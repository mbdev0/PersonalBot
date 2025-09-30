package sell

import (
	"math/big"
	"pump_fun/internal/core/position"
)

type TakeProfitPrice struct {
	takeProfitPrice float64
	amountToSell    float64
}

func NewTakeProfitPrice(tp float64, sellAmount float64) *TakeProfitPrice {
	return &TakeProfitPrice{takeProfitPrice: tp, amountToSell: sellAmount}
}

func NewTakeProfitPercentage(entryPrice, percentage float64, sellAmount float64) *TakeProfitPrice {
	targetPrice := entryPrice * (1 + percentage)
	return &TakeProfitPrice{takeProfitPrice: targetPrice, amountToSell: sellAmount}

}

func NewTakeProfitMarketCap(tpMcap float64, solPrice float64, totalSupply float64, sellAmount float64) *TakeProfitPrice {
	marketCapSol := tpMcap / solPrice
	tokenPrice := marketCapSol / totalSupply
	return &TakeProfitPrice{tokenPrice, sellAmount}
}

func (s *TakeProfitPrice) CheckIfPositionHasHit(p *position.PositionMessage) bool {
	return p.CurrentPrice.Cmp(big.NewFloat(s.takeProfitPrice)) != -1
}

func (s *TakeProfitPrice) SellAmount() float64 {
	return s.amountToSell
}
