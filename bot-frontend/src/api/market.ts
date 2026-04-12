import { API_BASE } from '@/config/urls';
import type { Market } from '@/types/market';

export async function GetSolanaPrice() {
  const url = API_BASE + '/market/sol-price';
  const resp = await fetch(url);
  if (!resp.ok) {
    console.error('failure to get sol price');
  }

  const data: Market = await resp.json();
  return data;
}
