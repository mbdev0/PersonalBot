import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';

import {
  Select,
  SelectTrigger,
  SelectValue,
  SelectContent,
  SelectGroup,
  SelectLabel,
  SelectItem,
} from '@/components/ui/select';
import { Check, Edit, Trash2, X } from 'lucide-react';
import { useState } from 'react';
import { Label } from '@/components/ui/label';
import type {
  SellStrategy,
  SellStrategyTypeOptions,
  SellStrategyType,
} from '@/features/tasks/types/sellStrategies';

interface SellStrategyEditProps {
  name: 'Take Profit' | 'Stop Loss';
  sellStrategy: SellStrategy;
  sellStrategyTypes: SellStrategyTypeOptions[];
  onSave: (sellStrat: SellStrategy) => void;
  onDelete: (id: string) => void;
}

export function SellStrategyEdit({
  name,
  sellStrategy,
  sellStrategyTypes,
  onSave,
  onDelete,
}: SellStrategyEditProps) {
  const [type, setType] = useState<string>(sellStrategy.type);
  const [value, setValue] = useState(sellStrategy.value);
  const [sellAmount, setSellAmount] = useState(sellStrategy.sell_amount);
  const [isEditable, setIsEditable] = useState(false);

  const isPercentageType = type.includes('percentage');

  const handleEdit = () => {
    setIsEditable(true);
  };

  const handleSave = () => {
    const entryValue = isPercentageType ? value / 100 : value;
    onSave({
      id: sellStrategy.id,
      type: type as SellStrategyType,
      value: entryValue,
      sell_amount: sellAmount,
    });
    setIsEditable(false);
  };

  const handleCancel = () => {
    setType(sellStrategy.type);
    setValue(sellStrategy.value);
    setSellAmount(sellStrategy.sell_amount);
    setIsEditable(false);
  };

  return (
    <div className="grid grid-cols-[140px_120px_120px_auto] grid-rows-2 gap-x-4 gap-y-2 p-4">
      <Label htmlFor="sell_strategy_type">{name} Type</Label>
      <Label htmlFor="value">{isPercentageType ? 'Value (%)' : 'Value'}</Label>
      <Label htmlFor="sell_amount">Sell Amount (%)</Label>
      <div></div>

      <Select disabled={!isEditable} value={type} onValueChange={setType}>
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
        value={isPercentageType ? value * 100 : value}
        onChange={(e) => setValue(e.target.valueAsNumber)}
        disabled={!isEditable}
        min={0}
      ></Input>

      <Input
        type="number"
        placeholder="Sell Amount"
        value={sellAmount}
        max={100}
        min={0}
        onChange={(e) => setSellAmount(e.target.valueAsNumber)}
        disabled={!isEditable}
      ></Input>

      <div className="flex pr-5 items-center gap-2">
        {!isEditable ? (
          <>
            <Button type="button" size="sm" onClick={handleEdit}>
              <Edit className="h-4 w-4" />
            </Button>
            <Button
              type="button"
              size="sm"
              variant="destructive"
              onClick={() => onDelete(sellStrategy.id)}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </>
        ) : (
          <>
            <Button type="button" size="sm" onClick={handleSave}>
              <Check className="h-4 w-4" />
            </Button>
            <Button type="button" size="sm" variant="outline" onClick={handleCancel}>
              <X className="h-4 w-4" />
            </Button>
          </>
        )}
      </div>
    </div>
  );
}
