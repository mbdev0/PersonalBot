import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { useFormData } from '@/features/tasks/hooks/useFormData';
import type { Wallet } from '@/features/wallets/types/wallet';

interface FormDataLoaderProps {
  children: (wallets: Wallet[], rpcGroups: RPCGroupDashboardRow[]) => React.ReactNode;
}

export function FormDataLoader({ children }: FormDataLoaderProps) {
  const { wallets, rpcGroups, isPending, isError, error } = useFormData();

  if (isPending) return <div className="text-center py-8">Loading...</div>;
  if (isError) return <div className="text-center py-8 text-red-600">{error?.message}</div>;
  if (!wallets.length)
    return <div className="text-center py-8">No wallets available. Create a wallet first!</div>;
  if (!rpcGroups.length)
    return (
      <div className="text-center py-8">No RPC Groups available. Create a RPC Group first!</div>
    );

  return <>{children(wallets, rpcGroups)}</>;
}
