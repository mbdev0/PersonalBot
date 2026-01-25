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
      <DialogContent className={cn('min-w-7xl', className)}>{children}</DialogContent>
    </Dialog>
  );
}
