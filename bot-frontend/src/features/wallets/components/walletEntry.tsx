import { useState } from 'react';
import { Chain } from '../types/wallet';
import { postWalletMutation } from '../hooks/useWallets';
import {
  getPubkeyFromPrivateKey,
  isValidSolanaPrivateKey,
} from '../../../utils/crypto/private_key';
import {
  Select,
  SelectContent,
  SelectTrigger,
  SelectValue,
  SelectItem,
} from '@/components/ui/select';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';

function WalletEntry({ onCompletion }: { onCompletion: () => void }) {
  const [walletName, setWalletName] = useState('');
  const [public_key, setPublicKey] = useState('');
  const [chain, setChain] = useState(Chain.Solana);
  const [walletPrivateKey, setWalletPrivateKey] = useState('');
  const mutation = postWalletMutation();

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
    <div className="space-y-5">
      <div className="grid grid-cols-2">
        <div className="space-y-2">
          <Label>Wallet Chain</Label>
          <Select value={chain} onValueChange={(e) => setChain(Chain[e as Chain])}>
            <SelectTrigger id="taskType">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem key={'solana'} value={Chain.Solana}>
                Solana
              </SelectItem>
              <SelectItem key={'bsc'} value={Chain.BSC}>
                BSC
              </SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="space-y-2 w-2xs">
          <Label>Wallet Name</Label>
          <Input
            type="text"
            name="walletName"
            id="walletName"
            placeholder="WalletName"
            value={walletName}
            onChange={(e) => setWalletName(e.target.value)}
          />
        </div>
      </div>

      <div className="grid grid-cols-2">
        <div className="space-y-2 w-8/12">
          <Label>Wallet Address</Label>
          <Input
            type="text"
            name="public_key"
            id="public_key"
            placeholder="WalletAddress"
            value={public_key}
            disabled={true}
          />
        </div>

        <div>
          <Label>Private Key</Label>
          <div className="space-y-2 w-10/12 flex gap-3">
            <Input
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

            <Button className="check_button" onClick={() => handleCheck(walletPrivateKey)}>
              Check
            </Button>
          </div>
        </div>
      </div>

      {!!error && (
        <div className="text-red-700">
          <Label>Error: {error}</Label>
        </div>
      )}

      <div className="space-x-2 flex justify-end">
        <Button
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
        </Button>

        <Button
          onClick={() => {
            onCompletion();
          }}
        >
          Cancel
        </Button>
      </div>
    </div>
  );
}

export default WalletEntry;
