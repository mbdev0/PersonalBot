import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { SlippageEntry } from './fields/slippage';
import { TokenAddressEntry } from './fields/tokenAddress';
import { WalletSelector } from './fields/walletSelector';
import { SellEntry } from './fields/sell';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import {
  COMPUTE_UNITS_DEFAULT,
  SELL_AMOUNT_DEFAULT,
  SELL_FEE_DEFAULT,
  SLIPPAGE_DEFAULT,
} from '../../utils/constants';
import { useAddStrategy } from '../../hooks/useStrategy';
import type { SellStrategyTaskPost } from '../../types/strategies/strategyTaskPost';
import { FormDataLoader } from './fields/formDataLoader';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { RetryEntry } from './fields/retries';

interface SellTaskEntryProps {
  program: string;
  onClose: () => void;
}

export function SellTaskEntry({ onClose, program }: SellTaskEntryProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <SellTaskForm onClose={onClose} wallets={wallets} rpcGroups={rpcGroups} program={program} />
      )}
    </FormDataLoader>
  );
}

interface SellTaskFormProps {
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
  program: string;
}

function SellTaskForm({ onClose, wallets, rpcGroups, program }: SellTaskFormProps) {
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [tokenAddress, setTokenAddress] = useState('');
  const [sellAmount, setSellAmount] = useState(SELL_AMOUNT_DEFAULT);
  const [sellFee, setSellFee] = useState(SELL_FEE_DEFAULT);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);
  const [retriesEnabled, setRetriesFunctionality] = useState(false);
  const [retries, setRetries] = useState<number>(0);
  const [retriesDelayMs, setRetriesDelayMs] = useState<number>(0);

  const wallet = selectedWallet ?? wallets[0];
  const rpcGroup = selectedRpcGroup ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const taskBody: SellStrategyTaskPost = {
      trading_type: 'SELL',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      sell_amount: sellAmount / 100,
      sell_fee: sellFee,
      rpc_group_id: rpcGroup?.id,
      program: program,
    };

    if (retriesEnabled) {
      taskBody.retries = retries;
      taskBody.retry_delay_ms = retriesDelayMs;
    }

    postMutation.mutate(taskBody, { onSuccess: onClose });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <Card className="p-3">
          <h2>Sell Settings</h2>
          <div className="space-y-4">
            <div className="grid grid-cols-[120px_120px] gap-4">
              <SellEntry
                sellAmount={sellAmount}
                onSellAmountChange={setSellAmount}
                sellFee={sellFee}
                onSellFeeChange={setSellFee}
              />
            </div>
            <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress} />
          </div>
        </Card>

        <Card className="p-3">
          <h2>Task Options</h2>
          <div className="flex flex-row flex-wrap gap-4">
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

      {postMutation.isError && (
        <p className="text-sm text-destructive">{postMutation.error.message}</p>
      )}

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!wallet || !tokenAddress}>
          Create Task
        </Button>
      </div>
    </form>
  );
}
