//create form
//we will need a state for each of these?
//what will it do:
//  type of task
//  slippage
//  compute units
//  wallet_addy_priv_key
//  token_addy
//  buy_amount
//  buy_fee

import { useEffect, useState } from 'react';
import { useWallets } from '../../wallets/hooks/useWallets';
import type { WalletDto } from '../../wallets/types/wallet';
import { postTasksMutation } from '../hooks/useTasks';

export function TaskEntry({ onClose }: { onClose: () => void }) {
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

  useEffect(() => {
    if (data && data.length > 0) {
      setWallet(data[0]);
    }
  }, [data]);

  const postMutation = postTasksMutation();

  //TODO: should change to a struct? or interface?
  const taskBody =
    taskType === 'Buy'
      ? {
          type: taskType,
          slippage: parseFloat(slippage) / 100,
          compute_units: parseInt(computeUnits),
          wallet_name: wallet?.wallet_name,
          token_address: tokenAddress,
          buy_amount: parseFloat(buyAmount),
          buy_fee: parseFloat(buyFee),
        }
      : {
          type: taskType,
          slippage: parseFloat(slippage) / 100,
          compute_units: parseInt(computeUnits),
          wallet_name: wallet?.wallet_name,
          token_address: tokenAddress,
          sell_amount: parseFloat(sellAmount) / 100,
          sell_fee: parseFloat(sellFee),
        };

  return (
    <>
      {isPending && <div>Loading wallets...</div>}
      {isError && <div>Error loading wallets: {error.message}</div>}
      {!isPending && data?.length === 0 && <div>No wallets available. Create a wallet first!</div>}

      {data && data.length > 0 && (
        <>
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
              onChange={(e) =>
                setWallet(data?.find((w) => w.wallet_name === e.target.value) ?? null)
              }
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
              if (!wallet) {
                alert('Please select a wallet');
                return;
              }

              postMutation.mutate({
                ...taskBody,
                wallet_name: wallet.wallet_name,
              });

              onClose();
            }}
            disabled={!wallet || isPending}
          >
            Create Task
          </button>
          <button onClick={onClose}>Close</button>
        </>
      )}
    </>
  );
}

//WORKS!
async function create_task_api_call(task_body: any) {
  // console.log(private_key, token_address);
  const url = 'http://localhost:9090/api/tasks/create';

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(task_body),
  });

  if (!response.ok) {
    console.log('Error whilst sending request:', response.status);
    console.log(response);
    return;
  }
  console.log('response: ', response.json);
}
