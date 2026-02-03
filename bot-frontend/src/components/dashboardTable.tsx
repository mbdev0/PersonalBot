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
import { TaskRowType, type DisplayRow } from '@/features/tasks/types/tableRows';

export function DashboardTable({
  data,
  setEditingRow,
}: {
  data: DisplayRow[];
  setEditingRow: (row: DisplayRow | null) => void;
}) {
  const rowActions = useRowActions(setEditingRow);

  const table = useReactTable({
    data,
    columns,
    meta: { rowActions },
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
                  style={{ width: header.column.columnDef.size }}
                  className="table-header-cell text-center"
                  key={header.id}
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
              const isChild = row.original.type === TaskRowType.Task && row.original.strategyId;
              const isParent = row.getCanExpand();
              const isExpanded = isParent && row.getIsExpanded();
              return (
                <TableRow
                  key={row.id}
                  className={`text-center border-0 group transition-colors duration-200 ${
                    isChild ? 'bg-foreground/1.5' : ''
                  } ${isExpanded ? 'bg-foreground/2' : ''}`}
                  data-state={row.getIsSelected() && 'selected'}
                >
                  {row.getVisibleCells().map((cell, cellIndex) => (
                    <TableCell
                      key={cell.id}
                      className={`py-4 px-6 text-center text-[13px] font-medium ${
                        isChild
                          ? cellIndex === 0
                            ? 'pl-14 relative before:content-[""] before:absolute before:left-8 before:top-0 before:bottom-0 before:w-px before:bg-foreground/10'
                            : 'text-foreground/60 text-[12px]'
                          : 'text-foreground/80'
                      }`}
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
