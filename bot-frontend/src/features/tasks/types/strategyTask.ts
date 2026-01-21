interface StrategyTaskDto {
  tradingType: string;
  id: number;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  wallet_address: string;
  filters: filters;
  sell_strategies: sellStrategies[];
  sell_fee: number;
}

interface StrategyTask {
  tradingType: string;
  id: number;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  wallet_address: string;
  filters: filters;
  sell_strategies: sellStrategies[];
  sell_fee: number;
}

interface filters {
  has_website?: boolean;
  has_twitter?: boolean;
  has_telegram?: boolean;
  dev_wallet?: string;
}

interface sellStrategies {
  type: string;
  value: number;
  sell_amount: number;
}

//Sent to the API
interface StrategyTaskPostDto {
  tradingType: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters: filters;
  sell_strategies: sellStrategies[];
  sell_fee: number;
}

//This will be what is used inside the React App
interface StrategyTaskPost {
  tradingType: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters: filters;
  sell_strategies: sellStrategies[];
  sell_fee: number;
}

interface StrategyTaskPut {
  id: number;
  tradingType: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters: filters;
  sell_strategies: sellStrategies[];
  sell_fee: number;
}

interface StrategyTaskPutDto {
  tradingType: string;
  buy_amount: number;
  buy_fee: number;
  compute_units: number;
  slippage: number;
  wallet_name: string;
  filters: filters;
  sell_strategies: sellStrategies[];
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
