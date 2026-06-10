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
  const { id, ...rest } = src;
  void id;
  return { ...rest, sell_strategies: rest.sell_strategies.map(mapSellStrategyPostToDto) };
}

function mapSellPutToDto(src: SellStrategyTaskPut): SellStrategyTaskPutDto {
  const { id, ...rest } = src;
  void id;
  return { ...rest };
}

function mapSpamPutToDto(src: SpamStrategyTaskPut): SpamStrategyTaskPutDto {
  const { id, ...rest } = src;
  void id;
  return { ...rest };
}

function mapBuyPutToDto(src: BuyStrategyTaskPut): BuyStrategyTaskPutDto {
  const { id, ...rest } = src;
  void id;
  return {
    ...rest,
    sell_fee: rest.sell_fee ?? 0,
    sell_strategies: rest.sell_strategies.map(mapSellStrategyPostToDto),
  };
}
