import type { TaskDto } from '../task';

export interface StrategyWSMessage {
  strategy_msg?: StrategyMessage;
  error?: string;
}

interface StrategyMessage {
  id: number;
  event: StrategyEventType;
  task?: TaskDto;
  state?: string;
  message?: string;
}

export const StrategyEventType = {
  STATUS_UPDATE: 'STATUS_UPDATE',
  MESSAGE_UPDATE: 'MESSAGE_UPDATE',
  TASK_CREATION: 'TASK_CREATION',
} as const;

export type StrategyEventType = (typeof StrategyEventType)[keyof typeof StrategyEventType];

export interface StrategySendWSMessage {
  type: 'Subscribe' | 'Unsubscribe';
  id: number;
}
