import { Button } from '@/components/ui/button';
import { Field, FieldDescription, FieldGroup, FieldLabel } from '@/components/ui/field';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';

interface RPCGroupEntryProps {
  onCompletion: () => void;
}

export function RpcGroupEntry({ onCompletion }: RPCGroupEntryProps) {
  return (
    <div className="p-4 space-y-5">
      <FieldGroup>
        <Field className="max-w-1/5">
          <FieldLabel>RPC Group Name</FieldLabel>
          <FieldDescription>Enter your name for the RPC group</FieldDescription>
          <Input id="rpc-group-name" type="text" placeholder="RPC Group Name" />
        </Field>
        <Field>
          <FieldLabel>RPCs</FieldLabel>
          <FieldDescription>Format: HTTPS,WS</FieldDescription>
          <Textarea
            placeholder="HTTPS,WS - Comma Seperated + New Line"
            className="min-h-96"
          ></Textarea>
        </Field>
      </FieldGroup>

      <div className="space-x-2 flex justify-end">
        <Button
          onClick={() => {
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
