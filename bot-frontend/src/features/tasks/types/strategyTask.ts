interface StrategyTaskDto {
  trading_type: string;
  id: number;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  wallet_address: string;
  filters: Filters;
  sell_strategies: SellStrategy[];
  sell_fee: number;
}

interface StrategyTask {
  trading_type: string;
  id: number;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  wallet_address: string;
  filters: Filters;
  sell_strategies: SellStrategy[];
  sell_fee: number;
}

export interface Filters {
  has_website?: boolean;
  has_twitter?: boolean;
  has_telegram?: boolean;
  dev_wallet?: string;
}

export interface SellStrategy {
  type: SellStrategyType;
  value: number;
  sell_amount: number;
}

export type SellStrategyType = TakeProfitType | StopLossType;

export enum TakeProfitType {
  Price = 'take_profit_price',
  Percentage = 'take_profit_percentage',
  Marketcap = 'take_profit_marketcap',
}

export enum StopLossType {
  Price = 'stop_loss_price',
  Percentage = 'stop_loss_percentage',
  Marketcap = 'stop_loss_marketcap',
}

export enum SellStrategyTypeLabel {}

//Sent to the API
interface StrategyTaskPostDto {
  trading_type: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters?: Filters;
  sell_strategies: SellStrategy[];
  sell_fee: number;
}

//This will be what is used inside the React App
interface StrategyTaskPost {
  trading_type: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters?: Filters;
  sell_strategies: SellStrategy[];
  sell_fee: number;
}

interface StrategyTaskPut {
  id: number;
  trading_type: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters: Filters;
  sell_strategies: SellStrategy[];
  sell_fee: number;
}

interface StrategyTaskPutDto {
  trading_type: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters: Filters;
  sell_strategies: SellStrategy[];
  sell_fee: number;
}

export {
  type StrategyTaskDto,
  type StrategyTask,
  type StrategyTaskPost,
  type StrategyTaskPostDto,
  type StrategyTaskPut,
  type StrategyTaskPutDto,
};
