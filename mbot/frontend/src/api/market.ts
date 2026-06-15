import { API_BASE } from '@/config/urls';
import type { Market } from '@/types/market';
import { toast } from 'sonner';

const MAX_ATTEMPTS = 5;

export async function GetSolanaPrice() {
  let delay = 10000;

  for (let attempt = 1; attempt <= MAX_ATTEMPTS; attempt++) {
    try {
      const resp = await fetch(API_BASE + '/market/sol-price');
      if (!resp.ok) {
        const text = await resp.text();
        throw new Error(text || `Request failed (${resp.status})`);
      }
      return (await resp.json()) as Market;
    } catch (e) {
      const error = e instanceof Error ? e : new Error('Unknown error');
      toast.error('Failed to fetch SOL price', { description: error.message });
      if (attempt < MAX_ATTEMPTS) {
        await new Promise((resolve) => setTimeout(resolve, delay));
        delay *= 2;
      } else {
        throw error;
      }
    }
  }
}
