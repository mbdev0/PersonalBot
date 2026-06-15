import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface SellEntryProps {
  sellAmount: number;
  onSellAmountChange: (val: number) => void;
  sellFee: number;
  onSellFeeChange: (val: number) => void;
}

export function SellEntry({
  sellAmount,
  onSellAmountChange,
  sellFee,
  onSellFeeChange,
}: SellEntryProps) {
  return (
    <>
      <div className="space-y-2">
        <Label htmlFor="sellAmount">Sell Amount (%)</Label>
        <Input
          id="sellAmount"
          type="number"
          value={sellAmount}
          max={100}
          min={0}
          onChange={(e) => onSellAmountChange(e.target.valueAsNumber)}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="sellFee">Sell Fee</Label>
        <Input
          id="sellFee"
          type="number"
          value={sellFee}
          onChange={(e) => onSellFeeChange(e.target.valueAsNumber)}
        />
      </div>
    </>
  );
}
