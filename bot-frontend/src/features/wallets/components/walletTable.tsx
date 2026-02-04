import { useState } from 'react';
import { deleteWalletMutation, useWallets } from '../hooks/useWallets';
import WalletUpdate from './walletUpdate';
import { type Wallet } from '../types/wallet';
import { BotDialog } from '@/components/botDialog';
import { DialogHeader } from '@/components/ui/dialog';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Pencil, Trash2 } from 'lucide-react';

function WalletTable() {
  const { isPending, isError, data, error } = useWallets();
  const [editingWallet, setEditingWallet] = useState<Wallet | null>(null);
  const deleteMutation = deleteWalletMutation();

  if (isPending) {
    return <div className="loading">Loading...</div>;
  }

  if (isError) {
    return <div className="error"> Error: {error.message}</div>;
  }

  return (
    <div className="wallet_table">
      <BotDialog isOpen={!!editingWallet} onClose={() => setEditingWallet(null)}>
        <DialogHeader className="font-semibold text-foreground">Edit Wallet</DialogHeader>

        {editingWallet && (
          <WalletUpdate
            wallet={editingWallet}
            onSuccess={() => setEditingWallet(null)}
            onCancel={() => setEditingWallet(null)}
          ></WalletUpdate>
        )}
      </BotDialog>
      <div className="table-container">
        <Table>
          <TableHeader>
            <TableRow className="hover:bg-transparent">
              <TableHead className="table-header-cell text-center">Wallet Name</TableHead>
              <TableHead className="table-header-cell text-center">Chain</TableHead>
              <TableHead className="table-header-cell text-center">Public Key</TableHead>
              <TableHead className="table-header-cell text-center">Actions</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {data?.map((wallet) => (
              <TableRow key={wallet.wallet_name} className="group">
                <TableCell className="table-body-cell text-center">{wallet.wallet_name}</TableCell>
                <TableCell className="table-body-cell text-center">{wallet.chain}</TableCell>
                <TableCell className="table-body-cell text-center">
                  {`${wallet.public_key.slice(0, 6)}...${wallet.public_key.slice(-6)}`}
                </TableCell>
                <TableCell className="py-4 px-6 text-center">
                  <div className="flex gap-1.5 justify-center">
                    <button
                      onClick={() => setEditingWallet(wallet)}
                      className="action-button-edit rounded-xs"
                    >
                      <Pencil className="action-icon" />
                    </button>
                    <button
                      onClick={() => deleteMutation.mutate(wallet.id)}
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
    </div>
  );
}

export default WalletTable;
