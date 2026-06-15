import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useRef, useState } from 'react';
import type {
  Dashboard,
  DashboardRow,
  StrategyDashboardRow,
  TaskDashboardRow,
} from '../types/dashboard';
import {
  StrategyEventType,
  type StrategySendWSMessage,
  type StrategyWSMessage,
} from '../types/strategies/strategyWebsocket';
import { mapTaskDtoToTask } from '../mapper/taskMapper';
import { STRATEGY_WS } from '@/config/urls';
import { toast } from 'sonner';
import type { SendPositionWSMessage } from '../types/positionWebsocket';

export const useStrategyWebsocket = (
  sendPositionWsMessage: (msg: SendPositionWSMessage) => void
) => {
  const client = useQueryClient();
  const websocket = useRef<WebSocket | undefined>(undefined);
  const subscribedTasks = useRef<Set<number>>(new Set<number>());
  const wasConnected = useRef(false);
  const [isWebsocketOpen, setWebsocketStatus] = useState(false);

  useEffect(() => {
    let cancelled = false;
    const ws = new WebSocket(STRATEGY_WS);
    websocket.current = ws;

    ws.onopen = () => {
      console.log('Strategy WebSocket connected');
      wasConnected.current = true;
      setWebsocketStatus(true);
    };

    ws.onmessage = (event) => {
      try {
        const data: StrategyWSMessage = JSON.parse(event.data);

        if (data.error) {
          console.error('WebSocket error:', data.error);
          toast.error(`websocket error: ${data.error}`);
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

          const updatedStrategy: Partial<StrategyDashboardRow> = {};

          if (
            data.strategy_msg.event === StrategyEventType.STATUS_UPDATE &&
            data.strategy_msg.state
          ) {
            updatedStrategy.state = data.strategy_msg.state;

            if (isTerminalState(data.strategy_msg.state)) {
              subscribedTasks.current.delete(data.strategy_msg.id);
            }
          }

          if (
            data.strategy_msg.event === StrategyEventType.TASK_CREATION &&
            data.strategy_msg.task
          ) {
            const task = mapTaskDtoToTask(data.strategy_msg.task);
            task.sell_fee = strategyTask.data.sell_fee;

            const taskRow: TaskDashboardRow = {
              program: task.program,
              type: 'task',
              data: task,
              id: task.task_id,
              ws_message: task.message,
              state: task.state.task_state,
            };

            updatedStrategy.children = [...strategyTask.children, taskRow];
          }

          if (
            data.strategy_msg.event === StrategyEventType.MESSAGE_UPDATE &&
            data.strategy_msg.message
          ) {
            updatedStrategy.ws_message = data.strategy_msg.message;
          }

          const updatedStrategyTask: DashboardRow = {
            ...strategyTask,
            ...updatedStrategy,
          };

          const updatedDashboardRows: DashboardRow[] = [
            ...oldData.rows.slice(0, strategyTaskIdx),
            updatedStrategyTask,
            ...oldData.rows.slice(strategyTaskIdx + 1),
          ];

          if (
            updatedStrategyTask.data.trading_type === 'BUY' &&
            updatedStrategyTask.state === 'Done'
          ) {
            updatedStrategyTask.tx_message = updatedStrategyTask.ws_message;
            sendPositionWsMessage({ id: updatedStrategyTask.data.buy_task_id, type: 'Subscribe' });
          }

          if (
            updatedStrategyTask.data.trading_type === 'SELL' &&
            updatedStrategyTask.state === 'Done'
          ) {
            updatedStrategyTask.tx_message = updatedStrategyTask.ws_message;
            sendPositionWsMessage({ id: updatedStrategyTask.data.sell_task_id, type: 'Subscribe' });
          }

          return {
            ...oldData,
            rows: updatedDashboardRows,
          };
        });
      } catch (error) {
        console.error('Failed to process WebSocket message:', error);
        toast.error(`failed to process websocket message: ${error}`);
      }
    };

    ws.onclose = (event) => {
      if (cancelled) return;
      console.log('Strategy WebSocket closed:', event.code);
      if (wasConnected.current) {
        toast.info('Disconnected from strategy websocket');
        wasConnected.current = false;
      } else {
        toast.error('Unable to connect to strategy websocket');
      }
      setWebsocketStatus(false);
    };

    ws.onerror = (error) => {
      console.error('Strategy WebSocket error:', error);
    };

    return () => {
      cancelled = true;
      ws.close();
    };
  }, [client, sendPositionWsMessage]);

  const sendStrategyWSMessage = useCallback((msg: StrategySendWSMessage) => {
    if (msg.type === 'Subscribe' && subscribedTasks.current.has(msg.id)) {
      return;
    }

    if (websocket.current?.readyState === WebSocket.OPEN) {
      websocket.current.send(JSON.stringify(msg));

      if (msg.type === 'Subscribe') {
        subscribedTasks.current.add(msg.id);
      } else if (msg.type === 'Unsubscribe') {
        subscribedTasks.current.delete(msg.id);
      }
    } else {
      console.warn('WebSocket not open, cannot send message');
      toast.error('WebSocket not open, cannot send message');
    }
  }, []);

  return {
    sendStrategyWSMessage,
    isStrategyWebsocketOpen: isWebsocketOpen,
  };
};

function isTerminalState(state: string): boolean {
  return ['Done', 'Task Failed', 'Transaction Failed', 'CANCELLED', 'SUCCESS', 'FAILED'].includes(state);
}
