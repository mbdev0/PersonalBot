import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import {
  deleteRPCGroup,
  getRPCDashboard,
  getRPCGroup,
  postRPCGroup,
  updateRPCGroup,
} from '../api/rpcGroups';
import type { RPCGroupPost, RPCGroupPut } from '../types/rpcGroup';
import { toast } from 'sonner';

export function useRpcGroupDashboard() {
  return useQuery({
    queryKey: ['rpcGroupDashboard'],
    queryFn: getRPCDashboard,
  });
}

export function usePostRPCGroup() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (rpcGroup: RPCGroupPost) => postRPCGroup(rpcGroup),
    onSuccess: () => {
      client.invalidateQueries({ queryKey: ['rpcGroupDashboard'] });
    },
  });
}

export function useRpcGroup({ id }: { id: number }) {
  return useQuery({
    queryKey: ['rpcGroup', id],
    queryFn: () => getRPCGroup(id),
  });
}

export function useUpdateRPCGroup() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (rpcGroup: RPCGroupPut) => updateRPCGroup(rpcGroup),
    onSuccess: (_, variables) => {
      client.invalidateQueries({ queryKey: ['rpcGroupDashboard'] });
      client.invalidateQueries({ queryKey: ['rpcGroup', variables.id] });
    },
  });
}

export function useDeleteRPCGroup() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (id: number) => deleteRPCGroup(id),
    onSuccess: (_, id) => {
      client.invalidateQueries({ queryKey: ['rpcGroupDashboard'] });
      client.removeQueries({ queryKey: ['rpcGroup', id] });
    },
    onError: (e) => toast.error('Failure to delete RPC Group', { description: e.message }),
  });
}
