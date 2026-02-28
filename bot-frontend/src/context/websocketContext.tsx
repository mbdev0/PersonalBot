import { useStateWebsocket } from '@/features/tasks/hooks/useStateWebsocket';
import type { SendTaskWSMessage } from '@/features/tasks/types/task_websocket';
import { createContext, type ReactNode } from 'react';

export const WebSocketContext = createContext<{ send: (msg: SendTaskWSMessage) => void } | null>(
  null
);

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const { send } = useStateWebsocket();

  return <WebSocketContext.Provider value={{ send }}>{children}</WebSocketContext.Provider>;
}
