import type { StrategyTaskDto, StrategyTask } from '../../types/strategies/strategyTask';

export function mapStrategyTaskDtoToStrategyTask(src: StrategyTaskDto): StrategyTask {
  return {
    ...src,
  };
}
