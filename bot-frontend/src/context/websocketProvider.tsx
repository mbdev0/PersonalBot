import { useStrategyWebsocket } from '@/features/tasks/hooks/useStrategyWebsocket';
import { useTaskWebsocket } from '@/features/tasks/hooks/useTaskWebsocket';
import { useMemo, type ReactNode } from 'react';
import { WebSocketContext } from './websocketContext';

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const { send, taskWebsocketOpen } = useTaskWebsocket();
  const { sendStrategyWSMessage, isStrategyWebsocketOpen } = useStrategyWebsocket();

  const value = useMemo(() => {
    const websocketOpen = taskWebsocketOpen && isStrategyWebsocketOpen;
    return { send, sendStrategyWSMessage, websocketOpen };
  }, [send, sendStrategyWSMessage, taskWebsocketOpen, isStrategyWebsocketOpen]);

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}
