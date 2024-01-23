package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"
	"sublinks/sublinks-federation/internal/lemmy"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gorilla/mux"
)

func (s Server) SetupActivityRoutes() {
	s.Router.HandleFunc("/activities/{action}/{id}", s.getActivityHandler).Methods("GET")
}

func (s Server) getActivityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var content []byte
	switch vars["action"] {
	case "create":
		obj, err := s.GetPostActivityObject(vars["id"])
		if err != nil {
			s.Logger.Error("Error reading object", err)
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

func (s Server) GetPostActivityObject(id string) (*activitypub.Post, error) {
	ctx := context.Background()
	c := lemmy.GetLemmyClient(ctx)
	post, err := c.GetPost(ctx, id)
	if err != nil {
		s.Logger.Error("Error reading post", err)
		return nil, err
	}
	return activitypub.ConvertPostToApub(post), nil
}
