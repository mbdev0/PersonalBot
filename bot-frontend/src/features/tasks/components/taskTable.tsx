import { useState } from 'react';
import { type Task } from '../types/task';
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
  const [expandedIds, setExpandedIds] = useState<Set<number>>(new Set());
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const actions = useRowActions(setEditingTask);

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

      <BotDialog isOpen={!!editingTask} onClose={() => setEditingTask(null)}>
        <DialogHeader>
          <DialogTitle>Edit Task</DialogTitle>
        </DialogHeader>
        {editingTask && (
          <TaskUpdate task={editingTask} onClose={() => setEditingTask(null)}></TaskUpdate>
        )}
      </BotDialog>
    </div>
  );
}

function buildRows(data: Dashboard, expandedIds: Set<number>): Row[] {
  const rows: Row[] = [];
  data.strategies.forEach((s) => {
    rows.push({
      id: s.id,
      type: TaskRowType.Strategy,
      depth: 0,
      data: s,
    });

    if (expandedIds.has(s.id)) {
      const tasks = data.tasksByStrategy[s.id] || [];
      tasks.forEach((t) => {
        rows.push({
          id: t.task_id,
          strategyId: s.id,
          type: TaskRowType.Task,
          depth: 1,
          data: t,
        });
      });
    }
  });

  data.manualTasks.forEach((t) => {
    rows.push({
      id: t.task_id,
      type: TaskRowType.Task,
      depth: 0,
      data: t,
    });
  });

  return rows;
}
