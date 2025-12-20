import { useState } from 'react';
import { TaskTable } from './taskTable';
import Modal from '../../../components/modal';
import { TaskEntry } from './taskEntry';

export function TaskDashboard() {
  const [isModalShowing, setModalShowing] = useState(false);
  return (
    <div className="task_dashboard">
      <h3>Tasks</h3>
      <TaskTable />
      <Modal isOpen={isModalShowing} onClose={() => setModalShowing(false)}>
        <TaskEntry onClose={() => setModalShowing(false)}></TaskEntry>
      </Modal>
      <button className="add_task" onClick={() => setModalShowing(true)}>
        Add Task
      </button>
    </div>
  );
}
