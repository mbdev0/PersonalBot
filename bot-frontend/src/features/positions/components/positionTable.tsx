import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import type { PositionDashboard } from '../types/position';
import { formatAddress } from '@/utils/crypto/address_shortner';

interface PositionTableProps {
  data: PositionDashboard;
}

export function PositionTable({ data }: PositionTableProps) {
  return (
    <div className="table-container">
      <Table className="table-fixed w-full">
        <TableHeader>
          <TableRow className="hover:bg-transparent">
            <TableHead className="table-header-cell text-center">Coin</TableHead>
            <TableHead className="table-header-cell text-center">Average MCap Entry ($) </TableHead>
            <TableHead className="table-header-cell text-center">Average MCap Exit ($)</TableHead>
            <TableHead className="table-header-cell text-center">PNL (SOL)</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {data?.map((position) => (
            <TableRow key={position.coin} className="group">
              <TableCell className="table-body-cell text-center">
                <div className=" flex items-center justify-center gap-2">
                  {formatAddress(position.coin, 6, 6)}
                  <span
                    className="text-xs cursor-pointer text-blue-400 hover:underline"
                    onClick={() =>
                      window.open(
                        `https://axiom.trade/meme/${position.address_for_url}?chain=sol`,
                        '_blank',
                        'noopener,noreferrer'
                      )
                    }
                  >
                    axiom ↗
                  </span>
                </div>
              </TableCell>
              <TableCell className="table-body-cell text-center">
                {`$${position.average_market_cap_entry}`}
              </TableCell>
              <TableCell className="table-body-cell text-center">
                {`$${position.average_market_cap_exit}`}
              </TableCell>
              <TableCell className="table-body-cell text-center">{`${position.total_pnl} ◎`}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );
}
