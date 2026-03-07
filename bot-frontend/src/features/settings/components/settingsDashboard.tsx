//TODO: what should we put here?
// discord webhook
// should send webhooks on tx send
// send webhooks on tx fail

import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { PageHeader } from '@/features/tasks/components/pageHeader';

export function SettingsDashboard() {
  return (
    <div className="space-y-6">
      <PageHeader>Settings</PageHeader>

      <div className="rounded-lg border border-foreground/10 bg-foreground/3 p-6 space-y-6 w-140">
        <p className="text-sm font-semibold text-foreground/50 uppercase">Notifications</p>

        <div className="space-y-2">
          <Label>Discord Webhook</Label>
          <div className="flex gap-2">
            <Input placeholder="https://discord.com/api/webhooks/..." />
            <Button className="hover:bg-green-700 hover:opacity-70" variant="default">
              Save
            </Button>
            <Button variant="outline">Test</Button>
          </div>
        </div>

        <div className="space-y-3">
          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium">Failed Transactions</p>
              <p className="text-xs text-foreground/40">Notify when a transaction fails</p>
            </div>
            <Switch />
          </div>

          <div className="flex items-center justify-between">
            <div>
              <p className="text-sm font-medium">Successful Transactions</p>
              <p className="text-xs text-foreground/40">Notify when a transaction confirms</p>
            </div>
            <Switch />
          </div>
        </div>
      </div>
    </div>
  );
}
