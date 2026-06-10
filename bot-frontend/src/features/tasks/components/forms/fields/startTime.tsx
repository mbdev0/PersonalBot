import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';

export interface StartTimeProps {
  isStartTimeEnabled: boolean;
  startTime: Date;
  onStartTimeChange: (date: Date) => void;
}

export function StartTime({ startTime, onStartTimeChange, isStartTimeEnabled }: StartTimeProps) {
  return (
    <div className="flex flex-col">
      <Label>Start Time</Label>
      <p className="text-xs text-foreground/40">hh:mm:ss</p>

      <div className="mt-auto">
        <Input
          type="time"
          disabled={isStartTimeEnabled}
          id="time-picker-optional"
          value={startTime.toTimeString().slice(0, 8)}
          step="1"
          onChange={(e) => {
            const date = e.target.valueAsDate;
            if (date == null) return;

            const next = new Date();
            next.setHours(date.getUTCHours(), date.getUTCMinutes(), date.getUTCSeconds(), 0);

            if (next <= new Date()) {
              next.setDate(next.getDate() + 1);
            }

            onStartTimeChange(next);
          }}
          className="appearance-none bg-background [&::-webkit-calendar-picker-indicator]:hidden [&::-webkit-calendar-picker-indicator]:appearance-none"
        />
      </div>
    </div>
  );
}
