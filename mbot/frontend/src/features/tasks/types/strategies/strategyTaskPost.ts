import type { Filters } from '../filters';
import type { SellStrategy } from '../sellStrategies';
import type { TradingType } from './strategyTask';

interface BaseStrategyTaskPost {
  program: string;
  trading_type: TradingType;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  rpc_group_id: number;
  retries?: number;
  retry_delay_ms?: number;
}

export interface AFKStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategy[];
}

export interface BuyStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee?: number;
  sell_strategies: SellStrategy[];
  token_address: string;
}

export interface SellStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
}

export interface SpamStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'SPAM';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  token_address: string;
  no_of_tasks: number;
  start_time: number;
}

export type StrategyTaskPost =
  | AFKStrategyTaskPost
  | BuyStrategyTaskPost
  | SellStrategyTaskPost
  | SpamStrategyTaskPost;
