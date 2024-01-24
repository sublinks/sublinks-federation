package http

import (
	"net/http"
)

func (server *Server) SetupApubRoutes() {
	server.Router.HandleFunc("/users/{user}/inbox", server.getInboxHandler).Methods("GET")
	server.Router.HandleFunc("/users/{user}/inbox", server.postInboxHandler).Methods("POST")
	server.Router.HandleFunc("/users/{user}/outbox", server.getOutboxHandler).Methods("GET")
	server.Router.HandleFunc("/users/{user}/outbox", server.postOutboxHandler).Methods("POST")
}

func (server *Server) getInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (server *Server) postInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (server *Server) getOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (server *Server) postOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}
