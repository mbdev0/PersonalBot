package subscriptionhub

import "pump_fun/internal/core/tasks"

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
