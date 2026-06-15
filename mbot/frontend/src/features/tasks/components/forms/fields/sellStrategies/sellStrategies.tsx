import { Card } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import { SellStrategyEdit } from './sellStrategyEdit';
import { SellStrategyEntry } from './sellStrategyEntry';
import {
  STOP_LOSS_OPTIONS,
  TAKE_PROFIT_OPTIONS,
  type SellStrategy,
  type SellStrategyPost,
} from '@/features/tasks/types/sellStrategies';

interface SellStrategiesProps {
  sellStrategies: SellStrategy[];
  setSellStrategies: (sellStrats: SellStrategy[]) => void;
}

export function SellStrategies({ sellStrategies, setSellStrategies }: SellStrategiesProps) {
  const onAddSellStrategy = (strat: SellStrategyPost) => {
    const strategyWithId: SellStrategy = {
      id: crypto.randomUUID(),
      ...strat,
    };

    setSellStrategies([strategyWithId, ...sellStrategies]);
  };

  const handleSaveStrategy = (updated: SellStrategy) => {
    setSellStrategies(sellStrategies.map((s) => (s.id === updated.id ? updated : s)));
  };

  const handleDeleteStrategy = (id: string) => {
    setSellStrategies(sellStrategies.filter((s) => s.id !== id));
  };

  const takeProfitStrategies = sellStrategies.filter((s) => s.type.startsWith('take_profit_'));
  const stopLossStrategies = sellStrategies.filter((s) => s.type.startsWith('stop_loss_'));

  return (
    <div>
      <Card className="p-2">
        <h2>Sell Strategies</h2>
        <div className="space-y-0">
          <Label htmlFor="take_profit">Take Profit</Label>
          <SellStrategyEntry
            name="Take Profit"
            sellStrategyTypes={TAKE_PROFIT_OPTIONS}
            onAddStrategy={onAddSellStrategy}
          ></SellStrategyEntry>

          {takeProfitStrategies.map((strategy) => (
            <SellStrategyEdit
              key={strategy.id}
              name="Take Profit"
              sellStrategy={strategy}
              sellStrategyTypes={TAKE_PROFIT_OPTIONS}
              onSave={handleSaveStrategy}
              onDelete={handleDeleteStrategy}
            />
          ))}
        </div>

        <div className="space-y-0">
          <Label htmlFor="stop_loss">Stop Loss</Label>
          <SellStrategyEntry
            name="Stop Loss"
            sellStrategyTypes={STOP_LOSS_OPTIONS}
            onAddStrategy={onAddSellStrategy}
          ></SellStrategyEntry>

          {stopLossStrategies.map((strategy) => (
            <SellStrategyEdit
              key={strategy.id}
              name="Stop Loss"
              sellStrategy={strategy}
              sellStrategyTypes={STOP_LOSS_OPTIONS}
              onSave={handleSaveStrategy}
              onDelete={handleDeleteStrategy}
            />
          ))}
        </div>
      </Card>
    </div>
  );
}
