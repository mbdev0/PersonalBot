package subscriptionhub

import (
	"fmt"
	"pump_fun/internal/core/tasks"
	"pump_fun/pkg/logger"
	"sync"
)

type Subscription struct {
	taskId        string
	taskEventChan chan tasks.TaskEvent
	cancel        func()
}

func (s *Subscription) Id() string {
	return s.taskId
}

func (s *Subscription) Chan() chan tasks.TaskEvent {
	return s.taskEventChan
}

func (s *Subscription) Cancel() func() {
	return s.cancel
}

type Hub struct {
	//id -> {subscription chan, cancel func}
	subscriptions map[string]Subscription
	// id -> {task event}
	last       map[string]tasks.TaskEvent
	mu         *sync.RWMutex
	bufferSize int
}

func (h *Hub) New() {
	h.subscriptions = map[string]Subscription{}
	h.last = map[string]tasks.TaskEvent{}
	h.mu = &sync.RWMutex{}
	h.bufferSize = 1000
}

func (h *Hub) Subscribe(task tasks.Task) (*Subscription, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// if there is an existing subscription, return an error -> we don't want multiple subs per task
	if _, ok := h.subscriptions[task.Id()]; ok {
		return nil, fmt.Errorf("an existing subscription is attatched to task: %s", task.Id())
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
		return nil, fmt.Errorf("error whilst making the subscription for task id: %s", task.Id())
	}

	logger.Information(h.subscriptions)
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

func (h *Hub) Publish(task tasks.Task, event tasks.TaskEvent) {
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
		h.mu.Lock()
		defer h.mu.Unlock()

		sub := h.subscriptions[t.Id()]
		close(sub.taskEventChan)

		delete(h.subscriptions, t.Id())
		delete(h.last, t.Id())
	}
}
