import { useStrategyWebsocket } from '@/features/tasks/hooks/useStrategyWebsocket';
import { useTaskWebsocket } from '@/features/tasks/hooks/useTaskWebsocket';
import { useMemo, type ReactNode } from 'react';
import { WebSocketContext } from './websocketContext';
import { usePositionWebsocket } from '@/features/tasks/hooks/usePositionWebsocket';

export function WebsocketProvider({ children }: { children: ReactNode }) {
  const { sendPositionWsMessage, isPositionWsOpen } = usePositionWebsocket();
  const { send, taskWebsocketOpen } = useTaskWebsocket(sendPositionWsMessage);
  const { sendStrategyWSMessage, isStrategyWebsocketOpen } =
    useStrategyWebsocket(sendPositionWsMessage);

  const value = useMemo(() => {
    const websocketOpen = taskWebsocketOpen && isStrategyWebsocketOpen && isPositionWsOpen;
    return { send, sendStrategyWSMessage, sendPositionWsMessage, websocketOpen };
  }, [
    send,
    sendStrategyWSMessage,
    sendPositionWsMessage,
    taskWebsocketOpen,
    isStrategyWebsocketOpen,
    isPositionWsOpen,
  ]);

  return <WebSocketContext.Provider value={value}>{children}</WebSocketContext.Provider>;
}
