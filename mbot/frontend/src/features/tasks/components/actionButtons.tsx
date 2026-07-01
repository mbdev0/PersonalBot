import type { RowActions } from '../types/rowActions';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import { StrategyTaskState } from '../types/strategies/strategyTask';
import { Play, Square, Pencil, Trash2, Copy, TerminalSquare } from 'lucide-react';
import { isFailure } from '../types/taskState';
import { QuickSellButtons } from './quickSellButtons';
import { cn } from '@/lib/utils';
import type { Settings } from '@/features/settings/types/settings';

interface ActionButtonProps {
  row: DisplayRow;
  rowActions: RowActions;
  settings?: Settings;
}

export function ActionButtons({ row, rowActions, settings }: ActionButtonProps) {
  const isDone = row.state === 'Done' || row.state === StrategyTaskState.success;

  const isRerunnable = isFailure(row.state);
  const isFailed = row.state === 'Task Failed' || row.state === StrategyTaskState.failed;

  const isCompletedBuyStrategy =
    isDone && row.type === TaskRowType.Strategy && row.data.trading_type === 'BUY';

  const isTerminal = isDone || isFailed;

  const isRunning =
    !isDone &&
    !isFailed &&
    row.state != StrategyTaskState.create &&
    row.state != StrategyTaskState.cancelled;

  return (
    <div className="grid-cols-1 space-y-2">
      <div className={cn('flex', isCompletedBuyStrategy ? 'justify-evenly' : 'gap-1.5')}>
        {isRunning && !isRerunnable ? (
          <button className="stop_task action-button-stop" onClick={() => rowActions.onStop(row)}>
            <Square className="action-icon" />
          </button>
        ) : (
          <button
            className={`start_task action-button-start ${
              isTerminal ? 'opacity-50 cursor-not-allowed' : ''
            }`}
            disabled={isTerminal}
            onClick={() => rowActions.onStart(row)}
          >
            <Play className="action-icon fill-current" />
          </button>
        )}

        <button className="edit action-button-edit" onClick={() => rowActions.onEdit(row)}>
          <Pencil className="action-icon" />
        </button>

        <button className="delete action-button-delete" onClick={() => rowActions.onDelete(row)}>
          <Trash2 className="action-icon" />
        </button>

        <button
          className="duplicate action-button-neutral"
          onClick={() => rowActions.onDuplicate(row)}
        >
          <Copy className="action-icon" />
        </button>

        <button
          className="terminal action-button-neutral"
          onClick={() => rowActions.onOpenTerminal(row)}
        >
          <TerminalSquare className="action-icon" />
        </button>
      </div>

      {isCompletedBuyStrategy && (
        <QuickSellButtons row={row} rowActions={rowActions} settings={settings} />
      )}
    </div>
  );
}
