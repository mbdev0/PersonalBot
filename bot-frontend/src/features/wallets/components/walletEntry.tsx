// we make a wallet entry component here

import { useState } from 'react';

// TODO: WILL BE A SHARED COMPONENTS BETWEEN UPDATE AND CREATING
function WalletEntry() {
  const [walletName, setWalletName] = useState('');
  const [walletPrivateKey, setWalletPrivateKey] = useState('');
  const [chain, setChain] = useState('');

  const walletBody = {
    wallet_name: walletName,
    chain: chain,
    private_key: walletPrivateKey,
  };

  return (
    <>
      <div className="task_type">
        <h3>Task Type</h3>
        <select value={chain} onChange={(e) => setChain(e.target.value)}>
          <option value="Solana">Solana</option>
        </select>
      </div>
      <div className="wallet_name">
        <h3>Wallet Name</h3>
        <input
          type="text"
          name="walletName"
          id="walletName"
          placeholder="WalletName"
          value={walletName}
          onChange={(e) => setWalletName(e.target.value)}
        />
      </div>
      <div className="wallet_private_key">
        <h3>Wallet Name</h3>
        <input
          type="text"
          name="wallet_private_key"
          id="wallet_private_key"
          placeholder="Wallet Private Key"
          value={walletPrivateKey}
          onChange={(e) => setWalletPrivateKey(e.target.value)}
        />
      </div>

      <button
        onClick={(_) => {
          create_wallet_api_call(walletBody);
        }}
      >
        Create Wallet
      </button>
    </>
  );
}

// TODO: MOVE THIS TO A POST MUTATION
async function create_wallet_api_call(wallet_body: any) {
  // console.log(private_key, token_address);
  const url = 'http://localhost:9090/api/wallet/wallets';

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(wallet_body),
  });

  if (!response.ok) {
    console.log('Error whilst sending request:', response.status);
    console.log(response);
    return;
  }
  console.log('response: ', response.statusText);
}

export default WalletEntry;
