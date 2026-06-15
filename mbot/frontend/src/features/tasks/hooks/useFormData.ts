import { useRpcGroupDashboard } from '@/features/rpc-groups/hooks/rpcGroups';
import { useWallets } from '@/features/wallets/hooks/useWallets';

export function useFormData() {
  const {
    isPending: walletsPending,
    isError: walletsError,
    data: wallets,
    error: walletErr,
  } = useWallets();
  const {
    isPending: rpcPending,
    isError: rpcError,
    data: rpcGroups,
    error: rpcErr,
  } = useRpcGroupDashboard();

  return {
    wallets: wallets ?? [],
    rpcGroups: rpcGroups ?? [],
    isPending: walletsPending || rpcPending,
    isError: walletsError || rpcError,
    error: walletErr ?? rpcErr,
  };
}
