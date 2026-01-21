import type {
  StrategyTaskPost,
  StrategyTaskPostDto,
  StrategyTaskPut,
  StrategyTaskPutDto,
} from '../types/strategyTask';

export function MapStrategyToPostDto(src: StrategyTaskPost): StrategyTaskPostDto {
  return {
    tradingType: src.tradingType,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_fee: src.sell_fee,
    compute_units: src.compute_units,
    slippage: src.slippage,
    wallet_name: src.wallet_name,
    filters: src.filters,
    sell_strategies: src.sell_strategies,
  };
}
export function MapStrategyToPutDto(src: StrategyTaskPut): StrategyTaskPutDto {
  return {
    tradingType: src.tradingType,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_fee: src.sell_fee,
    compute_units: src.compute_units,
    slippage: src.slippage,
    wallet_name: src.wallet_name,
    filters: src.filters,
    sell_strategies: src.sell_strategies,
  };
}
