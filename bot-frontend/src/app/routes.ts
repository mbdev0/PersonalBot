import { PositionDashboard } from '@/features/positions/components/positionDashboard';
import { RPCGroupDashboard } from '@/features/rpc-groups/components/rpcGroupsDashboard';
import { SettingsDashboard } from '@/features/settings/components/settingsDashboard';
import { TaskDashboard } from '@/features/tasks/components/taskDashboard';
import WalletDashboard from '@/features/wallets/components/walletDashboard';
import { ChartLine, Coins, Router, Settings, Wallet } from 'lucide-react';

interface RouteConfig {
  title: string;
  icon: React.ElementType;
  element: React.ComponentType;
  nav?: 'main' | 'footer';
}

export const RouteData: Record<string, RouteConfig> = {
  '/': {
    title: 'Task Dashboard',
    icon: Coins,
    element: TaskDashboard,
    nav: 'main',
  },
  '/wallets': {
    title: 'Wallets',
    icon: Wallet,
    element: WalletDashboard,
    nav: 'main',
  },
  '/settings': {
    title: 'Settings',
    icon: Settings,
    element: SettingsDashboard,
    nav: 'footer',
  },
  '/rpc': {
    title: 'RPC Groups',
    icon: Router,
    element: RPCGroupDashboard,
    nav: 'main',
  },
  '/positions': {
    title: 'Positions',
    icon: ChartLine,
    element: PositionDashboard,
    nav: 'main',
  },
};

//top bar labels
export const PAGE_LABELS = (path: string) => RouteData[path].title ?? '';

//sidebar routes
export const NavMain = Object.entries(RouteData)
  .filter((item) => item[1].nav === 'main')
  .map((item) => {
    return { title: item[1].title, url: item[0], icon: item[1].icon };
  });

export const NavFooter = Object.entries(RouteData)
  .filter((item) => item[1].nav === 'footer')
  .map((item) => {
    return { title: item[1].title, url: item[0], icon: item[1].icon };
  });

//app routes
export const RouteElements = Object.entries(RouteData).map(([url, route]) => {
  return { url, element: route.element };
});
