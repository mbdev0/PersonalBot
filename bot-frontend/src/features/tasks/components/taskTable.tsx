import { useState } from 'react';
import { useTaskDashboard } from '../hooks/useTaskDashboard';
import { type DisplayRow } from '../types/tableRows';
import { BotDialog } from '../../../components/botDialog';
import { TaskUpdate } from './taskUpdate';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { mapDashboardRowToRow } from '../mapper/rowMapper';
import { DashboardTable } from '@/components/dashboardTable';

export function TaskTable() {
  const { isPending, isError, data, error } = useTaskDashboard();
  const [editingRow, setEditingRow] = useState<DisplayRow | null>(null);

  if (isPending) {
    return <div className="loading_tasks">Loading Tasks...</div>;
  }

  if (isError) {
    return <div className="loading_task_error">Error whilst loading tasks: {error.message}</div>;
  }

  const rows = data.rows.map(mapDashboardRowToRow) ?? [];

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
