import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getSettings, postSettings, testDiscordWebhook } from '../api/settings';
import type { Settings } from '../types/settings';

export function usePostSettings() {
  const client = useQueryClient();

  return useMutation({
    mutationFn: (settings: Settings) => postSettings(settings),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ['settings'] });
    },
  });
}

export function useSettings() {
  return useQuery({
    queryFn: () => getSettings(),
    queryKey: ['settings'],
  });
}

export function useTestDiscordWebhook() {
  return useMutation({
    mutationFn: (discordWebhook: string) => testDiscordWebhook(discordWebhook),
    onError: () => {
      console.error('unable to send to discord webhook');
    },
  });
}
