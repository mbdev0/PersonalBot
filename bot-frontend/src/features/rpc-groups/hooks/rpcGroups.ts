import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import { getRPCDashboard, getRPCGroup, postRPCGroup } from '../api/rpcGroups';
import type { RPCGroupPost } from '../types/rpcGroup';

export function useRpcGroupDashboard() {
  useQuery({
    queryKey: ['rpcGroupDashboard'],
    queryFn: getRPCDashboard,
  });
}

export function usePostRPCGroup() {
  const client = useQueryClient();
  useMutation({
    mutationFn: (rpcGroup: RPCGroupPost) => postRPCGroup(rpcGroup),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ['rpcGroupDashboard'] });
    },
  });
}

export function useRpcGroup({ id }: { id: number }) {
  useQuery({
    queryKey: ['rpcGroup', id],
    queryFn: () => getRPCGroup(id),
  });
}
