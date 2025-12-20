import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { getTasks, postTask } from "../api/tasks";
import type { Task, TaskPost } from "../types/task";

export function useTasks(){
  return useQuery(
    {
      queryKey:['tasks'],
      queryFn: getTasks
    }
  )
}

export function postTasksMutation(){
  const client = useQueryClient();

  return useMutation({
    mutationFn: (task: TaskPost) => postTask(task) ,
    onSuccess() {
      client.invalidateQueries({queryKey: ['tasks']})
    }
  })
}
