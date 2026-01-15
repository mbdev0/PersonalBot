package strategies

import "personal_bot/internal/core/tasks"

type StrategyMessage struct {
	Id    int64      `json:"id"`
	Event string     `json:"event"`
	Task  tasks.Task `json:"task"`
}
