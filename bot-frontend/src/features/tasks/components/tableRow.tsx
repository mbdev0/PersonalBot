import { useState } from 'react';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import { ActionButtons } from './actionButtons';
import type { RowActions } from '../types/rowActions';
import { TableCell, TableRow } from '@/components/ui/table';

interface TableRowProps {
  tableRow: DisplayRow;
  rowActions: RowActions;
}

export function UnifiedTaskRow({ tableRow, rowActions }: TableRowProps) {
  return (
    <TableRow>
      <TableCell>{generateTableName(tableRow)}</TableCell>
      <TableCell>{generateRowType(tableRow)}</TableCell>
      <TableCell>{tableRow.wsMessage}</TableCell>
      <TableCell>{tableRow.state}</TableCell>
      <TableCell>
        <ActionButtons row={tableRow} rowActions={rowActions}></ActionButtons>
      </TableCell>
    </TableRow>
  );
}

function generateTableName(row: DisplayRow): string {
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

function generateRowType(row: DisplayRow): string {
  if (row.type == TaskRowType.Task) {
    return row.data.type;
  }
  return row.data.trading_type;
}
