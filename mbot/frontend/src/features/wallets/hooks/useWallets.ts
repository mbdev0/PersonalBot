import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { deleteWallet, getWallets, postWallet, updateWallet } from '../api/wallets';
import type { WalletPost, WalletPut } from '../types/wallet';
import { toast } from 'sonner';

export function useWallets() {
  return useQuery({
    queryKey: ['wallets'],
    queryFn: () => {
      return getWallets();
    },
  });
}

export function usePutWallet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (wallet: WalletPut) => updateWallet(wallet),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}

export function usePostWallet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (wallet: WalletPost) => postWallet(wallet),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
  });
}

export function useDeleteWallet() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: string) => deleteWallet(id),
    onSuccess() {
      queryClient.invalidateQueries({ queryKey: ['wallets'] });
    },
    onError: (e) => toast.error('Failure to delete wallet', { description: e.message }),
  });
}
