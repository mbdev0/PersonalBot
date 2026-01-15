package sell

import (
	"math/big"
	"personal_bot/internal/core/position"
)

type StopLossPrice struct {
	stopLossPrice float64
	amountToSell  float64
}

func NewStopLossPrice(price float64, sellAmount float64) *StopLossPrice {
	return &StopLossPrice{stopLossPrice: price, amountToSell: sellAmount}
}

func NewStopLossPercentage(entryPrice float64, percentage float64, sellAmount float64) *StopLossPrice {
	targetPrice := entryPrice * (1 + percentage)
	return &StopLossPrice{stopLossPrice: targetPrice, amountToSell: sellAmount}
}

func NewStopLossMarketCap(marketCap float64, solPrice float64, totalSupply float64, sellAmount float64) *StopLossPrice {
	marketCapSol := marketCap / solPrice
	tokenPrice := marketCapSol / totalSupply
	return &StopLossPrice{stopLossPrice: tokenPrice, amountToSell: sellAmount}
}

func (s *StopLossPrice) CheckIfPositionHasHit(p *position.PositionMessage) bool {
	return p.CurrentPrice.Cmp(big.NewFloat(s.stopLossPrice)) != 1
}

func (s *StopLossPrice) SellAmount() float64 {
	return s.amountToSell
}
