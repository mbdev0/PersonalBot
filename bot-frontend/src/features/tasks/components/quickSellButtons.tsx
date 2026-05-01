import type { Settings } from '@/features/settings/types/settings';
import type { RowActions } from '../types/rowActions';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import type { TaskPostDto } from '../types/task';

interface QuickSellButtonProps {
  row: DisplayRow;
  rowActions: RowActions;
  settings?: Settings;
}

export function QuickSellButtons({ row, rowActions, settings }: QuickSellButtonProps) {
  const button1 = settings?.quick_sell_buttons.button_1 ?? 0.25;
  const button2 = settings?.quick_sell_buttons.button_2 ?? 0.5;
  const button3 = settings?.quick_sell_buttons.button_3 ?? 0.75;
  const button4 = settings?.quick_sell_buttons.button_4 ?? 1;

  return (
    <div className="flex gap-1.5">
      <button
        className="quicksell quick-sell-button-base  hover:bg-green-400/10 hover:ring-green-400/25"
        onClick={() => spawnAndRunTask(row, button1, rowActions)}
      >
        {button1 * 100}%
      </button>
      <button
        className="quicksell quick-sell-button-base  hover:bg-green-500/10 hover:ring-green-500/25"
        onClick={() => spawnAndRunTask(row, button2, rowActions)}
      >
        {button2 * 100}%
      </button>
      <button
        className="quicksell quick-sell-button-base  hover:bg-green-600/10 hover:ring-green-600/25"
        onClick={() => spawnAndRunTask(row, button3, rowActions)}
      >
        {button3 * 100}%
      </button>
      <button
        className="quicksell quick-sell-button-base  hover:bg-green-700/10 hover:ring-green-700/25"
        onClick={() => spawnAndRunTask(row, button4, rowActions)}
      >
        {button4 * 100}%
      </button>
    </div>
  );
}

function spawnAndRunTask(row: DisplayRow, amountToSell: number, rowActions: RowActions) {
  if (
    !(
      (row.type === TaskRowType.Strategy && row.data.trading_type === 'BUY') ||
      (row.type === TaskRowType.Task && row.data.type === 'Buy')
    )
  ) {
    return;
  }

  console.log(row.data.sell_fee);
  console.log(row);

  const baseTask = {
    type: 'Sell' as const,
    slippage: row.data.slippage,
    compute_units: row.data.compute_units,
    wallet_name: row.data.wallet_name,
    sell_amount: amountToSell,
    sell_fee: row.data.sell_fee,
    rpc_group_id: row.data.rpc_group_id,
  };

  if (row.type === TaskRowType.Strategy && row.data.trading_type === 'BUY') {
    //create sell task in here
    const sellTask: TaskPostDto = {
      ...baseTask,
      token_address: row.data.token_address,
      strategy_id: row.id,
      sell_position_id: row.data.buy_task_id,
    };
    rowActions.onQuickSell(sellTask);
    return;
  }

  if (row.type === TaskRowType.Task && row.data.type === 'Buy') {
    const sellTask: TaskPostDto = {
      ...baseTask,
      token_address: row.data.token_address,
      strategy_id: row.data.strategy_id || 0,
      sell_position_id: row.id,
    };
    rowActions.onQuickSell(sellTask);
    return;
  }
}
