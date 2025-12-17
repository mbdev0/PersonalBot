// this will be the dashboard component - the main view
// we should expect to see all the wallets that exist on the api here
// what do we need to create?
// a table to see the existing wallets
// an add button
// a button to update a wallet (next to each row)

import WalletTable from './walletTable';
import './walletDashboard.css';
import WalletEntry from './walletEntry';
import Modal from '../../../components/modal';
import { useState } from 'react';
//this should be a composition of the following components:
// - a wallet entry menu
// - an update menu X
// - a wallet card menu - view 1 wallet X
// - a table to see the wallets X

function WalletDashboard() {
  // we will start with creating a table
  const [isAddModalShowing, setAddModal] = useState(false);

  return (
    <div className="wallet_dashboard">
      <h2>Wallet Dashboard</h2>
      <WalletTable />
      <Modal isOpen={isAddModalShowing} onClose={() => setAddModal(false)}>
        <WalletEntry onCompletion={() => setAddModal(false)}></WalletEntry>
      </Modal>
      <button className="add_wallet" onClick={() => setAddModal(true)}>
        Add Wallet
      </button>
    </div>
  );
}

export default WalletDashboard;
