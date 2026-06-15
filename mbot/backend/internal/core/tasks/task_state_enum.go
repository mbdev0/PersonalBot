package tasks

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

func (ts TaskState) IsSuccess() bool {
	return ts == TaskDone
}

func (ts TaskState) IsCreate() bool {
	return ts == TaskCreate
}

func (ts TaskState) IsCancel() bool {
	return ts == TaskCancel
}

func (ts TaskState) IsAbleToRun() bool {
	return ts == TaskCreate || ts.IsRetryable()
}
