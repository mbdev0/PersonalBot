import type { Wallet } from '@/features/wallets/types/wallet';
import { FormDataLoader } from './fields/formDataLoader';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import React, { useState } from 'react';
import { useUpdateStrategy } from '../../hooks/useStrategy';
import { Card } from '@/components/ui/card';
import { TokenAddressEntry } from './fields/tokenAddress';
import { BuyEntry } from './fields/buy';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { SlippageEntry } from './fields/slippage';
import { WalletSelector } from './fields/walletSelector';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import { RetryEntry } from './fields/retries';
import { StartTime } from './fields/startTime';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { NumberOfTasksEntry } from './fields/noOfTasks';
import { Button } from '@/components/ui/button';
import { Switch } from '@/components/ui/switch';
import type { SpamStrategyTask } from '../../types/strategies/strategyTask';
import type { SpamStrategyTaskPut } from '../../types/strategies/strategyTaskPut';

interface SpamTaskEditProps {
  program: string;
  onClose: () => void;
  task: SpamStrategyTask;
}

export function SpamTaskEdit({ program, onClose, task }: SpamTaskEditProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <SpamTaskEditForm
          program={program}
          onClose={onClose}
          wallets={wallets}
          rpcGroups={rpcGroups}
          task={task}
        ></SpamTaskEditForm>
      )}
    </FormDataLoader>
  );
}

interface SpamTaskEditFormProps {
  program: string;
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
  task: SpamStrategyTask;
}

function SpamTaskEditForm({ program, onClose, wallets, rpcGroups, task }: SpamTaskEditFormProps) {
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [sellFee, setSellFee] = useState(task.sell_fee);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);

  const [retriesEnabled, setRetriesFunctionality] = useState(() => {
    return task.retries != null;
  });

  const [retries, setRetries] = useState<number>(task.retries ?? 0);
  const [retriesDelayMs, setRetriesDelayMs] = useState<number>(task.retry_delay_ms ?? 0);
  const [startTime, setStartTime] = useState<Date>(() => {
    const date = new Date(1781050980852);
    return date;
  });

  const [startTimeEnabled, setStartTimeEnabled] = useState(task.start_time != 0);

  const [numberOfTasks, setNumberOfTasks] = useState(task.no_of_tasks);

  const wallet =
    selectedWallet ?? wallets.find((w) => w.wallet_name === task.wallet_name) ?? wallets[0];
  const rpcGroup =
    selectedRpcGroup ?? rpcGroups.find((rg) => rg.name === task.rpc_group) ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const payload: SpamStrategyTaskPut = {
      id: task.id,
      program: program,
      trading_type: 'SPAM',
      buy_fee: buyFee,
      buy_amount: buyAmount,
      token_address: tokenAddress,
      no_of_tasks: numberOfTasks,
      sell_fee: sellFee,
      compute_units: computeUnits,
      slippage: slippage / 100,
      wallet_name: wallet.wallet_name,
      rpc_group_id: rpcGroup.id,
      start_time: 0,
    };

    if (retriesEnabled) {
      payload.retries = retries;
      payload.retry_delay_ms = retriesDelayMs;
    }

    if (startTimeEnabled) {
      payload.start_time = startTime.valueOf();
    }

    putMutation.mutate(payload, { onSuccess: onClose });
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="flex gap-4">
        <Card className="flex flex-row flex-wrap p-4">
          {/*buy options*/}
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

          <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress} />
        </Card>

        <Card className="flex flex-row flex-wrap p-4">
          {/*{' '}*/}
          {/*task options*/}
          <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
          <SlippageEntry slippage={slippage} onChange={setSlippage} />
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

          <div className="flex gap-2">
            <div className="flex flex-col">
              <Label>Schedule Task?</Label>
              <div className="flex flex-1 items-center justify-center">
                <Switch checked={startTimeEnabled} onCheckedChange={setStartTimeEnabled} />
              </div>
            </div>
            <StartTime
              isStartTimeEnabled={!startTimeEnabled}
              startTime={startTime}
              onStartTimeChange={setStartTime}
            />
          </div>
          <NumberOfTasksEntry numberOfTasks={numberOfTasks} onChange={setNumberOfTasks} />
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
