package sell

import (
	"personal_bot/backend/internal/core/position"
)

type Strategy interface {
	CheckIfPositionHasHit(p *position.PositionMessage) bool
	SellAmount() float64
}
