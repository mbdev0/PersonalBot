package tasks

type TaskState int

const (
	TaskStateCreated TaskState = iota
	TaskStateRunning
	TaskStateTransactionSent
	TaskStateCompleted
	TaskStateTransactionFailed
	TaskStateFailed
	TaskStateValidationFailed
	TaskStateCancelled
	TaskStateTimeout
	TaskStateUnknown
)

func (s TaskState) ToString() string {
	switch s {
	case TaskStateCreated: //done
		return "Created"
	case TaskStateRunning: //done
		return "Running"
	case TaskStateTransactionSent: //done
		return "Transaction Sent"
	case TaskStateCompleted: //done
		return "Completed"
	case TaskStateTransactionFailed: //done
		return "Transaction Failed"
	case TaskStateFailed: //done
		return "Failed"
	case TaskStateValidationFailed: //done
		return "Validation Failed"
	case TaskStateCancelled: //done
		return "Cancelled"
	case TaskStateTimeout: //done
		return "Timeout"
	default:
		return "Unknown"
	}
}

var StateTransitions map[TaskState][]TaskState = map[TaskState][]TaskState{
	TaskStateCreated: {
		TaskStateRunning, TaskStateValidationFailed,
	},
	TaskStateRunning: {
		TaskStateTransactionSent, TaskStateTransactionFailed, TaskStateTimeout, TaskStateCancelled, TaskStateUnknown,
	},
	TaskStateTransactionSent: {
		TaskStateCompleted, TaskStateFailed, TaskStateTimeout, TaskStateCancelled,
	},
	TaskStateTransactionFailed: {
		TaskStateFailed,
	},
	TaskStateUnknown: {
		TaskStateFailed,
	},
	TaskStateFailed: {
		TaskStateCreated, TaskStateCancelled,
	},
}

type State struct {
	TaskState TaskState
	Error     string
}

func (s *State) SetState(state TaskState) {
	s.TaskState = state
}
