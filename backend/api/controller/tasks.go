package controller

import (
	"personal_bot/api/dto"
	"personal_bot/api/mapper"
	"personal_bot/internal/core/models/wallets"
	"personal_bot/internal/core/tasks"
	subscriptionhub "personal_bot/internal/services/subscription_hub"
	taskservice "personal_bot/internal/services/task_service"
	"personal_bot/internal/services/wallet"

	"golang.org/x/net/context"
)

type TaskController struct {
	TaskService   *taskservice.TaskService
	WalletService *wallet.Service
}

func (tc *TaskController) CreateTask(ctx context.Context, requestTask dto.RequestTask) (*dto.ResponseTask, error) {
	// get the req struct
	// map to buy task or sell task depending

	wallet, err := tc.WalletService.GetByName(ctx, requestTask.WalletAddressName)
	if err != nil {
		return nil, err
	}

	//we should check if the position exists if it occurs in the request task (selling)
	newTask, err := mapper.MapRequestTaskToTask(&requestTask, wallet)
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
	err = tc.TaskService.DeleteTask(id)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TaskController) TransitionTask(id int64, newState string) (err error) {
	state, err := tasks.ParseStateString(newState)
	if err != nil {
		return err
	}

	err = tc.TaskService.TransitionTask(id, state)
	if err != nil {
		return err
	}

	return nil
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
