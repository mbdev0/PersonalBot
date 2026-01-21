import { useEffect, useState } from 'react';
import { useWallets } from '../../wallets/hooks/useWallets';
import type { Wallet } from '../../wallets/types/wallet';
import { useAddTask } from '../hooks/useTasks';
import { WalletSelector } from './entry/fields/walletSelector';
import { SlippageEntry } from './entry/fields/slippage';
import { TokenAddressEntry } from './entry/fields/tokenAddress';
import { ComputeUnitsEntry } from './entry/fields/computeUnits';
import { BuyEntry } from './entry/fields/buy';

interface TaskEntryProps {
  onClose: () => void;
}

export function TaskEntry({ onClose }: TaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();

  const [taskType, setTaskType] = useState('Buy');
  const [slippage, setSlippage] = useState('20');
  const [computeUnits, setComputeUnits] = useState('100000');
  const [tokenAddress, setTokenAddress] = useState('');
  const [buyAmount, setBuyAmount] = useState('1.0');
  const [buyFee, setBuyFee] = useState('0.1');
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

    const body_base = {
      type: taskType,
      slippage: parseFloat(slippage) / 100,
      compute_units: parseInt(computeUnits),
      wallet_name: wallet.wallet_name,
      token_address: tokenAddress,
    };

    const taskBody =
      taskType === 'Buy'
        ? {
            ...body_base,
            buy_amount: parseFloat(buyAmount),
            buy_fee: parseFloat(buyFee),
          }
        : {
            ...body_base,
            sell_amount: parseFloat(sellAmount) / 100,
            sell_fee: parseFloat(sellFee),
          };

    postMutation.mutate({
      ...taskBody,
    });

    onClose();
  };

  return (
    <>
      <div className="task_entry_form">
        <div className="task_type">
          <h3>Task Type</h3>
          <select value={taskType} onChange={(e) => setTaskType(e.target.value)}>
            <option value="Buy">Buy</option>
            <option value="Sell">Sell</option>
          </select>
        </div>

        <SlippageEntry onChange={setSlippage} slippage={slippage}></SlippageEntry>

        <ComputeUnitsEntry
          onChange={setComputeUnits}
          computeUnits={computeUnits}
        ></ComputeUnitsEntry>

        <WalletSelector
          selectedWallet={wallet}
          onChange={setSelectedWallet}
          isError={isError}
          isPending={isPending}
          error={error}
          data={data}
        ></WalletSelector>

        <TokenAddressEntry onChange={setTokenAddress}></TokenAddressEntry>

        {taskType === 'Sell' && (
          <div className="sell_settings">
            <div className="sell_amount">
              <h3>Sell Amount</h3>
              <input
                type="text"
                name="sell_amount"
                id="sell_amount"
                placeholder="Sell Amount"
                value={sellAmount}
                onChange={(e) => setSellAmount(e.target.value)}
              />
            </div>

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
          </div>
        )}

        {taskType === 'Buy' && (
          <BuyEntry
            onBuyAmountChange={setBuyAmount}
            buyAmount={buyAmount}
            onBuyFeeChange={setBuyFee}
            buyFee={buyFee}
          ></BuyEntry>
        )}

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
    </>
  );
}
