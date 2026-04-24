import { useMemo, useState, useSyncExternalStore } from 'react';
import { useTaskDashboard } from '../hooks/useTaskDashboard';
import { type DisplayRow } from '../types/tableRows';
import { BotDialog } from '../../../components/botDialog';
import { TaskUpdate } from './taskUpdate';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { mapDashboardRowToRow } from '../mapper/rowMapper';
import { DashboardTable } from '@/components/dashboardTable';
import { key, positionStore } from '@/stores/positionStore';

export function TaskTable() {
  const { isPending, isError, data, error } = useTaskDashboard();
  const [editingRow, setEditingRow] = useState<DisplayRow | null>(null);
  const positionVersion = useSyncExternalStore(positionStore.subscribe, positionStore.getSnapshot);

  const rows = useMemo(() => {
    return (
      data?.rows.map((row) => {
        void positionVersion;
        if (row.type !== 'strategy') return mapDashboardRowToRow(row);

        if (row.data.trading_type === 'BUY') {
          return mapDashboardRowToRow({
            ...row,
            tx_message: row.tx_message ?? row.ws_message,
            ws_message: positionStore.get(key(row.id, row.data.buy_task_id)) ?? row.ws_message,
          });
        }

        if (row.data.trading_type === 'SELL') {
          return mapDashboardRowToRow({
            ...row,
            tx_message: row.tx_message ?? row.ws_message,
            ws_message: positionStore.get(key(row.id, row.data.sell_task_id)) ?? row.ws_message,
          });
        }

        return mapDashboardRowToRow({
          ...row,
          children: row.children.map((child) => ({
            ...child,
            tx_message: child.tx_message ?? child.ws_message,
            ws_message: positionStore.get(key(row.id, child.id)) ?? child.ws_message,
          })),
        });
      }) ?? []
    );
  }, [data?.rows, positionVersion]);

  if (isPending) {
    return <div className="loading_tasks">Loading Tasks...</div>;
  }

  if (isError) {
    return <div className="loading_task_error">Error whilst loading tasks: {error.message}</div>;
  }

  return (
    <div className="table-container">
      <DashboardTable data={rows} setEditingRow={setEditingRow} />

      <BotDialog isOpen={!!editingRow} onClose={() => setEditingRow(null)}>
        <DialogHeader>
          <DialogTitle className="font-semibold">Edit Task</DialogTitle>
        </DialogHeader>
        {editingRow && (
          <TaskUpdate row={editingRow} onClose={() => setEditingRow(null)}></TaskUpdate>
        )}
      </BotDialog>
    </div>
  );
}
