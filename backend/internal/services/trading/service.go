package trading

import (
	"context"
	"fmt"
	"maps"
	"math"
	"personal_bot/app/iterable"
	"personal_bot/infrastructure/persistence/repository"
	rpcgroupsModel "personal_bot/internal/core/rpc_groups"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/tasks"
	rpcgroups "personal_bot/internal/services/rpc_groups"

	"personal_bot/internal/services/subscription_hub/strategy"
	taskservice "personal_bot/internal/services/task_service"
	tradingStrategy "personal_bot/internal/services/trading/strategy"
	"personal_bot/pkg/logger"
	"slices"
	"sync"
	"time"
)

type Service struct {
	strategy        *tradingStrategy.Strategy
	tasks           map[int64]strategies.Task
	running         map[int64]context.CancelFunc
	mu              *sync.Mutex
	subhub          *strategy.SubscriptionHub
	repo            *repository.TradingRepository
	iterable        *iterable.Iterable
	taskService     *taskservice.TaskService
	rpcGroupService *rpcgroups.Service
}

func NewTradingService(strategy *tradingStrategy.Strategy, sh *strategy.SubscriptionHub, tr *repository.TradingRepository, ts *taskservice.TaskService, rgs *rpcgroups.Service) *Service {

	return &Service{
		strategy:        strategy,
		tasks:           map[int64]strategies.Task{},
		running:         map[int64]context.CancelFunc{},
		mu:              &sync.Mutex{},
		subhub:          sh,
		repo:            tr,
		iterable:        iterable.NewIterable(),
		taskService:     ts,
		rpcGroupService: rgs,
	}
}

func (s *Service) Create(ctx context.Context, st strategies.Task) (task strategies.Task, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if st.StrategyTaskId() == 0 {
		st.SetId(s.iterable.ID())
	}

	if _, ok := s.tasks[st.StrategyTaskId()]; ok {
		return nil, fmt.Errorf("task already exists with id: %d", st.StrategyTaskId())
	}

	_, err = s.rpcGroupService.Load(ctx, st.RPCGroupId())
	if err != nil {
		return nil, err
	}

	//if we get a buy/sell strategy we precreate the task to speed up running
	if st.StrategyType() == strategies.BUY {
		err := s.initBuyStrategy(st)
		if err != nil {
			return nil, err
		}
	} else if st.StrategyType() == strategies.SELL {
		err := s.initSellStrategy(st)
		if err != nil {
			return nil, err
		}
	}

	s.tasks[st.StrategyTaskId()] = st

	return st, nil
}

func (s *Service) initBuyStrategy(st strategies.Task) (err error) {
	buyStrategy, ok := st.(*strategies.Buy)
	if !ok {
		return fmt.Errorf("expected strategy task to be a Buy Strategy - ended up getting a different type")
	}

	node, err := s.rpcGroupService.GetNode(st.RPCGroupId())
	if err != nil {
		return err
	}

	bt := s.createBuyTask(*buyStrategy, node)
	createdTask, err := s.taskService.Create(bt)
	if err != nil {
		return err
	}

	buyStrategy.BuyTaskId = createdTask.Id()
	buyStrategy.PositionId = createdTask.Id()
	return nil
}

func (s *Service) createBuyTask(buyTask strategies.Buy, node rpcgroupsModel.GroupItem) *tasks.BuyTask {
	bt := tasks.NewBuyTask(buyTask.Wallet, buyTask.Token,
		[]tasks.Option{
			tasks.WithProgram(buyTask.GetProgram()),
			tasks.WithSlippage(buyTask.Slippage),
			tasks.WithComputeUnits(uint32(buyTask.ComputeUnits)),
			tasks.WithStrategyId(buyTask.StrategyTaskId()),
			tasks.WithRPCGroupId(buyTask.RPCGroup.Id),
			tasks.WithHttpNode(node.Http),
			tasks.WithWS(node.WS),
		},
		[]tasks.BuyOption{
			tasks.WithBuyAmount(buyTask.BuyAmount),
			tasks.WithBuyFee(buyTask.BuyFee)},
	)

	return bt
}

func (s *Service) initSellStrategy(st strategies.Task) (err error) {
	sellStrategy, ok := st.(*strategies.Sell)
	if !ok {
		return fmt.Errorf("expected strategy task to be a Sell Strategy - ended up getting a different type")
	}

	node, err := s.rpcGroupService.GetNode(sellStrategy.RPCGroupId())
	if err != nil {
		return err
	}

	sellTask := s.createSellTask(*sellStrategy, node)
	createdTask, err := s.taskService.Create(sellTask)
	if err != nil {
		return err
	}

	sellStrategy.SellTaskId = createdTask.Id()
	return nil
}

