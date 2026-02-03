import { Coins, Wallet, Bot } from 'lucide-react';
import { ThemeProvider } from './components/themeProvider';
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarProvider,
} from './components/ui/sidebar';
import { TaskDashboard } from './features/tasks/components/taskDashboard';
import { BrowserRouter, Link, Route, Routes } from 'react-router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import WalletDashboard from './features/wallets/components/walletDashboard';

const queryClient = new QueryClient();

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <QueryClientProvider client={queryClient}>
        <BrowserRouter>
          <SidebarProvider defaultOpen={false}>
            <AppSidebar />

            <div className="min-h-screen bg-background">
              <div className="max-w-7xl mx-auto py-8">
                <header>
                  <h1 className="text-3xl font-bold">Trading Bot</h1>
                </header>
                <Routes>
                  <Route path="/" element={<TaskDashboard />} />
                  <Route path="/wallets" element={<WalletDashboard />} />
                </Routes>
              </div>
            </div>
          </SidebarProvider>
        </BrowserRouter>
      </QueryClientProvider>
    </ThemeProvider>
  );
}

export default App;

export function AppSidebar() {
  const navMain = [
    {
      title: 'Task Dashboard',
      url: '/',
      icon: Coins,
      isActive: true,
    },
    {
      title: 'Wallets',
      url: '/wallets',
      icon: Wallet,
      isActive: false,
    },
  ];

  return (
    <Sidebar
      collapsible="none"
      className="w-[calc(var(--sidebar-width-icon)+1px)] border-r min-h-screen"
    >
      <SidebarHeader>
        <SidebarMenu>
          <SidebarMenuItem>
            <SidebarMenuButton
              tooltip={{ children: 'Bot - by Mustafa' }}
              size="lg"
              asChild
              className="md:h-8 md:p-0"
            >
              <a href="#">
                <div className="bg-sidebar-primary text-sidebar-primary-foreground flex aspect-square size-8 items-center justify-center rounded-lg">
                  <Bot className="size-4" />
                </div>
              </a>
            </SidebarMenuButton>
          </SidebarMenuItem>
        </SidebarMenu>
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent className="px-1.5 md:px-0">
            <SidebarMenu>
              {navMain.map((item) => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton
                    tooltip={{
                      children: item.title,
                      hidden: false,
                    }}
                    isActive={location.pathname === item.url}
                    asChild
                  >
                    <Link to={item.url}>
                      {' '}
                      <item.icon />
                      <span>{item.title}</span>
                    </Link>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarFooter />
    </Sidebar>
  );
}
