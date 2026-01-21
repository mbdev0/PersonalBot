import type { Wallet, WalletDto } from '../../../../wallets/types/wallet';

interface WalletSelectorProps {
  selectedWallet: WalletDto | null;
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
    <div className="wallet">
      <h3>Wallet</h3>
      <select
        value={selectedWallet?.wallet_name}
        disabled={isPending || isError}
        onChange={(e) => onChange(data?.find((w) => w.wallet_name === e.target.value) ?? null)}
      >
        {data?.map((wallet) => (
          <option key={wallet.id} value={wallet.wallet_name}>
            {wallet.wallet_name}
          </option>
        ))}
      </select>
    </div>
  );
}
