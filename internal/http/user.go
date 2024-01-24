package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"
	"sublinks/sublinks-federation/internal/lemmy"
)

func (server *Server) SetupUserRoutes() {
	server.Router.HandleFunc("/u/{user}", server.getUserInfoHandler).Methods("GET")
}

func (server *Server) getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	server.Logger.Info(fmt.Sprintf("Looking up user %server", vars["user"]))
	user, err := c.GetUser(ctx, vars["user"])
	if err != nil {
		server.Logger.Error("Error reading user", err)
		return
	}

	userLd := activitypub.ConvertUserToApub(user, r.Host)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(userLd, "", "  ")
	w.Write(content)
}
