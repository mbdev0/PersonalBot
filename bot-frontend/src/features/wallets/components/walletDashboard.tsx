import WalletTable from './walletTable';
import './walletDashboard.css';
import WalletEntry from './walletEntry';
import { useState } from 'react';
import { BotDialog } from '@/components/botDialog';
import { DialogHeader } from '@/components/ui/dialog';

function WalletDashboard() {
  // we will start with creating a table
  const [isAddModalShowing, setAddModal] = useState(false);

  return (
    <div className="wallet_dashboard">
      <h2>Wallet Dashboard</h2>
      <WalletTable />
      <BotDialog isOpen={isAddModalShowing} onClose={() => setAddModal(false)}>
        <DialogHeader className="font-bold text-foreground">Add Wallet</DialogHeader>
        <WalletEntry onCompletion={() => setAddModal(false)}></WalletEntry>
      </BotDialog>
      <button className="add_wallet" onClick={() => setAddModal(true)}>
        Add Wallet
      </button>
    </div>
  );
}

export default WalletDashboard;
