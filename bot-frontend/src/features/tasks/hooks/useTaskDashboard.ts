import { useQuery } from '@tanstack/react-query';
import { getDashboard } from '../api/taskDashboard';
import { useEffect } from 'react';
import type { StrategyDashboardRow } from '../types/dashboard';
import { useWebsocketSend } from '@/hooks/useWebsocketSend';
import type { SendTaskWSMessage } from '../types/task_websocket';
import { isTerminal } from '../types/task';

export function useTaskDashboard() {
  const send = useWebsocketSend();

  const query = useQuery({
    queryKey: ['dashboard'],
    queryFn: getDashboard,
  });

  useEffect(() => {
    if (!query.data) {
      return;
    }

    query.data.rows.forEach((row) => {
      if (row.type === 'strategy') {
        switch (row.data.trading_type) {
          case 'BUY':
            if (!isTerminal(row.state)) send({ type: 'Subscribe', id: row.data.buy_task_id });
            break;
          case 'SELL':
            if (!isTerminal(row.state)) send({ type: 'Subscribe', id: row.data.sell_task_id });
            break;
          default:
            subscribeToTaskChildren(row, send);
        }
      }
    });
  }, [query.data]);

  return query;
}

function subscribeToTaskChildren(
  row: StrategyDashboardRow,
  send: (msg: SendTaskWSMessage) => void
): void {
  if (row.type !== 'strategy') {
    return;
  }

  if (row.data.trading_type === 'BUY' || row.data.trading_type === 'SELL') {
    return;
  }

  row.children.forEach((taskRow) => {
    if (!isTerminal(taskRow.state)) send({ type: 'Subscribe', id: taskRow.id });
  });
}
