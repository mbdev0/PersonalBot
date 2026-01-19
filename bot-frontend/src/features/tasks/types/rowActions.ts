import type { Row } from './tableRows';

interface RowActions {
  onStart: (row: Row) => void;
  onStop: (row: Row) => void;
  onEdit: (row: Row) => void;
  onDelete: (row: Row) => void;
  onDuplicate: (row: Row) => void;
}

export { type RowActions };
