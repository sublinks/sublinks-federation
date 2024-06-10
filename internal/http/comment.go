package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"
	"sublinks/sublinks-federation/internal/model"

	"github.com/gorilla/mux"
)

func (server *Server) SetupCommentRoutes() {
	server.Logger.Debug("Setting up comment routes")
	server.Router.HandleFunc("/comment/{commentId}", server.getCommentHandler).Methods("GET")
}

func (server *Server) getCommentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	comment := model.Comment{UrlStub: vars["commentId"]}
	err := server.Database.Find(&comment)
	if err != nil {
		server.Logger.Error(fmt.Sprintf("Error reading comment: %+v %s", comment, err), err)
		return
	}
	commentLd := activitypub.ConvertCommentToNote(&comment)
	commentLd.Context = activitypub.GetContext()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(commentLd, "", "  ")
	_, err = w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}
