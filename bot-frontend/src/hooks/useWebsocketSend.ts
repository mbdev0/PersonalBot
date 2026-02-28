import { WebSocketContext } from '@/context/websocketContext';
import { useContext } from 'react';

export function useWebsocketSend() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('no provider for context');
  }
  return context.send;
}
