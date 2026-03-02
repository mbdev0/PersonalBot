import { useQuery } from '@tanstack/react-query';
import { getDashboard } from '../api/taskDashboard';
import { useEffect } from 'react';
import type { StrategyDashboardRow } from '../types/dashboard';
import { useStrategyWebsocketSend, useTaskWebsocketSend } from '@/hooks/useWebsocketSend';
import type { SendTaskWSMessage } from '../types/taskWebsocket';
import { isTerminal } from '../types/task';

export function useTaskDashboard() {
  const send = useTaskWebsocketSend();
  const strategySend = useStrategyWebsocketSend();

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
        if (row.state == 'SUCCESS' || row.state == 'Done') {
          return;
        }
        strategySend({ type: 'Subscribe', id: row.data.id });
        if (row.children.length > 0) {
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
