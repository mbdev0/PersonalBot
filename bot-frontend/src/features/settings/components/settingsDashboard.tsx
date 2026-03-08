//TODO: what should we put here?
// discord webhook
// should send webhooks on tx send
// send webhooks on tx fail

import { PageHeader } from '@/features/tasks/components/pageHeader';
import { useSettings } from '../hooks/useSettings';
import { Notifcation } from './notification';

export function SettingsDashboard() {
  const { data } = useSettings();

  return (
    <div className="space-y-6">
      <PageHeader>Settings</PageHeader>
      {data && <Notifcation data={data}></Notifcation>}
    </div>
  );
}
