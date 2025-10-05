package tasks

type TaskState int

const (
	TaskCreate TaskState = iota
	TaskRun
	TaskValidationFailed
	TxInstructionBuild
	TxInstructionBuildFailed
	TxBuild
	TxBuildFailed
	TxSend
	TxSendFailed
	TxConfirm
	TxFailed
	TaskDone
	TaskCancel
	TaskUnknown
	TaskFail
)

func (ts TaskState) ToString() string {
	switch ts {
	case TaskCreate:
		return "Create"
	case TaskRun:
		return "Run"
	case TaskValidationFailed:
		return "Task Validation Failed"
	case TxInstructionBuild:
		return "Transaction Instruction Build"
	case TxInstructionBuildFailed:
		return "Transaction Instruction Build Failed"
	case TxBuild:
		return "Transaction Build"
	case TxBuildFailed:
		return "Transaction Build Failed"
	case TxSend:
		return "Transaction Send"
	case TxSendFailed:
		return "Transaction Send Failed"
	case TxConfirm:
		return "Transaction Confirm"
	case TxFailed:
		return "Transaction Failed"
	case TaskDone:
		return "Task Done"
	case TaskCancel:
		return "Task Cancelled"
	case TaskUnknown:
		return "Task Unknown"
	case TaskFail:
		return "Task Failed"
	}

	return "Invalid State"
}

func ParseStateString(state string) (TaskState, error) {
	switch state {
	case "TaskCreate":
		return TaskCreate, nil
	case "TaskRun":
		return TaskRun, nil
	case "TxInstructionBuild":
		return TxInstructionBuild, nil
	case "TxBuild":
		return TxBuild, nil
	case "TxSend":
		return TxSend, nil
	case "TxConfirm":
		return TxConfirm, nil
	case "Done":
		return TaskDone, nil
	case "TaskCancel":
		return TaskCancel, nil
	case "TaskValidationFailed":
		return TaskValidationFailed, nil
	case "TxInstructionBuildFailed":
		return TxInstructionBuildFailed, nil
	case "TxBuildFailed":
		return TxBuildFailed, nil
	case "TxSendFailed":
		return TxSendFailed, nil
	case "TxFailed":
		return TxFailed, nil
	case "TaskFail":
		return TaskFail, nil
	default:
		return TaskUnknown, nil
	}
}
