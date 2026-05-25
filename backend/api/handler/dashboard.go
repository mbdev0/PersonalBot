package handler

import (
	"cmp"
	"encoding/json"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/api/dto"
	"personal_bot/pkg/logger"
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

	err = json.NewEncoder(w).Encode(dashboardResp)
	if err != nil {
		logger.Error(err)
	}
}

func (dh *DashboardHandler) generateDashboardResponse(strategies []dto.TradingTaskResponse, allTasks []dto.ResponseTask) dto.DashboardResponseDto {
	dashboardResponse := dto.DashboardResponseDto{}
	dashboardResponse.New()
	strategyToBuy := map[int64][]dto.ResponseTask{}
	sellPositionToSellTask := map[int64][]dto.ResponseTask{}

	for _, task := range allTasks {
		if task.Type == "Buy" && task.StrategyId != nil {
			strategyToBuy[*task.StrategyId] = append(strategyToBuy[*task.StrategyId], task)
		} else if task.Type == "Sell" && task.SellPositionId != nil {
			sellPositionToSellTask[*task.SellPositionId] = append(sellPositionToSellTask[*task.SellPositionId], task)
		}
	}

	for _, st := range strategies {
		switch st.Type {
		case dto.BUY:
			var sellTasks []dto.ResponseTask
			if st.BuyTaskId != nil {
				sellTasks = sellPositionToSellTask[*st.BuyTaskId]
			}
			dashboardResponse.Rows = append(dashboardResponse.Rows, dh.createBuyStrategyRow(st, sellTasks))
		case dto.SELL:
			dashboardResponse.Rows = append(dashboardResponse.Rows, dh.createSellStrategyRow(st))
		default:
			// AFK strategy — buy tasks as children, each with sell task grandchildren
			buyTasks := strategyToBuy[st.Id]
			children := make([]dto.ChildRow, 0, len(buyTasks))

			for _, buyTask := range buyTasks {
				buyTask.SellFee = st.SellFee

				sellTasks := sellPositionToSellTask[buyTask.TaskId]
				grandChildren := make([]dto.SellTasks, 0, len(sellTasks))

				for _, sellTask := range sellTasks {
					grandChildren = append(grandChildren, dto.SellTasks{
						Program:   string(sellTask.Program),
						Type:      string(dto.TASK),
						Id:        sellTask.TaskId,
						WsMessage: sellTask.Message,
						State:     sellTask.State.TaskState,
						Data:      sellTask,
					})
				}

				slices.SortFunc(grandChildren, func(a, b dto.SellTasks) int {
					return cmp.Compare(a.Data.TimeCreated, b.Data.TimeCreated)
				})

				children = append(children, dto.ChildRow{
					Program:   string(buyTask.Program),
					Type:      string(dto.TASK),
					Id:        buyTask.TaskId,
					WsMessage: buyTask.Message,
					State:     buyTask.State.TaskState,
					Data:      buyTask,
					Children:  grandChildren,
				})
			}

			slices.SortFunc(children, func(a, b dto.ChildRow) int {
				return cmp.Compare(a.Data.TimeCreated, b.Data.TimeCreated)
			})

			dashboardResponse.Rows = append(dashboardResponse.Rows, dto.TableRow{
				Program:   st.Program,
				Type:      string(dto.STRATEGY),
				Id:        st.Id,
				WsMessage: st.Message,
				State:     st.State,
				Data:      st,
				Children:  children,
			})
		}
	}

	slices.SortFunc(dashboardResponse.Rows, func(a, b dto.TableRow) int {
		return cmp.Compare(a.Data.TimeCreated, b.Data.TimeCreated)
	})

	return dashboardResponse
}

func (dh *DashboardHandler) createSellStrategyRow(st dto.TradingTaskResponse) dto.TableRow {
	return dto.TableRow{
		Program:   string(st.Program),
		Type:      string(dto.STRATEGY),
		Id:        st.Id,
		WsMessage: st.Message,
		State:     st.State,
		Data:      st,
	}
}

func (dh *DashboardHandler) createBuyStrategyRow(st dto.TradingTaskResponse, sellTasks []dto.ResponseTask) dto.TableRow {
	sellRows := make([]dto.ChildRow, 0, len(sellTasks))

	for _, child := range sellTasks {
		child.SellFee = st.SellFee
		childRow := dto.ChildRow{
			Program:   string(child.Program),
			Type:      string(dto.TASK),
			Id:        child.TaskId,
			WsMessage: child.Message,
			State:     child.State.TaskState,
			Data:      child,
		}

		sellRows = append(sellRows, childRow)
	}

	slices.SortFunc(sellRows, func(a, b dto.ChildRow) int {
		return cmp.Compare(a.Data.TimeCreated, b.Data.TimeCreated)
	})

	return dto.TableRow{
		Program:   st.Program,
		Type:      string(dto.STRATEGY),
		Id:        st.Id,
		WsMessage: st.Message,
		State:     st.State,
		Data:      st,
		Children:  sellRows,
	}
}
