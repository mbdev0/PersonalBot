import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface SlippageEntryProps {
  slippage: number;
  onChange: (slippage: number) => void;
}

export function SlippageEntry({ slippage, onChange }: SlippageEntryProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor="slippage">Slippage (%)</Label>
      <Input
        type="number"
        name="slippage"
        id="slippage"
        placeholder="Slippage"
        value={slippage}
        onChange={(e) => onChange(e.target.valueAsNumber)}
      ></Input>
    </div>
  );
}
