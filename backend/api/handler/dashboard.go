package handler

import (
	"cmp"
	"encoding/json"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/api/dto"
	"slices"
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
	strategies, err := dh.strategyController.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	allTasks, err := dh.taskController.GetAllTasks()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dashboardResp := dh.generateDashboardResponse(strategies, allTasks)

	json.NewEncoder(w).Encode(dashboardResp)
}

func (dh *DashboardHandler) generateDashboardResponse(strategies []dto.TradingTaskResponse, allTasks []dto.ResponseTask) dto.DashboardResponseDto {
	dashboardResponse := dto.DashboardResponseDto{}
	dashboardResponse.New()

	for _, st := range strategies {

		childrenRows := []dto.ChildRow{}

		//TODO: terrible when we have 1000s of tasks -> maybe future improvement? - dont improve performance prematurely
		if st.Type == dto.AFK {
			for _, t := range allTasks {

				if !dh.shouldCreateRowFor(t, st) {
					continue
				}

				childRow := dto.ChildRow{
					Type:      string(dto.TASK),
					Id:        t.TaskId,
					WsMessage: t.Message,
					State:     t.State.TaskState,
					Data:      t,
				}

				childrenRows = append(childrenRows, childRow)
			}
		}

		slices.SortFunc(childrenRows, func(a, b dto.ChildRow) int {
			return cmp.Compare(a.Data.TimeCreated, b.Data.TimeCreated)
		})

		tbr := dto.TableRow{
			Type:      string(dto.STRATEGY),
			Id:        st.Id,
			WsMessage: st.Message,
			State:     st.State,
			Data:      st,
			Children:  childrenRows,
		}

		dashboardResponse.Rows = append(dashboardResponse.Rows, tbr)
	}

	slices.SortFunc(dashboardResponse.Rows, func(a, b dto.TableRow) int {
		return cmp.Compare(a.Data.TimeCreated, b.Data.TimeCreated)
	})

	return dashboardResponse
}

func (dh *DashboardHandler) shouldCreateRowFor(t dto.ResponseTask, st dto.TradingTaskResponse) bool {
	if t.Type == string(dto.Sell) {
		return false
	}

	if t.StrategyId == nil || *t.StrategyId != st.Id {
		return false
	}

	return true
}
