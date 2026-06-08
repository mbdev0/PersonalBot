import type {
  StrategyTaskPost,
  StrategyTaskPostDto,
  AFKStrategyTaskPost,
  AFKStrategyTaskPostDto,
  BuyStrategyTaskPost,
  BuyStrategyTaskPostDto,
  SellStrategyTaskPost,
  SellStrategyTaskPostDto,
  SpamStrategyTaskPostDto,
  SpamStrategyTaskPost,
} from '../../types/strategies/strategyTaskPost';
import { mapSellStrategyPostToDto } from './sellStrategyMapper';

export function mapStrategyTaskToPostDto(src: StrategyTaskPost): StrategyTaskPostDto {
  switch (src.trading_type) {
    case 'AFK':
      return mapAFKPostToDto(src);
    case 'BUY':
      return mapBuyPostToDto(src);
    case 'SELL':
      return mapSellPostToDto(src);
    case 'SPAM':
      return mapSpamPostToDto(src);
  }
}

function mapAFKPostToDto(src: AFKStrategyTaskPost): AFKStrategyTaskPostDto {
  return {
    ...src,
    sell_strategies: src.sell_strategies.map(mapSellStrategyPostToDto),
  };
}

function mapBuyPostToDto(src: BuyStrategyTaskPost): BuyStrategyTaskPostDto {
  return {
    ...src,
    sell_fee: src.sell_fee ?? 0,
    sell_strategies: src.sell_strategies.map(mapSellStrategyPostToDto),
  };
}

function mapSellPostToDto(src: SellStrategyTaskPost): SellStrategyTaskPostDto {
  return {
    ...src,
    trading_type: 'SELL',
  };
}

function mapSpamPostToDto(src: SpamStrategyTaskPost): SpamStrategyTaskPostDto {
  return {
    ...src,
    trading_type: 'SPAM',
  };
}
