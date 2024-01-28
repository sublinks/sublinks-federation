package http

import (
	"context"
	"encoding/json"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"
	"sublinks/sublinks-federation/internal/lemmy"

	"github.com/gorilla/mux"
)

func (server *Server) SetupPostRoutes() {
	server.Router.HandleFunc("/post/{postId}", server.getPostHandler).Methods("GET")
}

func (server *Server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, vars["postId"])
	if err != nil {
		server.Logger.Error("Error reading post", err)
		return
	}
	postLd := activitypub.ConvertPostToApub(post)
	postLd.Context = activitypub.GetContext()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	_, err = w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}
