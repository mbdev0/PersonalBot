import type { DisplayRow } from './tableRows';
import type { TaskPostDto } from './task';

interface RowActions {
  onStart: (row: DisplayRow) => void;
  onStop: (row: DisplayRow) => void;
  onEdit: (row: DisplayRow) => void;
  onDelete: (row: DisplayRow) => void;
  onDuplicate: (row: DisplayRow) => void;
  onQuickSell: (task: TaskPostDto) => void;
  onOpenTerminal: (row: DisplayRow) => void;
}

export { type RowActions };
