import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  type SellStrategyCreate,
  type SellStrategyType,
  type SellStrategyTypeOptions,
} from '@/features/tasks/types/strategyTask';
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

interface SellStrategyEntryProps {
  name: 'Take Profit' | 'Stop Loss';
  sellStrategyTypes: SellStrategyTypeOptions[];
  onAddStrategy: (sellStrat: SellStrategyCreate) => void;
}

export function SellStrategyEntry({
  name,
  sellStrategyTypes,
  onAddStrategy,
}: SellStrategyEntryProps) {
  const [type, setType] = useState<string>(sellStrategyTypes[0].value);
  const [value, setValue] = useState(20);
  const [sellAmount, setSellAmount] = useState(50);

  const isPercentageType = type.includes('percentage');

  const onButtonClick = () => {
    const entryValue = isPercentageType ? value / 100 : value;
    onAddStrategy({
      type: type as SellStrategyType,
      value: entryValue,
      sell_amount: sellAmount,
    });
    setType(sellStrategyTypes[0].value);
    setValue(20);
    setSellAmount(50);
  };

  return (
    <div className="flex flex-row justify-evenly items-end">
      <div className="flex flex-col items-center">
        <Label htmlFor="sell_strategy_type">{name} Type</Label>
        <Select value={type} onValueChange={setType}>
          <SelectTrigger className="w-full max-w-48">
            <SelectValue placeholder="Select a strategy" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectLabel>{name} Strategies</SelectLabel>
              {sellStrategyTypes.map((strategyType) => (
                <SelectItem key={strategyType.value} value={strategyType.value}>
                  {strategyType.label}
                </SelectItem>
              ))}
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      <div className="flex flex-col items-center">
        <Label htmlFor="value">{isPercentageType ? 'Value (%)' : 'Value'}</Label>
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
        <Button type="button" size="sm" onClick={onButtonClick}>
          <Plus />
        </Button>
      </div>
    </div>
  );
}
