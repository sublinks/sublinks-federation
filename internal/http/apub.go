package http

import (
	"net/http"
)

func (s Server) SetupApubRoutes() {
	s.Router.HandleFunc("/users/{user}/inbox", s.getInboxHandler).Methods("GET")
	s.Router.HandleFunc("/users/{user}/inbox", s.postInboxHandler).Methods("POST")
	s.Router.HandleFunc("/users/{user}/outbox", s.getOutboxHandler).Methods("GET")
	s.Router.HandleFunc("/users/{user}/outbox", s.postOutboxHandler).Methods("POST")
}

func (s Server) getInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (s Server) postInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (s Server) getOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (s Server) postOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}
