import type { RowActions } from '../types/rowActions';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import { useState } from 'react';
import { StrategyTaskState } from '../types/strategies/strategyTask';
import { Play, Square, Pencil, Trash2, Copy } from 'lucide-react';

export interface ActionButtonProps {
  row: DisplayRow;
  rowActions: RowActions;
}

export function ActionButtons({ row, rowActions }: ActionButtonProps) {
  const [isRunning, setIsRunning] = useState(
    row.type === TaskRowType.Task
      ? row.state === 'TaskRun'
      : row.state === StrategyTaskState.running
  );

  return (
    <div className="flex gap-1.5">
      {isRunning ? (
        <button
          className="stop_task action-button-stop"
          onClick={() => {
            setIsRunning(false);
            rowActions.onStop(row);
          }}
        >
          <Square className="action-icon" />
        </button>
      ) : (
        <button
          className="start_task action-button-start"
          onClick={() => {
            setIsRunning(true);
            rowActions.onStart(row);
          }}
        >
          <Play className="action-icon fill-current" />
        </button>
      )}
      <button
        className="edit action-button-edit"
        onClick={() => {
          rowActions.onEdit(row);
        }}
      >
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
