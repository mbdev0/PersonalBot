import './App.css';
import { TaskDashboard } from './features/tasks/components/taskDashboard';
import WalletDashboard from './features/wallets/components/walletDashboard';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <div className="min-h-screen bg-slate-500">
        <div className="max-w-7xl mx-auto px-4 py-8">
          <header className="mb-8">
            <h1 className="text-3xl font-bold text-slate-900">Trading Bot</h1>
            <p className="text-slate-600">Tasks</p>
          </header>
          <div className="bg-slate-900 rounded-xl shadow-lg border border-slate-950 text-slate-200">
            <TaskDashboard />
            {/* <WalletDashboard /> */}
          </div>
        </div>
      </div>
    </QueryClientProvider>
  );
}

export default App;
