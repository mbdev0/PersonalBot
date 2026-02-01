import type { StrategyTask } from './strategies/strategyTask';
import type { Task } from './task';

interface StrategyRow {
  id: number;
  type: TaskRowType.Strategy;
  depth: 0;
  data: StrategyTask;
}

interface TaskRow {
  id: number;
  type: TaskRowType.Task;
  depth: number;
  data: Task;
  strategyId?: number;
}

enum TaskRowType {
  Task = 'task',
  Strategy = 'strategy',
}

type Row = StrategyRow | TaskRow;

export { type Row, TaskRowType };
