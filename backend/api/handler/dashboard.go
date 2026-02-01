package handler

import (
	"encoding/json"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/api/dto"
)

// TODO: CLEAN UP
type RowType string

const (
	STRATEGY RowType = "strategy"
	TASK     RowType = "task"
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

	dashboardResponse := dto.DashboardResponseDto{}
	for _, st := range strategies {

		childrenRows := []dto.ChildRow{}

		//TODO: terrible when we have 1000s of tasks -> maybe future improvement? - dont improve performance prematurely
		if st.Type == dto.AFK {
			for _, t := range allTasks {
				if t.Type == string(dto.Sell) {
					continue
				}

				if t.StrategyId == nil || *t.StrategyId != st.Id {
					continue
				}

				childRow := dto.ChildRow{
					Type:      string(TASK),
					Id:        t.TaskId,
					WsMessage: "",
					State:     t.State.TaskState,
					Data:      t,
				}

				childrenRows = append(childrenRows, childRow)
			}
		}

		tbr := dto.TableRow{
			Type:      string(STRATEGY),
			Id:        st.Id,
			WsMessage: "",
			State:     st.State,
			Data:      st,
			Children:  childrenRows,
		}

		dashboardResponse.Rows = append(dashboardResponse.Rows, tbr)
	}

	json.NewEncoder(w).Encode(dashboardResponse)
}
