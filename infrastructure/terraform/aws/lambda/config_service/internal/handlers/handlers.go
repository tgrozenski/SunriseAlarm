package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"myproject/internal/dynamodb"
	"myproject/internal/models"
	"myproject/internal/validation"
)

type ConfigHandler struct {
	store dynamodb.ConfigStore
}

func NewConfigHandler(store dynamodb.ConfigStore) *ConfigHandler {
	return &ConfigHandler{store: store}
}

func (h *ConfigHandler) GetConfig(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	deviceID := vars["deviceId"]
	if deviceID == "" {
		writeError(w, "deviceId is required", http.StatusBadRequest)
		return
	}

	config, err := h.store.GetConfig(r.Context(), deviceID)
	if err != nil {
		if err == dynamodb.ErrConfigNotFound {
			writeError(w, "config not found", http.StatusNotFound)
			return
		}
		writeError(w, "failed to retrieve config", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(config)
}

func (h *ConfigHandler) PutConfig(w http.ResponseWriter, r *http.Request) {
	var config models.UserConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if err := validation.ValidateUserConfig(&config); err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.store.PutConfig(r.Context(), &config); err != nil {
		writeError(w, "failed to save config", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
