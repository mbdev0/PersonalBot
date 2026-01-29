import { useState } from 'react';
import { TaskRowType, type Row } from '../types/tableRows';
import { ActionButtons } from './actionButtons';
import type { RowActions } from '../types/rowActions';
import { TableCell, TableRow } from '@/components/ui/table';

interface TableRowProps {
  tableRow: Row;
  rowActions: RowActions;
}

export function UnifiedTaskRow({ tableRow, rowActions }: TableRowProps) {
  const [isRunning, setIsRunning] = useState(false);
  return (
    <TableRow>
      <TableCell>{generateTableName(tableRow)}</TableCell>
      <TableCell>{generateRowType(tableRow)}</TableCell>
      <TableCell>cool w.s. message goes here</TableCell>
      <TableCell>{generateState(tableRow, isRunning)}</TableCell>
      <TableCell>
        <ActionButtons
          row={tableRow}
          isRunning={isRunning}
          setIsRunning={setIsRunning}
          rowActions={rowActions}
        ></ActionButtons>
      </TableCell>
    </TableRow>
  );
}

function generateTableName(row: Row): string {
  if (row.type == TaskRowType.Strategy) {
    return `Task - ${row.id}`;
  }

  const task = row.data;
  //if there's no strategy id then. we can say the task is a quick buy (straight to the task api)
  if (task.strategy_id == undefined) {
    const tokenAddress = `${row.data.token_address.slice(0, 5)}...${row.data.token_address.slice(-5)}`;
    return `Quick ${row.data.type} - ${tokenAddress}`;
  }

  return `BUY - ${row.data.token_address}`;
}

function generateState(row: Row, isRunning: boolean): string {
  if (row.type == TaskRowType.Task) {
    return row.data.state.task_state.toLowerCase();
  }

  if (isRunning) {
    return 'running';
  } else {
    return 'stopped';
  }
}

function generateRowType(row: Row): string {
  if (row.type == TaskRowType.Task) {
    return row.data.type;
  }
  return row.data.trading_type;
}
