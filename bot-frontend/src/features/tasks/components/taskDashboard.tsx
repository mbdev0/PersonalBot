import { useState } from 'react';
import { TaskTable } from './taskTable';
import { TaskEntry } from './taskEntry';
import { Button } from '@/components/ui/button';
import { BotDialog } from '../../../components/botDialog';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';

export function TaskDashboard() {
  const [isBotDialogShowing, setBotDialogShowing] = useState(false);
  return (
    <div className="task_dashboard">
      <div className="flex justify-end px-4 p-1">
        <Button
          className=" p-2 rounded-xl border border-accent-foreground "
          onClick={() => setBotDialogShowing(true)}
        >
          Add Task
        </Button>
      </div>

      <TaskTable />
      <BotDialog isOpen={isBotDialogShowing} onClose={() => setBotDialogShowing(false)}>
        <DialogHeader>
          <DialogTitle className="font-bold text-foreground">Add Task</DialogTitle>
        </DialogHeader>
        <TaskEntry onClose={() => setBotDialogShowing(false)}></TaskEntry>
      </BotDialog>
    </div>
  );
}
