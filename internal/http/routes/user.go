package routes

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

func SetupUserRoutes(r *mux.Router) {
	r.HandleFunc("/u/{user}", getUserInfoHandler).Methods("GET")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	log.Println(fmt.Sprintf("Looking up user %s", vars["user"]))
	user, err := c.GetUser(ctx, vars["user"])
	if err != nil {
		log.Println("Error reading user", err)
		return
	}

	userLd := activitypub.ConvertUserToApub(user, r.Host)
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(userLd, "", "  ")
	w.Write(content)
}
