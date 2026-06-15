import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface TokenAddressEntryProps {
  value: string;
  onChange: (val: string) => void;
}

export function TokenAddressEntry({ value, onChange }: TokenAddressEntryProps) {
  return (
    <div className="space-y-2">
      <Label htmlFor="tokenAddress">Token Address</Label>
      <Input
        id="tokenAddress"
        type="text"
        value={value}
        placeholder="Token Address"
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
