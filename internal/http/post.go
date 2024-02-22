package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sublinks/sublinks-federation/internal/activitypub"
	"sublinks/sublinks-federation/internal/model"

	"github.com/gorilla/mux"
)

func (server *Server) SetupPostRoutes() {
	server.Router.HandleFunc("/post/{postId}", server.getPostHandler).Methods("GET")
}

func (server *Server) getPostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post := model.Post{UrlStub: vars["postId"]}
	err := server.Database.Find(&post)
	if err != nil {
		server.Logger.Error(fmt.Sprintf("Error reading post: %+v %s", post, err), err)
		return
	}
	postLd := activitypub.ConvertPostToPage(&post)
	postLd.Context = activitypub.GetContext()
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/activity+json")
	content, _ := json.MarshalIndent(postLd, "", "  ")
	_, err = w.Write(content)
	if err != nil {
		server.Logger.Error("Error writing response", err)
	}
}
