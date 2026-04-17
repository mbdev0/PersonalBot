import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { PositionDashboard } from '../types/position';

interface PositionTableProps {
  data: PositionDashboard;
}

export function PositionTable({ data }: PositionTableProps) {
  return (
    <div className="table-container">
      <Table>
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="table-header-cell text-center">Coin</TableHead>
            <TableHead className="table-header-cell text-center">Average MCap Entry</TableHead>
            <TableHead className="table-header-cell text-center">Average MCap Exit</TableHead>
            <TableHead className="table-header-cell text-center">PNL</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data?.map((position) => (
            <TableRow key={position.coin} className="group">
              <TableCell className="table-body-cell text-center">{position.coin}</TableCell>
              <TableCell className="table-body-cell text-center">
                {position.average_market_cap_entry}
              </TableCell>
              <TableCell className="table-body-cell text-center">
                {position.average_market_cap_exit}
              </TableCell>
              <TableCell className="table-body-cell text-center">{position.total_pnl}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
