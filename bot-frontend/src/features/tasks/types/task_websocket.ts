export interface TaskWSMessage {
  task_event: TaskEvent;
  error?: string;
}

interface TaskEvent {
  task_id: number;
  strategy_id: number;
  state: string;
  time: string;
  message: string;
  event_type: string;
}
