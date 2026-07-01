import { useAddTask, useCreateAndRunTask, useDeleteTask, useTransitionTask } from './useTasks';
import type { RowActions } from '../types/rowActions';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import {
  useAddStrategy,
  useDeleteStrategy,
  useOpenTerminal,
  useStartStrategy,
  useStopStrategy,
} from './useStrategy';
import { toast } from 'sonner';
import { mapTaskToTaskPost } from '../mapper/taskMapper';
import type { TaskPostDto } from '../types/task';

export function useRowActions(setEditingRow: (row: DisplayRow | null) => void): RowActions {
  const deleteMutation = useDeleteTask();
  const duplicateMutation = useAddTask();
  const transitionMutation = useTransitionTask();
  const startStrategy = useStartStrategy();
  const stopStrategy = useStopStrategy();
  const createAndRunTaskMutation = useCreateAndRunTask();

  const deleteStrategyMutation = useDeleteStrategy();
  const duplicateStrategyMutation = useAddStrategy();
  const openTerminal = useOpenTerminal();

  return {
    onStart: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        transitionMutation.mutate(
          { id: row.id, taskAction: 'Run' },
          {
            onError: (e) =>
              toast.error(`Failed to start task ${row.id}`, { description: e.message }),
          }
        );
      } else {
        startStrategy.mutate(row.id, {
          onError: (e) => toast.error(`Failed to start task ${row.id}`, { description: e.message }),
        });
      }
    },
    onStop: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        transitionMutation.mutate(
          { id: row.id, taskAction: 'Stop' },
          {
            onError: (e) =>
              toast.error(`Failed to stop task ${row.id}`, { description: e.message }),
          }
        );
      } else {
        stopStrategy.mutate(row.id, {
          onError: (e) => toast.error(`Failed to stop task ${row.id}`, { description: e.message }),
        });
      }
    },
    onEdit: (row: DisplayRow) => {
      setEditingRow(row);
    },
    onDelete: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        deleteMutation.mutate(row.id, {
          onError: (e) =>
            toast.error(`Failed to delete task ${row.id}`, { description: e.message }),
        });
      } else {
        deleteStrategyMutation.mutate(row.id, {
          onError: (e) =>
            toast.error(`Failed to delete task ${row.id}`, { description: e.message }),
        });
      }
    },
    onDuplicate: (row: DisplayRow) => {
      if (row.type === TaskRowType.Task) {
        duplicateMutation.mutate(mapTaskToTaskPost(row.data), {
          onError: (e) =>
            toast.error(`Failed to duplicate task ${row.id}`, { description: e.message }),
        });
      } else {
        duplicateStrategyMutation.mutate(row.data, {
          onError: (e) =>
            toast.error(`Failed to duplicate task ${row.id}`, { description: e.message }),
        });
      }
    },
    onQuickSell: (task: TaskPostDto) => {
      createAndRunTaskMutation.mutate(task, {
        onError: (e) =>
          toast.error(`Failed to create task ${task.type}`, { description: e.message }),
      });
    },

    onOpenTerminal: (row: DisplayRow) => {
      openTerminal.mutate(row.id, {
        onError: (e) =>
          toast.error(`Failed to open terminal ${row.id}`, { description: e.message }),
      });
    },
  };
}
