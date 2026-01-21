import { useState } from 'react';
import { useWallets } from '../../wallets/hooks/useWallets';
import type { Wallet } from '../../wallets/types/wallet';
import { useAddTask } from '../hooks/useTasks';
import { WalletSelector } from './entry/fields/walletSelector';
import { SlippageEntry } from './entry/fields/slippage';
import { TokenAddressEntry } from './entry/fields/tokenAddress';
import { ComputeUnitsEntry } from './entry/fields/computeUnits';
import { BuyEntry } from './entry/fields/buy';
import { BuyTaskEntry } from './entry/buyTaskEntry';
import { SellTaskEntry } from './entry/sellTaskEntry';
import { AFKTaskEntry } from './entry/afkTaskEntry';

interface TaskEntryProps {
  onClose: () => void;
}

export function TaskEntry({ onClose }: TaskEntryProps) {
  const [taskType, setTaskType] = useState('Buy');
  return (
    <>
      <div className="task_entry_form">
        <div className="task_type">
          <h3>Task Type</h3>
          <select value={taskType} onChange={(e) => setTaskType(e.target.value)}>
            <option value="Buy">Buy</option>
            <option value="Sell">Sell</option>
            <option value="AFK">AFK</option>
          </select>
        </div>

        {taskType === 'Buy' && <BuyTaskEntry onClose={onClose}></BuyTaskEntry>}
        {taskType === 'Sell' && <SellTaskEntry onClose={onClose}></SellTaskEntry>}
        {taskType === 'AFK' && <AFKTaskEntry onClose={onClose}></AFKTaskEntry>}
      </div>
    </>
  );
}
