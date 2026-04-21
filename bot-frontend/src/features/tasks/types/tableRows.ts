import type { StrategyTask } from './strategies/strategyTask';
import type { Task } from './task';

enum TaskRowType {
  Task = 'task',
  Strategy = 'strategy',
}

interface StrategyRow {
  id: number;
  type: TaskRowType.Strategy;
  state: string;
  ws_message: string;
  tx_message?: string;
  data: StrategyTask;
  subRows: TaskRow[];
}

interface TaskRow {
  id: number;
  type: TaskRowType.Task;
  state: string;
  ws_message: string;
  tx_message?: string;
  data: Task;
  strategyId?: number;
  subRows?: never;
}

type DisplayRow = StrategyRow | TaskRow;

export { type DisplayRow, TaskRowType };
