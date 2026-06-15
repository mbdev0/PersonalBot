import WalletTable from './walletTable';
import WalletEntry from './walletEntry';
import { useState } from 'react';
import { BotDialog } from '@/components/botDialog';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';

function WalletDashboard() {
  // we will start with creating a table
  const [isAddModalShowing, setAddModal] = useState(false);

  return (
    <div className="wallet_dashboard space-y-8 py-4">
      <div className="flex items-center justify-end">
        <Button
          className="h-9 px-4 text-[13px] text-accent-foreground font-medium bg-foreground/5 hover:bg-foreground/10 border-0 ring-1 ring-foreground/20 hover:ring-foreground/30 transition-all duration-200"
          onClick={() => setAddModal(true)}
        >
          Add Wallet
        </Button>
      </div>
      <WalletTable />
      <BotDialog isOpen={isAddModalShowing} onClose={() => setAddModal(false)}>
        <DialogHeader>
          <DialogTitle className="font-semibold text-foreground">Add Wallet</DialogTitle>
        </DialogHeader>
        <WalletEntry onCompletion={() => setAddModal(false)}></WalletEntry>
      </BotDialog>
    </div>
  );
}

export default WalletDashboard;
