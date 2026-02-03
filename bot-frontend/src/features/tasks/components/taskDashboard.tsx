import { useState } from 'react';
import { TaskTable } from './taskTable';
import { TaskEntry } from './taskEntry';
import { Button } from '@/components/ui/button';
import { BotDialog } from '../../../components/botDialog';
import { DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { PageHeader } from './pageHeader';

export function TaskDashboard() {
  const [isBotDialogShowing, setBotDialogShowing] = useState(false);
  return (
    <div className="task_dashboard space-y-8">
      <div className="flex items-center justify-between">
        <div>
          <PageHeader>Tasks</PageHeader>
        </div>
        <Button
          className="h-9 px-4 text-[13px] text-accent-foreground font-medium bg-foreground/5 hover:bg-foreground/10 border-0 ring-1 ring-foreground/20 hover:ring-foreground/30 transition-all duration-200"
          onClick={() => setBotDialogShowing(true)}
        >
          Add Task
        </Button>
      </div>

      <TaskTable />
      <BotDialog isOpen={isBotDialogShowing} onClose={() => setBotDialogShowing(false)}>
        <DialogHeader>
          <DialogTitle className="font-semibold text-foreground">Add Task</DialogTitle>
        </DialogHeader>
        <TaskEntry onClose={() => setBotDialogShowing(false)}></TaskEntry>
      </BotDialog>
    </div>
  );
}
