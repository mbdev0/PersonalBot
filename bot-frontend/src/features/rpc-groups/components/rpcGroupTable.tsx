import { BotDialog } from '@/components/botDialog';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Eye, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { RpcGroupUpdate } from './rpcGroupUpdate';
import type { RPCGroupDashboardRow } from '../types/rpcGroup';
import { useDeleteRPCGroup, useRpcGroupDashboard } from '../hooks/rpcGroups';

export function RpcGroupTable() {
  const [editingRpcGroup, setEditingRpcGroup] = useState<RPCGroupDashboardRow | null>();
  const { isPending, isError, data, error } = useRpcGroupDashboard();
  const deleteMutation = useDeleteRPCGroup();

  if (isPending) {
    return <div className="loading">Loading...</div>;
  }

  if (isError) {
    return <div className="error"> Error: {error.message}</div>;
  }

  return (
    <>
      {editingRpcGroup && (
        <BotDialog onClose={() => setEditingRpcGroup(null)} isOpen={!!editingRpcGroup}>
          <RpcGroupUpdate
            editingRowId={editingRpcGroup.id}
            onCompletion={() => setEditingRpcGroup(null)}
          ></RpcGroupUpdate>
        </BotDialog>
      )}

      <div className="table-container">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="table-header-cell text-center">RPC Group Name</TableHead>
              <TableHead className="table-header-cell text-center">No. Of RPC's</TableHead>
              <TableHead className="table-header-cell text-center">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data?.map((rpcGroup) => (
              <TableRow key={rpcGroup.id} className="group">
                <TableCell className="table-body-cell text-center">{rpcGroup.name}</TableCell>
                <TableCell className="table-body-cell text-center">{rpcGroup.number}</TableCell>
                <TableCell className="py-4 px-6 text-center">
                  <div className="flex gap-1.5 justify-center">
                    <button
                      onClick={() => setEditingRpcGroup(rpcGroup)}
                      className="action-button-edit rounded-xs"
                    >
                      <Eye className="action-icon" />
                    </button>
                    <button
                      onClick={() => deleteMutation.mutate(rpcGroup.id)}
                      className="action-button-delete rounded-xs"
                    >
                      <Trash2 className="action-icon" />
                    </button>
                  </div>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>
    </>
  );
}
