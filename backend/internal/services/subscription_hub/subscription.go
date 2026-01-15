package subscriptionhub

import "personal_bot/internal/core/tasks"

type Subscription struct {
	taskId        int64
	taskEventChan chan tasks.TaskEvent
	cancel        func()
}

func (s *Subscription) Id() int64 {
	return s.taskId
}

func (s *Subscription) Chan() chan tasks.TaskEvent {
	return s.taskEventChan
}

func (s *Subscription) Cancel() func() {
	return s.cancel
}
