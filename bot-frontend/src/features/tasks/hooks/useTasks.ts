import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { deleteTask, getTasks, postTask, putTask, transitionTask } from '../api/tasks';
import type { Task, TaskAction, TaskPost, TaskPut } from '../types/task';

export function useTasks() {
  return useQuery({
    queryKey: ['tasks'],
    queryFn: getTasks,
  });
}

export function useAddTask() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (task: TaskPost) => postTask(task),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['tasks'] });
    },
    onError(e) {
      console.error('error: ', e);
    },
  });
}

export function useUpdateTask() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (task: TaskPut) => putTask(task),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['tasks'] });
    },
  });
}

export function useDeleteTask() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => deleteTask(id),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['tasks'] });
    },
  });
}

export function useTransitionTask() {
  return useMutation({
    mutationFn: ({ id, taskAction }: { id: number; taskAction: TaskAction }) =>
      transitionTask(id, taskAction),
  });
}
