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
    program: src.program,
    type: 'strategy',
    id: src.id,
    ws_message: src.ws_message,
    state: src.state,
    data: src.data,
    children: src.children.map(mapTaskRowDtoToTaskRow),
  };
}

function mapTaskRowDtoToTaskRow(src: TaskDashboardRowDto): TaskDashboardRow {
  return {
    program: src.program,
    type: 'task',
    id: src.id,
    ws_message: src.ws_message,
    state: src.state,
    data: mapTaskDtoToTask(src.data),
    children: src.children?.map(mapTaskRowDtoToTaskRow),
  };
}
