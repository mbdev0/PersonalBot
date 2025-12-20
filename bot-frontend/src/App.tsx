import './App.css';
import { TaskDashboard } from './features/tasks/components/taskDashboard';
import WalletDashboard from './features/wallets/components/walletDashboard';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <h1>Trading Bot</h1>
      {/* <CreateTask />/ */}
      {/* <WalletDashboard /> */}
      {/* <TaskTable /> */}
      {<TaskDashboard />}
    </QueryClientProvider>
  );
}

export default App;
