package http

import (
	"encoding/json"
	"net/http"
)

func (s Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Request("", r)
		next.ServeHTTP(w, r)
	})
}

type RequestError struct {
	Msg string `json:"message"`
}

func (s Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.Logger.Request("404 Not Found", r)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.Marshal(RequestError{Msg: "not found"})
	w.Write(content)
}

func (s Server) notAllowedMethod(w http.ResponseWriter, r *http.Request) {
	s.Logger.Request("405 Method Not Allowed", r)
	w.WriteHeader(http.StatusNotFound)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.Marshal(RequestError{Msg: "method not allowed"})
	w.Write(content)
}
