import type { StrategyTaskDto, StrategyTask } from './strategies/strategyTask';
import type { Task, TaskDto } from './task';

interface DashboardDto {
  rows: DashboardRowDto[];
}

interface BaseDashboardRowDto {
  type: string;
  id: number;
  ws_message: string;
  state: string;
}

export interface StrategyDashboardRowDto extends BaseDashboardRowDto {
  type: 'strategy';
  data: StrategyTaskDto;
  children: TaskDashboardRowDto[];
}

export interface TaskDashboardRowDto extends BaseDashboardRowDto {
  type: 'manual';
  data: TaskDto;
  children?: TaskDashboardRowDto[];
}

type DashboardRowDto = StrategyDashboardRowDto | TaskDashboardRowDto;

interface BaseDashboardRow {
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

interface Dashboard {
  rows: DashboardRow[];
}

export { type DashboardDto, type DashboardRowDto, type Dashboard };
