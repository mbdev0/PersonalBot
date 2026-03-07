import { useStrategyWebsocket } from '@/features/tasks/hooks/useStrategyWebsocket';
import { useTaskWebsocket } from '@/features/tasks/hooks/useTaskWebsocket';
import { useMemo, type ReactNode } from 'react';
import { WebSocketContext } from './websocketContext';

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const { send } = useTaskWebsocket();
  const { sendStrategyWSMessage } = useStrategyWebsocket();

  const value = useMemo(() => ({ send, sendStrategyWSMessage }), [send, sendStrategyWSMessage]);

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}
