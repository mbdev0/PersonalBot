export interface RPCGroup {
  id: number;
  name: string;
  group: string;
  creation_time: string;
}

export interface RPC {
  http: string;
  ws: string;
}

export interface RPCGroupDashboardRow {
  id: number;
  name: string;
  number: number;
  creation_time: string;
}

export interface RPCGroupPost {
  name: string;
  group: string;
}

export interface RPCGroupPut {
  id: number;
  name: string;
  group: string;
}
