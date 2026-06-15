import type { DashboardRow, TaskDashboardRow } from '../types/dashboard';
import { TaskRowType, type DisplayRow, type TaskRow } from '../types/tableRows';

export function mapDashboardRowToRow(row: DashboardRow): DisplayRow {
  if (row.type === 'strategy') {
    return {
      program: row.program,
      id: row.id,
      type: TaskRowType.Strategy,
      state: row.state,
      ws_message: row.ws_message,
      tx_message: row.tx_message,
      data: row.data,
      subRows: row.children.map((child) => mapTaskRowToRow(child, row.id)),
    };
  }

  return mapTaskRowToRow(row);
}

function mapTaskRowToRow(row: TaskDashboardRow, strategyId?: number): TaskRow {
  return {
    program: row.program,
    id: row.id,
    type: TaskRowType.Task,
    state: row.state,
    ws_message: row.ws_message,
    tx_message: row.tx_message,
    data: row.data,
    strategyId,
    subRows: row.children?.map((child) => mapTaskRowToRow(child)),
  };
}
