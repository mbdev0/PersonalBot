import { useQueryClient } from '@tanstack/react-query';
import { useCallback, useEffect, useRef, useState } from 'react';

import { POSITION_WS } from '@/config/urls';
import type { PositionWebsocketMessage, SendPositionWSMessage } from '../types/positionWebsocket';
import { key, positionStore } from '@/stores/positionStore';

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
      const data: PositionWebsocketMessage = JSON.parse(event.data);
      positionStore.setWsMessage(
        key(data.position_msg.strategy_id, data.position_msg.buy_task_id),
        prettifyWsMessage(data)
      );
    };

    ws.onclose = () => {
      console.log('disconnected');
      setWebsocketStatus(false);
    };

    return () => {
      ws.close();
    };
  }, [client]);

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
