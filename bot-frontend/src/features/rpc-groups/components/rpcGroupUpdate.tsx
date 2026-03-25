import { Button } from '@/components/ui/button';
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useRef, useState } from 'react';
import type { RPCGroup } from '../types/rpcGroup';

interface RPCGroupUpdateProps {
  editingRow: RPCGroup;
  onCompletion: () => void;
}

export function RpcGroupUpdate({ editingRow, onCompletion }: RPCGroupUpdateProps) {
  const [rpcGroupName, SetGroupName] = useState(editingRow.rpc_group_name);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  return (
    <div className="p-4 space-y-5">
      <FieldGroup>
        <Field className="max-w-1/5">
          <FieldLabel>RPC Group Name</FieldLabel>
          <FieldDescription>Enter your name for the RPC group</FieldDescription>
          <Input
            id="rpc-group-name"
            type="text"
            placeholder="RPC Group Name"
            value={rpcGroupName}
            onChange={(e) => {
              SetGroupName(e.currentTarget.value);
            }}
          />
        </Field>
        <Field>
          <FieldLabel>RPCs</FieldLabel>
          <FieldDescription>Format: HTTPS,WS</FieldDescription>
          <Textarea
            ref={textareaRef}
            placeholder="HTTPS,WS - Comma Seperated + New Line"
            className="min-h-96"
            defaultValue={setupTextArea(editingRow.rpcs)}
          ></Textarea>
        </Field>
      </FieldGroup>

      <div className="space-x-2 flex justify-end">
        <Button
          onClick={() => {
            // const updatedText = textareaRef.current?.value ?? '';
            // console.log(updatedText);
            onCompletion();
          }}
        >
          Save
        </Button>
        <Button
          onClick={() => {
            onCompletion();
          }}
        >
          Cancel
        </Button>
      </div>
    </div>
  );
}

function setupTextArea(
  rpcPairs: {
    http: string;
    ws: string;
  }[]
): string {
  let string = '';
  rpcPairs.forEach((rpc) => {
    string += rpc.http + ',' + rpc.ws + '\n';
  });

  return string;
}
