import { useState } from 'react';
import type { Wallet } from '../../../wallets/types/wallet';
import { useAddStrategy } from '../../hooks/useStrategy';
import { useWallets } from '../../../wallets/hooks/useWallets';
import { SlippageEntry } from './fields/slippage';
import { ComputeUnitsEntry } from './fields/computeUnits';
import { WalletSelector } from './fields/walletSelector';
import { BuyEntry } from './fields/buy';
import type { StrategyTaskPost } from '../../types/strategyTask';

interface AFKTaskEntryProps {
  onClose: () => void;
}

export function AFKTaskEntry({ onClose }: AFKTaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();

  const [slippage, setSlippage] = useState('20');
  const [computeUnits, setComputeUnits] = useState('100000');
  const [buyAmount, setBuyAmount] = useState('1.0');
  const [buyFee, setBuyFee] = useState('0.1');
  const [sellFee, setSellFee] = useState('0.1');
  const [selectedWallet, setSelectedWallet] = useState<Wallet | null>(null);

  const postMutation = useAddStrategy();
  const wallet = selectedWallet ?? data?.[0] ?? null;

  const handleSubmit = () => {
    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const strategyBody: StrategyTaskPost = {
      tradingType: 'AFK',
      slippage: parseFloat(slippage) / 100,
      compute_units: parseInt(computeUnits),
      wallet_name: wallet.wallet_name,
      buy_amount: parseFloat(buyAmount),
      buy_fee: parseFloat(buyFee),
      sell_fee: parseFloat(sellFee),
      sell_strategies: [
        {
          type: 'take_profit_percentage',
          value: 0.2,
          sell_amount: 1,
        },
      ],
      filters: {
        has_twitter: true,
      },
    };

    postMutation.mutate({
      ...strategyBody,
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

      <BuyEntry
        buyAmount={buyAmount}
        onBuyAmountChange={setBuyAmount}
        buyFee={buyFee}
        onBuyFeeChange={setBuyFee}
      ></BuyEntry>

      <div className="sell_fee">
        <h3>Sell Fee</h3>
        <input
          type="text"
          name="sell_fee"
          id="sell_fee"
          placeholder="Sell Fee"
          value={sellFee}
          onChange={(e) => setSellFee(e.target.value)}
        />
      </div>

      <button
        onClick={() => {
          handleSubmit();
        }}
        disabled={!wallet}
      >
        Create Task
      </button>
      <button onClick={onClose}>Close</button>
    </div>
  );
}
