package routes

import (
	"context"
	"encoding/json"
	"net/http"
	"sublinks/federation/internal/activitypub"
	"sublinks/federation/internal/lemmy"
	"sublinks/federation/internal/logging/logger"

	"github.com/gorilla/mux"
)

func SetupPostRoutes(r *mux.Router) {
	r.HandleFunc("/post/{postId}", getPostHandler).Methods("GET")
}

func getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, vars["postId"])
	if err != nil {
		logger.GetLogger().Println("Error reading post", err)
		return
	}
	postLd := activitypub.ConvertPostToApub(post)
	postLd.Context = activitypub.GetContext()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	w.Write(content)
}
