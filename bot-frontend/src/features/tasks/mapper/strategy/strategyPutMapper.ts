import type {
  StrategyTaskPut,
  StrategyTaskPutDto,
  AFKStrategyTaskPut,
  AFKStrategyTaskPutDto,
  BuyStrategyTaskPut,
  BuyStrategyTaskPutDto,
  SellStrategyTaskPut,
  SellStrategyTaskPutDto,
} from '../../types/strategies/strategyTaskPut';
import { mapSellStrategyDtoToPost, mapSellStrategyPostToDto } from './sellStrategyMapper';

export function mapStrategyTaskToPutDto(src: StrategyTaskPut): StrategyTaskPutDto {
  switch (src.trading_type) {
    case 'AFK':
      return mapAFKPutToDto(src);
    case 'BUY':
      return mapBuyPutToDto(src);
    case 'SELL':
      return mapSellPutToDto(src);
  }
}

function mapAFKPutToDto(src: AFKStrategyTaskPut): AFKStrategyTaskPutDto {
  return {
    trading_type: 'AFK',
    compute_units: src.compute_units,
    slippage: src.slippage,
    wallet_name: src.wallet_name,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_fee: src.sell_fee,
    filters: src.filters,
    sell_strategies: src.sell_strategies.map(mapSellStrategyPostToDto),
  };
}

function mapBuyPutToDto(src: BuyStrategyTaskPut): BuyStrategyTaskPutDto {
  return {
    trading_type: 'BUY',
    compute_units: src.compute_units,
    slippage: src.slippage,
    wallet_name: src.wallet_name,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_fee: src.sell_fee,
    token_address: src.token_address,
    sell_strategies: src.sell_strategies.map(mapSellStrategyPostToDto),
  };
}

function mapSellPutToDto(src: SellStrategyTaskPut): SellStrategyTaskPutDto {
  return {
    trading_type: 'SELL',
    compute_units: src.compute_units,
    slippage: src.slippage,
    wallet_name: src.wallet_name,
    sell_amount: src.sell_amount,
    sell_fee: src.sell_fee,
    token_address: src.token_address,
  };
}

export function mapStrategyTaskPutDtoToPut(src: StrategyTaskPutDto, id: number): StrategyTaskPut {
  switch (src.trading_type) {
    case 'AFK':
      return {
        id,
        trading_type: 'AFK',
        compute_units: src.compute_units,
        slippage: src.slippage,
        wallet_name: src.wallet_name,
        buy_amount: src.buy_amount,
        buy_fee: src.buy_fee,
        sell_fee: src.sell_fee,
        filters: src.filters,
        sell_strategies: src.sell_strategies.map(mapSellStrategyDtoToPost),
      };

    case 'BUY':
      return {
        id,
        trading_type: 'BUY',
        compute_units: src.compute_units,
        slippage: src.slippage,
        wallet_name: src.wallet_name,
        buy_amount: src.buy_amount,
        buy_fee: src.buy_fee,
        sell_fee: src.sell_fee,
        token_address: src.token_address,
        sell_strategies: src.sell_strategies.map(mapSellStrategyDtoToPost),
      };

    case 'SELL':
      return {
        id,
        trading_type: 'SELL',
        compute_units: src.compute_units,
        slippage: src.slippage,
        wallet_name: src.wallet_name,
        sell_amount: src.sell_amount,
        sell_fee: src.sell_fee,
        token_address: src.token_address,
      };
  }
}
