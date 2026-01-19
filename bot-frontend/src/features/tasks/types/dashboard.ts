import type { StrategyTask, StrategyTaskDto } from './strategyTask';
import type { Task, TaskDto } from './task';

interface DashboardDto {
  strategies: StrategyTaskDto[];
  tasksByStrategy: TasksByStrategyDto;
  manualTasks: TaskDto[];
}

interface TasksByStrategyDto {
  [index: string]: TaskDto[];
}

interface Dashboard {
  strategies: StrategyTask[];
  tasksByStrategy: TasksByStrategy;
  manualTasks: Task[];
}

interface TasksByStrategy {
  [index: string]: Task[];
}

export { type DashboardDto, type Dashboard };
