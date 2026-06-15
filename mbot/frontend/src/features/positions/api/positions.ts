import { API_BASE } from '@/config/urls';
import type { PositionDashboard } from '../types/position';

export async function getPositionDashboard() {
  const resp = await fetch(`${API_BASE}/position/dashboard`, {
    method: 'GET',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  if (!resp.ok) {
    const text = await resp.text();
    throw new Error(text || `Request failed (${resp.status})`);
  }

  const dashboard: PositionDashboard = await resp.json();
  return dashboard;
}
