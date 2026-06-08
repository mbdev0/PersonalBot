import type {
  StrategyTaskPut,
  StrategyTaskPutDto,
  AFKStrategyTaskPut,
  AFKStrategyTaskPutDto,
  BuyStrategyTaskPut,
  BuyStrategyTaskPutDto,
  SellStrategyTaskPut,
  SellStrategyTaskPutDto,
  SpamStrategyTaskPut,
  SpamStrategyTaskPutDto,
} from '../../types/strategies/strategyTaskPut';
import { mapSellStrategyPostToDto } from './sellStrategyMapper';

export function mapStrategyTaskToPutDto(src: StrategyTaskPut): StrategyTaskPutDto {
  switch (src.trading_type) {
    case 'AFK':
      return mapAFKPutToDto(src);
    case 'BUY':
      return mapBuyPutToDto(src);
    case 'SELL':
      return mapSellPutToDto(src);
    case 'SPAM':
      return mapSpamPutToDto(src);
  }
}

function mapAFKPutToDto(src: AFKStrategyTaskPut): AFKStrategyTaskPutDto {
  return {
    ...src,
    sell_strategies: src.sell_strategies.map(mapSellStrategyPostToDto),
  };
}

function mapBuyPutToDto(src: BuyStrategyTaskPut): BuyStrategyTaskPutDto {
  return {
    ...src,
    sell_fee: src.sell_fee ?? 0,
    sell_strategies: src.sell_strategies.map(mapSellStrategyPostToDto),
  };
}

function mapSellPutToDto(src: SellStrategyTaskPut): SellStrategyTaskPutDto {
  return {
    ...src,
  };
}

function mapSpamPutToDto(src: SpamStrategyTaskPut): SpamStrategyTaskPutDto {
  return {
    ...src,
  };
}
