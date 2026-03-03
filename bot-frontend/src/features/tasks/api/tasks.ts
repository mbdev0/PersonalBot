import { API_BASE } from '../../../config/urls';
import { mapTaskDtoToTask, mapTaskToPostDto, mapTaskToPutDto } from '../mapper/taskMapper';
import type { Task, TaskAction, TaskDto, TaskPost, TaskPut } from '../types/task';
const TASK_BASE = '/tasks/task';

export async function getTasks(): Promise<Task[]> {
  const url = API_BASE + '/tasks/task';
  const response = await fetch(url);
  if (!response.ok) {
    throw new Error('Error whilst getting tasks');
  }

  const taskDtos: TaskDto[] = await response.json();

  const mappedTasks = taskDtos.map(mapTaskDtoToTask);

  return mappedTasks;
}

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

export async function putTask(task: TaskPut): Promise<number> {
  const mappedTasks = mapTaskToPutDto(task);
  const url = API_BASE + TASK_BASE + `/${task.id}`;

  const response = await fetch(url, {
    method: 'PUT',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(mappedTasks),
  });

  if (!response.ok) {
    throw new Error(`Updating task ${task.id} failed: ${response.status}`);
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
