export interface PositionWebsocketMessage {
  position_msg: PositionMessage;
}

interface PositionMessage {
  message_type: 'PositionCreated' | 'PositionUpdate' | 'PositionStopped';
  buy_task_id: number;
  strategy_id: number;
  unrealized_profit: string;
  realized_profit: string;
  total_pnl: string;
  market_cap: string;
  current_price: string;
  initial_solana_amount: string;
  remaining_tokens: string;
  message: 'Position Created' | 'Position Closed' | 'Position Update';
}

export interface SendPositionWSMessage {
  type: 'Subscribe' | 'Unsubscribe';
  id: number;
}
