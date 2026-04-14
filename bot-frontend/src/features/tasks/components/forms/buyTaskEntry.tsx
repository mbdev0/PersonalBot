import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { TokenAddressEntry } from './fields/tokenAddress';
import { BuyEntry } from './fields/buy';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import {
  BUY_AMOUNT_DEFAULT,
  BUY_FEE_DEFAULT,
  COMPUTE_UNITS_DEFAULT,
  SELL_FEE_DEFAULT,
  SLIPPAGE_DEFAULT,
} from '../../utils/constants';
import { useAddStrategy } from '../../hooks/useStrategy';
import type { BuyStrategyTaskPost } from '../../types/strategies/strategyTaskPost';
import type { SellStrategy } from '../../types/sellStrategies';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { SellStrategies } from './fields/sellStrategies/sellStrategies';
import type { RPCGroupDashboardRow } from '@/features/rpc-groups/types/rpcGroup';
import { RPCGroupSelector } from './fields/rpcGroupSelector';
import { FormDataLoader } from './fields/formDataLoader';

interface BuyTaskEntryProps {
  onClose: () => void;
}

// BuyTaskEntry.tsx
export function BuyTaskEntry({ onClose }: BuyTaskEntryProps) {
  return (
    <FormDataLoader>
      {(wallets, rpcGroups) => (
        <BuyTaskForm onClose={onClose} wallets={wallets} rpcGroups={rpcGroups} />
      )}
    </FormDataLoader>
  );
}

interface BuyTaskFormProps {
  onClose: () => void;
  wallets: Wallet[];
  rpcGroups: RPCGroupDashboardRow[];
}

function BuyTaskForm({ onClose, wallets, rpcGroups }: BuyTaskFormProps) {
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [tokenAddress, setTokenAddress] = useState('');
  const [buyAmount, setBuyAmount] = useState(BUY_AMOUNT_DEFAULT);
  const [buyFee, setBuyFee] = useState(BUY_FEE_DEFAULT);
  const [sellFee, setSellFee] = useState<number | undefined>(SELL_FEE_DEFAULT);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [selectedRpcGroup, setSelectedRpcGroup] = useState<RPCGroupDashboardRow | null>(null);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>([]);

  const wallet = selectedWallet ?? wallets[0];
  const rpcGroup = selectedRpcGroup ?? rpcGroups[0];

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const taskBody: BuyStrategyTaskPost = {
      trading_type: 'BUY',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      buy_amount: buyAmount,
      buy_fee: buyFee,
      sell_fee: sellFee,
      sell_strategies: sellStrategies,
      rpc_group_id: rpcGroup?.id,
    };

    postMutation.mutate(taskBody, { onSuccess: onClose });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="grid grid-cols-2 gap-4">
        <Card className="p-3">
          <h2>Buy Settings</h2>
          <div className="space-y-4">
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
            <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress} />
          </div>
        </Card>

        <Card className="p-3">
          <h2>Task Options</h2>
          <div className="grid grid-cols-[120px_120px_130px] gap-4">
            <SlippageEntry slippage={slippage} onChange={setSlippage} />
            <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
            <WalletSelector selectedWallet={wallet} onChange={setSelectedWallet} />
            <RPCGroupSelector selectedRpcGroup={rpcGroup} onChange={setSelectedRpcGroup} />
          </div>
        </Card>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <SellStrategies sellStrategies={sellStrategies} setSellStrategies={setSellStrategies} />
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
