import { GetSolanaPrice } from '@/api/market';
import { SOLANA_PRICE_REFETCH } from '@/config/constants';
import { useQuery } from '@tanstack/react-query';

export function useSolanaPrice() {
  return useQuery({
    queryKey: ['solana_price'],
    queryFn: GetSolanaPrice,
    refetchInterval: SOLANA_PRICE_REFETCH,
    retry: false,
  });
}
