import { API_BASE } from "../../../config/urls";
import { mapTaskDtoToTask, mapTaskToDto, mapTaskToPostDto } from "../mapper/taskMapper";
import type { Task, TaskDto, TaskPost } from "../types/task";

export async function getTasks(): Promise<Task[]> {
  const url = API_BASE + '/tasks/task'
  const response = await fetch(url)
  if (!response.ok){
    throw new Error('Error whilst getting tasks')
  }

  const taskDtos: TaskDto[] = await response.json()

  const mappedTasks = taskDtos.map(mapTaskDtoToTask)

  return mappedTasks
}

export async function postTask(task: TaskPost) {
  const mappedTasks = mapTaskToPostDto(task)
  const url = API_BASE + '/tasks/create'

  const response = await fetch(url, {
    method: 'POST',
    headers: {
      Accept: 'application/json',
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(mappedTasks)
  })

  if (!response.ok){
    throw new Error(`Failed to add new task: ${response.status}`)
  }
}