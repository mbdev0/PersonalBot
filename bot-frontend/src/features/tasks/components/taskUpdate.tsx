import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { BuyTaskEdit } from './forms/buyTaskEdit';
import { SellTaskEdit } from './forms/sellTaskEdit';
import { TaskRowType, type DisplayRow } from '../types/tableRows';
import { AFKTaskEdit } from './forms/afkTaskEdit';
import { useState } from 'react';
import { PROGRAMS } from '../types/programs';

interface TaskUpdateProps {
  row: DisplayRow;
  onClose: () => void;
}

export function TaskUpdate({ row, onClose }: TaskUpdateProps) {
  const taskType = row.type === TaskRowType.Task ? row.data.type : row.data.trading_type;
  const [program, setProgram] = useState(row.program);

  return (
    <div className="space-y-6">
      <div className="flex gap-4">
        <div className="space-y-2">
          <Label htmlFor="taskType">Task Type</Label>
          <Select disabled={true} value={taskType}>
            <SelectTrigger id="taskType">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="BUY">Buy</SelectItem>
              <SelectItem value="SELL">Sell</SelectItem>
              <SelectItem value="AFK">AFK</SelectItem>
            </SelectContent>
          </Select>
        </div>

        <div className="flex flex-col items-center space-y-2 ">
          <Label htmlFor="program">Program</Label>
          <Select value={program} onValueChange={setProgram}>
            <SelectTrigger id="program">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {PROGRAMS.map((p) => {
                return <SelectItem value={p.value}>{p.label}</SelectItem>;
              })}
            </SelectContent>
          </Select>
        </div>
      </div>

      {row.type === TaskRowType.Strategy && row.data.trading_type === 'BUY' && (
        <BuyTaskEdit task={row.data} onClose={onClose} program={program}></BuyTaskEdit>
      )}
      {row.type === TaskRowType.Strategy && row.data.trading_type === 'SELL' && (
        <SellTaskEdit task={row.data} onClose={onClose} program={program}></SellTaskEdit>
      )}

      {row.type === TaskRowType.Strategy && row.data.trading_type === 'AFK' && (
        <AFKTaskEdit task={row.data} onClose={onClose} program={program}></AFKTaskEdit>
      )}
    </div>
  );
}
