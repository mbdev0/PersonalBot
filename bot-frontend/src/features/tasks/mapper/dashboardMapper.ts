import type { Dashboard, DashboardDto } from '../types/dashboard';

export function MapDashboardDtoToDashboard(src: DashboardDto): Dashboard {
  return {
    strategies: src.strategies,
    tasksByStrategy: src.tasksByStrategy,
    manualTasks: src.manualTasks,
  };
}
