import { API_BASE } from '../../../config/urls';
import type { StrategyTask } from '../types/strategies/strategyTask';
import type { StrategyTaskPost } from '../types/strategies/strategyTaskPost';
import type { StrategyTaskPut } from '../types/strategies/strategyTaskPut';

const STRATEGY_BASE = '/trading';

function prepareBody(task: StrategyTaskPost | StrategyTaskPut) {
  const body: Record<string, unknown> = { ...task };

  if ('sell_strategies' in task && task.sell_strategies) {
    body.sell_strategies = task.sell_strategies.map(({ type, value, sell_amount }) => ({
      type,
      value,
      sell_amount,
    }));
  }

  if ('sell_fee' in task && task.sell_fee == null) {
    body.sell_fee = 0;
  }

  return body;
}

export async function postStrategy(task: StrategyTaskPost) {
  const url = API_BASE + STRATEGY_BASE + '/create';

  const response = await fetch(url, {
    method: 'POST',
    headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(prepareBody(task)),
  });

  if (!response.ok) {
    throw new Error(`Failed to add new task: ${response.status}`);
  }

  const createdStrategyTask: StrategyTask = await response.json();
  return createdStrategyTask;
}

export async function deleteStrategy(id: number) {
  const url = API_BASE + STRATEGY_BASE + `/task/${id}`;

  const response = await fetch(url, { method: 'DELETE' });

  if (!response.ok) {
    throw new Error(`deleting task ${id} failed: ${response.status}`);
  }

  return response.status;
}

export async function putStrategy(task: StrategyTaskPut) {
  const { id, ...rest } = task;
  const url = API_BASE + STRATEGY_BASE + `/task/${id}`;

  const response = await fetch(url, {
    method: 'PUT',
    headers: { Accept: 'application/json', 'Content-Type': 'application/json' },
    body: JSON.stringify(prepareBody(rest as StrategyTaskPost)),
  });

  if (!response.ok) {
    throw new Error(`Updating task ${id} failed: ${response.status}`);
  }

  return response.status;
}

export async function startStrategy(id: number) {
  const url = API_BASE + STRATEGY_BASE + `/task/start/${id}`;
  const response = await fetch(url, { method: 'GET' });

  if (!response.ok) {
    throw new Error(`starting task ${id} failed: ${response.status}`);
  }

  return response.status;
}

export async function stopStrategy(id: number) {
  const url = API_BASE + STRATEGY_BASE + `/task/stop/${id}`;
  const response = await fetch(url, { method: 'GET' });

  if (!response.ok) {
    throw new Error(`stopping task ${id} failed: ${response.status}`);
  }

  return response.status;
}
