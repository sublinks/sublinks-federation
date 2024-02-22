package worker

import (
	"encoding/json"
	"os"
	"strings"
	"sublinks/sublinks-federation/internal/log"
	"sublinks/sublinks-federation/internal/model"
	"sublinks/sublinks-federation/internal/repository"
)

type PostWorker struct {
	log.Logger
	repository.Repository
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
	err = w.Repository.Save(post)
	if err != nil {
		w.Logger.Error("Error saving post: %s", err)
		return err
	}
	return nil
}
