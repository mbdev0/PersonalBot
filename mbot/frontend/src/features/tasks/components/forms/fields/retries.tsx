import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

interface RetryEntryProps {
  isRetryEnabled: boolean;
  maxRetries: number;
  onMaxRetryChange: (retries: number) => void;
  retryDelayMs: number;
  onRetryDelayChange: (delay: number) => void;
}

export function RetryEntry({
  isRetryEnabled,
  maxRetries,
  onMaxRetryChange,
  retryDelayMs,
  onRetryDelayChange,
}: RetryEntryProps) {
  return (
    <>
      <div className="flex gap-2">
        <div className="flex-col">
          <Label>Max Retries</Label>
          <p className="text-xs text-foreground/40">
            Keep as 0 if you want to run infinite retries
          </p>

          <div className="mt-auto">
            <Input
              id="max_retries"
              type="number"
              value={maxRetries}
              min={0}
              disabled={isRetryEnabled}
              onChange={(e) => onMaxRetryChange(e.target.valueAsNumber)}
            ></Input>
          </div>
        </div>

        <div className="flex flex-col">
          <Label>Retry Delay(MS)</Label>

          <div className="mt-auto">
            <Input
              id="retry_delay"
              type="number"
              value={retryDelayMs}
              min={0}
              disabled={isRetryEnabled}
              onChange={(e) => onRetryDelayChange(e.target.valueAsNumber)}
            ></Input>
          </div>
        </div>
      </div>
    </>
  );
}
