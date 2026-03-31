import { API_BASE } from '../../../config/urls';
import { mapTaskToPostDto } from '../mapper/taskMapper';
import type { TaskAction, TaskPost } from '../types/task';
const TASK_BASE = '/tasks/task';

export async function postTask(task: TaskPost): Promise<number> {
  const mappedTasks = mapTaskToPostDto(task);
  const url = API_BASE + TASK_BASE;

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(mappedTasks),
  });

  if (!response.ok) {
    throw new Error(`Failed to add new task: ${response.status}`);
  }

  return response.status;
}

export async function deleteTask(id: number): Promise<number> {
  const url = API_BASE + TASK_BASE + `/${id}`;

  const response = await fetch(url, {
    method: 'DELETE',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error(`Deletion failed of task: ${id}, response: ${response.status}`);
  }

  return response.status;
}

export async function transitionTask(id: number, taskState: TaskAction): Promise<number> {
  const url = API_BASE + `/tasks/transition/${id}`;
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ action: taskState }),
  });

  if (!response.ok) {
    throw new Error(
      `error whilst transitioning states for task ${id}, response: ${response.status}`
    );
  }

  return response.status;
}
