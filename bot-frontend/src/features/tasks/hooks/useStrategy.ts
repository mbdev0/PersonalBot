import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { StrategyTaskPost, StrategyTaskPut } from '../types/strategyTask';
import {
  deleteStrategy,
  postStrategy,
  putStrategy,
  startStrategy,
  stopStrategy,
} from '../api/strategyTasks';

export function useAddStrategy() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (task: StrategyTaskPost) => postStrategy(task),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError(e) {
      console.error('error: ', e);
    },
  });
}

export function useUpdateStrategy() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (task: StrategyTaskPut) => putStrategy(task),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useDeleteStrategy() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => deleteStrategy(id),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
  });
}

export function useStartStrategy({ id }: { id: number }) {
  return useQuery({
    queryKey: ['startStrategy', id],
    queryFn: () => startStrategy(id),
  });
}

export function useStopStrategy({ id }: { id: number }) {
  return useQuery({
    queryKey: ['stopStrategy', id],
    queryFn: () => stopStrategy(id),
  });
}
