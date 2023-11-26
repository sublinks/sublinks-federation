package routes

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupApubRoutes(r *mux.Router) {
	r.HandleFunc("/users/{user}/inbox", getInboxHandler).Methods("GET")
	r.HandleFunc("/users/{user}/inbox", postInboxHandler).Methods("POST")
	r.HandleFunc("/users/{user}/outbox", getOutboxHandler).Methods("GET")
	r.HandleFunc("/users/{user}/outbox", postOutboxHandler).Methods("POST")
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
