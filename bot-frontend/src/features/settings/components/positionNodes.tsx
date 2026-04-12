import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import type { Settings } from '../types/settings';
import { Label } from '@/components/ui/label';
import { useState } from 'react';
import { usePostSettings } from '../hooks/useSettings';

interface PositionNodeProps {
  data: Settings;
}

export function PositionNode({ data }: PositionNodeProps) {
  const mutation = usePostSettings();
  const [httpNode, setHttpNode] = useState(data.position_nodes.http_node);
  const [wsNode, setWSNode] = useState(data.position_nodes.ws_node);

  const update = (patch: Partial<Settings>) => {
    if (!data) return;
    mutation.mutate({ ...data, ...patch });
  };

  return (
    <div className="rounded-lg border border-foreground/10 bg-foreground/3 p-6 space-y-6 flex-1">
      <p className="text-sm font-semibold text-foreground/50 uppercase">Position Update Node</p>

      <div className="flex flex-col space-y-2">
        <Label>HTTP</Label>
        <Input
          key={data.position_nodes.http_node}
          defaultValue={data.position_nodes.http_node}
          placeholder="HTTP Node"
          onChange={(e) => setHttpNode(e.target.value)}
        />

        <Label>Websocket</Label>
        <Input
          key={data.position_nodes.ws_node}
          defaultValue={data.position_nodes.ws_node}
          placeholder="WS Node"
          onChange={(e) => setWSNode(e.target.value)}
        />
      </div>

      <div className="flex justify-end mt-4">
        <Button
          className="hover:bg-green-700 hover:opacity-70"
          variant="default"
          onClick={() => update({ position_nodes: { http_node: httpNode, ws_node: wsNode } })}
        >
          Save
        </Button>
      </div>
    </div>
  );
}
