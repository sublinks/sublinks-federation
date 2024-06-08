package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"

	"github.com/gorilla/mux"
)

func (server *Server) SetupUserRoutes() {
	server.Router.HandleFunc("/u/{user}", server.getUserInfoHandler).Methods("GET")
}

func (server *Server) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	server.Logger.Info(fmt.Sprintf("Looking up user %s", vars["user"]))
	user := server.ServiceManager.GetUserService().GetById(vars["user"])
	if user == nil {
		server.Logger.Error("User not found", nil)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	userLd := activitypub.ConvertActorToPerson(user)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(userLd, "", "  ")
	_, err := w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}
