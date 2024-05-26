package worker

import (
	"encoding/json"
	"os"
	"strings"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/service/posts"
)

type PostWorker struct {
	log.Logger
	Service posts.PostService
}

func (w *PostWorker) Process(msg []byte) error {
	post := model.Post{}
	err := json.Unmarshal(msg, &post)
	hostname := os.Getenv("HOSTNAME")
	post.UrlStub = strings.Replace(hostname, post.Id, "", 1)
	if err != nil {
		w.Logger.Error("Error unmarshalling post: %s", err)
		return err
	}
	res := w.Service.Save(&post)
	if !res {
		w.Logger.Error("Error saving post", nil)
		return err
	}
	return nil
}
