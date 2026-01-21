interface BuyEntryProps {
  buyAmount: string;
  onBuyAmountChange: (val: string) => void;
  buyFee: string;
  onBuyFeeChange: (val: string) => void;
}

export function BuyEntry({ buyAmount, onBuyAmountChange, buyFee, onBuyFeeChange }: BuyEntryProps) {
  return (
    <div className="buy_settings">
      <div className="buy_amount">
        <h3>Buy Amount</h3>
        <input
          type="text"
          name="buy_amount"
          id="buy_amount"
          placeholder="Buy Amount"
          value={buyAmount}
          onChange={(e) => onBuyAmountChange(e.target.value)}
        />
      </div>

      <div className="buy_fee">
        <h3>Buy Fee</h3>
        <input
          type="text"
          name="buy_fee"
          id="buy_fee"
          placeholder="Buy Fee"
          value={buyFee}
          onChange={(e) => onBuyFeeChange(e.target.value)}
        />
      </div>
    </div>
  );
}
