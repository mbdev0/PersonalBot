import { useState } from 'react';
import type { Task } from '../types/task';

interface TaskRowProps {
  task: Task;
  onStart: () => void;
  onStop: () => void;
  onEdit: () => void;
  onDelete: () => void;
  onDuplicate: () => void;
}

export function TaskRow({ task, onStart, onStop, onEdit, onDelete, onDuplicate }: TaskRowProps) {
  const [isRunning, setIsRunning] = useState(false);
  return (
    <tr>
      <td>{task.task_id}</td>
      <td>
        {`Task - ${task.task_id}`}
        <button className="edit_task_name">✏️</button>
      </td>
      <td>{task.type}</td>
      <td>{task.state.task_state}</td>
      <td>...{/*TODO: REPLACE THIS WITH MESSAGE component AFTER W.S*/}</td>
      <td>
        {/*TODO: make this a actions component */}
        {isRunning ? (
          <button
            className="stop_task"
            onClick={() => {
              setIsRunning(false);
              onStop();
            }}
          >
            ⏹️
          </button>
        ) : (
          <button
            className="start_task"
            onClick={() => {
              setIsRunning(true);
              onStart();
            }}
          >
            ▶️
          </button>
        )}

        <button
          className="edit"
          onClick={() => {
            onEdit();
          }}
        >
          ✏️
        </button>
        <button className="delete" onClick={() => onDelete()}>
          🗑️
        </button>
        <button className="duplicate" onClick={() => onDuplicate()}>
          ⿻
        </button>
      </td>
    </tr>
  );
}
