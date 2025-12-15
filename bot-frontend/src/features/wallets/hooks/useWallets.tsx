import { useQuery } from '@tanstack/react-query';
import { getWallets } from '../api/wallets';

export function useWallets() {
  return useQuery({
    queryKey: ['wallets'],
    queryFn: getWallets,
  });
}
