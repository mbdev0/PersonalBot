package logger

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type TaskLogger struct {
	messages chan string
	file     string
	done     chan struct{}
}

func logPath(taskId int64) string {
	exe, _ := os.Executable()
	return filepath.Join(filepath.Dir(exe), fmt.Sprintf("logs/task-%d-logs", taskId))
}

func NewTaskLogger(taskId int64) (*TaskLogger, error) {
	exe, _ := os.Executable()
	logsDir := filepath.Join(filepath.Dir(exe), "logs")
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return nil, err
	}

	file, err := os.Create(logPath(taskId))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	messagesChan := make(chan string, 100)
	tl := TaskLogger{
		messages: messagesChan,
		file:     file.Name(),
		done:     make(chan struct{}),
	}

	return &tl, nil
}

func (tl *TaskLogger) Start(ctx context.Context) error {
	logsFile, err := os.OpenFile(tl.file, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	go func() {
		defer logsFile.Close()
		defer close(tl.done)
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-tl.messages:
				fmt.Fprintln(logsFile, msg)
			}
		}
	}()

	return nil
}

type LogLevel string

const (
	SUCCESS LogLevel = "Success"
	ERROR   LogLevel = "Error"
	FAILURE LogLevel = "Failure"
	INFO    LogLevel = "Info"
)

func (tl *TaskLogger) write(level LogLevel, msg string) {
	message := fmt.Sprintf("%s | %s | %s", level, time.Now().Format(time.RFC3339), msg)
	tl.messages <- message
}

func (tl *TaskLogger) Success(args ...any) {
	msg := fmt.Sprint(args...)
	tl.write(SUCCESS, msg)
}

func (tl *TaskLogger) Error(args ...any) {
	msg := fmt.Sprint(args...)
	tl.write(ERROR, msg)
}

func (tl *TaskLogger) Failure(args ...any) {
	msg := fmt.Sprint(args...)
	tl.write(FAILURE, msg)

}

func (tl *TaskLogger) Information(args ...any) {
	msg := fmt.Sprint(args...)
	tl.write(INFO, msg)
}

func (tl *TaskLogger) TaskInformation(id int64, args ...any) {
	msg := fmt.Sprintf("[Task %d] %s", id, fmt.Sprint(args...))
	tl.write(INFO, msg)
}

func (tl *TaskLogger) TaskFailure(id int64, args ...any) {
	msg := fmt.Sprintf("[Task %d] %s", id, fmt.Sprint(args...))
	tl.write(FAILURE, msg)
}

func (tl *TaskLogger) TaskError(id int64, args ...any) {
	msg := fmt.Sprintf("[Task %d] %s", id, fmt.Sprint(args...))
	tl.write(ERROR, msg)
}

func (tl *TaskLogger) TaskSuccess(id int64, args ...any) {
	msg := fmt.Sprintf("[Task %d] %s", id, fmt.Sprint(args...))
	tl.write(SUCCESS, msg)
}

func (tl *TaskLogger) Done() <-chan struct{} {
	return tl.done
}
