// import { TaskType, type Task, type TaskDto } from "../types/task";

// function MapTaskDtoToTask(src: TaskDto): Task{
//     if (!isTaskType(src.state.task_state)) {
//         console.log("state is not a task state: ", src.state.task_state)
//         throw new Error("state is not a task state");
//     }

//     return {
//         task_id: src.task_id,
//         type: src.type,
//         slippage: src.slippage,
//         compute_units: src.compute_units,
//         token_address: src.token_address,
//         buy_amount: src.buy_amount,
//         buy_fee: src.buy_fee,
//         sell_amount: src.sell_amount,
//         sell_posiiton_id: src.sell_posiiton_id,
//         sell_fee: src.sell_fee,
//         state: {error: src.state.error, task_state: src.state.task_state}
//     }
// }

// function MapTaskToDto(src: Task): TaskDto {
//     return {
//         task_id: src.task_id,
//         type: src.type,
//         slippage: src.slippage,
//         compute_units: src.compute_units,
//         token_address: src.token_address,
//         buy_amount: src.buy_amount,
//         buy_fee: src.buy_fee,
//         sell_amount: src.sell_amount,
//         sell_posiiton_id: src.sell_posiiton_id,
//         sell_fee: src.sell_fee,
//         state: {error: src.state.error, task_state: src.state.task_state}
//     }
// }

// function isTaskType(value: string): value is TaskType{
//     return Object.values(TaskType).includes(value as TaskType)
// }
