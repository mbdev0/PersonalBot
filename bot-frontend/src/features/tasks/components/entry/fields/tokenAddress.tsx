interface TokenAddressEntryProps {
  onChange: (val: string) => void;
}

export function TokenAddressEntry({ onChange }: TokenAddressEntryProps) {
  return (
    <div className="token_address">
      <h3>Token Address</h3>
      <input
        type="text"
        name="token_address"
        id="token_address"
        placeholder="Token Address"
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
