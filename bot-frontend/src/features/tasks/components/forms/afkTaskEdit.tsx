import { useWallets } from '@/features/wallets/hooks/useWallets';
import { useUpdateStrategy } from '../../hooks/useStrategy';
import { type StrategyTask, type StrategyTaskPut } from '../../types/strategyTask';
import { useState } from 'react';
import { Button } from '@/components/ui/button';
import type { Wallet } from '@/features/wallets/types/wallet';
import { BuyEntry } from './fields/buy';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { SlippageEntry } from './fields/slippage';
import { WalletSelector } from './fields/walletSelector';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { FiltersEntry } from './fields/filters';
import { Card } from '@/components/ui/card';
import { SellStrategies } from './fields/sellStrategies/sellStrategies';
import type { SellStrategy } from '../../types/sellStrategies';

interface AFKTaskEditProps {
  task: StrategyTask;
  onClose: () => void;
}

export function AFKTaskEdit({ task, onClose }: AFKTaskEditProps) {
  const { isPending, isError, data, error } = useWallets();
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [sellFee, setSellFee] = useState(task.sell_fee);
  const [filters, setFilters] = useState(task.filters);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>(task.sell_strategies);

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
      trading_type: 'AFK',
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: sellFee,
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: selectedWallet?.wallet_name ?? wallet.wallet_name,
      filters: filters,
      sell_strategies: sellStrategies,
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
        <Card className="p-3">
          <h2>Buy/Sell Settings</h2>
          <div className="grid grid-cols-[120px_120px_120px] gap-4">
            <BuyEntry
              buyAmount={buyAmount}
              onBuyAmountChange={setBuyAmount}
              buyFee={buyFee}
              onBuyFeeChange={setBuyFee}
            />
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
          </div>
        </Card>

        <Card className="p-3">
          <h2>Task Options</h2>
          <div className="grid grid-cols-[120px_120px_130px] gap-4">
            <SlippageEntry slippage={slippage} onChange={setSlippage} />
            <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
            <WalletSelector
              selectedWallet={selectedWallet ?? wallet}
              onChange={handleWalletChange}
            />
          </div>
        </Card>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <SellStrategies sellStrategies={sellStrategies} setSellStrategies={setSellStrategies} />

        <Card className="p-2">
          <h2>Filters</h2>
          <FiltersEntry filters={filters} onFiltersChange={setFilters} />
        </Card>
      </div>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!wallet}>
          Update Task
        </Button>
      </div>
    </form>
  );
}
