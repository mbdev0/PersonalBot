import { Card } from '@/components/ui/card';
import { Label } from '@/components/ui/label';
import {
  STOP_LOSS_OPTIONS,
  TAKE_PROFIT_OPTIONS,
  type SellStrategyCreate,
  type SellStrategy,
} from '@/features/tasks/types/strategyTask';
import { SellStrategyEdit } from './sellStrategyEdit';
import { SellStrategyEntry } from './sellStrategyEntry';

interface SellStrategiesProps {
  sellStrategies: SellStrategy[];
  setSellStrategies: (sellStrats: SellStrategy[]) => void;
}

export function SellStrategies({ sellStrategies, setSellStrategies }: SellStrategiesProps) {
  const onAddSellStrategy = (strat: SellStrategyCreate) => {
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
    <>
      <Label htmlFor="sell strategies">Sell Strategies</Label>
      <Card>
        <Card>
          <Label>Take Profit</Label>

          <div className="flex flex-col justify-evenly space-y-4">
            <SellStrategyEntry
              name="Take Profit"
              sellStrategyTypes={TAKE_PROFIT_OPTIONS}
              onAddStrategy={onAddSellStrategy}
            ></SellStrategyEntry>

            <div className="">
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
          </div>
        </Card>

        <Card>
          <Label>Stop Loss</Label>
          <div className="flex flex-col justify-evenly space-y-4">
            <SellStrategyEntry
              name="Stop Loss"
              sellStrategyTypes={STOP_LOSS_OPTIONS}
              onAddStrategy={onAddSellStrategy}
            ></SellStrategyEntry>

            <div className="space-y-2">
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
          </div>
        </Card>
      </Card>
    </>
  );
}
