package subscriptionhub

import (
	"fmt"
	"pump_fun/internal/core/tasks"
	"sync"
)

type Subscription struct {
	taskEventChan chan tasks.TaskEvent
	cancel        func()
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
	last map[string]tasks.TaskEvent
	mu   *sync.RWMutex
}

func (h *Hub) New() {
	h.subscriptions = map[string]Subscription{}
	h.last = map[string]tasks.TaskEvent{}
	h.mu = &sync.RWMutex{}
}

func (h *Hub) Subscribe(task tasks.Task) *Subscription {
	h.mu.Lock()
	defer h.mu.Unlock()
	sub, ok := h.subscriptions[task.Id()]
	if !ok {
		subChan := make(chan tasks.TaskEvent, 1000)
		cancel := h.cancel(task)
		h.subscriptions[task.Id()] = Subscription{subChan, cancel}
	}

	return &sub
}

func (h *Hub) Unsubcribe(task tasks.Task) error {
	h.mu.RLock()
	defer h.mu.RUnlock()
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
	sub, ok := h.subscriptions[task.Id()]
	if !ok {
		return
	}

	sub.taskEventChan <- event
	h.last[task.Id()] = event
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
