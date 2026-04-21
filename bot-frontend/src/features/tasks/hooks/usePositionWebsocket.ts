import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useRef, useState } from 'react';
import type {
  Dashboard,
  DashboardRow,
  StrategyDashboardRow,
  TaskDashboardRow,
} from '../types/dashboard';
import { POSITION_WS } from '@/config/urls';
import type { PositionWebsocketMessage, SendPositionWSMessage } from '../types/positionWebsocket';

// ONLY FOR CHILDREN
export const usePositionWebsocket = () => {
  const client = useQueryClient();
  const websocket = useRef<WebSocket | undefined>(undefined);
  const subscribedTasks = useRef<Set<number>>(new Set<number>());
  const [websocketOpen, setWebsocketStatus] = useState(false);

  useEffect(() => {
    const ws = new WebSocket(POSITION_WS);
    websocket.current = ws;

    ws.onopen = () => {
      console.log('connected');
      setWebsocketStatus(true);
    };

    ws.onmessage = (event) => {
      // on every message we want to update the dashboard - this is for just task states/messages though
      const data: PositionWebsocketMessage = JSON.parse(event.data);
      client.setQueryData(['dashboard'], (oldData: Dashboard | undefined) => {
        if (!oldData) return oldData;

        if (!data) {
          return oldData;
        }

        // we basically want to update the message's for buy/sell strategy tasks or child tasks
        // to do this we need strategy_id in position, otherwise we'd have to search through every single strategy task to find our position_id
        // search for row with strategy id,
        const strategyRowIdx = oldData.rows.findIndex(
          (row) => row.type === 'strategy' && row.id === data.position_msg.strategy_id
        );

        if (strategyRowIdx === -1) {
          return oldData;
        }

        const strategyRow = oldData.rows[strategyRowIdx];

        if (!strategyRow || strategyRow == undefined) {
          return oldData;
        }

        if (strategyRow.type !== 'strategy') {
          return oldData;
        }

        const updatedRow: DashboardRow = ((strategyRow: DashboardRow) => {
          if (
            strategyRow.type === 'strategy' &&
            (strategyRow.data.trading_type === 'BUY' || strategyRow.data.trading_type === 'SELL')
          ) {
            return { ...strategyRow, ws_message: prettifyWsMessage(data) };
          } else {
            if (!strategyRow.children || strategyRow.children.length == 0) {
              return { ...strategyRow };
            }

            const childRow = strategyRow.children.find((child) => {
              return child.id === data.position_msg.buy_task_id;
            });

            if (!childRow) {
              return { ...strategyRow };
            }

            const updatedChildRow: TaskDashboardRow = {
              ...childRow,
              ws_message: prettifyWsMessage(data),
            };
            const updatedStrategyRow: StrategyDashboardRow = {
              ...strategyRow,
              children: strategyRow.children.map((child) =>
                child.id === data.position_msg.buy_task_id ? updatedChildRow : child
              ),
            };

            return updatedStrategyRow;
          }
        })(strategyRow);

        const updatedDashboardRows: DashboardRow[] = [
          ...oldData.rows.slice(0, strategyRowIdx),
          updatedRow,
          ...oldData.rows.slice(strategyRowIdx + 1),
        ];

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

  const sendPositionWsMessage = useCallback((msg: SendPositionWSMessage) => {
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
    sendPositionWsMessage,
    isPositionWsOpen: websocketOpen,
  };
};

function prettifyWsMessage(ws_message: PositionWebsocketMessage) {
  /* 
    [ unrealized PNL: ______ |  market cap: _______ | total pnl: ________ ]
  */

  const unrealizedPnl = toSignificantDecimals(ws_message.position_msg.unrealized_profit);
  const marketCap = parseFloat(ws_message.position_msg.market_cap).toFixed(2);
  const totalPnl = toSignificantDecimals(ws_message.position_msg.total_pnl);

  return `unrealised pnl: ${unrealizedPnl} | market cap: $${marketCap} | total pnl ${totalPnl}`;
}

function toSignificantDecimals(value: string, extraDigits = 3): string {
  const num = parseFloat(value);
  if (num === 0) return '0';

  const absNum = Math.abs(num);
  const firstSignificantDecimal = Math.max(0, Math.ceil(-Math.log10(absNum)));

  return num.toFixed(firstSignificantDecimal + extraDigits);
}
