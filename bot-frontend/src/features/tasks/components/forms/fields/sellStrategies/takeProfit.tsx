import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { TakeProfitType, type SellStrategy } from '@/features/tasks/types/strategyTask';
import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectLabel,
  SelectItem,
} from '@/components/ui/select';
import { Plus } from 'lucide-react';
import { useState } from 'react';
import { Label } from '@/components/ui/label';

interface TakeProfitEntryProps {
  sellStrategyTypes: TakeProfitType[];
  onAddStrategy: (sellStrat: SellStrategy) => void;
}

export function TakeProfitEntry({ sellStrategyTypes, onAddStrategy }: TakeProfitEntryProps) {
  const [type, setType] = useState<TakeProfitType>(TakeProfitType.Marketcap);
  const [value, setValue] = useState(20);
  const [sellAmount, setSellAmount] = useState(50);

  const onButtonClick = () => {
    const entryValue = type === TakeProfitType.Percentage ? value / 100 : value;
    onAddStrategy({ type: type, value: entryValue, sell_amount: sellAmount });
    setType(TakeProfitType.Marketcap);
    setValue(20);
    setSellAmount(50);
  };

  return (
    <div className="flex flex-row justify-evenly items-end">
      <div className="flex flex-col items-center">
        <Label htmlFor="sell_strategy_type">Take Profit Type</Label>
        <Select value={type} onValueChange={(val) => setType(val as TakeProfitType)}>
          <SelectTrigger className="w-full max-w-48">
            <SelectValue placeholder="Select a strategy" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>Take Profit Strategies</SelectLabel>
              {sellStrategyTypes.map((strategyType) => (
                <SelectItem key={strategyType} value={strategyType}>
                  {generateTakeProfitLabel(strategyType)}
                </SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      <div className="flex flex-col items-center">
        <Label htmlFor="value">{generateValueLabel(type)}</Label>
        <Input
          type="number"
          placeholder="Value"
          value={value}
          onChange={(e) => setValue(e.target.valueAsNumber)}
        ></Input>
      </div>

      <div className="flex flex-col items-center">
        <Label htmlFor="sell_amount">Sell Amount (%)</Label>
        <Input
          type="number"
          placeholder="Sell Amount"
          value={sellAmount}
          max={100}
          min={0}
          onChange={(e) => setSellAmount(e.target.valueAsNumber)}
        ></Input>
      </div>

      <div className="flex flex-row gap-4">
        <Button type="button" onClick={onButtonClick}>
          <Plus />
        </Button>
      </div>
    </div>
  );
}

function generateTakeProfitLabel(label: string) {
  return label
    .replaceAll('take_profit_', '')
    .toLowerCase()
    .split(' ')
    .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
    .join(' ');
}

function generateValueLabel(type: TakeProfitType) {
  if (type === TakeProfitType.Percentage) {
    return 'Value (%)';
  }
  return 'Value';
}
