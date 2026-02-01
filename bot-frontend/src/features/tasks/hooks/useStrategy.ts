import { useMutation, useQueryClient } from '@tanstack/react-query';
import {
  deleteStrategy,
  postStrategy,
  putStrategy,
  startStrategy,
  stopStrategy,
} from '../api/strategyTasks';
import type { StrategyTaskPost } from '../types/strategies/strategyTaskPost';
import type { StrategyTaskPut } from '../types/strategies/strategyTaskPut';

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

export function useStartStrategy() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => startStrategy(id),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError(e) {
      console.error('error whilst starting strategy: ', e);
    },
  });
}

export function useStopStrategy() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => stopStrategy(id),
    onSuccess() {
      client.invalidateQueries({ queryKey: ['dashboard'] });
    },
    onError(e) {
      console.error('error whilst stopping strategy: ', e);
    },
  });
}
