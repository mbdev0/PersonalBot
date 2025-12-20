import { useState } from 'react';
import { useTasks } from '../hooks/useTasks';
import './taskTable.css';

export function TaskTable() {
  const { isPending, isError, data, error } = useTasks();
  const [isRunning, setIsRunning] = useState(false);

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
            {/* could this be frontend only? - NO!*/}
            <th scope="column">Task Name</th>
            <th scope="column">Task Type</th>
            <th scope="column">Status</th>
            <th scope="column">Messages</th>
            <th scope="column">Actions</th>
          </tr>
        </thead>
        <tbody>
          {data?.map((task) => (
            <tr key={task.task_id}>
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
                  <button className="stop_task" onClick={() => setIsRunning(false)}>
                    ⏹️
                  </button>
                ) : (
                  <button className="start_task" onClick={() => setIsRunning(true)}>
                    ▶️
                  </button>
                )}
                <button className="edit">✏️</button>
                <button className="delete">🗑️</button>
                <button className="duplicate">⿻</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
