// this will be the dashboard component - the main view
// we should expect to see all the wallets that exist on the api here
// what do we need to create?
// a table to see the existing wallets
// an add button
// a button to update a wallet (next to each row)

import WalletTable from './walletTable';
//this should be a composition of the following components:
// - a wallet entry menu
// - an update menu
// - a wallet card menu - view 1 wallet
// - a table to see the wallets

function WalletDashboard() {
  // we will start with creating a table
  return (
    <>
      <h2>Wallet Dashboard</h2>
      <WalletTable />
    </>
  );
}

export default WalletDashboard;
