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
      <DialogContent className={cn('max-w-2xl', className)}>{children}</DialogContent>
    </Dialog>
  );
}
