import { cn } from '@/lib/utils';
import { type ReactNode } from 'react';

interface PageHeaderProps {
  children: ReactNode;
  className?: string;
}

export function PageHeader({ children, className }: PageHeaderProps) {
  return (
    <div>
      <h1 className={cn('text-2xl font-semibold text-foreground tracking-tight', className)}>
        {children}
      </h1>
    </div>
  );
}
