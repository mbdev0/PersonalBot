import { useStateWebsocket } from '@/features/tasks/hooks/useStateWebsocket';
import type { SendTaskWSMessage } from '@/features/tasks/types/task_websocket';
import { createContext, useContext, type ReactNode } from 'react';

const WebSocketContext = createContext<{ send: (msg: SendTaskWSMessage) => void } | null>(null);

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const { send } = useStateWebsocket();

  return <WebSocketContext.Provider value={{ send }}>{children}</WebSocketContext.Provider>;
}

export function useWebsocketSend() {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error('no provider for context');
  }

  return context.send;
}
