package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"participating-online/sublinks-federation/internal/activitypub"
	"participating-online/sublinks-federation/internal/lemmy"

	"github.com/gorilla/mux"
)

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, vars["postId"])
	if err != nil {
		log.Println("Error reading post", err)
		return
	}
	postLd := activitypub.ConvertPostToApub(post)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	w.Write(content)
}

func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	log.Println(fmt.Sprintf("Looking up user %s", vars["user"]))
	user, err := c.GetUser(ctx, vars["user"])
	if err != nil {
		log.Println("Error reading user", err)
		return
	}

	userLd := activitypub.ConvertUserToApub(user)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(userLd, "", "  ")
	w.Write(content)
}

func GetInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func PostInboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func GetOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}

func PostOutboxHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
}
