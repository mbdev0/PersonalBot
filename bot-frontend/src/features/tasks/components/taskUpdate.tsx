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
import { TaskRowType, type Row } from '../types/tableRows';
import { AFKTaskEdit } from './forms/afkTaskEdit';

interface TaskUpdateProps {
  row: Row;
  onClose: () => void;
}

export function TaskUpdate({ row, onClose }: TaskUpdateProps) {
  const taskType = row.type === TaskRowType.Task ? row.data.type : row.data.tradingType;

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <Label htmlFor="taskType">Task Type</Label>
        <Select disabled={true} value={taskType}>
          <SelectTrigger id="taskType">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="Buy">Buy</SelectItem>
            <SelectItem value="Sell">Sell</SelectItem>
            <SelectItem value="AFK">AFK</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {row.type === TaskRowType.Task && row.data.type === 'Buy' && (
        <BuyTaskEdit task={row.data} onClose={onClose}></BuyTaskEdit>
      )}
      {row.type === TaskRowType.Task && row.data.type === 'Sell' && (
        <SellTaskEdit task={row.data} onClose={onClose}></SellTaskEdit>
      )}

      {row.type === TaskRowType.Strategy && row.data.tradingType === 'AFK' && (
        <AFKTaskEdit task={row.data} onClose={onClose}></AFKTaskEdit>
      )}
    </div>
  );
}
