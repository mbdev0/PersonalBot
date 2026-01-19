import type { RowActions } from '../types/rowActions';
import type { Row } from '../types/tableRows';

export interface ActionButtonProps {
  row: Row;
  isRunning: boolean;
  setIsRunning: (arg: boolean) => void;
  rowActions: RowActions;
}

export function ActionButtons({ row, isRunning, setIsRunning, rowActions }: ActionButtonProps) {
  return (
    <>
      {isRunning ? (
        <button
          className="stop_task"
          onClick={() => {
            setIsRunning(false);
            rowActions.onStop(row);
          }}
        >
          ⏹️
        </button>
      ) : (
        <button
          className="start_task"
          onClick={() => {
            setIsRunning(true);
            rowActions.onStart(row);
          }}
        >
          ▶️
        </button>
      )}

      <button
        className="edit"
        onClick={() => {
          rowActions.onEdit(row);
        }}
      >
        ✏️
      </button>
      <button className="delete" onClick={() => rowActions.onDelete(row)}>
        🗑️
      </button>
      <button className="duplicate" onClick={() => rowActions.onDuplicate(row)}>
        ⿻
      </button>
    </>
  );
}
