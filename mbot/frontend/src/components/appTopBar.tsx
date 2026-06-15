import { useLocation } from 'react-router';
import { cn } from '@/lib/utils';
import { useSolanaPrice } from '@/hooks/useMarket';
import type { Market } from '@/types/market';
import { useContext } from 'react';
import { WebSocketContext } from '@/context/websocketContext';
import { PAGE_LABELS } from '@/app/routes';

export function TopBar() {
  const { pathname } = useLocation();
  const { data } = useSolanaPrice();

  return (
    <header className="relative flex h-12 items-center px-4 border-b-2 border-foreground/5">
      <span className="absolute left-1/2 -translate-x-1/2 text-sm font-semibold">
        {PAGE_LABELS(pathname)}
      </span>

      <div className="ml-auto flex items-center gap-8">
        <SolPrice data={data}></SolPrice>
        <ServerStatus />
      </div>
    </header>
  );
}

function SolPrice({ data }: { data: Market | undefined }) {
  return (
    <span className="text-xs text-muted-foreground">
      SOL{' '}
      <span className="text-foreground font-medium">
        {!data ? '-' : `$${data.solana_price.toFixed(2)}`}
      </span>
    </span>
  );
}

function ServerStatus() {
  const ctx = useContext(WebSocketContext);

  if (!ctx) {
    console.error('no provider for context');
    return;
  }
  return (
    <div className="flex items-center gap-1.5">
      <div
        className={cn(
          'h-1.5 w-1.5 rounded-full ',
          ctx.websocketOpen ? 'animate-pulse bg-green-500' : 'bg-red-500'
        )}
      />
      <span className="text-xs text-muted-foreground">
        {ctx.websocketOpen ? 'Connected' : 'Disconnected'}
      </span>
    </div>
  );
}
