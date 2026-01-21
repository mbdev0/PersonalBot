import { useState } from 'react';
import { useAddTask } from '../../hooks/useTasks';
import { useWallets } from '../../../wallets/hooks/useWallets';
import type { Wallet } from '../../../wallets/types/wallet';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { SlippageEntry } from './fields/slippage';
import { TokenAddressEntry } from './fields/tokenAddress';
import { WalletSelector } from './fields/walletSelector';
import { SellEntry } from './fields/sell';

interface SellTaskEntryProps {
  onClose: () => void;
}

export function SellTaskEntry({ onClose }: SellTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();

  const [slippage, setSlippage] = useState('20');
  const [computeUnits, setComputeUnits] = useState('100000');
  const [tokenAddress, setTokenAddress] = useState('');
  const [sellAmount, setSellAmount] = useState('20');
  const [sellFee, setSellFee] = useState('0.1');
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const postMutation = useAddTask();

  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = () => {
    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody = {
      type: 'Sell',
      slippage: parseFloat(slippage) / 100,
      compute_units: parseInt(computeUnits),
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      sell_amount: parseFloat(sellAmount) / 100,
      sell_fee: parseFloat(sellFee),
    };

    postMutation.mutate({
      ...taskBody,
    });

    onClose();
  };

  return (
    <div className="sellEntry">
      <SlippageEntry onChange={setSlippage} slippage={slippage}></SlippageEntry>

      <ComputeUnitsEntry onChange={setComputeUnits} computeUnits={computeUnits}></ComputeUnitsEntry>

      <WalletSelector
        selectedWallet={wallet}
        onChange={setSelectedWallet}
        isError={isError}
        isPending={isPending}
        error={error}
        data={data}
      ></WalletSelector>

      <TokenAddressEntry value={tokenAddress} onChange={setTokenAddress}></TokenAddressEntry>

      <SellEntry
        sellAmount={sellAmount}
        onSellAmountChange={setSellAmount}
        sellFee={sellFee}
        onSellFeeChange={setSellFee}
      ></SellEntry>

      <button
        onClick={() => {
          handleSubmit();
        }}
        disabled={!wallet || !tokenAddress}
      >
        Create Task
      </button>
      <button onClick={onClose}>Close</button>
    </div>
  );
}
