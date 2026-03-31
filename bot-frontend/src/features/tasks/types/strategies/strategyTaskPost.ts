import type { Filters } from '../filters';
import type { SellStrategyDto, SellStrategyPost } from '../sellStrategies';
import type { TradingType } from './strategyTask';

interface BaseStrategyTaskCreateDto {
  trading_type: TradingType;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  rpc_group_id: number;
}

export interface AFKStrategyTaskPostDto extends BaseStrategyTaskCreateDto {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategyDto[];
}

export interface BuyStrategyTaskPostDto extends BaseStrategyTaskCreateDto {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  sell_strategies: SellStrategyDto[];
  token_address: string;
}

export interface SellStrategyTaskPostDto extends BaseStrategyTaskCreateDto {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
}

export type StrategyTaskPostDto =
  | AFKStrategyTaskPostDto
  | BuyStrategyTaskPostDto
  | SellStrategyTaskPostDto;

interface BaseStrategyTaskPost {
  trading_type: TradingType;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  rpc_group_id: number;
}

export interface AFKStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategyPost[];
}

export interface BuyStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee?: number;
  sell_strategies: SellStrategyPost[];
  token_address: string;
}

export interface SellStrategyTaskPost extends BaseStrategyTaskPost {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
}

export type StrategyTaskPost = AFKStrategyTaskPost | BuyStrategyTaskPost | SellStrategyTaskPost;
