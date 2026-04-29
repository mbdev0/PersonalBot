import { type Task, type TaskDto, type TaskPost, type TaskPostDto } from '../types/task';

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
    message: src.message,
    rpc_group_id: src.rpc_group_id,
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
    rpc_group_id: src.rpc_group_id,
    strategy_id: src.strategy_id,
  };
}

export function mapTaskToTaskPost(src: Task): TaskPostDto {
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
    rpc_group_id: src.rpc_group_id,
    strategy_id: src.strategy_id ?? 0,
  };
}
