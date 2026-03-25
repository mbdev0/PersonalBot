import { BotDialog } from '@/components/botDialog';
import { PageHeader } from '@/features/tasks/components/pageHeader';
import { useState } from 'react';
import { RpcGroupTable } from './rpcGroupTable';
import { RpcGroupEntry } from './rpcGroupEntry';
import { Button } from '@/components/ui/button';

export function RPCGroupDashboard() {
  const [isAddModalShowing, setAddModal] = useState(false);

  return (
    <div className="space-y-8">
      {/* in here we'll have our add button + our add modal */}
      <div className="flex justify-between">
        <PageHeader>RPC Groups</PageHeader>
        <Button
          className="h-9 px-4 text-[13px] text-accent-foreground font-medium bg-foreground/5 hover:bg-foreground/10 border-0 ring-1 ring-foreground/20 hover:ring-foreground/30 transition-all duration-200"
          onClick={() => setAddModal(true)}
        >
          Add RPC Group
        </Button>
      </div>

      <BotDialog
        isOpen={isAddModalShowing}
        onClose={() => {
          setAddModal(false);
        }}
      >
        <RpcGroupEntry onCompletion={() => setAddModal(false)}></RpcGroupEntry>
      </BotDialog>

      <RpcGroupTable />
    </div>
  );
}
