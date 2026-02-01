import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useWallets } from '../../../wallets/hooks/useWallets';
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

interface BuyTaskEntryProps {
  onClose: () => void;
}

export function BuyTaskEntry({ onClose }: BuyTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();
  const postMutation = useAddStrategy();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [tokenAddress, setTokenAddress] = useState('');
  const [buyAmount, setBuyAmount] = useState(BUY_AMOUNT_DEFAULT);
  const [buyFee, setBuyFee] = useState(BUY_FEE_DEFAULT);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const [sellFee, setSellFee] = useState<number | undefined>(SELL_FEE_DEFAULT);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>([]);

  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

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
    };

    postMutation.mutate(taskBody);
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
          </div>
        </Card>
      </div>

      <div className="grid grid-cols-2 gap-4">
        <SellStrategies sellStrategies={sellStrategies} setSellStrategies={setSellStrategies} />
      </div>

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
