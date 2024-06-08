package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"

	"github.com/gorilla/mux"
)

func (server *Server) SetupCommunityRoutes() {
	server.Router.HandleFunc("/c/{community}", server.getCommunityInfoHandler).Methods("GET")
}

func (server *Server) getCommunityInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	server.Logger.Info(fmt.Sprintf("Looking up community %s", vars["community"]))
	community := server.ServiceManager.GetCommunityService().GetById(vars["community"])
	if community == nil {
		server.Logger.Error("Community not found", nil)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	communityLd := activitypub.ConvertActorToGroup(community)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(communityLd, "", "  ")
	_, err := w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}
