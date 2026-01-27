import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useAddStrategy } from '../../hooks/useStrategy';
import { useWallets } from '../../../wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { BuyEntry } from './fields/buy';
import {
  type Filters,
  type SellStrategy,
  type SellStrategyCreate,
  type StrategyTaskPost,
} from '../../types/strategyTask';
import { Button } from '@/components/ui/button';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import {
  BUY_AMOUNT_DEFAULT,
  BUY_FEE_DEFAULT,
  COMPUTE_UNITS_DEFAULT,
  SELL_FEE_DEFAULT,
  SLIPPAGE_DEFAULT,
} from '../../utils/constants';
import { FiltersEntry } from './fields/filters';
import { Card } from '@/components/ui/card';
import { SellStrategies } from './fields/sellStrategies/sellStrategies';

interface AFKTaskEntryProps {
  onClose: () => void;
}

export function AFKTaskEntry({ onClose }: AFKTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [buyAmount, setBuyAmount] = useState(BUY_AMOUNT_DEFAULT);
  const [buyFee, setBuyFee] = useState(BUY_FEE_DEFAULT);
  const [sellFee, setSellFee] = useState(SELL_FEE_DEFAULT);
  const [filters, setFilters] = useState<Filters>({});
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>([]);

  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const strategyBody: StrategyTaskPost = {
      trading_type: 'AFK',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: sellFee,
      sell_strategies: sellStrategies,
      filters: filters,
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

      <div>
        <Label htmlFor="filters">Filters</Label>
        <Card>
          <FiltersEntry filters={filters} onFiltersChange={setFilters}></FiltersEntry>
        </Card>
      </div>

      <div>
        <SellStrategies
          sellStrategies={sellStrategies}
          setSellStrategies={setSellStrategies}
        ></SellStrategies>
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
