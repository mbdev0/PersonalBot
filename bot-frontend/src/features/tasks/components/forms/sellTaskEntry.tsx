import { useState } from 'react';
import { useAddTask } from '../../hooks/useTasks';
import { useWallets } from '../../../wallets/hooks/useWallets';
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

interface SellTaskEntryProps {
  onClose: () => void;
}

export function SellTaskEntry({ onClose }: SellTaskEntryProps) {
  const postMutation = useAddTask();
  const { data: wallets } = useWallets();

  const [slippage, setSlippage] = useState(SLIPPAGE_DEFAULT);
  const [computeUnits, setComputeUnits] = useState(COMPUTE_UNITS_DEFAULT);
  const [tokenAddress, setTokenAddress] = useState('');
  const [sellAmount, setSellAmount] = useState(SELL_AMOUNT_DEFAULT);
  const [sellFee, setSellFee] = useState(SELL_FEE_DEFAULT);
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);

  const effectiveWallet = selectedWallet ?? wallets?.[0] ?? null;

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!effectiveWallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody = {
      type: 'Sell',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: effectiveWallet.wallet_name,
      token_address: tokenAddress,
      sell_amount: sellAmount / 100,
      sell_fee: sellFee,
    };

    postMutation.mutate(taskBody);
    onClose();
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
          <div className="grid grid-cols-[120px_120px_130px] gap-4">
            <SlippageEntry slippage={slippage} onChange={setSlippage} />
            <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
            <WalletSelector selectedWallet={effectiveWallet} onChange={setSelectedWallet} />
          </div>
        </Card>
      </div>

      <div className="flex justify-end gap-3 pt-4">
        <Button type="button" variant="outline" onClick={onClose}>
          Cancel
        </Button>
        <Button type="submit" disabled={!effectiveWallet || !tokenAddress}>
          Create Task
        </Button>
      </div>
    </form>
  );
}
