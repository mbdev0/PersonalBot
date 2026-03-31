import { useState } from 'react';
import { BuyTaskEntry } from './forms/buyTaskEntry';
import { SellTaskEntry } from './forms/sellTaskEntry';
import { AFKTaskEntry } from './forms/afkTaskEntry';
import { Label } from '@/components/ui/label';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface TaskEntryProps {
  onClose: () => void;
}

export function TaskEntry({ onClose }: TaskEntryProps) {
  const [taskType, setTaskType] = useState('Buy');

  return (
    <div className="space-y-6">
      <div className="space-y-2">
        <Label htmlFor="taskType">Task Type</Label>
        <Select value={taskType} onValueChange={setTaskType}>
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

      {taskType === 'Buy' && <BuyTaskEntry onClose={onClose} />}
      {taskType === 'Sell' && <SellTaskEntry onClose={onClose} />}
      {taskType === 'AFK' && <AFKTaskEntry onClose={onClose} />}
    </div>
  );
}
