import {
  flexRender,
  getCoreRowModel,
  getExpandedRowModel,
  useReactTable,
} from '@tanstack/react-table';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { useRowActions } from '@/features/tasks/hooks/useRowActions';
import { columns } from '@/features/tasks/types/columns';
import type { DisplayRow } from '@/features/tasks/types/tableRows';
import { useSettings } from '@/features/settings/hooks/useSettings';

export function DashboardTable({
  data,
  setEditingRow,
}: {
  data: DisplayRow[];
  setEditingRow: (row: DisplayRow | null) => void;
}) {
  const rowActions = useRowActions(setEditingRow);
  const settings = useSettings();
  const settingsData = settings.data;

  const table = useReactTable({
    data,
    columns,
    meta: { rowActions, settingsData },
    getCoreRowModel: getCoreRowModel(),
    getExpandedRowModel: getExpandedRowModel(),
    getSubRows: (row) => row.subRows,
  });

  return (
    <div className="overflow-hidden">
      <Table className="table-fixed">
        <TableHeader>
          {table.getHeaderGroups().map((headerGroup) => (
            <TableRow key={headerGroup.id} className="hover:bg-transparent">
              {headerGroup.headers.map((header) => (
                <TableHead
                  key={header.id}
                  style={{ width: header.column.columnDef.size }}
                  className="table-header-cell text-center"
                >
                  {header.isPlaceholder
                    ? null
                    : flexRender(header.column.columnDef.header, header.getContext())}
                </TableHead>
              ))}
            </TableRow>
          ))}
        </TableHeader>
        <TableBody>
          {table.getRowModel().rows.length ? (
            table.getRowModel().rows.map((row) => {
              const depth = row.depth;

              return (
                <TableRow
                  key={row.id}
                  className={`text-center border-0 transition-colors duration-200 ${
                    depth === 1 ? 'bg-foreground/1.5' : ''
                  } ${depth === 2 ? 'bg-foreground/[0.008]' : ''}`}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell
                      key={cell.id}
                      className={[
                        'py-4 text-[13px] font-medium text-center',
                        depth === 0 ? 'text-foreground/80 px-6' : '',
                        depth === 1 ? 'text-foreground/60 text-[12px] px-8' : '',
                        depth === 2 ? 'text-foreground/40 text-[12px] px-12' : '',
                      ]
                        .filter(Boolean)
                        .join(' ')}
                    >
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </TableCell>
                  ))}
                </TableRow>
              );
            })
          ) : (
            <TableRow className="hover:bg-transparent">
              <TableCell
                colSpan={columns.length}
                className="h-32 text-center text-[13px] text-muted-foreground/40"
              >
                No tasks yet
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </div>
  );
}
