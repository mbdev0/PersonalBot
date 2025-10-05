package tasks

type Transistion struct {
	From    TaskState
	Next    TaskState
	OnError TaskState
}
