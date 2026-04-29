package controller

import (
	"context"
	"fmt"
	"personal_bot/api/dto"
	"personal_bot/api/mapper"
	"personal_bot/internal/core/wallets"
	rpcgroups "personal_bot/internal/services/rpc_groups"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	taskservice "personal_bot/internal/services/task_service"
	"personal_bot/internal/services/wallet"
)

type TaskController struct {
	TaskService     *taskservice.TaskService
	WalletService   *wallet.Service
	RPCGroupService *rpcgroups.Service
}

func NewTaskController(ts *taskservice.TaskService, walletService *wallet.Service, rpcGroupService *rpcgroups.Service) *TaskController {
	return &TaskController{
		TaskService:     ts,
		WalletService:   walletService,
		RPCGroupService: rpcGroupService,
	}
}

func (tc *TaskController) CreateTask(ctx context.Context, requestTask dto.RequestTask) (*dto.ResponseTask, error) {
	// get the req struct
	// map to buy task or sell task depending

	wallet, err := tc.WalletService.GetByName(ctx, requestTask.WalletAddressName)
	if err != nil {
		return nil, err
	}

	_, err = tc.RPCGroupService.Load(ctx, requestTask.RPCGroupId)
	if err != nil {
		return nil, err
	}

	rpcGroup, err := tc.RPCGroupService.GetNode(requestTask.RPCGroupId)
	if err != nil {
		return nil, err
	}

	//we should check if the position exists if it occurs in the request task (selling)
	newTask, err := mapper.MapRequestTaskToTask(&requestTask, wallet, rpcGroup)
	if err != nil {
		return nil, err
	}

	createdTask, err := tc.TaskService.Create(newTask)
	if err != nil {
		return nil, err
	}

	response, err := mapper.MapTaskToReponseTask(createdTask)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (tc *TaskController) GetTask(id int64) (*dto.ResponseTask, error) {
	task, err := tc.TaskService.GetTaskWith(id)
	if err != nil {
		return nil, err
	}

	response, err := mapper.MapTaskToReponseTask(task)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (tc *TaskController) GetAllTasks() ([]dto.ResponseTask, error) {
	allTasks := tc.TaskService.GetAllTasks()

	response := make([]dto.ResponseTask, 0, len(allTasks))
	for _, task := range allTasks {
		responseObj, err := mapper.MapTaskToReponseTask(task)
		if err != nil {
			return nil, err
		}
		response = append(response, *responseObj)
	}

	return response, nil
}

func (tc *TaskController) UpdateTask(ctx context.Context, id int64, reqTask dto.PatchRequestTask) (*dto.ResponseTask, error) {
	// we call update with => id + newTask

	task, err := tc.TaskService.GetTaskWith(id)
	if err != nil {
		return nil, err
	}

	var wallet *wallets.SolanaWallet
	if reqTask.WalletAddressName != nil {
		walletResp, err := tc.WalletService.GetByName(ctx, *reqTask.WalletAddressName)
		if err != nil {
			return nil, err
		}
		wallet = &walletResp
	}

	mappedPatch, err := mapper.MapReqPatchToPatch(reqTask, task.Type(), wallet)
	if err != nil {
		return nil, err
	}

	updated, err := tc.TaskService.UpdateTask(task, mappedPatch)
	if err != nil {
		return nil, err
	}

	mappedUpdatedTask, err := mapper.MapTaskToReponseTask(updated)
	if err != nil {
		return nil, err
	}
	return mappedUpdatedTask, nil
}

func (tc *TaskController) DeleteTask(id int64) (err error) {
	//TODO: unload the rpc group on deletion
	err = tc.TaskService.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskController) TransitionTask(id int64, action dto.ActionType) (err error) {
	switch action {
	case dto.Run:
		if err = tc.TaskService.StartTask(id); err != nil {
			return err
		}
	case dto.Stop:
		if err = tc.TaskService.StopTask(id); err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid state passed in - use either Run or Stop")
	}
	return nil
}

func (tc *TaskController) CreateAndRunTask(ctx context.Context, task dto.RequestTask) error {
	createdTask, err := tc.CreateTask(ctx, task)
	if err != nil {
		return err
	}

	return tc.TransitionTask(createdTask.TaskId, dto.Run)
}

func (tc *TaskController) Subscribe(id int64) (*subscriptionhub.Subscription, error) {
	task, err := tc.TaskService.GetTaskWith(id)
	if err != nil {
		return nil, err
	}

	c, err := tc.TaskService.Subscribe(task)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (tc *TaskController) Unsubcribe(id int64) error {
	task, err := tc.TaskService.GetTaskWith(id)
	if err != nil {
		return err
	}

	err = tc.TaskService.Unsubscribe(task)
	if err != nil {
		return err
	}
	return nil
}

func (tc *TaskController) TestEP() string {
	return "Test successful"
}
