import type {
  StrategyTaskPost,
  StrategyTaskPostDto,
  AFKStrategyTaskPost,
  AFKStrategyTaskPostDto,
  BuyStrategyTaskPost,
  BuyStrategyTaskPostDto,
  SellStrategyTaskPost,
  SellStrategyTaskPostDto,
} from '../../types/strategies/strategyTaskPost';
import { mapSellStrategyPostToDto, mapSellStrategyDtoToPost } from './sellStrategyMapper';

export function mapStrategyTaskToPostDto(src: StrategyTaskPost): StrategyTaskPostDto {
  switch (src.trading_type) {
    case 'AFK':
      return mapAFKPostToDto(src);
    case 'BUY':
      return mapBuyPostToDto(src);
    case 'SELL':
      return mapSellPostToDto(src);
  }
}

function mapAFKPostToDto(src: AFKStrategyTaskPost): AFKStrategyTaskPostDto {
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

function mapBuyPostToDto(src: BuyStrategyTaskPost): BuyStrategyTaskPostDto {
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

function mapSellPostToDto(src: SellStrategyTaskPost): SellStrategyTaskPostDto {
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

export function mapStrategyTaskPostDtoToPost(src: StrategyTaskPostDto): StrategyTaskPost {
  switch (src.trading_type) {
    case 'AFK':
      return {
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
