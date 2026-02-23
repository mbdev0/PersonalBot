import { type ColumnDef, type RowData } from '@tanstack/react-table';
import { ActionButtons } from '../components/actionButtons';
import { TaskRowType, type DisplayRow } from './tableRows';
import type { RowActions } from './rowActions';
import { formatAddress } from '@/utils/crypto/address_shortner';
import { ArrowRight, ArrowDown } from 'lucide-react';

//this extends the table meta to add a field called row actions so TS doesn't complain
declare module '@tanstack/react-table' {
  interface TableMeta<TData extends RowData = RowData> {
    rowActions: RowActions;
  }
}

export const columns: ColumnDef<DisplayRow>[] = [
  {
    id: 'expand',
    header: '',
    size: 10,
    cell: ({ row }) => {
      return row.getCanExpand() ? (
        <button
          onClick={row.getToggleExpandedHandler()}
          className="h-6 w-6 flex items-center justify-center rounded-md hover:bg-foreground/8 transition-[background-color,box-shadow] duration-200 ring-1 ring-foreground/5 hover:ring-foreground/15"
        >
          {row.getIsExpanded() ? (
            <ArrowDown className="h-3.5 w-3.5 text-foreground/60" />
          ) : (
            <ArrowRight className="h-3.5 w-3.5 text-foreground/60" />
          )}
        </button>
      ) : (
        ''
      );
    },
  },
  {
    accessorKey: 'task_type',
    header: 'Task Type',
    cell: ({ row }) => {
      {
        return row.original.type === TaskRowType.Task
          ? row.original.data.type
          : row.original.data.trading_type;
      }
    },
  },
  {
    accessorKey: 'task_name',

    header: 'Task Name',
    cell: ({ row }) => {
      if (row.original.type === TaskRowType.Task) {
        if (row.original.strategyId) {
          return `${row.original.data.type} ${formatAddress(row.original.data.token_address)}`;
        }
      } else {
        if (row.original.data.trading_type === 'AFK') {
          return `Task - ${row.original.data.trading_type} ${row.original.id}`;
        }
        return `Quick ${row.original.data.trading_type} - ${formatAddress(row.original.data.token_address)}`;
      }
    },
  },
  {
    accessorKey: 'message',
    header: 'Message',
    cell: ({ row }) => {
      return (
        <div className="whitespace-normal break-all text-left px-2 max-w-xs">
          {row.original.ws_message}
        </div>
      );
    },
  },
  {
    accessorKey: 'status',
    header: 'Status',
    cell: ({ row }) => {
      return row.original.state;
    },
  },
  {
    accessorKey: 'actions',
    header: 'Actions',
    cell: ({ row, table }) => {
      return (
        <div className="flex justify-center">
          {row.original.type === TaskRowType.Strategy && (
            <ActionButtons row={row.original} rowActions={table.options.meta!.rowActions} />
          )}
        </div>
      );
    },
  },
];
