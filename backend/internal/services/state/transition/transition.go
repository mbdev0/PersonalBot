package transition

import (
	"fmt"
	"personal_bot/internal/core/tasks"
	"personal_bot/pkg/logger"
)

var StateTransitions = map[tasks.TaskState]tasks.Transistion{
	//current : next, on error
	tasks.TaskCreate:         {Next: tasks.TaskRun, OnError: tasks.TaskFail},
	tasks.TaskRun:            {Next: tasks.TxInstructionBuild, OnError: tasks.TaskValidationFailed},
	tasks.TxInstructionBuild: {Next: tasks.TxBuild, OnError: tasks.TxInstructionBuildFailed},
	tasks.TxBuild:            {Next: tasks.TxSend, OnError: tasks.TxBuildFailed},
	tasks.TxSend:             {Next: tasks.TxConfirm, OnError: tasks.TxSendFailed},
	tasks.TxConfirm:          {Next: tasks.TaskDone, OnError: tasks.TxFailed},
	tasks.TaskCancel:         {Next: tasks.TaskCreate, OnError: tasks.TaskFail},
}

func IsAbleToTransitionTo(nextState tasks.TaskState, task tasks.Task) bool {
	transition, ok := StateTransitions[task.State().TaskState]
	if !ok {
		return false
	}

	if nextState == tasks.TaskCancel {
		return true
	}

	if transition.Next != nextState {
		return false
	}

	return true
}

func IsRetryableState(state tasks.TaskState) bool {
	switch state {
	case tasks.TaskValidationFailed, tasks.TxInstructionBuildFailed, tasks.TxBuildFailed, tasks.TxSendFailed, tasks.TxFailed, tasks.TaskFail, tasks.TaskCancel:
		return true
	default:
		return false
	}
}

func AutoTransitionTask(task tasks.Task, err error) error {
	transition, ok := StateTransitions[task.State().TaskState]
	if !ok {
		return fmt.Errorf("no next transition ")
	}

	if err != nil {
		logger.Error("we got an error whilst transitioning states")
		task.SetState(tasks.State{TaskState: transition.OnError, Error: err.Error()})
		return nil
	}

	logger.Information("successful transition to: ", transition.Next.ToString())
	task.SetState(tasks.State{TaskState: transition.Next})

	return nil
}
