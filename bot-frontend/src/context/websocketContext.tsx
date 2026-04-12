import type { StrategySendWSMessage } from '@/features/tasks/types/strategies/strategyWebsocket';
import type { SendTaskWSMessage } from '@/features/tasks/types/taskWebsocket';
import { createContext } from 'react';

export const WebSocketContext = createContext<{
  send: (msg: SendTaskWSMessage) => void;
  sendStrategyWSMessage: (msg: StrategySendWSMessage) => void;
  websocketOpen: boolean;
} | null>(null);
