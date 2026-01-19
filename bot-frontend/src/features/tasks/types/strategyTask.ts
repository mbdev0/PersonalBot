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

export { type StrategyTaskDto, type StrategyTask };
