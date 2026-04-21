import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteTask, postTask, transitionTask } from '../api/tasks';
import type { TaskAction, TaskPost } from '../types/task';
import { toast } from 'sonner';
import type { Dashboard } from '../types/dashboard';

export function useAddTask() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (task: TaskPost) => postTask(task),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError(e) {
      console.error('error: ', e);
    },
  });
}

export function useDeleteTask() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => deleteTask(id),
    onSuccess(_, id) {
      client.setQueryData(['dashboard'], (old: Dashboard | undefined) => {
        if (!old) return old;
        return { ...old, rows: old.rows.filter((row) => row.id !== id) };
      });
    },
    onError: (e) => toast.error('Failure to delete task', { description: e.message }),
  });
}

export function useTransitionTask() {
  return useMutation({
    mutationFn: ({ id, taskAction }: { id: number; taskAction: TaskAction }) =>
      transitionTask(id, taskAction),
    onError: (e) => toast.error('Failure to transition task', { description: e.message }),
  });
}
