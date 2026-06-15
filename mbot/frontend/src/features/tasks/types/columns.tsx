import { type ColumnDef, type RowData } from '@tanstack/react-table';
import { ActionButtons } from '../components/actionButtons';
import { TaskRowType, type DisplayRow } from './tableRows';
import type { RowActions } from './rowActions';
import { formatAddress } from '@/utils/crypto/address_shortner';
import { ArrowRight, ArrowDown } from 'lucide-react';
import { ChildRowActionButtons } from '../components/childRowActionButtons';
import { MessageCell } from '../components/messageCell';
import type { Settings } from '@/features/settings/types/settings';

// this extends the table meta to add a field called row actions so TS doesn't complain
// eslint complaining compaining about no unused vars, but for the interface to work we must copy the same
// interface signature
declare module '@tanstack/react-table' {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface TableMeta<TData extends RowData = RowData> {
    rowActions: RowActions;
    settingsData?: Settings;
  }
}

export const columns: ColumnDef<DisplayRow>[] = [
  {
    accessorKey: 'task_type',
    header: 'Task Type',
    size: 140,
    cell: ({ row }) => {
      const canExpand = row.getCanExpand();
      const isExpanded = row.getIsExpanded();

      return (
        <div className="flex items-center gap-2">
          {canExpand ? (
            <button
              onClick={row.getToggleExpandedHandler()}
              className="h-6 w-6 shrink-0 flex items-center justify-center rounded-md hover:bg-foreground/8 transition-colors ring-1 ring-foreground/5 hover:ring-foreground/15"
            >
              {isExpanded ? (
                <ArrowDown className="h-3.5 w-3.5 text-foreground/60" />
              ) : (
                <ArrowRight className="h-3.5 w-3.5 text-foreground/60" />
              )}
            </button>
          ) : (
            <div className="w-6 shrink-0" />
          )}
          <span>
            {row.original.type === TaskRowType.Task
              ? row.original.data.type
              : row.original.data.trading_type}
          </span>
        </div>
      );
    },
  },
  {
    accessorKey: 'task_name',

    header: 'Task Name',
    size: 180,
    cell: ({ row }) => {
      if (row.original.type === TaskRowType.Task) {
        if (row.original.strategyId) {
          return `${row.original.data.type} ${formatAddress(row.original.data.token_address)}`;
        }
        if (row.original.data.sell_amount && row.original.data.type === 'Sell') {
          return `${row.original.data.type} ${row.original.data.sell_amount * 100}%`;
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
    size: 480,
    cell: ({ row }) => <MessageCell message={row.original.ws_message ?? ''} />,
  },
  {
    accessorKey: 'status',
    header: 'Status',
    size: 250,
    cell: ({ row }) => {
      return row.original.state;
    },
  },
  {
    accessorKey: 'actions',
    header: 'Actions',
    size: 160,
    cell: ({ row, table }) => {
      return (
        <div className="flex justify-center">
          {row.original.type === TaskRowType.Strategy && (
            <ActionButtons
              row={row.original}
              rowActions={table.options.meta!.rowActions}
              settings={table.options.meta!.settingsData}
            />
          )}
          {row.original.type === TaskRowType.Task && (
            <ChildRowActionButtons
              row={row.original}
              rowActions={table.options.meta!.rowActions}
              settings={table.options.meta!.settingsData}
            />
          )}
        </div>
      );
    },
  },
];
