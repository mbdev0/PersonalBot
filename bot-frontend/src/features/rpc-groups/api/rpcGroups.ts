import { API_BASE } from '@/config/urls';
import type { RPCGroup, RPCGroupDashboardRow, RPCGroupPost, RPCGroupPut } from '../types/rpcGroup';

const baseRPCUrl = API_BASE + '/rpc_groups/';

export async function getRPCDashboard(): Promise<RPCGroupDashboardRow[]> {
  const resp = await fetch(baseRPCUrl, {
    method: 'GET',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  const data: RPCGroupDashboardRow[] = await resp.json();
  return data;
}

export async function postRPCGroup(rpcGroup: RPCGroupPost) {
  const body = JSON.stringify(rpcGroup);

  const resp = await fetch(baseRPCUrl, {
    method: 'POST',
    body: body,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  // const data = await resp.json();
  return resp.status;
}

export async function updateRPCGroup(rpcGroup: RPCGroupPut) {
  const body = JSON.stringify(rpcGroup);

  const resp = await fetch(baseRPCUrl + rpcGroup.id, {
    method: 'PUT',
    body: body,
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  return resp.status;
}

export async function getRPCGroup(id: number) {
  const resp = await fetch(baseRPCUrl + id, {
    method: 'GET',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  if (!resp.ok) {
    return;
  }

  const data: RPCGroup = await resp.json();
  return data;
}

export async function deleteRPCGroup(id: number) {
  const resp = await fetch(baseRPCUrl + id, {
    method: 'DELETE',
  });

  if (!resp.ok) {
    return;
  }
  return resp.status;
}
