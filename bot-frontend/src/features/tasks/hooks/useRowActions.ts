import { useMemo } from 'react';
import { useAddTask, useDeleteTask, useTransitionTask } from './useTasks';
import type { RowActions } from '../types/rowActions';
import { TaskType, type Task } from '../types/task';
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

  const deleteStrategyMutation = useDeleteStrategy();
  const duplicateStrategyMutation = useAddStrategy();

  return useMemo(
    () => ({
      onStart: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          transitionMutation.mutate({ id: row.id, taskAction: TaskType.task_run });
        } else {
          const { error } = useStartStrategy({ id: row.id });

          if (error) {
            console.log('error whilst starting task');
          }
        }
      },
      onStop: (row: Row) => {
        if (row.type === TaskRowType.Task) {
          transitionMutation.mutate({ id: row.id, taskAction: TaskType.task_cancel });
        } else {
          const { error } = useStopStrategy({ id: row.id });

          if (error) {
            console.log('error whilst starting task');
          }
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
