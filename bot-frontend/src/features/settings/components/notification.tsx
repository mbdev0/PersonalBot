import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import type { Settings } from '../types/settings';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { usePostSettings, useTestDiscordWebhook } from '../hooks/useSettings';
import { useState } from 'react';

interface NotifcationProps {
  data: Settings;
}

export function Notifcation({ data }: NotifcationProps) {
  const mutation = usePostSettings();
  const testMutation = useTestDiscordWebhook();
  const [discordWebhook, setDiscordWebhook] = useState(data.discord_webhook);

  const update = (patch: Partial<Settings>) => {
    if (!data) return;
    mutation.mutate({ ...data, ...patch });
  };

  return (
    <div className="rounded-lg border border-foreground/10 bg-foreground/3 p-6 space-y-6 w-140">
      <p className="text-sm font-semibold text-foreground/50 uppercase">Notifications</p>

      <div className="space-y-2">
        <Label>Discord Webhook</Label>
        <div className="flex gap-2">
          <Input
            key={data?.discord_webhook}
            defaultValue={data?.discord_webhook ?? ''}
            onChange={(e) => setDiscordWebhook(e.target.value)}
          />

          <Button
            className="hover:bg-green-700 hover:opacity-70"
            variant="default"
            onClick={() => update({ discord_webhook: discordWebhook })}
          >
            Save
          </Button>
          <Button variant="outline" onClick={() => testMutation.mutate(discordWebhook)}>
            Test
          </Button>
        </div>
      </div>

      <div className="space-y-3">
        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium">Failed Transactions</p>
            <p className="text-xs text-foreground/40">Notify when a transaction fails</p>
          </div>
          <Switch
            disabled={data.discord_webhook == ''}
            checked={data.discord_webhook == '' ? false : data.send_on_fail}
            onCheckedChange={(checked) => update({ send_on_fail: checked })}
          />
        </div>

        <div className="flex items-center justify-between">
          <div>
            <p className="text-sm font-medium">Successful Transactions</p>
            <p className="text-xs text-foreground/40">Notify when a transaction confirms</p>
          </div>
          <Switch
            disabled={data.discord_webhook == ''}
            checked={data.discord_webhook == '' ? false : data.send_on_success}
            onCheckedChange={(checked) => update({ send_on_success: checked })}
          />
        </div>
      </div>
    </div>
  );
}
