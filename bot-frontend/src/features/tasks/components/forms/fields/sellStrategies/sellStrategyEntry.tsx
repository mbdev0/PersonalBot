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
import { Plus, Trash2 } from 'lucide-react';
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

  const onClear = () => {
    setType(sellStrategyTypes[0].value);
    setValue(20);
    setSellAmount(50);
  };

  return (
    <div className="grid grid-cols-[140px_120px_120px_auto] grid-rows-2 gap-x-4 gap-y-2 p-4">
      <Label htmlFor="sell_strategy_type">{name} Type</Label>
      <Label htmlFor="value">{isPercentageType ? 'Value (%)' : 'Value'}</Label>
      <Label htmlFor="sell_amount">Sell Amount (%)</Label>
      <div></div>

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

      <Input
        type="number"
        placeholder="Value"
        value={value}
        min={0}
        onChange={(e) => setValue(e.target.valueAsNumber)}
      ></Input>
      <Input
        type="number"
        placeholder="Sell Amount"
        value={sellAmount}
        max={100}
        min={0}
        onChange={(e) => setSellAmount(e.target.valueAsNumber)}
      ></Input>

      <div className="flex pr-5 items-center gap-2">
        <Button type="button" size="sm" onClick={onButtonClick}>
          <Plus />
        </Button>

        <Button type="button" size="sm" onClick={onClear}>
          <Trash2 />
        </Button>
      </div>
    </div>
  );
}
