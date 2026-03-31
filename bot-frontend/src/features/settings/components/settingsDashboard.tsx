import { PageHeader } from '@/features/tasks/components/pageHeader';
import { useSettings } from '../hooks/useSettings';
import { Notifcation } from './notification';
import { PositionNode } from './positionNodes';

export function SettingsDashboard() {
  const { data } = useSettings();

  return (
    <div className="space-y-6">
      <PageHeader>Settings</PageHeader>
      {data && (
        <div className="flex space-x-4">
          <Notifcation data={data}></Notifcation>
          <PositionNode data={data}></PositionNode>
        </div>
      )}
    </div>
  );
}
