import { useQueryClient } from '@tanstack/react-query';
import { useEffect, useRef } from 'react';
import type {
  Dashboard,
  DashboardRow,
  StrategyDashboardRow,
  TaskDashboardRow,
} from '../types/dashboard';
import type { TaskWSMessage } from '../types/task_websocket';

//this is for states for tasks only
const useStateWebsocket = () => {
  const client = useQueryClient();
  const websocket = useRef<WebSocket | undefined>(undefined);

  useEffect(() => {
    const ws = new WebSocket('ws://127.0.0.1:9090/api/tasks/subscribe');
    websocket.current = ws;

    ws.onopen = () => {
      console.log('connected');
    };

    ws.onmessage = (event) => {
      // on every message we want to update the dashboard - this is for just task states/messages though
      const data: TaskWSMessage = JSON.parse(event.data);
      client.setQueryData(['dashboard'], (oldData: Dashboard | undefined) => {
        if (!oldData) return oldData;

        const strategyTaskIdx = oldData.rows.findIndex((v) => v.id === data.task_event.strategy_id);

        if (strategyTaskIdx === -1) {
          console.error('unable to find strategy task with id: ', data.task_event.strategy_id);
          return oldData;
        }

        const strategyTask = oldData.rows[strategyTaskIdx];
        if (strategyTask.type != 'strategy') {
          console.error('when looking for strategy task - we ended up getting a task for the row');
          return oldData;
        }

        if (shouldMoveStateAndMessageUp(strategyTask)) {
          const updated: StrategyDashboardRow = {
            ...strategyTask,
            state: data.task_event.state,
            ws_message: data.task_event.message,
          };

          const updatedRows: DashboardRow[] = [
            ...oldData.rows.slice(0, strategyTaskIdx),
            updated,
            ...oldData.rows.slice(strategyTaskIdx + 1),
          ];

          return {
            ...oldData,
            rows: updatedRows,
          };
        }

        const childTaskIdx = strategyTask.children.findIndex(
          (v) => v.id === data.task_event.task_id
        );

        if (childTaskIdx === -1) {
          console.error('unable to find task_id with id: ', data.task_event.task_id);
          return oldData;
        }

        const childTask = strategyTask.children[childTaskIdx];
        const updatedChildRow: TaskDashboardRow = {
          ...childTask,
          state: data.task_event.state,
          ws_message: data.task_event.message,
        };

        //we need to now update the child rows array
        const updatedChildRows: TaskDashboardRow[] = [
          ...strategyTask.children.slice(0, childTaskIdx),
          updatedChildRow,
          ...strategyTask.children.slice(childTaskIdx + 1),
        ];

        const updatedDashboardRow: StrategyDashboardRow = {
          ...strategyTask,
          children: updatedChildRows,
        };
        // update the strategy rows  array
        const updatedDashboardRows: DashboardRow[] = [
          ...oldData.rows.slice(0, strategyTaskIdx),
          updatedDashboardRow,
          ...oldData.rows.slice(strategyTaskIdx + 1),
        ];

        return {
          ...oldData,
          rows: updatedDashboardRows,
        };
      });
    };

    ws.onclose = () => {
      console.log('disconnected');
    };

    return () => {
      ws.close();
    };
  }, [client]);

  //we want to return a send object so we can subscribe to tasks//unsubscribe to tasks

  const send = (type: 'Subscribe' | 'Unsubscribe', task_id: number) => {
    if (websocket.current?.OPEN) {
      websocket.current.send(JSON.stringify({ type: type, task_id: task_id }));
    } else {
      console.error('unable to send message to websocket');
    }
  };

  return {
    send,
  };
};

function shouldMoveStateAndMessageUp(row: DashboardRow): boolean {
  return (
    row.type === 'strategy' && (row.data.trading_type === 'BUY' || row.data.trading_type === 'SELL')
  );
}
