package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/gorilla/mux"
)

func (server *Server) SetupActivityRoutes() {
	server.Logger.Debug("Setting up activity routes")
	server.Router.HandleFunc("/activities/{action}/{id}", server.getActivityHandler).Methods("GET")
}

func (server *Server) getActivityHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var content []byte
	switch vars["action"] {
	case "create":
		obj, err := server.GetPostActivityObject(vars["id"])
		if err != nil {
			server.Logger.Error("Error reading object", err)
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
	default:
		server.Logger.Error(fmt.Sprintf("action %s not found", vars["action"]), nil)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	_, err := w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}

func (server *Server) GetPostActivityObject(id string) (*activitypub.Page, error) {
	post := server.ServiceManager.PostService().GetById(id)
	return activitypub.ConvertPostToPage(post), nil
}
