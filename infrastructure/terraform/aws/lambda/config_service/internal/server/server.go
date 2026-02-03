package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"myproject/internal/dynamodb"
	"myproject/internal/handlers"
	"myproject/internal/notification"
)

func NewServer(store dynamodb.ConfigStore, secretManager notification.SecretManager) http.Handler {
	router := mux.NewRouter()
	configHandler := handlers.NewConfigHandler(store)
	alarmHandler := handlers.NewAlarmHandler(store, nil, secretManager)

	router.HandleFunc("/config/{deviceId}", configHandler.GetConfig).Methods("GET")
	router.HandleFunc("/config", configHandler.PutConfig).Methods("POST")
	router.HandleFunc("/check_alarm", alarmHandler.CheckAlarm).Methods("GET", "POST")
	router.HandleFunc("/health", healthCheck).Methods("GET")

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
