interface ComputeUnitsEntryProps {
  computeUnits: string;
  onChange: (val: string) => void;
}

export function ComputeUnitsEntry({ computeUnits, onChange }: ComputeUnitsEntryProps) {
  return (
    <div className="compute_units">
      <h3>Compute Units</h3>
      <input
        type="text"
        name="compute_units"
        id="compute_units"
        placeholder="Compute Units"
        value={computeUnits}
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
