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
import type { RPCGroup } from '../types/rpcGroup';


const data: RPCGroup[] = [
  {
    rpc_group_name: 'main',
    rpcs: [
      {
        http: 'https://abc.com',
        ws: 'ws://abc.com',
      },
      {
        http: 'https://abc.com',
        ws: 'ws://abc.com',
      },
    ],
  },
  {
    rpc_group_name: 'main2',
    rpcs: [
      {
        http: 'https://abc.com',
        ws: 'ws://abc.com',
      },
      {
        http: 'https://abc.com',
        ws: 'ws://abc.com',
      },
      {
        http: 'https://abc.com',
        ws: 'ws://abc.com',
      },
    ],
  },
];

export function RpcGroupTable() {
  //in here we'll have our edit modal + onClick of eye, setEditingWallet
  //botdialog will show up if setEditingWallet is set, hidden on null
  const [editingRpcGroup, setEditingRpcGroup] = useState<RPCGroup | null>();

  return (
    <>
      {editingRpcGroup && (
        <BotDialog onClose={() => setEditingRpcGroup(null)} isOpen={!!editingRpcGroup}>
          <RpcGroupUpdate
            editingRow={editingRpcGroup}
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
              <TableRow key={rpcGroup.rpc_group_name} className="group">
                <TableCell className="table-body-cell text-center">
                  {rpcGroup.rpc_group_name}
                </TableCell>
                <TableCell className="table-body-cell text-center">
                  {rpcGroup.rpcs.length}
                </TableCell>
                <TableCell className="py-4 px-6 text-center">
                  <div className="flex gap-1.5 justify-center">
                    <button
                      onClick={() => setEditingRpcGroup(rpcGroup)}
                      className="action-button-edit rounded-xs"
                    >
                      <Eye className="action-icon" />
                    </button>
                    <button
                      // onClick={() => deleteMutation.mutate(wallet.id)}
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
