import { Dialog, DialogContent } from '@/components/ui/dialog';
import { cn } from '@/lib/utils';

interface BotDialogProps {
  onClose: () => void;
  isOpen: boolean;
  className?: string;
  children: React.ReactNode;
}

export function BotDialog({ onClose, isOpen, className, children }: BotDialogProps) {
  return (
    <Dialog open={!!isOpen} onOpenChange={(open) => !open && onClose()}>
      <DialogContent
        className={cn(
          'no-scrollbar -mx-4 min-w-4/5 max-h-3/4 overflow-y-auto px-4 border-0',
          className
        )}
      >
        {children}
      </DialogContent>
    </Dialog>
  );
}
