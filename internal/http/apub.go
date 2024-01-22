package http

import (
	"net/http"
)

func (s Server) SetupApubRoutes() {
	s.HandleFunc("/users/{user}/inbox", getInboxHandler).Methods("GET")
	s.HandleFunc("/users/{user}/inbox", postInboxHandler).Methods("POST")
	s.HandleFunc("/users/{user}/outbox", getOutboxHandler).Methods("GET")
	s.HandleFunc("/users/{user}/outbox", postOutboxHandler).Methods("POST")
}

func getInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func postInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func getOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func postOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}
