package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Health struct {
	Status  string                     `json:"status"`
	Message string                     `json:"message,omitempty"`
	Queue   map[string]map[string]bool `json:"queue"`
	DB      bool                       `json:"db"`
}

func (server *Server) SetupInternalRoutes() {
	server.Logger.Info("Setting up internal routes")
	server.Router.HandleFunc("/internal/health", server.getHealthHandler).Methods("GET")

	// Keeping this for the moment to not break existing health check
	server.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods("GET")
}

func (server *Server) getHealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	health := Health{Status: "ok"}
	health.DB = server.Database.Ping()
	health.Queue = server.Queue.Status()
	if !health.DB {
		health.Status = "error"
		health.Message = "Database is not ok"
		w.WriteHeader(http.StatusInternalServerError)
	} else if !server.checkQueueHealth(&health) {
		w.WriteHeader(http.StatusInternalServerError)
	}
	if health.Status == "ok" {
		w.WriteHeader(http.StatusOK)
	}
	content, _ := json.MarshalIndent(health, "", "  ")
	_, err := w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}

func (server *Server) checkQueueHealth(health *Health) bool {
	// No consumers is a problem. However, it's (probably) ok if we have no producers
	// as we might not have sent any messages yet.
	if len(health.Queue["consumers"]) == 0 {
		health.Status = "error"
		health.Message = "No consumers"
		return false
	}
	for queueType, status := range health.Queue {
		for queueName, ok := range status {
			if !ok {
				health.Status = "error"
				health.Message = fmt.Sprintf("%s queue '%s' is not ok", queueType, queueName)
				return false
			}
		}
	}
	return true
}
