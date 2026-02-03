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
        <button onClick={row.getToggleExpandedHandler()} style={{ cursor: 'pointer' }}>
          {row.getIsExpanded() ? <ArrowDown /> : <ArrowRight />}
        </button>
      ) : (
        ''
      );
    },
  },
  {
    accessorKey: 'task_type',
    header: () => <div className="font-extrabold">Task Type</div>,
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

    header: () => <div className="font-extrabold">Task Name</div>,
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
    header: () => <div className="font-extrabold">Message</div>,
    cell: ({ row }) => {
      return row.original.wsMessage;
    },
  },
  {
    accessorKey: 'status',
    header: () => <div className="font-extrabold">Status</div>,
    cell: ({ row }) => {
      return row.original.state;
    },
  },
  {
    accessorKey: 'actions',
    header: () => <div className="font-extrabold">Actions</div>,
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
