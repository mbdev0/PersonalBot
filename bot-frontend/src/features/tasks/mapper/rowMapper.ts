import type { DashboardRow } from '../types/dashboard';
import { TaskRowType, type DisplayRow } from '../types/tableRows';

export function mapDashboardRowToRow(row: DashboardRow): DisplayRow {
  if (row.type === 'strategy') {
    return {
      id: row.id,
      type: TaskRowType.Strategy,
      state: row.state,
      ws_message: row.ws_message,
      data: row.data,
      subRows: row.children.map((child) => ({
        id: child.id,
        type: TaskRowType.Task as const,
        state: child.state,
        ws_message: child.ws_message,
        data: child.data,
        strategyId: row.id,
      })),
    };
  }

  return {
    id: row.id,
    type: TaskRowType.Task,
    state: row.state,
    ws_message: row.ws_message,
    data: row.data,
  };
}
