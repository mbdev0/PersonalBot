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
import { useStrategyWebsocketSend } from '@/hooks/useWebsocketSend';
import { toast } from 'sonner';

export function useAddStrategy() {
  const client = useQueryClient();
  const strategySend = useStrategyWebsocketSend();

  return useMutation({
    mutationFn: (task: StrategyTaskPost) => postStrategy(task),
    onSuccess: (response) => {
      client.invalidateQueries({ queryKey: ['dashboard'] });
      strategySend({ type: 'Subscribe', id: response.id });
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
    onError: (e) => {
      toast.error(`Failure to delete strategy task: ${e.message}`);
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
      toast.error(`Failure to start strategy task: ${e.message}`);
    },
  });
}

export function useStopStrategy() {
  return useMutation({
    mutationFn: (id: number) => stopStrategy(id),
    onSuccess() {},
    onError(e) {
      console.error('error whilst stopping strategy: ', e);
      toast.error(`Failure to stop strategy task: ${e.message}`);
    },
  });
}
