import { useQuery } from '@tanstack/react-query';
import { getDashboard } from '../api/taskDashboard';

export function useTaskDashboard() {
  return useQuery({
    queryKey: ['tasks'],
    queryFn: getDashboard,
  });
}
