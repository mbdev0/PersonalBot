import { Label } from '@/components/ui/label';
import { usePositionDashboard } from '../hooks/positions';
import { PositionTable } from './positionTable';

export function PositionDashboard() {
  const { isError, data, error, isLoading } = usePositionDashboard();

  if (isError) {
    const string = `Error whilst trying to load Position Dashboard - Error: ${error}`;
    return <Label>{string}</Label>;
  }

  if (isLoading) {
    return <Label>Loading Position Dashboard...</Label>;
  }

  return (
    <div className="space-y-8 py-4">
      <div></div>

      {data && <PositionTable data={data}></PositionTable>}
    </div>
  );
}
