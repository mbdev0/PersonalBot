import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useRef, useState } from 'react';
import type { Dashboard, StrategyDashboardRow, TaskDashboardRow } from '../types/dashboard';
import type { SendTaskWSMessage, TaskWSMessage } from '../types/taskWebsocket';
import { TASK_WS } from '@/config/urls';
import type { SendPositionWSMessage } from '../types/positionWebsocket';

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
      const data: TaskWSMessage = JSON.parse(event.data);

      client.setQueryData(['dashboard'], (oldData: Dashboard | undefined) => {
        if (!oldData || data.error) return oldData;

        const strategyTaskIdx = oldData.rows.findIndex((v) => v.id === data.task_event.strategy_id);

        if (strategyTaskIdx === -1) {
          console.error('unable to find strategy task with id: ', data.task_event.strategy_id);
          return oldData;
        }

        const strategyTask = oldData.rows[strategyTaskIdx];
        if (strategyTask.type !== 'strategy') {
          console.error('expected strategy row but got task row');
          return oldData;
        }

        const { updatedChildren, updatedTask } = updateTaskInChildren(
          strategyTask.children,
          data.task_event.task_id,
          data,
          sendPositionWsMessage
        );

        if (!updatedTask) {
          console.error('unable to find task_id: ', data.task_event.task_id);
          return oldData;
        }

        // subscribe to position if buy task just completed
        if (
          updatedTask.state === 'Done' &&
          updatedTask.ws_message.includes('successfully confirmed transaction')
        ) {
          sendPositionWsMessage({ id: updatedTask.id, type: 'Subscribe' });
        }

        const updatedDashboardRow: StrategyDashboardRow = {
          ...strategyTask,
          children: updatedChildren,
        };

        return {
          ...oldData,
          rows: [
            ...oldData.rows.slice(0, strategyTaskIdx),
            updatedDashboardRow,
            ...oldData.rows.slice(strategyTaskIdx + 1),
          ],
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
  }, [client, sendPositionWsMessage]);

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

function applyTaskEvent(task: TaskDashboardRow, event: TaskWSMessage): TaskDashboardRow {
  const updated = { ...task };
  switch (event.task_event.event_type) {
    case 'StatusUpdate':
      updated.state = event.task_event.state.task_state;
      break;
    case 'ProgressMessage':
      updated.ws_message = event.task_event.message;
      break;
  }
  return updated;
}

function updateTaskInChildren(
  children: TaskDashboardRow[],
  taskId: number,
  event: TaskWSMessage,
  sendPositionWsMessage: (msg: SendPositionWSMessage) => void
): { updatedChildren: TaskDashboardRow[]; updatedTask: TaskDashboardRow | null } {
  let updatedTask: TaskDashboardRow | null = null;

  // depth 1 — direct children (buy tasks on AFK, sell tasks on BUY strategy)
  const directIdx = children.findIndex((v) => v.id === taskId);
  if (directIdx !== -1) {
    updatedTask = applyTaskEvent(children[directIdx], event);
    return {
      updatedChildren: [
        ...children.slice(0, directIdx),
        updatedTask,
        ...children.slice(directIdx + 1),
      ],
      updatedTask,
    };
  }

  // depth 2 — grandchildren (sell tasks under AFK buy tasks)
  const updatedChildren = children.map((buyTask) => {
    if (!buyTask.children?.length) return buyTask;

    const sellIdx = buyTask.children.findIndex((v) => v.id === taskId);
    if (sellIdx === -1) return buyTask;

    updatedTask = applyTaskEvent(buyTask.children[sellIdx], event);
    return {
      ...buyTask,
      children: [
        ...buyTask.children.slice(0, sellIdx),
        updatedTask,
        ...buyTask.children.slice(sellIdx + 1),
      ],
    };
  });

  return { updatedChildren, updatedTask };
}
