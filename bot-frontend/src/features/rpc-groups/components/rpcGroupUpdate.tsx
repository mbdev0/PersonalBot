import { Button } from '@/components/ui/button';
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { useRef, useState } from 'react';
import { useRpcGroup, useUpdateRPCGroup } from '../hooks/rpcGroups';
import type { RPCGroup } from '../types/rpcGroup';

interface RPCGroupUpdateProps {
  editingRowId: number;
  onCompletion: () => void;
}

export function RpcGroupUpdate({ editingRowId, onCompletion }: RPCGroupUpdateProps) {
  const { isPending, isError, data, error } = useRpcGroup({ id: editingRowId });

  if (isPending) {
    return <div className="loading">Loading...</div>;
  }

  if (isError) {
    return <div className="error"> Error: {error.message}</div>;
  }

  if (data) {
    return <RPCGroupUpdateForm data={data} onCompletion={onCompletion}></RPCGroupUpdateForm>;
  }
}

interface RPCGroupUpdateFormProps {
  data: RPCGroup;
  onCompletion: () => void;
}

function RPCGroupUpdateForm({ data, onCompletion }: RPCGroupUpdateFormProps) {
  const updateMutation = useUpdateRPCGroup();
  const [groupName, SetGroupName] = useState(data.name);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  return (
    <>
      <div className="p-4 space-y-5">
        <FieldGroup>
          <Field className="max-w-1/5">
            <FieldLabel>RPC Group Name</FieldLabel>
            <FieldDescription>Enter your name for the RPC group</FieldDescription>
            <Input
              id="rpc-group-name"
              type="text"
              placeholder="RPC Group Name"
              key={data.name}
              value={groupName}
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
              className="min-h-96 whitespace-pre-wrap"
              defaultValue={data.group}
            ></Textarea>
          </Field>
        </FieldGroup>

        <div className="space-x-2 flex justify-end">
          <Button
            onClick={() => {
              updateMutation.mutate(
                {
                  id: data.id,
                  name: groupName,
                  group: textareaRef.current?.value ?? '',
                },
                {
                  onSuccess: () => {
                    onCompletion();
                  },
                }
              );
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
    </>
  );
}