func (s *Service) createSellTask(sell strategies.Sell, rpcGroup rpcgroupsModel.GroupItem) *tasks.SellTask {
	st := tasks.NewSellTask(
		sell.GetWallet(),
		sell.Token,
		[]tasks.Option{
			tasks.WithProgram(sell.GetProgram()),
			tasks.WithComputeUnits(uint32(sell.GetComputeUnits())),
			tasks.WithSlippage(sell.GetSlippage()),
			tasks.WithStrategyId(sell.StrategyTaskId()),
			tasks.WithRPCGroupId(sell.RPCGroupId()),
			tasks.WithHttpNode(rpcGroup.Http),
			tasks.WithWS(rpcGroup.WS),
		},
		[]tasks.SellOption{
			tasks.WithSellAmount(sell.SellAmount),
			tasks.WithSellFee(sell.SellFee),
		},
	)

	return st
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.tasks[id]; !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	taskCancel, ok := s.running[id]
	if ok {
		taskCancel()
	}

	task := s.tasks[id]
	delete(s.tasks, id)

	err := s.rpcGroupService.Unload(task.RPCGroupId())
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	success, err := s.repo.Delete(ctx, id)
	if err != nil {
		s.tasks[id] = task
		return fmt.Errorf("failed to delete from database: %w", err)
	}

	if !success {
		s.tasks[id] = task
		return fmt.Errorf(" unsuccessful delete from database")
	}

	if task.StrategyType() == strategies.SELL {
		sell, ok := task.(*strategies.Sell)
		if !ok {
			return fmt.Errorf("unable to cast task id %d, to type sell strategy", task.StrategyTaskId())
		}

		err := s.taskService.DeleteTask(sell.SellTaskId)
		if err != nil {
			return err
		}

		return nil
	}

	tasks := s.taskService.GetTasksWithStrategyId(id)

	for _, task := range tasks {
		err = s.taskService.DeleteTask(task.Id())
		if err != nil {
			return err
		}
	}

	return nil

}

func (s *Service) GetBy(id int64) (strategies.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("task not found with id: %d", id)
	}

	return task, nil
}

func (s *Service) GetAll() []strategies.Task {
	s.mu.Lock()
	defer s.mu.Unlock()

	allTasks := make([]strategies.Task, 0, len(s.tasks))
	for _, val := range s.tasks {
		allTasks = append(allTasks, val)
	}

	return allTasks
}

