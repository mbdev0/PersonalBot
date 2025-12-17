import { useState } from 'react';
import { Chain } from '../types/wallet';
import getEnumKeys from '../../../utils/enum_helper';
import { PostWalletMutation } from '../hooks/useWallets';
import {
  getPubkeyFromPrivateKey,
  isValidSolanaPrivateKey,
} from '../../../utils/crypto/private_key';
import './walletUpdate.css';

function WalletEntry({ onCompletion }: { onCompletion: () => void }) {
  const [walletName, setWalletName] = useState('');
  const [public_key, setPublicKey] = useState('');
  const [chain, setChain] = useState<Chain>(Chain.Solana);
  const [walletPrivateKey, setWalletPrivateKey] = useState('');
  const mutation = PostWalletMutation();

  //temporary - how should we show errors? toaster (sonner - shadcn) + field errors -> validation only?
  const [error, setError] = useState('');

  const handleCheck = (priv_key: string) => {
    if (!isValidSolanaPrivateKey(priv_key)) {
      setError('Private key is not valid');
      return;
    }

    const pubkey = getPubkeyFromPrivateKey(priv_key);
    if (!pubkey) {
      setError('Failed to process the private key');
      return;
    }

    setPublicKey(pubkey);
    setWalletPrivateKey(priv_key);
    setError('');
  };

  return (
    <div className="wallet-update-form">
      <h2>Add New Wallet</h2>
      <div className="wallet_chain">
        <h4>Wallet Chain</h4>
        {/* https://stackoverflow.com/a/72883012 */}
        <select
          value={chain}
          onChange={(e) => setChain(Chain[e.target.value as keyof typeof Chain])}
        >
          {getEnumKeys(Chain).map((key, index) => (
            <option key={index} value={Chain[key]}>
              {key}
            </option>
          ))}
        </select>
      </div>
      <div className="wallet_name">
        <h4>Wallet Name</h4>
        <input
          type="text"
          name="walletName"
          id="walletName"
          placeholder="WalletName"
          value={walletName}
          onChange={(e) => setWalletName(e.target.value)}
        />
      </div>
      <div className="public_key">
        <h4>Wallet Address</h4>
        <input
          type="text"
          name="public_key"
          id="public_key"
          placeholder="WalletAddress"
          value={public_key}
          disabled={true}
        />
      </div>
      <div className="wallet_private_key">
        <h4>Private Key</h4>
        <input
          type="text"
          name="wallet_private_key"
          id="wallet_private_key"
          placeholder="Wallet Private Key"
          value={walletPrivateKey}
          onChange={(e) => {
            setWalletPrivateKey(e.target.value);
            setError('');
          }}
        />
        <button className="check_button" onClick={() => handleCheck(walletPrivateKey)}>
          Check
        </button>
      </div>

      {!!error && (
        <div className="error">
          <h4>Error: {error}</h4>
        </div>
      )}

      <div className="submission">
        <button
          onClick={() => {
            if (!isValidSolanaPrivateKey(walletPrivateKey)) {
              setError('Invalid private key');
              return;
            }

            const pubkey = getPubkeyFromPrivateKey(walletPrivateKey);
            if (!pubkey) {
              setError('Failed to derive public key');
              return;
            }

            mutation.mutate(
              {
                wallet_name: walletName,
                chain: chain,
                private_key: walletPrivateKey,
                public_key: public_key,
              },
              {
                onSuccess: () => {
                  onCompletion();
                },
                onError: (error) => {
                  console.log(error);
                },
              }
            );
          }}
          disabled={!walletPrivateKey || !isValidSolanaPrivateKey(walletPrivateKey)}
        >
          Add Wallet
        </button>

        <button
          onClick={() => {
            onCompletion();
          }}
        >
          Cancel
        </button>
      </div>
    </div>
  );
}

export default WalletEntry;
