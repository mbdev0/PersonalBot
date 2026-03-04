import { ThemeProvider } from './components/themeProvider';
import { SidebarInset, SidebarProvider } from './components/ui/sidebar';
import { TaskDashboard } from './features/tasks/components/taskDashboard';
import { BrowserRouter, Route, Routes } from 'react-router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import WalletDashboard from './features/wallets/components/walletDashboard';
import { AppSidebar } from './components/appSidebar';
import { WebsocketProvider } from './context/websocketProvider';

const queryClient = new QueryClient();

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <QueryClientProvider client={queryClient}>
        <WebsocketProvider>
          <BrowserRouter>
            <SidebarProvider defaultOpen={false}>
              <AppSidebar />
              <SidebarInset className="h-svh overflow-hidden">
                <div className="flex justify-center w-full h-full overflow-y-auto">
                  <div className="max-w-11/12 w-full py-8">
                    <header className="mb-8">
                      <h1 className="text-3xl font-bold">Trading Bot</h1>
                    </header>
                    <Routes>
                      <Route path="/" element={<TaskDashboard />} />
                      <Route path="/wallets" element={<WalletDashboard />} />
                    </Routes>
                  </div>
                </div>
              </SidebarInset>
            </SidebarProvider>
          </BrowserRouter>
        </WebsocketProvider>
      </QueryClientProvider>
    </ThemeProvider>
  );
}

export default App;
