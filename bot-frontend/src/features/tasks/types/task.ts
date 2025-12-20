interface TaskDto {
  task_id: string;
  type: string;
  slippage: number;
  compute_units: number;
  token_address: string;
  wallet_name: string;
  buy_amount: number | undefined;
  buy_fee: number | undefined;
  sell_amount: number | undefined;
  sell_fee: number | undefined;
  sell_posiiton_id: string | undefined;
  state: StateDto;
}

interface StateDto {
  task_state: string;
  error: string;
}

interface Task {
  task_id: string;
  type: string;
  slippage: number;
  compute_units: number;
  token_address: string;
  wallet_name: string;
  buy_amount: number | undefined;
  buy_fee: number | undefined;
  sell_amount: number | undefined;
  sell_fee: number | undefined;
  sell_posiiton_id: string | undefined;
  state: State;
}

interface TaskPost{
  type: string;
  slippage: number;
  compute_units: number;
  token_address: string;
  wallet_name: string;
  buy_amount?: number;
  buy_fee?: number;
  sell_amount?: number;
  sell_fee?: number;
}

interface TaskPostDto{
  type: string;
  slippage: number;
  compute_units: number;
  token_address: string;
  wallet_name: string;
  buy_amount: number | undefined;
  buy_fee: number | undefined;
  sell_amount: number | undefined;
  sell_fee: number | undefined;
}

interface State {
  task_state: string;
  error: string;
}

enum TaskType {
  task_create = "TaskCreate",
  task_run = "TaskRun",
  tx_instruction_build = "TxInstructionBuild",
  tx_build = "TxBuild",
  tx_send = "TxSend",
  tx_confirm = "TxConfirm",
  task_done = "Done",
  task_cancel = "TaskCancel",
  task_validation_failed = "TaskValidationFailed",
  tx_instruction_build_failed = "TxInstructionBuildFailed",
  tx_build_failed = "TxBuildFailed",
  tx_send_failed = "TxSendFailed",
  tx_failed = "TxFailed",
}

export { type TaskDto, type Task, TaskType, type TaskPostDto, type TaskPost };
