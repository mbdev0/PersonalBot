import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getWallets, updateWallet } from '../api/wallets';
import type { WalletPut } from '../types/wallet';

export function useWallets() {
  return useQuery({
    queryKey: ['wallets'],
    queryFn: getWallets,
  });
}

export function PutWalletMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (wallet: WalletPut) => updateWallet(wallet),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}
