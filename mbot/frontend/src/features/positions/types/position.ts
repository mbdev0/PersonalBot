export interface PositionRow {
  total_pnl: string;
  average_market_cap_entry: string;
  average_market_cap_exit: string;
  coin: string;
  address_for_url: string;
}

export type PositionDashboard = PositionRow[];
