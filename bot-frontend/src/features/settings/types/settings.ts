export interface Settings {
  discord_webhook: string;
  send_on_fail: boolean;
  send_on_success: boolean;
  position_nodes: PositionNodeSettings;
}

export interface PositionNodeSettings {
  http_node: string;
  ws_node: string;
}
