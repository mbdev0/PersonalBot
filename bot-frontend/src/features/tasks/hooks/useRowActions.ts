import { useAddTask, useDeleteTask, useTransitionTask } from './useTasks';
import type { RowActions } from '../types/rowActions';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import {
  useAddStrategy,
  useDeleteStrategy,
  useStartStrategy,
  useStopStrategy,
} from './useStrategy';

export function useRowActions(setEditingRow: (row: DisplayRow | null) => void): RowActions {
  const deleteMutation = useDeleteTask();
  const duplicateMutation = useAddTask();
  const transitionMutation = useTransitionTask();
  const startStrategy = useStartStrategy();
  const stopStrategy = useStopStrategy();

  const deleteStrategyMutation = useDeleteStrategy();
  const duplicateStrategyMutation = useAddStrategy();

  return {
    onStart: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        transitionMutation.mutate({ id: row.id, taskAction: 'Run' });
      } else {
        startStrategy.mutate(row.id);
      }
    },
    onStop: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        transitionMutation.mutate({ id: row.id, taskAction: 'Stop' });
      } else {
        stopStrategy.mutate(row.id);
      }
    },
    onEdit: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        setEditingRow(row);
      } else {
        setEditingRow(row);
      }
    },
    onDelete: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        deleteMutation.mutate(row.id);
      } else {
        deleteStrategyMutation.mutate(row.id);
        console.log('strategy task delete not implemented');
      }
    },
    onDuplicate: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        duplicateMutation.mutate(row.data);
      } else {
        duplicateStrategyMutation.mutate(row.data);
        console.log('strategy task start not implemented');
      }
    },
  };
}
