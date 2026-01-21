interface SlippageEntryProps {
  slippage: string;
  onChange: (slippage: string) => void;
}

export function SlippageEntry({ slippage, onChange }: SlippageEntryProps) {
  return (
    <div className="slippage">
      <h3>Slippage</h3>
      <input
        type="text"
        name="slippage"
        id="slippage"
        placeholder="Slippage"
        value={slippage}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
