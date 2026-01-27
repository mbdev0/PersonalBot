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
  sell_strategies: SellStrategyDto[];
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
  id: string;
  type: SellStrategyType;
  value: number;
  sell_amount: number;
}

export type SellStrategyCreate = Omit<SellStrategy, 'id'>;

export interface SellStrategyDto {
  id: string;
  type: SellStrategyType;
  value: number;
  sell_amount: number;
}

export const SELL_STRATEGY_TYPES = {
  TAKE_PROFIT_PERCENTAGE: 'take_profit_percentage',
  TAKE_PROFIT_MARKETCAP: 'take_profit_marketcap',
  TAKE_PROFIT_PRICE: 'take_profit_price',
  STOP_LOSS_PERCENTAGE: 'stop_loss_percentage',
  STOP_LOSS_MARKETCAP: 'stop_loss_marketcap',
  STOP_LOSS_PRICE: 'stop_loss_price',
} as const;

export type SellStrategyType =
  | 'take_profit_percentage'
  | 'take_profit_marketcap'
  | 'take_profit_price'
  | 'stop_loss_percentage'
  | 'stop_loss_marketcap'
  | 'stop_loss_price';

export interface SellStrategyTypeOptions {
  value: string;
  label: string;
}

export const TAKE_PROFIT_OPTIONS: SellStrategyTypeOptions[] = [
  { value: SELL_STRATEGY_TYPES.TAKE_PROFIT_PERCENTAGE, label: 'Percentage' },
  { value: SELL_STRATEGY_TYPES.TAKE_PROFIT_MARKETCAP, label: 'Market Cap' },
  { value: SELL_STRATEGY_TYPES.TAKE_PROFIT_PRICE, label: 'Price' },
] as const;

export const STOP_LOSS_OPTIONS: SellStrategyTypeOptions[] = [
  { value: SELL_STRATEGY_TYPES.STOP_LOSS_PERCENTAGE, label: 'Percentage' },
  { value: SELL_STRATEGY_TYPES.STOP_LOSS_MARKETCAP, label: 'Market Cap' },
  { value: SELL_STRATEGY_TYPES.STOP_LOSS_PRICE, label: 'Price' },
] as const;

//Sent to the API
interface StrategyTaskPostDto {
  trading_type: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters?: Filters;
  sell_strategies: SellStrategyDto[];
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
  sell_strategies: SellStrategyCreate[];
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
  sell_strategies: SellStrategyDto[];
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
