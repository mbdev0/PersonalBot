package handler

import "net/http"

type TradingHandler struct {
}

func NewTradingHandler() http.Handler {
	mux := http.NewServeMux()
	handler := &TradingHandler{}
	handler.registerRoutes(mux)
	return mux
}

func (th *TradingHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /afk", th.afkMonitoring)
}

func (th *TradingHandler) afkMonitoring(w http.ResponseWriter, r *http.Request) {
	//we need to take in a body of the afk req
	// should have buy strats (e.g. amount, slippage, cu, etc. )

}
