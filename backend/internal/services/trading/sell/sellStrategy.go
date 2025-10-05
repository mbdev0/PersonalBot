package sell

import (
	"pump_fun/internal/core/position"
)

type Strategy interface {
	CheckIfPositionHasHit(p *position.PositionMessage) bool
	SellAmount() float64
}
