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
  wsMessage: string;
  data: StrategyTask;
  subRows: TaskRow[];
}

interface TaskRow {
  id: number;
  type: TaskRowType.Task;
  state: string;
  wsMessage: string;
  data: Task;
  strategyId?: number;
  subRows?: never;
}

type DisplayRow = StrategyRow | TaskRow;

export { type DisplayRow, type StrategyRow, type TaskRow, TaskRowType };
