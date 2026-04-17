import { API_BASE } from '@/config/urls';
import type { Market } from '@/types/market';
import { toast } from 'sonner';

export async function GetSolanaPrice() {
  const url = API_BASE + '/market/sol-price';
  const resp = await fetch(url);

  if (!resp.ok) {
    const text = await resp.text();
    toast.error(text || `Request failed (${resp.status})`);
    return;
  }

  const data: Market = await resp.json();
  return data;
}
