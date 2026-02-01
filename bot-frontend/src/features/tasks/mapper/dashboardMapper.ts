// features/dashboard/mapper/dashboardMapper.ts
import { mapTaskDtoToTask } from '../../tasks/mapper/taskMapper';
import type {
  DashboardDto,
  Dashboard,
  DashboardRowDto,
  DashboardRow,
  StrategyDashboardRow,
  StrategyDashboardRowDto,
  TaskDashboardRow,
  TaskDashboardRowDto,
} from '../types/dashboard';
import { mapStrategyTaskDtoToStrategyTask } from './strategy/strategyMapper';

export function MapDashboardDtoToDashboard(src: DashboardDto): Dashboard {
  return {
    rows: src.rows.map(mapDashboardRowDtoToDashboardRow),
  };
}

function mapDashboardRowDtoToDashboardRow(src: DashboardRowDto): DashboardRow {
  if (src.type === 'strategy') {
    return mapStrategyRowDtoToStrategyRow(src);
  }

  return mapTaskRowDtoToTaskRow(src);
}

function mapStrategyRowDtoToStrategyRow(src: StrategyDashboardRowDto): StrategyDashboardRow {
  return {
    type: 'strategy',
    id: src.id,
    wsMessage: src.wsMessage,
    state: src.state,
    data: mapStrategyTaskDtoToStrategyTask(src.data),
    children: src.children.map(mapTaskRowDtoToTaskRow),
  };
}

function mapTaskRowDtoToTaskRow(src: TaskDashboardRowDto): TaskDashboardRow {
  return {
    type: 'task',
    id: src.id,
    wsMessage: src.wsMessage,
    state: src.state,
    data: mapTaskDtoToTask(src.data),
  };
}
