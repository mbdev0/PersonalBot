export interface RPCGroup {
  rpc_group_name: string;
  rpcs: RPC[];
}

export interface RPC {
  http: string;
  ws: string;
}

export interface RPCGroupDashboard {
  rows: RPCGroupDashboardRow[];
}

export interface RPCGroupDashboardRow {
  rpc_group_name: string;
  num_of_rpcs: number;
}

export interface RPCGroupPost {
  rpc_group_name: string;
  rpc: string[];
}
