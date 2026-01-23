import './App.css';
import { TaskDashboard } from './features/tasks/components/taskDashboard';
import WalletDashboard from './features/wallets/components/walletDashboard';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-background">
        <div className="max-w-7xl mx-auto py-8">
          <header>
            <h1 className="text-3xl font-bold">Trading Bot</h1>
            <p>Tasks</p>
          </header>
          <div className="rounded-xl shadow-lg border bg-secondary p-5">
            <TaskDashboard />
            {/* <WalletDashboard /> */}
          </div>
        </div>
      </div>
    </QueryClientProvider>
  );
}

export default App;
