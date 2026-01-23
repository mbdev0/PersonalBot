import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useAddTask } from '../../hooks/useTasks';
import { useWallets } from '../../../wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { TokenAddressEntry } from './fields/tokenAddress';
import { BuyEntry } from './fields/buy';
import { Button } from '@/components/ui/button';

interface BuyTaskEntryProps {
  onClose: () => void;
}

export function BuyTaskEntry({ onClose }: BuyTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();
  const postMutation = useAddTask();

  const [slippage, setSlippage] = useState(20);
  const [computeUnits, setComputeUnits] = useState(100000);
  const [tokenAddress, setTokenAddress] = useState('');
  const [buyAmount, setBuyAmount] = useState(1.5);
  const [buyFee, setBuyFee] = useState(0.1);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);

  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody = {
      type: 'Buy',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      buy_amount: buyAmount,
      buy_fee: buyFee,
    };

    postMutation.mutate(taskBody);
    onClose();
  };

  if (isPending) return <div className="text-center py-8">Loading wallets...</div>;
  if (isError) return <div className="text-center py-8 text-red-600">Error: {error.message}</div>;
  if (data?.length === 0)
    return <div className="text-center py-8">No wallets available. Create a wallet first!</div>;

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <SlippageEntry slippage={slippage} onChange={setSlippage} />
        <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
      </div>

      <WalletSelector
        selectedWallet={wallet}
        onChange={setSelectedWallet}
        isError={isError}
        isPending={isPending}
        error={error}
        data={data}
      />

      <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress} />

      <div className="grid grid-cols-2 gap-4">
        <BuyEntry
          buyAmount={buyAmount}
          onBuyAmountChange={setBuyAmount}
          buyFee={buyFee}
          onBuyFeeChange={setBuyFee}
        />
      </div>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!wallet || !tokenAddress}>
          Create Task
        </Button>
      </div>
    </form>
  );
}
