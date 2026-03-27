package strategy

import (
	"fmt"
	"personal_bot/internal/core/strategies"
	"personal_bot/internal/core/tasks"
	"personal_bot/pkg/logger"
	"sync"
)

type Subscription struct {
	Sub_id  int64
	SubChan chan strategies.StrategyMessage
	cancel  func()
}

type SubscriptionHub struct {
	subscriptions map[int64]*Subscription
	bufferSize    int
	mu            *sync.Mutex
	last          map[int64]strategies.StrategyMessage
}

func NewSubscriptionHub() *SubscriptionHub {
	return &SubscriptionHub{
		subscriptions: map[int64]*Subscription{},
		bufferSize:    1000,
		mu:            &sync.Mutex{},
		last:          map[int64]strategies.StrategyMessage{},
	}
}

func (sh *SubscriptionHub) Subscribe(taskId int64) (*Subscription, error) {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	if _, ok := sh.subscriptions[taskId]; ok {
		return nil, fmt.Errorf("an existing subscription is attatched to task: %d", taskId)
	}

	subChan := make(chan strategies.StrategyMessage, sh.bufferSize)
	if _, ok := sh.last[taskId]; ok {
		subChan <- sh.last[taskId]
	}

	cancel := sh.cancel(taskId)

	sub := &Subscription{Sub_id: taskId, SubChan: subChan, cancel: cancel}

	sh.subscriptions[taskId] = sub
	return sub, nil

}

func (sh *SubscriptionHub) Unsubscribe(taskId int64) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	if _, ok := sh.subscriptions[taskId]; !ok {
		return fmt.Errorf("task not found for deletion: %d", taskId)
	}

	sh.subscriptions[taskId].cancel()
	return nil
}

func (sh *SubscriptionHub) cancel(id int64) func() {
	return func() {
		sub := sh.subscriptions[id]
		close(sub.SubChan)
		delete(sh.subscriptions, id)
		delete(sh.last, id)
	}
}

func (sh *SubscriptionHub) PublishTakeCreation(id int64, task tasks.Task) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	msg := strategies.StrategyMessage{
		Id:    id,
		Event: strategies.TaskCreation,
		Task:  task,
	}

	logger.Information("publishing task creation")

	sh.last[id] = msg

	sub, ok := sh.subscriptions[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}
	sub.SubChan <- msg
	return nil
}

func (sh *SubscriptionHub) PublishStateUpdate(id int64, state string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	msg := strategies.StrategyMessage{
		Id:    id,
		Event: strategies.StatusUpdate,
		State: &state,
	}

	sh.last[id] = msg

	sub, ok := sh.subscriptions[id]
	if !ok {
		//if there isnt a subscription, we can just set the last message - we dont need a sub to actually publish a message
		// we can just keep the last message safe
		return nil
	}
	sub.SubChan <- msg
	return nil
}

func (sh *SubscriptionHub) PublishProgressMessage(id int64, message string) error {
	sh.mu.Lock()
	defer sh.mu.Unlock()

	msg := strategies.StrategyMessage{
		Id:      id,
		Event:   strategies.MessageUpdate,
		Message: &message,
	}

	sh.last[id] = msg

	sub, ok := sh.subscriptions[id]
	if !ok {
		return fmt.Errorf("task not found with id: %d", id)
	}
	sub.SubChan <- msg
	return nil
}
