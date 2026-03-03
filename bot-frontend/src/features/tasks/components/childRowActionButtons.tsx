import type { RowActions } from '../types/rowActions';
import { type DisplayRow } from '../types/tableRows';
import { Play, Square, Trash2 } from 'lucide-react';
import { isFailure, TaskState, isRunning } from '../types/taskState';

export interface ActionButtonProps {
  row: DisplayRow;
  rowActions: RowActions;
}

export function ChildRowActionButtons({ row, rowActions }: ActionButtonProps) {
  const isDone = row.state === TaskState.TASK_DONE;
  const isFailed = isFailure(row.state);
  const isTerminal = isDone || isFailed;

  return (
    <div className="flex gap-1.5">
      {isRunning(row.state) ? (
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

      <button className="delete action-button-delete" onClick={() => rowActions.onDelete(row)}>
        <Trash2 className="action-icon" />
      </button>
    </div>
  );
}
