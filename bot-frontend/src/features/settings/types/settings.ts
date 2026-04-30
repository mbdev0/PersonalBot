export interface Settings {
  discord_webhook: string;
  send_on_fail: boolean;
  send_on_success: boolean;
  position_nodes: PositionNodeSettings;
  quick_sell_buttons: QuickSellButtons;
}

export interface PositionNodeSettings {
  http_node: string;
  ws_node: string;
}

export interface QuickSellButtons {
  button_1: number;
  button_2: number;
  button_3: number;
  button_4: number;
}
