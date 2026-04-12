import { useSettings } from '../hooks/useSettings';
import { Notifcation } from './notification';
import { PositionNode } from './positionNodes';

export function SettingsDashboard() {
  const { data } = useSettings();

  return (
    <div className="space-y-6 py-4">
      {data && (
        <div className="flex gap-4">
          <Notifcation data={data}></Notifcation>
          <PositionNode data={data}></PositionNode>
        </div>
      )}
    </div>
  );
}
