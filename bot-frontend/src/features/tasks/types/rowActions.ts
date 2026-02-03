import type { DisplayRow } from './tableRows';

interface RowActions {
  onStart: (row: DisplayRow) => void;
  onStop: (row: DisplayRow) => void;
  onEdit: (row: DisplayRow) => void;
  onDelete: (row: DisplayRow) => void;
  onDuplicate: (row: DisplayRow) => void;
}

export { type RowActions };
