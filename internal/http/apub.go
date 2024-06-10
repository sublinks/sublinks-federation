package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (server *Server) SetupApubRoutes() {
	server.Logger.Debug("Setting up Apub routes")
	server.Router.HandleFunc("/{type}/{id}/inbox", server.getInboxHandler).Methods("GET")
	server.Router.HandleFunc("/{type}/{id}/inbox", server.postInboxHandler).Methods("POST")
	server.Router.HandleFunc("/{type}/{id}/outbox", server.getOutboxHandler).Methods("GET")
	server.Router.HandleFunc("/{type}/{id}/outbox", server.postOutboxHandler).Methods("POST")
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
	vars := mux.Vars(r)
	switch vars["type"] {
	case "u":
		break
	case "c":
		break
	default:
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func (server *Server) postOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}
