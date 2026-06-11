import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useAddStrategy } from '../../hooks/useStrategy';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { BuyEntry } from './fields/buy';
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
import type { SellStrategy } from '../../types/sellStrategies';
import type { Filters } from '../../types/filters';
import type { AFKStrategyTaskPost } from '../../types/strategies/strategyTaskPost';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { FormDataLoader } from './fields/formDataLoader';
import { Switch } from '@/components/ui/switch';
import { RetryEntry } from './fields/retries';

interface AFKTaskEntryProps {
  program: string;
  onClose: () => void;
}

export function AFKTaskEntry({ onClose, program }: AFKTaskEntryProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <AFKTaskForm onClose={onClose} wallets={wallets} rpcGroups={rpcGroups} program={program} />
      )}
    </FormDataLoader>
  );
}

interface AFKTaskFormProps {
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
  program: string;
}

function AFKTaskForm({ onClose, wallets, rpcGroups, program }: AFKTaskFormProps) {
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [buyAmount, setBuyAmount] = useState(BUY_AMOUNT_DEFAULT);
  const [buyFee, setBuyFee] = useState(BUY_FEE_DEFAULT);
  const [sellFee, setSellFee] = useState(SELL_FEE_DEFAULT);
  const [filters, setFilters] = useState<Filters>({});
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>([]);
  const [retriesEnabled, setRetriesFunctionality] = useState(false);
  const [retries, setRetries] = useState<number>(0);
  const [retriesDelayMs, setRetriesDelayMs] = useState<number>(0);

  const wallet = selectedWallet ?? wallets[0]; // ← no optional chaining, guaranteed by FormDataLoader
  const rpcGroup = selectedRpcGroup ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const strategyBody: AFKStrategyTaskPost = {
      program: program,
      trading_type: 'AFK',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: sellFee,
      sell_strategies: sellStrategies,
      filters,
      rpc_group_id: rpcGroup?.id,
    };

    if (retriesEnabled) {
      strategyBody.retries = retries;
      strategyBody.retry_delay_ms = retriesDelayMs;
    }

    postMutation.mutate(strategyBody, { onSuccess: onClose }); // ← close only on success
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <Card className="p-3">
          <h2>Buy/Sell Settings</h2>
          <div className="flex flex-wrap gap-4">
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

          <div className="flex flex-wrap gap-4">
            <SlippageEntry slippage={slippage} onChange={setSlippage} />
            <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
            <WalletSelector selectedWallet={wallet} onChange={setSelectedWallet} />
            <RPCGroupSelector selectedRpcGroup={rpcGroup} onChange={setSelectedRpcGroup} />
            <div className="flex flex-row gap-2">
              <div className="flex flex-col">
                <Label>Enable Retries? </Label>
                <div className="flex flex-1 items-center justify-center">
                  <Switch checked={retriesEnabled} onCheckedChange={setRetriesFunctionality} />
                </div>
              </div>
              <RetryEntry
                isRetryEnabled={!retriesEnabled}
                maxRetries={retries}
                onMaxRetryChange={setRetries}
                retryDelayMs={retriesDelayMs}
                onRetryDelayChange={setRetriesDelayMs}
              />
            </div>
          </div>
        </Card>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <SellStrategies
          sellStrategies={sellStrategies}
          setSellStrategies={setSellStrategies}
        ></SellStrategies>

        <Card className="p-2">
          <h2>Filters</h2>
          <FiltersEntry filters={filters} onFiltersChange={setFilters}></FiltersEntry>
        </Card>
      </div>

      {postMutation.isError && (
        <p className="text-sm text-destructive">{postMutation.error.message}</p>
      )}

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
