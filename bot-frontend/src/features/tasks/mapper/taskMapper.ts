import {
  type Task,
  type TaskDto,
  type TaskPost,
  type TaskPostDto,
  type TaskPut,
  type TaskPutDto,
} from '../types/task';

export function mapTaskDtoToTask(src: TaskDto): Task {
  return {
    task_id: src.task_id,
    type: src.type,
    slippage: src.slippage,
    compute_units: src.compute_units,
    token_address: src.token_address,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_amount: src.sell_amount,
    sell_posiiton_id: src.sell_posiiton_id,
    sell_fee: src.sell_fee,
    wallet_name: src.wallet_name,
    strategy_id: src.strategy_id,
    state: { error: src.state.error, task_state: src.state.task_state },
  };
}

export function mapTaskToDto(src: Task): TaskDto {
  return {
    task_id: src.task_id,
    type: src.type,
    slippage: src.slippage,
    compute_units: src.compute_units,
    token_address: src.token_address,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_amount: src.sell_amount,
    sell_posiiton_id: src.sell_posiiton_id,
    sell_fee: src.sell_fee,
    wallet_name: src.wallet_name,
    strategy_id: src.strategy_id,
    state: { error: src.state.error, task_state: src.state.task_state },
  };
}

export function mapTaskToPostDto(src: TaskPost): TaskPostDto {
  return {
    type: src.type,
    slippage: src.slippage,
    compute_units: src.compute_units,
    wallet_name: src.wallet_name,
    token_address: src.token_address,
    buy_amount: src.buy_amount,
    buy_fee: src.buy_fee,
    sell_amount: src.sell_amount,
    sell_fee: src.sell_fee,
  };
}

export function mapTaskToPutDto(src: TaskPut): TaskPutDto {
  const { id, ...rest } = src;
  return rest;
}
