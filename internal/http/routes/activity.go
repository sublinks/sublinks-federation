package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"sublinks/federation/internal/activitypub"
	"sublinks/federation/internal/lemmy"

	"fmt"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gorilla/mux"
)

func SetupActivityRoutes(r *mux.Router) {
	r.HandleFunc("/activities/{action}/{id}", getActivityHandler).Methods("GET")
}

func getActivityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var content []byte
	switch vars["action"] {
	case "create":
		obj, err := GetPostActivityObject(vars["id"])
		if err != nil {
			log.Println("Error reading object", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		content, _ = json.MarshalIndent(
			activitypub.NewActivity(
				r.RequestURI,
				cases.Title(language.English).String(vars["action"]),
				obj.AttributedTo,
				obj.To,
				obj.Cc,
				obj.Audience,
				obj,
			), "", "  ")

		break
	default:
		error.Error(fmt.Errorf("action %s not found", vars["action"]))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	w.Write(content)
}

func GetPostActivityObject(id string) (*activitypub.Post, error) {
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, id)
	if err != nil {
		log.Println("Error reading post", err)
		return nil, err
	}
	return activitypub.ConvertPostToApub(post), nil
}
