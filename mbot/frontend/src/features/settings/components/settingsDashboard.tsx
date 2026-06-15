import { toast } from 'sonner';
import { useSettings } from '../hooks/useSettings';
import { Notifcation } from './notification';
import { PositionNode } from './positionNodes';
import { QuickSellButtons } from './quickSellButtons';

export function SettingsDashboard() {
  const { data, isError, error } = useSettings();

  if (isError) {
    toast.error(`error whilst loading settings dashboard: ${error}`);
    return;
  }

  return (
    <div className="space-y-6 py-4">
      {data && (
        <div className="flex gap-4">
          <Notifcation data={data}></Notifcation>
          <PositionNode data={data}></PositionNode>
          <QuickSellButtons data={data} />
        </div>
      )}
    </div>
  );
}
