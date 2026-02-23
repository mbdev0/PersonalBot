export interface TaskWSMessage {
  task_event: TaskEvent;
  error?: string;
}

interface TaskEvent {
  task_id: number;
  strategy_id: number;
  state: TaskState;
  time: string;
  message: string;
  event_type: 'StatusUpdate' | 'ProgressMessage';
}

interface TaskState {
  task_state: string;
  error: string;
}

export interface SendTaskWSMessage {
  type: 'Subscribe' | 'Unsubscribe';
  id: number;
}
