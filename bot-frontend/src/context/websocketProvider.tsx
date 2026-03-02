import { useStrategyWebsocket } from '@/features/tasks/hooks/useStrategyWebsocket';
import { useTaskWebsocket } from '@/features/tasks/hooks/useTaskWebsocket';
import type { ReactNode } from 'react';
import { WebSocketContext } from './websocketContext';

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const { send } = useTaskWebsocket();
  const { sendStrategyWSMessage } = useStrategyWebsocket();

  return (
    <WebSocketContext.Provider value={{ send, sendStrategyWSMessage }}>
      {children}
    </WebSocketContext.Provider>
  );
}
