interface StrategyTaskDto {
  tradingType: string;
  id: number;
  buyAmount: number;
  buyFee: number;
  computeUnits: number;
  slipapge: number;
  walletName: string;
  walletAddress: string;
  filters: filters;
  sellStrategies: sellStrategies[];
  sellFee: number;
}

interface StrategyTask {
  tradingType: string;
  id: number;
  buyAmount: number;
  buyFee: number;
  computeUnits: number;
  slipapge: number;
  walletName: string;
  walletAddress: string;
  filters: filters;
  sellStrategies: sellStrategies[];
  sellFee: number;
}

interface filters {
  hasWebsite: boolean | undefined;
  hasTwitter: boolean | undefined;
  hasTelegram: boolean | undefined;
  devWallet: string | undefined;
}

interface sellStrategies {
  type: string;
  value: number;
  sellAmount: number;
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
  sell_strategies: sellStrategies;
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
  sell_strategies: sellStrategies;
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
  sell_strategies: sellStrategies;
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
  sell_strategies: sellStrategies;
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
