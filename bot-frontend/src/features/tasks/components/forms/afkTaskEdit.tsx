import { useUpdateStrategy } from '../../hooks/useStrategy';
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
import type { AFKStrategyTask } from '../../types/strategies/strategyTask';
import type { AFKStrategyTaskPut } from '../../types/strategies/strategyTaskPut';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import { FormDataLoader } from './fields/formDataLoader';
import { Switch } from '@/components/ui/switch';
import { RetryEntry } from './fields/retries';

interface AFKTaskEditProps {
  task: AFKStrategyTask;
  onClose: () => void;
  program: string;
}
export function AFKTaskEdit({ task, onClose, program }: AFKTaskEditProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <AFKEditForm
          task={task}
          onClose={onClose}
          wallets={wallets}
          rpcGroups={rpcGroups}
          program={program}
        />
      )}
    </FormDataLoader>
  );
}

interface AFKEditFormProps {
  task: AFKStrategyTask;
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
  program: string;
}

function AFKEditForm({ task, onClose, wallets, rpcGroups, program }: AFKEditFormProps) {
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [sellFee, setSellFee] = useState(task.sell_fee);
  const [filters, setFilters] = useState(task.filters);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>(task.sell_strategies);
  const [retriesEnabled, setRetriesFunctionality] = useState(task.retries != null);
  const [retries, setRetries] = useState<number>(task.retries ?? 0);
  const [retriesDelayMs, setRetriesDelayMs] = useState<number>(task.retry_delay_ms ?? 0);

  // derive current values from task — fall back to task's existing wallet/rpc
  const wallet =
    selectedWallet ?? wallets.find((w) => w.wallet_name === task.wallet_name) ?? wallets[0];
  const rpcGroup =
    selectedRpcGroup ?? rpcGroups.find((rg) => rg.name === task.rpc_group) ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const taskBody: AFKStrategyTaskPut = {
      program: program,
      trading_type: 'AFK',
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: sellFee,
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      filters,
      sell_strategies: sellStrategies,
      id: task.id,
      rpc_group_id: rpcGroup.id,
    };

    if (retriesEnabled) {
      taskBody.retries = retries;
      taskBody.retry_delay_ms = retriesDelayMs;
    }

    putMutation.mutate(taskBody, { onSuccess: onClose });
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
        <SellStrategies sellStrategies={sellStrategies} setSellStrategies={setSellStrategies} />
        <Card className="p-2">
          <h2>Filters</h2>
          <FiltersEntry filters={filters} onFiltersChange={setFilters} />
        </Card>
      </div>

      {putMutation.isError && (
        <p className="text-sm text-destructive">{putMutation.error.message}</p>
      )}

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
