interface SellEntryProps {
  sellAmount: string;
  onSellAmountChange: (val: string) => void;
  sellFee: string;
  onSellFeeChange: (val: string) => void;
}

export function SellEntry({
  sellAmount,
  onSellAmountChange,
  sellFee,
  onSellFeeChange,
}: SellEntryProps) {
  return (
    <div className="sell_settings">
      <div className="sell_amount">
        <h3>Sell Amount %</h3>
        <input
          type="text"
          name="sell_amountt"
          id="sell_amount"
          placeholder="Sell Amount"
          value={sellAmount}
          onChange={(e) => onSellAmountChange(e.target.value)}
        />
      </div>

      <div className="sell_fee">
        <h3>Sell Fee</h3>
        <input
          type="text"
          name="sell_fee"
          id="sell_fee"
          placeholder="Sell Fee"
          value={sellFee}
          onChange={(e) => onSellFeeChange(e.target.value)}
        />
      </div>
    </div>
  );
}
