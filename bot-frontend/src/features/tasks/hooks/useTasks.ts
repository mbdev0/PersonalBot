import { useMutation, useQueryClient } from '@tanstack/react-query';
import { deleteTask, postTask, transitionTask } from '../api/tasks';
import type { TaskAction, TaskPost } from '../types/task';

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
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useTransitionTask() {
  return useMutation({
    mutationFn: ({ id, taskAction }: { id: number; taskAction: TaskAction }) =>
      transitionTask(id, taskAction),
  });
}
