import { useState } from 'react';
import { useTaskDashboard } from '../hooks/useTaskDashboard';
import type { Dashboard } from '../types/dashboard';
import { TaskRowType, type Row } from '../types/tableRows';
import { UnifiedTaskRow } from './tableRow';
import { useRowActions } from '../hooks/useRowActions';
import { Table, TableBody, TableHead, TableHeader, TableRow } from '@/components/ui/table';
import { BotDialog } from '../../../components/botDialog';
import { TaskUpdate } from './taskUpdate';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';

export function TaskTable() {
  const { isPending, isError, data, error } = useTaskDashboard();
  //TODO: this will probably require a data table to allow expanding of children from tanstack see https://ui.shadcn.com/docs/components/radix/data-table
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());
  const [editingRow, setEditingRow] = useState<Row | null>(null);
  const actions = useRowActions(setEditingRow);

  if (isPending) {
    return <div className="loading_tasks">Loading Tasks...</div>;
  }

  if (isError) {
    return <div className="loading_task_error">Error whilst loading tasks: {error.message}</div>;
  }

  const rows = buildRows(data, expandedIds);

  return (
    <div className="bg-background text-foreground rounded-lg shadow-sm p-6">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="font-extrabold">Task Name</TableHead>
            <TableHead className="font-extrabold">Task Type</TableHead>
            <TableHead className="font-extrabold">Message</TableHead>
            <TableHead className="font-extrabold">Status</TableHead>
            <TableHead className="font-extrabold">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((task) => (
            <UnifiedTaskRow key={task.id} tableRow={task} rowActions={actions} />
          ))}
        </TableBody>
      </Table>

      <BotDialog isOpen={!!editingRow} onClose={() => setEditingRow(null)}>
        <DialogHeader>
          <DialogTitle>Edit Task</DialogTitle>
        </DialogHeader>
        {editingRow && (
          <TaskUpdate row={editingRow} onClose={() => setEditingRow(null)}></TaskUpdate>
        )}
      </BotDialog>
    </div>
  );
}

function buildRows(data: Dashboard, expandedIds: Set<number>): Row[] {
  const rows: Row[] = [];
  data.rows.forEach((r) => {
    if (r.type == 'strategy') {
      rows.push({
        id: r.id,
        type: TaskRowType.Strategy,
        depth: 0,
        data: r.data,
      });

      if (r.children.length > 0) {
        r.children.forEach((c) => {
          rows.push({
            id: c.id,
            type: TaskRowType.Task,
            depth: 1,
            data: c.data,
          });
        });
      }
    }
  });

  return rows;
}
