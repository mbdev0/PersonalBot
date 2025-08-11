package transaction

import "pump_fun/internal/core/tasks"

type Transaction interface {
	BuildInstructions() error
	BuildTransaction() error
	SendTransaction() error
	ConfirmTransaction() error
	GetTask() tasks.Task
}
