import { useState } from 'react';
import type { Task } from '../../types/task';
import { useUpdateTask } from '../../hooks/useTasks';
import { useWallets } from '@/features/wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import type { Wallet } from '@/features/wallets/types/wallet';
import { TokenAddressEntry } from './fields/tokenAddress';
import { BuyEntry } from './fields/buy';
import { BUY_AMOUNT_DEFAULT, BUY_FEE_DEFAULT } from '../../utils/constants';
import { Button } from '@/components/ui/button';

interface BuyTaskEntryProps {
  task: Task;
  onClose: () => void;
}

export function BuyTaskEdit({ task, onClose }: BuyTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();

  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [selectedWallet, setWallet] = useState(task.wallet_name);

  const wallet = data?.find((w) => w.wallet_name === selectedWallet) ?? null;

  const handleWalletChange = (wallet: Wallet | null) => {
    if (wallet) {
      setWallet(wallet.wallet_name);
    }
  };

  const putMutation = useUpdateTask();

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody = {
      type: 'Buy',
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet.wallet_name ?? selectedWallet,
      token_address: tokenAddress,
      buy_amount: buyAmount,
      buy_fee: buyFee,
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
        <SlippageEntry slippage={slippage} onChange={setSlippage}></SlippageEntry>
        <ComputeUnitsEntry
          computeUnits={computeUnits}
          onChange={setComputeUnits}
        ></ComputeUnitsEntry>
      </div>

      <WalletSelector
        selectedWallet={wallet}
        onChange={handleWalletChange}
        isPending={isPending}
        isError={isError}
        data={data}
        error={error}
      ></WalletSelector>

      <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress}></TokenAddressEntry>

      <div className="grid grid-cols-2 gap-4">
        <BuyEntry
          buyAmount={buyAmount ?? BUY_AMOUNT_DEFAULT}
          onBuyAmountChange={setBuyAmount}
          buyFee={buyFee ?? BUY_FEE_DEFAULT}
          onBuyFeeChange={setBuyFee}
        ></BuyEntry>
      </div>

      <div className="flex justify-end gap-4 pt-4">
        <Button variant="outline" onClick={onClose}>
          Close
        </Button>

        <Button type="submit">Edit</Button>
      </div>
    </form>
  );
}
