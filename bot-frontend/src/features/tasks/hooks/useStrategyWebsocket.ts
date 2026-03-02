import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useRef } from 'react';
import type { Dashboard, DashboardRow, StrategyDashboardRow } from '../types/dashboard';
import {
  StrategyEventType,
  type StrategySendWSMessage,
  type StrategyWSMessage,
} from '../types/strategies/strategyWebsocket';

export const useStrategyWebsocket = () => {
  const client = useQueryClient();
  const websocket = useRef<WebSocket | undefined>(undefined);
  const subscribedTasks = useRef<Set<number>>(new Set<number>());

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:9090/api/trading/subscribe');
    websocket.current = ws;

    ws.onopen = () => {
      console.log('Strategy WebSocket connected');
    };

    ws.onmessage = (event) => {
      try {
        const data: StrategyWSMessage = JSON.parse(event.data);

        if (data.error) {
          console.error('WebSocket error:', data.error);
          return;
        }

        client.setQueryData(['dashboard'], (oldData: Dashboard | undefined) => {
          if (!oldData) return oldData;

          if (!data.strategy_msg) {
            console.error('No strategy message in WebSocket response');
            return oldData;
          }

          const strategyTaskIdx = oldData.rows.findIndex((v) => v.id === data.strategy_msg?.id);

          if (strategyTaskIdx === -1) {
            console.error('Strategy not found:', data.strategy_msg.id);
            return oldData;
          }

          const strategyTask = oldData.rows[strategyTaskIdx];

          if (strategyTask.type !== 'strategy') {
            console.error('Expected strategy row, got task row');
            return oldData;
          }

          const updates: Partial<StrategyDashboardRow> = {};

          if (
            data.strategy_msg.event === StrategyEventType.STATUS_UPDATE &&
            data.strategy_msg.state
          ) {
            updates.state = data.strategy_msg.state;

            if (isTerminalState(data.strategy_msg.state)) {
              subscribedTasks.current.delete(data.strategy_msg.id);
            }
          }

          if (
            data.strategy_msg.event === StrategyEventType.MESSAGE_UPDATE &&
            data.strategy_msg.message
          ) {
            updates.ws_message = data.strategy_msg.message;
          }

          const updatedStrategyTask: DashboardRow = {
            ...strategyTask,
            ...updates,
          };

          const updatedDashboardRows: DashboardRow[] = [
            ...oldData.rows.slice(0, strategyTaskIdx),
            updatedStrategyTask,
            ...oldData.rows.slice(strategyTaskIdx + 1),
          ];

          return {
            ...oldData,
            rows: updatedDashboardRows,
          };
        });
      } catch (error) {
        console.error('Failed to process WebSocket message:', error);
      }
    };

    ws.onclose = (event) => {
      console.log('Strategy WebSocket closed:', event.code);
    };

    ws.onerror = (error) => {
      console.error('Strategy WebSocket error:', error);
    };

    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close(1000, 'Component unmounting');
      }
    };
  }, [client]);

  const sendStrategyWSMessage = useCallback((msg: StrategySendWSMessage) => {
    if (msg.type === 'Subscribe' && subscribedTasks.current.has(msg.id)) {
      return;
    }
    console.log(subscribedTasks.current);
    console.log('sending subscription for task: ', msg.id);

    if (websocket.current?.readyState === WebSocket.OPEN) {
      websocket.current.send(JSON.stringify(msg));

      if (msg.type === 'Subscribe') {
        subscribedTasks.current.add(msg.id);
      } else if (msg.type === 'Unsubscribe') {
        subscribedTasks.current.delete(msg.id);
      }
    } else {
      console.warn('WebSocket not open, cannot send message');
    }
  }, []);

  return {
    sendStrategyWSMessage,
  };
};

function isTerminalState(state: string): boolean {
  return ['Done', 'Task Failed', 'Transaction Failed'].includes(state);
}
