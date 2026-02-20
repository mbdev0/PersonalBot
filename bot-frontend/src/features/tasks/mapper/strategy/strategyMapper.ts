import type { StrategyTaskDto, StrategyTask } from '../../types/strategies/strategyTask';

export function mapStrategyTaskDtoToStrategyTask(src: StrategyTaskDto): StrategyTask {
  switch (src.trading_type) {
    case 'AFK':
      return {
        trading_type: 'AFK',
        id: src.id,
        wallet_name: src.wallet_name,
        wallet_address: src.wallet_address,
        compute_units: src.compute_units,
        slippage: src.slippage,
        buy_amount: src.buy_amount,
        buy_fee: src.buy_fee,
        sell_fee: src.sell_fee,
        filters: src.filters,
        sell_strategies: src.sell_strategies,
      };

    case 'BUY':
      return {
        trading_type: 'BUY',
        id: src.id,
        wallet_name: src.wallet_name,
        wallet_address: src.wallet_address,
        compute_units: src.compute_units,
        slippage: src.slippage,
        buy_amount: src.buy_amount,
        buy_fee: src.buy_fee,
        sell_fee: src.sell_fee,
        token_address: src.token_address,
        sell_strategies: src.sell_strategies,
        buy_task_id: src.buy_task_id,
        position_id: src.position_id,
      };

    case 'SELL':
      return {
        trading_type: 'SELL',
        id: src.id,
        wallet_name: src.wallet_name,
        wallet_address: src.wallet_address,
        compute_units: src.compute_units,
        slippage: src.slippage,
        sell_amount: src.sell_amount,
        sell_fee: src.sell_fee,
        token_address: src.token_address,
        sell_task_id: src.sell_task_id,
      };
  }
}
