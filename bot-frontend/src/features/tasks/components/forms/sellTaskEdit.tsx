import { useState } from 'react';
import type { Task } from '../../types/task';
import { useUpdateTask } from '../../hooks/useTasks';
import { useWallets } from '@/features/wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import type { Wallet } from '@/features/wallets/types/wallet';
import { TokenAddressEntry } from './fields/tokenAddress';
import { SELL_AMOUNT_DEFAULT, SELL_FEE_DEFAULT } from '../../utils/constants';
import { Button } from '@/components/ui/button';
import { Card } from '@/components/ui/card';
import { SellEntry } from './fields/sell';

interface SellTaskEditProps {
  task: Task;
  onClose: () => void;
}

export function SellTaskEdit({ task, onClose }: SellTaskEditProps) {
  const { isPending, isError, data, error } = useWallets();
  const putMutation = useUpdateTask();

  if (!task.sell_amount) {
    return (
      <div className="text-center py-8 text-red-600">
        Error: task.sell_amount was undefined - double check the logic
      </div>
    );
  }

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [sellAmount, setSellAmount] = useState(task.sell_amount * 100);
  const [sellFee, setSellFee] = useState(task.sell_fee);
  const [selectedWallet, setWallet] = useState(task.wallet_name);

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

    const taskBody = {
      type: 'Sell',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      sell_amount: (sellAmount ?? SELL_AMOUNT_DEFAULT) / 100, // Convert back to decimal
      sell_fee: sellFee ?? SELL_FEE_DEFAULT,
    };

    putMutation.mutate({ ...taskBody, id: task.task_id });
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
          <div className="grid grid-cols-[120px_120px_130px] gap-4">
            <SlippageEntry slippage={slippage} onChange={setSlippage} />
            <ComputeUnitsEntry computeUnits={computeUnits} onChange={setComputeUnits} />
            <WalletSelector
              selectedWallet={wallet}
              onChange={handleWalletChange}
              isPending={isPending}
              isError={isError}
              data={data}
              error={error}
            />
          </div>
        </Card>
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
