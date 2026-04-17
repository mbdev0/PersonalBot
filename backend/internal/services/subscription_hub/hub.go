package subscriptionhub

import (
	"fmt"
	"personal_bot/internal/core/tasks"
	"personal_bot/pkg/logger"
	"sync"
	"time"
)

type Hub struct {
	//id -> {subscription chan, cancel func}
	subscriptions map[int64]Subscription
	// id -> {task event}
	last       map[int64]tasks.TaskEvent
	mu         *sync.RWMutex
	bufferSize int
}

func NewTaskSubscriptionHub() *Hub {
	return &Hub{
		subscriptions: map[int64]Subscription{},
		last:          map[int64]tasks.TaskEvent{},
		mu:            &sync.RWMutex{},
		bufferSize:    1000,
	}
}

func (h *Hub) Subscribe(task tasks.Task) (*Subscription, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// if there is an existing subscription, return an error -> we don't want multiple subs per task
	if _, ok := h.subscriptions[task.Id()]; ok {
		return nil, fmt.Errorf("an existing subscription is attatched to task: %d", task.Id())
	}

	subChan := make(chan tasks.TaskEvent, h.bufferSize)
	cancel := h.cancel(task)
	h.subscriptions[task.Id()] = Subscription{task.Id(), subChan, cancel}
	sub, ok := h.subscriptions[task.Id()]

	if last, ok := h.last[task.Id()]; ok {
		select {
		case sub.taskEventChan <- last:
		default:
			logger.Error("error whilst adding the last take event to the subscription")
		}
	}

	if !ok {
		return nil, fmt.Errorf("error whilst making the subscription for task id: %d", task.Id())
	}

	return &sub, nil
}

func (h *Hub) Unsubcribe(task tasks.Task) error {

	h.mu.Lock()
	defer h.mu.Unlock()
	sub, ok := h.subscriptions[task.Id()]
	if !ok {
		return fmt.Errorf("task not subscribed too")
	}

	sub.cancel()
	return nil
}

func (h *Hub) Last(task tasks.Task) (*tasks.TaskEvent, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	last, ok := h.last[task.Id()]
	if !ok {
		return nil, fmt.Errorf("task hasn't published anything yet")
	}

	return &last, nil
}

func (h *Hub) publish(task tasks.Task, event tasks.TaskEvent) {
	h.mu.Lock()
	defer h.mu.Unlock()

	//always keep a track of the last state
	h.last[task.Id()] = event

	sub, ok := h.subscriptions[task.Id()]
	if !ok {
		return
	}

	select {
	case sub.taskEventChan <- event:
	default:
		logger.InfoMessage("cannot pass event onto subscription channel")
	}

}

func (h *Hub) cancel(t tasks.Task) func() {
	return func() {
		sub := h.subscriptions[t.Id()]
		close(sub.taskEventChan)

		delete(h.subscriptions, t.Id())
		delete(h.last, t.Id())
	}
}

func (h *Hub) PublishStateChange(task tasks.Task) {

	wsState := tasks.StateDto{
		TaskState: task.State().TaskState.ToString(),
		Error:     task.State().Error,
	}

	taskEvent := tasks.TaskEvent{
		TaskId:    task.Id(),
		State:     wsState,
		Time:      time.Now().Local().String(),
		EventType: tasks.StateUpdate,
	}

	if task.GetStrategyId() != nil {
		taskEvent.StrategyId = *task.GetStrategyId()
	}

	h.publish(task, taskEvent)
}

func (h *Hub) PublishMessage(task tasks.Task, message string) {
	wsState := tasks.StateDto{
		TaskState: task.State().TaskState.ToString(),
		Error:     task.State().Error,
	}

	taskEvent := tasks.TaskEvent{
		TaskId:    task.Id(),
		State:     wsState,
		Time:      time.Now().Local().String(),
		Message:   message,
		EventType: tasks.ProgressMessage,
	}

	if task.GetStrategyId() != nil {
		taskEvent.StrategyId = *task.GetStrategyId()
	}

	task.SetMessage(message)
	h.publish(task, taskEvent)
}
