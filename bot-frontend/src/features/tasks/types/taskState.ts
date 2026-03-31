export const TaskState = {
  CREATE: 'Create',
  VALIDATING: 'Validating',
  VALIDATION_FAIL: 'Validation Failed',
  TX_INSTRUCTION_BUILD: 'Building Instructions',
  TX_INSTRUCTION_BUILD_FAIL: 'Instruction Build Failed',
  TX_BUILD: 'Building Transaction',
  TX_BUILD_FAIL: 'Transaction Build Failed',
  TX_SEND: 'Sending Transaction',
  TX_SEND_FAIL: 'Transaction Send Failed',
  TX_CONFIRM: 'Confirming Transaction',
  TX_CONFIRM_FAIL: 'Transaction Confirmation Failed',
  TASK_UPDATING_POSITION: 'Updating Position',
  TASK_UPDATING_POSITION_FAIL: 'Position Update Failed',
  TASK_DONE: 'Done',
  TASK_CANCEL: 'Cancelled',
  TASK_FAIL: 'Failed',
  TASK_UNKNOWN: 'Unknown',
} as const;

export type TaskState = (typeof TaskState)[keyof typeof TaskState];

const RUNNING_STATES: Set<string> = new Set([
  TaskState.VALIDATING,
  TaskState.TX_INSTRUCTION_BUILD,
  TaskState.TX_BUILD,
  TaskState.TX_SEND,
  TaskState.TX_CONFIRM,
  TaskState.TASK_UPDATING_POSITION,
]);

const FAILURE_STATES: Set<string> = new Set([
  TaskState.VALIDATION_FAIL,
  TaskState.TX_INSTRUCTION_BUILD_FAIL,
  TaskState.TX_BUILD_FAIL,
  TaskState.TX_SEND_FAIL,
  TaskState.TX_CONFIRM_FAIL,
  TaskState.TASK_UPDATING_POSITION_FAIL,
  TaskState.TASK_FAIL,
]);

export function isRunning(state: string): boolean {
  return RUNNING_STATES.has(state);
}

export function isFailure(state: string): boolean {
  return FAILURE_STATES.has(state);
}
