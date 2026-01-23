import { useWallets } from '@/features/wallets/hooks/useWallets';
import { useUpdateStrategy } from '../../hooks/useStrategy';
import type { StrategyTask, StrategyTaskPut } from '../../types/strategyTask';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import type { Wallet } from '@/features/wallets/types/wallet';
import { BuyEntry } from './fields/buy';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { SlippageEntry } from './fields/slippage';
import { WalletSelector } from './fields/walletSelector';
import { Label } from '@radix-ui/react-label';
import { Input } from '@/components/ui/input';

interface AFKTaskEditProps {
  task: StrategyTask;
  onClose: () => void;
}

export function AFKTaskEdit({ task, onClose }: AFKTaskEditProps) {
  const { isPending, isError, data, error } = useWallets();
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [sellFee, setSellFee] = useState(task.sell_fee);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);

  const wallet = data?.find((w) => w.wallet_name === task.wallet_name) ?? null;

  const handleWalletChange = (wallet: Wallet | null) => {
    if (wallet) {
      setSelectedWallet(wallet);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody: StrategyTaskPut = {
      tradingType: 'AFK',
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: 0, //TODO:
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      filters: task.filters, //TODO
      sell_strategies: [], //TODO
      id: task.id,
    };

    putMutation.mutate(taskBody);
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
        onChange={handleWalletChange}
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

      {/* IS class name needed? */}
      <div className="space-y-2">
        <Label htmlFor="sellFee">Sell Fee</Label>
        <Input
          id="sellFee"
          type="number"
          value={sellFee}
          onChange={(e) => setSellFee(e.target.valueAsNumber)}
        />
      </div>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!selectedWallet}>
          Update Task
        </Button>
      </div>
    </form>
  );
}
