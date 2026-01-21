import { useState } from 'react';
import { TaskTable } from './taskTable';
import Modal from '../../../components/modal';
import { TaskEntry } from './taskEntry';

export function TaskDashboard() {
  const [isModalShowing, setModalShowing] = useState(false);
  return (
    <div className="task_dashboard">
      <div className="flex justify-end px-4 p-1">
        <button
          className="add_task border-slate-600 bg-slate-500 p-2 rounded-xl "
          onClick={() => setModalShowing(true)}
        >
          Add Task
        </button>
      </div>

      <TaskTable />
      <Modal isOpen={isModalShowing} onClose={() => setModalShowing(false)}>
        <TaskEntry onClose={() => setModalShowing(false)}></TaskEntry>
      </Modal>
    </div>
  );
}
