import type { StrategyTask } from './strategies/strategyTask';
import type { Task, TaskDto } from './task';

interface BaseDashboardRow {
  program: string;
  type: string;
  id: number;
  ws_message: string;
  tx_message?: string;
  state: string;
}

export interface StrategyDashboardRow extends BaseDashboardRow {
  type: 'strategy';
  data: StrategyTask;
  children: TaskDashboardRow[];
}

export interface TaskDashboardRow extends BaseDashboardRow {
  type: 'task';
  data: Task;
  children?: TaskDashboardRow[];
}

export type DashboardRow = StrategyDashboardRow | TaskDashboardRow;

export interface StrategyDashboardRowDto {
  program: string;
  type: 'strategy';
  id: number;
  ws_message: string;
  state: string;
  data: StrategyTask;
  children: TaskDashboardRowDto[];
}

export interface TaskDashboardRowDto {
  program: string;
  type: 'manual';
  id: number;
  ws_message: string;
  state: string;
  data: TaskDto;
  children?: TaskDashboardRowDto[];
}

export type DashboardRowDto = StrategyDashboardRowDto | TaskDashboardRowDto;

interface Dashboard {
  rows: DashboardRow[];
}

interface DashboardDto {
  rows: DashboardRowDto[];
}

export { type DashboardDto, type Dashboard };
