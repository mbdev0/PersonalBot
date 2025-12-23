import { useEffect, useState } from 'react';
import { useWallets } from '../../wallets/hooks/useWallets';
import type { Task } from '../types/task';
import { useUpdateTask } from '../hooks/useTasks';

interface TaskUpdateProps {
  task: Task;
  onClose: () => void;
}

export function TaskUpdate({ task, onClose }: TaskUpdateProps) {
  const { isPending, isError, data, error } = useWallets();

  const [taskType, setTaskType] = useState(task.type);
  const [slippage, setSlippage] = useState(task.slippage * 100);
  const [computeUnits, setComputeUnits] = useState(task.compute_units);
  const [tokenAddress, setTokenAddress] = useState(task.token_address);
  const [buyAmount, setBuyAmount] = useState(task.buy_amount);
  const [buyFee, setBuyFee] = useState(task.buy_fee);
  const [sellAmount, setSellAmount] = useState(task.sell_amount);
  const [sellFee, setSellFee] = useState(task.sell_fee);
  const [wallet, setWallet] = useState(task.wallet_name);
  const putMutation = useUpdateTask();

  useEffect(() => {
    if (data && data.length > 0 && !wallet) {
      setWallet(data[0].wallet_name);
    }
  }, [data, wallet]);

  const handleSubmit = () => {
    if (!wallet) {
      alert('Please select a wallet');
      return;
    }

    const body_base = {
      type: taskType,
      slippage: slippage / 100,
      compute_units: computeUnits,
      wallet_name: wallet,
      token_address: tokenAddress,
    };

    const taskBody =
      taskType === 'Buy'
        ? {
            ...body_base,
            buy_amount: buyAmount,
            buy_fee: buyFee,
          }
        : {
            ...body_base,
            sell_amount: sellAmount != null ? sellAmount / 100 : 0,
            sell_fee: sellFee,
          };

    putMutation.mutate({
      ...taskBody,
      id: task.task_id,
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
            type="number"
            name="slippage"
            id="slippage"
            placeholder="Slippage"
            value={slippage}
            onChange={(e) => setSlippage(e.target.valueAsNumber)}
          />
        </div>

        <div className="compute_units">
          <h3>Compute Units</h3>
          <input
            type="number"
            name="compute_units"
            id="compute_units"
            placeholder="Compute Units"
            value={computeUnits}
            onChange={(e) => setComputeUnits(e.target.valueAsNumber)}
          />
        </div>

        <div className="wallet">
          <h3>Wallet</h3>
          <select
            value={wallet}
            disabled={isPending || isError}
            onChange={(e) => setWallet(e.target.value)}
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
            value={tokenAddress}
            onChange={(e) => setTokenAddress(e.target.value)}
          />
        </div>

        {taskType === 'Sell' && (
          <div className="sell_settings">
            <div className="sell_amount">
              <h3>Sell Amount</h3>
              <input
                type="number"
                name="sell_amount"
                id="sell_amount"
                placeholder="Sell Amount"
                value={sellAmount}
                onChange={(e) => setSellAmount(e.target.valueAsNumber)}
              />
            </div>

            <div className="sell_fee">
              <h3>Sell Fee</h3>
              <input
                type="number"
                name="sell_fee"
                id="sell_fee"
                placeholder="Sell Fee"
                value={sellFee}
                onChange={(e) => setSellFee(e.target.valueAsNumber)}
              />
            </div>
          </div>
        )}

        {taskType === 'Buy' && (
          <div className="buy_settings">
            <div className="buy_amount">
              <h3>Buy Amount</h3>
              <input
                type="number"
                name="buy_amount"
                id="buy_amount"
                placeholder="Buy Amount"
                value={buyAmount}
                onChange={(e) => setBuyAmount(e.target.valueAsNumber)}
              />
            </div>

            <div className="buy_fee">
              <h3>Buy Fee</h3>
              <input
                type="number"
                name="buy_fee"
                id="buy_fee"
                placeholder="Buy Fee"
                value={buyFee}
                onChange={(e) => setBuyFee(e.target.valueAsNumber)}
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
          Update Task
        </button>
        <button onClick={onClose}>Close</button>
      </div>
    </>
  );
}
