import { API_BASE } from '@/config/urls';
import type { Settings } from '../types/settings';

const settingsUrl = API_BASE + '/settings/';

export async function testDiscordWebhook(discordWebhook: string) {
  const testDiscordWebhookBody = {
    discord_webhook: discordWebhook,
  };

  const resp = await fetch(settingsUrl + 'test_webhook', {
    body: JSON.stringify(testDiscordWebhookBody),
    method: 'POST',
  });

  if (!resp.ok) {
    const text = await resp.text();
    throw new Error(text || `Request failed (${resp.status})`);
  }

  return resp.status;
}

export async function postSettings(settings: Settings) {
  const resp = await fetch(settingsUrl, {
    body: JSON.stringify(settings),
    method: 'POST',
  });

  if (!resp.ok) {
    const text = await resp.text();
    throw new Error(text || `Request failed (${resp.status})`);
  }

  return resp.status;
}

export async function getSettings(): Promise<Settings> {
  const resp = await fetch(settingsUrl, {
    method: 'GET',
  });

  if (!resp.ok) {
    const text = await resp.text();
    throw new Error(text || `Request failed (${resp.status})`);
  }

  const settingsData: Settings = await resp.json();

  return settingsData;
}
