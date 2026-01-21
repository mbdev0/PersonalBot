//This will be for Strategy tasks, e.g. AFK, Spam etc.

import { API_BASE } from '../../../config/urls';
import { MapStrategyToPostDto, MapStrategyToPutDto } from '../mapper/strategyMapper';
import type { StrategyTaskPost, StrategyTaskPut } from '../types/strategyTask';

const STRATEGY_BASE = '/trading';

export async function postStrategy(task: StrategyTaskPost) {
  const mappedTasks = MapStrategyToPostDto(task);
  const url = API_BASE + STRATEGY_BASE + '/create';

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(mappedTasks),
  });

  if (!response.ok) {
    throw new Error(`Failed to add new task: ${response.status}`);
  }

  return response.status;
}

export async function deleteStrategy(id: number) {
  const url = API_BASE + STRATEGY_BASE + `/task/${id}`;

  const response = await fetch(url, {
    method: 'DELETE',
  });

  if (!response.ok) {
    throw new Error(`deleting task ${id} failed: ${response.status}`);
  }

  return response.status;
}

export async function putStrategy(task: StrategyTaskPut) {
  const mappedTasks = MapStrategyToPutDto(task);
  const url = API_BASE + STRATEGY_BASE + `/task/${task.id}`;

  const response = await fetch(url, {
    method: 'PUT',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(mappedTasks),
  });

  if (!response.ok) {
    throw new Error(`Updating task ${task.id} failed: ${response.status}`);
  }

  return response.status;
}

export async function startStrategy(id: number) {
  const url = API_BASE + STRATEGY_BASE + `/task/start/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error(`starting task ${id} failed: ${response.status}`);
  }

  return response.status;
}

export async function stopStrategy(id: number) {
  const url = API_BASE + STRATEGY_BASE + `/task/stop/${id}`;

  const response = await fetch(url, {
    method: 'GET',
  });

  if (!response.ok) {
    throw new Error(`starting task ${id} failed: ${response.status}`);
  }

  return response.status;
}
