import { useMemo } from 'react';
import { useAddTask, useDeleteTask, useTransitionTask } from './useTasks';
import type { RowActions } from '../types/rowActions';
import { TaskType } from '../types/task';
import { TaskRowType, type Row } from '../types/tableRows';
import {
  useAddStrategy,
  useDeleteStrategy,
  useStartStrategy,
  useStopStrategy,
} from './useStrategy';

export function useRowActions(setEditingRow: (row: Row | null) => void): RowActions {
  const deleteMutation = useDeleteTask();
  const duplicateMutation = useAddTask();
  const transitionMutation = useTransitionTask();
  const startStrategy = useStartStrategy();
  const stopStrategy = useStopStrategy();

  const deleteStrategyMutation = useDeleteStrategy();
  const duplicateStrategyMutation = useAddStrategy();

  return useMemo(
    () => ({
      onStart: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          transitionMutation.mutate({ id: row.id, taskAction: TaskType.task_run });
        } else {
          startStrategy.mutate(row.id);
        }
      },
      onStop: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          transitionMutation.mutate({ id: row.id, taskAction: TaskType.task_cancel });
        } else {
          stopStrategy.mutate(row.id);
        }
      },
      onEdit: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          setEditingRow(row);
        } else {
          setEditingRow(row);
        }
      },
      onDelete: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          deleteMutation.mutate(row.id);
        } else {
          deleteStrategyMutation.mutate(row.id);
          console.log('strategy task delete not implemented');
        }
      },
      onDuplicate: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          duplicateMutation.mutate(row.data);
        } else {
          duplicateStrategyMutation.mutate(row.data);
          console.log('strategy task start not implemented');
        }
      },
    }),
    []
  );
}
