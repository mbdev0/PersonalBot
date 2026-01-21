interface TokenAddressEntryProps {
  value: string;
  onChange: (val: string) => void;
}

export function TokenAddressEntry({ value, onChange }: TokenAddressEntryProps) {
  return (
    <div className="token_address">
      <h3>Token Address</h3>
      <input
        type="text"
        name="token_address"
        id="token_address"
        placeholder="Token Address"
        value={value}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
