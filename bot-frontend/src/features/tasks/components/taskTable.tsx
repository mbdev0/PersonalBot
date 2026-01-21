import { useState } from 'react';
import Modal from '../../../components/modal';
import { TaskUpdate } from './taskUpdate';
import { type Task } from '../types/task';
import './taskTable.css';
import { useTaskDashboard } from '../hooks/useTaskDashboard';
import type { Dashboard } from '../types/dashboard';
import { TaskRowType, type Row } from '../types/tableRows';
import { UnifiedTaskRow } from './tableRow';
import { useRowActions } from '../hooks/useRowActions';
import { Table, TableBody, TableHead, TableHeader, TableRow } from '@/components/ui/table';

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
    <div className="bg-slate-950 rounded-lg border shadow-sm p-6">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead className="text-slate-50">Task Name</TableHead>
            <TableHead className="text-slate-50">Task Type</TableHead>
            <TableHead className="text-slate-50">Message</TableHead>
            <TableHead className="text-slate-50">Status</TableHead>
            <TableHead className="text-slate-50">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((task) => (
            <UnifiedTaskRow key={task.id} tableRow={task} rowActions={actions} />
          ))}
        </TableBody>
      </Table>
      {/* <table>
        <thead>
          <tr>
            <th scope="column">Task Name</th>
            <th scope="column">Task Type</th>
            <th scope="column">Message</th>
            <th scope="column">Status</th>
            <th scope="column">Actions</th>
          </tr>
        </thead>
        <tbody>
          {rows.map((task) => (
            <TableRow key={task.id} tableRow={task} rowActions={actions} />
          ))}
        </tbody>
      </table>

      {editingTask && (
        <Modal isOpen={!!editingTask} onClose={() => setEditingTask(null)}>
          <TaskUpdate task={editingTask} onClose={() => setEditingTask(null)} />
        </Modal>
      )} */}
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
