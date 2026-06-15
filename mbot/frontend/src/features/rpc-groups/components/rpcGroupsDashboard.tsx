import { BotDialog } from '@/components/botDialog';
import { useState } from 'react';
import { RpcGroupTable } from './rpcGroupTable';
import { RpcGroupEntry } from './rpcGroupEntry';
import { Button } from '@/components/ui/button';

export function RPCGroupDashboard() {
  const [isAddModalShowing, setAddModal] = useState(false);

  return (
    <div className="space-y-8 py-4">
      <div className="flex justify-end">
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
