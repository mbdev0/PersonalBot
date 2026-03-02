import { useStrategyWebsocket } from '@/features/tasks/hooks/useStrategyWebsocket';
import { useTaskWebsocket } from '@/features/tasks/hooks/useTaskWebsocket';
import type { StrategySendWSMessage } from '@/features/tasks/types/strategies/strategyWebsocket';
import type { SendTaskWSMessage } from '@/features/tasks/types/taskWebsocket';
import { createContext, type ReactNode } from 'react';

export const WebSocketContext = createContext<{
  send: (msg: SendTaskWSMessage) => void;
  sendStrategyWSMessage: (msg: StrategySendWSMessage) => void;
} | null>(null);


