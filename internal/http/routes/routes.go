package routes

import (
	"encoding/json"
	"net/http"
	"sublinks/federation/internal/logging"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	SetupUserRoutes(r)
	SetupPostRoutes(r)
	SetupApubRoutes(r)
	SetupActivityRoutes(r)
	r.NotFoundHandler = http.HandlerFunc(notFound)
	r.MethodNotAllowedHandler = http.HandlerFunc(notAllowedMethod)
	r.Use(loggingMiddleware)
	return r
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logging.Request("", r)
		next.ServeHTTP(w, r)
	})
}

type RequestError struct {
	Msg string `json:"message"`
}

func notFound(w http.ResponseWriter, r *http.Request) {
	logging.Request("404 Not Found", r)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.Marshal(RequestError{Msg: "not found"})
	w.Write(content)
}

func notAllowedMethod(w http.ResponseWriter, r *http.Request) {
	logging.Request("405 Method Not Allowed", r)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.Marshal(RequestError{Msg: "method not allowed"})
	w.Write(content)
}
