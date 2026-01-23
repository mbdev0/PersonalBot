import { Button } from '@/components/ui/button';
import type { RowActions } from '../types/rowActions';
import type { Row } from '../types/tableRows';
import { ButtonGroup } from '@/components/ui/button-group';

export interface ActionButtonProps {
  row: Row;
  isRunning: boolean;
  setIsRunning: (arg: boolean) => void;
  rowActions: RowActions;
}

export function ActionButtons({ row, isRunning, setIsRunning, rowActions }: ActionButtonProps) {
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
