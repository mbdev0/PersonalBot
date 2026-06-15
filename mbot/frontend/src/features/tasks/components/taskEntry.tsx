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
import { PROGRAMS } from '../types/programs';
import { SpamTaskEntry } from './forms/spamTaskEntry';

interface TaskEntryProps {
  onClose: () => void;
}

export function TaskEntry({ onClose }: TaskEntryProps) {
  const [taskType, setTaskType] = useState('Buy');
  const [program, setProgram] = useState('PumpfunNative');

  return (
    <div className="space-y-6">
      <div className="flex gap-4">
        <div className="flex flex-col items-center space-y-2 ">
          <Label htmlFor="taskType">Task Type</Label>
          <Select value={taskType} onValueChange={setTaskType}>
            <SelectTrigger id="taskType">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="Buy">Buy</SelectItem>
              <SelectItem value="Sell">Sell</SelectItem>
              <SelectItem value="AFK">AFK</SelectItem>
              <SelectItem value="Spam">Spam</SelectItem>
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

      {taskType === 'Buy' && <BuyTaskEntry onClose={onClose} program={program} />}
      {taskType === 'Sell' && <SellTaskEntry onClose={onClose} program={program} />}
      {taskType === 'AFK' && <AFKTaskEntry onClose={onClose} program={program} />}
      {taskType === 'Spam' && <SpamTaskEntry onClose={onClose} program={program} />}
    </div>
  );
}
