import type { DashboardRow } from '../types/dashboard';
import { TaskRowType, type DisplayRow } from '../types/tableRows';

export function mapDashboardRowToRow(row: DashboardRow): DisplayRow {
  if (row.type === 'strategy') {
    return {
      id: row.id,
      type: TaskRowType.Strategy,
      state: row.state,
      wsMessage: row.wsMessage,
      data: row.data,
      subRows: row.children.map((child) => ({
        id: child.id,
        type: TaskRowType.Task as const,
        state: child.state,
        wsMessage: child.wsMessage,
        data: child.data,
        strategyId: row.id,
      })),
    };
  }

  return {
    id: row.id,
    type: TaskRowType.Task,
    state: row.state,
    wsMessage: row.wsMessage,
    data: row.data,
  };
}
