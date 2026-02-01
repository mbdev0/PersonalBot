import type { SellStrategyPost, SellStrategyDto } from '../../types/sellStrategies';

export function mapSellStrategyPostToDto(src: SellStrategyPost): SellStrategyDto {
  return {
    id: crypto.randomUUID(),
    type: src.type,
    value: src.value,
    sell_amount: src.sell_amount,
  };
}

export function mapSellStrategyDtoToPost(src: SellStrategyDto): SellStrategyPost {
  return {
    type: src.type,
    value: src.value,
    sell_amount: src.sell_amount,
  };
}
