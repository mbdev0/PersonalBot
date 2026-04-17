export interface PositionRow {
  total_pnl: string;
  average_market_cap_entry: string;
  average_market_cap_exit: string;
  coin: string;
}

export type PositionDashboard = PositionRow[];
