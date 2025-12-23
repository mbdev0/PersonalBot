import { useEffect, useState } from 'react';
import { useWallets } from '../../wallets/hooks/useWallets';
import type { WalletDto } from '../../wallets/types/wallet';
import { useAddTask } from '../hooks/useTasks';

interface TaskEntryProps {
  onClose: () => void;
}

export function TaskEntry({ onClose }: TaskEntryProps) {
  const { isPending, isError, data, error } = useWallets();

  const [taskType, setTaskType] = useState('Buy');
  const [slippage, setSlippage] = useState('20');
  const [computeUnits, setComputeUnits] = useState('100000');
  const [tokenAddress, setTokenAddress] = useState('2gvhTGcWJFUmPcXyUWM1KxbaqQKryMPKjiZHNNg7pump');
  const [buyAmount, setBuyAmount] = useState('1.0');
  const [buyFee, setBuyFee] = useState('0.1');
  const [sellAmount, setSellAmount] = useState('20');
  const [sellFee, setSellFee] = useState('0.1');
  const [wallet, setWallet] = useState<WalletDto | null>(null);
  const postMutation = useAddTask();

  useEffect(() => {
    if (data && data.length > 0 && !wallet) {
      setWallet(data[0]);
    }
  }, [data, wallet]);

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
      {isPending && <div>Loading wallets...</div>}
      {isError && <div>Error loading wallets: {error.message}</div>}
      {!isPending && data?.length === 0 && <div>No wallets available. Create a wallet first!</div>}

      <div className="task_entry_form">
        <div className="task_type">
          <h3>Task Type</h3>
          <select value={taskType} onChange={(e) => setTaskType(e.target.value)}>
            <option value="Buy">Buy</option>
            <option value="Sell">Sell</option>
          </select>
        </div>

        <div className="slippage">
          <h3>Slippage</h3>
          <input
            type="text"
            name="slippage"
            id="slippage"
            placeholder="Slippage"
            value={slippage}
            onChange={(e) => setSlippage(e.target.value)}
          />
        </div>

        <div className="compute_units">
          <h3>Compute Units</h3>
          <input
            type="text"
            name="compute_units"
            id="compute_units"
            placeholder="Compute Units"
            value={computeUnits}
            onChange={(e) => setComputeUnits(e.target.value)}
          />
        </div>

        <div className="wallet">
          <h3>Wallet</h3>
          <select
            value={wallet?.wallet_name}
            disabled={isPending || isError}
            onChange={(e) => setWallet(data?.find((w) => w.wallet_name === e.target.value) ?? null)}
          >
            {data?.map((key, _) => (
              <option value={key.wallet_name}>{key.wallet_name}</option>
            ))}
          </select>
        </div>

        <div className="token_address">
          <h3>Token Address</h3>
          <input
            type="text"
            name="token_address"
            id="token_address"
            placeholder="Token Address"
            onChange={(e) => setTokenAddress(e.target.value)}
          />
        </div>

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
          <div className="buy_settings">
            <div className="buy_amount">
              <h3>Buy Amount</h3>
              <input
                type="text"
                name="buy_amount"
                id="buy_amount"
                placeholder="Buy Amount"
                value={buyAmount}
                onChange={(e) => setBuyAmount(e.target.value)}
              />
            </div>

            <div className="buy_fee">
              <h3>Buy Fee</h3>
              <input
                type="text"
                name="buy_fee"
                id="buy_fee"
                placeholder="Buy Fee"
                value={buyFee}
                onChange={(e) => setBuyFee(e.target.value)}
              />
            </div>
          </div>
        )}

        <button
          onClick={() => {
            handleSubmit();
          }}
          disabled={!wallet || isPending}
        >
          Create Task
        </button>
        <button onClick={onClose}>Close</button>
      </div>
    </>
  );
}
