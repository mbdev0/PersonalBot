// internal/core/tasks/state.go
package tasks

import "fmt"

type TaskState string

const (
	TaskCreate TaskState = "TaskCreate"

	TaskValidating     TaskState = "TaskValidating"
	TaskValidationFail TaskState = "TaskValidationFail"

	TxInstructionBuild     TaskState = "TxInstructionBuild"
	TxInstructionBuildFail TaskState = "TxInstructionBuildFail"

	TxBuild     TaskState = "TxBuild"
	TxBuildFail TaskState = "TxBuildFail"

	TxSend     TaskState = "TxSend"
	TxSendFail TaskState = "TxSendFail"

	TxConfirm     TaskState = "TxConfirm"
	TxConfirmFail TaskState = "TxConfirmFail"

	TaskUpdatingPosition     TaskState = "TaskUpdatingPosition"
	TaskUpdatingPositionFail TaskState = "TaskUpdatingPositionFail"

	TaskDone   TaskState = "TaskDone"
	TaskCancel TaskState = "TaskCancel"
	TaskFail   TaskState = "TaskFail"

	TaskUnknown TaskState = "TaskUnknown"
)

func (ts TaskState) ToString() string {
	switch ts {
	case TaskCreate:
		return "Create"
	case TaskValidating:
		return "Validating"
	case TaskValidationFail:
		return "Validation Failed"
	case TxInstructionBuild:
		return "Building Instructions"
	case TxInstructionBuildFail:
		return "Instruction Build Failed"
	case TxBuild:
		return "Building Transaction"
	case TxBuildFail:
		return "Transaction Build Failed"
	case TxSend:
		return "Sending Transaction"
	case TxSendFail:
		return "Transaction Send Failed"
	case TxConfirm:
		return "Confirming Transaction"
	case TxConfirmFail:
		return "Transaction Confirmation Failed"
	case TaskUpdatingPosition:
		return "Updating Position"
	case TaskUpdatingPositionFail:
		return "Position Update Failed"
	case TaskDone:
		return "Done"
	case TaskCancel:
		return "Cancelled"
	case TaskFail:
		return "Failed"
	case TaskUnknown:
		return "Unknown"
	default:
		return "Invalid State"
	}
}

func ParseStateString(state string) (TaskState, error) {
	switch state {
	case "TaskCreate":
		return TaskCreate, nil
	case "TaskValidating":
		return TaskValidating, nil
	case "TaskValidationFail", "TaskValidationFailed":
		return TaskValidationFail, nil
	case "TxInstructionBuild":
		return TxInstructionBuild, nil
	case "TxInstructionBuildFail", "TxInstructionBuildFailed":
		return TxInstructionBuildFail, nil
	case "TxBuild":
		return TxBuild, nil
	case "TxBuildFail", "TxBuildFailed":
		return TxBuildFail, nil
	case "TxSend":
		return TxSend, nil
	case "TxSendFail", "TxSendFailed":
		return TxSendFail, nil
	case "TxConfirm":
		return TxConfirm, nil
	case "TxConfirmFail", "TxFailed":
		return TxConfirmFail, nil
	case "TaskUpdatingPosition":
		return TaskUpdatingPosition, nil
	case "TaskUpdatingPositionFail":
		return TaskUpdatingPositionFail, nil
	case "TaskDone", "Done":
		return TaskDone, nil
	case "TaskCancel":
		return TaskCancel, nil
	case "TaskFail":
		return TaskFail, nil
	default:
		return TaskUnknown, fmt.Errorf("unknown task state: %s", state)
	}
}

func (ts TaskState) IsTerminal() bool {
	return ts == TaskDone || ts == TaskFail
}

func (ts TaskState) IsError() bool {
	return ts == TaskValidationFail ||
		ts == TxInstructionBuildFail ||
		ts == TxBuildFail ||
		ts == TxSendFail ||
		ts == TxConfirmFail ||
		ts == TaskUpdatingPositionFail ||
		ts == TaskFail
}

func (ts TaskState) IsRetryable() bool {
	return ts == TaskCancel || ts.IsError() && !ts.IsTerminal()
}

func (ts TaskState) IsAbleToRun() bool {
	return ts == TaskCreate || ts.IsRetryable()
}
