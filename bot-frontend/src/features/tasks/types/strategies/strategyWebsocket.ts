import type { TaskDto } from '../task';

export interface StrategyWSMessage {
  strategy_msg?: StrategyMessage;
  error?: string;
}

interface StrategyMessage {
  id: number;
  event: string;
  task?: TaskDto;
  state?: string;
  message?: string;
}

export interface StrategySendWSMessage {
  type: 'Subscribe' | 'Unsubscribe';
  id: number;
}
