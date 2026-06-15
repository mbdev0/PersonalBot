import { Label } from '@/components/ui/label';

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { useRpcGroupDashboard } from '@/features/rpc-groups/hooks/rpcGroups';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';

interface RPCGroupSelectorProps {
  selectedRpcGroup: RPCGroupDashboardRow | null;
  onChange: (rpcGroup: RPCGroupDashboardRow | null) => void;
}

export function RPCGroupSelector({ selectedRpcGroup, onChange }: RPCGroupSelectorProps) {
  const { isPending, isError, data, error } = useRpcGroupDashboard();

  const effectiveRPCGroup = selectedRpcGroup ?? data?.[0] ?? null;

  if (isPending) return <div>Loading RPC Groups...</div>;
  if (isError) return <div>Error: {error?.message}</div>;
  if (data?.length === 0) return <div>No RPC Groups. Create one first!</div>;

  return (
    <div className="space-y-2">
      <Label htmlFor="rpc_group">RPC Group</Label>
      <Select
        value={effectiveRPCGroup?.name}
        onValueChange={(e) => onChange(data?.find((rg) => rg.name === e) ?? null)}
      >
        <SelectTrigger id="rpc_group">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {data?.map((rg) => (
            <SelectItem key={rg.name} value={rg.name}>
              {rg.name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
