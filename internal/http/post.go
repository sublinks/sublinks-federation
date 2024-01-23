package http

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"
	"sublinks/sublinks-federation/internal/lemmy"
)

func (s Server) SetupPostRoutes() {
	s.Router.HandleFunc("/post/{postId}", s.getPostHandler).Methods("GET")
}

func (s Server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, vars["postId"])
	if err != nil {
		s.Logger.Error("Error reading post", err)
		return
	}
	postLd := activitypub.ConvertPostToApub(post)
	postLd.Context = activitypub.GetContext()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	w.Write(content)
}
