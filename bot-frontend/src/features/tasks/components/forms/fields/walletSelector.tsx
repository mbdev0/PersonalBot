import { Label } from '@/components/ui/label';
import type { Wallet, WalletDto } from '../../../../wallets/types/wallet';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';

interface WalletSelectorProps {
  selectedWallet: Wallet | null;
  onChange: (wallet: Wallet | null) => void;
  isPending: boolean;
  isError: boolean;
  data: Wallet[] | undefined;
  error: Error | null;
}

export function WalletSelector({
  selectedWallet,
  onChange,
  isPending,
  isError,
  data,
  error,
}: WalletSelectorProps) {
  if (isPending) return <div>Loading wallets...</div>;
  if (isError) return <div>Error: {error?.message}</div>;
  if (data?.length === 0) return <div>No wallets. Create one first!</div>;

  return (
    <div className="space-y-2">
      <Label htmlFor="wallet">Wallet</Label>
      <Select
        value={selectedWallet?.wallet_name}
        onValueChange={(e) => onChange(data?.find((w) => w.wallet_name === e) ?? null)}
      >
        <SelectTrigger id="wallet">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {data?.map((w) => (
            <SelectItem key={w.wallet_name} value={w.wallet_name}>
              {w.wallet_name}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
