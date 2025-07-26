package tasks

import "fmt"

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

func (ts TaskState) ToString() string {
	switch ts {
	case TaskStateCreated:
		return "Created"
	case TaskStateRunning:
		return "Running"
	case TaskStateTransactionSent:
		return "TransactionSent"
	case TaskStateCompleted:
		return "Completed"
	case TaskStateTransactionFailed:
		return "TransactionFailed"
	case TaskStateFailed:
		return "Failed"
	case TaskStateValidationFailed:
		return "ValidationFailed"
	case TaskStateCancelled:
		return "Cancelled"
	case TaskStateTimeout:
		return "TimedOut"
	case TaskStateUnknown:
		return "Unknown"
	}

	return "Invalid State"
}

func ParseStateString(in string) (TaskState, error) {
	switch in {
	case TaskStateCreated.ToString():
		return TaskStateCreated, nil
	case TaskStateRunning.ToString():
		return TaskStateRunning, nil
	case TaskStateTransactionSent.ToString():
		return TaskStateTransactionSent, nil
	case TaskStateCompleted.ToString():
		return TaskStateCompleted, nil
	case TaskStateTransactionFailed.ToString():
		return TaskStateTransactionFailed, nil
	case TaskStateFailed.ToString():
		return TaskStateFailed, nil
	case TaskStateValidationFailed.ToString():
		return TaskStateValidationFailed, nil
	case TaskStateCancelled.ToString():
		return TaskStateCancelled, nil
	case TaskStateTimeout.ToString():
		return TaskStateTimeout, nil
	case TaskStateUnknown.ToString():
		return TaskStateUnknown, nil
	}

	return TaskStateUnknown, fmt.Errorf("invalid task state for: %v", in)
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