func (s *Service) Update(id int64, task strategies.Task) (strategies.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	oldTask, ok := s.tasks[id]
	if !ok {
		return nil, fmt.Errorf("unable to find task: %d", id)
	}

	s.tasks[id] = task
	task.SetId(id)
	task.SetTimeCreated(oldTask.GetTimeCreated())

	var err error
	//any tasks that are precreated, we must update
	if task.StrategyType() == strategies.BUY {
		err = s.updateBuyTask(task)
	} else if task.StrategyType() == strategies.SELL {
		err = s.updateSellTask(task)
	} else if task.StrategyType() == strategies.SPAM {
		err = s.updateSpamTask(task)
	}

	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *Service) updateBuyTask(task strategies.Task) error {
	allTasks := s.taskService.GetTasksWithStrategyId(task.StrategyTaskId())
	if len(allTasks) == 0 || len(allTasks) > 1 {
		return nil
	}

	buyPatchPtr, ok := task.(*strategies.Buy)
	if !ok {
		return fmt.Errorf("invalid casting to buy patch, check if this is being called in the right place")
	}
	buyTask := *buyPatchPtr
	cu := uint32(math.Round(buyTask.ComputeUnits))

	taskPatch := tasks.BuyPatch{
		Program:        &buyTask.Program,
		Wallet:         &buyTask.Wallet,
		Token:          &buyTask.Token,
		Amount:         buyTask.BuyAmount,
		Fee:            &buyTask.BuyFee,
		Slippage:       &buyTask.Slippage,
		ComputeUnit:    &cu,
		Retries:        buyTask.Retries,
		RetriesDelayMs: buyTask.RetriesDelayMS,
	}

	_, err := s.taskService.UpdateTask(allTasks[0], &taskPatch)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) updateSellTask(task strategies.Task) error {
	allTasks := s.taskService.GetTasksWithStrategyId(task.StrategyTaskId())

	if len(allTasks) == 0 || len(allTasks) > 1 {
		return nil
	}

	sellPatchPtr, ok := task.(*strategies.Sell)
	if !ok {
		return fmt.Errorf("invalid casting to sell patch, check if this is being called in the right place")
	}

	sellTask := *sellPatchPtr

	computeUnits := uint32(math.Round(sellTask.ComputeUnits))

	taskPatch := tasks.SellPatch{
		Program:        &sellTask.Program,
		Wallet:         &sellTask.Wallet,
		Token:          &sellTask.Token,
		Amount:         &sellTask.SellAmount,
		Fee:            &sellTask.SellFee,
		Slippage:       &sellTask.Slippage,
		ComputeUnit:    &computeUnits,
		Retries:        sellTask.Retries,
		RetriesDelayMs: sellTask.RetriesDelayMS,
	}

	_, err := s.taskService.UpdateTask(allTasks[0], &taskPatch)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) updateSpamTask(task strategies.Task) error {
	spamTask, ok := task.(*strategies.Spam)
	if !ok {
		return fmt.Errorf("wrong task passed in, check updating logic")
	}

	allSubTasks := s.taskService.GetTasksWithStrategyId(spamTask.StrategyTaskId())
	cu := uint32(math.Round(spamTask.ComputeUnits))
	taskPatch := tasks.BuyPatch{
		Program:        &spamTask.Program,
		Wallet:         &spamTask.Wallet,
		Token:          &spamTask.Token,
		Amount:         spamTask.BuyAmount,
		Fee:            &spamTask.BuyFee,
		Slippage:       &spamTask.Slippage,
		ComputeUnit:    &cu,
		Retries:        spamTask.Retries,
		RetriesDelayMs: spamTask.RetriesDelayMS,
	}

	for _, t := range allSubTasks {
		_, err := s.taskService.UpdateTask(t, &taskPatch)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Start(id int64) error {
	task, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	_, isRunning := s.running[id]
	if isRunning {
		return fmt.Errorf("task is already running %d", id)
	}

	ctxCancel, cancel := context.WithCancel(context.Background())
	err := s.strategy.Run(ctxCancel, task)
	if err != nil {
		cancel()
		return err
	}

	s.running[task.StrategyTaskId()] = cancel

	return nil
}

func (s *Service) Stop(id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// this will call a tasks cancel
	strategyTask, ok := s.tasks[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}

	taskCancel, ok := s.running[id]
	if !ok {
		return fmt.Errorf("task not running with id: %d", id)
	}

	childTasks := s.taskService.GetTasksWithStrategyId(strategyTask.StrategyTaskId())

	for _, child := range childTasks {
		err := s.taskService.StopTask(child.Id())
		if err != nil {
			logger.Error(err)
		}
	}

	taskCancel()
	delete(s.running, strategyTask.StrategyTaskId())
	return nil
}

func (s *Service) Subscribe(taskId int64) (*strategy.Subscription, error) {
	sub, err := s.subhub.Subscribe(taskId)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *Service) Unsubscribe(taskId int64) {
	s.subhub.Unsubscribe(taskId)
}

// loads tasks into memory
func (s *Service) LoadFromDB(ctx context.Context) error {
	ttFromDb, err := s.repo.GetAllTasks(ctx)
	if err != nil {
		return err
	}

	maxId := s.repo.GetMaxId(ctx)
	s.iterable.SetIterable(maxId)

	for _, tt := range ttFromDb {
		s.tasks[tt.StrategyTaskId()] = tt
		_, err = s.rpcGroupService.Load(ctx, tt.RPCGroupId())
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}

// pushes all in-memory tasks to db
func (s *Service) Shutdown(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	deleteSuccess, err := s.repo.DeleteAll(ctx)
	if err != nil {
		return err
	}

	if !deleteSuccess {
		return fmt.Errorf("failed to wipe table for insertion")
	}

	tasksToSave := slices.Collect(maps.Values(s.tasks))

	success, err := s.repo.AddAllTasks(ctx, tasksToSave)
	if err != nil {
		return fmt.Errorf("error whilst shutting down: %w", err)
	}

	if !success {
		return fmt.Errorf("unsuccesful graceful shutdown in tasks service")
	}

	return nil
}
