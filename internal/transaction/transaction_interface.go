package transaction

import "pump_fun/internal/models/tasks"

type Transaction interface {
	BuildInstructions() error
	BuildTransaction() error
	SendTransaction() error
	ConfirmTransaction() error
	GetTask() tasks.Task
}
