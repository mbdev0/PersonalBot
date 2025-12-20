import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { deleteWallet, getWallets, postWallet, updateWallet } from '../api/wallets';
import type { WalletPost, WalletPut } from '../types/wallet';

export function useWallets() {
  return useQuery({
    queryKey: ['wallets'],
    queryFn: () => {
      return getWallets();
    },
  });
}

export function putWalletMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (wallet: WalletPut) => updateWallet(wallet),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}

export function postWalletMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (wallet: WalletPost) => postWallet(wallet),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}

export function deleteWalletMutation() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteWallet(id),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}
