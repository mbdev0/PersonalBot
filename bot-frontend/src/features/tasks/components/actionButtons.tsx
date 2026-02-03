import { Button } from '@/components/ui/button';
import type { RowActions } from '../types/rowActions';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import { ButtonGroup } from '@/components/ui/button-group';
import { useState } from 'react';
import { StrategyTaskState } from '../types/strategies/strategyTask';

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
    <>
      <ButtonGroup>
        <ButtonGroup>
          {isRunning ? (
            <Button
              className="stop_task"
              onClick={() => {
                setIsRunning(false);
                rowActions.onStop(row);
              }}
            >
              ⏹️
            </Button>
          ) : (
            <Button
              className="start_task"
              onClick={() => {
                setIsRunning(true);
                rowActions.onStart(row);
              }}
            >
              ▶️
            </Button>
          )}
        </ButtonGroup>
        <ButtonGroup>
          <Button
            className="edit"
            onClick={() => {
              rowActions.onEdit(row);
            }}
          >
            ✏️
          </Button>
        </ButtonGroup>
        <ButtonGroup>
          <Button className="delete" onClick={() => rowActions.onDelete(row)}>
            🗑️
          </Button>
        </ButtonGroup>

        <ButtonGroup>
          <Button className="duplicate" onClick={() => rowActions.onDuplicate(row)}>
            ⿻
          </Button>
        </ButtonGroup>
      </ButtonGroup>
    </>
  );
}
