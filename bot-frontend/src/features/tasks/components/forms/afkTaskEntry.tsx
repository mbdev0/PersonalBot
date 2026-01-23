import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useAddStrategy } from '../../hooks/useStrategy';
import { useWallets } from '../../../wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { BuyEntry } from './fields/buy';
import type { StrategyTaskPost } from '../../types/strategyTask';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';

interface AFKTaskEntryProps {
  onClose: () => void;
}

export function AFKTaskEntry({ onClose }: AFKTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(20);
  const [computeUnits, setComputeUnits] = useState(100000);
  const [buyAmount, setBuyAmount] = useState(1.0);
  const [buyFee, setBuyFee] = useState(0.1);
  const [sellFee, setSellFee] = useState(0.1);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);

  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const strategyBody: StrategyTaskPost = {
      tradingType: 'AFK',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: sellFee,
      sell_strategies: [
        {
          type: 'take_profit_percentage',
          value: 0.2,
          sell_amount: 1,
        },
      ],
      filters: {
        has_twitter: true,
      },
    };

    postMutation.mutate(strategyBody);
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

      <div className="grid grid-cols-2 gap-4">
        <BuyEntry
          buyAmount={buyAmount}
          onBuyAmountChange={setBuyAmount}
          buyFee={buyFee}
          onBuyFeeChange={setBuyFee}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="sell_fee">Sell Fee</Label>
        <Input
          id="sell_fee"
          type="number"
          placeholder="0.00"
          value={sellFee}
          onChange={(e) => setSellFee(e.target.valueAsNumber)}
        />
      </div>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!wallet}>
          Create Task
        </Button>
      </div>
    </form>
  );
}
