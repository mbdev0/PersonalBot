import type { Wallet } from '@/features/wallets/types/wallet';
import { FormDataLoader } from './fields/formDataLoader';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import React, { useState } from 'react';
import { useAddStrategy } from '../../hooks/useStrategy';
import {
  BUY_AMOUNT_DEFAULT,
  BUY_FEE_DEFAULT,
  COMPUTE_UNITS_DEFAULT,
  SELL_FEE_DEFAULT,
  SLIPPAGE_DEFAULT,
} from '../../utils/constants';
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
import type { SpamStrategyTaskPost } from '../../types/strategies/strategyTaskPost';
import { Switch } from '@/components/ui/switch';

interface SpamTaskEntryProps {
  program: string;
  onClose: () => void;
}

export function SpamTaskEntry({ program, onClose }: SpamTaskEntryProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <SpamTaskEntryForm
          program={program}
          onClose={onClose}
          wallets={wallets}
          rpcGroups={rpcGroups}
        ></SpamTaskEntryForm>
      )}
    </FormDataLoader>
  );
}

interface SpamTaskEntryFormProps {
  program: string;
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
}

function SpamTaskEntryForm({ program, onClose, wallets, rpcGroups }: SpamTaskEntryFormProps) {
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [buyAmount, setBuyAmount] = useState(BUY_AMOUNT_DEFAULT);
  const [buyFee, setBuyFee] = useState(BUY_FEE_DEFAULT);
  const [tokenAddress, setTokenAddress] = useState('');
  const [sellFee, setSellFee] = useState(SELL_FEE_DEFAULT);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);
  const [retriesEnabled, setRetriesFunctionality] = useState(false);
  const [retries, setRetries] = useState<number>(0);
  const [retriesDelayMs, setRetriesDelayMs] = useState<number>(0);
  const [startTime, setStartTime] = useState<Date>(() => {
    const date = new Date();
    date.setHours(date.getHours() + 1);
    date.setSeconds(0);
    return date;
  });
  const [startTimeEnabled, setStartTimeEnabled] = useState(false);

  const [numberOfTasks, setNumberOfTasks] = useState(1);

  const wallet = selectedWallet ?? wallets[0];
  const rpcGroup = selectedRpcGroup ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const payload: SpamStrategyTaskPost = {
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

    postMutation.mutate(payload, { onSuccess: onClose });
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
          Create Task
        </Button>
      </div>
    </form>
  );
}
