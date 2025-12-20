import { useState } from 'react';
import { deleteWalletMutation, useWallets } from '../hooks/useWallets';
import Modal from '../../../components/modal';
import WalletUpdate from './walletUpdate';
import { type Wallet } from '../types/wallet';
import './walletTable.css';

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
      <Modal isOpen={!!editingWallet} onClose={() => setEditingWallet(null)}>
        {editingWallet && (
          <WalletUpdate
            wallet={editingWallet}
            onSuccess={() => setEditingWallet(null)}
            onCancel={() => setEditingWallet(null)}
          ></WalletUpdate>
        )}
      </Modal>
      <table>
        <thead>
          <tr>
            <th scope="col">Wallet Name</th>
            <th scope="col">Chain</th>
            <th scope="col">Public Key</th>
            <th scope="col">Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.map((wallet) => (
            <tr key={wallet.wallet_name}>
              <td>{wallet.wallet_name}</td>
              <td>{wallet.chain}</td>
              <td>{`${wallet.public_key.slice(0, 4)}...${wallet.public_key.slice(-4)}`}</td>
              <td>
                <button onClick={() => setEditingWallet(wallet)}>✏️</button>
                <button onClick={() => deleteMutation.mutate(wallet.id)}>🗑️</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}

export default WalletTable;
