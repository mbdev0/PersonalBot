package handler

import (
	"encoding/json"
	"net/http"
	"personal_bot/api/controller"
	"personal_bot/internal/core/settings"
	"personal_bot/pkg/logger"
)

type SettingsHandler struct {
	settingsController *controller.SettingsController
}

func NewSettingsHandler(settingsContoller *controller.SettingsController) *http.ServeMux {
	mux := http.NewServeMux()
	settingsHandler := SettingsHandler{
		settingsController: settingsContoller,
	}
	settingsHandler.registerRoutes(mux)

	return mux
}

func (s *SettingsHandler) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /", s.GetSettings)
	mux.HandleFunc("POST /", s.PostSettings)
	mux.HandleFunc("POST /test_webhook", s.TestDiscordEndpoint)

}

func (s *SettingsHandler) GetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.settingsController.GetSettings(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(settings)
	if err != nil {
		logger.Error(err)
	}

}

func (s *SettingsHandler) PostSettings(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var reqSettings settings.Settings
	err := decoder.Decode(&reqSettings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// TODO:  update the notifier with the values
	err = s.settingsController.PostSettings(r.Context(), reqSettings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}

func (s *SettingsHandler) TestDiscordEndpoint(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var testWebhook settings.TestWebhook
	err := decoder.Decode(&testWebhook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = s.settingsController.TestDiscordEndpoint(testWebhook.DiscordWebhook)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusNoContent)
}
