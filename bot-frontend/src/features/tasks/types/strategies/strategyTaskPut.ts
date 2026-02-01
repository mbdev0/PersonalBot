import type { Filters } from '../filters';
import type { SellStrategyDto, SellStrategyPost } from '../sellStrategies';
import type { TradingType } from './strategyTask';

interface BaseStrategyTaskPutDto {
  trading_type: TradingType;
  compute_units: number;
  slippage: number;
  wallet_name: string;
}

export interface AFKStrategyTaskPutDto extends BaseStrategyTaskPutDto {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategyDto[];
}

export interface BuyStrategyTaskPutDto extends BaseStrategyTaskPutDto {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  sell_strategies: SellStrategyDto[];
  token_address: string;
}

export interface SellStrategyTaskPutDto extends BaseStrategyTaskPutDto {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
}

export type StrategyTaskPutDto =
  | AFKStrategyTaskPutDto
  | BuyStrategyTaskPutDto
  | SellStrategyTaskPutDto;

interface BaseStrategyTaskPut {
  id: number;
  trading_type: TradingType;
  compute_units: number;
  slippage: number;
  wallet_name: string;
}

export interface AFKStrategyTaskPut extends BaseStrategyTaskPut {
  trading_type: 'AFK';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  filters: Filters;
  sell_strategies: SellStrategyPost[];
}

export interface BuyStrategyTaskPut extends BaseStrategyTaskPut {
  trading_type: 'BUY';
  buy_amount: number;
  buy_fee: number;
  sell_fee: number;
  sell_strategies: SellStrategyPost[];
  token_address: string;
}

export interface SellStrategyTaskPut extends BaseStrategyTaskPut {
  trading_type: 'SELL';
  sell_amount: number;
  sell_fee: number;
  token_address: string;
}

export type StrategyTaskPut = AFKStrategyTaskPut | BuyStrategyTaskPut | SellStrategyTaskPut;
