import { useState } from 'react';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import type { Wallet } from '@/features/wallets/types/wallet';
import { TokenAddressEntry } from './fields/tokenAddress';
import { BuyEntry } from './fields/buy';
import { BUY_AMOUNT_DEFAULT, BUY_FEE_DEFAULT } from '../../utils/constants';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { useUpdateStrategy } from '../../hooks/useStrategy';
import type { BuyStrategyTaskPut } from '../../types/strategies/strategyTaskPut';
import type { BuyStrategyTask } from '../../types/strategies/strategyTask';
import type { SellStrategy } from '../../types/sellStrategies';
import { SellStrategies } from './fields/sellStrategies/sellStrategies';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { FormDataLoader } from './fields/formDataLoader';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import { Switch } from '@/components/ui/switch';
import { RetryEntry } from './fields/retries';

interface BuyTaskEditProps {
  task: BuyStrategyTask;
  onClose: () => void;
  program: string;
}

export function BuyTaskEdit({ task, onClose, program }: BuyTaskEditProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <BuyEditForm
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

interface BuyEditFormProps {
  task: BuyStrategyTask;
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
  program: string;
}

function BuyEditForm({ task, onClose, wallets, rpcGroups, program }: BuyEditFormProps) {
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [sellFee, setSellFee] = useState<number | undefined>(task.sell_fee);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>(task.sell_strategies);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);
  const [retriesEnabled, setRetriesFunctionality] = useState(task.retries != null);
  const [retries, setRetries] = useState<number>(task.retries ?? 0);
  const [retriesDelayMs, setRetriesDelayMs] = useState<number>(task.retry_delay_ms ?? 0);

  const wallet =
    selectedWallet ?? wallets.find((w) => w.wallet_name === task.wallet_name) ?? wallets[0];
  const rpcGroup =
    selectedRpcGroup ?? rpcGroups.find((rg) => rg.name === task.rpc_group) ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const taskBody: BuyStrategyTaskPut = {
      program: program,
      id: task.id,
      trading_type: 'BUY',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_strategies: sellStrategies,
      sell_fee: sellFee,
      rpc_group_id: rpcGroup.id,
    };

    if (retriesEnabled) {
      taskBody.retries = retries;
      taskBody.retry_delay_ms = retriesDelayMs;
    }

    console.log(taskBody);
    putMutation.mutate(taskBody, { onSuccess: onClose });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <Card className="p-3">
          <h2>Buy Settings</h2>
          <div className="space-y-4">
            <div className="flex flex-wrap gap-4">
              <BuyEntry
                buyAmount={buyAmount ?? BUY_AMOUNT_DEFAULT}
                onBuyAmountChange={setBuyAmount}
                buyFee={buyFee ?? BUY_FEE_DEFAULT}
                onBuyFeeChange={setBuyFee}
              />
              <div className="space-y-2">
                <Label htmlFor="sell_fee">Sell Fee</Label>
                <Input
                  id="sell_fee"
                  type="number"
                  value={sellFee}
                  onChange={(e) => setSellFee(e.target.valueAsNumber)}
                />
              </div>
            </div>
            <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress} />
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
      </div>

      {putMutation.isError && (
        <p className="text-sm text-destructive">{putMutation.error.message}</p>
      )}

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!wallet || !tokenAddress}>
          Update Task
        </Button>
      </div>
    </form>
  );
}
