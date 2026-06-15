import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface NumberOfTasks {
  numberOfTasks: number;
  onChange: (tasksNum: number) => void;
}

export function NumberOfTasksEntry({ numberOfTasks, onChange }: NumberOfTasks) {
  return (
    <div>
      <Label> Number of Sub-Tasks </Label>
      <p className="text-xs text-foreground/40">Keep as 0 if you want to run infinite tasks</p>
      <Input
        id="number-of-tasks"
        value={numberOfTasks}
        type="number"
        onChange={(e) => onChange(e.target.valueAsNumber)}
        min={0}
      ></Input>
    </div>
  );
}
