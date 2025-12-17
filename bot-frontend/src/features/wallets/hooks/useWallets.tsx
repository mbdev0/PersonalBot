import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getWallets, postWallet, updateWallet } from '../api/wallets';
import type { WalletPost, WalletPut } from '../types/wallet';

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

export function PostWalletMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (wallet: WalletPost) => postWallet(wallet),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}
