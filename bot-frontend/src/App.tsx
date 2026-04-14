import { ThemeProvider } from './components/themeProvider';
import { SidebarInset, SidebarProvider } from './components/ui/sidebar';
import { BrowserRouter } from 'react-router';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { AppSidebar } from './components/appSidebar';
import { WebsocketProvider } from './context/websocketProvider';
import { TopBar } from './components/appTopBar';
import { Toaster } from 'sonner';
import { AppRoutes } from './app/appRoutes';

const queryClient = new QueryClient();

function App() {
  return (
    <ThemeProvider defaultTheme="dark" storageKey="ui-theme">
      <QueryClientProvider client={queryClient}>
        <WebsocketProvider>
          <BrowserRouter>
            <SidebarProvider defaultOpen={false}>
              <AppSidebar />
              <SidebarInset className="h-svh overflow-hidden flex flex-col">
                <TopBar />
                <Toaster position="bottom-right" />
                <div className="flex justify-center w-full h-full overflow-y-auto">
                  <div className="max-w-11/12 w-full">
                    <AppRoutes />
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
