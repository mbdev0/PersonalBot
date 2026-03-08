import type { RowActions } from '../types/rowActions';
import { type DisplayRow } from '../types/tableRows';
import { StrategyTaskState } from '../types/strategies/strategyTask';
import { Play, Square, Pencil, Trash2, Copy } from 'lucide-react';
import { isFailure } from '../types/taskState';

export interface ActionButtonProps {
  row: DisplayRow;
  rowActions: RowActions;
}

export function ActionButtons({ row, rowActions }: ActionButtonProps) {
  const isDone = row.state === 'Done' || row.state === StrategyTaskState.success;

  const isRerunnable = isFailure(row.state);
  const isFailed = row.state === 'Task Failed' || row.state === StrategyTaskState.failed;

  const isTerminal = isDone || isFailed;

  const isRunning =
    !isDone &&
    !isFailed &&
    row.state != StrategyTaskState.create &&
    row.state != StrategyTaskState.cancelled;

  return (
    <div className="flex gap-1.5">
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
    </div>
  );
}
