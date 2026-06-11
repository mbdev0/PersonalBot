import { useState } from 'react';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import type { Wallet } from '@/features/wallets/types/wallet';
import { TokenAddressEntry } from './fields/tokenAddress';
import { SELL_AMOUNT_DEFAULT, SELL_FEE_DEFAULT } from '../../utils/constants';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { SellEntry } from './fields/sell';
import type { SellStrategyTaskPut } from '../../types/strategies/strategyTaskPut';
import { useUpdateStrategy } from '../../hooks/useStrategy';
import { FormDataLoader } from './fields/formDataLoader';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import type { SellStrategyTask } from '../../types/strategies/strategyTask';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { RetryEntry } from './fields/retries';

interface SellTaskEditProps {
  task: SellStrategyTask;
  onClose: () => void;
  program: string;
}

export function SellTaskEdit({ task, onClose, program }: SellTaskEditProps) {
  if (!task.sell_amount) {
    return (
      <div className="text-center py-8 text-red-600">
        Error: task.sell_amount was undefined - double check the logic
      </div>
    );
  }

  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <SellEditForm
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

interface SellEditFormProps {
  task: SellStrategyTask;
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
  program: string;
}

function SellEditForm({ task, onClose, wallets, rpcGroups, program }: SellEditFormProps) {
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [sellAmount, setSellAmount] = useState(task.sell_amount * 100);
  const [sellFee, setSellFee] = useState(task.sell_fee);
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

    const taskBody: SellStrategyTaskPut = {
      program: program,
      id: task.id,
      trading_type: 'SELL',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      sell_amount: (sellAmount ?? SELL_AMOUNT_DEFAULT) / 100,
      sell_fee: sellFee ?? SELL_FEE_DEFAULT,
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
          <h2>Sell Settings</h2>
          <div className="space-y-4">
            <div className="grid grid-cols-[120px_120px] gap-4">
              <SellEntry
                sellAmount={sellAmount ?? SELL_AMOUNT_DEFAULT}
                onSellAmountChange={setSellAmount}
                sellFee={sellFee ?? SELL_FEE_DEFAULT}
                onSellFeeChange={setSellFee}
              />
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
