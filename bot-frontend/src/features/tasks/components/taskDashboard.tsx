import { useState } from 'react';
import { TaskTable } from './taskTable';
import { TaskEntry } from './taskEntry';
import { Button } from '@/components/ui/button';
import { BotDialog } from '../../../components/botDialog';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';

export function TaskDashboard() {
  const [isModalShowing, setModalShowing] = useState(false);
  return (
    <div className="task_dashboard">
      <div className="flex justify-end px-4 p-1">
        <Button
          className=" p-2 rounded-xl border border-accent-foreground "
          onClick={() => setModalShowing(true)}
        >
          Add Task
        </Button>
      </div>

      <TaskTable />
      <BotDialog isOpen={isModalShowing} onClose={() => setModalShowing(false)}>
        <DialogHeader>
          <DialogTitle className="font-bold text-foreground">Add Task</DialogTitle>
        </DialogHeader>
        <TaskEntry onClose={() => setModalShowing(false)}></TaskEntry>
      </BotDialog>
    </div>
  );
}
