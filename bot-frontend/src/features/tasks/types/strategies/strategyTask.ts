import type { Filters } from '../filters';
import type { SellStrategyDto } from '../sellStrategies';

export type TradingType = 'AFK' | 'BUY' | 'SELL' | 'SPAM';

interface BaseStrategyTaskDto {
  trading_type: TradingType;
  id: number;
  wallet_name: string;
  wallet_address: string;
  compute_units: number;
  slippage: number;
  rpc_group: string;
  rpc_group_id: number;
  program: string;
  retries: number;
  retry_delay_ms: number;
}

export interface AFKStrategyTaskDto extends BaseStrategyTaskDto {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategyDto[];
}

export interface BuyStrategyTaskDto extends BaseStrategyTaskDto {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  sell_strategies: SellStrategyDto[];
  token_address: string;
  buy_task_id: number;
  position_id: number;
}

export interface SellStrategyTaskDto extends BaseStrategyTaskDto {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
  sell_task_id: number;
}

export interface SpamStrategyTaskDto extends BaseStrategyTaskDto {
  trading_type: 'SPAM';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  token_address: string;
  no_of_tasks: number;
  start_time: number;
}

export type StrategyTaskDto =
  | AFKStrategyTaskDto
  | BuyStrategyTaskDto
  | SellStrategyTaskDto
  | SpamStrategyTaskDto;

interface BaseStrategyTask {
  trading_type: TradingType;
  id: number;
  wallet_name: string;
  wallet_address: string;
  compute_units: number;
  slippage: number;
  rpc_group: string;
  rpc_group_id: number;
  program: string;
  retries: number;
  retry_delay_ms: number;
}

export interface AFKStrategyTask extends BaseStrategyTask {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategyDto[];
}

export interface BuyStrategyTask extends BaseStrategyTask {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  sell_strategies: SellStrategyDto[];
  token_address: string;
  buy_task_id: number;
  position_id: number;
}

export interface SellStrategyTask extends BaseStrategyTask {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
  sell_task_id: number;
}

export interface SpamStrategyTask extends BaseStrategyTask {
  trading_type: 'SPAM';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  token_address: string;
  no_of_tasks: number;
  start_time: number;
}

export type StrategyTask = AFKStrategyTask | BuyStrategyTask | SellStrategyTask | SpamStrategyTask;

export enum StrategyTaskState {
  create = 'CREATED',
  failed = 'FAILED',
  success = 'SUCCESS',
  cancelled = 'CANCELLED',
}
