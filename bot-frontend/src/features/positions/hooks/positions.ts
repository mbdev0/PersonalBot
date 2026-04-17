import { useQuery } from '@tanstack/react-query';
import { getPositionDashboard } from '../api/positions';

export function usePositionDashboard() {
  return useQuery({
    queryKey: ['position'],
    queryFn: getPositionDashboard,
    refetchInterval: 10000,
  });
}
