package handler

import (
	"encoding/json"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/api/dto"
)

type DashboardHandler struct {
	strategyController *controller.StrategyController
	taskController     *controller.TaskController
}

func NewDashboardHandler(sc *controller.StrategyController, tc *controller.TaskController) http.Handler {
	mux := http.NewServeMux()
	dashboardHandler := &DashboardHandler{strategyController: sc, taskController: tc}
	dashboardHandler.registerRoutes(mux)

	return mux
}

func (dh *DashboardHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", dh.GetDashboard)
}

func (dh *DashboardHandler) GetDashboard(w http.ResponseWriter, r *http.Request) {
	strategies, _ := dh.strategyController.GetAll()
	allTasks, _ := dh.taskController.GetAllTasks()

	tasksByStrategy := make(map[int64][]dto.ResponseTask)
	manualTasks := []dto.ResponseTask{}

	for _, task := range allTasks {
		if task.StrategyId != nil {
			tasksByStrategy[*task.StrategyId] = append(tasksByStrategy[*task.StrategyId], task)
		} else {
			manualTasks = append(manualTasks, task)
		}
	}

	response := dto.DashboardResponse{
		Strategies:      strategies,
		TasksByStrategy: tasksByStrategy,
		ManualTasks:     manualTasks,
	}

	json.NewEncoder(w).Encode(response)
}
