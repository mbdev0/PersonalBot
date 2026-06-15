import type { StrategyTask } from './strategies/strategyTask';
import type { Task } from './task';

enum TaskRowType {
  Task = 'task',
  Strategy = 'strategy',
}

interface StrategyRow {
  id: number;
  program: string;
  type: TaskRowType.Strategy;
  state: string;
  ws_message: string;
  tx_message?: string;
  data: StrategyTask;
  subRows: TaskRow[];
}

export interface TaskRow {
  id: number;
  program: string;
  type: TaskRowType.Task;
  state: string;
  ws_message: string;
  tx_message?: string;
  data: Task;
  strategyId?: number;
  subRows?: TaskRow[];
}

type DisplayRow = StrategyRow | TaskRow;

export { type DisplayRow, TaskRowType };
