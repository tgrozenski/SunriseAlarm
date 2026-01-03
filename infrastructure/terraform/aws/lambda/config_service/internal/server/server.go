package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"myproject/internal/dynamodb"
	"myproject/internal/handlers"
)

func NewServer(store dynamodb.ConfigStore) http.Handler {
	router := mux.NewRouter()
	handler := handlers.NewConfigHandler(store)

	router.HandleFunc("/config/{deviceId}", handler.GetConfig).Methods("GET")
	router.HandleFunc("/config", handler.PutConfig).Methods("POST")
	router.HandleFunc("/health", healthCheck).Methods("GET")

	return router
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
