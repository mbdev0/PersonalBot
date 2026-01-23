import type { Task } from '../types/task';
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

interface TaskUpdateProps {
  task: Task;
  onClose: () => void;
}

export function TaskUpdate({ task, onClose }: TaskUpdateProps) {
  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <Label htmlFor="taskType">Task Type</Label>
        <Select disabled={true} value={task.type}>
          <SelectTrigger id="taskType">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="Buy">Buy</SelectItem>
            <SelectItem value="Sell">Sell</SelectItem>
          </SelectContent>
        </Select>
      </div>

      {task.type === 'Buy' && <BuyTaskEdit task={task} onClose={onClose}></BuyTaskEdit>}
      {task.type === 'Sell' && <SellTaskEdit task={task} onClose={onClose}></SellTaskEdit>}
    </div>
  );
}
