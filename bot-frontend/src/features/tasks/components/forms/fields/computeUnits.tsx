import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface ComputeUnitsEntryProps {
  computeUnits: number;
  onChange: (val: number) => void;
}

export function ComputeUnitsEntry({ computeUnits, onChange }: ComputeUnitsEntryProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor="compute_units">Compute Units</Label>
      <Input
        type="number"
        id="compute_units"
        placeholder="Compute Units"
        value={computeUnits}
        onChange={(e) => onChange(e.target.valueAsNumber)}
      ></Input>
    </div>
  );
}
