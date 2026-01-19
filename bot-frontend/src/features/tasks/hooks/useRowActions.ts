import { useMemo } from 'react';
import { useAddTask, useDeleteTask, useTransitionTask } from './useTasks';
import type { RowActions } from '../types/rowActions';
import { TaskType, type Task } from '../types/task';
import { TaskRowType, type Row } from '../types/tableRows';

export function useRowActions(setEditingWallet: (task: Task | null) => void): RowActions {
  const deleteMutation = useDeleteTask();
  const duplicateMutation = useAddTask();
  const transitionMutation = useTransitionTask();

  return useMemo(
    () => ({
      onStart: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          transitionMutation.mutate({ id: row.id, taskAction: TaskType.task_run });
        } else {
          console.log('strategy task start not implemented');
        }
      },
      onStop: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          transitionMutation.mutate({ id: row.id, taskAction: TaskType.task_cancel });
        } else {
          console.log('strategy task stop not implemented');
        }
      },
      onEdit: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          setEditingWallet(row.data);
        } else {
          console.log('strategy task edit not implemented');
        }
      },
      onDelete: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          deleteMutation.mutate(row.id);
        } else {
          console.log('strategy task delete not implemented');
        }
      },
      onDuplicate: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          duplicateMutation.mutate(row.data);
        } else {
          console.log('strategy task start not implemented');
        }
      },
    }),
    []
  );
}
