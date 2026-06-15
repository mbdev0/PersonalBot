type Listener = () => void;

export const key = (strategyId: number, buyTaskId: number) => `${strategyId}:${buyTaskId}`;

let version = 0;
const positionMap = new Map<string, string>();

let listeners: Listener[] = [];

export const positionStore = {
  setWsMessage(id: string, message: string) {
    positionMap.set(id, message);
    version++;
    emitChanges();
  },
  subscribe(listener: Listener) {
    listeners = [...listeners, listener];
    return () => {
      listeners = listeners.filter((l) => l !== listener);
    };
  },
  getSnapshot() {
    return version;
  },
  get(id: string) {
    return positionMap.get(id);
  },
};

function emitChanges() {
  for (const listener of listeners) {
    listener();
  }
}
