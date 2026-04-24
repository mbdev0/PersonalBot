import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useRef, useState } from 'react';
import type {
  Dashboard,
  DashboardRow,
  StrategyDashboardRow,
  TaskDashboardRow,
} from '../types/dashboard';
import type { SendTaskWSMessage, TaskWSMessage } from '../types/taskWebsocket';
import { TASK_WS } from '@/config/urls';
import type { SendPositionWSMessage } from '../types/positionWebsocket';

// ONLY FOR CHILDREN
export const useTaskWebsocket = (sendPositionWsMessage: (msg: SendPositionWSMessage) => void) => {
  const client = useQueryClient();
  const websocket = useRef<WebSocket | undefined>(undefined);
  const subscribedTasks = useRef<Set<number>>(new Set<number>());
  const [websocketOpen, setWebsocketStatus] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(TASK_WS);
    websocket.current = ws;

    ws.onopen = () => {
      console.log('connected');
      setWebsocketStatus(true);
    };

    ws.onmessage = (event) => {
      // on every message we want to update the dashboard - this is for just task states/messages though
      const data: TaskWSMessage = JSON.parse(event.data);
      client.setQueryData(['dashboard'], (oldData: Dashboard | undefined) => {
        if (!oldData) return oldData;

        if (data.error) {
          return oldData;
        }

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
        };

        switch (data.task_event.event_type) {
          case 'StatusUpdate':
            updatedChildRow.state = data.task_event.state.task_state;
            break;
          case 'ProgressMessage':
            updatedChildRow.ws_message = data.task_event.message;
            break;
        }

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

        if (
          updatedChildRow.state === 'Done' &&
          updatedChildRow.ws_message.includes('successfully confirmed transaction')
        ) {
          updatedChildRow.tx_message = updatedChildRow.ws_message;
          sendPositionWsMessage({ id: updatedChildRow.id, type: 'Subscribe' });
        }

        return {
          ...oldData,
          rows: updatedDashboardRows,
        };
      });
    };

    ws.onclose = () => {
      console.log('disconnected');
      setWebsocketStatus(false);
    };

    return () => {
      ws.close();
    };
  }, [client]);

  //we want to return a send object so we can subscribe to tasks//unsubscribe to tasks

  const send = useCallback((msg: SendTaskWSMessage) => {
    if (msg.type === 'Subscribe' && subscribedTasks.current.has(msg.id)) {
      return;
    }

    if (websocket.current?.readyState == WebSocket.OPEN) {
      websocket.current.send(JSON.stringify(msg));
      if (msg.type === 'Subscribe') {
        subscribedTasks.current.add(msg.id);
      } else {
        subscribedTasks.current.delete(msg.id);
      }
    } else {
      console.error('unable to send message to websocket');
    }
  }, []);

  return {
    send,
    taskWebsocketOpen: websocketOpen,
  };
};
