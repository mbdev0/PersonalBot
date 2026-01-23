import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface BuyEntryProps {
  buyAmount: number;
  onBuyAmountChange: (val: number) => void;
  buyFee: number;
  onBuyFeeChange: (val: number) => void;
}

export function BuyEntry({ buyAmount, onBuyAmountChange, buyFee, onBuyFeeChange }: BuyEntryProps) {
  return (
    <>
      <div className="space-y-2">
        <Label htmlFor="buyAmount">Buy Amount</Label>
        <Input
          id="buyAmount"
          type="number"
          value={buyAmount}
          onChange={(e) => onBuyAmountChange(e.target.valueAsNumber)}
        />
      </div>

      <div className="space-y-2">
        <Label htmlFor="buyFee">Buy Fee</Label>
        <Input
          id="buyFee"
          type="number"
          value={buyFee}
          onChange={(e) => onBuyFeeChange(e.target.valueAsNumber)}
        />
      </div>
    </>
  );
}
