import { API_BASE } from '../../../config/urls';
import { MapDashboardDtoToDashboard } from '../mapper/dashboardMapper';
import type { Dashboard, DashboardDto } from '../types/dashboard';

export async function getDashboard(): Promise<Dashboard> {
  const url = API_BASE + '/dashboard';
  const response = await fetch(url);

  if (!response.ok) {
    throw new Error('Error whilst getting tasks');
  }

  const dashboardDto: DashboardDto = await response.json();

  const dashboard = MapDashboardDtoToDashboard(dashboardDto);

  return dashboard;
}
