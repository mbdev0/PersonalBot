import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useAddTask } from '../../hooks/useTasks';
import { useWallets } from '../../../wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { TokenAddressEntry } from './fields/tokenAddress';
import { BuyEntry } from './fields/buy';

interface BuyTaskEntryProps {
  onClose: () => void;
}

export function BuyTaskEntry({ onClose }: BuyTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();

  const [slippage, setSlippage] = useState('20');
  const [computeUnits, setComputeUnits] = useState('100000');
  const [tokenAddress, setTokenAddress] = useState('');
  const [buyAmount, setBuyAmount] = useState('1.0');
  const [buyFee, setBuyFee] = useState('0.1');
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);
  const postMutation = useAddTask();

  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = () => {
    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const taskBody = {
      type: 'Buy',
      slippage: parseFloat(slippage) / 100,
      compute_units: parseInt(computeUnits),
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
      buy_amount: parseFloat(buyAmount),
      buy_fee: parseFloat(buyFee),
    };

    postMutation.mutate({
      ...taskBody,
    });

    onClose();
  };

  return (
    <div className="buyEntry">
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

      <BuyEntry
        buyAmount={buyAmount}
        onBuyAmountChange={setBuyAmount}
        buyFee={buyFee}
        onBuyFeeChange={setBuyFee}
      ></BuyEntry>

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
