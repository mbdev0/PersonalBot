package controller

type TaskController struct{}

func (tc *TaskController) GetTasks() []string {
	tasks := []string{"Task1", "Task2", "Task3"} // Example tasks
	return tasks
}

func (tc *TaskController) TestEP() string {
	return "Test successful"
}
