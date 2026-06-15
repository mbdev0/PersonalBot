import type {
  AFKStrategyTaskPost,
  BuyStrategyTaskPost,
  SellStrategyTaskPost,
  SpamStrategyTaskPost,
} from './strategyTaskPost';

export type AFKStrategyTaskPut = AFKStrategyTaskPost & { id: number };
export type BuyStrategyTaskPut = BuyStrategyTaskPost & { id: number };
export type SellStrategyTaskPut = SellStrategyTaskPost & { id: number };
export type SpamStrategyTaskPut = SpamStrategyTaskPost & { id: number };

export type StrategyTaskPut =
  | AFKStrategyTaskPut
  | BuyStrategyTaskPut
  | SellStrategyTaskPut
  | SpamStrategyTaskPut;
