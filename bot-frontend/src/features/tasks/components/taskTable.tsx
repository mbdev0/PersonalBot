import { useState } from 'react';
import { useAddTask, useDeleteTask, useTasks, useTransitionTask } from '../hooks/useTasks';
import './taskTable.css';
import Modal from '../../../components/modal';
import { TaskUpdate } from './taskUpdate';
import { TaskType, type Task } from '../types/task';
import { TaskRow } from './taskRow';

export function TaskTable() {
  const { isPending, isError, data, error } = useTasks();
  const [editingTask, setEditingTask] = useState<Task | null>(null);
  const deleteMutation = useDeleteTask();
  const duplicateMutation = useAddTask();
  const transitionMutation = useTransitionTask();

  if (isPending) {
    return <div className="loading_tasks">Loading Tasks...</div>;
  }

  if (isError) {
    return <div className="loading_task_error">Error whilst loading tasks: {error.message}</div>;
  }

  return (
    <div className="task_table">
      <table>
        <thead>
          <tr>
            <th scope="column">Task ID</th>
            <th scope="column">Task Name</th>
            <th scope="column">Task Type</th>
            <th scope="column">Status</th>
            <th scope="column">Messages</th>
            <th scope="column">Actions</th>
          </tr>
        </thead>
        <tbody>
          {data
            ?.sort(function (a, b) {
              return a.task_id - b.task_id;
            })
            .map((task) => (
              <TaskRow
                key={task.task_id}
                task={task}
                onEdit={() => setEditingTask(task)}
                onStart={() =>
                  transitionMutation.mutate({ id: task.task_id, taskAction: TaskType.task_run })
                }
                onStop={() =>
                  transitionMutation.mutate({ id: task.task_id, taskAction: TaskType.task_cancel })
                }
                onDelete={() => deleteMutation.mutate(task.task_id)}
                onDuplicate={() => duplicateMutation.mutate(task)}
              ></TaskRow>
            ))}
        </tbody>
      </table>

      {editingTask && (
        <Modal isOpen={!!editingTask} onClose={() => setEditingTask(null)}>
          <TaskUpdate task={editingTask} onClose={() => setEditingTask(null)}></TaskUpdate>
        </Modal>
      )}
    </div>
  );
}
