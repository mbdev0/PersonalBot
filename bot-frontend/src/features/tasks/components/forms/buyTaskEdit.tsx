import { useState } from 'react';
import { useWallets } from '@/features/wallets/hooks/useWallets';
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

interface BuyTaskEditProps {
  task: BuyStrategyTask;
  onClose: () => void;
}

export function BuyTaskEdit({ task, onClose }: BuyTaskEditProps) {
  const { isPending, isError, data, error } = useWallets();
  const putMutation = useUpdateStrategy();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [selectedWallet, setWallet] = useState(task.wallet_name);
  const [sellFee, setSellFee] = useState<number | undefined>(task.sell_fee);
  const [sellStrategies, setSellStrategies] = useState<SellStrategy[]>(task.sell_strategies);

  const wallet = data?.find((w) => w.wallet_name === selectedWallet) ?? null;

  const handleWalletChange = (wallet: Wallet | null) => {
    if (wallet) {
      setWallet(wallet.wallet_name);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody: BuyStrategyTaskPut = {
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
    };

    putMutation.mutate({ ...taskBody, id: task.id });
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
          <div className="grid grid-cols-[120px_120px_130px] gap-4">
            <SlippageEntry slippage={slippage} onChange={setSlippage} />
            <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
            <WalletSelector selectedWallet={wallet} onChange={handleWalletChange} />
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
          Update Task
        </Button>
      </div>
    </form>
  );
}
